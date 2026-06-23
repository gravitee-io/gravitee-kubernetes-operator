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

import { defineConfig } from "@playwright/test";
import { fileURLToPath } from "node:url";
import path from "node:path";

const __dirname = path.dirname(fileURLToPath(import.meta.url));

// Optional provisioner-lane filter, set by `scripts/e2e.mjs` from
// `--provision-with <p>` (or directly via the env var in CI). Built as a
// CASE-SENSITIVE RegExp on purpose: Playwright's --grep CLI flag is
// case-insensitive, so a bare `@gko` there would also match every `@GKO-1234`
// Xray tag and select the whole suite. A case-sensitive `@gko` matches only the
// lowercase provisioner tag that the matrix arms and *-gko-only files carry.
const provisioner = process.env["E2E_PROVISIONER"]?.trim().toLowerCase();
if (provisioner) {
  console.log(`[e2e] provisioner lane: @${provisioner}`);
}

export default defineConfig({
  globalSetup: "./global-setup.ts",
  testDir: "./tests",
  // `*.test.ts` are plain test files; `*.scenario.ts` are provisioner-matrix
  // files that expand into one test per provisioner via forEachProvisioner.
  testMatch: ["**/*.test.ts", "**/*.scenario.ts"],
  grep: provisioner ? new RegExp(String.raw`@${provisioner}\b`) : undefined,
  timeout: 30_000,
  retries: 0,
  workers: 1,
  reporter: [
    ["html", { open: "never" }],
    ["list"],
    ["junit", { outputFile: path.join(__dirname, "../playwright-results/results.xml") }],
  ],
  projects: [
    {
      name: "platform-test",
    },
  ],
});
