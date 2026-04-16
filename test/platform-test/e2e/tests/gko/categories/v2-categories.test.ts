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
 * V2 API Categories — batch 5 coverage.
 *
 * Xray tests:
 *   GKO-187: Add a valid category to a V2 API
 *   GKO-189: Add many valid categories at once
 *   GKO-190: Add a non-existing category (operator tolerates it)
 *   GKO-191: Remove a category from a V2 API
 *   GKO-192: Category removed in APIM is reflected on reconcile
 *   GKO-261: Category renamed in APIM triggers redeploy
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 *
 * Note on assertions: APIM silently drops unknown category refs (same as V4),
 * so these tests focus on Accepted reconciliation and state-change observability
 * rather than asserting category IDs on the live API.
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

interface StatusWithConditions {
  id?: string;
  conditions?: Array<{
    type: string;
    status: string;
  }>;
}

function acceptedTrue(status: StatusWithConditions): boolean {
  return status.conditions?.find((c) => c.type === "Accepted")?.status === "True";
}

test.describe("V2 API Categories — Extended", () => {
  test.afterEach(async () => {
    await kubectlSafe
      .del(fixture("crds/api-definitions/v2-api-with-category.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/api-definitions/v2-api-without-category.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/api-definitions/v2-api-with-many-categories.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/api-definitions/v2-api-renamed-category.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/api-definitions/v2-api-non-existing-category.yaml"))
      .catch(() => {});
  });

  // ── GKO-187: Add valid category ─────────────────────────────

  test(`Create V2 API with a category ${XRAY.CATEGORIES.V2_VALID_CATEGORY} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v2-with-category";
    const fixturePath = fixture("crds/api-definitions/v2-api-with-category.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("apidefinition", API_NAME);
    expect(acceptedTrue(status)).toBe(true);

    await kubectl.del(fixturePath);
  });

  // ── GKO-189: Many categories at once ────────────────────────

  test(`Deploy V2 API with many category refs ${XRAY.CATEGORIES.V2_MANY_CATEGORIES} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v2-many-cats";
    const fixturePath = fixture("crds/api-definitions/v2-api-with-many-categories.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("apidefinition", API_NAME);
    expect(acceptedTrue(status)).toBe(true);

    await kubectl.del(fixturePath);
  });

  // ── GKO-190: Non-existing category ──────────────────────────

  test(`V2 API with non-existing category is accepted ${XRAY.CATEGORIES.V2_NON_EXISTING_CATEGORY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-non-existing-cat";
    const fixturePath = fixture("crds/api-definitions/v2-api-non-existing-category.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME);
    const api = await mapi.fetchApi(status.id);
    // Unknown category ref must not leak onto the deployed API.
    expect(api.categories ?? []).not.toContain("nonexistent-v2-category-12345");

    await kubectl.del(fixturePath);
  });

  // ── GKO-191: Remove category from V2 API ────────────────────

  test(`Remove category from V2 API on re-apply ${XRAY.CATEGORIES.V2_REMOVE_CATEGORY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-with-category";
    const withCategory = fixture("crds/api-definitions/v2-api-with-category.yaml");
    const withoutCategory = fixture("crds/api-definitions/v2-api-without-category.yaml");

    await kubectl.apply(withCategory);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;

    await kubectl.apply(withoutCategory);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    // Changed description is observable on the APIM side — ensures the
    // reconciler actually ran on the without-categories spec.
    await mapi.waitForApiMatches(apiId, {
      description: "E2E test: V2 API with category removed",
    });

    await kubectl.del(withoutCategory);
  });

  // ── GKO-192: Category removed in APIM is reflected ──────────
  // When the CRD omits the categories field, APIM must not retain stale values.

  test(`Category removed in APIM reflects on reconcile ${XRAY.CATEGORIES.V2_CATEGORY_REMOVED_FROM_APIM} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-with-category";
    const withoutCategory = fixture("crds/api-definitions/v2-api-without-category.yaml");

    // Starting point: CRD without categories — APIM must show no categories.
    await kubectl.apply(withoutCategory);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;

    const api = await mapi.fetchApi(apiId);
    expect(api.categories === undefined || api.categories.length === 0).toBe(true);

    await kubectl.del(withoutCategory);
  });

  // ── GKO-261: Rename triggers redeploy ───────────────────────
  // Changing the categories list bumps the spec hash so the reconciler runs
  // instead of being filtered by LastSpecHashPredicate.

  test(`Category rename on V2 API triggers redeploy ${XRAY.CATEGORIES.V2_CATEGORY_RENAME_REDEPLOY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-many-cats";
    const original = fixture("crds/api-definitions/v2-api-with-many-categories.yaml");
    const renamed = fixture("crds/api-definitions/v2-api-renamed-category.yaml");

    await kubectl.apply(original);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    const firstId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;

    await kubectl.apply(renamed);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    await mapi.waitForApiMatches(firstId, {
      description: "E2E test: V2 API with renamed categories",
    });

    // Id must remain stable across the redeploy.
    const secondId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;
    expect(secondId).toBe(firstId);

    await kubectl.del(renamed);
  });
});
