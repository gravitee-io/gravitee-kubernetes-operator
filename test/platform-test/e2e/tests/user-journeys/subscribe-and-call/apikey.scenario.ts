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
 * V4 API-Key plan subscriptions, exercised through every supported provisioner
 * (currently GKO and Terraform) from a single shared intent. Each
 * `forEachProvisioner` block defines the platform behaviour once; the matrix runs
 * it against each provisioner, tags the arms (e.g. `@gko` / `@terraform`), and
 * carries the original per-provisioner Xray ids so coverage tracking is unchanged.
 *
 * Each scenario's key values are defined in a block scope and passed as the
 * initial params, so `provision()` applies the real keys directly (no throwaway
 * apply). Provisioner-specific behaviour (admission/templating for GKO; drift/
 * redaction for TF) lives in apikey-gko-only.test.ts / apikey-tf-only.test.ts.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test, expect } from "../../../setup.js";
import { XRAY, TAGS, since } from "../../../helpers/tags.js";
import { forEachProvisioner } from "../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../helpers/provisioner-env.js";
import { subscriptionYaml } from "../../../../src/provisioners/index.js";
import {
  APIKEY_GKO,
  gkoApplyApiKeys,
  tfApiKeyVars,
  uniqueKey,
  RUN_ID,
  type ApiKeyParams,
} from "./params.js";

/** This journey folder — fixtures are co-located under gko/ and terraform/. */
const here = path.dirname(fileURLToPath(import.meta.url));

/** GKO factory for the standard single-api-key-plan API + a given subscription name. */
function apikeyGko(subName: string) {
  return gkoScenario<ApiKeyParams>({
    manifests: [APIKEY_GKO.apiManifest, APIKEY_GKO.appManifest],
    roles: {
      api: APIKEY_GKO.apiName,
      application: APIKEY_GKO.appName,
      subscription: subName,
    },
    dynamicRoles: ["subscription"],
    contextPath: APIKEY_GKO.contextPath,
    applyParams: gkoApplyApiKeys(subName),
  });
}

/** Terraform factory for the parameterized apikey-custom fixture. */
function apikeyTfCustom(hridSuffix: string) {
  return tfScenario<ApiKeyParams>({
    fixture: path.join(here, "terraform/apikey-custom"),
    toVars: tfApiKeyVars(hridSuffix),
    // Lets provisioned.remove("subscription") drop the subscription from the
    // desired state and re-apply (the count-gated resource in apikey-custom).
    removeVars: { subscription: { create_subscription: false } },
  });
}

// api-key plan subscriptions ship in APIM 4.12, for both provisioners.
const REGRESSION = [TAGS.REGRESSION, since("4.12")];

// ── Auto-generated single key, reachable via gateway ──────────────────────────
// GKO splits this into a count check (2825) + a gateway check (2826); the TF
// auto test (2879) does both. The shared body covers both, for every provisioner arm.
forEachProvisioner<ApiKeyParams>(
  {
    title: "Single api-key generated on api-key plan subscription, reachable via gateway",
    provisioners: {
      gko: apikeyGko("e2e-sub-apikey-single"),
      terraform: tfScenario<ApiKeyParams>({ fixture: path.join(here, "terraform/apikey-auto") }),
    },
    xray: {
      gko: [XRAY.SUBSCRIPTIONS.V4_APIKEY_SINGLE_KEY, XRAY.SUBSCRIPTIONS.V4_APIKEY_GATEWAY_CALL],
      terraform: XRAY.TERRAFORM.APIKEY_AUTO_GENERATED,
    },
    tags: REGRESSION,
    timeoutMs: { gko: 60_000 },
  },
  async ({ provisioned, mapi, gateway }) => {
    const apiId = await provisioned.apiId();
    const subId = await provisioned.subscriptionId();
    const ctx = await provisioned.contextPath();

    // Exactly one active key, and no asynchronous second insert (APIM-13686).
    const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
    const apiKey = active.key;
    expect(apiKey).toBeTruthy();
    await new Promise((r) => setTimeout(r, 2_000));
    expect(await mapi.listActiveSubscriptionApiKeys(apiId, subId)).toHaveLength(1);

    await test.step("Gateway rejects without the key, accepts with it", async () => {
      await gateway.assertResponds(ctx, { status: 401 });
      await gateway.assertResponds(ctx, {
        status: 200,
        headers: { "X-Gravitee-Api-Key": apiKey },
      });
    });
  },
  {}, // no keys -> APIM auto-generates one
);

