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
 * Subscription immutability — batch 7.
 *
 * Xray tests:
 *   GKO-1460: Changing the plan of an existing Subscription CR must be
 *             rejected / not take effect. The APIM-side subscription must
 *             continue to point at the original plan.
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

const V4_API = "crds/api-v4-definitions/v4-proxy-api-jwt-oauth2-plans.yaml";
const APP = "crds/applications/application-simple.yaml";
const SUB_A = "crds/subscriptions/subscription-plan-a.yaml";
const SUB_B = "crds/subscriptions/subscription-plan-b.yaml";

const API_NAME = "e2e-v4-sub-immutability";
const APP_NAME = "e2e-app-simple";
const SUB_NAME = "e2e-sub-plan-change";

test.describe("Subscription — immutability", () => {
  test.afterEach(async () => {
    await kubectlSafe.del(fixture(SUB_B)).catch(() => {});
    await kubectlSafe.del(fixture(SUB_A)).catch(() => {});
    await kubectlSafe.del(fixture(APP)).catch(() => {});
    await kubectlSafe.del(fixture(V4_API)).catch(() => {});
  });

  test(`Changing subscription plan is blocked ${XRAY.SUBSCRIPTIONS.SUBSCRIPTION_IMMUTABILITY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    await test.step("Deploy API, application, and subscription to JWT plan", async () => {
      await kubectl.apply(fixture(V4_API));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
      await kubectl.apply(fixture(APP));
      await kubectl.waitForCondition("application", APP_NAME, "Accepted");
      await kubectl.apply(fixture(SUB_A));
      await kubectl.waitForCondition("subscription", SUB_NAME, "Accepted");
    });

    const apiId = (
      await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME)
    ).id;
    const subId = (
      await kubectl.getStatus<{ id: string }>("subscription", SUB_NAME)
    ).id;

    const plans = await mapi.listApiPlans(apiId);
    const jwtPlanId = plans.find((p) => p.name === "JWT plan")?.id;
    expect(jwtPlanId, "expected JWT plan to exist").toBeTruthy();

    await test.step("APIM subscription points at JWT plan", async () => {
      const sub = await mapi.fetchSubscription(apiId, subId);
      const rawPlan = (sub as unknown as { plan?: unknown }).plan;
      const planId =
        typeof rawPlan === "string"
          ? rawPlan
          : (rawPlan as { id?: string } | undefined)?.id;
      expect(planId).toBe(jwtPlanId);
    });

    await test.step("Attempt to change the subscription plan is blocked", async () => {
      // The plan change can be rejected at admission OR accepted by
      // the CR but ignored by the operator. Either way, APIM must keep
      // the original plan — we tolerate both paths and only assert the
      // APIM-side invariant.
      try {
        await kubectl.apply(fixture(SUB_B));
      } catch {
        // admission-level rejection is fine
      }

      // Give the operator time to attempt (and reject) the change.
      await new Promise((r) => setTimeout(r, 5_000));

      const sub = await mapi.fetchSubscription(apiId, subId);
      const rawPlan = (sub as unknown as { plan?: unknown }).plan;
      const planId =
        typeof rawPlan === "string"
          ? rawPlan
          : (rawPlan as { id?: string } | undefined)?.id;
      expect(planId, "APIM subscription plan must remain JWT").toBe(jwtPlanId);
    });

    // Cleanup in reverse dependency order: sub → app → api.
    await kubectl.del(fixture(SUB_A));
    await kubectl.del(fixture(APP));
    await kubectl.del(fixture(V4_API));
  });
});
