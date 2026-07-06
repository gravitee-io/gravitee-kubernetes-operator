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
 * Journey: assign a category to an API.
 *
 * As an API producer, I organise my API in the portal by assigning it to a
 * category. Categories are an inline attribute of apim_apiv4 (no standalone
 * Terraform resource and no GKO Category CRD): an API can only *reference* a
 * category that already exists in the environment, and APIM silently drops
 * references to unknown categories. So the category is created ONCE as a
 * provisioner-agnostic precondition via mAPI, and both drivers then assign it by
 * key (spec.categories / apim_apiv4.categories) and remove it when stripped.
 *
 * This is the same inline-attribute pattern as label-an-api (the reference
 * journey), with the added category precondition.
 *
 * Fixtures are co-located in this folder.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import { forEachProvisioner } from "../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../helpers/provisioner-env.js";
import type { Category } from "../../../../src/types/apim.js";

const here = path.dirname(fileURLToPath(import.meta.url));

// Category slug. e2e-prefixed for isolation in the shared APIM env, matching the
// e2e-sa-* / e2e-v4-* convention for test-created resources. Passed as both name
// and key; the beforeEach guard asserts APIM kept this exact key, and both static
// fixtures reference it verbatim.
const CATEGORY_KEY = "e2e-portal-category";

/** The single knob: whether the API carries the category. */
interface CategoryParams {
  withCategories: boolean;
}

// The precondition category, created once per test in beforeEach (mapi is a
// test-scoped fixture, so it is unreachable from beforeAll). Module scope is safe
// because the suite runs serially (workers: 1).
let category: Category | undefined;

// Create the referenced category BEFORE the test body provisions the API: for
// both drivers the API (and its category-ref resolution) is created inside the
// body's provision() call, and APIM drops references to categories that do not
// yet exist. A top-level beforeEach always runs before the test body.
test.beforeEach(async ({ mapi }) => {
  category = await mapi.createCategory({ key: CATEGORY_KEY, name: CATEGORY_KEY });
  // Fail fast (instead of with a confusing "category never lands" timeout) if
  // APIM's slugification produced a key other than the one the fixtures use.
  expect(category.key, "APIM category key must match the fixtures").toBe(CATEGORY_KEY);
});

// Registered BEFORE forEachProvisioner so it runs AFTER the generator's own
// teardown (Playwright runs afterEach in reverse registration order): the API is
// destroyed first, then its category precondition is removed.
test.afterEach(async ({ mapi }) => {
  if (category) await mapi.deleteCategory(category.id).catch(() => {});
  category = undefined;
});

forEachProvisioner<CategoryParams>(
  {
    title: "Assign a category to a V4 API",
    provisioners: {
      gko: gkoScenario<CategoryParams>({
        manifests: [],
        roles: { api: "categorized-api" },
        dynamicRoles: ["api"],
        applyParams: async (k, params) => {
          await k.apply(
            path.join(here, params.withCategories ? "gko/api-with-categories.yaml" : "gko/api-without-categories.yaml"),
          );
        },
      }),
      terraform: tfScenario<CategoryParams>({
        fixture: path.join(here, "terraform"),
        toVars: (params) => ({ with_categories: params.withCategories }),
      }),
    },
    xray: {
      // The GKO arm covers both the assign (@GKO-267) and remove (@GKO-270)
      // category cases; this journey is their sole coverage.
      gko: [XRAY.CATEGORIES.VALID_CATEGORY_V4, XRAY.CATEGORIES.REMOVE_CATEGORY_V4],
      terraform: XRAY.TERRAFORM.API_CATEGORIES_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 60_000 },
  },
  async ({ provisioned, mapi }) => {
    const apiId = await provisioned.apiId();
    const cat = category!;

    await test.step("Category assigned through the provisioner lands in APIM", async () => {
      await expect
        .poll(
          async () => {
            const assigned = (await mapi.fetchApi(apiId)).categories ?? [];
            // APIM may echo the reference by key or by UUID depending on version;
            // normalise a UUID back to the key so the expectation is
            // representation-independent.
            return assigned.map((c) => (c === cat.id ? cat.key : c)).slice().sort();
          },
          { timeout: 30_000, message: "API category reaches APIM" },
        )
        .toEqual([cat.key].sort());
    });

    await test.step("Stripping the category removes it in APIM", async () => {
      await provisioned.update({ withCategories: false });
      await expect
        .poll(async () => (await mapi.fetchApi(apiId)).categories ?? [], {
          timeout: 30_000,
          message: "API category removed",
        })
        .toEqual([]);
    });
  },
  { withCategories: true },
);
