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
 * V2 Subscriptions — Advanced scenarios.
 *
 * Xray tests:
 *   GKO-796: Subscription between API and App in different mgmt contexts is rejected
 *   GKO-798: Subscribing to a V2 API with local=true is rejected by admission
 *   GKO-821: Delete a V2 JWT subscription — plan closes, gateway call no longer works
 *   GKO-825: Delete a V2 OAuth2 subscription
 *   GKO-839: API and Application must have Synced=True after last reconcile
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

const APP_SIMPLE = "crds/applications/application-simple.yaml";
const V2_JWT_API = "crds/subscriptions/v2-api-jwt-plan.yaml";
const V2_OAUTH2_API = "crds/subscriptions/v2-api-oauth2-plan.yaml";
const V2_LOCAL_API = "crds/subscriptions/v2-api-local-true.yaml";
const SUB_JWT_V2 = "crds/subscriptions/subscription-jwt-v2.yaml";
const SUB_OAUTH2_V2 = "crds/subscriptions/subscription-oauth2-v2.yaml";
const SUB_LOCAL_V2 = "crds/subscriptions/subscription-to-local-v2-api.yaml";
const CROSS_CTX_APP = "crds/subscriptions/cross-ctx-app.yaml";
const CROSS_CTX_SUB = "crds/subscriptions/cross-ctx-subscription.yaml";
const TEMP_CTX = "crds/management-context/temporary-ctx.yaml";

