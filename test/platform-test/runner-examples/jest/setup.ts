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

// In a real project you would import from "@gravitee/platform-test".
// Here we use a relative path since this example lives inside the library repo.
import {
  loadGraviteeConfig,
  createMapiFromConfig,
  type Mapi,
  type Gateway,
} from "../../dist/index.js";

import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));

let mapi: Mapi;
let gateway: Gateway;

/**
 * Initialise shared Mapi and Gateway clients from config.yaml.
 * Call this in beforeAll() of each test file.
 *
 * Config resolution order:
 *   1. Environment variables (GRAVITEE_BASE_URL, GRAVITEE_GATEWAY_URL, etc.)
 *   2. config.yaml in the platform-test directory
 */
export async function initClients(): Promise<{ mapi: Mapi; gateway: Gateway }> {
  if (mapi) return { mapi, gateway };

  const configPath = path.resolve(__dirname, "../../config.yaml");
  const config = await loadGraviteeConfig(configPath);

  mapi = createMapiFromConfig(config);
  gateway = mapi.gateway({
    baseUrl: config.gateway?.baseUrl ?? "http://localhost:30082",
  });

  return { mapi, gateway };
}
