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
 * Journey: control who can discover an API in the portal, then retire it.
 *
 * As an API producer, I decide who can see my API in the developer portal
 * (`visibility`) and whether it is listed there at all (`lifecycleState`), then
 * delete it when it is retired. Both fields are portal-side concerns: the
 * regression that matters is that they change IN PLACE on the existing API and
 * that neither of them takes the API off the gateway — only `state` does, which
 * is `publish-api-and-serve-traffic`'s story.
 *
 * The portal's own rendering is a downstream consumer of these two fields; this
 * journey asserts them at the source (the APIM record) so no browser fixture is
 * needed.
 *
 * Fixtures are co-located in this folder. Admission-time rules around the same
 * fields (no plans + STARTED, context-path conflict, missing required fields)
 * are Kubernetes-layer behaviour and stay in tests/gko/admission-webhook.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test } from "../../../../setup.js";
import { XRAY, TAGS } from "../../../../helpers/tags.js";
import { forEachProvisioner } from "../../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../../helpers/provisioner-env.js";
import { assertProvisioner } from "../../../../../src/provisioners/index.js";
import type { ApiLifecycleState, ApiVisibility } from "../../../../../src/types/apim.js";

const here = path.dirname(fileURLToPath(import.meta.url));

/** The portal-facing state of the API, re-provisioned on every update(). */
interface PortalStateParams {
  visibility: ApiVisibility;
  lifecycleState: ApiLifecycleState;
}

/** The three combinations the portal actually distinguishes. */
const MATRIX: PortalStateParams[] = [
  { visibility: "PRIVATE", lifecycleState: "PUBLISHED" },
  { visibility: "PUBLIC", lifecycleState: "PUBLISHED" },
  { visibility: "PUBLIC", lifecycleState: "UNPUBLISHED" },
];

/** One GKO manifest per combination, all under the same CR name. */
function gkoManifest(params: PortalStateParams): string {
  const slug = `${params.visibility}-${params.lifecycleState}`.toLowerCase();
  return path.join(here, `gko/api-${slug}.yaml`);
}

forEachProvisioner<PortalStateParams>(
  {
    title: "Configure portal visibility and lifecycle state, then retire the API",
    provisioners: {
      gko: gkoScenario<PortalStateParams>({
        // The API is the parameterized resource, so there is no static
        // manifest: provision applies the first matrix entry and every
        // update() re-applies another variant over the same CR name.
        manifests: [],
        roles: { api: "portal-visibility-api" },
        dynamicRoles: ["api"],
        contextPath: "/portal-visibility-api",
        applyParams: async (k, params) => {
          await k.apply(gkoManifest(params));
        },
      }),
      terraform: tfScenario<PortalStateParams>({
        fixture: path.join(here, "terraform"),
        toVars: (params) => ({
          visibility: params.visibility,
          lifecycle_state: params.lifecycleState,
        }),
        removeVars: { api: { create_api: false } },
        // `[0]` is the count index of the gated resource; view.read() matches
        // the state address exactly.
        addresses: { api: "apim_apiv4.api[0]" },
      }),
    },
    xray: {
      gko: [
        XRAY.API_LIFECYCLE.V4_VISIBILITY_PRIVATE,
        XRAY.API_LIFECYCLE.V4_VISIBILITY_PUBLIC,
        XRAY.API_LIFECYCLE.V4_PUBLISHED_IN_PORTAL,
        XRAY.API_LIFECYCLE.V4_UNPUBLISHED_NOT_IN_PORTAL,
        XRAY.API_LIFECYCLE.V4_PORTAL_VISIBILITY_RULES,
        XRAY.API_LIFECYCLE.DELETE_V4_API,
      ],
      terraform: XRAY.TERRAFORM.API_VISIBILITY_LIFECYCLE_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 120_000 },
  },
  async ({ provisioned, mapi, gateway }) => {
    const apiId = await provisioned.apiId();
    const ctx = await provisioned.contextPath();

    for (const [index, params] of MATRIX.entries()) {
      const label = `${params.visibility} + ${params.lifecycleState}`;
      await test.step(`APIM records ${label}`, async () => {
        // Entry 0 is what provision() already applied.
        if (index > 0) await provisioned.update(params);
        await mapi.waitForApiMatches(
          apiId,
          { visibility: params.visibility, lifecycleState: params.lifecycleState, state: "STARTED" },
          { timeoutMs: 30_000 },
        );
      });
    }

    await test.step("Unlisting from the portal leaves the API on the gateway", async () => {
      // The API is UNPUBLISHED at this point (last matrix entry) but still
      // STARTED, so traffic must be unaffected.
      await gateway.assertResponds(ctx, { status: 200 });
    });

    await test.step("Retiring the API removes it from APIM and the gateway", async () => {
      await provisioned.remove("api");
      await assertProvisioner(provisioned, "api", "gone");
      await mapi.waitForApiAbsent(apiId, { timeoutMs: 30_000 });
      await gateway.assertResponds(ctx, { status: 404 });
    });
  },
  MATRIX[0],
);
