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
 * Journey: share an application with a teammate.
 *
 * As an application developer, I share my application with a teammate, promote
 * them when they take over its subscriptions, and remove them when they leave.
 * Members are an inline attribute on both drivers (`spec.members` /
 * `apim_application.members`) with no standalone membership resource, so the
 * journey asserts the application's member list in APIM after each change.
 *
 * As with API members, the member must be a REAL user or it is silently dropped,
 * so the journey creates a `gravitee`-source service account as a precondition.
 * The primary owner is a member too and is expected in every readout.
 *
 * Fixtures are co-located in this folder. The member cases only the operator can
 * produce — a member naming a non-existent user, group or role, a member missing
 * `source`, and a member declared WITHOUT a role (the Automation API marks it
 * non-nullable, so Terraform rejects omitting it) — stay in
 * tests/gko/applications.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test, expect } from "../../../../setup.js";
import { XRAY, TAGS } from "../../../../helpers/tags.js";
import { forEachProvisioner } from "../../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../../helpers/provisioner-env.js";

const here = path.dirname(fileURLToPath(import.meta.url));

/** Named by both fixtures; the members endpoint identifies a member by display name. */
const MEMBER_SOURCE_ID = "e2e-sa-app-member";
const MEMBER_DISPLAY_NAME = `${MEMBER_SOURCE_ID} Service`;

/** APIM's primary owner in the test environment, a member of every application. */
const PRIMARY_OWNER = "admin";

interface AppMemberParams {
  role: "USER" | "OWNER";
  withMember: boolean;
  notifyMembers: boolean;
}

interface AppMemberStage {
  label: string;
  params: AppMemberParams;
  /** displayName -> role, for every member APIM should report. */
  expected: Record<string, string>;
}

const STAGES: AppMemberStage[] = [
  {
    label: "shared with the USER role",
    params: { role: "USER", withMember: true, notifyMembers: false },
    expected: { [PRIMARY_OWNER]: "PRIMARY_OWNER", [MEMBER_DISPLAY_NAME]: "USER" },
  },
  {
    label: "promoted to OWNER, with member notifications on",
    params: { role: "OWNER", withMember: true, notifyMembers: true },
    expected: { [PRIMARY_OWNER]: "PRIMARY_OWNER", [MEMBER_DISPLAY_NAME]: "OWNER" },
  },
  {
    label: "revoked, leaving only the primary owner",
    params: { role: "USER", withMember: false, notifyMembers: false },
    expected: { [PRIMARY_OWNER]: "PRIMARY_OWNER" },
  },
];

/** One GKO manifest per stage, all under the same CR name. */
function gkoManifest(params: AppMemberParams): string {
  if (!params.withMember) return path.join(here, "gko/application-no-members.yaml");
  return path.join(here, `gko/application-role-${params.role.toLowerCase()}.yaml`);
}

// Created once per test in beforeEach (mapi is a test-scoped fixture, so it is
// unreachable from beforeAll). Idempotent, so repeat runs are safe.
test.beforeEach(async ({ mapi }) => {
  await mapi.createServiceAccount(MEMBER_SOURCE_ID);
});

forEachProvisioner<AppMemberParams>(
  {
    title: "Share an application with a teammate",
    provisioners: {
      gko: gkoScenario<AppMemberParams>({
        // The application carries the membership, so it is the parameterized
        // resource rather than a static manifest.
        manifests: [],
        roles: { application: "shared-application" },
        dynamicRoles: ["application"],
        applyParams: async (k, params) => {
          await k.apply(gkoManifest(params));
        },
      }),
      terraform: tfScenario<AppMemberParams>({
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
        XRAY.APPLICATIONS_MEMBERS.APP_ADD_MEMBER_ROLE_NAME,
        XRAY.APPLICATIONS_MEMBERS.APP_CHANGE_MEMBER_ROLE,
        XRAY.APPLICATIONS_MEMBERS.APP_REMOVE_MEMBER,
      ],
      terraform: XRAY.TERRAFORM.APPLICATION_MEMBERS_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 90_000 },
  },
  async ({ provisioned, mapi }) => {
    const appId = await provisioned.applicationId();

    for (const [index, stage] of STAGES.entries()) {
      await test.step(`Application ${stage.label}`, async () => {
        // Stage 0 is what provision() already applied.
        if (index > 0) await provisioned.update(stage.params);

        // Compare the WHOLE member set: a role change that added a second entry
        // instead of updating the existing one only shows up as an extra key.
        await expect
          .poll(
            async () => {
              const members = await mapi.listApplicationMembers(appId);
              return Object.fromEntries(
                members.map((m) => [m.displayName ?? m.id, m.role ?? "(none)"]),
              );
            },
            { timeout: 30_000, message: `application members after: ${stage.label}` },
          )
          .toEqual(stage.expected);
      });
    }
  },
  STAGES[0].params,
);
