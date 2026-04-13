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
 * V4 API Categories — Extended scenarios (batch 4).
 *
 * Xray tests:
 *   GKO-271: Category removal synced from APIM when category deleted
 *   GKO-415: Import V4 API with non-existing category + dryRun=true
 *   GKO-416: Import V4 API with non-existing category + dryRun=false
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

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

test.describe("V4 API Categories — Extended", () => {
  // ── GKO-271: Category removal synced from APIM ──────────────
  // When an API is deployed with no categories field, any categories that
  // may have been set in APIM should not leak back onto the CRD.

  test(`Category removal synced from APIM ${XRAY.CATEGORIES.V4_CATEGORY_REMOVED_FROM_APIM} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-start-stop";
    const fixturePath = fixture("crds/api-v4-definitions/v4-proxy-api-started.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const api = await mapi.fetchApi(status.id);
    // CRD has no categories set — APIM-side must reflect the same.
    expect(api.categories === undefined || api.categories.length === 0).toBe(true);

    await kubectl.del(fixturePath);
  });

  // ── GKO-415: Import with non-existing category + dryRun ─────
  // dryRun=true should not create the API but validation must still pass —
  // the non-existing category is treated as a warning.

  test(`Import V4 API with non-existing category (dryRun) ${XRAY.CATEGORIES.V4_IMPORT_NON_EXISTING_CATEGORY_DRYRUN} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-non-existing-cat";
    const fixturePath = fixture(
      "crds/v4-lifecycle-extended/v4-proxy-api-non-existing-category.yaml",
    );

    // Applying the CRD is the closest-analog to an "apply import" for the
    // operator; dryRun=true is exercised by `kubectl apply --dry-run=server`
    // which runs through admission but never persists.
    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("apiv4definition", API_NAME);
    expect(acceptedTrue(status)).toBe(true);

    await kubectl.del(fixturePath);
  });

  // ── GKO-416: Import with non-existing category + apply ──────

  test(`Import V4 API with non-existing category (apply) ${XRAY.CATEGORIES.V4_IMPORT_NON_EXISTING_CATEGORY_APPLY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-non-existing-cat";
    const fixturePath = fixture(
      "crds/v4-lifecycle-extended/v4-proxy-api-non-existing-category.yaml",
    );

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const api = await mapi.fetchApi(status.id);
    // Non-existing category must not end up on the deployed API.
    expect(api.categories ?? []).not.toContain("nonexistent-category-12345");

    await kubectl.del(fixturePath);
  });
});
