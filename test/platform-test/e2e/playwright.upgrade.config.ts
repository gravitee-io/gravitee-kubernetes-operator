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

/**
 * Dedicated config for the upgrade SURVIVAL specs (`tests/upgrade/*.spec.ts`).
 *
 * These run in two phases across SEPARATE processes (before and after an in-place
 * GKO + APIM upgrade) and are intentionally kept out of the normal suite: the
 * default config's testMatch is `*.test.ts` / `*.scenario.ts`, so these `*.spec.ts`
 * files are never collected there. The two phases are selected by the
 * `e2e:upgrade:before` / `e2e:upgrade:after` npm scripts via a filename filter.
 *
 * globalSetup is the same as the normal suite (infra checks + the shared dev-ctx).
 * The timeout is larger because a survival step waits on the operator to reconcile
 * resources carried over a version change.
 */
export default defineConfig({
  globalSetup: "./global-setup.ts",
  testDir: "./tests/upgrade",
  testMatch: ["**/survival.*.spec.ts"],
  timeout: 120_000,
  retries: 0,
  workers: 1,
  reporter: [
    ["html", { open: "never" }],
    ["list"],
    ["junit", { outputFile: path.join(__dirname, "../playwright-results/results.xml") }],
  ],
  projects: [{ name: "upgrade" }],
});
