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
 * Journey: publish a message (event) API and choose how consumers read it.
 *
 * As an event API producer, I publish a MESSAGE (event) API and expose it over
 * the entrypoint my consumers need — polling (http-get), publishing (http-post),
 * a server-sent-event stream, or a websocket — and I change that choice, along
 * with the API's description and version, without recreating the API.
 *
 * The entrypoint type is a variant of ONE story, so it is a variant table here
 * rather than a folder per type. Each variant also carries its own description
 * and version, which makes every re-apply prove the descriptive fields change in
 * place too. The webhook subscription entrypoint is constant across the variants.
 *
 * Fixtures are co-located in this folder. Actually consuming messages (SSE /
 * webhook delivery) is not covered anywhere in the suite: no test receives
 * messages today, which needs a subscriber fixture the cluster does not provide.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test, expect } from "../../../../setup.js";
import type { ApiV4 } from "../../../../../src/types/apim.js";
import { XRAY, TAGS } from "../../../../helpers/tags.js";
import { forEachProvisioner } from "../../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../../helpers/provisioner-env.js";

const here = path.dirname(fileURLToPath(import.meta.url));

/** The HTTP entrypoint the API exposes, plus the revision that ships with it. */
interface MessageApiParams {
  httpEntrypoint: "http-get" | "http-post" | "sse" | "websocket";
  version: string;
}

/**
 * The four HTTP entrypoint types a MESSAGE API can expose. Each carries a
 * distinct version so every re-apply also has descriptive fields to move.
 */
const VARIANTS: MessageApiParams[] = [
  { httpEntrypoint: "http-get", version: "1.0.0" },
  { httpEntrypoint: "http-post", version: "2.0.0" },
  { httpEntrypoint: "sse", version: "3.0.0" },
  { httpEntrypoint: "websocket", version: "4.0.0" },
];

/** Declared verbatim in both fixtures, so the assertion can be shared. */
function description(params: MessageApiParams): string {
  return `V4 MESSAGE (event) API exposing a ${params.httpEntrypoint} entrypoint`;
}

/** Every entrypoint type APIM should report, sorted. */
function expectedEntrypoints(params: MessageApiParams): string[] {
  return [params.httpEntrypoint, "webhook"].sort();
}

forEachProvisioner<MessageApiParams>(
  {
    title: "Publish a message (event) API over each consumer entrypoint",
    provisioners: {
      gko: gkoScenario<MessageApiParams>({
        // The API is the parameterized resource: provision applies the first
        // variant, update() re-applies another over the same CR name.
        manifests: [],
        roles: { api: "message-api" },
        dynamicRoles: ["api"],
        applyParams: async (k, params) => {
          await k.apply(path.join(here, `gko/message-api-${params.httpEntrypoint}.yaml`));
        },
      }),
      terraform: tfScenario<MessageApiParams>({
        fixture: path.join(here, "terraform"),
        toVars: (params) => ({
          http_entrypoint: params.httpEntrypoint,
          api_version: params.version,
        }),
      }),
    },
    xray: {
      gko: [
        XRAY.MESSAGE_APIS.DEPLOY_V4_MSG_SYNC_MGMT,
        XRAY.MESSAGE_APIS.DEPLOY_V4_MSG_SYNC_K8S,
        XRAY.MESSAGE_APIS.HTTP_GET_ENTRYPOINT,
        XRAY.MESSAGE_APIS.HTTP_POST_ENTRYPOINT,
        XRAY.MESSAGE_APIS.SSE_ENTRYPOINT,
        XRAY.MESSAGE_APIS.WEBHOOK_ENTRYPOINT,
        XRAY.MESSAGE_APIS.WEBSOCKET_ENTRYPOINT,
        // Every variant reads from the same mock message endpoint group.
        XRAY.MESSAGE_APIS.MOCK_ENDPOINT,
        XRAY.API_LIFECYCLE.UPDATE_V4_MESSAGE_API,
      ],
      terraform: XRAY.TERRAFORM.MESSAGE_API_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 120_000 },
  },
  async ({ provisioned, mapi }) => {
    const apiId = await provisioned.apiId();

    for (const [index, params] of VARIANTS.entries()) {
      await test.step(`MESSAGE API exposed over ${params.httpEntrypoint}`, async () => {
        // Variant 0 is what provision() already applied.
        if (index > 0) await provisioned.update(params);

        // One atomic poll: type/state, the descriptive fields and the entrypoint
        // set all have to agree on the SAME read, otherwise a half-applied
        // update could pass on a lucky interleaving of two reads.
        await expect
          .poll(
            async () => {
              const api = (await mapi.fetchApi(apiId)) as ApiV4;
              return {
                type: api.type,
                state: api.state,
                apiVersion: api.apiVersion,
                description: api.description,
                entrypoints: api.listeners
                  .flatMap((l) => l.entrypoints ?? [])
                  .map((e) => e.type)
                  .sort(),
              };
            },
            {
              timeout: 30_000,
              message: `MESSAGE API is started and exposes ${params.httpEntrypoint}`,
            },
          )
          .toEqual({
            type: "MESSAGE",
            state: "STARTED",
            apiVersion: params.version,
            description: description(params),
            entrypoints: expectedEntrypoints(params),
          });
      });
    }
  },
  VARIANTS[0],
);
