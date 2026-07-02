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
 * apim_group resource + data source driven by the Gravitee Terraform Provider.
 *
 * Covers GKO-2564 ("[TF] – Manage groups"), built on top of GKO-2551 which
 * added the Automation API group endpoints. The Terraform provider writes
 * groups through the Automation API, so every group created here lands in
 * APIM with `origin: KUBERNETES`. Verification reads back through the v1
 * management API (see mapi.listGroups / fetchGroupMembers).
 *
 * Note on members: the only pre-existing memory-source user is `admin`, the
 * org primary owner — and APIM rejects adding the primary owner as a group
 * member (it surfaces as an HTTP 500). The member test therefore provisions a
 * dedicated `gravitee`-source service account, a non-PO user that resolves
 * cleanly. See mapi.createServiceAccount.
 *
 * Preconditions:
 *   - APIM is running
 *   - terraform CLI is installed (>= 1.3 for optional() in object types)
 *   - the locally-built provider mirror exposes the apim_group resource
 *     (scripts/build-tf-provider.sh)
 */

import { test, expect } from "../../setup.js";
import { XRAY, TAGS } from "../../helpers/tags.js";
import * as terraform from "../../helpers/terraform.js";
import type { TfWorkspace } from "../../helpers/terraform.js";

// ── Resource lifecycle (shared workspace, read-only assertions) ──────

test.describe("Terraform — Groups · Resource lifecycle @since-4.12", () => {
  let ws: TfWorkspace;
  let groupHrid: string;

  test.beforeAll(async () => {
    // init + apply + output are sequential terraform invocations, each capped
    // at terraform.TF_TIMEOUT_MS. The hook timeout must exceed their combined
    // ceiling so terraform's own timeout fires first instead of Playwright
    // orphaning a running terraform process. See TF_WORKSPACE_TIMEOUT_MS.
    test.setTimeout(terraform.TF_WORKSPACE_TIMEOUT_MS);
    ws = await terraform.initWorkspace("groups/lifecycle");
    await terraform.apply(ws);
    groupHrid = await terraform.output(ws, "group_hrid");
  });

  test.afterAll(async () => {
    test.setTimeout(terraform.TF_WORKSPACE_TIMEOUT_MS);
    if (ws) await terraform.destroyWorkspace(ws);
  });

  // Moved to the terraform arm of the shared tests/user-journeys/create-group-with-member/create-group-with-member.scenario.ts
  // (create group -> lands in APIM with origin KUBERNETES), tagged @GKO-2865.

  test(`notify_members is propagated to APIM ${XRAY.TERRAFORM.GROUP_NOTIFY_MEMBERS} ${TAGS.REGRESSION}`, async ({
    mapi,
  }) => {
    // notify_members (TF) is the inverse of disable_membership_notifications
    // (APIM wire). The fixture default is notify_members = true.
    await mapi.waitForGroupMatches(groupHrid, { disable_membership_notifications: false });
  });

  test(`terraform plan shows no changes ${XRAY.TERRAFORM.GROUP_IDEMPOTENT} ${TAGS.REGRESSION}`, async () => {
    const result = await terraform.plan(ws);
    expect(result.hasChanges).toBe(false);
  });
});

// ── In-place update ─────────────────────────────────────────────────

test.describe("Terraform — Groups · Update @since-4.12", () => {
  test(`terraform apply updates a group in place ${XRAY.TERRAFORM.GROUP_UPDATE} ${TAGS.REGRESSION}`, async ({
    mapi,
  }) => {
    // Three sequential applies plus init/destroy — size the timeout above the
    // per-process TF_TIMEOUT_MS ceiling for that many calls.
    test.setTimeout(terraform.TF_TIMEOUT_MS * 9);

    const HRID = "e2e-tf-group-update";
    const ws = await terraform.initWorkspace("groups/lifecycle");
    try {
      let groupId: string;

      await test.step("Create the baseline group", async () => {
        await terraform.writeVars(ws, {
          hrid_suffix: "update",
          group_name: "e2e-tf-group-update",
          notify_members: true,
        });
        await terraform.apply(ws);
        groupId = await terraform.output(ws, "group_id");
        await mapi.waitForGroupMatches(HRID, { name: "e2e-tf-group-update" });
      });

      await test.step("Rename the group — same id, no replacement", async () => {
        await terraform.writeVars(ws, {
          hrid_suffix: "update",
          group_name: "e2e-tf-group-update-renamed",
          notify_members: true,
        });
        await terraform.apply(ws);
        // hrid is unchanged, so the deterministic group id must be stable —
        // the update is in-place, not a destroy + recreate.
        const idAfter = await terraform.output(ws, "group_id");
        expect(idAfter).toBe(groupId);
        await mapi.waitForGroupMatches(HRID, { name: "e2e-tf-group-update-renamed" });
      });

      await test.step("Flip notify_members to false", async () => {
        await terraform.writeVars(ws, {
          hrid_suffix: "update",
          group_name: "e2e-tf-group-update-renamed",
          notify_members: false,
        });
        await terraform.apply(ws);
        await mapi.waitForGroupMatches(HRID, { disable_membership_notifications: true });
      });
    } finally {
      await terraform.destroyWorkspace(ws);
    }
  });
});

