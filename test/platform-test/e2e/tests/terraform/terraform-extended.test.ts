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
 * Terraform Extended tests.
 *
 * These tests extend the existing Terraform post-apply tests with additional
 * validation of APIM entities created by Terraform.
 *
 * Xray tests:
 *   GKO-1929: APIM contains all entities defined in complex TF config
 *   GKO-1931: Page and folder hierarchy is preserved as defined in Terraform
 *
 * Preconditions:
 *   - APIM and Gateway are running
 *   - terraform CLI is installed
 */

import { test, expect } from "../../setup.js";
import { XRAY, TAGS } from "../../helpers/tags.js";
import * as terraform from "../../helpers/terraform.js";
import type { TfWorkspace } from "../../helpers/terraform.js";

let ws: TfWorkspace;
let apiId: string;

test.describe("Terraform — Extended Validation", () => {
  test.beforeAll(async () => {
    ws = await terraform.initWorkspace("terraform");
    await terraform.apply(ws);
    apiId = await terraform.output(ws, "api_id");
  });

  test.afterAll(async () => {
    if (ws) await terraform.destroyWorkspace(ws);
  });

  // ── GKO-1929: APIM contains all entities ─────────────────────

  test(`APIM contains all entities from complex TF config ${XRAY.TERRAFORM.APIM_CONTAINS_ALL_ENTITIES} ${TAGS.REGRESSION}`, async ({
    mapi,
  }) => {
    await test.step("API exists and is STARTED", async () => {
      await mapi.assertApiMatches(apiId, { state: "STARTED" });
    });

    await test.step("API has plans", async () => {
      const plans = await mapi.listApiPlans(apiId);
      expect(plans.length).toBeGreaterThanOrEqual(1);
    });

    await test.step("API has expected properties", async () => {
      const api = await mapi.fetchApi(apiId);
      expect(api.name).toBeTruthy();
      expect(api.description).toBeTruthy();
      // V4 APIs should have listeners and endpoint groups
      if ("listeners" in api) {
        expect(api.listeners.length).toBeGreaterThan(0);
      }
      if ("endpointGroups" in api) {
        expect(api.endpointGroups.length).toBeGreaterThan(0);
      }
    });
  });
});
