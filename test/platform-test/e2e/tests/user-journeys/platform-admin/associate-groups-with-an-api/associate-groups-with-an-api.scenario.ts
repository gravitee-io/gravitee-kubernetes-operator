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
 * Journey: put an API under a group's ownership, then detach it.
 *
 * As a platform admin, I hand an API to a team by associating it with that
 * team's group, and detach it when ownership moves. `groups` is an inline
 * attribute of the API on both drivers (a list of group names, HRIDs or UUIDs)
 * with no standalone association resource.
 *
 * The group is created by the SAME provisioner as the API, and must reconcile
 * first: APIM silently drops a reference to a group that does not exist yet, so
 * the API is applied only after the group is ready (GKO orders it through
 * `dynamicRoles`, Terraform through the `apim_group.group.name` reference).
 * That silent drop is also why the assertion resolves the group's APIM id and
 * looks for it on the API record, rather than trusting the apply to have worked.
 *
 * Fixtures are co-located in this folder. A reference to a group that does NOT
 * exist is operator behaviour (a `.status` warning naming the missing group) and
 * stays in tests/gko/members.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test, expect } from "../../../../setup.js";
import { XRAY, TAGS } from "../../../../helpers/tags.js";
import { forEachProvisioner } from "../../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../../helpers/provisioner-env.js";

const here = path.dirname(fileURLToPath(import.meta.url));

/** The single knob: whether the API is associated with the group. */
interface GroupAssociationParams {
  withGroup: boolean;
}

forEachProvisioner<GroupAssociationParams>(
  {
    title: "Associate a group with an API",
    provisioners: {
      gko: gkoScenario<GroupAssociationParams>({
        // Only the Group is static. The API is dynamic so the provisioner
        // applies it AFTER the Group reaches Accepted, and re-applies it on
        // every update().
        manifests: [path.join(here, "gko/group.yaml")],
        roles: { group: "api-owning-group", api: "api-in-a-group" },
        dynamicRoles: ["api"],
        applyParams: async (k, params) => {
          await k.apply(
            path.join(here, params.withGroup ? "gko/api-with-group.yaml" : "gko/api-without-group.yaml"),
          );
        },
      }),
      terraform: tfScenario<GroupAssociationParams>({
        fixture: path.join(here, "terraform"),
        toVars: (params) => ({ with_group: params.withGroup }),
      }),
    },
    xray: {
      gko: [
        XRAY.MEMBERS.V4_CREATE_EXISTING_GROUP,
        XRAY.MEMBERS.V4_ADD_GROUP_REFS,
        XRAY.MEMBERS.V4_ADD_GROUP_REFS_VARIANT,
      ],
      terraform: XRAY.TERRAFORM.API_GROUPS_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 120_000 },
  },
  async ({ provisioned, mapi }) => {
    const apiId = await provisioned.apiId();
    const groupId = await provisioned.groupId();

    await test.step("The API is recorded under the group", async () => {
      await expect
        .poll(async () => (await mapi.fetchApi(apiId)).groups ?? [], {
          timeout: 30_000,
          message: "API is associated with the group in APIM",
        })
        .toEqual([groupId]);
    });

    await test.step("Detaching the API clears the association", async () => {
      await provisioned.update({ withGroup: false });
      await expect
        .poll(async () => (await mapi.fetchApi(apiId)).groups ?? [], {
          timeout: 30_000,
          message: "API is no longer associated with the group",
        })
        .toEqual([]);
    });
  },
  { withGroup: true },
);
