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
 * Journey: provision an API referencing a dictionary property.
 *
 * As an API producer, I provision an API that references a dictionary property,
 * resolved at the gateway. A MANUAL dictionary created through any provisioner is
 * deployed and resolvable at the gateway: an API whose transform-headers flow
 * injects `{#dictionaries['<hrid>']['env']}` sees the value reflected by the echo
 * endpoint. A second scenario proves an in-place property update propagates to
 * the gateway.
 *
 * Fixtures are co-located in this folder. DYNAMIC dictionaries and the
 * undeploy/delete lifecycle live in the manage-dynamic-dictionary journey; the
 * remaining provisioner-specific behaviour (GKO JOLT, delete, admission,
 * templating) stays in the per-provisioner suite under tests/gko/dictionaries.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test, expect } from "../../../setup.js";
import { loadGraviteeConfig, poll } from "../../../../src/index.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import { forEachProvisioner } from "../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../helpers/provisioner-env.js";
import {
  manualDictionaryYaml,
  tfManualDictVars,
  MANUAL_INITIAL,
  type ManualDictParams,
} from "./params.js";

const here = path.dirname(fileURLToPath(import.meta.url));

async function gatewayBaseUrl(): Promise<string> {
  const config = await loadGraviteeConfig(path.resolve(here, "../../../../config.yaml"));
  return config.gateway?.baseUrl ?? "http://localhost:30082";
}

/** Poll until the gateway resolves the dictionary value into the X-Dict-Env header. */
async function assertDictHeader(
  baseUrl: string,
  ctx: string,
  expected: string,
  timeoutMs = 30_000,
): Promise<void> {
  await poll(
    async () => {
      const res = await fetch(`${baseUrl}${ctx}`);
      if (res.status !== 200) {
        const errorBody = await res.text().catch(() => "<no body>");
        throw new Error(`gateway returned ${res.status} for ${ctx}: ${errorBody}`);
      }
      const body = (await res.json()) as { headers?: Record<string, string> };
      const value = body.headers?.["X-Dict-Env"] ?? body.headers?.["x-dict-env"];
      expect(value, "dictionary value resolved in the gateway response header").toBe(expected);
    },
    { timeoutMs, intervalMs: 2_000, description: `dictionary value "${expected}" resolves in the gateway response` },
  );
}

forEachProvisioner(
  {
    title: "Manual dictionary value resolves in the gateway response",
    provisioners: {
      gko: gkoScenario<void>({
        // Apply the dictionary first, then the API that references it by HRID.
        manifests: [path.join(here, "gko/dictionary.yaml"), path.join(here, "gko/api.yaml")],
        roles: {
          api: "dictionary-consumer-api",
          dictionary: { kind: "dictionary", name: "manual-dictionary" },
        },
        contextPath: "/dictionary-consumer-api",
      }),
      terraform: tfScenario<void>({ fixture: path.join(here, "terraform") }),
    },
    xray: {
      gko: XRAY.DICTIONARIES.CREATE_AND_RESOLVE,
      terraform: XRAY.TERRAFORM.DICTIONARY_MANUAL_RESOLVE,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 60_000 },
  },
  async ({ provisioned }) => {
    const baseUrl = await gatewayBaseUrl();
    const ctx = await provisioned.contextPath();
    // The shared invariant: the gateway resolves the dictionary value into the
    // header the echo endpoint reflects back.
    await assertDictHeader(baseUrl, ctx, "test");
  },
);

// ── Updating a manual dictionary property propagates to the gateway ────────────
forEachProvisioner<ManualDictParams>(
  {
    title: "Updating a manual dictionary property propagates to the gateway",
    provisioners: {
      // The API is static; the dictionary is rendered from params so update()
      // re-applies it (GKO) / re-applies with new tfvars (Terraform).
      gko: gkoScenario<ManualDictParams>({
        manifests: [path.join(here, "gko/api.yaml")],
        roles: {
          api: "dictionary-consumer-api",
          dictionary: { kind: "dictionary", name: "manual-dictionary" },
        },
        dynamicRoles: ["dictionary"],
        contextPath: "/dictionary-consumer-api",
        applyParams: async (k, p) => {
          await k.applyString(manualDictionaryYaml(p));
        },
      }),
      terraform: tfScenario<ManualDictParams>({
        fixture: path.join(here, "terraform"),
        toVars: tfManualDictVars,
      }),
    },
    xray: {
      gko: XRAY.DICTIONARIES.MANUAL_UPDATE_PROPAGATES,
      terraform: XRAY.TERRAFORM.DICTIONARY_MANUAL_UPDATE,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 90_000 },
  },
  async ({ provisioned }) => {
    const baseUrl = await gatewayBaseUrl();
    const ctx = await provisioned.contextPath();
    await assertDictHeader(baseUrl, ctx, "test");

    await test.step("Change the property value and re-apply", async () => {
      await provisioned.update({ envValue: "staging" });
      await assertDictHeader(baseUrl, ctx, "staging", 45_000);
    });
  },
  MANUAL_INITIAL,
);
