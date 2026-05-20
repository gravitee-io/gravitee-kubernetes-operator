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
 * V4 API-Key plan subscriptions driven by the Gravitee Terraform Provider.
 *
 * Covers GKO-2560 / terraform-provider-apim PR #122 which replaced the legacy
 * scalar `custom_api_key` attribute with an `api_keys` list (each entry: a
 * required Sensitive `key` + optional RFC3339 `expire_at`). These tests
 * mirror the GKO CRD-driven coverage in v4-subscriptions-apikey.test.ts so
 * that both write paths into the Automation API are exercised against the
 * same APIM behaviours (rotation, expireAt round-trip, gateway routing).
 *
 * Preconditions:
 *   - APIM and Gateway are running
 *   - terraform CLI is installed (>= 1.3 for optional() in object types)
 */

import { test, expect } from "../../setup.js";
import { XRAY, TAGS } from "../../helpers/tags.js";
import * as terraform from "../../helpers/terraform.js";
import type { TfWorkspace } from "../../helpers/terraform.js";

// APIM enforces uniqueness on api-key values per API across active and
// revoked states, and APIM MongoDB persists across cluster lifecycle on the
// local setup, so re-running with hardcoded key values yields
// "API Key already exists" 400s on re-runs. A per-process suffix makes every
// run pick fresh values. See [[apim_apikey_value_uniqueness]].
const RUN_ID = `${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 6)}`;

/** Generate a unique api-key value of at least 32 chars. */
function uniqueKey(prefix: string): string {
  return `${prefix}-${RUN_ID}`.padEnd(32, "0");
}

interface KeyEntry {
  key: string;
  // null (rather than undefined-then-omitted) keeps the JSON shape stable
  // across calls and is treated as "no expireAt" by the TF provider's
  // optional(string) attribute.
  expire_at?: string | null;
}

/** Build the JSON-encoded keys list passed to terraform via auto.tfvars.json. */
function tfKeys(keys: KeyEntry[]): KeyEntry[] {
  return keys.map((k) => ({ key: k.key, expire_at: k.expire_at ?? null }));
}

