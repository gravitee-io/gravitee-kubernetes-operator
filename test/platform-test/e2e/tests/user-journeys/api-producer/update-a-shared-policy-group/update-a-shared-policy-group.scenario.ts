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
 * Journey: update a shared policy group.
 *
 * As an API producer, I define a shared policy group once and change its policy
 * step as my needs change, without the change landing as a second resource or as
 * an undeployed draft.
 *
 * Both stages assert the WHOLE declared group — description, phase, api type and
 * the step's own configuration — rather than only that the resource survived.
 * The update changes the step name and the header value it injects, so an update
 * that reaches APIM but silently no-ops fails here.
 *
 * Id stability is asserted by construction rather than by comparing two ids: the
 * second stage polls the ORIGINAL id for the NEW content. A group replaced
 * instead of updated makes that id 404, and a no-op update leaves the old header
 * value behind. Re-reading the id from the handle would prove neither — it is
 * cached per role, so the second call is not a fresh read.
 *
 * `lifecycleState` is asserted per stage because an update that APIM accepts but
 * never deploys is invisible to every consumer of the group.
 *
 * Fixtures are co-located in this folder. Reusing a shared policy group inside an
 * API flow, and therefore anything at the gateway, belongs to the neighbouring
 * `reuse-shared-policy-group` journey, which stays blocked.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test } from "../../../../setup.js";
import type { DeepPartial, SharedPolicyGroup } from "../../../../../src/index.js";
import { XRAY, TAGS } from "../../../../helpers/tags.js";
import { forEachProvisioner } from "../../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../../helpers/provisioner-env.js";

const here = path.dirname(fileURLToPath(import.meta.url));

/** The one knob: whether the group carries its updated step and description. */
interface SpgParams {
  updated: boolean;
}

/**
 * Declared identically in both fixtures. APIM adds `scope` to a transform-headers
 * configuration on its own, which partial matching tolerates.
 */
function expectedSpg(params: SpgParams): DeepPartial<SharedPolicyGroup> {
  return {
    apiType: "PROXY",
    phase: "REQUEST",
    lifecycleState: "DEPLOYED",
    description: params.updated
      ? "Shared policy group after its update"
      : "Shared policy group as first authored",
    steps: [
      {
        name: params.updated ? "Inject the updated tracking header" : "Inject the tracking header",
        enabled: true,
        policy: "transform-headers",
        configuration: {
          addHeaders: [
            { name: "X-SPG-Test", value: params.updated ? "spg-header-updated" : "spg-header" },
          ],
        },
      },
    ],
  };
}

forEachProvisioner<SpgParams>(
  {
    title: "Update a shared policy group's step and description",
    provisioners: {
      gko: gkoScenario<SpgParams>({
        // The group itself carries the parameterized step, so it is applied by
        // applyParams rather than as a static manifest.
        manifests: [],
        roles: { sharedPolicyGroup: "update-spg" },
        dynamicRoles: ["sharedPolicyGroup"],
        applyParams: async (k, params) => {
          const variant = params.updated ? "updated" : "initial";
          await k.apply(path.join(here, `gko/spg-${variant}.yaml`));
        },
      }),
      terraform: tfScenario<SpgParams>({
        fixture: path.join(here, "terraform"),
        toVars: (params) => ({ updated: params.updated }),
      }),
    },
    xray: {
      gko: XRAY.SHARED_POLICY_GROUPS.UPDATE_SPG,
      terraform: XRAY.TERRAFORM.SPG_UPDATE_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 90_000 },
  },
  async ({ provisioned, mapi }) => {
    const spgId = await provisioned.sharedPolicyGroupId();

    await test.step("APIM records the shared policy group as declared", async () => {
      await mapi.waitForSharedPolicyGroupMatches(spgId, expectedSpg({ updated: false }), {
        timeoutMs: 30_000,
      });
    });

    await test.step("Updating the step rewrites the same shared policy group", async () => {
      await provisioned.update({ updated: true });
      await mapi.waitForSharedPolicyGroupMatches(spgId, expectedSpg({ updated: true }), {
        timeoutMs: 30_000,
      });
    });
  },
  { updated: false },
);
