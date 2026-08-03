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
 * Journey: apply a policy to a flow.
 *
 * As an API producer, I add a transform-headers policy to my API's response
 * phase, change which header it adds, and remove it again.
 *
 * The assertion is at the GATEWAY, not only in the API definition: APIM will
 * happily record a flow the gateway never applies, so a control-plane-only check
 * proves nothing about the policy actually running. The definition is checked
 * too, so a failure says whether the flow never reached APIM or reached it and
 * was not executed.
 *
 * Fixtures are co-located in this folder. Plan-level policy validation (general
 * conditions referencing a page that does not exist) stays GKO-only under
 * tests/gko/policies.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test, expect } from "../../../../setup.js";
import type { ApiV4 } from "../../../../../src/types/apim.js";
import { XRAY, TAGS } from "../../../../helpers/tags.js";
import { forEachProvisioner } from "../../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../../helpers/provisioner-env.js";

const here = path.dirname(fileURLToPath(import.meta.url));

interface PolicyParams {
  withPolicy: boolean;
  flowName: string;
  headerName: string;
  headerValue: string;
}

/** Matches gko/api-with-policy.yaml. */
const INITIAL: PolicyParams = {
  withPolicy: true,
  flowName: "Add custom header",
  headerName: "X-E2E-Policy",
  headerValue: "applied",
};

/** Matches gko/api-with-updated-policy.yaml. */
const UPDATED: PolicyParams = {
  withPolicy: true,
  flowName: "Add updated custom header",
  headerName: "X-E2E-Policy-Updated",
  headerValue: "reapplied",
};

const REMOVED: PolicyParams = { ...INITIAL, withPolicy: false };

/** One GKO manifest per stage, all under the same CR name. */
function gkoManifest(params: PolicyParams): string {
  if (!params.withPolicy) return path.join(here, "gko/api-without-policy.yaml");
  const variant = params.headerName === UPDATED.headerName ? "api-with-updated-policy" : "api-with-policy";
  return path.join(here, `gko/${variant}.yaml`);
}

forEachProvisioner<PolicyParams>(
  {
    title: "Apply a policy to an API flow",
    provisioners: {
      gko: gkoScenario<PolicyParams>({
        // The API carries the flow, so it is the parameterized resource.
        manifests: [],
        roles: { api: "policy-api" },
        dynamicRoles: ["api"],
        contextPath: "/policy-api",
        applyParams: async (k, params) => {
          await k.apply(gkoManifest(params));
        },
      }),
      terraform: tfScenario<PolicyParams>({
        fixture: path.join(here, "terraform"),
        toVars: (params) => ({
          with_policy: params.withPolicy,
          flow_name: params.flowName,
          header_name: params.headerName,
          header_value: params.headerValue,
        }),
      }),
    },
    xray: {
      gko: [
        XRAY.POLICIES.DEPLOY_V4_WITH_POLICY,
        XRAY.POLICIES.UPDATE_POLICY,
        XRAY.POLICIES.REMOVE_POLICY,
      ],
      terraform: XRAY.TERRAFORM.API_POLICY_FLOW_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 120_000 },
  },
  async ({ provisioned, mapi, gateway }) => {
    const apiId = await provisioned.apiId();
    const ctx = await provisioned.contextPath();

    /** The API's flow names as APIM records them. */
    const flowNames = async () => ((await mapi.fetchApi(apiId)) as ApiV4).flows?.map((f) => f.name) ?? [];

    await test.step("The policy is recorded and the gateway applies it", async () => {
      await expect
        .poll(flowNames, { timeout: 30_000, message: "flow reaches APIM" })
        .toEqual([INITIAL.flowName]);
      await gateway.assertRespondsWithHeaders(
        ctx,
        { [INITIAL.headerName]: INITIAL.headerValue },
        { status: 200 },
      );
    });

    await test.step("Rewriting the policy changes what the gateway adds", async () => {
      await provisioned.update(UPDATED);
      await expect
        .poll(flowNames, { timeout: 30_000, message: "rewritten flow reaches APIM" })
        .toEqual([UPDATED.flowName]);
      await gateway.assertRespondsWithHeaders(
        ctx,
        // The old header must be GONE, not merely joined by the new one.
        { [UPDATED.headerName]: UPDATED.headerValue, [INITIAL.headerName]: null },
        { status: 200 },
      );
    });

    await test.step("Removing the flow stops the gateway adding the header", async () => {
      await provisioned.update(REMOVED);
      await expect
        .poll(flowNames, { timeout: 30_000, message: "flow removed from APIM" })
        .toEqual([]);
      await gateway.assertRespondsWithHeaders(
        ctx,
        { [UPDATED.headerName]: null },
        { status: 200 },
      );
    });
  },
  INITIAL,
);
