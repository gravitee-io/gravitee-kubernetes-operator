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
 * V2 API Lifecycle & Members.
 *
 * Xray tests:
 *   GKO-1065: Update API path on a V2 CRD
 *   GKO-202:  Add member with role name to a V2 API
 *   GKO-204:  Add member with non-existing role to a V2 API
 *   GKO-205:  Create a V2 API with a non-existing member
 *   GKO-207:  Create a V2 API with a non-existing group
 *   GKO-208:  Create a V2 API with an existing group
 *   GKO-216:  Remove member from a V2 API
 *   GKO-258:  PRIMARY_OWNER role in members section
 *   GKO-308:  Change member role of a V2 API
 *   GKO-393:  Add member without role to a V2 API (defaults applied)
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 *   - Xray provisions an "e2e-group-with-member" group
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

interface StatusWithConditions {
  id?: string;
  conditions?: Array<{
    type: string;
    status: string;
  }>;
}

function acceptedTrue(status: StatusWithConditions): boolean {
  return status.conditions?.find((c) => c.type === "Accepted")?.status === "True";
}

test.describe("V2 API Members — Extended", () => {
  // ── GKO-1065: Update API path ────────────────────────────────

  test(`Update V2 API virtual_host path ${XRAY.V2_API_LIFECYCLE.V2_UPDATE_API_PATH} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-export";
    const base = fixture("crds/import-export/v2-api-export.yaml");
    const updated = fixture("crds/api-definitions/v2-api-updated-path.yaml");

    await kubectl.apply(base);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;

    await kubectl.apply(updated);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    await mapi.waitForApiMatches(apiId, { name: API_NAME });

    await kubectl.del(updated);
  });

  // ── GKO-202: Add member with role name ───────────────────────

  test(`Add member with role name to V2 API ${XRAY.MEMBERS.V2_ADD_MEMBER_WITH_ROLE} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-with-members";
    const fixturePath = fixture("crds/members/v2-api-with-members.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    const apiId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;
    await mapi.waitForApiMatches(apiId, { name: API_NAME });

    await kubectl.del(fixturePath);
  });

  // ── GKO-204: Non-existing role ───────────────────────────────

  test(`V2 API with non-existing role ${XRAY.MEMBERS.V2_NON_EXISTING_ROLE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v2-bad-role";
    const fixturePath = fixture("crds/members/v2-api-non-existing-role.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("apidefinition", API_NAME);
    expect(acceptedTrue(status)).toBe(true);

    await kubectl.del(fixturePath);
  });

  // ── GKO-205: Non-existing member ─────────────────────────────

  test(`V2 API with non-existing member ${XRAY.MEMBERS.V2_NON_EXISTING_MEMBER} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v2-bad-member";
    const fixturePath = fixture("crds/members/v2-api-non-existing-member.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("apidefinition", API_NAME);
    expect(acceptedTrue(status)).toBe(true);

    await kubectl.del(fixturePath);
  });

  // ── GKO-207: Non-existing group ──────────────────────────────

  test(`V2 API with non-existing group ${XRAY.MEMBERS.V2_NON_EXISTING_GROUP} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v2-bad-group";
    const fixturePath = fixture("crds/members/v2-api-non-existing-group.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("apidefinition", API_NAME);
    expect(acceptedTrue(status)).toBe(true);

    await kubectl.del(fixturePath);
  });

  // ── GKO-208: Existing group ──────────────────────────────────

  test(`V2 API with existing group ${XRAY.MEMBERS.V2_EXISTING_GROUP} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-with-groups";
    const fixturePath = fixture("crds/members/v2-api-with-groups.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    const apiId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;
    await mapi.waitForApiMatches(apiId, { name: API_NAME });

    await kubectl.del(fixturePath);
  });

  // ── GKO-216: Remove member from V2 API ───────────────────────

  test(`Remove member from V2 API ${XRAY.MEMBERS.V2_REMOVE_MEMBER} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-with-members";
    const withMembers = fixture("crds/members/v2-api-with-members.yaml");
    const removed = fixture("crds/members/v2-api-member-removed.yaml");

    await kubectl.apply(withMembers);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;

    await kubectl.apply(removed);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    await mapi.waitForApiMatches(apiId, { name: API_NAME });

    await kubectl.del(removed);
  });

  // ── GKO-258: PRIMARY_OWNER in members section ───────────────
  // Declaring a second PRIMARY_OWNER in members is rejected by the webhook —
  // variant: a single PO entry (with the same source as the mgmt-ctx user)
  // should go through reconciliation. Error-path is covered by GKO-569.

  test(`V2 API with PRIMARY_OWNER member role is accepted ${XRAY.MEMBERS.V2_PO_IN_MEMBERS} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v2-po-member";
    const fixturePath = fixture("crds/members/v2-api-po-in-members.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("apidefinition", API_NAME);
    expect(acceptedTrue(status)).toBe(true);

    await kubectl.del(fixturePath);
  });

  // ── GKO-308: Change member role ──────────────────────────────

  test(`Change member role of V2 API ${XRAY.MEMBERS.V2_CHANGE_MEMBER_ROLE} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-with-members";
    const withMembers = fixture("crds/members/v2-api-with-members.yaml");
    const reviewer = fixture("crds/members/v2-api-member-reviewer.yaml");

    await kubectl.apply(withMembers);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;

    await kubectl.apply(reviewer);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    await mapi.waitForApiMatches(apiId, { name: API_NAME });

    await kubectl.del(reviewer);
  });

  // ── GKO-393: Member without role ─────────────────────────────

  test(`V2 API with member missing role field ${XRAY.MEMBERS.V2_MEMBER_NO_ROLE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v2-member-no-role";
    const fixturePath = fixture("crds/members/v2-api-member-no-role.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("apidefinition", API_NAME);
    expect(acceptedTrue(status)).toBe(true);

    await kubectl.del(fixturePath);
  });
});
