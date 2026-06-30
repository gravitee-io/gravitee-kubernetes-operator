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
 * Use case: stand up a message (event) API.
 *
 * A V4 MESSAGE API created through either provisioner is recorded in APIM as
 * type MESSAGE and reaches STARTED. The entrypoint-type matrix
 * (SSE/webhook/websocket consumption) and policy variants stay GKO-only under
 * tests/gko/message-apis.
 *
 * Fixtures live in fixtures/use-cases/consume-message-api/.
 */

import { expect } from "../../../setup.js";
import type { ApiV4 } from "../../../../src/types/apim.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import { forEachProvisioner } from "../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../helpers/provisioner-env.js";

const FIXTURE = "use-cases/consume-message-api";

forEachProvisioner(
  {
    title: "Stand up a message (event) API",
    provisioners: {
      gko: gkoScenario<void>({
        manifests: [`${FIXTURE}/gko/message-api.yaml`],
        roles: { api: "e2e-uc-message" },
      }),
      terraform: tfScenario<void>({ fixture: `${FIXTURE}/terraform` }),
    },
    xray: {
      gko: [XRAY.MESSAGE_APIS.DEPLOY_V4_MSG_SYNC_MGMT, XRAY.MESSAGE_APIS.DEPLOY_V4_MSG_SYNC_K8S],
      terraform: XRAY.TERRAFORM.MESSAGE_API_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 60_000 },
  },
  async ({ provisioned, mapi }) => {
    const apiId = await provisioned.apiId();

    // The shared invariant: APIM records a started MESSAGE API, whichever driver
    // authored it.
    await expect
      .poll(
        async () => {
          const api = (await mapi.fetchApi(apiId)) as ApiV4;
          return { type: api.type, state: api.state };
        },
        { timeout: 30_000, message: "MESSAGE API is started in APIM" },
      )
      .toMatchObject({ type: "MESSAGE", state: "STARTED" });
  },
);
