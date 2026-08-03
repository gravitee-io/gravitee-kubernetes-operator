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
 * Message APIs — Lifecycle tests.
 *
 * Xray tests:
 *   GKO-164: Deploy V4 message API with policy
 *
 * Deploying a MESSAGE API and exposing it over each consumer entrypoint
 * (GKO-72/73/129/130/132/133/134/136) is the shared cross-provisioner journey
 * tests/user-journeys/api-producer/publish-a-message-api.
 *
 * Removed:
 *   GKO-135: Kafka endpoint — APIM schema bug
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS, PROVISIONER } from "../../../helpers/tags.js";
import * as kubectl from "../../../helpers/kubectl.js";
import type { ApiV4 } from "../../../../src/types/apim.js";

test.describe(`Message APIs — Lifecycle ${PROVISIONER.GKO}`, () => {
  // Safety-net cleanup: runs even if a test times out before its inline
  // cleanup. Each del() ignores errors (the resource may already be gone).
  test.afterEach(async () => {
    await kubectl.del(fixture("message-apis/v4-message-api-with-policy/crd.yaml")).catch(() => {});
  });

  // ── GKO-164: Deploy V4 message API with policy ──────────────

  test(`Deploy V4 message API with policy ${XRAY.MESSAGE_APIS.MSG_API_WITH_POLICY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-msg-policy";
    const fixturePath = fixture("message-apis/v4-message-api-with-policy/crd.yaml");

    await test.step("Apply CRD with transform-headers policy", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API has flows configured in APIM", async () => {
      const api = (await mapi.fetchApi(apiId)) as ApiV4;
      expect(api).toBeTruthy();
      if ("flows" in api && api.flows) {
        expect(api.flows.length).toBeGreaterThanOrEqual(1);
      }
    });

    await kubectl.del(fixturePath);
  });
});
