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
 * Manual dictionary resolution, shared across provisioners. A MANUAL dictionary
 * created through any provisioner must be deployed and resolvable at the gateway:
 * an API whose transform-headers flow injects `{#dictionaries['<hrid>']['env']}`
 * sees the value reflected by the echo endpoint.
 *
 * The shared invariant (gateway resolves the value to "test") is provisioner-
 * agnostic; each arm references its own dictionary HRID in its own fixtures
 * (GKO `default-e2e-dict-manual`, Terraform `e2e-tf-dict-resolve`). The
 * provisioner-specific dictionary behaviour (GKO dynamic/JOLT, delete, admission,
 * templating) stays in the per-provisioner suite under tests/gko/dictionaries.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { expect } from "../../../setup.js";
import { loadGraviteeConfig, poll } from "../../../../src/index.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import { forEachProvisioner } from "../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../helpers/provisioner-env.js";

const __dirname = path.dirname(fileURLToPath(import.meta.url));

async function gatewayBaseUrl(): Promise<string> {
  const config = await loadGraviteeConfig(path.resolve(__dirname, "../../../../config.yaml"));
  return config.gateway?.baseUrl ?? "http://localhost:30082";
}

forEachProvisioner(
  {
    title: "Manual dictionary value resolves in the gateway response",
    provisioners: {
      gko: gkoScenario<void>({
        // Apply the dictionary first, then the API that references it by HRID.
        manifests: [
          "dictionaries/dictionary-manual/crd.yaml",
          "dictionaries/api-with-dictionary/crd.yaml",
        ],
        roles: {
          api: "e2e-api-with-dict",
          // Full { kind, name } form: there is no shorthand convention for the
          // dictionary role, and we only need provision() to await it Accepted.
          dictionary: { kind: "dictionary", name: "e2e-dict-manual" },
        },
        contextPath: "/e2e-api-with-dict",
      }),
      terraform: tfScenario<void>({ fixture: "dictionaries/manual-resolve" }),
    },
    xray: {
      gko: XRAY.DICTIONARIES.CREATE_AND_RESOLVE,
      terraform: XRAY.TERRAFORM.DICTIONARY_MANUAL_RESOLVE,
    },
    tags: [TAGS.REGRESSION],
    // Dictionaries land via the Automation API and the v4 transform-headers flow
    // resolution path both ship in 4.12, for both provisioners.
    since: { gko: "4.12", terraform: "4.12" },
    // GKO provisions two CRs then polls the gateway; give it headroom over 30s.
    timeoutMs: { gko: 60_000 },
  },
  async ({ provisioned }) => {
    const baseUrl = await gatewayBaseUrl();
    const ctx = await provisioned.contextPath();

    // The shared invariant: the gateway resolves the dictionary value into the
    // header the echo endpoint reflects back. Poll for eventual consistency
    // (operator/provider apply -> APIM -> gateway sync). A non-200 throws with
    // the status + body so that on timeout the failure says WHY the gateway
    // never resolved, rather than a bare "received undefined". `poll` retries on
    // throw, so a transient non-200 during gateway sync is not a hard failure.
    await poll(
      async () => {
        const res = await fetch(`${baseUrl}${ctx}`);
        if (res.status !== 200) {
          const errorBody = await res.text().catch(() => "<no body>");
          throw new Error(`gateway returned ${res.status} for ${ctx}: ${errorBody}`);
        }
        const body = (await res.json()) as { headers?: Record<string, string> };
        const value = body.headers?.["X-Dict-Env"] ?? body.headers?.["x-dict-env"];
        expect(value, "dictionary value resolved in the gateway response header").toBe("test");
      },
      { timeoutMs: 30_000, intervalMs: 2_000, description: "dictionary value resolves in the gateway response" },
    );
  },
);
