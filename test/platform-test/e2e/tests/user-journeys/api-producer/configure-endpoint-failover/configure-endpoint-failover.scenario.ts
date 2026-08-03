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
 * Journey: give an API a failover policy.
 *
 * As an API producer, I put a failover policy on my API's endpoint group so a
 * slow or failing backend is retried instead of surfacing as an error, and I
 * tighten the retry budget later without republishing the API.
 *
 * Failover is a nested block on the API, which is where the regression risk is:
 * a whole block can round-trip to APIM with default values silently substituted,
 * so the journey asserts every field it declares, and asserts them again after
 * an update that changes exactly one of them.
 *
 * Fixtures are co-located in this folder. Exercising a real endpoint failure at
 * the gateway needs a backend the suite can take down on demand, which the test
 * cluster does not provide; the traffic assertion here only proves the policy
 * does not break the happy path.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test } from "../../../../setup.js";
import { XRAY, TAGS } from "../../../../helpers/tags.js";
import { forEachProvisioner } from "../../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../../helpers/provisioner-env.js";

const here = path.dirname(fileURLToPath(import.meta.url));

/** The one knob: how many times the gateway retries a failing endpoint. */
interface FailoverParams {
  maxRetries: 3 | 5;
}

/** Declared identically in both fixtures; only maxRetries is parameterized. */
const FAILOVER_POLICY = {
  enabled: true,
  slowCallDuration: 500,
  openStateDuration: 5000,
  maxFailures: 5,
  perSubscription: false,
};

forEachProvisioner<FailoverParams>(
  {
    title: "Configure endpoint failover on an API",
    provisioners: {
      gko: gkoScenario<FailoverParams>({
        // The API carries the parameterized failover block, so it is applied by
        // applyParams rather than as a static manifest.
        manifests: [],
        roles: { api: "failover-api" },
        dynamicRoles: ["api"],
        contextPath: "/failover-api",
        applyParams: async (k, params) => {
          await k.apply(path.join(here, `gko/api-retries-${params.maxRetries}.yaml`));
        },
      }),
      terraform: tfScenario<FailoverParams>({
        fixture: path.join(here, "terraform"),
        toVars: (params) => ({ max_retries: params.maxRetries }),
      }),
    },
    xray: {
      gko: XRAY.API_LIFECYCLE.FAILOVER_V4_PROXY,
      terraform: XRAY.TERRAFORM.API_FAILOVER_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 90_000 },
  },
  async ({ provisioned, mapi, gateway }) => {
    const apiId = await provisioned.apiId();
    const ctx = await provisioned.contextPath();

    await test.step("APIM records the whole failover policy", async () => {
      await mapi.waitForApiMatches(
        apiId,
        { state: "STARTED", failover: { ...FAILOVER_POLICY, maxRetries: 3 } },
        { timeoutMs: 30_000 },
      );
    });

    await test.step("The API with failover still serves traffic", async () => {
      await gateway.assertResponds(ctx, { status: 200 });
    });

    await test.step("Tightening the retry budget updates the policy in place", async () => {
      await provisioned.update({ maxRetries: 5 });
      await mapi.waitForApiMatches(
        apiId,
        { failover: { ...FAILOVER_POLICY, maxRetries: 5 } },
        { timeoutMs: 30_000 },
      );
    });
  },
  { maxRetries: 3 },
);
