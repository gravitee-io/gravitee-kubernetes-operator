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
 * Journey: operate a DYNAMIC dictionary through its lifecycle.
 *
 * As an API producer, I run a DYNAMIC dictionary whose HTTP provider polls an
 * upstream and whose values an API resolves at the gateway. The same intent runs
 * against both provisioners: a dictionary created through any provisioner is
 * started, its value resolvable via `{#dictionaries['<hrid>']['X-Test-Specific']}`
 * in a transform-headers flow, and its lifecycle transitions (provider update,
 * deployed=false, delete) reach the gateway.
 *
 * Fixtures are co-located: `gko/api.yaml` (static consumer API) + `params.ts`
 * (the DYNAMIC Dictionary CR rendered from params, applied via applyParams so an
 * update re-applies a changed CR) for GKO; `terraform/main.tf` for Terraform.
 *
 * Provisioner-specific dictionary behaviour that has no cross-provisioner
 * meaning (GKO admission validation, GKO secret-templating, plain CR delete)
 * stays in the per-provisioner suite under tests/gko/dictionaries.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test, expect } from "../../../setup.js";
import { loadGraviteeConfig, poll } from "../../../../src/index.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import { forEachProvisioner } from "../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../helpers/provisioner-env.js";
import {
  dictionaryYaml,
  tfDynamicDictVars,
  GKO_DICT_NAME,
  INITIAL,
  type DynamicDictParams,
} from "./params.js";

const here = path.dirname(fileURLToPath(import.meta.url));

async function gatewayBaseUrl(): Promise<string> {
  const config = await loadGraviteeConfig(path.resolve(here, "../../../../config.yaml"));
  return config.gateway?.baseUrl ?? "http://localhost:30082";
}

function headerValue(headers: Record<string, string> | undefined, name: string): string | undefined {
  if (!headers) return undefined;
  // HTTP headers are case-insensitive; the echo endpoint may return any casing.
  const target = name.toLowerCase();
  const key = Object.keys(headers).find((k) => k.toLowerCase() === target);
  return key ? headers[key] : undefined;
}

/** GKO factory: static consumer API + the DYNAMIC dictionary applied from params. */
function dynDictGko() {
  return gkoScenario<DynamicDictParams>({
    manifests: [path.join(here, "gko/api.yaml")],
    roles: {
      api: "dyn-dictionary-consumer-api",
      dictionary: { kind: "dictionary", name: GKO_DICT_NAME },
    },
    dynamicRoles: ["dictionary"],
    contextPath: "/dyn-dictionary-consumer-api",
    applyParams: async (k, p) => {
      await k.applyString(dictionaryYaml(p));
    },
  });
}

/** Terraform factory: the co-located dynamic-dictionary fixture, params -> tfvars. */
function dynDictTf() {
  return tfScenario<DynamicDictParams>({
    fixture: path.join(here, "terraform"),
    toVars: tfDynamicDictVars,
    // remove("dictionary") drops the count-gated dictionary from desired state.
    removeVars: { dictionary: { create_dictionary: false } },
  });
}

/** Poll until the gateway resolves the dictionary value into the X-Env header. */
async function assertResolves(baseUrl: string, ctx: string, expected: string): Promise<void> {
  await poll(
    async () => {
      const res = await fetch(`${baseUrl}${ctx}`);
      if (res.status !== 200) {
        const errorBody = await res.text().catch(() => "<no body>");
        throw new Error(`gateway returned ${res.status} for ${ctx}: ${errorBody}`);
      }
      const body = (await res.json()) as { headers?: Record<string, string> };
      expect(headerValue(body.headers, "X-Env")).toBe(expected);
    },
    { timeoutMs: 60_000, intervalMs: 3_000, description: `dictionary resolves "${expected}"` },
  );
}

/**
 * Poll until the gateway no longer resolves `previous`. Once the dictionary is
 * gone (deleted) or stopped (deployed=false), the EL expression can no longer
 * find the entry: depending on policy this surfaces as a 500 (EL eval throws) or
 * a 200 with the value missing/changed. Both prove the previously-served value
 * is gone. Poll across ≥ 2 trigger cycles so we catch "stopped", not "paused".
 */
async function assertStopsResolving(baseUrl: string, ctx: string, previous: string): Promise<void> {
  await poll(
    async () => {
      const res = await fetch(`${baseUrl}${ctx}`);
      if (res.status === 500) return;
      if (res.status !== 200) {
        const errorBody = await res.text().catch(() => "<no body>");
        throw new Error(`gateway returned ${res.status} for ${ctx}: ${errorBody}`);
      }
      const body = (await res.json()) as { headers?: Record<string, string> };
      expect(headerValue(body.headers, "X-Env")).not.toBe(previous);
    },
    { timeoutMs: 30_000, intervalMs: 3_000, description: `dictionary no longer resolves "${previous}"` },
  );
}