// ── Destroy ─────────────────────────────────────────────────────────

test.describe("Terraform — Groups · Destroy @since-4.12", () => {
  test(`terraform destroy removes the group from APIM ${XRAY.TERRAFORM.GROUP_DESTROY} ${TAGS.REGRESSION}`, async ({
    mapi,
  }) => {
    test.setTimeout(terraform.TF_WORKSPACE_TIMEOUT_MS);

    const HRID = "e2e-tf-group-destroy";
    const ws = await terraform.initWorkspace("groups/lifecycle");
    try {
      await terraform.writeVars(ws, { hrid_suffix: "destroy", group_name: HRID });
      await terraform.apply(ws);
      await mapi.waitForGroupMatches(HRID, { name: HRID });

      await terraform.destroy(ws);
      await mapi.waitForGroupAbsent(HRID);
    } finally {
      // destroy already ran; destroyWorkspace re-runs it as a no-op and
      // removes the temp dir.
      await terraform.destroyWorkspace(ws);
    }
  });
});

// ── Members ─────────────────────────────────────────────────────────

test.describe("Terraform — Groups · Members @since-4.12", () => {
  test(`member lifecycle — add and remove members ${XRAY.TERRAFORM.GROUP_MEMBER_LIFECYCLE} ${TAGS.REGRESSION}`, async ({
    mapi,
  }) => {
    // init + 3 applies + destroy — size the timeout above the per-process
    // TF_TIMEOUT_MS ceiling for that many calls.
    test.setTimeout(terraform.TF_TIMEOUT_MS * 9);

    const SA_A = "e2e-sa-tf-group-a";
    const SA_B = "e2e-sa-tf-group-b";
    const member = (sa: string) => ({ source: "gravitee", source_id: sa, roles: { API: "USER" } });

    const ws = await terraform.initWorkspace("groups/members");
    try {
      let groupId: string;

      await test.step("Apply a group with one resolvable member", async () => {
        await mapi.createServiceAccount(SA_A);
        await mapi.createServiceAccount(SA_B);
        await terraform.writeVars(ws, { members: [member(SA_A)] });
        await terraform.apply(ws);
        groupId = await terraform.output(ws, "group_id");

        const members = await mapi.fetchGroupMembers(groupId);
        expect(members).toHaveLength(1);
        expect(members[0]?.displayName).toContain(SA_A);
        expect(members[0]?.roles?.["API"]).toBe("USER");
      });

      await test.step("Add a second member", async () => {
        await terraform.writeVars(ws, { members: [member(SA_A), member(SA_B)] });
        await terraform.apply(ws);

        const names = (await mapi.fetchGroupMembers(groupId)).map((m) => m.displayName ?? "");
        expect(names).toHaveLength(2);
        expect(names.some((n) => n.includes(SA_A))).toBe(true);
        expect(names.some((n) => n.includes(SA_B))).toBe(true);
      });

      await test.step("Remove a member", async () => {
        await terraform.writeVars(ws, { members: [member(SA_B)] });
        await terraform.apply(ws);

        const members = await mapi.fetchGroupMembers(groupId);
        expect(members).toHaveLength(1);
        expect(members[0]?.displayName).toContain(SA_B);
      });
    } finally {
      await terraform.destroyWorkspace(ws);
    }
  });

  // ACCEPTED LIMITATION (GKO-2862, Won't Do). A group member declared with a subset of
  // role scopes (e.g. `roles = { API = "USER" }`) is stored by APIM with the
  // other scopes' default roles expanded in (`APPLICATION` and `INTEGRATION`
  // are added server-side). The provider's Read surfaces that expanded map
  // into state, so `terraform plan` reports a perpetual diff and a
  // group-with-members is never idempotent — every apply re-expands, every
  // plan wants to strip. GKO-2862 was resolved "Won't Do", so this perpetual
  // diff is accepted behaviour, not a bug to be fixed. This test codifies that
  // contract: a group WITH members reports a diff right after apply, and the
  // diff is the server-expanded role scopes. If the provider ever round-trips
  // members cleanly, this fails loudly and we revisit GKO-2862.
  // Re-verified 2026-06-08 against provider main@c97d698.
  test(
    `terraform plan keeps reporting changes for server-expanded member roles ${XRAY.TERRAFORM.GROUP_MEMBERS_PERPETUAL_DIFF} ${TAGS.REGRESSION}`,
    async ({ mapi }) => {
      test.setTimeout(terraform.TF_TIMEOUT_MS * 8);

      const SA = "e2e-sa-tf-group-idem";
      const ws = await terraform.initWorkspace("groups/members");
      try {
        await mapi.createServiceAccount(SA);
        await terraform.writeVars(ws, {
          members: [{ source: "gravitee", source_id: SA, roles: { API: "USER" } }],
        });
        await terraform.apply(ws);

        // Per GKO-2862 (Won't Do), a plan right after apply is NOT a no-op:
        // APIM expanded the declared `API` role into the other scopes, and the
        // provider's Read surfaces them, so the plan proposes stripping the
        // server-added `APPLICATION`/`INTEGRATION` scopes.
        const result = await terraform.plan(ws);
        expect(result.hasChanges).toBe(true);
        expect(result.stdout).toMatch(/APPLICATION|INTEGRATION/);
      } finally {
        await terraform.destroyWorkspace(ws);
      }
    },
  );

  test(`non-resolvable member is accepted and silently ignored ${XRAY.TERRAFORM.GROUP_WITH_MEMBER} ${TAGS.REGRESSION}`, async ({
    mapi,
  }) => {
    test.setTimeout(terraform.TF_WORKSPACE_TIMEOUT_MS);

    const ws = await terraform.initWorkspace("groups/members");
    try {
      // The provider documents that members without an IDP entry are
      // ignored. The apply must still succeed and leave the group with zero
      // resolved members.
      await terraform.writeVars(ws, {
        members: [{ source: "memory", source_id: "e2e-ghost-user-xyz", roles: { API: "USER" } }],
      });
      await terraform.apply(ws);

      const groupId = await terraform.output(ws, "group_id");
      const members = await mapi.fetchGroupMembers(groupId);
      expect(members).toHaveLength(0);
    } finally {
      await terraform.destroyWorkspace(ws);
    }
  });
});

