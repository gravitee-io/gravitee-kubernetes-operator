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
 * V4 Proxy API Lifecycle tests.
 *
 * Xray tests:
 *   GKO-71:  Deploy V4 Proxy API with syncFrom Management
 *   GKO-140: Delete a V4 API
 *   GKO-165: Create V4 proxy API with missing required fields
 *   GKO-176: Should not deploy when no changes are made to V4 CRD
 *   GKO-212: API is re-deployed when applying the same CRD after a delete
 *   GKO-469: Context path already exists
 *   GKO-502: No plans + STARTED → error
 *   GKO-503: No plans + STOPPED → OK
 *   GKO-859: Add Failover to V4 Proxy API
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

test.describe("V4 Proxy API — Lifecycle", () => {
  // ── GKO-71: Deploy with syncFrom Management ──────────────────

  test(`Deploy V4 Proxy API with syncFrom Management ${XRAY.API_LIFECYCLE.DEPLOY_V4_SYNC_FROM_MGMT} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
    gateway,
  }) => {
    const API_NAME = "e2e-v4-sync-mgmt";
    const API_PATH = "/e2e-v4-sync-mgmt";

    await test.step("Apply CRD with syncFrom Management", async () => {
      await kubectl.apply(fixture("crds/api-v4-definitions/v4-proxy-api-sync-from-mgmt.yaml"));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API is STARTED and reachable on gateway", async () => {
      await mapi.assertApiMatches(apiId, { name: API_NAME, state: "STARTED" });
      await gateway.assertResponds(API_PATH, { status: 200 });
    });

    await kubectl.del(fixture("crds/api-v4-definitions/v4-proxy-api-sync-from-mgmt.yaml"));
  });

  // ── GKO-140: Delete a V4 API ─────────────────────────────────

  test(`Delete a V4 API ${XRAY.API_LIFECYCLE.DELETE_V4_API} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
    gateway,
  }) => {
    const API_NAME = "e2e-v4-sync-mgmt";
    const API_PATH = "/e2e-v4-sync-mgmt";
    const fixturePath = fixture("crds/api-v4-definitions/v4-proxy-api-sync-from-mgmt.yaml");

    await test.step("Deploy API", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("Verify API exists and is reachable", async () => {
      await mapi.assertApiStarted(apiId);
      await gateway.assertResponds(API_PATH, { status: 200 });
    });

    await test.step("Delete the CRD", async () => {
      await kubectl.del(fixturePath);
      await kubectl.waitForDeletion("apiv4definition", API_NAME);
    });

    await test.step("API is gone from APIM", async () => {
      await mapi.assertApiHttpStatus(apiId, 404);
    });

    await test.step("Gateway no longer responds for the path", async () => {
      await gateway.assertResponds(API_PATH, { status: 404 });
    });
  });

  // ── GKO-165: Missing required fields ─────────────────────────

  test(`Create V4 proxy API with missing required fields ${XRAY.API_LIFECYCLE.MISSING_REQUIRED_FIELDS_V4_PROXY} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(
      fixture("crds/api-v4-definitions/v4-proxy-api-invalid.yaml"),
    );
    expect(stderr.toLowerCase()).toMatch(/denied|rejected|invalid|error/);
  });

  // ── GKO-212: Re-deploy after delete ──────────────────────────

  test(`API is re-deployed when applying the same CRD after a delete ${XRAY.API_LIFECYCLE.REDEPLOY_AFTER_DELETE} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
    gateway,
  }) => {
    const API_NAME = "e2e-v4-sync-mgmt";
    const API_PATH = "/e2e-v4-sync-mgmt";
    const fixturePath = fixture("crds/api-v4-definitions/v4-proxy-api-sync-from-mgmt.yaml");

    await test.step("Deploy, verify, then delete", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
      await gateway.assertResponds(API_PATH, { status: 200 });
      await kubectl.del(fixturePath);
      await kubectl.waitForDeletion("apiv4definition", API_NAME);
    });

    await test.step("Re-deploy the same CRD", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API is reachable again", async () => {
      await mapi.assertApiStarted(apiId);
      await gateway.assertResponds(API_PATH, { status: 200 });
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-469: Context path already exists ─────────────────────

  test(`Context path conflict is rejected ${XRAY.API_LIFECYCLE.CONTEXT_PATH_CONFLICT_V4} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const fixturePath1 = fixture("crds/api-v4-definitions/v4-proxy-api-sync-from-mgmt.yaml");
    const fixturePath2 = fixture("crds/api-v4-definitions/v4-proxy-api-conflict-path.yaml");

    await test.step("Deploy first API", async () => {
      await kubectl.apply(fixturePath1);
      await kubectl.waitForCondition("apiv4definition", "e2e-v4-sync-mgmt", "Accepted");
    });

    await test.step("Second API with same path is rejected", async () => {
      const stderr = await kubectl.applyExpectFailure(fixturePath2);
      expect(stderr.toLowerCase()).toContain("context path");
    });

    await kubectl.del(fixturePath1);
    // Clean up the failed one too in case it was partially created
    await kubectl.del(fixturePath2);
  });

  // ── GKO-502: No plans + STARTED → error ─────────────────────

  test(`V4 API with no plans and STARTED state fails ${XRAY.API_LIFECYCLE.NO_PLANS_STARTED_V4} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyExpectFailure(
      fixture("crds/api-v4-definitions/v4-proxy-api-no-plans-started.yaml"),
    );
    expect(stderr.toLowerCase()).toContain("plan");
  });

  // ── GKO-503: No plans + STOPPED → OK ────────────────────────

  test(`V4 API with no plans and STOPPED state succeeds ${XRAY.API_LIFECYCLE.NO_PLANS_STOPPED_V4} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-no-plans-stopped";
    const fixturePath = fixture("crds/api-v4-definitions/v4-proxy-api-no-plans-stopped.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    await mapi.assertApiStopped(status.id);

    await kubectl.del(fixturePath);
  });

  // ── GKO-859: Failover configuration ─────────────────────────

  test(`V4 Proxy API with failover configuration ${XRAY.API_LIFECYCLE.FAILOVER_V4_PROXY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
    gateway,
  }) => {
    const API_NAME = "e2e-v4-failover";
    const API_PATH = "/e2e-v4-failover";
    const fixturePath = fixture("crds/api-v4-definitions/v4-proxy-api-with-failover.yaml");

    await test.step("Deploy API with failover", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API has failover config in APIM", async () => {
      await mapi.assertApiMatches(apiId, {
        name: API_NAME,
        state: "STARTED",
        failover: {
          enabled: true,
          maxRetries: 3,
        },
      });
    });

    await test.step("API is reachable on gateway", async () => {
      await gateway.assertResponds(API_PATH, { status: 200 });
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-83: Deploy V4 API with labels ────────────────────────

  test(`Deploy V4 API with labels and categories ${XRAY.API_LIFECYCLE.DEPLOY_V4_WITH_LABELS_CATEGORIES} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-labels-cats";
    const fixturePath = fixture("crds/api-v4-definitions/v4-proxy-api-with-labels-categories.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const api = await mapi.fetchApi(status.id);

    expect(api.labels).toBeTruthy();
    expect(api.labels!.length).toBeGreaterThanOrEqual(2);
    expect(api.labels).toContain("e2e-label-1");
    expect(api.labels).toContain("e2e-label-2");

    await kubectl.del(fixturePath);
  });

  // ── GKO-141: Update a V4 message API ─────────────────────────

  test(`Update a V4 message API ${XRAY.API_LIFECYCLE.UPDATE_V4_MESSAGE_API} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-message-mock";
    const createFixture = fixture("crds/api-v4-definitions/v4-message-api-mock.yaml");
    const updateFixture = fixture("crds/api-v4-definitions/v4-message-api-mock-updated.yaml");

    await test.step("Deploy message API", async () => {
      await kubectl.apply(createFixture);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("Update message API", async () => {
      await kubectl.apply(updateFixture);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("Changes are reflected in APIM", async () => {
      await mapi.waitForApiMatches(apiId, {
        description: "E2E test: V4 Message API updated description",
        apiVersion: "2.0.0",
      });
    });

    await kubectl.del(updateFixture);
  });

  // ── GKO-142: Update V4 message API with missing fields ───────

  test(`Update V4 message API with missing fields is rejected ${XRAY.API_LIFECYCLE.UPDATE_V4_MESSAGE_MISSING_FIELDS} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-message-mock";
    const createFixture = fixture("crds/api-v4-definitions/v4-message-api-mock.yaml");

    await test.step("Deploy valid message API", async () => {
      await kubectl.apply(createFixture);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("Update with invalid CRD (missing listeners) is rejected", async () => {
      const stderr = await kubectl.applyExpectFailure(
        fixture("crds/api-v4-definitions/v4-proxy-api-invalid.yaml"),
      );
      expect(stderr.toLowerCase()).toMatch(/denied|rejected|invalid|error/);
    });

    await kubectl.del(createFixture);
  });

  // ── GKO-176: No-op when CRD unchanged ───────────────────────

  test(`No-op when CRD is reapplied without changes ${XRAY.API_LIFECYCLE.NO_DEPLOY_WHEN_NO_CHANGES} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-sync-mgmt";
    const fixturePath = fixture("crds/api-v4-definitions/v4-proxy-api-sync-from-mgmt.yaml");

    await test.step("First apply", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status1 = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);

    await test.step("Re-apply the same CRD", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status2 = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);

    await test.step("API ID and state remain unchanged", async () => {
      expect(status2.id).toBe(status1.id);
      await mapi.assertApiStarted(status2.id);
    });

    await kubectl.del(fixturePath);
  });
});
