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
 * Journey: secure an API with a plan.
 *
 * As an API producer, I secure my API with a JWT plan and an OAuth2 plan. Both
 * plan security types must land in APIM, whichever provisioner created the API.
 * Gateway enforcement with real tokens (subscription flows) stays GKO-only; this
 * journey proves the plan security configuration round-trips through both drivers.
 *
 * Fixtures are co-located in this folder. Plan admission validation (e.g. general
 * conditions, GKO-238) stays GKO-only under tests/gko/policies.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import { forEachProvisioner } from "../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../helpers/provisioner-env.js";

const here = path.dirname(fileURLToPath(import.meta.url));

forEachProvisioner(
  {
    title: "Secure an API with a JWT plan and an OAuth2 plan",
    provisioners: {
      gko: gkoScenario<void>({
        manifests: [path.join(here, "gko/api-with-plans.yaml")],
        roles: { api: "secured-api" },
      }),
      terraform: tfScenario<void>({ fixture: path.join(here, "terraform") }),
    },
    xray: {
      gko: [XRAY.PLANS.OAUTH2_PLAN_V4, XRAY.PLANS.JWT_PLAN_V4],
      terraform: XRAY.TERRAFORM.API_SECURE_WITH_PLAN_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 60_000 },
  },
  async ({ provisioned, mapi }) => {
    const apiId = await provisioned.apiId();

    // The shared invariant: APIM records a published plan for each declared
    // security type, regardless of which driver authored it.
    await expect
      .poll(
        async () => {
          const plans = await mapi.listApiPlans(apiId);
          return {
            jwt: plans.some((p) => p.security?.type === "JWT"),
            oauth2: plans.some((p) => p.security?.type === "OAUTH2"),
          };
        },
        { timeout: 30_000, message: "JWT and OAuth2 plans are present in APIM" },
      )
      .toMatchObject({ jwt: true, oauth2: true });
  },
);
