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
 * Dictionaries Lifecycle tests.
 *
 * Xray tests:
 *   GKO-2903: Create a manual dictionary and verify resolution in API response
 *   GKO-2904: Create a dynamic dictionary and verify resolution in API response
 *   GKO-2905: Delete a dictionary
 *   GKO-2912: Admission webhook rejects DYNAMIC dictionary with manual field set
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { poll, loadGraviteeConfig } from "../../../../src/index.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectl from "../../../helpers/kubectl.js";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));

const DICT_NAME = "e2e-dict-manual";
const API_NAME = "e2e-api-with-dict";
const API_PATH = "/e2e-api-with-dict";

async function gatewayBaseUrl(): Promise<string> {
  const config = await loadGraviteeConfig(path.resolve(__dirname, "../../../../config.yaml"));
  return config.gateway?.baseUrl ?? "http://localhost:30082";
}

test.describe("Dictionaries — Lifecycle", () => {
  // Safety-net cleanup: runs even if a test times out before its inline
  // cleanup. Each del() ignores errors (the resource may already be gone).
  test.afterEach(async () => {
    for (const f of [
      "crds/dictionaries/api-with-dictionary.yaml",
      "crds/dictionaries/api-with-dynamic-dictionary.yaml",
      "crds/dictionaries/dictionary-manual.yaml",
      "crds/dictionaries/dictionary-dynamic.yaml",
    ]) {
      await kubectl.del(fixture(f)).catch(() => {});
    }
  });

  test(`Create dictionary and resolve in API header ${XRAY.DICTIONARIES.CREATE_AND_RESOLVE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const dictFixture = fixture("crds/dictionaries/dictionary-manual.yaml");
    const apiFixture = fixture("crds/dictionaries/api-with-dictionary.yaml");

    await test.step("Apply dictionary CRD", async () => {
      await kubectl.apply(dictFixture);
      await kubectl.waitForCondition("dictionary", DICT_NAME, "Accepted");
    });

    await test.step("Dictionary has an ID in status", async () => {
      const status = await kubectl.getStatus<{ id: string }>("dictionary", DICT_NAME);
      expect(status.id).toBeTruthy();
    });

    await test.step("Apply API that uses the dictionary in a header", async () => {
      await kubectl.apply(apiFixture);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("Gateway resolves dictionary value in echo response", async () => {
      const baseUrl = await gatewayBaseUrl();
      await poll(
        async () => {
          const res = await fetch(`${baseUrl}${API_PATH}`);
          expect(res.status).toBe(200);
          const body = (await res.json()) as { headers?: Record<string, string> };
          expect(body.headers).toBeDefined();
          expect(body.headers!["X-Dict-Env"] ?? body.headers!["x-dict-env"]).toBe("test");
        },
        { timeoutMs: 30_000, intervalMs: 2_000, description: "dictionary header resolved" },
      );
    });

    await kubectl.del(apiFixture);
    await kubectl.del(dictFixture);
  });

  // ── GKO-2564: Dynamic dictionary — resolve echo header via JOLT ──

  test(`Create dynamic dictionary and resolve in API header ${XRAY.DICTIONARIES.DYNAMIC_RESOLVE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const dictFixture = fixture("crds/dictionaries/dictionary-dynamic.yaml");
    const apiFixture = fixture("crds/dictionaries/api-with-dynamic-dictionary.yaml");
    const DYN_DICT_NAME = "e2e-dict-dynamic";
    const DYN_API_NAME = "e2e-api-with-dyn-dict";
    const DYN_API_PATH = "/e2e-api-with-dyn-dict";

    await test.step("Apply dynamic dictionary CRD", async () => {
      await kubectl.apply(dictFixture);
      await kubectl.waitForCondition("dictionary", DYN_DICT_NAME, "Accepted");
    });

    await test.step("Dictionary has an ID in status", async () => {
      const status = await kubectl.getStatus<{ id: string }>("dictionary", DYN_DICT_NAME);
      expect(status.id).toBeTruthy();
    });

    await test.step("Apply API that uses the dynamic dictionary in a header", async () => {
      await kubectl.apply(apiFixture);
      await kubectl.waitForCondition("apiv4definition", DYN_API_NAME, "Accepted");
    });

    await test.step("Gateway resolves dynamic dictionary value in echo response", async () => {
      const baseUrl = await gatewayBaseUrl();
      // Dynamic dictionaries need time: trigger interval (5s) + gateway sync.
      await poll(
        async () => {
          const res = await fetch(`${baseUrl}${DYN_API_PATH}`);
          expect(res.status).toBe(200);
          const body = (await res.json()) as { headers?: Record<string, string> };
          expect(body.headers).toBeDefined();
          expect(body.headers!["X-Env"] ?? body.headers!["x-env"]).toBe("ABCDEF");
        },
        { timeoutMs: 60_000, intervalMs: 3_000, description: "dynamic dictionary header resolved" },
      );
    });

    await kubectl.del(apiFixture);
    await kubectl.del(dictFixture);
  });

  // ── GKO-2563: Delete a dictionary ──────────────────────────────

  test(`Delete a dictionary ${XRAY.DICTIONARIES.DELETE_DICTIONARY} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const dictFixture = fixture("crds/dictionaries/dictionary-manual.yaml");

    await test.step("Create dictionary", async () => {
      await kubectl.apply(dictFixture);
      await kubectl.waitForCondition("dictionary", DICT_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("dictionary", DICT_NAME);
    expect(status.id).toBeTruthy();

    await test.step("Delete the dictionary CRD", async () => {
      await kubectl.del(dictFixture);
      await kubectl.waitForDeletion("dictionary", DICT_NAME);
    });
  });

  // ── Admission: DYNAMIC + manual field is rejected ─────────────

  test(`Admission webhook rejects DYNAMIC dictionary with manual field set ${XRAY.DICTIONARIES.ADMISSION_REJECTS_DYNAMIC_WITH_MANUAL} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const invalidFixture = fixture("crds/dictionaries/dictionary-dynamic-invalid.yaml");

    const stderr = await kubectl.applyExpectFailure(invalidFixture);
    expect(stderr).toMatch(/dictionary type is DYNAMIC but 'manual' field is set/);
  });
});
