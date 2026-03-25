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
 * Reconciliation Status & Conditions tests.
 *
 * Xray tests:
 *   GKO-1387: Accepted condition set to False when reconciliation fails
 *   GKO-1388: Accepted condition is not False on successful reconciliation
 *   GKO-1389: Accepted condition updates on configuration changes
 *   GKO-1390: processingStatus is still present in the status field
 *   GKO-1392: ResolvedRefs condition is no longer present
 *   GKO-1445: Idempotent reconciliation
 *   GKO-1446: Status conditions reflect reconciliation state
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

interface StatusWithConditions {
  id: string;
  processingStatus?: string;
  conditions?: Array<{
    type: string;
    status: string;
    reason?: string;
    message?: string;
    lastTransitionTime?: string;
  }>;
}

test.describe("Reconciliation — Status & Conditions", () => {
  // ── GKO-1388: Accepted not False on success ──────────────────

  test(`Accepted condition not False on successful reconciliation ${XRAY.DEPLOYMENT_RECONCILIATION.ACCEPTED_NOT_FALSE_ON_SUCCESS} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-start-stop";
    const fixturePath = fixture("crds/api-v4-definitions/v4-proxy-api-started.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("apiv4definition", API_NAME);

    const acceptedCondition = status.conditions?.find((c) => c.type === "Accepted");
    expect(acceptedCondition).toBeTruthy();
    expect(acceptedCondition!.status).toBe("True");

    await kubectl.del(fixturePath);
  });

  // ── GKO-1390: processingStatus present ───────────────────────

  test(`processingStatus is still present in status ${XRAY.DEPLOYMENT_RECONCILIATION.PROCESSING_STATUS_PRESENT} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-start-stop";
    const fixturePath = fixture("crds/api-v4-definitions/v4-proxy-api-started.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("apiv4definition", API_NAME);
    expect(status.processingStatus).toBeTruthy();

    await kubectl.del(fixturePath);
  });

  // ── GKO-1446: Conditions reflect reconciliation state ────────

  test(`Status conditions reflect reconciliation state ${XRAY.DEPLOYMENT_RECONCILIATION.STATUS_CONDITIONS_REFLECT_STATE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-start-stop";
    const fixturePath = fixture("crds/api-v4-definitions/v4-proxy-api-started.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("apiv4definition", API_NAME);

    await test.step("Conditions are populated", async () => {
      expect(status.conditions).toBeTruthy();
      expect(status.conditions!.length).toBeGreaterThan(0);
    });

    await test.step("Accepted condition has required fields", async () => {
      const accepted = status.conditions!.find((c) => c.type === "Accepted");
      expect(accepted).toBeTruthy();
      expect(accepted!.status).toBeTruthy();
      expect(accepted!.lastTransitionTime).toBeTruthy();
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-1445: Idempotent reconciliation ──────────────────────

  test(`Idempotent reconciliation for repeated apply ${XRAY.DEPLOYMENT_RECONCILIATION.IDEMPOTENT_RECONCILIATION} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-start-stop";
    const fixturePath = fixture("crds/api-v4-definitions/v4-proxy-api-started.yaml");

    await test.step("First apply", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status1 = await kubectl.getStatus<StatusWithConditions>("apiv4definition", API_NAME);

    await test.step("Re-apply the same CRD", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status2 = await kubectl.getStatus<StatusWithConditions>("apiv4definition", API_NAME);

    await test.step("API ID is unchanged", async () => {
      expect(status2.id).toBe(status1.id);
    });

    await test.step("API state is still STARTED", async () => {
      await mapi.assertApiStarted(status2.id);
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-1387: Accepted=False on failure ──────────────────────

  test(`Accepted condition is False when reconciliation fails ${XRAY.DEPLOYMENT_RECONCILIATION.ACCEPTED_FALSE_ON_FAILURE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-bad-endpoint";
    const fixturePath = fixture("crds/api-v4-definitions/v4-proxy-api-bad-endpoint-type.yaml");

    await test.step("Apply CRD that passes webhook but fails reconciliation", async () => {
      await kubectl.apply(fixturePath);
      // Don't waitForCondition("Accepted") — it will never become True
      await new Promise((r) => setTimeout(r, 5_000));
    });

    await test.step("Accepted condition is False with error reason", async () => {
      const status = await kubectl.getStatus<StatusWithConditions>("apiv4definition", API_NAME);
      const accepted = status.conditions?.find((c) => c.type === "Accepted");
      expect(accepted).toBeTruthy();
      expect(accepted!.status).toBe("False");
      expect(accepted!.reason).toBe("ControlPlaneError");
      expect(accepted!.message).toBeTruthy();
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-1389: Accepted updates on config changes ─────────────

  test(`Accepted condition updates on configuration changes ${XRAY.DEPLOYMENT_RECONCILIATION.ACCEPTED_UPDATES_ON_CHANGE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-reconcile";
    const createFixture = fixture("crds/api-v4-definitions/v4-proxy-api-reconcile.yaml");
    const updateFixture = fixture("crds/api-v4-definitions/v4-proxy-api-reconcile-updated.yaml");

    await test.step("Apply initial CRD", async () => {
      await kubectl.apply(createFixture);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status1 = await kubectl.getStatus<StatusWithConditions>("apiv4definition", API_NAME);
    const accepted1 = status1.conditions?.find((c) => c.type === "Accepted");
    expect(accepted1?.status).toBe("True");

    await test.step("Apply updated CRD", async () => {
      await kubectl.apply(updateFixture);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status2 = await kubectl.getStatus<StatusWithConditions>("apiv4definition", API_NAME);
    const accepted2 = status2.conditions?.find((c) => c.type === "Accepted");
    expect(accepted2?.status).toBe("True");

    await kubectl.del(updateFixture);
  });

  // ── GKO-1392: ResolvedRefs condition ─────────────────────────

  // GKO-1392 is forward-looking: ResolvedRefs is deprecated and will be removed from status
  // in a future GKO release. Re-enable once the operator no longer sets this condition.
  test.fixme(`ResolvedRefs condition is not present in status ${XRAY.DEPLOYMENT_RECONCILIATION.RESOLVED_REFS_NOT_PRESENT} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-start-stop";
    const fixturePath = fixture("crds/api-v4-definitions/v4-proxy-api-started.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("apiv4definition", API_NAME);

    // GKO-1392: ResolvedRefs condition is no longer present in status
    const resolvedRefs = status.conditions?.find((c) => c.type === "ResolvedRefs");
    expect(resolvedRefs).toBeUndefined();

    await kubectl.del(fixturePath);
  });
});
