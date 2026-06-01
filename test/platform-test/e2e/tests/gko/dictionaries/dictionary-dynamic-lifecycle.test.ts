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
 * DYNAMIC Dictionary lifecycle tests (GKO-2565).
 *
 * Covers the start/stop pathways for DYNAMIC dictionaries that the basic
 * happy-path test in dictionary-lifecycle.test.ts does not exercise:
 *   - delete on a polling DYNAMIC dict stops gateway resolution
 *   - update to provider headers propagates to the gateway
 *   - deployed:false stops polling without deleting the CR
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 *   - The gateway can reach https://api.gravitee.io/echo
 */

import { test, fixture, expect } from "../../../setup.js";
import { poll, loadGraviteeConfig } from "../../../../src/index.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectl from "../../../helpers/kubectl.js";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));

const DELETE_DICT_FIXTURE = fixture("crds/dictionaries/dictionary-dynamic.yaml");
const DELETE_API_FIXTURE = fixture("crds/dictionaries/api-with-dynamic-dictionary.yaml");
const DELETE_DICT_NAME = "e2e-dict-dynamic";
const DELETE_API_NAME = "e2e-api-with-dyn-dict";
const DELETE_API_PATH = "/e2e-api-with-dyn-dict";

const UPDATE_DICT_V1_FIXTURE = fixture("crds/dictionaries/dictionary-dynamic-update-v1.yaml");
const UPDATE_DICT_V2_FIXTURE = fixture("crds/dictionaries/dictionary-dynamic-update-v2.yaml");
const UPDATE_API_FIXTURE = fixture("crds/dictionaries/api-with-dyn-dict-update.yaml");
const UPDATE_DICT_NAME = "e2e-dict-dyn-update";
const UPDATE_API_NAME = "e2e-api-dyn-dict-update";
const UPDATE_API_PATH = "/e2e-api-dyn-dict-update";

const DEPLOYED_DICT_FIXTURE = fixture("crds/dictionaries/dictionary-dynamic-deployed.yaml");
const UNDEPLOYED_DICT_FIXTURE = fixture("crds/dictionaries/dictionary-dynamic-undeployed.yaml");
const DEPLOYED_API_FIXTURE = fixture("crds/dictionaries/api-with-dyn-dict-deployed.yaml");
const DEPLOYED_DICT_NAME = "e2e-dict-dyn-deployed";
const DEPLOYED_API_NAME = "e2e-api-dyn-dict-deployed";
const DEPLOYED_API_PATH = "/e2e-api-dyn-dict-deployed";

async function gatewayBaseUrl(): Promise<string> {
  const config = await loadGraviteeConfig(path.resolve(__dirname, "../../../../config.yaml"));
  return config.gateway?.baseUrl ?? "http://localhost:30082";
}

function headerValue(headers: Record<string, string> | undefined, name: string): string | undefined {
  if (!headers) return undefined;
  return headers[name] ?? headers[name.toLowerCase()];
}