// ── Custom api-key value honored end-to-end ───────────────────────────────────
{
  const CUSTOM_KEY = uniqueKey("custom-apikey");
  forEachProvisioner<ApiKeyParams>(
    {
      title: "Custom api-key value is honored end-to-end",
      provisioners: {
        gko: apikeyGko("e2e-sub-apikey-custom"),
        terraform: apikeyTfCustom("custom-value"),
      },
      xray: { gko: XRAY.SUBSCRIPTIONS.V4_APIKEY_CUSTOM_KEY, terraform: XRAY.TERRAFORM.APIKEY_CUSTOM_VALUE },
      tags: REGRESSION,
      timeoutMs: { gko: 60_000 },
    },
    async ({ provisioned, mapi, gateway }) => {
      const apiId = await provisioned.apiId();
      const subId = await provisioned.subscriptionId();
      const ctx = await provisioned.contextPath();

      // The discriminator: APIM persists exactly the declared value.
      const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
      expect(active.key).toBe(CUSTOM_KEY);

      await test.step("Gateway accepts the custom key, rejects without it", async () => {
        await gateway.assertResponds(ctx, {
          status: 200,
          headers: { "X-Gravitee-Api-Key": CUSTOM_KEY },
        });
        await gateway.assertResponds(ctx, { status: 401 });
      });
    },
    { keys: [{ key: CUSTOM_KEY }] },
  );
}

// ── expireAt is propagated to APIM ────────────────────────────────────────────
{
  const KEY = uniqueKey("expire-apikey");
  // 30 min keeps the test deterministic (admission rejects expireAt < 1 min).
  const expireAt = new Date(Date.now() + 30 * 60 * 1_000).toISOString();
  forEachProvisioner<ApiKeyParams>(
    {
      title: "expireAt on a custom api-key is propagated to APIM",
      provisioners: {
        gko: apikeyGko("e2e-sub-apikey-expire"),
        terraform: apikeyTfCustom("expire-at"),
      },
      xray: { gko: XRAY.SUBSCRIPTIONS.V4_APIKEY_EXPIRE_AT, terraform: XRAY.TERRAFORM.APIKEY_EXPIRE_AT },
      tags: REGRESSION,
      timeoutMs: { gko: 60_000 },
    },
    async ({ provisioned, mapi }) => {
      const apiId = await provisioned.apiId();
      const subId = await provisioned.subscriptionId();

      const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
      expect(active.key).toBe(KEY);
      expect(active.expireAt).toBeTruthy();

      // APIM may serialise expireAt as RFC3339 or epoch-ms; route through Date.
      const apimMs = new Date(active.expireAt as unknown as string | number).getTime();
      expect(Math.abs(apimMs - new Date(expireAt).getTime())).toBeLessThan(2_000);
    },
    { keys: [{ key: KEY, expireAt }] },
  );
}

// ── Length boundaries: 32 and 256 char keys are accepted ──────────────────────
{
  const k32 = uniqueKey("bnd-32"); // padded to exactly 32 chars
  const prefix = `bnd-256-${RUN_ID}-`;
  const k256 = prefix + "y".repeat(256 - prefix.length);
  forEachProvisioner<ApiKeyParams>(
    {
      title: "Api-keys at the 32 and 256 char boundaries are accepted",
      provisioners: {
        gko: apikeyGko("e2e-sub-apikey-boundary"),
        terraform: apikeyTfCustom("boundary"),
      },
      xray: { gko: XRAY.SUBSCRIPTIONS.V4_APIKEY_LENGTH_BOUNDARIES, terraform: XRAY.TERRAFORM.APIKEY_LENGTH_BOUNDARIES },
      tags: REGRESSION,
      timeoutMs: { gko: 90_000 },
    },
    async ({ provisioned, mapi }) => {
      const apiId = await provisioned.apiId();
      const subId = await provisioned.subscriptionId();

      await test.step("32-char key accepted (lower boundary)", async () => {
        expect(k32).toHaveLength(32);
        const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
        expect(active.key).toBe(k32);
      });

      await test.step("256-char key accepted (upper boundary)", async () => {
        expect(k256).toHaveLength(256);
        await provisioned.update({ keys: [{ key: k256 }] });
        await expect
          .poll(
            async () => {
              const { data } = await mapi.listSubscriptionApiKeys(apiId, subId);
              return data.some((k) => k.key === k256 && !k.revoked && !k.expired);
            },
            { timeout: 30_000, message: "256-char key reaches active state" },
          )
          .toBe(true);
      });
    },
    { keys: [{ key: k32 }] },
  );
}

