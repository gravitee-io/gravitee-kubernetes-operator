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
 * Provision with the provisioner layer, then assert with the shared assertions.
 *
 * This is the model the GKO e2e suite itself uses. The point of showing it here
 * is that NONE of it is Playwright-specific: `GkoProvisioner` and the `mapi` /
 * `gateway` assertions depend only on `node:*`, `yaml` and native `fetch`, so
 * the same body works under any runner (see ../../jest for the same test).
 *
 * The suite's own `forEachProvisioner` helper additionally runs one body against
 * EVERY provisioner and adds safety-net teardown; it lives in the e2e layer
 * because it binds to Playwright's `test`. Outside that layer you drive the
 * provisioner directly, as below.
 *
 * Requires a cluster with APIM + the GKO operator. Point MANIFEST at a CR of
 * your own:
 *
 *   MANIFEST=/path/to/api.yaml API_NAME=my-api npx playwright test provision-and-assert
 */

import { test } from "@playwright/test";
import { initClients } from "../setup.js";
import { GkoProvisioner, assertProvisioner } from "../../../dist/provisioners/index.js";
import type { Mapi, Gateway, Provisioned } from "../../../dist/index.js";

const MANIFEST = process.env["MANIFEST"];
const API_NAME = process.env["API_NAME"];
const CONTEXT_PATH = process.env["API_PATH"] ?? "/my-api";

let mapi: Mapi;
let gateway: Gateway;
let provisioned: Provisioned<void> | undefined;

test.beforeAll(async () => {
  if (!MANIFEST || !API_NAME) {
    throw new Error("MANIFEST and API_NAME env vars are required");
  }
  ({ mapi, gateway } = await initClients());
});

// Teardown lives in a hook, not in the test body: a body that times out never
// reaches its own cleanup, and a leaked API poisons every later run.
test.afterAll(async () => {
  await provisioned?.destroy();
});

test.describe("Provision through GKO, assert against APIM", () => {
  test("the API is started in APIM and serves traffic", async () => {
    const provisioner = new GkoProvisioner<void>({
      manifests: [MANIFEST!],
      roles: { api: API_NAME! },
      contextPath: CONTEXT_PATH,
    });

    // provision() applies the manifests and waits for each role's Accepted
    // condition, so a successful return already means "GKO landed it".
    provisioned = await provisioner.provision();

    // The provisioner's own record: "did MY layer land this?" — answered from
    // the CR status, never from APIM.
    await assertProvisioner(provisioned, "api", "applied");

    // The platform's answer: "is the resource actually right?" — identical
    // whichever provisioner created it.
    await mapi.assertApiStarted(await provisioned.apiId());
    await gateway.assertResponds(await provisioned.contextPath(), { status: 200 });
  });
});
