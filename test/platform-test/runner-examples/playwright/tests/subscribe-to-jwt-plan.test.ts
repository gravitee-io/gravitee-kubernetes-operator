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
 * Subscribe to JWT Plan test.
 *
 * Mirrors the Chainsaw test at:
 *   test/e2e/chainsaw/tests/scenarios/subscribeToJwtPlan/
 *
 * Preconditions:
 *   - APIM and Gateway are running
 *   - A v4 API with a JWT plan is deployed
 *   - An application and subscription exist
 *   - API_PATH env var is set (default: /jwt-demo)
 *   - JWT_TOKEN env var contains a valid JWT for the plan
 *
 * The test checks gateway behaviour with and without a valid JWT.
 */

import { test } from "@playwright/test";
import { initClients } from "../setup.js";
import type { Gateway } from "../../../dist/index.js";

const API_PATH = process.env["API_PATH"] ?? "/jwt-demo";
const JWT_TOKEN = process.env["JWT_TOKEN"]!;

let gateway: Gateway;

test.beforeAll(async () => {
  ({ gateway } = await initClients());
});

test.describe("JWT plan subscription", () => {
  test("should return 401 without a token", async () => {
    await gateway.assertResponds(API_PATH, { status: 401 });
  });

  test("should return 200 with a valid JWT", async () => {
    await gateway.assertResponds(API_PATH, {
      status: 200,
      headers: { Authorization: `Bearer ${JWT_TOKEN}` },
    });
  });

  test("should return 401 after subscription is removed", async () => {
    // This assertion assumes the subscription has been deleted externally
    // between the previous test and this one.
    await gateway.assertResponds(API_PATH, {
      status: 401,
      headers: { Authorization: `Bearer ${JWT_TOKEN}` },
    });
  });
});
