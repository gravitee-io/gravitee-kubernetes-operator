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

// Subscriptions that the new tests build inline with kubectl.applyString.
// Listed here so afterEach can delete them tolerantly even on test timeout.
const INLINE_SUB_NAMES = [
  "e2e-sub-apikey-custom",
  "e2e-sub-apikey-boundary",
  "e2e-sub-apikey-expire",
  "e2e-sub-apikey-secret",
  "e2e-sub-apikey-rotation-instant",
  "e2e-sub-apikey-rotation-gradual",
  "e2e-sub-apikey-reactivation",
  "e2e-sub-apikey-idempotent",
  "e2e-sub-apikey-staggered",
];

// APIM enforces uniqueness on api-key *values* per API across active and
// revoked states (and APIM MongoDB persists across cluster lifecycle on the
// local setup), so re-running a test with hardcoded key values yields
// "API Key already exists" 400s on the second run. A per-process suffix
// makes every test run pick fresh values without polluting earlier ones.
//
// Date.now() alone collides if two processes start in the same millisecond
// (e.g. CI re-runs against a shared cluster, or a future workers>1 config),
// so we mix in a 4-char random suffix.
const RUN_ID = `${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 6)}`;

/** Generate a unique api-key value of at least 32 chars. */
function uniqueKey(prefix: string): string {
  return `${prefix}-${RUN_ID}`.padEnd(32, "0");
}

const SECRET_NAME = `e2e-apikey-secret-${RUN_ID}`;

/**
 * Build a Subscription manifest targeting the api-key plan API + simple app.
 * Used by tests that need to mutate spec.apiKeys across calls (rotation,
 * reactivation, length boundaries, expireAt, secret-sourced keys).
 *
 * Uses single-quoted YAML scalars so values containing backticks (e.g. the
 * GKO templating syntax `[[ secret \`name/key\` ]]`) round-trip cleanly.
 */
function subscriptionYaml(
  name: string,
  keys: Array<{ key: string; expireAt?: string }>,
): string {
  const keysYaml = keys
    .map((k) => {
      const escaped = k.key.replaceAll("'", "''");
      let line = `    - key: '${escaped}'`;
      if (k.expireAt) line += `\n      expireAt: "${k.expireAt}"`;
      return line;
    })
    .join("\n");
  return [
    "apiVersion: gravitee.io/v1alpha1",
    "kind: Subscription",
    "metadata:",
    `  name: ${name}`,
    "spec:",
    "  api:",
    `    name: "${API_NAME}"`,
    `    kind: "ApiV4Definition"`,
    `  plan: "ApiKey"`,
    "  application:",
    `    name: "${APP_NAME}"`,
    "  apiKeys:",
    keysYaml,
    "",
  ].join("\n");
}

/** Build a Secret manifest carrying a single `apiKey` field. */
function apiKeySecretYaml(name: string, value: string): string {
  return [
    "apiVersion: v1",
    "kind: Secret",
    "metadata:",
    `  name: ${name}`,
    "type: Opaque",
    "stringData:",
    `  apiKey: '${value.replaceAll("'", "''")}'`,
    "",
  ].join("\n");
}

