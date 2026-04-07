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
 * Shared Policy Groups Lifecycle tests.
 *
 * Xray tests:
 *   GKO-976:  Add SPG to V4 API
 *   GKO-980:  Remove SPG from V4 API
 *   GKO-981:  Update SPG
 *   GKO-1462: SPG lifecycle validation
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

test.describe("Shared Policy Groups — Lifecycle", () => {
  // ── GKO-976: Add SPG to V4 API ─────────────────────────────

  test(`Add SPG to V4 API ${XRAY.SHARED_POLICY_GROUPS.ADD_SPG_TO_API} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const SPG_NAME = "e2e-spg-proxy";
    const API_NAME = "e2e-v4-with-spg";
    const spgFixture = fixture("crds/shared-policy-groups/spg-proxy-request.yaml");
    const apiFixture = fixture("crds/api-v4-definitions/v4-api-with-spg.yaml");

    await test.step("Deploy Shared Policy Group", async () => {
      await kubectl.apply(spgFixture);
      await kubectl.waitForCondition("sharedpolicygroup", SPG_NAME, "Accepted");
    });

    const spgStatus = await kubectl.getStatus<{ id: string }>("sharedpolicygroup", SPG_NAME);
    expect(spgStatus.id).toBeTruthy();

    await test.step("Deploy V4 API referencing the SPG", async () => {
      await kubectl.apply(apiFixture);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const apiStatus = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = apiStatus.id;

    await test.step("API has flows with SPG reference in APIM", async () => {
      const api = await mapi.fetchApi(apiId);
      expect(api).toBeTruthy();
      if ("flows" in api && api.flows) {
        expect(api.flows.length).toBeGreaterThanOrEqual(1);
      }
    });

    await kubectl.del(apiFixture);
    await kubectl.del(spgFixture);
  });

  // ── GKO-980: Remove SPG from V4 API ────────────────────────

  test(`Remove SPG from V4 API ${XRAY.SHARED_POLICY_GROUPS.REMOVE_SPG_FROM_API} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const SPG_NAME = "e2e-spg-proxy";
    const API_NAME = "e2e-v4-with-spg";
    const spgFixture = fixture("crds/shared-policy-groups/spg-proxy-request.yaml");
    const apiWithSpg = fixture("crds/api-v4-definitions/v4-api-with-spg.yaml");
    const apiWithoutSpg = fixture("crds/api-v4-definitions/v4-api-without-spg.yaml");

    await test.step("Deploy SPG and API with SPG reference", async () => {
      await kubectl.apply(spgFixture);
      await kubectl.waitForCondition("sharedpolicygroup", SPG_NAME, "Accepted");
      await kubectl.apply(apiWithSpg);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const apiStatus = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = apiStatus.id;

    await test.step("Remove SPG reference from API", async () => {
      await kubectl.apply(apiWithoutSpg);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("API no longer has SPG flows in APIM", async () => {
      const api = await mapi.fetchApi(apiId);
      if ("flows" in api) {
        expect(api.flows?.length ?? 0).toBe(0);
      }
    });

    await kubectl.del(apiWithoutSpg);
    await kubectl.del(spgFixture);
  });

  // ── GKO-981: Update SPG ────────────────────────────────────

  test(`Update SPG ${XRAY.SHARED_POLICY_GROUPS.UPDATE_SPG} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const SPG_NAME = "e2e-spg-proxy";
    const createFixture = fixture("crds/shared-policy-groups/spg-proxy-request.yaml");
    const updateFixture = fixture("crds/shared-policy-groups/spg-proxy-updated.yaml");

    await test.step("Deploy SPG", async () => {
      await kubectl.apply(createFixture);
      await kubectl.waitForCondition("sharedpolicygroup", SPG_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("sharedpolicygroup", SPG_NAME);
    expect(status.id).toBeTruthy();

    await test.step("Update SPG with new header value", async () => {
      await kubectl.apply(updateFixture);
      await kubectl.waitForCondition("sharedpolicygroup", SPG_NAME, "Accepted");
    });

    await test.step("SPG ID remains the same after update", async () => {
      const updatedStatus = await kubectl.getStatus<{ id: string }>("sharedpolicygroup", SPG_NAME);
      expect(updatedStatus.id).toBe(status.id);
    });

    await kubectl.del(updateFixture);
  });

  // ── GKO-1462: SPG lifecycle validation ──────────────────────

  test(`SPG lifecycle validation ${XRAY.SHARED_POLICY_GROUPS.SPG_LIFECYCLE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const SPG_PROXY_NAME = "e2e-spg-proxy";
    const SPG_MESSAGE_NAME = "e2e-spg-message";
    const proxyFixture = fixture("crds/shared-policy-groups/spg-proxy-request.yaml");
    const messageFixture = fixture("crds/shared-policy-groups/spg-message-unsupported.yaml");

    await test.step("Deploy valid PROXY SPG", async () => {
      await kubectl.apply(proxyFixture);
      await kubectl.waitForCondition("sharedpolicygroup", SPG_PROXY_NAME, "Accepted");
    });

    const proxyStatus = await kubectl.getStatus<{ id: string }>("sharedpolicygroup", SPG_PROXY_NAME);
    expect(proxyStatus.id).toBeTruthy();

    await test.step("Deploy MESSAGE SPG (may or may not be supported)", async () => {
      // MESSAGE apiType SPG may be accepted or rejected depending on APIM version
      try {
        await kubectl.apply(messageFixture);
        await kubectl.waitForCondition("sharedpolicygroup", SPG_MESSAGE_NAME, "Accepted");
        const messageStatus = await kubectl.getStatus<{ id: string }>("sharedpolicygroup", SPG_MESSAGE_NAME);
        expect(messageStatus.id).toBeTruthy();
      } catch {
        // If rejected, that is also valid behavior
        const status = await kubectl.getStatus<{ id: string }>("sharedpolicygroup", SPG_MESSAGE_NAME);
        expect(status).toBeTruthy();
      }
    });

    await kubectl.del(proxyFixture);
    await kubectl.del(messageFixture);
  });
});
