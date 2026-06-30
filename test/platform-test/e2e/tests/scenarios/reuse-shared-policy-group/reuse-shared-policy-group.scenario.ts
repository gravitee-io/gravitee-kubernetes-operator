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
 * Use case: reuse a shared policy group across an API.
 *
 * Currently BLOCKED on both provisioners, so both arms are `pending` (a visible
 * gap, never a green-washed test). The intended invariant is provisioner-agnostic:
 * an API whose request flow invokes a Shared Policy Group runs that SPG at the
 * gateway, and detaching it stops it.
 *
 * Blockers (see fixtures/use-cases/reuse-shared-policy-group/README.md):
 *   - GKO: the documented `sharedPolicyGroupRef` form (gko/api-with-spg.yaml) is
 *     rejected by the admission webhook — the SPG ref is resolved to the SPG
 *     crossId by the reconciler but NOT before the APIM dry-run (GKO-3001).
 *   - Terraform: `apim_shared_policy_group` exposes only `id`, not the crossId,
 *     and only the crossId actually executes the SPG at the gateway.
 *
 * When GKO-3001 (and the TF crossId gap) are fixed, restore the `gko`/`terraform`
 * factories and run the body below. The fixtures already use the correct forms.
 */

import { expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import { forEachProvisioner } from "../../../helpers/for-each-provisioner.js";

/** Does the API have a request flow that invokes a shared-policy-group step? */
interface ApiWithFlows {
  flows?: Array<{ request?: Array<{ policy?: string }> }>;
}
function reusesSpg(api: ApiWithFlows): boolean {
  return (api.flows ?? []).some((f) =>
    (f.request ?? []).some((s) => s.policy === "shared-policy-group-policy"),
  );
}

/** The single knob: whether the API reuses the shared policy group. */
interface SpgParams {
  withSpg: boolean;
}

forEachProvisioner<SpgParams>(
  {
    title: "Reuse a shared policy group across an API",
    // No provisioner can currently stand up a working SPG reuse end-to-end.
    provisioners: {},
    pending: {
      gko: "GKO-3001: admission rejects sharedPolicyGroupRef (SPG ref not resolved before the APIM dry-run)",
      terraform: "apim_shared_policy_group exposes no crossId, and only the crossId executes the SPG at the gateway",
    },
    xray: {
      gko: [XRAY.SHARED_POLICY_GROUPS.ADD_SPG_TO_API, XRAY.SHARED_POLICY_GROUPS.REMOVE_SPG_FROM_API],
      terraform: XRAY.TERRAFORM.SPG_REUSE_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
  },
  // Intended assertion, run once the blockers above are lifted: APIM records the
  // API flow invoking the SPG, and detaching removes it.
  async ({ provisioned, mapi }) => {
    const apiId = await provisioned.apiId();
    await expect
      .poll(async () => reusesSpg((await mapi.fetchApi(apiId)) as ApiWithFlows), { timeout: 30_000 })
      .toBe(true);
    await provisioned.update({ withSpg: false });
    await expect
      .poll(async () => reusesSpg((await mapi.fetchApi(apiId)) as ApiWithFlows), { timeout: 30_000 })
      .toBe(false);
  },
  { withSpg: true },
);
