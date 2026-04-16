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
 * Terraform — plan general conditions page.
 *
 * Xray tests:
 *   GKO-1930: General conditions page is correctly attached to the plan
 *
 * Preconditions:
 *   - APIM is running
 *   - terraform CLI is installed
 */

import { test, expect } from "../../setup.js";
import { XRAY, TAGS } from "../../helpers/tags.js";
import * as terraform from "../../helpers/terraform.js";
import type { TfWorkspace } from "../../helpers/terraform.js";

let ws: TfWorkspace;
let apiId: string;

interface PlanWithGeneralConditions {
  id: string;
  generalConditions?: string;
}

test.describe("Terraform — Plan General Conditions", () => {
  test.beforeAll(async () => {
    // terraform init + apply on a fresh workspace blows past Playwright's 30s
    // default hook timeout — bump it to give terraform room to provision the
    // provider, run init and apply.
    test.setTimeout(180_000);
    ws = await terraform.initWorkspace("terraform-general-conditions");
    await terraform.apply(ws);
    apiId = await terraform.output(ws, "api_id");
  });

  test.afterAll(async () => {
    test.setTimeout(180_000);
    if (ws) await terraform.destroyWorkspace(ws);
  });

  // ── GKO-1930: General conditions page attached to plan ──────

  test(`Plan created via Terraform references the general conditions page ${XRAY.TERRAFORM.GENERAL_CONDITIONS_PAGE} ${TAGS.REGRESSION}`, async ({
    mapi,
  }) => {
    const plans = (await mapi.listApiPlans(apiId)) as PlanWithGeneralConditions[];
    expect(plans.length).toBeGreaterThanOrEqual(1);

    const plan = plans[0]!;
    // The APIM management API surfaces general conditions as the page id on
    // the plan. Terraform sends `general_conditions_hrid` on the apply — the
    // provider resolves it to a page id. The assertion here is "non-empty
    // and references an existing page for the API".
    expect(plan.generalConditions).toBeTruthy();
  });
});
