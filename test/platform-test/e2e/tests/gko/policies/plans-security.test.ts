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
 *   GKO-162: Add OAuth2 plan to V4 API
 *   GKO-163: Add JWT plan to V4 API
 *   GKO-238: General conditions in V4 plan
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

test.describe("Plans — Security Types", () => {
  // ── GKO-162: Add OAuth2 plan to V4 API ──────────────────────

  test(`Add OAuth2 plan to V4 API ${XRAY.PLANS.OAUTH2_PLAN_V4} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-oauth2-plan";
    const fixturePath = fixture("crds/api-v4-definitions/v4-proxy-api-oauth2-plan.yaml");

    await test.step("Deploy API with OAuth2 plan", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("Plan exists in APIM with security type OAUTH2", async () => {
      const plans = await mapi.listApiPlans(apiId);
      expect(plans.length).toBeGreaterThanOrEqual(1);
      const oauth2Plan = plans.find((p: { security: { type: string } }) => p.security?.type === "OAUTH2");
      expect(oauth2Plan).toBeTruthy();
      expect(oauth2Plan!.name).toBe("OAuth2 plan");
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-163: Add JWT plan to V4 API ─────────────────────────

  test(`Add JWT plan to V4 API ${XRAY.PLANS.JWT_PLAN_V4} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-jwt-plan";
    const fixturePath = fixture("crds/api-v4-definitions/v4-proxy-api-jwt-plan.yaml");

    await test.step("Deploy API with JWT plan", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("Plan exists in APIM with security type JWT", async () => {
      const plans = await mapi.listApiPlans(apiId);
      expect(plans.length).toBeGreaterThanOrEqual(1);
      const jwtPlan = plans.find((p: { security: { type: string } }) => p.security?.type === "JWT");
      expect(jwtPlan).toBeTruthy();
      expect(jwtPlan!.name).toBe("JWT plan");
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-238: General conditions in V4 plan ──────────────────
  // Plan references a non-existing page as generalConditions.
  // The admission webhook rejects the CRD at apply time.

  test(`General conditions referencing non-existing page fails reconciliation ${XRAY.PLANS.GENERAL_CONDITIONS} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const fixturePath = fixture("crds/api-v4-definitions/v4-proxy-api-general-conditions.yaml");

    await test.step("Apply is rejected by admission webhook", async () => {
      const stderr = await kubectl.applyExpectFailure(fixturePath);
      expect(stderr.toLowerCase()).toContain("general conditions");
    });
  });
});
