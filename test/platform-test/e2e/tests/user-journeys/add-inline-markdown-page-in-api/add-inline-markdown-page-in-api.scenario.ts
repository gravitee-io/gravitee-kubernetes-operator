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
 * Journey: add an inline markdown page to an API.
 *
 * As an API producer, I ship API documentation alongside the API definition. An
 * inline markdown page is an attribute of apim_apiv4 (spec.pages / pages[]) — there
 * is no standalone Page Terraform resource — yet the journey is fully expressible
 * through both provisioners: a page set through either driver lands in APIM and is
 * removed when stripped. Inline page fetchers (pages[].source) are also expressible
 * on both drivers and are a feasible follow-up; only V2 documentation has no
 * Terraform path (no apim_apiv2) and stays GKO-only.
 *
 * This is the same inline-attribute pattern as label-an-api, except the payload is
 * a nested object list, so the with/without fixtures and the assertion compare a
 * page object ({ name, type, content, published }) rather than a plain string.
 *
 * Fixtures are co-located in this folder.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import { forEachProvisioner } from "../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../helpers/provisioner-env.js";
import type { Page } from "../../../../src/types/apim.js";

const here = path.dirname(fileURLToPath(import.meta.url));

// The inline page both fixtures attach. Content is compared trimmed because the
// two drivers' block-scalar / heredoc syntaxes differ on the trailing newline.
const PAGE = {
  name: "Getting started",
  type: "MARKDOWN",
  content: "# Getting started\n\nCall `GET /` to reach the upstream echo endpoint.",
  published: true,
};

/** Project a fetched page to the fields under test (content normalised). */
function project(p: Page) {
  return {
    name: p.name,
    type: p.type,
    content: (p.content ?? "").trim(),
    published: p.published ?? false,
  };
}

/** The single knob: whether the API carries the inline page. */
interface PageParams {
  withPage: boolean;
}

forEachProvisioner<PageParams>(
  {
    title: "Add an inline markdown page to a V4 API",
    provisioners: {
      gko: gkoScenario<PageParams>({
        manifests: [],
        roles: { api: "documented-api" },
        dynamicRoles: ["api"],
        applyParams: async (k, params) => {
          await k.apply(path.join(here, params.withPage ? "gko/api-with-page.yaml" : "gko/api-without-page.yaml"));
        },
      }),
      terraform: tfScenario<PageParams>({
        fixture: path.join(here, "terraform"),
        toVars: (params) => ({ with_page: params.withPage }),
      }),
    },
    xray: {
      // The GKO arm reuses the end-to-end documentation-reconcile test (page lands
      // on create, removed on strip); this journey is its sole coverage.
      gko: XRAY.PAGES.V4_DOC_RECONCILED,
      terraform: XRAY.TERRAFORM.API_INLINE_PAGES_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 60_000 },
  },
  async ({ provisioned, mapi }) => {
    const apiId = await provisioned.apiId();

    await test.step("Page attached through the provisioner lands in APIM", async () => {
      await expect
        .poll(async () => (await mapi.listApiPages(apiId)).map(project), {
          timeout: 30_000,
          message: "API documentation page reaches APIM",
        })
        .toEqual([{ ...PAGE, content: PAGE.content.trim() }]);
    });

    await test.step("Stripping the page removes it in APIM", async () => {
      await provisioned.update({ withPage: false });
      await expect
        .poll(async () => (await mapi.listApiPages(apiId)).length, {
          timeout: 30_000,
          message: "API documentation page removed",
        })
        .toBe(0);
    });
  },
  { withPage: true },
);