// ── Replacement, drift & import ──────────────────────────────────────

test.describe("Terraform — Groups · Replacement, drift & import @since-4.12", () => {
  test(`changing the hrid replaces the group ${XRAY.TERRAFORM.GROUP_HRID_REPLACEMENT} ${TAGS.REGRESSION}`, async ({
    mapi,
  }) => {
    test.setTimeout(terraform.TF_TIMEOUT_MS * 9);

    const ws = await terraform.initWorkspace("groups/lifecycle");
    try {
      await terraform.writeVars(ws, { hrid_suffix: "hrid-a", group_name: "e2e-tf-group-hrid-a" });
      await terraform.apply(ws);
      await mapi.waitForGroupMatches("e2e-tf-group-hrid-a", { name: "e2e-tf-group-hrid-a" });

      await terraform.writeVars(ws, { hrid_suffix: "hrid-b", group_name: "e2e-tf-group-hrid-b" });
      const planned = await terraform.plan(ws);
      expect(planned.hasChanges).toBe(true);
      // hrid is a RequiresReplace attribute — the plan must be a destroy +
      // recreate, not an in-place update.
      expect(planned.stdout).toMatch(/must be replaced|forces replacement/);

      await terraform.apply(ws);
      await mapi.waitForGroupAbsent("e2e-tf-group-hrid-a");
      await mapi.waitForGroupMatches("e2e-tf-group-hrid-b", { name: "e2e-tf-group-hrid-b" });
    } finally {
      await terraform.destroyWorkspace(ws);
    }
  });

  test(`terraform plan detects an out-of-band group deletion ${XRAY.TERRAFORM.GROUP_DRIFT} ${TAGS.REGRESSION}`, async ({
    mapi,
  }) => {
    test.setTimeout(terraform.TF_WORKSPACE_TIMEOUT_MS);

    const HRID = "e2e-tf-group-drift";
    const ws = await terraform.initWorkspace("groups/lifecycle");
    try {
      await terraform.writeVars(ws, { hrid_suffix: "drift", group_name: HRID });
      await terraform.apply(ws);
      const groupId = await terraform.output(ws, "group_id");

      // Delete the group behind Terraform's back.
      await mapi.deleteGroup(groupId);
      await mapi.waitForGroupAbsent(HRID);

      // The provider's Read must notice the group is gone, so the next plan
      // proposes recreating it instead of reporting no changes.
      const planned = await terraform.plan(ws);
      expect(planned.hasChanges).toBe(true);
    } finally {
      await terraform.destroyWorkspace(ws);
    }
  });

  test(`an existing group can be imported into Terraform state ${XRAY.TERRAFORM.GROUP_IMPORT} ${TAGS.REGRESSION}`, async () => {
    test.setTimeout(terraform.TF_TIMEOUT_MS * 9);

    const HRID = "e2e-tf-group-import";
    const ws = await terraform.initWorkspace("groups/lifecycle");
    try {
      await terraform.writeVars(ws, { hrid_suffix: "import", group_name: HRID });
      await terraform.apply(ws);

      // Forget the resource, then re-import it from APIM. The provider's
      // import id is a JSON object of environment_id / hrid / organization_id.
      await terraform.tf(ws, ["state", "rm", "apim_group.test"]);
      await terraform.tf(ws, [
        "import",
        "-no-color",
        "apim_group.test",
        JSON.stringify({ environment_id: "DEFAULT", hrid: HRID, organization_id: "DEFAULT" }),
      ]);

      // A clean import populates every attribute, so the plan is a no-op.
      const planned = await terraform.plan(ws);
      expect(planned.hasChanges).toBe(false);
    } finally {
      await terraform.destroyWorkspace(ws);
    }
  });
});

