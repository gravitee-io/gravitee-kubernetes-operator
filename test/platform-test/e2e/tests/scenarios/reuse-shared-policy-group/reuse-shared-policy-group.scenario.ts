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
 * A Shared Policy Group is defined once and reused from a V4 API's request flow.
 * The shared invariant is provisioner-agnostic: whichever driver authors it,
 * APIM records the API with a flow that invokes the SPG via a
 * shared-policy-group-policy step; detaching the SPG removes that flow.
 *
 * NOTE: this asserts the SPG reuse at the APIM config level (matching the
 * original GKO-976/980 intent). End-to-end gateway execution of an SPG step is
 * NOT asserted here: a stronger gateway check (the SPG's transform-headers value
 * reflected by the echo backend) did not resolve for EITHER provisioner, which
 * points at an SPG deployment-lifecycle gap rather than a provisioner difference.
 * Tracked as a follow-up; see PARITY.md.
 *
 * Fixtures live in fixtures/use-cases/reuse-shared-policy-group/. SPG update
 * (id stability) and apiType validation (GKO-981/1462) stay GKO-only under
 * tests/gko/shared-policy-groups.
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import { forEachProvisioner } from "../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../helpers/provisioner-env.js";

/** Does the API have a request flow that invokes a shared-policy-group step? */
interface ApiWithFlows {
  flows?: Array<{ request?: Array<{ policy?: string }> }>;
}
function reusesSpg(api: ApiWithFlows): boolean {
  return (api.flows ?? []).some((f) =>
    (f.request ?? []).some((s) => s.policy === "shared-policy-group-policy"),
  );
}

const FIXTURE = "use-cases/reuse-shared-policy-group";
const GKO_API_WITH_SPG = fixture(`${FIXTURE}/gko/api-with-spg.yaml`);
const GKO_API_WITHOUT_SPG = fixture(`${FIXTURE}/gko/api-without-spg.yaml`);

/** The single knob: whether the API reuses the shared policy group. */
interface SpgParams {
  withSpg: boolean;
}

forEachProvisioner<SpgParams>(
  {
    title: "Reuse a shared policy group across an API",
    provisioners: {
      gko: gkoScenario<SpgParams>({
        // The SPG is static; the API (which references it) is the parameterized
        // resource so update() can detach the SPG by re-applying without the flow.
        manifests: [`${FIXTURE}/gko/shared-policy-group.yaml`],
        roles: {
          spg: { kind: "sharedpolicygroup", name: "e2e-uc-spg" },
          api: "e2e-uc-spg-api",
        },
        dynamicRoles: ["api"],
        contextPath: "/e2e-uc-spg-api",
        applyParams: async (k, params) => {
          await k.apply(params.withSpg ? GKO_API_WITH_SPG : GKO_API_WITHOUT_SPG);
        },
      }),
      terraform: tfScenario<SpgParams>({
        fixture: `${FIXTURE}/terraform`,
        toVars: (params) => ({ attach_spg: params.withSpg }),
      }),
    },
    xray: {
      gko: [XRAY.SHARED_POLICY_GROUPS.ADD_SPG_TO_API, XRAY.SHARED_POLICY_GROUPS.REMOVE_SPG_FROM_API],
      terraform: XRAY.TERRAFORM.SPG_REUSE_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 90_000 },
  },
  async ({ provisioned, mapi }) => {
    const apiId = await provisioned.apiId();

    await test.step("API reuses the SPG (flow invokes it in APIM)", async () => {
      await expect
        .poll(async () => reusesSpg((await mapi.fetchApi(apiId)) as ApiWithFlows), {
          timeout: 30_000,
          message: "API flow references the shared policy group",
        })
        .toBe(true);
    });

    await test.step("Detaching the SPG removes the flow reference", async () => {
      await provisioned.update({ withSpg: false });
      await expect
        .poll(async () => reusesSpg((await mapi.fetchApi(apiId)) as ApiWithFlows), {
          timeout: 30_000,
          message: "SPG flow reference removed from the API",
        })
        .toBe(false);
    });
  },
  { withSpg: true },
);
