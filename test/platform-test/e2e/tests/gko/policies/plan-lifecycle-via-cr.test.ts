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
 * Plan lifecycle via CR — batch 7.
 *
 * Xray tests:
 *   GKO-1459: Publishing and removing plans via the API CR syncs to APIM.
 *             Keyless + API Key plans are published; removing one leaves
 *             the other intact in APIM.
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

const TWO_PLANS = "crds/api-v4-definitions/v4-proxy-api-two-plans.yaml";
const ONE_PLAN = "crds/api-v4-definitions/v4-proxy-api-one-plan-remaining.yaml";
const API_NAME = "e2e-v4-two-plans";

test.describe("Plan lifecycle — via CR", () => {
  test.afterEach(async () => {
    await kubectlSafe.del(fixture(ONE_PLAN)).catch(() => {});
    await kubectlSafe.del(fixture(TWO_PLANS)).catch(() => {});
  });

  test(`Publish then remove plan via CR ${XRAY.PLANS.PLAN_LIFECYCLE_VIA_CR} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    await test.step("Deploy API with Keyless + API Key plans", async () => {
      await kubectl.apply(fixture(TWO_PLANS));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const apiId = (
      await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME)
    ).id;

    await test.step("Both plans exist and are PUBLISHED in APIM", async () => {
      await expect
        .poll(
          async () => {
            const plans = await mapi.listApiPlans(apiId);
            return plans
              .map((p) => ({ name: p.name, status: p.status }))
              .sort((a, b) => a.name.localeCompare(b.name));
          },
          { timeout: 30_000, intervals: [1_000] },
        )
        .toEqual([
          { name: "API Key plan", status: "PUBLISHED" },
          { name: "Keyless plan", status: "PUBLISHED" },
        ]);
    });

    await test.step("Remove API Key plan from CR", async () => {
      await kubectl.apply(fixture(ONE_PLAN));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("Only Keyless plan remains in APIM", async () => {
      await expect
        .poll(
          async () => {
            const plans = await mapi.listApiPlans(apiId);
            return plans.map((p) => p.name).sort();
          },
          { timeout: 30_000, intervals: [1_000] },
        )
        .toEqual(["Keyless plan"]);
    });

    await kubectl.del(fixture(ONE_PLAN));
  });
});
