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
 * API Lifecycle: Start / Stop tests.
 *
 * Xray tests:
 *   GKO-69:   Deploy a valid CRD for V4 Proxy API with syncFrom Kubernetes
 *   GKO-1464: Validate API lifecycle start/stop across V2, V4, and Native APIs
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

const API_NAME = "e2e-v4-start-stop";
const API_PATH = "/e2e-v4-start-stop";

test(`API Lifecycle — Start / Stop ${XRAY.API_LIFECYCLE.DEPLOY_V4_SYNC_K8S} ${XRAY.API_LIFECYCLE.START_STOP_V2_V4_NATIVE} ${TAGS.REGRESSION}`, async ({
  kubectl,
  mapi,
  gateway,
}) => {
  await test.step("Deploy V4 Proxy API with syncFrom Kubernetes", async () => {
    await kubectl.apply(fixture("crds/api-v4-definitions/v4-proxy-api-started.yaml"));
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
  });

  const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
  const apiId = status.id;

  await test.step("API responds 200 when STARTED", async () => {
    await mapi.waitForApiMatches(apiId, {
      name: API_NAME,
      state: "STARTED",
    });
    await gateway.assertResponds(API_PATH, { status: 200 });
  });

  await test.step("API responds 404 when STOPPED", async () => {
    await kubectl.apply(fixture("crds/api-v4-definitions/v4-proxy-api-stopped.yaml"));
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    await mapi.waitForApiStopped(apiId);
    await gateway.assertResponds(API_PATH, { status: 404 });
  });

  await test.step("API responds 200 again when re-STARTED", async () => {
    await kubectl.apply(fixture("crds/api-v4-definitions/v4-proxy-api-started.yaml"));
    await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");

    await mapi.waitForApiStarted(apiId);
    await gateway.assertResponds(API_PATH, { status: 200 });
  });

  await kubectl.del(fixture("crds/api-v4-definitions/v4-proxy-api-started.yaml"));
});
