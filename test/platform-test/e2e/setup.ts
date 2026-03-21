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

import { test as base } from "@playwright/test";
import {
  loadGraviteeConfig,
  createMapiFromConfig,
  type Mapi,
  type Gateway,
} from "../src/index.js";
import * as kubectl from "./helpers/kubectl.js";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));

type E2EFixtures = {
  mapi: Mapi;
  gateway: Gateway;
  kubectl: typeof kubectl;
};

let sharedMapi: Mapi | undefined;
let sharedGateway: Gateway | undefined;

async function initClients(): Promise<{ mapi: Mapi; gateway: Gateway }> {
  if (sharedMapi && sharedGateway) return { mapi: sharedMapi, gateway: sharedGateway };

  const configPath = path.resolve(__dirname, "../config.yaml");
  const config = await loadGraviteeConfig(configPath);

  sharedMapi = createMapiFromConfig(config);
  sharedGateway = sharedMapi.gateway({
    baseUrl: config.gateway?.baseUrl ?? "http://localhost:30082",
  });

  return { mapi: sharedMapi, gateway: sharedGateway };
}

export const test = base.extend<E2EFixtures>({
  mapi: async ({}, use) => {
    const { mapi } = await initClients();
    await use(mapi);
  },
  gateway: async ({}, use) => {
    const { gateway } = await initClients();
    await use(gateway);
  },
  kubectl: async ({}, use) => {
    await use(kubectl);
  },
});

export { expect } from "@playwright/test";

/** Resolve a fixture path relative to the e2e/fixtures directory. */
export function fixture(...segments: string[]): string {
  return path.resolve(__dirname, "fixtures", ...segments);
}
