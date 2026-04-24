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
 * V4 API-Key plan subscriptions.
 *
 * Covers the CRD-driven flow where a Subscription targets a plan of security
 * type API_KEY (enabled by GKO-2547 / PR #1642, merged 2026-04-22), and guards
 * against APIM-13686 where two API keys were generated for a single
 * subscription.
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

const API_APIKEY = "crds/api-v4-definitions/v4-proxy-api-apikey-plan.yaml";
const API_TWO_PLANS = "crds/api-v4-definitions/v4-proxy-api-two-plans.yaml";
const APP = "crds/applications/application-simple.yaml";
const SUB_APIKEY = "crds/subscriptions/subscription-apikey-v4.yaml";
const SUB_APIKEY_TWO_PLANS = "crds/subscriptions/subscription-apikey-v4-two-plans.yaml";

const API_NAME = "e2e-v4-apikey-plan";
const API_PATH = "/e2e-v4-apikey-plan";
const TWO_PLANS_API_NAME = "e2e-v4-two-plans";
const TWO_PLANS_API_PATH = "/e2e-v4-two-plans";
const APP_NAME = "e2e-app-simple";
const SUB_NAME = "e2e-sub-apikey-v4";
const SUB_TWO_PLANS_NAME = "e2e-sub-apikey-v4-two-plans";

test.describe("V4 API-Key Plan Subscriptions", () => {
  test.afterEach(async () => {
    // Reverse dependency order: subscriptions → applications → APIs
    await kubectlSafe.del(fixture(SUB_APIKEY)).catch(() => {});
    await kubectlSafe.del(fixture(SUB_APIKEY_TWO_PLANS)).catch(() => {});
    await kubectlSafe.del(fixture(APP)).catch(() => {});
    await kubectlSafe.del(fixture(API_APIKEY)).catch(() => {});
    await kubectlSafe.del(fixture(API_TWO_PLANS)).catch(() => {});
  });

  // ── Single API key generated on subscription (APIM-13686 regression) ──

  test(`Single API key generated on api-key plan subscription ${XRAY.SUBSCRIPTIONS.V4_APIKEY_SINGLE_KEY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    await kubectl.apply(fixture(API_APIKEY));
    await kubectl.apply(fixture(APP));
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");
    await kubectl.apply(fixture(SUB_APIKEY));
    await kubectl.waitForCondition("subscription", SUB_NAME, "Accepted");

    const { id: subId } = await kubectl.getStatus<{ id: string }>("subscription", SUB_NAME);
    const { id: apiId } = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);

    await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);

    // Re-poll after a short delay to rule out an asynchronous second insert
    // (APIM-13686 produced two keys ~1s apart).
    await new Promise((r) => setTimeout(r, 2_000));
    const result = await mapi.listSubscriptionApiKeys(apiId, subId);
    expect(result.pagination?.totalCount ?? result.data.length).toBe(1);
  });

  // ── Gateway accepts the generated key, rejects without it ──

  test(`Gateway accepts the generated api key and rejects without it ${XRAY.SUBSCRIPTIONS.V4_APIKEY_GATEWAY_CALL} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
    gateway,
  }) => {
    // Budget: gateway.assertResponds polls for up to 30s waiting for the
    // api-key plan to propagate from APIM to the gateway. Combined with the
    // kubectl waits that precede it, the default 30s Playwright test timeout
    // has no headroom and the test is killed mid-poll on slower CI runs.
    // Locally the gateway sync is sub-second so 30s is plenty; in CI it can
    // take 15–25s, pushing total test time well past 30s.
    test.setTimeout(60_000);

    await kubectl.apply(fixture(API_APIKEY));
    await kubectl.apply(fixture(APP));
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");
    await kubectl.apply(fixture(SUB_APIKEY));
    await kubectl.waitForCondition("subscription", SUB_NAME, "Accepted");

    const { id: subId } = await kubectl.getStatus<{ id: string }>("subscription", SUB_NAME);
    const { id: apiId } = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);

    await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
    const { data } = await mapi.listSubscriptionApiKeys(apiId, subId);
    const apiKey = data[0]?.key;
    expect(apiKey).toBeTruthy();

    await test.step("Gateway rejects call without api key header", async () => {
      await gateway.assertResponds(API_PATH, { status: 401 });
    });

    await test.step("Gateway accepts call with valid api key header", async () => {
      await gateway.assertResponds(API_PATH, {
        status: 200,
        headers: { "X-Gravitee-Api-Key": apiKey },
      });
    });
  });

  // ── Webhook accepts API_KEY subscriptions (GKO-2547 regression) ──

  test(`Admission webhook accepts api-key plan subscriptions ${XRAY.SUBSCRIPTIONS.V4_APIKEY_WEBHOOK_ACCEPTED} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    await kubectl.apply(fixture(API_APIKEY));
    await kubectl.apply(fixture(APP));
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");

    // Before GKO-2547, the admission webhook rejected API_KEY plan subscriptions
    // with "security type is not one of [JWT,OAUTH2,MTLS]". `kubectl apply` must
    // not fail and the Subscription must reach Accepted.
    await kubectl.apply(fixture(SUB_APIKEY));
    await kubectl.waitForCondition("subscription", SUB_NAME, "Accepted");

    const status = await kubectl.getStatus<{ id: string }>("subscription", SUB_NAME);
    expect(status.id).toBeTruthy();
  });

  // ── Key is revoked when the Subscription CRD is deleted ──

  test(`Api key is revoked when the subscription is deleted ${XRAY.SUBSCRIPTIONS.V4_APIKEY_KEY_REVOKED_ON_DELETE} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
    gateway,
  }) => {
    // Budget: see the comment on the GKO-2826 test above. This test has two
    // gateway polls (one pre-delete, one post-delete), so the headroom need
    // is even larger.
    test.setTimeout(60_000);

    await kubectl.apply(fixture(API_APIKEY));
    await kubectl.apply(fixture(APP));
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");
    await kubectl.apply(fixture(SUB_APIKEY));
    await kubectl.waitForCondition("subscription", SUB_NAME, "Accepted");

    const { id: subId } = await kubectl.getStatus<{ id: string }>("subscription", SUB_NAME);
    const { id: apiId } = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);

    await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
    const { data } = await mapi.listSubscriptionApiKeys(apiId, subId);
    const apiKey = data[0]?.key;
    expect(apiKey).toBeTruthy();

    // Confirm the key works before deleting the subscription.
    await gateway.assertResponds(API_PATH, {
      status: 200,
      headers: { "X-Gravitee-Api-Key": apiKey },
    });

    await kubectl.del(fixture(SUB_APIKEY));
    await kubectl.waitForDeletion("subscription", SUB_NAME);

    // After the subscription is gone, the previously valid key must stop working.
    // Status code depends on gateway policy (401, 403, or 404) — assertNotResponds
    // is the right fit here.
    await gateway.assertNotResponds(API_PATH, {
      notStatus: 200,
      headers: { "X-Gravitee-Api-Key": apiKey },
    });
  });

  // ── Multi-plan API: api-key plan coexists with keyless plan ──

  test(`Api-key plan coexists with keyless plan on the same API ${XRAY.SUBSCRIPTIONS.V4_APIKEY_MIXED_WITH_KEYLESS} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
    gateway,
  }) => {
    await kubectl.apply(fixture(API_TWO_PLANS));
    await kubectl.apply(fixture(APP));
    await kubectl.waitForCondition("apiv4definition", TWO_PLANS_API_NAME, "Accepted");
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");
    await kubectl.apply(fixture(SUB_APIKEY_TWO_PLANS));
    await kubectl.waitForCondition("subscription", SUB_TWO_PLANS_NAME, "Accepted");

    const { id: subId } = await kubectl.getStatus<{ id: string }>("subscription", SUB_TWO_PLANS_NAME);
    const { id: apiId } = await kubectl.getStatus<{ id: string }>("apiv4definition", TWO_PLANS_API_NAME);

    await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
    const { data } = await mapi.listSubscriptionApiKeys(apiId, subId);
    const apiKey = data[0]?.key;
    expect(apiKey).toBeTruthy();

    await test.step("Keyless plan accepts traffic with no api key header", async () => {
      await gateway.assertResponds(TWO_PLANS_API_PATH, { status: 200 });
    });

    await test.step("Api-key plan accepts traffic with a valid api key header", async () => {
      await gateway.assertResponds(TWO_PLANS_API_PATH, {
        status: 200,
        headers: { "X-Gravitee-Api-Key": apiKey },
      });
    });

    // Discriminator: a present but invalid api-key header is handled by the
    // api-key plan (not keyless), so the gateway must reject. Without this
    // step, the valid-key assertion above would pass even if api-key plan
    // resolution regressed and keyless silently handled every header-bearing
    // request.
    await test.step("Api-key plan rejects traffic with an invalid api key header", async () => {
      await gateway.assertNotResponds(TWO_PLANS_API_PATH, {
        notStatus: 200,
        headers: { "X-Gravitee-Api-Key": "bogus-invalid-key" },
      });
    });
  });
});