// ── Instant rotation: replace KEY_A with KEY_B ────────────────────────────────
{
  const KEY_A = uniqueKey("rot-instant-A");
  const KEY_B = uniqueKey("rot-instant-B");
  forEachProvisioner<ApiKeyParams>(
    {
      title: "Instant api-key rotation revokes old key and activates new key",
      provisioners: {
        gko: apikeyGko("e2e-sub-apikey-rotation-instant"),
        terraform: apikeyTfCustom("rot-instant"),
      },
      xray: { gko: XRAY.SUBSCRIPTIONS.V4_APIKEY_ROTATION_INSTANT, terraform: XRAY.TERRAFORM.APIKEY_ROTATION_INSTANT },
      tags: REGRESSION,
      timeoutMs: { gko: 150_000 },
    },
    async ({ provisioned, mapi, gateway }) => {
      const apiId = await provisioned.apiId();
      const subId = await provisioned.subscriptionId();
      const ctx = await provisioned.contextPath();

      const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
      expect(active.key).toBe(KEY_A);
      await gateway.assertResponds(ctx, { status: 200, headers: { "X-Gravitee-Api-Key": KEY_A } });

      await test.step("Replace KEY_A with KEY_B", async () => {
        await provisioned.update({ keys: [{ key: KEY_B }] });
        await expect
          .poll(
            async () => {
              const { data } = await mapi.listSubscriptionApiKeys(apiId, subId);
              return {
                revokedA: data.some((k) => k.key === KEY_A && k.revoked),
                activeB: data.some((k) => k.key === KEY_B && !k.revoked && !k.expired),
              };
            },
            { timeout: 30_000, message: "rotation reconcile swaps the keys" },
          )
          .toMatchObject({ revokedA: true, activeB: true });
      });

      await test.step("Gateway rejects the rotated-out key, accepts the rotated-in key", async () => {
        await gateway.assertNotResponds(ctx, { notStatus: 200, headers: { "X-Gravitee-Api-Key": KEY_A } });
        await gateway.assertResponds(ctx, { status: 200, headers: { "X-Gravitee-Api-Key": KEY_B } });
      });
    },
    { keys: [{ key: KEY_A }] },
  );
}

// ── Gradual rotation: two active keys, then deprecate the old ─────────────────
{
  const KEY_A = uniqueKey("rot-gradual-A");
  const KEY_B = uniqueKey("rot-gradual-B");
  forEachProvisioner<ApiKeyParams>(
    {
      title: "Gradual api-key rotation supports two active keys then deprecates the old",
      provisioners: {
        gko: apikeyGko("e2e-sub-apikey-rotation-gradual"),
        terraform: apikeyTfCustom("rot-gradual"),
      },
      xray: { gko: XRAY.SUBSCRIPTIONS.V4_APIKEY_ROTATION_GRADUAL, terraform: XRAY.TERRAFORM.APIKEY_ROTATION_GRADUAL },
      tags: REGRESSION,
      timeoutMs: { gko: 180_000 },
    },
    async ({ provisioned, mapi, gateway }) => {
      const apiId = await provisioned.apiId();
      const subId = await provisioned.subscriptionId();
      const ctx = await provisioned.contextPath();
      await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);

      await test.step("Add KEY_B alongside KEY_A; both active", async () => {
        await provisioned.update({ keys: [{ key: KEY_A }, { key: KEY_B }] });
        await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 2, { timeoutMs: 30_000 });
        await gateway.assertResponds(ctx, { status: 200, headers: { "X-Gravitee-Api-Key": KEY_A } });
        await gateway.assertResponds(ctx, { status: 200, headers: { "X-Gravitee-Api-Key": KEY_B } });
      });

      await test.step("Remove KEY_A; APIM revokes it and keeps KEY_B active", async () => {
        await provisioned.update({ keys: [{ key: KEY_B }] });
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
        await gateway.assertNotResponds(ctx, { notStatus: 200, headers: { "X-Gravitee-Api-Key": KEY_A } });
        await gateway.assertResponds(ctx, { status: 200, headers: { "X-Gravitee-Api-Key": KEY_B } });
      });
    },
    { keys: [{ key: KEY_A }] },
  );
}

