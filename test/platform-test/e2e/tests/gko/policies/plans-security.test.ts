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
 * Plans — Security Types tests.
 *
 * Xray tests:
 *   GKO-238: General conditions in V4 plan
 *
 * GKO-162/163 (OAuth2 / JWT plan security types) moved to the shared
 * cross-provisioner journey tests/scenarios/secure-api-with-plan — both plan
 * security types are now proven against GKO and Terraform. The admission test
 * below stays GKO-only.
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

test.describe("Plans — Security Types", () => {
  // GKO-162 (OAuth2 plan) and GKO-163 (JWT plan) are now covered by the
  // cross-provisioner journey tests/scenarios/secure-api-with-plan.

  // ── GKO-238: General conditions in V4 plan ──────────────────
  // Plan references a non-existing page as generalConditions.
  // The admission webhook rejects the CRD at apply time.

  test(`General conditions referencing non-existing page fails reconciliation ${XRAY.PLANS.GENERAL_CONDITIONS} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const fixturePath = fixture("plans/general-conditions/crd.yaml");

    await test.step("Apply is rejected by admission webhook", async () => {
      const stderr = await kubectl.applyExpectFailure(fixturePath);
      expect(stderr.toLowerCase()).toContain("general conditions");
    });
  });
});
