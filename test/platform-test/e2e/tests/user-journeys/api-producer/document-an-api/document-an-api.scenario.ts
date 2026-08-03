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
 * Journey: document an API.
 *
 * As an API producer, I ship documentation alongside the API definition, revise
 * it as the API changes — renaming the page, rewriting it, narrowing who can
 * read it — and remove it when it is obsolete.
 *
 * Documentation is an inline attribute of the API on both drivers (spec.pages,
 * keyed by hrid / pages[] with an hrid) — there is no standalone Page Terraform
 * resource. The page keeps its hrid across the revision, so a rename has to
 * update the existing page rather than replace it.
 *
 * Fixtures are co-located in this folder. Page fetchers (pages[].source) are
 * expressible on both drivers and are a feasible follow-up; their rejection
 * cases (a web fetcher with no URL, an invalid cron) and all V2 documentation
 * (no apim_apiv2) stay GKO-only under tests/gko/pages.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test, expect } from "../../../../setup.js";
import { XRAY, TAGS } from "../../../../helpers/tags.js";
import { forEachProvisioner } from "../../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../../helpers/provisioner-env.js";
import type { Page } from "../../../../../src/types/apim.js";

const here = path.dirname(fileURLToPath(import.meta.url));

/**
 * The page as each fixture declares it. Content is compared trimmed because the
 * two drivers' block-scalar / heredoc syntaxes differ on the trailing newline.
 */
const INITIAL_PAGE = {
  name: "Getting started",
  type: "MARKDOWN",
  content: "# Getting started\n\nCall `GET /` to reach the upstream echo endpoint.",
  published: true,
  visibility: "PUBLIC",
};

/** The same page (same hrid) renamed, rewritten and made private. */
const REVISED_PAGE = {
  name: "Quick start",
  type: "MARKDOWN",
  content: "# Quick start\n\nSend any request to the base path; the upstream echoes it back.",
  published: true,
  visibility: "PRIVATE",
};

/** Project a fetched page to the fields under test (content normalised). */
function project(p: Page) {
  return {
    name: p.name,
    type: p.type,
    content: (p.content ?? "").trim(),
    published: p.published ?? false,
    visibility: p.visibility,
  };
}

/** Which revision of the documentation the API carries, if any. */
interface PageParams {
  withPage: boolean;
  revised: boolean;
}

forEachProvisioner<PageParams>(
  {
    title: "Document a V4 API with an inline markdown page",
    provisioners: {
      gko: gkoScenario<PageParams>({
        manifests: [],
        roles: { api: "documented-api" },
        dynamicRoles: ["api"],
        applyParams: async (k, params) => {
          const manifest = !params.withPage
            ? "gko/api-without-page.yaml"
            : params.revised
              ? "gko/api-with-revised-page.yaml"
              : "gko/api-with-page.yaml";
          await k.apply(path.join(here, manifest));
        },
      }),
      terraform: tfScenario<PageParams>({
        fixture: path.join(here, "terraform"),
        toVars: (params) => ({ with_page: params.withPage, page_revised: params.revised }),
      }),
    },
    xray: {
      gko: [
        XRAY.PAGES.V4_DOC_RECONCILED,
        XRAY.PAGES.MARKDOWN_PAGE_CRUD_V4,
        XRAY.PAGES.MARKDOWN_PAGE_UPDATE_V4,
        XRAY.PAGES.V4_DOC_OPERATIONS,
        XRAY.PAGES.V4_DOC_RENAME,
        XRAY.PAGES.V4_DOC_VISIBILITY_PUBLIC,
      ],
      terraform: XRAY.TERRAFORM.API_INLINE_PAGES_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 60_000 },
  },
  async ({ provisioned, mapi }) => {
    const apiId = await provisioned.apiId();

    await test.step("Page shipped with the API lands in APIM", async () => {
      await expect
        .poll(async () => (await mapi.listApiPages(apiId)).map(project), {
          timeout: 30_000,
          message: "API documentation page reaches APIM",
        })
        .toEqual([INITIAL_PAGE]);
    });

    await test.step("Revising the page renames and rewrites it in place", async () => {
      await provisioned.update({ withPage: true, revised: true });
      // The whole page list, so a revision that CREATED a second page instead of
      // updating the existing one fails here rather than passing on a lucky read.
      await expect
        .poll(async () => (await mapi.listApiPages(apiId)).map(project), {
          timeout: 30_000,
          message: "revised documentation page reaches APIM",
        })
        .toEqual([REVISED_PAGE]);
    });

    await test.step("Stripping the page removes it in APIM", async () => {
      await provisioned.update({ withPage: false, revised: false });
      await expect
        .poll(async () => (await mapi.listApiPages(apiId)).length, {
          timeout: 30_000,
          message: "API documentation page removed",
        })
        .toBe(0);
    });
  },
  { withPage: true, revised: false },
);
