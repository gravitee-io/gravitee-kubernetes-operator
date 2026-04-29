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
 * Recovery after CR removal.
 *
 * Xray tests:
 *   GKO-1808: Removing the CR and reapplying it recreates the API in APIM
 *             with the same configuration, and status reaches Accepted=True.
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";
import type { ApiV4 } from "../../../../src/types/apim.js";

const V4_API = "crds/api-v4-definitions/v4-proxy-api-started.yaml";

test.describe("Recovery — reapplying configuration", () => {
  test.afterEach(async () => {
    await kubectlSafe.del(fixture(V4_API)).catch(() => {});
  });

  test(`CR recreated after removal reaches Accepted=True ${XRAY.DEPLOYMENT_RECONCILIATION.RECOVERY_REAPPLY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    // The full recovery cycle (apply → wait Accepted → delete → wait
    // for APIM cleanup → reapply → wait Accepted) is server-side bound
    // and routinely sits at ~25-30s before any verification runs. Same
    // pattern as operator-restart.test.ts.
    test.slow();

    const API_NAME = "e2e-v4-start-stop";

    await test.step("Initial apply", async () => {
      await kubectl.apply(fixture(V4_API));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const firstId = (
      await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME)
    ).id;
    expect(typeof firstId).toBe("string");

    await test.step("Delete CR and wait for APIM cleanup", async () => {
      await kubectl.del(fixture(V4_API));
      await kubectl.waitForDeletion("apiv4definition", API_NAME);
    });

    await test.step("Reapply same CR", async () => {
      await kubectl.apply(fixture(V4_API));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const secondId = (
      await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME)
    ).id;

    // The recovery promise is "same configuration", not just "exists" —
    // a regression that recreated the API with missing plans, wrong
    // state, or a broken endpoint config would silently pass a
    // name-only check. Assert state, version, listener path, and plan
    // presence to match the fixture.
    await test.step("Reapplied API matches the fixture configuration", async () => {
      await mapi.waitForApiMatches(
        secondId,
        {
          name: API_NAME,
          state: "STARTED",
          apiVersion: "1.0.0",
          definitionVersion: "V4",
          // APIM normalises the listener path with a trailing slash.
          listeners: [{ paths: [{ path: "/e2e-v4-start-stop/" }] }],
        },
        { description: "reapplied API has the original spec" },
      );
    });

    await test.step("Reapplied API has the keyless plan and http-proxy endpoint", async () => {
      const [plans, api] = await Promise.all([
        mapi.listApiPlans(secondId),
        mapi.fetchApi(secondId) as Promise<ApiV4>,
      ]);

      const keyless = plans.find((p) => p.name === "Free plan");
      expect(keyless, "expected Free plan to exist after reapply").toBeTruthy();
      expect(keyless?.security?.type).toBe("KEY_LESS");

      expect(api.endpointGroups?.length, "at least one endpoint group must exist").toBeGreaterThan(
        0,
      );
      expect(api.endpointGroups?.[0]?.endpoints?.[0]?.type).toBe("http-proxy");
    });

    await kubectl.del(fixture(V4_API));
  });
});