// ── Reactivation: a revoked key reappears in the spec ─────────────────────────
{
  const KEY_A = uniqueKey("reactivation-A");
  const KEY_B = uniqueKey("reactivation-B");
  forEachProvisioner<ApiKeyParams>(
    {
      title: "Previously revoked api-key is reactivated when re-added to spec",
      provisioners: {
        gko: apikeyGko("e2e-sub-apikey-reactivation"),
        terraform: apikeyTfCustom("reactivation"),
      },
      xray: { gko: XRAY.SUBSCRIPTIONS.V4_APIKEY_REACTIVATION, terraform: XRAY.TERRAFORM.APIKEY_REACTIVATION },
      tags: REGRESSION,
      timeoutMs: { gko: 180_000 },
    },
    async ({ provisioned, mapi, gateway }) => {
      const apiId = await provisioned.apiId();
      const subId = await provisioned.subscriptionId();
      const ctx = await provisioned.contextPath();
      await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);

      await test.step("Replace KEY_A with KEY_B; KEY_A becomes revoked", async () => {
        await provisioned.update({ keys: [{ key: KEY_B }] });
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
        await provisioned.update({ keys: [{ key: KEY_A }, { key: KEY_B }] });
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
        await gateway.assertResponds(ctx, { status: 200, headers: { "X-Gravitee-Api-Key": KEY_A } });
      });
    },
    { keys: [{ key: KEY_A }] },
  );
}

// ── Staggered expirations: zero-downtime rotation ─────────────────────────────
{
  const KEY_V1 = uniqueKey("stagger-v1");
  const KEY_V2 = uniqueKey("stagger-v2");
  const KEY_V3 = uniqueKey("stagger-v3");
  const expireV1 = new Date(Date.now() + 30 * 60 * 1_000).toISOString();
  const expireV2 = new Date(Date.now() + 90 * 60 * 1_000).toISOString();
  forEachProvisioner<ApiKeyParams>(
    {
      title: "Multi-key subscription with staggered expirations supports zero-downtime rotation",
      provisioners: {
        gko: apikeyGko("e2e-sub-apikey-staggered"),
        terraform: apikeyTfCustom("staggered"),
      },
      xray: { gko: XRAY.SUBSCRIPTIONS.V4_APIKEY_STAGGERED_EXPIRY, terraform: XRAY.TERRAFORM.APIKEY_STAGGERED_EXPIRY },
      tags: REGRESSION,
      timeoutMs: { gko: 180_000 },
    },
    async ({ provisioned, mapi, gateway }) => {
      const apiId = await provisioned.apiId();
      const subId = await provisioned.subscriptionId();
      const ctx = await provisioned.contextPath();
      await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 3, { timeoutMs: 30_000 });

      await test.step("APIM stores each key with the correct expireAt (or none)", async () => {
        const { data } = await mapi.listSubscriptionApiKeys(apiId, subId);
        const v1 = data.find((k) => k.key === KEY_V1);
        const v2 = data.find((k) => k.key === KEY_V2);
        const v3 = data.find((k) => k.key === KEY_V3);
        expect(v1, "v1 present").toBeDefined();
        expect(v2, "v2 present").toBeDefined();
        expect(v3, "v3 present").toBeDefined();
        const v1Ms = new Date(v1!.expireAt as unknown as string | number).getTime();
        const v2Ms = new Date(v2!.expireAt as unknown as string | number).getTime();
        expect(Math.abs(v1Ms - new Date(expireV1).getTime())).toBeLessThan(2_000);
        expect(Math.abs(v2Ms - new Date(expireV2).getTime())).toBeLessThan(2_000);
        expect(v3!.expireAt).toBeFalsy();
      });

      await test.step("Gateway accepts all three keys during the overlap", async () => {
        for (const key of [KEY_V1, KEY_V2, KEY_V3]) {
          await gateway.assertResponds(ctx, { status: 200, headers: { "X-Gravitee-Api-Key": key } });
        }
      });

      await test.step("Drop KEY_V1; APIM revokes it while v2 and v3 stay active", async () => {
        await provisioned.update({ keys: [{ key: KEY_V2, expireAt: expireV2 }, { key: KEY_V3 }] });
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
        await gateway.assertNotResponds(ctx, { notStatus: 200, headers: { "X-Gravitee-Api-Key": KEY_V1 } });
        await gateway.assertResponds(ctx, { status: 200, headers: { "X-Gravitee-Api-Key": KEY_V2 } });
        await gateway.assertResponds(ctx, { status: 200, headers: { "X-Gravitee-Api-Key": KEY_V3 } });
      });
    },
    { keys: [{ key: KEY_V1, expireAt: expireV1 }, { key: KEY_V2, expireAt: expireV2 }, { key: KEY_V3 }] },
  );
}