const SINCE = { gko: "4.12", terraform: "4.12" } as const;

// ── Dynamic dictionary value resolves at the gateway ──────────────────────────
forEachProvisioner<DynamicDictParams>(
  {
    title: "Dynamic dictionary value resolves at the gateway",
    provisioners: { gko: dynDictGko(), terraform: dynDictTf() },
    xray: {
      gko: XRAY.DICTIONARIES.DYNAMIC_RESOLVE,
      terraform: XRAY.TERRAFORM.DICTIONARY_DYNAMIC_RESOLVE,
    },
    tags: [TAGS.REGRESSION],
    since: SINCE,
    timeoutMs: { gko: 120_000 },
  },
  async ({ provisioned }) => {
    const baseUrl = await gatewayBaseUrl();
    const ctx = await provisioned.contextPath();
    await assertResolves(baseUrl, ctx, "ABCDEF");
  },
  INITIAL,
);

// ── Updating the provider config propagates to the gateway ────────────────────
forEachProvisioner<DynamicDictParams>(
  {
    title: "Updating a dynamic dictionary propagates the new value to the gateway",
    provisioners: { gko: dynDictGko(), terraform: dynDictTf() },
    xray: {
      gko: XRAY.DICTIONARIES.DYNAMIC_UPDATE_PROPAGATES,
      terraform: XRAY.TERRAFORM.DICTIONARY_DYNAMIC_UPDATE,
    },
    tags: [TAGS.REGRESSION],
    since: SINCE,
    timeoutMs: { gko: 150_000 },
  },
  async ({ provisioned }) => {
    const baseUrl = await gatewayBaseUrl();
    const ctx = await provisioned.contextPath();
    await assertResolves(baseUrl, ctx, "ABCDEF");

    await test.step("Change the provider header value and re-apply", async () => {
      await provisioned.update({ headerValue: "ZYXWVU", deployed: true });
      await assertResolves(baseUrl, ctx, "ZYXWVU");
    });
  },
  INITIAL,
);

// ── Setting deployed=false stops gateway resolution ───────────────────────────
forEachProvisioner<DynamicDictParams>(
  {
    title: "Setting deployed=false on a dynamic dictionary stops gateway resolution",
    provisioners: { gko: dynDictGko(), terraform: dynDictTf() },
    xray: {
      gko: XRAY.DICTIONARIES.DYNAMIC_DEPLOYED_FALSE_STOPS,
      terraform: XRAY.TERRAFORM.DICTIONARY_DYNAMIC_DEPLOYED_FALSE_STOPS,
    },
    tags: [TAGS.REGRESSION],
    since: SINCE,
    timeoutMs: { gko: 150_000 },
  },
  async ({ provisioned }) => {
    const baseUrl = await gatewayBaseUrl();
    const ctx = await provisioned.contextPath();
    await assertResolves(baseUrl, ctx, "ABCDEF");

    await test.step("Flip deployed=false; the gateway stops resolving", async () => {
      await provisioned.update({ headerValue: "ABCDEF", deployed: false });
      await assertStopsResolving(baseUrl, ctx, "ABCDEF");
    });
  },
  INITIAL,
);

// ── Deleting the dictionary stops gateway resolution ──────────────────────────
forEachProvisioner<DynamicDictParams>(
  {
    title: "Deleting a dynamic dictionary stops gateway resolution",
    provisioners: { gko: dynDictGko(), terraform: dynDictTf() },
    xray: {
      gko: XRAY.DICTIONARIES.DYNAMIC_DELETE_STOPS,
      terraform: XRAY.TERRAFORM.DICTIONARY_DYNAMIC_DELETE_STOPS,
    },
    tags: [TAGS.REGRESSION],
    since: SINCE,
    timeoutMs: { gko: 150_000 },
  },
  async ({ provisioned }) => {
    const baseUrl = await gatewayBaseUrl();
    const ctx = await provisioned.contextPath();
    await assertResolves(baseUrl, ctx, "ABCDEF");

    await test.step("Remove only the dictionary; the gateway stops resolving, the API stays up", async () => {
      await provisioned.remove("dictionary");
      await assertStopsResolving(baseUrl, ctx, "ABCDEF");
    });
  },
  INITIAL,
);
