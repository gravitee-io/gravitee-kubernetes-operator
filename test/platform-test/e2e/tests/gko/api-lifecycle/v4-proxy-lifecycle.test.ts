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
 * V4 Proxy API deployment & reconciliation, through the operator.
 *
 * Xray tests:
 *   GKO-71:  Deploy V4 Proxy API with syncFrom Management
 *   GKO-176: Should not deploy when no changes are made to V4 CRD
 *   GKO-212: API is re-deployed when applying the same CRD after a delete
 *
 * What a customer reaches through Terraform too moved to shared journeys:
 * deleting an API and the visibility/lifecycle matrix (GKO-140) to
 * configure-visibility-and-lifecycle, failover (GKO-859) to
 * configure-endpoint-failover, labels (GKO-83) to label-an-api, the message-API
 * update (GKO-141) to publish-a-message-api. The admission rejections
 * (GKO-165/469/502/503/142) moved to tests/gko/admission-webhook.
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS, PROVISIONER } from "../../../helpers/tags.js";
import * as kubectl from "../../../helpers/kubectl.js";

test.describe(`V4 Proxy API — Deployment & reconciliation ${PROVISIONER.GKO}`, () => {
  // Safety-net cleanup: runs even if a test times out before its inline
  // cleanup. The del() ignores errors (the resource may already be gone).
  // This name is shared across test files, so a leak here cascades into
  // unrelated suites — always remove it.
  test.afterEach(async () => {
    await kubectl.del(fixture("api-v4-definitions/sync-from-mgmt/crd.yaml")).catch(() => {});
  });

  // ── GKO-71: Deploy with syncFrom Management ──────────────────

  test(`Deploy V4 Proxy API with syncFrom Management ${XRAY.API_LIFECYCLE.DEPLOY_V4_SYNC_FROM_MGMT} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
    gateway,
  }) => {
    const API_NAME = "e2e-v4-sync-mgmt";
    const API_PATH = "/e2e-v4-sync-mgmt";

    await test.step("Apply CRD with syncFrom Management", async () => {
      await kubectl.apply(fixture("api-v4-definitions/sync-from-mgmt/crd.yaml"));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API is STARTED and reachable on gateway", async () => {
      await mapi.assertApiMatches(apiId, { name: API_NAME, state: "STARTED" });
      await gateway.assertResponds(API_PATH, { status: 200 });
    });

    await kubectl.del(fixture("api-v4-definitions/sync-from-mgmt/crd.yaml"));
  });

  // ── GKO-212: Re-deploy after delete ──────────────────────────

  test(`API is re-deployed when applying the same CRD after a delete ${XRAY.API_LIFECYCLE.REDEPLOY_AFTER_DELETE} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
    gateway,
  }) => {
    const API_NAME = "e2e-v4-sync-mgmt";
    const API_PATH = "/e2e-v4-sync-mgmt";
    const fixturePath = fixture("api-v4-definitions/sync-from-mgmt/crd.yaml");

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

  // ── GKO-176: No-op when CRD unchanged ───────────────────────

  // FIXME(GKO-2940): re-applying a V4 CR that omits the endpoint `secondary`
  // field (post-GKO-2857) churns metadata.generation while the Accepted
  // condition's observedGeneration lags, so kubectl wait (observedGeneration-
  // aware in kubectl >=1.31) blocks. Re-enable once GKO-2940 is fixed.
  test.fixme(`No-op when CRD is reapplied without changes ${XRAY.API_LIFECYCLE.NO_DEPLOY_WHEN_NO_CHANGES} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-sync-mgmt";
    const fixturePath = fixture("api-v4-definitions/sync-from-mgmt/crd.yaml");

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
