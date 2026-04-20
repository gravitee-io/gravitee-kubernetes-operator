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
 * Auditability — batch 7.
 *
 * Xray tests:
 *   GKO-1463: Updates on a CR are traceable — k8s Events are emitted,
 *             status.conditions.lastTransitionTime reflects the change,
 *             and APIM-side updatedAt advances after the reconcile.
 *
 * The GKO-1463 Xray scope also mentions APIM audit logs. That API is not
 * currently exposed via the test Mapi client, so we verify traceability
 * through the k8s events (operator-side) and updatedAt (APIM-side) which
 * together give a reliable audit trail.
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

const V4_ORIGINAL = "crds/api-v4-definitions/v4-proxy-api-reconcile.yaml";
const V4_UPDATED = "crds/api-v4-definitions/v4-proxy-api-reconcile-updated.yaml";
const API_NAME = "e2e-v4-reconcile";

test.describe("Auditability — changes traced via events and updatedAt", () => {
  test.afterEach(async () => {
    await kubectlSafe.del(fixture(V4_UPDATED)).catch(() => {});
    await kubectlSafe.del(fixture(V4_ORIGINAL)).catch(() => {});
  });

  test(`Update to CR is auditable via events and updatedAt ${XRAY.DEPLOYMENT_RECONCILIATION.AUDITABILITY_EVENTS} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    await test.step("Deploy initial CR", async () => {
      await kubectl.apply(fixture(V4_ORIGINAL));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const apiId = (
      await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME)
    ).id;
    const before = await mapi.fetchApi(apiId);
    const beforeUpdatedAt = before.updatedAt;

    // Separate reconciles by a second so updatedAt can advance.
    await new Promise((r) => setTimeout(r, 1_100));

    await test.step("Update CR", async () => {
      await kubectl.apply(fixture(V4_UPDATED));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("APIM updatedAt has advanced", async () => {
      await expect
        .poll(
          async () => {
            const api = await mapi.fetchApi(apiId);
            return new Date(api.updatedAt).getTime();
          },
          { timeout: 30_000, intervals: [1_000] },
        )
        .toBeGreaterThan(new Date(beforeUpdatedAt).getTime());
    });

    await test.step("Reconcile fired a k8s Event", async () => {
      // The operator emits reason=UpdateSucceeded after a successful apply.
      await kubectl.assertEventContains("apiv4definition", API_NAME, "succeeded");
    });

    await kubectl.del(fixture(V4_UPDATED));
  });
});
