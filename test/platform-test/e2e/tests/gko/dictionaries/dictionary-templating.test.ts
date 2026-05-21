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
 * Dictionary Templating test.
 *
 * Verifies that a dynamic dictionary can pull its provider URL and header
 * value from a Kubernetes Secret using the [[ secret `name/key` ]] notation.
 *
 * Xray tests:
 *   GKO-TBD-DICT-DYNAMIC-TPL: Dynamic dictionary with templated secret values
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { poll, loadGraviteeConfig } from "../../../../src/index.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));

const DICT_NAME = "e2e-dict-dyn-tpl";
const API_NAME = "e2e-api-with-dyn-dict-tpl";
const API_PATH = "/e2e-api-with-dyn-dict-tpl";

async function gatewayBaseUrl(): Promise<string> {
  const config = await loadGraviteeConfig(path.resolve(__dirname, "../../../../config.yaml"));
  return config.gateway?.baseUrl ?? "http://localhost:30082";
}

test.describe("Dictionaries — Templating", () => {
  // FIXME(GKO-2858): a secret-templated Dictionary cannot be deleted — the Dictionary
  // controller runs template.Compile before the IsBeingDeleted() check, so the finalizer
  // is never removed and cleanup hangs in Terminating. Un-fixme once GKO-2858 is fixed.
  // https://gravitee.atlassian.net/browse/GKO-2858
  test.fixme(`Dynamic dictionary with secret templates resolves in API header ${XRAY.DICTIONARIES.DYNAMIC_TEMPLATE_RESOLVE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const secretFixture = fixture("crds/dictionaries/dictionary-dynamic-tpl-secret.yaml");
    const dictFixture = fixture("crds/dictionaries/dictionary-dynamic-tpl.yaml");
    const apiFixture = fixture("crds/dictionaries/api-with-dynamic-dictionary-tpl.yaml");

    await test.step("Apply secret with provider URL and header value", async () => {
      await kubectl.apply(secretFixture);
    });

    await test.step("Apply dynamic dictionary using secret templates", async () => {
      await kubectl.apply(dictFixture);
      await kubectl.waitForCondition("dictionary", DICT_NAME, "Accepted");
    });

    await test.step("Dictionary has an ID in status", async () => {
      const status = await kubectl.getStatus<{ id: string }>("dictionary", DICT_NAME);
      expect(status.id).toBeTruthy();
    });

    await test.step("Apply API that uses the templated dictionary in a header", async () => {
      await kubectl.apply(apiFixture);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("Gateway resolves templated dictionary value in echo response", async () => {
      const baseUrl = await gatewayBaseUrl();
      await poll(
        async () => {
          const res = await fetch(`${baseUrl}${API_PATH}`);
          expect(res.status).toBe(200);
          const body = (await res.json()) as { headers?: Record<string, string> };
          expect(body.headers).toBeDefined();
          expect(body.headers!["X-Env"] ?? body.headers!["x-env"]).toBe("ABCDEF");
        },
        { timeoutMs: 60_000, intervalMs: 3_000, description: "templated dynamic dictionary header resolved" },
      );
    });

    await kubectl.del(apiFixture);
    await kubectl.del(dictFixture);
    await kubectl.del(secretFixture);
  });
});
