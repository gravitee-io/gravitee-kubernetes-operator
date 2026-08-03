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
 * Category reference tolerance on a V4 API.
 *
 * Xray tests:
 *   GKO-269: A reference to a category that does not exist is ignored
 *
 * Adding, rewriting and removing a policy on a flow (GKO-94/95/96) is the shared
 * journey tests/user-journeys/api-producer/apply-policies-to-a-flow.
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS, PROVISIONER } from "../../../helpers/tags.js";
import * as kubectl from "../../../helpers/kubectl.js";

test.describe(`Categories — reference tolerance ${PROVISIONER.GKO}`, () => {
  // Safety-net cleanup: runs even if a test times out before its inline
  // cleanup. Each del() ignores errors (the resource may already be gone).
  test.afterEach(async () => {
    await kubectl.del(fixture("categories/v4-with-labels/crd.yaml")).catch(() => {});
  });

  // ── GKO-269: Non-existing category ───────────────────────────
  // Assigning and removing a valid category (@GKO-267 / @GKO-270) is covered by
  // the cross-provisioner journey tests/user-journeys/api-producer/assign-categories-to-api/.
  // This case covers the GKO-only behaviour that an unknown category reference is
  // tolerated: the API still deploys.

  test(`Non-existing category is ignored ${XRAY.CATEGORIES.NON_EXISTING_CATEGORY_V4} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-labels-cats";
    const fixturePath = fixture("categories/v4-with-labels/crd.yaml");

    // The fixture has labels but no categories field — API should deploy fine
    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    await mapi.assertApiStarted(status.id);

    await kubectl.del(fixturePath);
  });
});
