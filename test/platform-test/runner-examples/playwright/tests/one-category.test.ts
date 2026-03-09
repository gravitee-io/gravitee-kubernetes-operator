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
 * One Category test.
 *
 * Mirrors the Chainsaw test at:
 *   test/e2e/chainsaw/tests/apis/categories/oneCategory/
 *
 * Preconditions:
 *   - APIM is running
 *   - An API exists with a category assigned
 *   - API_ID env var is set to the APIM API id
 *   - CATEGORY_NAME env var is set to the expected category name
 */

import { test } from "@playwright/test";
import { initClients } from "../setup.js";
import type { Mapi } from "../../../dist/index.js";

const API_ID = process.env["API_ID"]!;
const CATEGORY_NAME = process.env["CATEGORY_NAME"] ?? "my-category";

let mapi: Mapi;

test.beforeAll(async () => {
  ({ mapi } = await initClients());
});

test.describe("API with one category", () => {
  test("should have the expected category via partial match", async () => {
    await mapi.assertApiMatches(API_ID, {
      categories: [CATEGORY_NAME],
    });
  });
});
