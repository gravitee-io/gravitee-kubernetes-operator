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
 * Documentation page lifecycle tests (V4 only).
 *
 * Xray tests:
 *   GKO-279:  Add documentation page using a fetcher in V4 API CRD
 *   GKO-1933: Autofetched content is not incorrectly deleted after operator sync
 *
 * Adding, revising and removing an inline page (GKO-277/278) is the shared
 * journey tests/user-journeys/api-producer/document-an-api.
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import YAML from "yaml";
import { test, expect, fixture } from "../../../setup.js";
import { XRAY, PROVISIONER } from "../../../helpers/tags.js";
import * as kubectl from "../../../helpers/kubectl.js";

test.describe(`Page Lifecycle ${PROVISIONER.GKO}`, () => {
  // Safety-net cleanup: runs even if a test times out before its inline
  // cleanup. Each del() ignores errors (the resource may already be gone).
  test.afterEach(async () => {
    await kubectl.del(fixture("pages/v4-api-with-swagger-http-fetcher/crd.yaml")).catch(() => {});
  });
  test(`Swagger HTTP fetcher page ${XRAY.PAGES.FETCHER_PAGE_V4} ${XRAY.PAGES.AUTOFETCH_PRESERVED}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-swagger-http-fetcher";
    const FIXTURE = fixture("pages/v4-api-with-swagger-http-fetcher/crd.yaml");

    await test.step("Deploy API with swagger HTTP fetcher page", async () => {
      await kubectl.apply(FIXTURE);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("Exported CRD contains swagger page with HTTP fetcher source", async () => {
      const crdYaml = await mapi.exportApiCrd(apiId);
      const crd = YAML.parse(crdYaml);
      const swagger = crd?.spec?.pages?.swagger;
      expect(swagger).toBeDefined();
      expect(swagger.type).toBe("SWAGGER");
      expect(swagger.source?.type).toBe("http-fetcher");
      expect(swagger.source?.configuration?.url).toBe("https://petstore.swagger.io/v2/swagger.json");
      expect(swagger.source?.configuration?.autoFetch).toBe(true);
      expect(swagger.source?.configuration?.fetchCron).toBeDefined();
      expect(swagger.source?.configuration?.useSystemProxy).toBe(false);
    });

    await kubectl.del(FIXTURE);
  });
});
