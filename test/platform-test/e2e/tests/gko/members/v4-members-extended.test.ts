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
 * V4 API Members — Extended scenarios.
 *
 * Xray tests:
 *   GKO-244:  PrimaryOwner explicitly defined in CRD
 *   GKO-256:  Create V4 API with non-existing group
 *   GKO-306:  Primary owner via management-context user
 *   GKO-307:  Transfer primary owner
 *   GKO-658:  Take over primary owner via management-context user
 *
 * Granting, re-roling and revoking a member (GKO-213/247/253/259/402) is the
 * shared journey tests/user-journeys/platform-admin/manage-api-members, and
 * associating groups (GKO-257/314/1004) is associate-groups-with-an-api.
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 *   - Xray provisions an "e2e-group-with-member" group
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS, PROVISIONER } from "../../../helpers/tags.js";
import * as kubectl from "../../../helpers/kubectl.js";

interface StatusWithConditions {
  id?: string;
  conditions?: Array<{
    type: string;
    status: string;
    reason?: string;
    message?: string;
  }>;
}

const NON_EXISTING_GROUP = "members/v4-api-non-existing-group/crd.yaml";
const EXTRA_PO = "members/v4-api-extra-po/crd.yaml";

function acceptedTrue(status: StatusWithConditions): boolean {
  return status.conditions?.find((c) => c.type === "Accepted")?.status === "True";
}

const SYNC_FROM_MGMT = "crds/api-v4-definitions/v4-proxy-api-sync-from-mgmt.yaml";

test.describe(`V4 API Members — Extended ${PROVISIONER.GKO}`, () => {
  // Safety-net cleanup: runs even if a test times out before its inline
  // cleanup. Each del() ignores errors (the resource may already be gone).
  // e2e-v4-sync-mgmt is shared with other files, so a leak here cascades.
  test.afterEach(async () => {
    for (const f of [NON_EXISTING_GROUP, EXTRA_PO, SYNC_FROM_MGMT]) {
      await kubectl.del(fixture(f)).catch(() => {});
    }
  });

  // ── GKO-244: PrimaryOwner defined in CRD ─────────────────────
  // Declaring PO in members is rejected by the webhook in current GKO — this
  // is covered by GKO-569. The authoritative path is to set PO via the
  // management-context user (see GKO-306/658).

  test(`PrimaryOwner defined in CRD via members section ${XRAY.MEMBERS.V4_PO_DEFINED_IN_CRD} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const fixturePath = fixture(EXTRA_PO);

    // Extra-PO fixture carries a PRIMARY_OWNER entry in members — operator
    // accepts the CRD (validation happens at reconciliation, not admission).
    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", "e2e-v4-extra-po", "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("apiv4definition", "e2e-v4-extra-po");
    expect(acceptedTrue(status)).toBe(true);

    await kubectl.del(fixturePath);
  });

  // ── GKO-256: Create with non-existing group ──────────────────
  // Group validation happens during reconciliation; CRD should still be
  // Accepted (operator ignores unknown groups).

  test(`Create V4 API with non-existing group ${XRAY.MEMBERS.V4_CREATE_NON_EXISTING_GROUP} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-bad-group";
    const fixturePath = fixture(NON_EXISTING_GROUP);

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("apiv4definition", API_NAME);
    expect(acceptedTrue(status)).toBe(true);

    await kubectl.del(fixturePath);
  });

  // ── GKO-306: Primary owner via management-context user ──────
  // The management-context configured in dev-ctx is itself the primary owner
  // of every API the operator deploys. Verify that the sync-from-mgmt fixture
  // resolves to a non-empty primaryOwner.

  test(`Primary owner resolved via management-context user ${XRAY.MEMBERS.V4_PO_VIA_MGMT_CONTEXT} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-sync-mgmt";
    const fixturePath = fixture("api-v4-definitions/sync-from-mgmt/crd.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const api = await mapi.fetchApi(status.id);
    expect(api.primaryOwner).toBeDefined();
    expect(api.primaryOwner.id).toBeTruthy();

    await kubectl.del(fixturePath);
  });

  // ── GKO-307: Transfer primary owner ─────────────────────────
  // Transferring the PO is done by updating the ManagementContext credentials
  // (out of scope). This test asserts the operator deploys stably when the
  // CRD is reapplied — the PO transfer scenario is covered by GKO-658.

  // FIXME(GKO-2940): re-applying a V4 CR that omits the endpoint `secondary`
  // field (post-GKO-2857) churns metadata.generation while the Accepted
  // condition's observedGeneration lags, so kubectl wait (observedGeneration-
  // aware in kubectl >=1.31) blocks. Re-enable once GKO-2940 is fixed.
  test.fixme(`Primary owner stable across re-apply ${XRAY.MEMBERS.V4_TRANSFER_PRIMARY_OWNER} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-sync-mgmt";
    const fixturePath = fixture("api-v4-definitions/sync-from-mgmt/crd.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    const poBefore = (await mapi.fetchApi(
      (await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME)).id,
    )).primaryOwner.id;

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    const poAfter = (await mapi.fetchApi(
      (await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME)).id,
    )).primaryOwner.id;

    expect(poAfter).toBe(poBefore);

    await kubectl.del(fixturePath);
  });

  // ── GKO-658: Take over PO via management-context user ──────
  // Variant of GKO-306 — re-applying the CRD must not reset or remove the
  // primary owner even after an intervening delete + re-create cycle.

  test(`Take over primary owner via mgmt-context user ${XRAY.MEMBERS.V4_TAKE_OVER_PO_VIA_MGMT_CTX} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-sync-mgmt";
    const fixturePath = fixture("api-v4-definitions/sync-from-mgmt/crd.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    await kubectl.del(fixturePath);
    await kubectl.waitForDeletion("apiv4definition", API_NAME);

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const api = await mapi.fetchApi(status.id);
    expect(api.primaryOwner.id).toBeTruthy();

    await kubectl.del(fixturePath);
  });

});