test.describe("Terraform — V4 API-Key Plan Subscriptions", () => {
  // ── Auto-generated key when no api_keys block is set ──────────────

  test(`Auto-generated api-key when no api_keys block is set ${XRAY.TERRAFORM.APIKEY_AUTO_GENERATED} ${TAGS.REGRESSION}`, async ({
    mapi,
    gateway,
  }) => {
    test.setTimeout(120_000);

    const ws = await terraform.initWorkspace("terraform-apikey-auto");
    try {
      await terraform.apply(ws);
      const apiId = await terraform.output(ws, "api_id");
      const subId = await terraform.output(ws, "sub_id");
      const ctx = await terraform.output(ws, "api_context_path");

      // APIM auto-generates exactly one active key when api_keys is omitted.
      // Re-poll after a short delay to rule out an asynchronous second insert
      // (APIM-13686 produced two keys ~1s apart on the GKO path).
      const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
      const apiKey = active?.key;
      expect(apiKey).toBeTruthy();
      await new Promise((r) => setTimeout(r, 2_000));
      const stillSingle = await mapi.listActiveSubscriptionApiKeys(apiId, subId);
      expect(stillSingle).toHaveLength(1);

      await gateway.assertResponds(ctx, {
        status: 200,
        headers: { "X-Gravitee-Api-Key": apiKey },
      });
    } finally {
      await terraform.destroyWorkspace(ws);
    }
  });

  // ── Custom api-key value is honored end-to-end ────────────────────

  test(`Custom api-key value is honored end-to-end ${XRAY.TERRAFORM.APIKEY_CUSTOM_VALUE} ${TAGS.REGRESSION}`, async ({
    mapi,
    gateway,
  }) => {
    test.setTimeout(120_000);

    const ws = await terraform.initWorkspace("terraform-apikey-custom");
    try {
      const CUSTOM_KEY = uniqueKey("tf-custom-apikey");
      await terraform.writeVars(ws, {
        hrid_suffix: "custom-value",
        keys: tfKeys([{ key: CUSTOM_KEY }]),
      });
      await terraform.apply(ws);

      const apiId = await terraform.output(ws, "api_id");
      const subId = await terraform.output(ws, "sub_id");
      const ctx = await terraform.output(ws, "api_context_path");

      // The discriminator: APIM persists *exactly* the key declared in the
      // spec (not a regenerated value).
      const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
      expect(active.key).toBe(CUSTOM_KEY);

      await test.step("Gateway accepts the custom key", async () => {
        await gateway.assertResponds(ctx, {
          status: 200,
          headers: { "X-Gravitee-Api-Key": CUSTOM_KEY },
        });
      });

      await test.step("Gateway rejects calls with no api-key header", async () => {
        await gateway.assertResponds(ctx, { status: 401 });
      });
    } finally {
      await terraform.destroyWorkspace(ws);
    }
  });

  // ── expire_at is propagated to APIM ───────────────────────────────

  test(`expire_at on a custom api-key is propagated to APIM ${XRAY.TERRAFORM.APIKEY_EXPIRE_AT} ${TAGS.REGRESSION}`, async ({
    mapi,
  }) => {
    test.setTimeout(120_000);

    const ws = await terraform.initWorkspace("terraform-apikey-custom");
    try {
      const CUSTOM_KEY = uniqueKey("tf-expire-apikey");
      // 30 min keeps the test deterministic without depending on real expiry.
      const expireAt = new Date(Date.now() + 30 * 60 * 1_000).toISOString();
      await terraform.writeVars(ws, {
        hrid_suffix: "expire-at",
        keys: tfKeys([{ key: CUSTOM_KEY, expire_at: expireAt }]),
      });
      await terraform.apply(ws);

      const apiId = await terraform.output(ws, "api_id");
      const subId = await terraform.output(ws, "sub_id");

      const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
      expect(active.key).toBe(CUSTOM_KEY);
      expect(active.expireAt).toBeTruthy();

      // APIM may serialise expireAt as RFC3339 or as epoch-ms; route through
      // Date so the comparison is representation-agnostic.
      const apimMs = new Date(active.expireAt as unknown as string | number).getTime();
      const expectedMs = new Date(expireAt).getTime();
      expect(Math.abs(apimMs - expectedMs)).toBeLessThan(2_000);
    } finally {
      await terraform.destroyWorkspace(ws);
    }
  });

  // ── Mixed: api-key plan coexists with keyless plan ────────────────

  test(`Api-key plan coexists with keyless plan on the same API ${XRAY.TERRAFORM.APIKEY_MIXED_WITH_KEYLESS} ${TAGS.REGRESSION}`, async ({
    mapi,
    gateway,
  }) => {
    test.setTimeout(120_000);

    const ws = await terraform.initWorkspace("terraform-apikey-mixed");
    try {
      const CUSTOM_KEY = uniqueKey("tf-mixed-apikey");
      await terraform.writeVars(ws, {
        hrid_suffix: "mixed",
        keys: tfKeys([{ key: CUSTOM_KEY }]),
      });
      await terraform.apply(ws);

      const apiId = await terraform.output(ws, "api_id");
      const subId = await terraform.output(ws, "sub_id");
      const ctx = await terraform.output(ws, "api_context_path");

      const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
      expect(active.key).toBe(CUSTOM_KEY);

      await test.step("Keyless plan accepts traffic with no api-key header", async () => {
        await gateway.assertResponds(ctx, { status: 200 });
      });

      await test.step("Api-key plan accepts traffic with a valid api-key header", async () => {
        await gateway.assertResponds(ctx, {
          status: 200,
          headers: { "X-Gravitee-Api-Key": CUSTOM_KEY },
        });
      });

      // Discriminator: an invalid api-key header is handled by the api-key
      // plan (not keyless), so the gateway must reject. Without this step,
      // the valid-key assertion would pass even if api-key plan resolution
      // regressed and keyless silently handled every header-bearing request.
      await test.step("Api-key plan rejects traffic with an invalid api-key header", async () => {
        await gateway.assertNotResponds(ctx, {
          notStatus: 200,
          headers: { "X-Gravitee-Api-Key": "bogus-invalid-key" },
        });
      });
    } finally {
      await terraform.destroyWorkspace(ws);
    }
  });

  // ── Length boundaries: 32 and 256 char keys are accepted ──────────

  test(`Api-keys at the 32 and 256 char boundaries are accepted ${XRAY.TERRAFORM.APIKEY_LENGTH_BOUNDARIES} ${TAGS.REGRESSION}`, async ({
    mapi,
  }) => {
    test.setTimeout(180_000);

    const ws = await terraform.initWorkspace("terraform-apikey-custom");
    try {
      const k32 = uniqueKey("tf-bnd-32"); // padded to exactly 32 chars
      expect(k32).toHaveLength(32);
      await terraform.writeVars(ws, {
        hrid_suffix: "boundary",
        keys: tfKeys([{ key: k32 }]),
      });
      await terraform.apply(ws);

      const apiId = await terraform.output(ws, "api_id");
      const subId = await terraform.output(ws, "sub_id");

      await test.step("32-char key accepted (lower boundary)", async () => {
        const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
        expect(active.key).toBe(k32);
      });

      // Rotate in-place to a 256-char key (upper boundary). Same subscription,
      // just a new key, so APIM revokes k32 and activates k256.
      const prefix = `tf-bnd-256-${RUN_ID}-`;
      const k256 = prefix + "y".repeat(256 - prefix.length);
      expect(k256).toHaveLength(256);
      await terraform.writeVars(ws, {
        hrid_suffix: "boundary",
        keys: tfKeys([{ key: k256 }]),
      });
      await terraform.apply(ws);

      await test.step("256-char key accepted (upper boundary)", async () => {
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
    } finally {
      await terraform.destroyWorkspace(ws);
    }
  });

  // ── Length: out-of-bounds keys rejected (server-side) ─────────────
  //
  // The TF provider schema has no length validators on `key`, so the limit
  // (if any) is enforced server-side by the Automation API. A 31-char key
  // is below the GKO-CRD minLength=32 and is expected to be rejected by
  // APIM with a 400. The discriminator is "apply succeeded vs. apply
  // failed with a length-related error"; we don't assert on the exact
  // message wording since it comes from the server.

  test(`Out-of-bounds api-keys are rejected at apply time ${XRAY.TERRAFORM.APIKEY_LENGTH_REJECTED} ${TAGS.REGRESSION}`, async ({
    mapi,
  }) => {
    test.setTimeout(180_000);

    const ws = await terraform.initWorkspace("terraform-apikey-custom");
    try {
      // 31-char key (below minLength=32 on the CRD; APIM is expected to
      // enforce the same bound server-side via the Automation API).
      const shortKey = `tf-too-short-${RUN_ID}`.padEnd(31, "0").slice(0, 31);
      expect(shortKey).toHaveLength(31);
      await terraform.writeVars(ws, {
        hrid_suffix: "too-short",
        keys: tfKeys([{ key: shortKey }]),
      });

      await test.step("31-char key is rejected", async () => {
        const out = await terraform.applyExpectFailure(ws);
        // Scope to terms tied to length/validation rather than any generic
        // "error" so a provider-auth failure wouldn't accidentally pass.
        expect(out.toLowerCase()).toMatch(
          /minlength|too short|length|invalid|400|bad request|validation/,
        );
        // The api/app should still be created; the subscription is the part
        // that fails. Drop the workspace state so the next sub-step's apply
        // sees a clean plan.
      });

      // 257-char key (above maxLength=256). Re-using the same workspace —
      // anything created so far will be cleaned up by destroyWorkspace.
      const longPrefix = `tf-too-long-${RUN_ID}-`;
      const longKey = longPrefix + "z".repeat(257 - longPrefix.length);
      expect(longKey).toHaveLength(257);
      await terraform.writeVars(ws, {
        hrid_suffix: "too-long",
        keys: tfKeys([{ key: longKey }]),
      });

      await test.step("257-char key is rejected", async () => {
        const out = await terraform.applyExpectFailure(ws);
        expect(out.toLowerCase()).toMatch(
          /maxlength|too long|length|invalid|400|bad request|validation/,
        );
      });

      // Sanity: no active keys exist on either api+subscription combo
      // (best-effort; if the api/app weren't created we'll just skip).
      const apiIdOrEmpty = await terraform.output(ws, "api_id").catch(() => "");
      const subIdOrEmpty = await terraform.output(ws, "sub_id").catch(() => "");
      if (apiIdOrEmpty && subIdOrEmpty) {
        const active = await mapi
          .listActiveSubscriptionApiKeys(apiIdOrEmpty, subIdOrEmpty)
          .catch(() => []);
        expect(active).toHaveLength(0);
      }
    } finally {
      await terraform.destroyWorkspace(ws);
    }
  });

  // ── Instant rotation: replace KEY_A with KEY_B ────────────────────

  test(`Instant api-key rotation revokes old key and activates new key ${XRAY.TERRAFORM.APIKEY_ROTATION_INSTANT} ${TAGS.REGRESSION}`, async ({
    mapi,
    gateway,
  }) => {
    test.setTimeout(180_000);

    const ws = await terraform.initWorkspace("terraform-apikey-custom");
    try {
      const KEY_A = uniqueKey("tf-rot-instant-A");
      const KEY_B = uniqueKey("tf-rot-instant-B");

      await terraform.writeVars(ws, {
        hrid_suffix: "rot-instant",
        keys: tfKeys([{ key: KEY_A }]),
      });
      await terraform.apply(ws);

      const apiId = await terraform.output(ws, "api_id");
      const subId = await terraform.output(ws, "sub_id");
      const ctx = await terraform.output(ws, "api_context_path");

      const [active] = await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);
      expect(active.key).toBe(KEY_A);

      await gateway.assertResponds(ctx, {
        status: 200,
        headers: { "X-Gravitee-Api-Key": KEY_A },
      });

      await test.step("Replace KEY_A with KEY_B in the TF spec", async () => {
        await terraform.writeVars(ws, {
          hrid_suffix: "rot-instant",
          keys: tfKeys([{ key: KEY_B }]),
        });
        await terraform.apply(ws);
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
            { timeout: 30_000, message: "rotation reconcile swaps the keys" },
          )
          .toMatchObject({ revokedA: true, activeB: true });
      });

      await test.step("Gateway rejects the rotated-out key", async () => {
        await gateway.assertNotResponds(ctx, {
          notStatus: 200,
          headers: { "X-Gravitee-Api-Key": KEY_A },
        });
      });

      await test.step("Gateway accepts the rotated-in key", async () => {
        await gateway.assertResponds(ctx, {
          status: 200,
          headers: { "X-Gravitee-Api-Key": KEY_B },
        });
      });
    } finally {
      await terraform.destroyWorkspace(ws);
    }
  });

  // ── Gradual rotation: two active keys, then deprecate the old ─────

  test(`Gradual api-key rotation supports two active keys then deprecates the old ${XRAY.TERRAFORM.APIKEY_ROTATION_GRADUAL} ${TAGS.REGRESSION}`, async ({
    mapi,
    gateway,
  }) => {
    test.setTimeout(240_000);

    const ws = await terraform.initWorkspace("terraform-apikey-custom");
    try {
      const KEY_A = uniqueKey("tf-rot-gradual-A");
      const KEY_B = uniqueKey("tf-rot-gradual-B");

      await terraform.writeVars(ws, {
        hrid_suffix: "rot-gradual",
        keys: tfKeys([{ key: KEY_A }]),
      });
      await terraform.apply(ws);

      const apiId = await terraform.output(ws, "api_id");
      const subId = await terraform.output(ws, "sub_id");
      const ctx = await terraform.output(ws, "api_context_path");

      await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);

      await test.step("Add KEY_B alongside KEY_A; both are active", async () => {
        await terraform.writeVars(ws, {
          hrid_suffix: "rot-gradual",
          keys: tfKeys([{ key: KEY_A }, { key: KEY_B }]),
        });
        await terraform.apply(ws);
        await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 2, { timeoutMs: 30_000 });
      });

      await test.step("Gateway accepts both keys during the overlap", async () => {
        await gateway.assertResponds(ctx, {
          status: 200,
          headers: { "X-Gravitee-Api-Key": KEY_A },
        });
        await gateway.assertResponds(ctx, {
          status: 200,
          headers: { "X-Gravitee-Api-Key": KEY_B },
        });
      });

      await test.step("Remove KEY_A; APIM revokes KEY_A and keeps KEY_B active", async () => {
        await terraform.writeVars(ws, {
          hrid_suffix: "rot-gradual",
          keys: tfKeys([{ key: KEY_B }]),
        });
        await terraform.apply(ws);
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

      await test.step("Gateway rejects the deprecated key, accepts the new one", async () => {
        await gateway.assertNotResponds(ctx, {
          notStatus: 200,
          headers: { "X-Gravitee-Api-Key": KEY_A },
        });
        await gateway.assertResponds(ctx, {
          status: 200,
          headers: { "X-Gravitee-Api-Key": KEY_B },
        });
      });
    } finally {
      await terraform.destroyWorkspace(ws);
    }
  });

  // ── Reactivation: a revoked key reappears in the spec ─────────────

  test(`Previously revoked api-key is reactivated when re-added to spec ${XRAY.TERRAFORM.APIKEY_REACTIVATION} ${TAGS.REGRESSION}`, async ({
    mapi,
    gateway,
  }) => {
    test.setTimeout(240_000);

    const ws = await terraform.initWorkspace("terraform-apikey-custom");
    try {
      const KEY_A = uniqueKey("tf-reactivation-A");
      const KEY_B = uniqueKey("tf-reactivation-B");

      await terraform.writeVars(ws, {
        hrid_suffix: "reactivation",
        keys: tfKeys([{ key: KEY_A }]),
      });
      await terraform.apply(ws);

      const apiId = await terraform.output(ws, "api_id");
      const subId = await terraform.output(ws, "sub_id");
      const ctx = await terraform.output(ws, "api_context_path");

      await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);

      await test.step("Replace KEY_A with KEY_B; KEY_A becomes revoked", async () => {
        await terraform.writeVars(ws, {
          hrid_suffix: "reactivation",
          keys: tfKeys([{ key: KEY_B }]),
        });
        await terraform.apply(ws);
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
        await terraform.writeVars(ws, {
          hrid_suffix: "reactivation",
          keys: tfKeys([{ key: KEY_A }, { key: KEY_B }]),
        });
        await terraform.apply(ws);
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
        await gateway.assertResponds(ctx, {
          status: 200,
          headers: { "X-Gravitee-Api-Key": KEY_A },
        });
      });
    } finally {
      await terraform.destroyWorkspace(ws);
    }
  });

  // ── Staggered expirations: zero-downtime rotation ─────────────────

  test(`Multi-key subscription with staggered expirations supports zero-downtime rotation ${XRAY.TERRAFORM.APIKEY_STAGGERED_EXPIRY} ${TAGS.REGRESSION}`, async ({
    mapi,
    gateway,
  }) => {
    test.setTimeout(240_000);

    const ws = await terraform.initWorkspace("terraform-apikey-custom");
    try {
      // Three keys mirroring the canonical staged-rotation pattern from
      // GKO-2550: an early-expiring key on its way out, a longer-lived key,
      // and an evergreen key with no expire_at set.
      const KEY_V1 = uniqueKey("tf-stagger-v1");
      const KEY_V2 = uniqueKey("tf-stagger-v2");
      const KEY_V3 = uniqueKey("tf-stagger-v3");
      const expireV1 = new Date(Date.now() + 30 * 60 * 1_000).toISOString();
      const expireV2 = new Date(Date.now() + 90 * 60 * 1_000).toISOString();

      await terraform.writeVars(ws, {
        hrid_suffix: "staggered",
        keys: tfKeys([
          { key: KEY_V1, expire_at: expireV1 },
          { key: KEY_V2, expire_at: expireV2 },
          { key: KEY_V3 },
        ]),
      });
      await terraform.apply(ws);

      const apiId = await terraform.output(ws, "api_id");
      const subId = await terraform.output(ws, "sub_id");
      const ctx = await terraform.output(ws, "api_context_path");

      await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 3, { timeoutMs: 30_000 });

      await test.step("APIM stores each key with the correct expireAt (or none)", async () => {
        const { data } = await mapi.listSubscriptionApiKeys(apiId, subId);
        const v1 = data.find((k) => k.key === KEY_V1);
        const v2 = data.find((k) => k.key === KEY_V2);
        const v3 = data.find((k) => k.key === KEY_V3);
        expect(v1, "v1 must be present in APIM").toBeDefined();
        expect(v2, "v2 must be present in APIM").toBeDefined();
        expect(v3, "v3 must be present in APIM").toBeDefined();

        const v1Ms = new Date(v1!.expireAt as unknown as string | number).getTime();
        const v2Ms = new Date(v2!.expireAt as unknown as string | number).getTime();
        expect(Math.abs(v1Ms - new Date(expireV1).getTime())).toBeLessThan(2_000);
        expect(Math.abs(v2Ms - new Date(expireV2).getTime())).toBeLessThan(2_000);
        // v3 has no expire_at in the spec — APIM must not synthesize one.
        expect(v3!.expireAt).toBeFalsy();
      });

      await test.step("Gateway accepts all three keys during the overlap window", async () => {
        await gateway.assertResponds(ctx, {
          status: 200,
          headers: { "X-Gravitee-Api-Key": KEY_V1 },
        });
        await gateway.assertResponds(ctx, {
          status: 200,
          headers: { "X-Gravitee-Api-Key": KEY_V2 },
        });
        await gateway.assertResponds(ctx, {
          status: 200,
          headers: { "X-Gravitee-Api-Key": KEY_V3 },
        });
      });

      // Drop the soon-to-expire v1 ahead of its expiry — the standard
      // "rotate out the oldest key without service interruption" move.
      await test.step("Drop KEY_V1; APIM revokes it while v2 and v3 stay active", async () => {
        await terraform.writeVars(ws, {
          hrid_suffix: "staggered",
          keys: tfKeys([
            { key: KEY_V2, expire_at: expireV2 },
            { key: KEY_V3 },
          ]),
        });
        await terraform.apply(ws);
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
        await gateway.assertNotResponds(ctx, {
          notStatus: 200,
          headers: { "X-Gravitee-Api-Key": KEY_V1 },
        });
        await gateway.assertResponds(ctx, {
          status: 200,
          headers: { "X-Gravitee-Api-Key": KEY_V2 },
        });
        await gateway.assertResponds(ctx, {
          status: 200,
          headers: { "X-Gravitee-Api-Key": KEY_V3 },
        });
      });
    } finally {
      await terraform.destroyWorkspace(ws);
    }
  });

  // ── terraform destroy revokes the key ─────────────────────────────

  test(`Api key is revoked when terraform destroy removes the subscription ${XRAY.TERRAFORM.APIKEY_REVOKED_ON_DESTROY} ${TAGS.REGRESSION}`, async ({
    mapi,
    gateway,
  }) => {
    test.setTimeout(120_000);

    let ws: TfWorkspace | null = null;
    try {
      ws = await terraform.initWorkspace("terraform-apikey-custom");
      const CUSTOM_KEY = uniqueKey("tf-destroy-apikey");
      await terraform.writeVars(ws, {
        hrid_suffix: "revoked-on-destroy",
        keys: tfKeys([{ key: CUSTOM_KEY }]),
      });
      await terraform.apply(ws);

      const apiId = await terraform.output(ws, "api_id");
      const subId = await terraform.output(ws, "sub_id");
      const ctx = await terraform.output(ws, "api_context_path");

      await mapi.waitForSubscriptionApiKeyCount(apiId, subId, 1);

      // Confirm the key works before destroy.
      await gateway.assertResponds(ctx, {
        status: 200,
        headers: { "X-Gravitee-Api-Key": CUSTOM_KEY },
      });

      await terraform.destroy(ws);

      // After destroy, the previously valid key must stop working. The
      // gateway error code depends on policy (401/403/404), so
      // assertNotResponds is the right fit.
      await gateway.assertNotResponds(ctx, {
        notStatus: 200,
        headers: { "X-Gravitee-Api-Key": CUSTOM_KEY },
      });
    } finally {
      if (ws) await terraform.destroyWorkspace(ws);
    }
  });

  // ── TF-specific: plan is clean immediately after apply ────────────
  //
  // This catches drift where the provider round-trips api_keys via APIM
  // and gets back a value that doesn't match the desired state — leading
  // to perpetual "1 to change" on subsequent plans. There's no GKO
  // equivalent for this since k8s reconciliation is push-based, not
  // declarative-diff-based.

  test(`terraform plan reports no drift immediately after apply ${XRAY.TERRAFORM.APIKEY_PLAN_NO_DRIFT} ${TAGS.REGRESSION}`, async () => {
    test.setTimeout(180_000);

    const ws = await terraform.initWorkspace("terraform-apikey-custom");
    try {
      const KEY_A = uniqueKey("tf-no-drift-A");
      const KEY_B = uniqueKey("tf-no-drift-B");
      const expireB = new Date(Date.now() + 30 * 60 * 1_000).toISOString();
      await terraform.writeVars(ws, {
        hrid_suffix: "no-drift",
        keys: tfKeys([
          { key: KEY_A },
          { key: KEY_B, expire_at: expireB },
        ]),
      });
      await terraform.apply(ws);

      const { hasChanges, stdout } = await terraform.plan(ws);
      expect(hasChanges, `terraform plan reported drift:\n${stdout}`).toBe(false);

      // Re-applying the identical spec must also be a no-op.
      const applyOut = await terraform.apply(ws);
      // The provider should report 0 added/changed/destroyed for an
      // unchanged spec. Wording: "Apply complete! Resources: 0 added,
      // 0 changed, 0 destroyed."
      expect(applyOut).toMatch(/0 added.*0 changed.*0 destroyed/);
    } finally {
      await terraform.destroyWorkspace(ws);
    }
  });

  // ── TF-specific: key values are sensitive in plan output ──────────
  //
  // The TF provider schema marks `key` as Sensitive. Verify that a user
  // running `terraform plan` does not see their custom api-key values
  // printed in stdout — a regression in the Sensitive marking would leak
  // the keys to CI logs, ticket attachments, screenshots, etc.

  test(`Custom api-key values are redacted as sensitive in terraform plan output ${XRAY.TERRAFORM.APIKEY_SENSITIVE_IN_PLAN} ${TAGS.REGRESSION}`, async () => {
    test.setTimeout(120_000);

    const ws = await terraform.initWorkspace("terraform-apikey-custom");
    try {
      const SECRET_KEY = uniqueKey("tf-sensitive-leak-detect");
      await terraform.writeVars(ws, {
        hrid_suffix: "sensitive",
        keys: tfKeys([{ key: SECRET_KEY }]),
      });

      // Plan before apply: the create plan should also redact the value.
      const { stdout: preApplyPlan } = await terraform.plan(ws);
      expect(
        preApplyPlan.includes(SECRET_KEY),
        "pre-apply plan output leaked the api-key value",
      ).toBe(false);
      expect(preApplyPlan.toLowerCase()).toContain("sensitive");

      await terraform.apply(ws);

      // Post-apply plan: should be clean (no-drift), but if drift were
      // detected the key value still must not appear in stdout.
      const { stdout: postApplyPlan } = await terraform.plan(ws);
      expect(
        postApplyPlan.includes(SECRET_KEY),
        "post-apply plan output leaked the api-key value",
      ).toBe(false);
    } finally {
      await terraform.destroyWorkspace(ws);
    }
  });
});
