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
 * Terraform subscription errors & application lifecycle.
 *
 * Xray tests:
 *   GKO-1380: invalid subscription configuration in TF — apply fails with a
 *             clear message naming the offending api/plan reference.
 *
 * GKO-1383 (TF-based deletion of an application archives it in APIM) moved to
 * the shared cross-provisioner journey
 * tests/scenarios/register-and-retire-application (its Terraform arm).
 *
 * GKO-1381 (role-specific access via TF) is dropped — the harness lacks a
 * non-admin user; see tags.ts for details.
 *
 * Preconditions:
 *   - APIM and Gateway are running
 *   - terraform CLI is installed
 */

import { test, expect } from "../../setup.js";
import { XRAY, TAGS } from "../../helpers/tags.js";
import * as terraform from "../../helpers/terraform.js";
import type { TfWorkspace } from "../../helpers/terraform.js";

test.describe("Terraform — subscription errors & app delete", () => {
  // ── GKO-1380: invalid subscription is rejected by apply ──────

  test(`Invalid subscription config produces a clear apply error ${XRAY.TERRAFORM.INVALID_SUBSCRIPTION_CONFIG} ${TAGS.REGRESSION}`, async () => {
    // init + apply + destroy are sequential terraform invocations, each capped
    // at terraform.TF_TIMEOUT_MS. The test timeout must exceed their combined
    // ceiling so terraform's own timeout fires first instead of Playwright
    // orphaning a running terraform process (leaking the .tfstate lock).
    test.setTimeout(terraform.TF_WORKSPACE_TIMEOUT_MS);
    let ws: TfWorkspace | null = null;
    try {
      ws = await terraform.initWorkspace("subscriptions/invalid-sub");
      let stderr = "";
      let succeeded = false;
      try {
        await terraform.apply(ws);
        succeeded = true;
      } catch (err: unknown) {
        const e = err as { stderr?: string; stdout?: string; message?: string };
        stderr = `${e.stderr ?? ""}\n${e.stdout ?? ""}\n${e.message ?? ""}`;
      }
      expect(succeeded, "expected `terraform apply` to fail for invalid subscription").toBe(false);
      // Scoped to terms tied to the missing api_hrid/plan_hrid lookup. A
      // generic TF error (e.g. provider auth failure) would still include
      // "error" but not these — dropping the catch-all keeps the assertion
      // meaningful.
      expect(stderr.toLowerCase()).toMatch(/api_hrid|plan_hrid|does-not-exist|not found|404/);
    } finally {
      if (ws) await terraform.destroyWorkspace(ws);
    }
  });

  // GKO-1383 (Terraform destroy archives the application) is now covered by the
  // Terraform arm of tests/scenarios/register-and-retire-application.
});
