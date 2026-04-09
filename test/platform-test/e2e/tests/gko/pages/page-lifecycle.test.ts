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
 *   GKO-277:  Add/remove/update inline documentation pages to V4 API CRD
 *   GKO-278:  Update inline documentation pages in V4 API CRD
 *   GKO-279:  Add documentation page using a fetcher in V4 API CRD
 *   GKO-1933: Autofetched content is not incorrectly deleted after operator sync
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import YAML from "yaml";
import { test, expect, fixture } from "../../../setup.js";
import { XRAY } from "../../../helpers/tags.js";

test.describe("Page Lifecycle", () => {
  test(`Markdown page create, update, delete ${XRAY.PAGES.MARKDOWN_PAGE_CRUD_V4} ${XRAY.PAGES.MARKDOWN_PAGE_UPDATE_V4}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-markdown-page";
    const WITH_PAGE = fixture("crds/pages/v4-api-with-page-markdown.yaml");
    const UPDATED_PAGE = fixture("crds/pages/v4-api-with-updated-page-markdown.yaml");
    const WITHOUT_PAGE = fixture("crds/pages/v4-api-without-page-markdown.yaml");

    await test.step("Deploy API with markdown page", async () => {
      await kubectl.apply(WITH_PAGE);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("Exported CRD contains markdown page", async () => {
      const crdYaml = await mapi.exportApiCrd(apiId);
      const crd = YAML.parse(crdYaml);
      const pages = crd?.spec?.pages;
      expect(pages?.["markdown-page"]).toBeDefined();
      expect(pages["markdown-page"].type).toBe("MARKDOWN");
      expect(pages["markdown-page"].name).toBe("markdown-page");
      expect(pages["markdown-page"].parentId).toBeDefined();
      expect(pages["markdown-folder"]?.type).toBe("FOLDER");
      expect(pages["markdown-folder"]?.name).toBe("markdown-folder");
    });

    await test.step("Update markdown page content", async () => {
      await kubectl.apply(UPDATED_PAGE);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("Exported CRD reflects content update", async () => {
      const crdYaml = await mapi.exportApiCrd(apiId);
      const crd = YAML.parse(crdYaml);
      const content = crd?.spec?.pages?.["markdown-page"]?.content;
      expect(content).toBeDefined();
      expect(content).toContain("This is an update");
    });

    await test.step("Remove markdown page", async () => {
      await kubectl.apply(WITHOUT_PAGE);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("Exported CRD no longer contains markdown page but folder remains", async () => {
      const crdYaml = await mapi.exportApiCrd(apiId);
      const crd = YAML.parse(crdYaml);
      const pages = crd?.spec?.pages;
      expect(pages?.["markdown-page"]).toBeUndefined();
      expect(pages?.["markdown-folder"]?.type).toBe("FOLDER");
      expect(pages?.["markdown-folder"]?.name).toBe("markdown-folder");
    });

    await kubectl.del(WITHOUT_PAGE);
  });

  test(`Swagger HTTP fetcher page ${XRAY.PAGES.FETCHER_PAGE_V4} ${XRAY.PAGES.AUTOFETCH_PRESERVED}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-swagger-http-fetcher";
    const FIXTURE = fixture("crds/pages/v4-api-with-swagger-http-fetcher.yaml");

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
