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
 * Journey: create, rename, and delete a group.
 *
 * As a platform admin, I create a group so I can organise API members, rename it
 * as the team it represents changes, and delete it when the team is dissolved.
 * A group created through any provisioner lands in APIM via the Automation API
 * (origin KUBERNETES), and a rename must update it IN PLACE rather than
 * replacing it — the id is what API memberships point at.
 *
 * Fixtures are co-located in this folder. Provisioner-specific group behaviour
 * (GKO member reconciliation / admission; Terraform drift, import, data source,
 * hrid replacement, provider-side validation) stays in the per-provisioner
 * suites under tests/gko/groups and tests/terraform/groups.test.ts.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test, expect } from "../../../../setup.js";
import { XRAY, TAGS } from "../../../../helpers/tags.js";
import { forEachProvisioner } from "../../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../../helpers/provisioner-env.js";
import { assertProvisioner } from "../../../../../src/provisioners/index.js";

const here = path.dirname(fileURLToPath(import.meta.url));

/** The single knob the journey re-provisions with: the created vs renamed state. */
interface GroupParams {
  renamed: boolean;
}

/** Matches the name in gko/group-renamed.yaml; the assertion only checks the suffix. */
const RENAMED_TF = "simple-group-tf-renamed";

forEachProvisioner<GroupParams>(
  {
    title: "Create, rename, and delete a group",
    provisioners: {
      gko: gkoScenario<GroupParams>({
        manifests: [path.join(here, "gko/group.yaml")],
        roles: { group: "simple-group" },
        // provision applies the created manifest; update() re-applies the renamed
        // one over it. At provision params.renamed is false, so this is a no-op.
        applyParams: async (k, params) => {
          if (params.renamed) await k.apply(path.join(here, "gko/group-renamed.yaml"));
        },
      }),
      terraform: tfScenario<GroupParams>({
        fixture: path.join(here, "terraform"),
        toVars: (params) => ({ group_name: params.renamed ? RENAMED_TF : "simple-group-tf" }),
        removeVars: { group: { create_group: false } },
        // The `[0]` is the count index: a `count`-gated resource is addressed by
        // index in Terraform state, and `view.read()` matches the address exactly.
        addresses: { group: "apim_group.group[0]" },
      }),
    },
    xray: {
      gko: [XRAY.GROUPS.CREATE_WITH_MEMBER, XRAY.GROUPS.MODIFY_GROUP, XRAY.GROUPS.DELETE_GROUP],
      terraform: [
        XRAY.TERRAFORM.GROUP_CREATE,
        XRAY.TERRAFORM.GROUP_UPDATE,
        XRAY.TERRAFORM.GROUP_DESTROY,
        XRAY.TERRAFORM.GROUP_NOTIFY_MEMBERS,
      ],
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 60_000 },
  },
  async ({ provisioned, mapi }) => {
    const groupId = await provisioned.groupId();

    await test.step("Created group lands in APIM (origin KUBERNETES)", async () => {
      // Both provisioners write through the Automation API, so APIM records
      // origin KUBERNETES either way. notifyMembers is the inverse of APIM's
      // wire field, and both fixtures declare it false.
      await mapi.waitForGroupMatchesById(groupId, {
        origin: "KUBERNETES",
        disable_membership_notifications: true,
      });
    });

    await test.step("Rename updates the group in place, keeping its id", async () => {
      await provisioned.update({ renamed: true });
      // Poll on the id, not the hrid: a rename that REPLACED the group would
      // surface here as the id never showing the new name.
      await expect
        .poll(async () => (await mapi.fetchGroupById(groupId)).name, {
          timeout: 30_000,
          message: "renamed group is readable under the same id",
        })
        .toMatch(/renamed$/);
    });

    await test.step("Deleting the group removes it from APIM", async () => {
      await provisioned.remove("group");
      await assertProvisioner(provisioned, "group", "gone");
      await mapi.waitForGroupAbsentById(groupId, { timeoutMs: 30_000 });
    });
  },
  { renamed: false },
);
