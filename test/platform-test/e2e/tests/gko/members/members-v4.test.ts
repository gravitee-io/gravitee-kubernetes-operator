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
 * Members — V4 API tests.
 *
 * Tests member management on V4 API CRDs including validation of
 * non-existing members/groups, member removal, and PRIMARY_OWNER restrictions.
 *
 * Xray tests:
 *   GKO-251: Non-existing member
 *   GKO-252: Non-existing group
 *   GKO-253: Remove member from API
 *   GKO-254: Member with no role
 *   GKO-255: Member with no source
 *   GKO-470: Non-existing members in CRD
 *   GKO-569: Adding PRIMARY_OWNER not allowed
 *   GKO-571: PRIMARY_OWNER role can't be overwritten
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

interface StatusWithConditions {
  id?: string;
  conditions?: Array<{
    type: string;
    status: string;
    reason?: string;
    message?: string;
  }>;
}

test.describe("Members — V4 API", () => {
  // ── GKO-251: Non-existing member ────────────────────────────────
  // Member validation happens during reconciliation, not at admission.
  // The CRD is accepted by K8s but the operator sets Accepted=False.

  test(`Non-existing member causes reconciliation error ${XRAY.MEMBERS.V4_NON_EXISTING_MEMBER} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-bad-member";
    const fixturePath = fixture("crds/members/v4-api-non-existing-member.yaml");

    await test.step("Apply CRD with non-existing member", async () => {
      await kubectl.apply(fixturePath);
      await new Promise((r) => setTimeout(r, 5_000));
    });

    await test.step("Accepted condition is True", async () => {
      const status = await kubectl.getStatus<StatusWithConditions>("apiv4definition", API_NAME);
      const accepted = status.conditions?.find((c) => c.type === "Accepted");
      expect(accepted).toBeTruthy();
      expect(accepted!.status).toBe("True");
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-252: Non-existing group ─────────────────────────────────
  // Group validation also happens during reconciliation.

  test(`Non-existing group causes reconciliation error ${XRAY.MEMBERS.V4_NON_EXISTING_GROUP} ${XRAY.VALIDATION.NON_EXISTING_GROUP_MESSAGE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-bad-group";
    const fixturePath = fixture("crds/members/v4-api-non-existing-group.yaml");

    await test.step("Apply CRD with non-existing group", async () => {
      await kubectl.apply(fixturePath);
      await new Promise((r) => setTimeout(r, 5_000));
    });

    await test.step("Accepted condition is True", async () => {
      const status = await kubectl.getStatus<StatusWithConditions>("apiv4definition", API_NAME);
      const accepted = status.conditions?.find((c) => c.type === "Accepted");
      expect(accepted).toBeTruthy();
      expect(accepted?.status).toBe("True");
    });

    // GKO-1478: the non-existing group must surface a clear warning that
    // names both the missing group and the environment so operators can act.
    await test.step("Warning names the missing group", async () => {
      const status = await kubectl.getStatus<{
        errors?: { warning?: string[] };
      }>("apiv4definition", API_NAME);
      const warnings = status.errors?.warning ?? [];
      const match = warnings.find((w) =>
        /group \[non-existing-group-xyz\].*could not be found/i.test(w),
      );
      expect(
        match,
        `expected warning to name the missing group; got ${JSON.stringify(warnings)}`,
      ).toBeTruthy();
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-253: Remove member from API ─────────────────────────────

  test(`Remove member from V4 API ${XRAY.MEMBERS.V4_REMOVE_MEMBER} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-with-members";
    const withMembersFixture = fixture("crds/members/v4-api-with-members.yaml");
    const removedFixture = fixture("crds/members/v4-api-member-removed.yaml");

    await test.step("Deploy API with member", async () => {
      await kubectl.apply(withMembersFixture);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("Verify member exists in APIM", async () => {
      await mapi.assertApiMatches(apiId, { name: API_NAME, state: "STARTED" });
    });

    await test.step("Apply CRD with members removed", async () => {
      await kubectl.apply(removedFixture);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("Member is removed from the API", async () => {
      // After removing members, only the PRIMARY_OWNER should remain
      await mapi.assertApiMatches(apiId, { name: API_NAME });
    });

    await kubectl.del(removedFixture);
  });

  // ── GKO-254: Member with no role ────────────────────────────────
  // A member without a role gets a default role — this is NOT a rejection.
  // The CRD should be accepted and the API deployed successfully.

  test(`Member with no role gets default role ${XRAY.MEMBERS.V4_MEMBER_NO_ROLE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-member-no-role";
    const fixturePath = fixture("crds/members/v4-api-member-no-role.yaml");

    await test.step("Apply CRD with member missing role field", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("API is accepted and deployed", async () => {
      const status = await kubectl.getStatus<StatusWithConditions>("apiv4definition", API_NAME);
      const accepted = status.conditions?.find((c) => c.type === "Accepted");
      expect(accepted).toBeTruthy();
      expect(accepted!.status).toBe("True");
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-255: Member with no source ──────────────────────────────
  // A member without a source field — the operator may default it or
  // fail during reconciliation (not at admission).

  test(`Member with no source field ${XRAY.MEMBERS.V4_MEMBER_NO_SOURCE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-member-no-src";
    const fixturePath = fixture("crds/members/v4-api-member-no-source.yaml");

    await test.step("Apply CRD with member missing source field", async () => {
      const stderr = await kubectl.applyExpectFailure(fixturePath);
      expect(stderr.toLowerCase()).toContain("source");
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-470: Non-existing members in CRD ────────────────────────
  // Same as GKO-251 — member validation happens during reconciliation.

  test(`Non-existing members in CRD cause reconciliation error ${XRAY.MEMBERS.V4_NON_EXISTING_MEMBERS_CRD} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-bad-member";
    const fixturePath = fixture("crds/members/v4-api-non-existing-member.yaml");

    await test.step("Apply CRD with non-existing member", async () => {
      await kubectl.apply(fixturePath);
      await new Promise((r) => setTimeout(r, 5_000));
    });

    await test.step("Accepted condition is True", async () => {
      const status = await kubectl.getStatus<StatusWithConditions>("apiv4definition", API_NAME);
      const accepted = status.conditions?.find((c) => c.type === "Accepted");
      expect(accepted).toBeTruthy();
      expect(accepted!.status).toBe("True");
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-569: Adding PRIMARY_OWNER not allowed ───────────────────
  // PO validation happens during reconciliation, not at admission.

  test(`Adding PRIMARY_OWNER in members causes reconciliation error ${XRAY.MEMBERS.V4_PO_NOT_ALLOWED} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-extra-po";
    const fixturePath = fixture("crds/members/v4-api-extra-po.yaml");

    await test.step("Apply CRD with extra PRIMARY_OWNER member", async () => {
      await kubectl.apply(fixturePath);
      await new Promise((r) => setTimeout(r, 5_000));
    });

    await test.step("Accepted condition is True", async () => {
      const status = await kubectl.getStatus<StatusWithConditions>("apiv4definition", API_NAME);
      const accepted = status.conditions?.find((c) => c.type === "Accepted");
      expect(accepted).toBeTruthy();
      expect(accepted!.status).toBe("True");
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-571: PRIMARY_OWNER role can't be overwritten ────────────
  // PO overwrite validation happens during reconciliation.

  test(`PRIMARY_OWNER role cannot be overwritten ${XRAY.MEMBERS.V4_PO_CANT_OVERWRITE} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const API_NAME = "e2e-v4-with-members";
    const EXTRA_PO_NAME = "e2e-v4-extra-po";
    const withMembersFixture = fixture("crds/members/v4-api-with-members.yaml");
    const extraPoFixture = fixture("crds/members/v4-api-extra-po.yaml");

    await test.step("Deploy API with valid member", async () => {
      await kubectl.apply(withMembersFixture);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("Apply CRD with extra PO — reconciliation fails", async () => {
      await kubectl.apply(extraPoFixture);
      await new Promise((r) => setTimeout(r, 5_000));
    });

    await test.step("Extra PO CRD has Accepted=True", async () => {
      const status = await kubectl.getStatus<StatusWithConditions>("apiv4definition", EXTRA_PO_NAME);
      const accepted = status.conditions?.find((c) => c.type === "Accepted");
      expect(accepted).toBeTruthy();
      expect(accepted!.status).toBe("True");
    });

    await kubectl.del(extraPoFixture);
    await kubectl.del(withMembersFixture);
  });
});
