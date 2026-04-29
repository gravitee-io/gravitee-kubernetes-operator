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
 * Operator restart recovery.
 *
 * Xray tests:
 *   GKO-1451: Restarting the GKO operator pod mid-flight must not break
 *             existing CRs. After the rollout completes, the operator
 *             re-lists CRs and the existing API remains Accepted=True,
 *             and a follow-up CR update reconciles cleanly.
 */

import { readFileSync } from "node:fs";
import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

// Use a dedicated fixture (separate API name) so the slow post-restart
// cleanup window cannot collide with sibling tests that also use
// v4-proxy-api-reconcile.yaml (notably reconciliation.test.ts and
// auditability.test.ts).
const ORIGINAL = "crds/api-v4-definitions/v4-proxy-api-restart.yaml";
const API_NAME = "e2e-v4-restart";
const OPERATOR_DEPLOY = "gko-controller-manager";

test.describe("Operator restart — recovery", () => {
  test.afterEach(async () => {
    await kubectlSafe.del(fixture(ORIGINAL)).catch(() => {});
  });

  test(`Operator restart does not break existing CRs ${XRAY.DEPLOYMENT_RECONCILIATION.OPERATOR_RESTART_RECOVERY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    test.slow();

    await test.step("Deploy CR", async () => {
      await kubectl.apply(fixture(ORIGINAL));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const apiIdBefore = (
      await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME)
    ).id;

    await test.step("Restart operator and wait for rollout", async () => {
      await kubectl.rolloutRestart("deployment", OPERATOR_DEPLOY);
      await kubectl.waitForRollout("deployment", OPERATOR_DEPLOY, 180);
    });

    await test.step("Existing CR remains Accepted=True after restart", async () => {
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("APIM still sees the same API after restart", async () => {
      const api = await mapi.fetchApi(apiIdBefore);
      expect(api.name).toBe(API_NAME);
    });

    // The previous two assertions only confirm pre-restart state still
    // looks healthy; they don't prove the new operator process has
    // re-listed CRs or can reconcile *future* changes. Apply a spec
    // change and assert APIM picks it up — that's what proves the
    // reconcile loop is live again.
    await test.step("Post-restart spec change reconciles into APIM", async () => {
      const original = readFileSync(fixture(ORIGINAL), "utf8");
      const updated = original
        .replace(/version: "1\.0\.0"/, 'version: "1.0.1"')
        .replace(
          /description: "E2E test: Operator restart recovery"/,
          'description: "E2E test: Operator restart recovery (post-rollout update)"',
        );
      expect(updated, "version/description bump must take effect").not.toBe(original);

      // The validating webhook is part of the operator that just
      // restarted; the readiness probe can flip to Ready a few
      // seconds before the webhook server actually accepts TLS
      // connections. Retry briefly to absorb that race.
      const applyDeadline = Date.now() + 60_000;
      let lastErr: unknown;
      while (Date.now() < applyDeadline) {
        try {
          await kubectl.applyString(updated);
          lastErr = undefined;
          break;
        } catch (err) {
          lastErr = err;
          await new Promise((r) => setTimeout(r, 2_000));
        }
      }
      if (lastErr) throw lastErr;
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

      await mapi.waitForApiMatches(
        apiIdBefore,
        {
          apiVersion: "1.0.1",
          description: "E2E test: Operator restart recovery (post-rollout update)",
        },
        { description: "post-restart update reconciles into APIM" },
      );
    });
    // Cleanup handled by afterEach — the operator needs a few seconds after
    // the webhook endpoints recover before accepting deletes.
  });
});
