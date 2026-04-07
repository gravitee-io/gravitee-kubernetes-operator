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
 * Defaults: Namespace & SyncFrom default value tests.
 *
 * Xray tests:
 *   GKO-770: V4 API syncFrom defaults to MANAGEMENT when omitted
 *   GKO-765: V2 API local defaults to false when omitted
 *   GKO-463: contextRef namespace defaults to current namespace
 *   GKO-466: Valid name and namespace in contextRef
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

test.describe("Defaults — Namespace & SyncFrom", () => {
  // ── GKO-770: V4 syncFrom defaults to MANAGEMENT ────────────

  test(`V4 API syncFrom defaults to MANAGEMENT when omitted ${XRAY.DEFAULTS.V4_SYNC_FROM_MGMT_DEFAULT} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-no-sync";
    const fixturePath = fixture("crds/defaults/v4-api-no-sync-from.yaml");

    await test.step("Apply V4 API without definitionContext.syncFrom", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API is created in APIM (operator applied default syncFrom)", async () => {
      await mapi.assertApiMatches(apiId, { name: API_NAME });
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-765: V2 local defaults to false ─────────────────────

  test(`V2 API local defaults to false when omitted ${XRAY.DEFAULTS.V2_LOCAL_FALSE_DEFAULT} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v2-no-local";
    const fixturePath = fixture("crds/defaults/v2-api-no-local.yaml");

    await test.step("Apply V2 API without local field", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apidefinition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apidefinition", API_NAME);
    const apiId = status.id;

    await test.step("API is created in APIM (operator applied default local=false)", async () => {
      await mapi.assertApiMatches(apiId, { name: API_NAME });
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-463: Namespace defaults to current ──────────────────

  test(`contextRef namespace defaults to current namespace ${XRAY.DEFAULTS.NAMESPACE_DEFAULT} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-no-ns";
    const fixturePath = fixture("crds/defaults/v4-api-no-namespace-ctx.yaml");

    await test.step("Apply V4 API with contextRef without namespace", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API is created in APIM (namespace defaulted to current)", async () => {
      await mapi.assertApiMatches(apiId, { name: API_NAME });
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-466: Valid name and namespace in contextRef ─────────

  test(`Valid name and namespace in contextRef ${XRAY.DEFAULTS.VALID_NAME_NAMESPACE} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-valid-ctx";
    const fixturePath = fixture("crds/defaults/v4-api-valid-ctx.yaml");

    await test.step("Apply V4 API with explicit name and namespace in contextRef", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API is created in APIM", async () => {
      await mapi.assertApiMatches(apiId, { name: API_NAME });
    });

    await kubectl.del(fixturePath);
  });
});
