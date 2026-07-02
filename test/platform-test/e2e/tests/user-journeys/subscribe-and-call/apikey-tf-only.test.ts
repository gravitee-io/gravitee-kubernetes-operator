/**
 * Copyright (C) 2015 The Gravitee team (http://gravitee.io)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/**
 * Terraform-only api-key behaviour: declarative-diff drift, Sensitive-value
 * redaction in plan output, and server-side length rejection at apply time.
 * These are provider-level concerns with no GKO analog (k8s reconciliation is
 * push-based, not declarative-diff-based). The no-drift case uses the handle's
 * `checks` surface; the pre-apply plan and expected-failure cases use the raw
 * terraform helper because they do not fit the "provision succeeds" handle flow.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test, expect } from "../../../setup.js";
import { XRAY, TAGS, PROVISIONER } from "../../../helpers/tags.js";
import * as terraform from "../../../helpers/terraform.js";
import { tfScenario } from "../../../helpers/provisioner-env.js";
import { isTerraform, type Provisioned } from "../../../../src/provisioners/index.js";
import { tfApiKeyVars, uniqueKey, RUN_ID, type ApiKeyParams } from "./params.js";

/** Co-located terraform fixture for the custom-key subscription. */
const apikeyCustomFixture = path.join(path.dirname(fileURLToPath(import.meta.url)), "terraform/apikey-custom");

test.describe(`Terraform-only: api-key provider behaviour ${PROVISIONER.TERRAFORM} @since-4.12`, () => {
  let handle: Provisioned<ApiKeyParams> | undefined;
  let ws: terraform.TfWorkspace | undefined;

  test.afterEach(async () => {
    test.setTimeout(terraform.TF_WORKSPACE_TIMEOUT_MS);
    if (handle) {
      await handle.destroy().catch(() => {});
      handle = undefined;
    }
    if (ws) {
      await terraform.destroyWorkspace(ws).catch(() => {});
      ws = undefined;
    }
  });

  // ── plan is clean immediately after apply (no declarative-diff drift) ──
  test(`terraform plan reports no drift immediately after apply ${XRAY.TERRAFORM.APIKEY_PLAN_NO_DRIFT} ${TAGS.REGRESSION}`, async () => {
    test.setTimeout(terraform.TF_WORKSPACE_TIMEOUT_MS);
    const KEY_A = uniqueKey("tf-no-drift-A");
    const KEY_B = uniqueKey("tf-no-drift-B");
    const expireB = new Date(Date.now() + 30 * 60 * 1_000).toISOString();

    const provisioner = await tfScenario<ApiKeyParams>({
      fixture: apikeyCustomFixture,
      toVars: tfApiKeyVars("no-drift"),
    })();
    handle = await provisioner.provision({ keys: [{ key: KEY_A }, { key: KEY_B, expireAt: expireB }] });

    expect(isTerraform(handle.checks)).toBe(true);
    if (isTerraform(handle.checks)) {
      await handle.checks.assertNoDrift();
      await handle.checks.assertReapplyNoop();
    }
  });

  // ── key values are Sensitive (redacted) in plan output ──
  test(`Custom api-key values are redacted as sensitive in terraform plan output ${XRAY.TERRAFORM.APIKEY_SENSITIVE_IN_PLAN} ${TAGS.REGRESSION}`, async () => {
    test.setTimeout(terraform.TF_WORKSPACE_TIMEOUT_MS);
    const SECRET_KEY = uniqueKey("tf-sensitive-leak-detect");

    ws = await terraform.initWorkspace(apikeyCustomFixture);
    await terraform.writeVars(ws, tfApiKeyVars("sensitive")({ keys: [{ key: SECRET_KEY }] }));

    const { stdout: preApply } = await terraform.plan(ws);
    expect(preApply.includes(SECRET_KEY), "pre-apply plan leaked the api-key value").toBe(false);
    expect(preApply.toLowerCase()).toContain("sensitive");

    await terraform.apply(ws);
    const { stdout: postApply } = await terraform.plan(ws);
    expect(postApply.includes(SECRET_KEY), "post-apply plan leaked the api-key value").toBe(false);
  });

  // ── server-side length rejection at apply time ──
  test(`Out-of-bounds api-keys are rejected at apply time ${XRAY.TERRAFORM.APIKEY_LENGTH_REJECTED} ${TAGS.REGRESSION}`, async () => {
    test.setTimeout(terraform.TF_WORKSPACE_TIMEOUT_MS);
    ws = await terraform.initWorkspace(apikeyCustomFixture);

    await test.step("31-char key rejected (below the 32 minimum)", async () => {
      const shortKey = `tf-too-short-${RUN_ID}`.padEnd(31, "0").slice(0, 31);
      expect(shortKey).toHaveLength(31);
      await terraform.writeVars(ws!, tfApiKeyVars("too-short")({ keys: [{ key: shortKey }] }));
      const out = await terraform.applyExpectFailure(ws!);
      expect(out.toLowerCase()).toMatch(/minlength|too short|length|invalid|400|bad request|validation/);
    });

    await test.step("257-char key rejected (above the 256 maximum)", async () => {
      const longPrefix = `tf-too-long-${RUN_ID}-`;
      const longKey = longPrefix + "z".repeat(257 - longPrefix.length);
      expect(longKey).toHaveLength(257);
      await terraform.writeVars(ws!, tfApiKeyVars("too-long")({ keys: [{ key: longKey }] }));
      const out = await terraform.applyExpectFailure(ws!);
      expect(out.toLowerCase()).toMatch(/maxlength|too long|length|invalid|400|bad request|validation/);
    });
  });
});