// ── Schema validation (negative) ─────────────────────────────────────

test.describe("Terraform — Groups · Validation @since-4.12", () => {
  test(`an invalid hrid is rejected by the provider ${XRAY.TERRAFORM.GROUP_INVALID_HRID} ${TAGS.REGRESSION}`, async () => {
    test.setTimeout(terraform.TF_WORKSPACE_TIMEOUT_MS);

    const ws = await terraform.initWorkspace("groups/lifecycle");
    try {
      // The space and "!" violate the hrid pattern
      // ^[a-zA-Z0-9][a-zA-Z0-9_-]+[a-zA-Z0-9]$.
      await terraform.writeVars(ws, { hrid_suffix: "bad hrid!" });
      const output = await terraform.applyExpectFailure(ws);
      expect(output.toLowerCase()).toMatch(/hrid|pattern|invalid/);
    } finally {
      await terraform.destroyWorkspace(ws);
    }
  });

  test(`an over-long name is rejected by the provider ${XRAY.TERRAFORM.GROUP_INVALID_NAME} ${TAGS.REGRESSION}`, async () => {
    test.setTimeout(terraform.TF_WORKSPACE_TIMEOUT_MS);

    const ws = await terraform.initWorkspace("groups/lifecycle");
    try {
      // The group name must be 1–512 characters.
      await terraform.writeVars(ws, {
        hrid_suffix: "long-name",
        group_name: "x".repeat(513),
      });
      const output = await terraform.applyExpectFailure(ws);
      expect(output.toLowerCase()).toMatch(/name|length|512|invalid/);
    } finally {
      await terraform.destroyWorkspace(ws);
    }
  });
});

// ── Data source ─────────────────────────────────────────────────────

test.describe("Terraform — Groups · Data source @since-4.12", () => {
  test(`data "apim_group" reads a group back by hrid ${XRAY.TERRAFORM.GROUP_DATA_SOURCE} ${TAGS.REGRESSION}`, async () => {
    test.setTimeout(terraform.TF_WORKSPACE_TIMEOUT_MS);

    const ws = await terraform.initWorkspace("groups/datasource");
    try {
      await terraform.apply(ws);

      // The data source must surface exactly what the resource created.
      const [resourceId, resourceName, dsId, dsName, dsHrid, dsNotify] = await Promise.all([
        terraform.output(ws, "resource_id"),
        terraform.output(ws, "resource_name"),
        terraform.output(ws, "ds_id"),
        terraform.output(ws, "ds_name"),
        terraform.output(ws, "ds_hrid"),
        terraform.output(ws, "ds_notify_members"),
      ]);

      expect(dsId).toBe(resourceId);
      expect(dsName).toBe(resourceName);
      expect(dsHrid).toBe("e2e-tf-group-datasource");
      expect(dsNotify).toBe("false");
    } finally {
      await terraform.destroyWorkspace(ws);
    }
  });
});
