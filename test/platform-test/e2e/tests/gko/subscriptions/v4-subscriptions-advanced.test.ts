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
 * V4 Plans & Subscriptions — Advanced scenarios (batch 4).
 *
 * Xray tests:
 *   GKO-795: Subscription read-only enforcement — operator must reject mutations
 *            of immutable fields (plan, api, app) on an accepted subscription.
 *   GKO-822: Removing a subscribed V4 JWT plan from an API must be rejected by
 *            the admission webhook.
 *   GKO-826: Same as GKO-822, but for OAuth2 plans.
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { readFile } from "node:fs/promises";
import YAML from "yaml";
import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

const APP_SIMPLE = "crds/applications/application-simple.yaml";
const JWT_API = "crds/api-v4-definitions/v4-api-jwt-plan.yaml";
const JWT_API_PLAN_REMOVED = "crds/api-v4-definitions/v4-api-jwt-plan-removed.yaml";
const OAUTH2_API = "crds/api-v4-definitions/v4-api-oauth2-plan.yaml";
const OAUTH2_API_PLAN_REMOVED = "crds/api-v4-definitions/v4-api-oauth2-plan-removed.yaml";
const SUB_JWT = "crds/subscriptions/subscription-jwt-v4.yaml";
const SUB_OAUTH2 = "crds/subscriptions/subscription-oauth2-v4.yaml";

test.describe("V4 Plans & Subscriptions — Advanced", () => {
  test.afterEach(async () => {
    // Safety-net cleanup: subscriptions → applications → APIs
    await kubectlSafe.del(fixture(SUB_JWT)).catch(() => {});
    await kubectlSafe.del(fixture(SUB_OAUTH2)).catch(() => {});
    await kubectlSafe.del(fixture(APP_SIMPLE)).catch(() => {});
    await kubectlSafe.del(fixture(JWT_API)).catch(() => {});
    await kubectlSafe.del(fixture(OAUTH2_API)).catch(() => {});
  });

  // ── GKO-795: Subscription read-only enforcement ─────────────

  test(`Subscription CRD rejects mutation of immutable fields once Accepted ${XRAY.SUBSCRIPTIONS.V4_SUBSCRIPTION_READ_ONLY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const apiFixture = fixture(JWT_API);
    const appFixture = fixture(APP_SIMPLE);
    const subFixture = fixture(SUB_JWT);

    await kubectl.apply(apiFixture);
    await kubectl.apply(appFixture);
    await kubectl.waitForCondition("apiv4definition", "e2e-v4-jwt-api", "Accepted");
    await kubectl.waitForCondition("application", "e2e-app-simple", "Accepted");
    await kubectl.apply(subFixture);
    await kubectl.waitForCondition("subscription", "e2e-sub-jwt-v4", "Accepted");

    const subStatus = await kubectl.getStatus<{ id: string }>("subscription", "e2e-sub-jwt-v4");
    const apiId = (await kubectl.getStatus<{ id: string }>("apiv4definition", "e2e-v4-jwt-api")).id;
    expect(subStatus.id).toBeTruthy();

    await test.step("Mutating plan ref on accepted subscription is rejected", async () => {
      const raw = await readFile(subFixture, "utf8");
      const doc = YAML.parse(raw) as { spec: { plan: string } };
      doc.spec.plan = "illegal-plan";
      const stderr = await kubectl.applyStringExpectFailure(YAML.stringify(doc));
      expect(stderr.toLowerCase()).toMatch(/immutable|illegal/);
    });

    await test.step("Subscription remains Accepted after rejected mutation", async () => {
      await kubectl.waitForCondition("subscription", "e2e-sub-jwt-v4", "Accepted");
      await mapi.assertSubscriptionAccepted(apiId, subStatus.id);
    });

    // Cleanup in reverse dependency order (subscriptions → apps → APIs)
    await kubectl.del(subFixture);
    await kubectl.del(appFixture);
    await kubectl.del(apiFixture);
  });

  // ── GKO-822: JWT plan deletion with active subscription ────

  test(`Removing a subscribed V4 JWT plan is rejected by admission ${XRAY.SUBSCRIPTIONS.V4_JWT_PLAN_DELETION_WITH_SUB} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const apiFixture = fixture(JWT_API);
    const apiNoPlanFixture = fixture(JWT_API_PLAN_REMOVED);
    const appFixture = fixture(APP_SIMPLE);
    const subFixture = fixture(SUB_JWT);

    await kubectl.apply(apiFixture);
    await kubectl.apply(appFixture);
    await kubectl.waitForCondition("apiv4definition", "e2e-v4-jwt-api", "Accepted");
    await kubectl.waitForCondition("application", "e2e-app-simple", "Accepted");
    await kubectl.apply(subFixture);
    await kubectl.waitForCondition("subscription", "e2e-sub-jwt-v4", "Accepted");

    await test.step("Applying API without the subscribed JWT plan is rejected", async () => {
      const stderr = await kubectl.applyExpectFailure(apiNoPlanFixture);
      expect(stderr).toMatch(/Plan \[.*\] could not be found/i);
    });

    // Cleanup
    await kubectl.del(subFixture);
    await kubectl.del(appFixture);
    await kubectl.del(apiFixture);
  });

  // ── GKO-826: OAuth2 plan deletion with active subscription ─

  test(`Removing a subscribed V4 OAuth2 plan is rejected by admission ${XRAY.SUBSCRIPTIONS.V4_OAUTH2_PLAN_DELETION_WITH_SUB} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const apiFixture = fixture(OAUTH2_API);
    const apiNoPlanFixture = fixture(OAUTH2_API_PLAN_REMOVED);
    const appFixture = fixture(APP_SIMPLE);
    const subFixture = fixture(SUB_OAUTH2);

    await kubectl.apply(apiFixture);
    await kubectl.apply(appFixture);
    await kubectl.waitForCondition("apiv4definition", "e2e-v4-oauth2-api", "Accepted");
    await kubectl.waitForCondition("application", "e2e-app-simple", "Accepted");
    await kubectl.apply(subFixture);
    await kubectl.waitForCondition("subscription", "e2e-sub-oauth2-v4", "Accepted");

    await test.step("Applying API without the subscribed OAuth2 plan is rejected", async () => {
      const stderr = await kubectl.applyExpectFailure(apiNoPlanFixture);
      expect(stderr).toMatch(/Plan \[.*\] could not be found/i);
    });

    await kubectl.del(subFixture);
    await kubectl.del(appFixture);
    await kubectl.del(apiFixture);
  });
});