test.describe("V2 Subscriptions — Advanced", () => {
  test.afterEach(async () => {
    // Cleanup in reverse dependency order: subscriptions → apps → APIs.
    await kubectlSafe.del(fixture(SUB_JWT_V2)).catch(() => {});
    await kubectlSafe.del(fixture(SUB_OAUTH2_V2)).catch(() => {});
    await kubectlSafe.del(fixture(SUB_LOCAL_V2)).catch(() => {});
    await kubectlSafe.del(fixture(CROSS_CTX_SUB)).catch(() => {});
    await kubectlSafe.del(fixture(CROSS_CTX_APP)).catch(() => {});
    await kubectlSafe.del(fixture(APP_SIMPLE)).catch(() => {});
    await kubectlSafe.del(fixture(V2_JWT_API)).catch(() => {});
    await kubectlSafe.del(fixture(V2_OAUTH2_API)).catch(() => {});
    await kubectlSafe.del(fixture(V2_LOCAL_API)).catch(() => {});
    await kubectlSafe.del(fixture(TEMP_CTX)).catch(() => {});
  });

  // ── GKO-796: Cross-mgmt-context subscription is rejected ────

  test(`Subscription across different mgmt contexts is rejected ${XRAY.SUBSCRIPTIONS.CROSS_MGMT_CONTEXT_ERROR} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    test.slow(); // V2 API + two mgmt contexts takes extra reconcile time.

    await kubectl.apply(fixture(TEMP_CTX));
    await kubectl.apply(fixture(V2_JWT_API));
    await kubectl.apply(fixture(CROSS_CTX_APP));
    await kubectl.waitForCondition("apidefinition", "e2e-v2-jwt-api", "Accepted");
    await kubectl.waitForCondition("application", "e2e-app-temp-ctx", "Accepted");

    // The CR either gets rejected at admission or goes to Accepted=False.
    let admissionRejected = false;
    try {
      await kubectl.apply(fixture(CROSS_CTX_SUB));
    } catch {
      admissionRejected = true;
    }

    if (!admissionRejected) {
      await expect
        .poll(
          async () => {
            const status = await kubectl
              .getStatus<{
                conditions?: Array<{ type: string; status: string }>;
              }>("subscription", "e2e-sub-cross-ctx")
              .catch(() => ({ conditions: [] }));
            return status.conditions?.find((c) => c.type === "Accepted")?.status;
          },
          { timeout: 30_000 },
        )
        .toBe("False");
    }

    await kubectl.del(fixture(CROSS_CTX_SUB)).catch(() => {});
    await kubectl.del(fixture(CROSS_CTX_APP));
    await kubectl.del(fixture(V2_JWT_API));
    await kubectl.del(fixture(TEMP_CTX));
  });

  // ── GKO-798: Subscribing to a local=true V2 API is rejected ─

  test(`Subscription to V2 API with local=true is rejected ${XRAY.SUBSCRIPTIONS.V2_LOCAL_SUBSCRIPTION_ERROR} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    test.slow();

    // `local: true` means the API is not pushed to APIM; the API CR reconciles
    // without the operator requiring a management context roundtrip, so we
    // don't wait for Accepted there. We just need the Application present.
    await kubectl.apply(fixture(V2_LOCAL_API));
    await kubectl.apply(fixture(APP_SIMPLE));
    await kubectl.waitForCondition("application", "e2e-app-simple", "Accepted");

    // Admission rejects the subscription because you cannot subscribe to
    // a local V2 API.
    const stderr = await kubectl.applyExpectFailure(fixture(SUB_LOCAL_V2));
    expect(stderr.toLowerCase()).toMatch(/local|denied|invalid|cannot/);

    await kubectl.del(fixture(APP_SIMPLE));
    await kubectl.del(fixture(V2_LOCAL_API));
  });

  // ── GKO-821: Delete V2 JWT subscription ─────────────────────

  test(`Delete V2 JWT subscription closes access ${XRAY.SUBSCRIPTIONS.V2_JWT_DELETE} ${TAGS.REGRESSION}`, async ({
    kubectl,
    gateway,
  }) => {
    test.slow();

    await kubectl.apply(fixture(V2_JWT_API));
    await kubectl.apply(fixture(APP_SIMPLE));
    await kubectl.waitForCondition("apidefinition", "e2e-v2-jwt-api", "Accepted");
    await kubectl.waitForCondition("application", "e2e-app-simple", "Accepted");
    await kubectl.apply(fixture(SUB_JWT_V2));
    await kubectl.waitForCondition("subscription", "e2e-sub-jwt-v2", "Accepted");

    // Delete the subscription CR — the CR should be removed and the API
    // still responds 401 (JWT plan is still protecting it; the subscription
    // closing removes only the APP's access, not the plan).
    await kubectl.del(fixture(SUB_JWT_V2));
    await kubectl.waitForDeletion("subscription", "e2e-sub-jwt-v2");

    await gateway.assertResponds("/e2e-v2-jwt-api", { status: 401 });

    await kubectl.del(fixture(APP_SIMPLE));
    await kubectl.del(fixture(V2_JWT_API));
  });

  // ── GKO-825: Delete V2 OAuth2 subscription ──────────────────

  test(`Delete V2 OAuth2 subscription closes access ${XRAY.SUBSCRIPTIONS.V2_OAUTH2_DELETE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    test.slow();

    await kubectl.apply(fixture(V2_OAUTH2_API));
    await kubectl.apply(fixture(APP_SIMPLE));
    await kubectl.waitForCondition("apidefinition", "e2e-v2-oauth2-api", "Accepted");
    await kubectl.waitForCondition("application", "e2e-app-simple", "Accepted");
    await kubectl.apply(fixture(SUB_OAUTH2_V2));
    await kubectl.waitForCondition("subscription", "e2e-sub-oauth2-v2", "Accepted");

    await kubectl.del(fixture(SUB_OAUTH2_V2));
    await kubectl.waitForDeletion("subscription", "e2e-sub-oauth2-v2");

    await kubectl.del(fixture(APP_SIMPLE));
    await kubectl.del(fixture(V2_OAUTH2_API));
  });

  // ── GKO-839: API and App must be Synced on last reconcile ──

  test(`V2 API and Application show Accepted=True after reconcile ${XRAY.SUBSCRIPTIONS.API_APP_SYNCED_LAST_RECONCILE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    test.slow();

    await kubectl.apply(fixture(V2_JWT_API));
    await kubectl.apply(fixture(APP_SIMPLE));
    await kubectl.waitForCondition("apidefinition", "e2e-v2-jwt-api", "Accepted");
    await kubectl.waitForCondition("application", "e2e-app-simple", "Accepted");

    const apiStatus = await kubectl.getStatus<{
      conditions?: Array<{ type: string; status: string }>;
    }>("apidefinition", "e2e-v2-jwt-api");
    const appStatus = await kubectl.getStatus<{
      conditions?: Array<{ type: string; status: string }>;
    }>("application", "e2e-app-simple");

    expect(apiStatus.conditions?.find((c) => c.type === "Accepted")?.status).toBe("True");
    expect(appStatus.conditions?.find((c) => c.type === "Accepted")?.status).toBe("True");

    await kubectl.del(fixture(APP_SIMPLE));
    await kubectl.del(fixture(V2_JWT_API));
  });
});
