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
 * Journey: publish an API to the portal and serve traffic.
 *
 * As an API producer, I publish a proxy API to the portal and serve traffic. A V4
 * proxy API, published (PUBLIC + PUBLISHED) and STARTED, is reachable at the
 * gateway; stopping it takes it off the gateway; re-starting serves traffic again.
 *
 * Fixtures are co-located in this folder. The full visibility/lifecycle matrix and
 * V2/native lifecycle stay GKO-only under tests/gko/api-lifecycle (Terraform has no
 * apim_apiv2).
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import { forEachProvisioner } from "../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../helpers/provisioner-env.js";

const here = path.dirname(fileURLToPath(import.meta.url));

/** The single knob: serve traffic (STARTED) or take the API off the gateway (STOPPED). */
interface ApiStateParams {
  state: "STARTED" | "STOPPED";
}

forEachProvisioner<ApiStateParams>(
  {
    title: "Publish an API to the portal and serve traffic",
    provisioners: {
      gko: gkoScenario<ApiStateParams>({
        // The API itself is the parameterized resource (applied by applyParams),
        // so there is no static manifest: provision applies the STARTED variant,
        // update() swaps in the STARTED/STOPPED one.
        manifests: [],
        roles: { api: "published-api" },
        dynamicRoles: ["api"],
        contextPath: "/published-api",
        applyParams: async (k, params) => {
          await k.apply(path.join(here, params.state === "STOPPED" ? "gko/api-stopped.yaml" : "gko/api-started.yaml"));
        },
      }),
      terraform: tfScenario<ApiStateParams>({
        fixture: path.join(here, "terraform"),
        toVars: (params) => ({ state: params.state }),
      }),
    },
    xray: {
      gko: [XRAY.API_LIFECYCLE.DEPLOY_V4_SYNC_K8S, XRAY.API_LIFECYCLE.START_STOP_V2_V4_NATIVE],
      terraform: XRAY.TERRAFORM.API_PUBLISH_SERVE_TRAFFIC_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 90_000 },
  },
  async ({ provisioned, mapi, gateway }) => {
    const apiId = await provisioned.apiId();
    const ctx = await provisioned.contextPath();

    await test.step("Published, started API serves traffic and shows portal state", async () => {
      await mapi.waitForApiMatches(
        apiId,
        { state: "STARTED", visibility: "PUBLIC", lifecycleState: "PUBLISHED" },
        { timeoutMs: 30_000 },
      );
      await gateway.assertResponds(ctx, { status: 200 });
    });

    await test.step("Stopping the API takes it off the gateway", async () => {
      await provisioned.update({ state: "STOPPED" });
      await mapi.waitForApiStopped(apiId, { timeoutMs: 30_000 });
      await gateway.assertResponds(ctx, { status: 404 });
    });

    await test.step("Re-starting the API serves traffic again", async () => {
      await provisioned.update({ state: "STARTED" });
      await mapi.waitForApiStarted(apiId, { timeoutMs: 30_000 });
      await gateway.assertResponds(ctx, { status: 200 });
    });
  },
  { state: "STARTED" },
);