// ── Mixed: api-key plan coexists with a keyless plan on the same API ──────────
{
  const CUSTOM_KEY = uniqueKey("mixed-apikey");
  forEachProvisioner<ApiKeyParams>(
    {
      title: "Api-key plan coexists with keyless plan on the same API",
      provisioners: {
        gko: gkoScenario<ApiKeyParams>({
          manifests: [path.join(here, "gko/api-mixed.yaml")],
          roles: {
            api: "e2e-v4-two-plans",
            application: "e2e-app-simple",
            subscription: "e2e-sub-apikey-two-plans",
          },
          dynamicRoles: ["subscription"],
          contextPath: "/e2e-v4-two-plans",
          applyParams: async (kubectl, params: ApiKeyParams) => {
            await kubectl.applyString(
              subscriptionYaml({
                name: "e2e-sub-apikey-two-plans",
                apiName: "e2e-v4-two-plans",
                plan: "ApiKey",
                applicationName: "e2e-app-simple",
                keys: params.keys ?? [],
              }),
            );
          },
        }),
        terraform: tfScenario<ApiKeyParams>({
          fixture: path.join(here, "terraform/apikey-mixed"),
          toVars: tfApiKeyVars("mixed"),
        }),
      },
      xray: { gko: XRAY.SUBSCRIPTIONS.V4_APIKEY_MIXED_WITH_KEYLESS, terraform: XRAY.TERRAFORM.APIKEY_MIXED_WITH_KEYLESS },
      tags: REGRESSION,
      timeoutMs: { gko: 90_000 },
    },
    async ({ provisioned, mapi, gateway }) => {
      const apiId = await provisioned.apiId();
      const subId = await provisioned.subscriptionId();
      const ctx = await provisioned.contextPath();
      const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
      expect(active.key).toBe(CUSTOM_KEY);

      await test.step("Keyless plan accepts traffic with no api-key header", async () => {
        await gateway.assertResponds(ctx, { status: 200 });
      });
      await test.step("Api-key plan accepts traffic with a valid api-key header", async () => {
        await gateway.assertResponds(ctx, { status: 200, headers: { "X-Gravitee-Api-Key": CUSTOM_KEY } });
      });
      // Discriminator: an invalid api-key header is handled by the api-key plan
      // (not keyless), so the gateway must reject it.
      await test.step("Api-key plan rejects traffic with an invalid api-key header", async () => {
        await gateway.assertNotResponds(ctx, {
          notStatus: 200,
          headers: { "X-Gravitee-Api-Key": "bogus-invalid-key" },
        });
      });
    },
    { keys: [{ key: CUSTOM_KEY }] },
  );
}

// ── Removing only the subscription revokes its key while the API stays up ─────
// Each provisioner removes the subscription the way a user would: GKO deletes
// the Subscription CR; Terraform drops it from the desired state and re-applies.
// The API/app stay up, so a non-200 proves the KEY was revoked, not that the
// endpoint vanished.
{
  const CUSTOM_KEY = uniqueKey("revoke-apikey");
  forEachProvisioner<ApiKeyParams>(
    {
      title: "Removing the subscription revokes its api-key while the API stays up",
      provisioners: {
        gko: apikeyGko("e2e-sub-apikey-revoke"),
        terraform: apikeyTfCustom("revoke"),
      },
      xray: { gko: XRAY.SUBSCRIPTIONS.V4_APIKEY_KEY_REVOKED_ON_DELETE, terraform: XRAY.TERRAFORM.APIKEY_REVOKED_ON_DESTROY },
      tags: REGRESSION,
      timeoutMs: { gko: 90_000 },
    },
    async ({ provisioned, mapi, gateway }) => {
      const apiId = await provisioned.apiId();
      const subId = await provisioned.subscriptionId();
      const ctx = await provisioned.contextPath();
      await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
      await gateway.assertResponds(ctx, { status: 200, headers: { "X-Gravitee-Api-Key": CUSTOM_KEY } });

      await test.step("Remove only the subscription; the key stops working, the API stays up", async () => {
        await provisioned.remove("subscription");
        await gateway.assertNotResponds(ctx, { notStatus: 200, headers: { "X-Gravitee-Api-Key": CUSTOM_KEY } });
      });
    },
    { keys: [{ key: CUSTOM_KEY }] },
  );
}
