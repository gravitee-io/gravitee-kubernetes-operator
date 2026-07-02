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
 * Terraform workspace lifecycle for the e2e suite.
 *
 * The pure mechanics live in the runner-agnostic `@gravitee/platform-test`
 * library (`src/provisioners/engines/terraform-core.ts`) so the Terraform
 * provisioner can reuse them. This adapter keeps the e2e-specific bits the core
 * must not own: loading `config.yaml` for the APIM auth env, and resolving
 * fixture folder names via `fixture()`. It preserves the original
 * `initWorkspace(fixtureName)` signature and re-exports the rest unchanged so
 * existing terraform tests keep working.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { loadGraviteeConfig } from "../../src/cmd/config.js";
import { fixture } from "../setup.js";
import {
  initWorkspace as initWorkspaceCore,
  type TfWorkspace,
} from "../../src/provisioners/engines/terraform-core.js";

const __dirname = path.dirname(fileURLToPath(import.meta.url));

export {
  TF_TIMEOUT_MS,
  TF_WORKSPACE_TIMEOUT_MS,
  tf,
  apply,
  plan,
  writeVars,
  applyExpectFailure,
  output,
  destroy,
  destroyWorkspace,
} from "../../src/provisioners/engines/terraform-core.js";
export type { TfWorkspace } from "../../src/provisioners/engines/terraform-core.js";

/**
 * Build the APIM auth/server environment terraform needs, from `config.yaml`
 * (env vars override config fields, same as the rest of the suite).
 */
export async function terraformEnv(): Promise<Record<string, string>> {
  const configPath = path.resolve(__dirname, "../../config.yaml");
  const config = await loadGraviteeConfig(configPath);
  const baseUrl = config.apim?.baseUrl ?? "http://localhost:30083";

  return {
    ...(process.env as Record<string, string>),
    APIM_SERVER_URL: `${baseUrl}/automation`,
    APIM_USERNAME: config.apim?.auth?.username ?? "admin",
    APIM_PASSWORD: config.apim?.auth?.password ?? "admin",
  };
}

/**
 * Create a Terraform workspace from a fixture folder name (e.g.
 * "subscriptions/apikey-auto"). Resolves the fixture dir + APIM env, then
 * delegates to the core engine.
 */
export async function initWorkspace(fixtureName: string): Promise<TfWorkspace> {
  const env = await terraformEnv();
  // Co-located journey fixtures are passed as absolute paths and pass through
  // unchanged; legacy relative names are rooted at e2e/fixtures.
  const dir = path.isAbsolute(fixtureName) ? fixtureName : fixture(fixtureName);
  return initWorkspaceCore(dir, env);
}