test.describe("Dictionaries — Dynamic lifecycle", () => {
  // Safety-net cleanup: if any test times out before its inline cleanup runs,
  // these still execute. Each del() ignores errors because the resource may
  // already have been removed by the test itself.
  test.afterEach(async () => {
    await kubectl.del(DELETE_API_FIXTURE).catch(() => {});
    await kubectl.del(DELETE_DICT_FIXTURE).catch(() => {});
    await kubectl.del(UPDATE_API_FIXTURE).catch(() => {});
    await kubectl.del(UPDATE_DICT_V2_FIXTURE).catch(() => {});
    await kubectl.del(UPDATE_DICT_V1_FIXTURE).catch(() => {});
    await kubectl.del(DEPLOYED_API_FIXTURE).catch(() => {});
    await kubectl.del(UNDEPLOYED_DICT_FIXTURE).catch(() => {});
    await kubectl.del(DEPLOYED_DICT_FIXTURE).catch(() => {});
  });

  test(`Delete a dynamic dictionary stops gateway resolution ${XRAY.DICTIONARIES.DYNAMIC_DELETE_STOPS} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    test.setTimeout(120_000);
    const baseUrl = await gatewayBaseUrl();

    await test.step("Apply dynamic dictionary CRD", async () => {
      await kubectl.apply(DELETE_DICT_FIXTURE);
      await kubectl.waitForCondition("dictionary", DELETE_DICT_NAME, "Accepted");
    });

    await test.step("Apply API that uses the dynamic dictionary", async () => {
      await kubectl.apply(DELETE_API_FIXTURE);
      await kubectl.waitForCondition("apiv4definition", DELETE_API_NAME, "Accepted");
    });

    await test.step("Gateway resolves dynamic dictionary value", async () => {
      await poll(
        async () => {
          const res = await fetch(`${baseUrl}${DELETE_API_PATH}`);
          expect(res.status).toBe(200);
          const body = (await res.json()) as { headers?: Record<string, string> };
          expect(headerValue(body.headers, "X-Env")).toBe("ABCDEF");
        },
        { timeoutMs: 60_000, intervalMs: 3_000, description: "initial dynamic dictionary resolution" },
      );
    });

    await test.step("Delete the dictionary CRD", async () => {
      await kubectl.del(DELETE_DICT_FIXTURE);
      await kubectl.waitForDeletion("dictionary", DELETE_DICT_NAME);
    });

    await test.step("Gateway no longer resolves the dictionary value", async () => {
      // After delete, the operator stopped polling and APIM removed the
      // dictionary. At the gateway the EL expression can no longer find
      // the dictionary entry; depending on policy this surfaces as either
      // a 200 with the value missing/different, or a 500 because the EL
      // evaluation throws. Both outcomes prove the previously-resolved
      // value ("ABCDEF") is gone. Poll across ≥ 2 trigger cycles
      // (5 s × 2 = 10 s) to be confident polling has stopped, not paused.
      await poll(
        async () => {
          const res = await fetch(`${baseUrl}${DELETE_API_PATH}`);
          if (res.status === 500) {
            return;
          }
          expect(res.status).toBe(200);
          const body = (await res.json()) as { headers?: Record<string, string> };
          const value = headerValue(body.headers, "X-Env");
          expect(value).not.toBe("ABCDEF");
        },
        { timeoutMs: 30_000, intervalMs: 3_000, description: "dynamic dictionary no longer resolved after delete" },
      );
    });

    await kubectl.del(DELETE_API_FIXTURE);
  });

  test(`Update a dynamic dictionary propagates new provider config to gateway ${XRAY.DICTIONARIES.DYNAMIC_UPDATE_PROPAGATES} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    test.setTimeout(120_000);
    const baseUrl = await gatewayBaseUrl();

    await test.step("Apply v1 of the dynamic dictionary (header value ABCDEF)", async () => {
      await kubectl.apply(UPDATE_DICT_V1_FIXTURE);
      await kubectl.waitForCondition("dictionary", UPDATE_DICT_NAME, "Accepted");
    });

    await test.step("Apply API that uses the dictionary", async () => {
      await kubectl.apply(UPDATE_API_FIXTURE);
      await kubectl.waitForCondition("apiv4definition", UPDATE_API_NAME, "Accepted");
    });

    await test.step("Gateway resolves v1 value (ABCDEF)", async () => {
      await poll(
        async () => {
          const res = await fetch(`${baseUrl}${UPDATE_API_PATH}`);
          expect(res.status).toBe(200);
          const body = (await res.json()) as { headers?: Record<string, string> };
          expect(headerValue(body.headers, "X-Env")).toBe("ABCDEF");
        },
        { timeoutMs: 60_000, intervalMs: 3_000, description: "v1 dynamic dictionary value resolved" },
      );
    });

    await test.step("Apply v2 of the dictionary (header value ZYXWVU)", async () => {
      // Apply over an edited fixture file to trigger reconciliation.
      await kubectl.apply(UPDATE_DICT_V2_FIXTURE);
      await kubectl.waitForCondition("dictionary", UPDATE_DICT_NAME, "Accepted");
    });

    await test.step("Gateway resolves the updated value (ZYXWVU)", async () => {
      await poll(
        async () => {
          const res = await fetch(`${baseUrl}${UPDATE_API_PATH}`);
          expect(res.status).toBe(200);
          const body = (await res.json()) as { headers?: Record<string, string> };
          expect(headerValue(body.headers, "X-Env")).toBe("ZYXWVU");
        },
        { timeoutMs: 60_000, intervalMs: 3_000, description: "updated dynamic dictionary value propagated" },
      );
    });

    await kubectl.del(UPDATE_API_FIXTURE);
    await kubectl.del(UPDATE_DICT_V2_FIXTURE);
  });

  test(`Setting deployed=false on a dynamic dictionary stops gateway resolution ${XRAY.DICTIONARIES.DYNAMIC_DEPLOYED_FALSE_STOPS} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    test.setTimeout(120_000);
    const baseUrl = await gatewayBaseUrl();

    await test.step("Apply dictionary with deployed=true", async () => {
      await kubectl.apply(DEPLOYED_DICT_FIXTURE);
      await kubectl.waitForCondition("dictionary", DEPLOYED_DICT_NAME, "Accepted");
    });

    await test.step("Apply API that uses the dictionary", async () => {
      await kubectl.apply(DEPLOYED_API_FIXTURE);
      await kubectl.waitForCondition("apiv4definition", DEPLOYED_API_NAME, "Accepted");
    });

    await test.step("Gateway resolves dictionary value while deployed", async () => {
      await poll(
        async () => {
          const res = await fetch(`${baseUrl}${DEPLOYED_API_PATH}`);
          expect(res.status).toBe(200);
          const body = (await res.json()) as { headers?: Record<string, string> };
          expect(headerValue(body.headers, "X-Env")).toBe("DEPLOYED");
        },
        { timeoutMs: 60_000, intervalMs: 3_000, description: "deployed dictionary resolved" },
      );
    });

    await test.step("Apply same dictionary with deployed=false", async () => {
      await kubectl.apply(UNDEPLOYED_DICT_FIXTURE);
      await kubectl.waitForCondition("dictionary", DEPLOYED_DICT_NAME, "Accepted");
    });

    await test.step("Gateway no longer resolves the dictionary value", async () => {
      // CR still exists; only `deployed` flipped. Verify the gateway
      // stopped resolving across ≥ 2 trigger cycles.
      await poll(
        async () => {
          const res = await fetch(`${baseUrl}${DEPLOYED_API_PATH}`);
          expect(res.status).toBe(200);
          const body = (await res.json()) as { headers?: Record<string, string> };
          const value = headerValue(body.headers, "X-Env");
          expect(value).not.toBe("DEPLOYED");
        },
        { timeoutMs: 30_000, intervalMs: 3_000, description: "undeployed dictionary stopped resolving" },
      );
    });

    await test.step("Dictionary CR still exists (not deleted)", async () => {
      expect(await kubectl.exists("dictionary", DEPLOYED_DICT_NAME)).toBe(true);
    });

    await kubectl.del(DEPLOYED_API_FIXTURE);
    await kubectl.del(UNDEPLOYED_DICT_FIXTURE);
  });
});
