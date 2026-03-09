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
 * Start/Stop API lifecycle test.
 *
 * Mirrors the Chainsaw test at:
 *   test/e2e/chainsaw/tests/apis/state/startStopApi/
 *
 * Preconditions:
 *   - APIM and Gateway are running
 *   - A v4 API exists and is STARTED
 *   - API_ID env var is set to the APIM API id
 *   - API_PATH env var is set to the API's context path (default: /start-stop-test)
 *
 * The test verifies gateway responses at each lifecycle stage.
 * In a real scenario an operator or API call toggles the state
 * between assertions — here each test block asserts the current state.
 */

import { initClients } from "../setup.js";
import type { Mapi, Gateway } from "../../../dist/index.js";

const API_ID = process.env["API_ID"]!;
const API_PATH = process.env["API_PATH"] ?? "/start-stop-test";

let mapi: Mapi;
let gateway: Gateway;

beforeAll(async () => {
  ({ mapi, gateway } = await initClients());
});

describe("Start / Stop API lifecycle", () => {
  it("should respond 200 when the API is started", async () => {
    await mapi.assertApiStarted(API_ID);
    await gateway.assertResponds(API_PATH, { status: 200 });
  });

  it("should respond 404 when the API is stopped", async () => {
    await mapi.assertApiStopped(API_ID);
    await gateway.assertResponds(API_PATH, { status: 404 });
  });

  it("should respond 200 again when the API is re-started", async () => {
    await mapi.assertApiStarted(API_ID);
    await gateway.assertResponds(API_PATH, { status: 200 });
  });
});
