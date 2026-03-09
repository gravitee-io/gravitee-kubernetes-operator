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
 * Terraform post-apply verification.
 *
 * Inspired by the acceptance tests in the Terraform provider at:
 *   terraform-provider-apim/tests/acceptance/
 *
 * After running `terraform apply`, use this test to verify the API
 * is actually reachable on the gateway and matches expected state in APIM.
 *
 * Usage:
 *   API_ID=$(terraform output -raw api_id) \
 *   API_PATH=$(terraform output -raw api_context_path) \
 *   npx playwright test terraform-post-apply
 */

import { test } from "@playwright/test";
import { initClients } from "../setup.js";
import type { Mapi, Gateway, ApiState } from "../../../dist/index.js";

const API_ID = process.env["API_ID"]!;
const API_PATH = process.env["API_PATH"]!;
const EXPECTED_STATE = (process.env["EXPECTED_STATE"] ?? "STARTED") as ApiState;

let mapi: Mapi;
let gateway: Gateway;

test.beforeAll(async () => {
  if (!API_ID || !API_PATH) {
    throw new Error(
      "API_ID and API_PATH env vars are required. Set them from terraform output.",
    );
  }
  ({ mapi, gateway } = await initClients());
});

test.describe("Terraform post-apply verification", () => {
  test("API should exist in APIM and be in expected state", async () => {
    await mapi.assertApiState(API_ID, EXPECTED_STATE);
  });

  test("API should be reachable on the gateway", async () => {
    await gateway.assertResponds(API_PATH, { status: 200 });
  });
});
