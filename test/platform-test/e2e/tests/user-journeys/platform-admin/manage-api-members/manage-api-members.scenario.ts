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
 * Journey: grant, re-role and revoke access to an API.
 *
 * As a platform admin, I grant a teammate access to an API, adjust their role as
 * their responsibilities change, and revoke it when they move on. Members are an
 * inline attribute on both drivers (`spec.members` / `apim_apiv4.members`) with
 * no standalone membership resource, so the journey asserts the API's member list
 * in APIM after each change.
 *
 * The member has to be a REAL user: a member naming a user that does not exist is
 * silently dropped (the operator logs a warning and APIM records only the primary
 * owner), so the journey creates a `gravitee`-source service account as a
 * provisioner-agnostic precondition, the same way create-group-with-member does.
 * The primary owner is a member too and is expected in every readout.
 *
 * Fixtures are co-located in this folder. Member cases that only the operator can
 * produce stay in tests/gko/members: a member naming a non-existent user or
 * group, PRIMARY_OWNER declared in `members`, primary owner resolved from the
 * ManagementContext user, a member missing `source`, and a member declared
 * WITHOUT a role — the Automation API's member schema marks `role` non-nullable,
 * so the Terraform provider rejects omitting it and only GKO can default it.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test, expect } from "../../../../setup.js";
import { XRAY, TAGS } from "../../../../helpers/tags.js";
import { forEachProvisioner } from "../../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../../helpers/provisioner-env.js";

const here = path.dirname(fileURLToPath(import.meta.url));

/**
 * The service account both fixtures name as the member. `createServiceAccount`
 * builds its display name from firstname + lastname, and the v2 members endpoint
 * identifies a member by display name only (it returns no source/sourceId).
 */
const MEMBER_SOURCE_ID = "e2e-sa-api-member";
const MEMBER_DISPLAY_NAME = `${MEMBER_SOURCE_ID} Service`;

/** APIM's primary owner in the test environment, a member of every API. */
const PRIMARY_OWNER = "admin";

/** How the API declares the member, if at all. */
interface MemberParams {
  role: "USER" | "REVIEWER";
  withMember: boolean;
  notifyMembers: boolean;
}

interface MemberStage {
  label: string;
  params: MemberParams;
  /** displayName -> role name, for every member APIM should report. */
  expected: Record<string, string>;
}

/** Each stage must move the role to a DIFFERENT value than the one before it, or
 * a no-op would pass as a successful change. */
const STAGES: MemberStage[] = [
  {
    label: "granted with the USER role",
    params: { role: "USER", withMember: true, notifyMembers: false },
    expected: { [PRIMARY_OWNER]: "PRIMARY_OWNER", [MEMBER_DISPLAY_NAME]: "USER" },
  },
  {
    label: "re-roled to REVIEWER, with member notifications on",
    params: { role: "REVIEWER", withMember: true, notifyMembers: true },
    expected: { [PRIMARY_OWNER]: "PRIMARY_OWNER", [MEMBER_DISPLAY_NAME]: "REVIEWER" },
  },
  {
    label: "revoked, leaving only the primary owner",
    params: { role: "USER", withMember: false, notifyMembers: false },
    expected: { [PRIMARY_OWNER]: "PRIMARY_OWNER" },
  },
];

/** One GKO manifest per stage, all under the same CR name. */
function gkoManifest(params: MemberParams): string {
  if (!params.withMember) return path.join(here, "gko/api-no-members.yaml");
  return path.join(here, `gko/api-role-${params.role.toLowerCase()}.yaml`);
}

// Created once per test in beforeEach (mapi is a test-scoped fixture, so it is
// unreachable from beforeAll). Idempotent, so repeat runs are safe. The account
// is left behind on purpose: it is a shared, inert precondition, and deleting it
// would break a concurrent read of an API that still references it.
test.beforeEach(async ({ mapi }) => {
  await mapi.createServiceAccount(MEMBER_SOURCE_ID);
});

forEachProvisioner<MemberParams>(
  {
    title: "Grant, re-role and revoke access to an API",
    provisioners: {
      gko: gkoScenario<MemberParams>({
        // The API carries the membership, so it is the parameterized resource.
        manifests: [],
        roles: { api: "api-with-members" },
        dynamicRoles: ["api"],
        applyParams: async (k, params) => {
          await k.apply(gkoManifest(params));
        },
      }),
      terraform: tfScenario<MemberParams>({
        fixture: path.join(here, "terraform"),
        toVars: (params) => ({
          with_member: params.withMember,
          member_role: params.role,
          notify_members: params.notifyMembers,
        }),
      }),
    },
    xray: {
      gko: [
        XRAY.MEMBERS.V4_ADD_MEMBER_WITH_ROLE_NAME,
        XRAY.MEMBERS.V4_DUPLICATE_KEY_ON_ROLE_CHANGE,
        XRAY.MEMBERS.V4_REMOVE_MEMBER,
        XRAY.MEMBERS.V4_REMOVE_MEMBER_VARIANT,
        XRAY.MEMBERS.V4_NOTIFY_MEMBERS_ENABLED,
      ],
      terraform: XRAY.TERRAFORM.API_MEMBERS_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 120_000 },
  },
  async ({ provisioned, mapi }) => {
    const apiId = await provisioned.apiId();

    for (const [index, stage] of STAGES.entries()) {
      await test.step(`Membership ${stage.label}`, async () => {
        // Stage 0 is what provision() already applied.
        if (index > 0) await provisioned.update(stage.params);

        // Compare the WHOLE member set, not just the member under test: a role
        // change that added a second entry instead of updating the existing one
        // (GKO-259's duplicate-key failure mode) only shows up as an extra key.
        await expect
          .poll(
            async () => {
              const members = await mapi.listApiMembers(apiId);
              return Object.fromEntries(
                members.map((m) => [m.displayName ?? m.id, m.roles?.[0]?.name ?? "(none)"]),
              );
            },
            { timeout: 30_000, message: `API members after: ${stage.label}` },
          )
          .toEqual(stage.expected);
      });
    }
  },
  STAGES[0].params,
);