test.describe("V4 API-Key Plan Subscriptions", () => {
  test.afterEach(async () => {
    // Reverse dependency order: subscriptions → applications → APIs.
    // Subscriptions go FIRST so any GKO templating finalizers on a Secret
    // are released before we try to delete the Secret itself.
    await kubectlSafe.del(fixture(SUB_APIKEY)).catch(() => {});
    await kubectlSafe.del(fixture(SUB_APIKEY_TWO_PLANS)).catch(() => {});
    for (const name of INLINE_SUB_NAMES) {
      await kubectlSafe.deleteResource("subscription", name).catch(() => {});
    }
    await kubectlSafe.del(fixture(APP)).catch(() => {});
    await kubectlSafe.del(fixture(API_APIKEY)).catch(() => {});
    await kubectlSafe.del(fixture(API_TWO_PLANS)).catch(() => {});
    await kubectlSafe.deleteResource("secret", SECRET_NAME).catch(() => {});
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
    const active = await mapi.listActiveSubscriptionApiKeys(apiId, subId);
    expect(active).toHaveLength(1);
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

    const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
    const apiKey = active?.key;
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

    const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
    const apiKey = active?.key;
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

    const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
    const apiKey = active?.key;
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

  // ── Custom api-key value is honored end-to-end ──

  test(`Custom api-key value is honored end-to-end ${XRAY.SUBSCRIPTIONS.V4_APIKEY_CUSTOM_KEY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
    gateway,
  }) => {
    // Budget mirrors the GKO-2826 gateway test: the api-key plan needs time
    // to propagate from APIM to the gateway, especially on slower CI.
    test.setTimeout(60_000);

    const SUB = "e2e-sub-apikey-custom";
    const CUSTOM_KEY = uniqueKey("e2e-custom-apikey");

    await kubectl.apply(fixture(API_APIKEY));
    await kubectl.apply(fixture(APP));
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");
    await kubectl.applyString(subscriptionYaml(SUB, [{ key: CUSTOM_KEY }]));
    await kubectl.waitForCondition("subscription", SUB, "Accepted");

    const { id: subId } = await kubectl.getStatus<{ id: string }>("subscription", SUB);
    const { id: apiId } = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);

    // The discriminator: the key value persisted in APIM is *exactly* the one
    // declared in the spec (not a regenerated value).
    const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
    expect(active.key).toBe(CUSTOM_KEY);

    await test.step("Gateway accepts the custom key from the spec", async () => {
      await gateway.assertResponds(API_PATH, {
        status: 200,
        headers: { "X-Gravitee-Api-Key": CUSTOM_KEY },
      });
    });
  });

  // ── CRD schema rejects api-keys outside 32–256 char bounds ──

  test(`CRD schema rejects api-keys outside 32-256 char bounds ${XRAY.SUBSCRIPTIONS.V4_APIKEY_LENGTH_REJECTED} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    await kubectl.apply(fixture(API_APIKEY));
    await kubectl.apply(fixture(APP));
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");

    await test.step("31-char key rejected (below minLength=32)", async () => {
      const stderr = await kubectl.applyStringExpectFailure(
        subscriptionYaml("e2e-sub-apikey-too-short", [{ key: "a".repeat(31) }]),
      );
      expect(stderr).toMatch(/spec\.apiKeys.*key|minLength|too short|invalid/i);
    });

    await test.step("257-char key rejected (above maxLength=256)", async () => {
      const stderr = await kubectl.applyStringExpectFailure(
        subscriptionYaml("e2e-sub-apikey-too-long", [{ key: "a".repeat(257) }]),
      );
      expect(stderr).toMatch(/spec\.apiKeys.*key|maxLength|too long|invalid/i);
    });
  });

  // ── CRD schema accepts 32 and 256 char api-keys at the boundaries ──

  test(`CRD schema accepts 32 and 256 char api-keys at boundaries ${XRAY.SUBSCRIPTIONS.V4_APIKEY_LENGTH_BOUNDARIES} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    test.setTimeout(60_000);

    await kubectl.apply(fixture(API_APIKEY));
    await kubectl.apply(fixture(APP));
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");

    const { id: apiId } = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const SUB = "e2e-sub-apikey-boundary";

    await test.step("32-char key accepted (lower boundary)", async () => {
      const k32 = uniqueKey("e2e-bnd-32"); // padded to exactly 32 chars
      expect(k32).toHaveLength(32);
      await kubectl.applyString(subscriptionYaml(SUB, [{ key: k32 }]));
      await kubectl.waitForCondition("subscription", SUB, "Accepted");
      const { id: subId } = await kubectl.getStatus<{ id: string }>("subscription", SUB);
      const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
      expect(active.key).toBe(k32);
      await kubectl.deleteResource("subscription", SUB);
      await kubectl.waitForDeletion("subscription", SUB);
    });

    await test.step("256-char key accepted (upper boundary)", async () => {
      // 256 chars: prefix + RUN_ID + 'y' padding to fill to the cap.
      const prefix = `e2e-bnd-256-${RUN_ID}-`;
      const k256 = prefix + "y".repeat(256 - prefix.length);
      expect(k256).toHaveLength(256);
      await kubectl.applyString(subscriptionYaml(SUB, [{ key: k256 }]));
      await kubectl.waitForCondition("subscription", SUB, "Accepted");
      const { id: subId } = await kubectl.getStatus<{ id: string }>("subscription", SUB);
      const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
      expect(active.key).toBe(k256);
    });
  });

  // ── expireAt is propagated from spec to APIM ──

  test(`expireAt on a custom api-key is propagated to APIM ${XRAY.SUBSCRIPTIONS.V4_APIKEY_EXPIRE_AT} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    test.setTimeout(60_000);

    await kubectl.apply(fixture(API_APIKEY));
    await kubectl.apply(fixture(APP));
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");

    const { id: apiId } = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const SUB = "e2e-sub-apikey-expire";
    const CUSTOM_KEY = uniqueKey("e2e-expire-apikey");
    // Admission rejects expireAt < 1 minute from now (validateEndingAt).
    // 30 min keeps the test deterministic without depending on real expiry.
    const expireAt = new Date(Date.now() + 30 * 60 * 1_000).toISOString();

    await kubectl.applyString(subscriptionYaml(SUB, [{ key: CUSTOM_KEY, expireAt }]));
    await kubectl.waitForCondition("subscription", SUB, "Accepted");

    const { id: subId } = await kubectl.getStatus<{ id: string }>("subscription", SUB);
    const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
    expect(active.key).toBe(CUSTOM_KEY);
    expect(active.expireAt).toBeTruthy();

    // APIM may serialise expireAt as RFC3339 or as epoch-ms; route through
    // Date so the comparison is representation-agnostic.
    const apimMs = new Date(active.expireAt as unknown as string | number).getTime();
    const expectedMs = new Date(expireAt).getTime();
    expect(Math.abs(apimMs - expectedMs)).toBeLessThan(2_000);
  });

  // ── Custom api-key sourced from a Kubernetes Secret ──

  test(`Custom api-key sourced from a Kubernetes Secret ${XRAY.SUBSCRIPTIONS.V4_APIKEY_SECRET_SOURCED} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
    gateway,
  }) => {
    test.setTimeout(60_000);

    const SECRET_KEY_VALUE = uniqueKey("secret-source-apikey");
    const TEMPLATE_REF = `[[ secret \`${SECRET_NAME}/apiKey\` ]]`;
    const SUB = "e2e-sub-apikey-secret";

    await kubectl.applyString(apiKeySecretYaml(SECRET_NAME, SECRET_KEY_VALUE));
    await kubectl.apply(fixture(API_APIKEY));
    await kubectl.apply(fixture(APP));
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");

    await kubectl.applyString(subscriptionYaml(SUB, [{ key: TEMPLATE_REF }]));
    await kubectl.waitForCondition("subscription", SUB, "Accepted");

    const { id: subId } = await kubectl.getStatus<{ id: string }>("subscription", SUB);
    const { id: apiId } = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);

    // The discriminator: APIM stores the *resolved* secret value, not the
    // template literal. A regression in the GKO templating engine would leave
    // the literal `[[ secret ... ]]` in APIM (or fail admission entirely).
    const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
    expect(active.key).toBe(SECRET_KEY_VALUE);

    await test.step("Gateway accepts the resolved secret-sourced key", async () => {
      await gateway.assertResponds(API_PATH, {
        status: 200,
        headers: { "X-Gravitee-Api-Key": SECRET_KEY_VALUE },
      });
    });

    // Explicit teardown so the GKO finalizer on the Secret is released
    // before afterEach attempts to delete the Secret.
    await kubectl.deleteResource("subscription", SUB);
    await kubectl.waitForDeletion("subscription", SUB);
  });

  // ── Instant key rotation: replace key A with key B ──

  test(`Instant api-key rotation revokes old key and activates new key ${XRAY.SUBSCRIPTIONS.V4_APIKEY_ROTATION_INSTANT} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
    gateway,
  }) => {
    // Budget: gateway accept-A poll, then APIM rotation poll, then gateway
    // reject-A and accept-B polls. Each gateway poll is up to 30s on CI.
    test.setTimeout(120_000);

    const SUB = "e2e-sub-apikey-rotation-instant";
    const KEY_A = uniqueKey("rotation-instant-A");
    const KEY_B = uniqueKey("rotation-instant-B");

    await kubectl.apply(fixture(API_APIKEY));
    await kubectl.apply(fixture(APP));
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");

    await kubectl.applyString(subscriptionYaml(SUB, [{ key: KEY_A }]));
    await kubectl.waitForCondition("subscription", SUB, "Accepted");

    const { id: subId } = await kubectl.getStatus<{ id: string }>("subscription", SUB);
    const { id: apiId } = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);

    const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
    expect(active.key).toBe(KEY_A);

    await gateway.assertResponds(API_PATH, {
      status: 200,
      headers: { "X-Gravitee-Api-Key": KEY_A },
    });

    await test.step("Replace KEY_A with KEY_B in spec", async () => {
      await kubectl.applyString(subscriptionYaml(SUB, [{ key: KEY_B }]));
    });

    await test.step("APIM revokes KEY_A and activates KEY_B", async () => {
      await expect
        .poll(
          async () => {
            const { data } = await mapi.listSubscriptionApiKeys(apiId, subId);
            return {
              revokedA: data.some((k) => k.key === KEY_A && k.revoked),
              activeB: data.some((k) => k.key === KEY_B && !k.revoked && !k.expired),
            };
          },
          { timeout: 30_000, message: "rotation reconcile to swap keys" },
        )
        .toMatchObject({ revokedA: true, activeB: true });
    });

    await test.step("Gateway rejects the rotated-out key", async () => {
      await gateway.assertNotResponds(API_PATH, {
        notStatus: 200,
        headers: { "X-Gravitee-Api-Key": KEY_A },
      });
    });

    await test.step("Gateway accepts the rotated-in key", async () => {
      await gateway.assertResponds(API_PATH, {
        status: 200,
        headers: { "X-Gravitee-Api-Key": KEY_B },
      });
    });
  });

  // ── Gradual rotation: two active keys, then deprecate the old ──

  test(`Gradual api-key rotation supports two active keys then deprecates the old ${XRAY.SUBSCRIPTIONS.V4_APIKEY_ROTATION_GRADUAL} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
    gateway,
  }) => {
    // Budget: this test asserts gateway behaviour at 3 rotation states,
    // each requiring a fresh gateway poll cycle.
    test.setTimeout(150_000);

    const SUB = "e2e-sub-apikey-rotation-gradual";
    const KEY_A = uniqueKey("rotation-gradual-A");
    const KEY_B = uniqueKey("rotation-gradual-B");

    await kubectl.apply(fixture(API_APIKEY));
    await kubectl.apply(fixture(APP));
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");

    await kubectl.applyString(subscriptionYaml(SUB, [{ key: KEY_A }]));
    await kubectl.waitForCondition("subscription", SUB, "Accepted");

    const { id: subId } = await kubectl.getStatus<{ id: string }>("subscription", SUB);
    const { id: apiId } = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);

    await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);

    await test.step("Add KEY_B alongside KEY_A; both active", async () => {
      await kubectl.applyString(subscriptionYaml(SUB, [{ key: KEY_A }, { key: KEY_B }]));
      await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 2, { timeoutMs: 30_000 });
    });

    await test.step("Gateway accepts both keys during overlap", async () => {
      await gateway.assertResponds(API_PATH, {
        status: 200,
        headers: { "X-Gravitee-Api-Key": KEY_A },
      });
      await gateway.assertResponds(API_PATH, {
        status: 200,
        headers: { "X-Gravitee-Api-Key": KEY_B },
      });
    });

    await test.step("Remove KEY_A; APIM revokes KEY_A and keeps KEY_B active", async () => {
      await kubectl.applyString(subscriptionYaml(SUB, [{ key: KEY_B }]));
      await expect
        .poll(
          async () => {
            const { data } = await mapi.listSubscriptionApiKeys(apiId, subId);
            return {
              revokedA: data.some((k) => k.key === KEY_A && k.revoked),
              activeB: data.some((k) => k.key === KEY_B && !k.revoked && !k.expired),
            };
          },
          { timeout: 30_000, message: "gradual rotation final reconcile" },
        )
        .toMatchObject({ revokedA: true, activeB: true });
    });

    await test.step("Gateway rejects the deprecated key and accepts the new one", async () => {
      await gateway.assertNotResponds(API_PATH, {
        notStatus: 200,
        headers: { "X-Gravitee-Api-Key": KEY_A },
      });
      await gateway.assertResponds(API_PATH, {
        status: 200,
        headers: { "X-Gravitee-Api-Key": KEY_B },
      });
    });
  });

  // ── Reactivation: a previously revoked key reappears in the spec ──

  test(`Previously revoked api-key is reactivated when re-added to spec ${XRAY.SUBSCRIPTIONS.V4_APIKEY_REACTIVATION} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
    gateway,
  }) => {
    test.setTimeout(150_000);

    const SUB = "e2e-sub-apikey-reactivation";
    const KEY_A = uniqueKey("reactivation-A");
    const KEY_B = uniqueKey("reactivation-B");

    await kubectl.apply(fixture(API_APIKEY));
    await kubectl.apply(fixture(APP));
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");

    await kubectl.applyString(subscriptionYaml(SUB, [{ key: KEY_A }]));
    await kubectl.waitForCondition("subscription", SUB, "Accepted");

    const { id: subId } = await kubectl.getStatus<{ id: string }>("subscription", SUB);
    const { id: apiId } = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);

    await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);

    await test.step("Replace KEY_A with KEY_B, KEY_A becomes revoked", async () => {
      await kubectl.applyString(subscriptionYaml(SUB, [{ key: KEY_B }]));
      await expect
        .poll(
          async () => {
            const { data } = await mapi.listSubscriptionApiKeys(apiId, subId);
            return {
              revokedA: data.some((k) => k.key === KEY_A && k.revoked),
              activeB: data.some((k) => k.key === KEY_B && !k.revoked && !k.expired),
            };
          },
          { timeout: 30_000, message: "first rotation revokes KEY_A" },
        )
        .toMatchObject({ revokedA: true, activeB: true });
    });

    await test.step("Re-add KEY_A alongside KEY_B; KEY_A is reactivated", async () => {
      await kubectl.applyString(subscriptionYaml(SUB, [{ key: KEY_A }, { key: KEY_B }]));
      await expect
        .poll(
          async () => {
            const { data } = await mapi.listSubscriptionApiKeys(apiId, subId);
            return {
              activeA: data.some((k) => k.key === KEY_A && !k.revoked && !k.expired),
              activeB: data.some((k) => k.key === KEY_B && !k.revoked && !k.expired),
            };
          },
          { timeout: 30_000, message: "KEY_A reactivated alongside KEY_B" },
        )
        .toMatchObject({ activeA: true, activeB: true });
    });

    await test.step("Gateway accepts the reactivated key", async () => {
      await gateway.assertResponds(API_PATH, {
        status: 200,
        headers: { "X-Gravitee-Api-Key": KEY_A },
      });
    });
  });

  // ── Idempotent reconcile: re-applying same spec creates no duplicates ──

  test(`Re-applying same custom-key spec does not create extra keys ${XRAY.SUBSCRIPTIONS.V4_APIKEY_IDEMPOTENT} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    test.setTimeout(60_000);

    const SUB = "e2e-sub-apikey-idempotent";
    const KEY = uniqueKey("idempotent-key");

    await kubectl.apply(fixture(API_APIKEY));
    await kubectl.apply(fixture(APP));
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");

    await kubectl.applyString(subscriptionYaml(SUB, [{ key: KEY }]));
    await kubectl.waitForCondition("subscription", SUB, "Accepted");

    const { id: subId } = await kubectl.getStatus<{ id: string }>("subscription", SUB);
    const { id: apiId } = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);

    const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
    expect(active.key).toBe(KEY);

    // Re-apply the identical spec twice; each apply re-triggers admission and
    // (on real diffs) the reconciler. A regression that re-creates the key on
    // every reconcile would surface as extra entries — active or revoked.
    await kubectl.applyString(subscriptionYaml(SUB, [{ key: KEY }]));
    await kubectl.applyString(subscriptionYaml(SUB, [{ key: KEY }]));

    // Allow any spurious reconcile churn time to manifest.
    await new Promise((r) => setTimeout(r, 5_000));

    const { data } = await mapi.listSubscriptionApiKeys(apiId, subId);
    const activeCount = data.filter((k) => !k.revoked && !k.expired).length;
    expect(activeCount).toBe(1);
    // Scope the duplicate check to *this run's* key value. APIM MongoDB
    // persists across cluster lifecycle (and across test files), so other
    // runs' keys appear here as revoked entries — that is not a regression.
    // A regression that recreates the same key value on each reconcile would
    // surface as multiple entries sharing this run's KEY string.
    const sameValueEntries = data.filter((k) => k.key === KEY);
    expect(sameValueEntries).toHaveLength(1);
  });

  // ── Multi-key subscription with staggered expirations: zero-downtime ──

  test(`Multi-key subscription with staggered expirations supports zero-downtime rotation ${XRAY.SUBSCRIPTIONS.V4_APIKEY_STAGGERED_EXPIRY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
    gateway,
  }) => {
    // Budget: 3 gateway accept polls, then a rotation reconcile, then a
    // gateway reject + 2 gateway accept polls.
    test.setTimeout(150_000);

    const SUB = "e2e-sub-apikey-staggered";
    // Three keys mirroring the canonical staged-rotation pattern from
    // GKO-2550: an early-expiring key on its way out, a longer-lived key,
    // and an evergreen key with no expireAt set.
    const KEY_V1 = uniqueKey("staggered-v1");
    const KEY_V2 = uniqueKey("staggered-v2");
    const KEY_V3 = uniqueKey("staggered-v3");
    const expireV1 = new Date(Date.now() + 30 * 60 * 1_000).toISOString();
    const expireV2 = new Date(Date.now() + 90 * 60 * 1_000).toISOString();

    await kubectl.apply(fixture(API_APIKEY));
    await kubectl.apply(fixture(APP));
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    await kubectl.waitForCondition("application", APP_NAME, "Accepted");

    await kubectl.applyString(
      subscriptionYaml(SUB, [
        { key: KEY_V1, expireAt: expireV1 },
        { key: KEY_V2, expireAt: expireV2 },
        { key: KEY_V3 },
      ]),
    );
    await kubectl.waitForCondition("subscription", SUB, "Accepted");

    const { id: subId } = await kubectl.getStatus<{ id: string }>("subscription", SUB);
    const { id: apiId } = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);

    await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 3, { timeoutMs: 30_000 });

    await test.step("APIM stores each key with the correct expireAt (or none)", async () => {
      const { data } = await mapi.listSubscriptionApiKeys(apiId, subId);
      const v1 = data.find((k) => k.key === KEY_V1);
      const v2 = data.find((k) => k.key === KEY_V2);
      const v3 = data.find((k) => k.key === KEY_V3);
      expect(v1, "v1 must be present in APIM").toBeDefined();
      expect(v2, "v2 must be present in APIM").toBeDefined();
      expect(v3, "v3 must be present in APIM").toBeDefined();

      // Both v1 and v2 must round-trip their expireAt within ~2s slop.
      const v1Ms = new Date(v1!.expireAt as unknown as string | number).getTime();
      const v2Ms = new Date(v2!.expireAt as unknown as string | number).getTime();
      expect(Math.abs(v1Ms - new Date(expireV1).getTime())).toBeLessThan(2_000);
      expect(Math.abs(v2Ms - new Date(expireV2).getTime())).toBeLessThan(2_000);
      // v3 has no expireAt in the spec — APIM must not synthesize one.
      expect(v3!.expireAt).toBeFalsy();
    });

    await test.step("Gateway accepts all three keys during the overlap window", async () => {
      await gateway.assertResponds(API_PATH, {
        status: 200,
        headers: { "X-Gravitee-Api-Key": KEY_V1 },
      });
      await gateway.assertResponds(API_PATH, {
        status: 200,
        headers: { "X-Gravitee-Api-Key": KEY_V2 },
      });
      await gateway.assertResponds(API_PATH, {
        status: 200,
        headers: { "X-Gravitee-Api-Key": KEY_V3 },
      });
    });

    // Drop the soon-to-expire v1 ahead of its expiry — the standard
    // "rotate out the oldest key without service interruption" move.
    await test.step("Drop KEY_V1; APIM revokes it while v2 and v3 stay active", async () => {
      await kubectl.applyString(
        subscriptionYaml(SUB, [
          { key: KEY_V2, expireAt: expireV2 },
          { key: KEY_V3 },
        ]),
      );
      await expect
        .poll(
          async () => {
            const { data } = await mapi.listSubscriptionApiKeys(apiId, subId);
            return {
              v1Revoked: data.some((k) => k.key === KEY_V1 && k.revoked),
              v2Active: data.some((k) => k.key === KEY_V2 && !k.revoked && !k.expired),
              v3Active: data.some((k) => k.key === KEY_V3 && !k.revoked && !k.expired),
            };
          },
          { timeout: 30_000, message: "staggered rotation revokes v1, keeps v2+v3" },
        )
        .toMatchObject({ v1Revoked: true, v2Active: true, v3Active: true });
    });

    await test.step("Gateway rejects v1 but continues to accept v2 and v3", async () => {
      await gateway.assertNotResponds(API_PATH, {
        notStatus: 200,
        headers: { "X-Gravitee-Api-Key": KEY_V1 },
      });
      await gateway.assertResponds(API_PATH, {
        status: 200,
        headers: { "X-Gravitee-Api-Key": KEY_V2 },
      });
      await gateway.assertResponds(API_PATH, {
        status: 200,
        headers: { "X-Gravitee-Api-Key": KEY_V3 },
      });
    });
  });
});
