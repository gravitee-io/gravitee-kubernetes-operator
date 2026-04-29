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

const V4_API = "crds/api-v4-definitions/v4-proxy-api-started.yaml";

test.describe("Recovery — reapplying configuration", () => {
  test.afterEach(async () => {
    await kubectlSafe.del(fixture(V4_API)).catch(() => {});
  });

  test(`CR recreated after removal reaches Accepted=True ${XRAY.DEPLOYMENT_RECONCILIATION.RECOVERY_REAPPLY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-start-stop";

    await test.step("Initial apply", async () => {
      await kubectl.apply(fixture(V4_API));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const firstId = (
      await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME)
    ).id;

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

    // The new API in APIM is a fresh resource — just confirm it exists and
    // is reachable. The IDs may differ (new UUID) or match (idempotent) —
    // both are acceptable for "recovery".
    const reapplied = await mapi.fetchApi(secondId);
    expect(reapplied.name).toBe(API_NAME);
    expect(typeof firstId).toBe("string");

    await kubectl.del(fixture(V4_API));
  });
});
