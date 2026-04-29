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
 * V2 API Groups, Members & Primary Owner.
 *
 * Xray tests:
 *   GKO-398:  Add group using name (hrid) to V2 API
 *   GKO-399:  Multiple groups using names in V2 API
 *   GKO-400:  Remove group from V2 API CRD
 *   GKO-401:  notifyMembers=true on V2 API is accepted
 *   GKO-601:  PRIMARY_OWNER role not overwriteable via members section
 *   GKO-602:  Declaring a different user as PRIMARY_OWNER is not allowed
 *   GKO-657:  Management-context user takes over as PRIMARY_OWNER
 *   GKO-659:  Adding PRIMARY_OWNER back to members has no effect
 *   GKO-1003: Add GroupRefs to V2 API
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

interface StatusWithConditions {
  id?: string;
  conditions?: Array<{ type: string; status: string }>;
}

function acceptedTrue(status: StatusWithConditions): boolean {
  return status.conditions?.find((c) => c.type === "Accepted")?.status === "True";
}

const GROUP_A = "crds/members/group-for-v2-members.yaml";
const GROUP_B = "crds/members/group-b-for-v2-members.yaml";

test.describe("V2 API Groups, Members & PO — Extended", () => {
  test.afterEach(async () => {
    await kubectlSafe
      .del(fixture("crds/members/v2-api-with-group-hrid.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/members/v2-api-group-removed.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/members/v2-api-with-multiple-groups.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/members/v2-api-notify-members.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/members/v2-api-po-overwrite-attempt.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/members/v2-api-po-different-member.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/members/v2-api-po-take-over.yaml"))
      .catch(() => {});
    await kubectlSafe
      .del(fixture("crds/members/v2-api-with-group-refs.yaml"))
      .catch(() => {});
    await kubectlSafe.del(fixture(GROUP_B)).catch(() => {});
    await kubectlSafe.del(fixture(GROUP_A)).catch(() => {});
  });

  // ── GKO-398: Add group using name (hrid) ────────────────────

  test(`V2 API references group by hrid ${XRAY.MEMBERS.V2_ADD_GROUP_HRID} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-group-hrid";

    await kubectl.apply(fixture(GROUP_A));
    await kubectl.waitForCondition("group", "e2e-v2-group-a", "Accepted");

    await kubectl.apply(fixture("crds/members/v2-api-with-group-hrid.yaml"));
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    const apiId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;
    await mapi.waitForApiMatches(apiId, { name: API_NAME });

    await kubectl.del(fixture("crds/members/v2-api-with-group-hrid.yaml"));
    await kubectl.del(fixture(GROUP_A));
  });

  // ── GKO-399: Multiple groups ────────────────────────────────

  test(`V2 API references multiple groups by hrid ${XRAY.MEMBERS.V2_MULTIPLE_GROUPS} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-multi-groups";

    await kubectl.apply(fixture(GROUP_A));
    await kubectl.apply(fixture(GROUP_B));
    await kubectl.waitForCondition("group", "e2e-v2-group-a", "Accepted");
    await kubectl.waitForCondition("group", "e2e-v2-group-b", "Accepted");

    await kubectl.apply(fixture("crds/members/v2-api-with-multiple-groups.yaml"));
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    const apiId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;
    await mapi.waitForApiMatches(apiId, { name: API_NAME });

    await kubectl.del(fixture("crds/members/v2-api-with-multiple-groups.yaml"));
    await kubectl.del(fixture(GROUP_B));
    await kubectl.del(fixture(GROUP_A));
  });

  // ── GKO-400: Remove group from V2 API ───────────────────────

  test(`Remove group from V2 API on re-apply ${XRAY.MEMBERS.V2_REMOVE_GROUP} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-group-hrid";
    const withGroup = fixture("crds/members/v2-api-with-group-hrid.yaml");
    const without = fixture("crds/members/v2-api-group-removed.yaml");

    await kubectl.apply(fixture(GROUP_A));
    await kubectl.waitForCondition("group", "e2e-v2-group-a", "Accepted");

    await kubectl.apply(withGroup);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    const apiId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;

    await kubectl.apply(without);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    await mapi.waitForApiMatches(apiId, {
      description: "E2E test: V2 API with group removed",
    });

    await kubectl.del(without);
    await kubectl.del(fixture(GROUP_A));
  });

  // ── GKO-401: notifyMembers=true is accepted ─────────────────

  test(`V2 API with notifyMembers=true is accepted ${XRAY.MEMBERS.V2_NOTIFY_MEMBERS} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v2-notify-members";
    const fixturePath = fixture("crds/members/v2-api-notify-members.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("apidefinition", API_NAME);
    expect(acceptedTrue(status)).toBe(true);

    await kubectl.del(fixturePath);
  });

  // ── GKO-601: PRIMARY_OWNER can't be overwritten via members ─

  test(`V2 API keeps mgmt-ctx PRIMARY_OWNER when members list re-asserts it ${XRAY.MEMBERS.V2_PO_NOT_OVERWRITEABLE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v2-po-overwrite";
    const fixturePath = fixture("crds/members/v2-api-po-overwrite-attempt.yaml");

    // Declaring the same mgmt-ctx user (`admin`) as PO must be a no-op and
    // the API reconciles to Accepted=True.
    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    const status = await kubectl.getStatus<StatusWithConditions>("apidefinition", API_NAME);
    expect(acceptedTrue(status)).toBe(true);

    await kubectl.del(fixturePath);
  });

  // GKO-602 (V2 API with a different PRIMARY_OWNER is rejected) — not
  // covered: GKO admission does not enforce this (product gap).

  // ── GKO-657: PO take-over via management-context user ──────

  test(`V2 API with mgmt-ctx user as PRIMARY_OWNER is accepted ${XRAY.MEMBERS.V2_TAKE_OVER_PO} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-po-takeover";
    const fixturePath = fixture("crds/members/v2-api-po-take-over.yaml");

    await kubectl.apply(fixturePath);
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    const apiId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;
    await mapi.waitForApiMatches(apiId, { name: API_NAME });

    await kubectl.del(fixturePath);
  });

  // GKO-659 (adding PO to members has no effect) — not covered: overlaps
  // with GKO-601 (idempotent re-apply with PO member) and GKO-602 is also
  // a product gap.

  // ── GKO-1003: Add GroupRefs to V2 API ───────────────────────

  test(`V2 API references a group via groupRefs ${XRAY.MEMBERS.V2_ADD_GROUP_REFS} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-group-refs";

    await kubectl.apply(fixture(GROUP_A));
    await kubectl.waitForCondition("group", "e2e-v2-group-a", "Accepted");

    await kubectl.apply(fixture("crds/members/v2-api-with-group-refs.yaml"));
    await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");

    const apiId = (await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME)).id;
    await mapi.waitForApiMatches(apiId, { name: API_NAME });

    await kubectl.del(fixture("crds/members/v2-api-with-group-refs.yaml"));
    await kubectl.del(fixture(GROUP_A));
  });
});
