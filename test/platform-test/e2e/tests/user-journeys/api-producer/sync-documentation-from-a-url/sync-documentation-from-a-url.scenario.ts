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
 * Journey: sync API documentation from a URL.
 *
 * As an API producer, I keep my OpenAPI spec at a URL rather than pasting it
 * into the API definition: I declare a page whose content comes from an
 * http-fetcher, APIM pulls the spec in and refreshes it on a schedule, and the
 * fetched content survives later edits to the page.
 *
 * The fetcher is an inline attribute of the page on both drivers
 * (spec.pages.<hrid>.source / apim_apiv4.pages[].source), so neither driver
 * needs a standalone page resource. Fixtures declare no content at all — every
 * character of it comes from the fetch.
 *
 * github-fetcher stays GKO-only: admission pre-fetches the repository at apply
 * time and the test cluster has no GitHub credentials. Fetcher rejection cases
 * (a web fetcher with no URL, an invalid cron) and all V2 documentation live
 * under tests/gko/pages.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test, expect } from "../../../../setup.js";
import { XRAY, TAGS } from "../../../../helpers/tags.js";
import { forEachProvisioner } from "../../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../../helpers/provisioner-env.js";
import type { Page } from "../../../../../src/types/apim.js";

const here = path.dirname(fileURLToPath(import.meta.url));

const SPEC_URL = "https://petstore.swagger.io/v2/swagger.json";

/** The page as both fixtures declare it, before the rename. */
const FETCHER_PAGE = {
  name: "API reference",
  type: "SWAGGER",
  published: true,
  visibility: "PUBLIC",
  source: {
    type: "http-fetcher",
    configuration: {
      url: SPEC_URL,
      autoFetch: true,
      fetchCron: "*/10 * * * * *",
      useSystemProxy: false,
    },
  },
};

/** Project a fetched page to the declared fields under test. */
function project(p: Page) {
  return {
    name: p.name,
    type: p.type,
    published: p.published ?? false,
    visibility: p.visibility,
    source: p.source,
  };
}

/**
 * Shape of the fetched content, so the assertion holds whatever the upstream
 * spec says: it is valid JSON, it names itself, and it describes some paths.
 * Returns null while the content is still absent or not yet a full document.
 */
function fetchedSpec(p: Page | undefined): { title: string; pathCount: number } | null {
  if (!p?.content) return null;
  try {
    const spec = JSON.parse(p.content) as { info?: { title?: string }; paths?: object };
    const title = spec.info?.title ?? "";
    if (!title || !spec.paths) return null;
    return { title, pathCount: Object.keys(spec.paths).length };
  } catch {
    return null;
  }
}

/** Which revision of the fetched documentation the API carries, if any. */
interface FetcherParams {
  withPage: boolean;
  renamed: boolean;
}

forEachProvisioner<FetcherParams>(
  {
    title: "Document a V4 API with a page fetched from a URL",
    provisioners: {
      gko: gkoScenario<FetcherParams>({
        manifests: [],
        roles: { api: "fetched-docs-api" },
        dynamicRoles: ["api"],
        applyParams: async (k, params) => {
          const manifest = !params.withPage
            ? "gko/api-without-page.yaml"
            : params.renamed
              ? "gko/api-with-renamed-fetcher-page.yaml"
              : "gko/api-with-fetcher-page.yaml";
          await k.apply(path.join(here, manifest));
        },
      }),
      terraform: tfScenario<FetcherParams>({
        fixture: path.join(here, "terraform"),
        toVars: (params) => ({ with_page: params.withPage, page_renamed: params.renamed }),
      }),
    },
    xray: {
      gko: [XRAY.PAGES.FETCHER_PAGE_V4, XRAY.PAGES.AUTOFETCH_PRESERVED],
      terraform: XRAY.TERRAFORM.API_PAGE_FETCHER_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 60_000 },
  },
  async ({ provisioned, mapi }) => {
    const apiId = await provisioned.apiId();

    await test.step("Page declaring an http-fetcher lands in APIM", async () => {
      await expect
        .poll(async () => (await mapi.listApiPages(apiId)).map(project), {
          timeout: 30_000,
          message: "fetcher page reaches APIM",
        })
        .toEqual([FETCHER_PAGE]);
    });

    await test.step("APIM fetches the spec into the page content", async () => {
      await expect
        .poll(async () => fetchedSpec((await mapi.listApiPages(apiId))[0]), {
          timeout: 30_000,
          message: "spec content is fetched from the URL",
        })
        .toMatchObject({ title: expect.any(String), pathCount: expect.any(Number) });

      const spec = fetchedSpec((await mapi.listApiPages(apiId))[0]);
      expect(spec?.pathCount).toBeGreaterThan(0);
    });

    await test.step("Renaming the page keeps the fetched content", async () => {
      await provisioned.update({ withPage: true, renamed: true });
      // The whole page list plus the content, so a rename that replaced the page
      // (losing what was fetched) fails here rather than passing on the name alone.
      await expect
        .poll(
          async () => {
            const pages = await mapi.listApiPages(apiId);
            return { count: pages.length, name: pages[0]?.name, spec: fetchedSpec(pages[0]) };
          },
          { timeout: 30_000, message: "renamed page keeps its fetched content" },
        )
        .toMatchObject({
          count: 1,
          name: "API reference (v2)",
          spec: { title: expect.any(String) },
        });
    });

    await test.step("Stripping the page removes it in APIM", async () => {
      await provisioned.update({ withPage: false, renamed: false });
      await expect
        .poll(async () => (await mapi.listApiPages(apiId)).length, {
          timeout: 30_000,
          message: "fetched documentation page removed",
        })
        .toBe(0);
    });
  },
  { withPage: true, renamed: false },
);
