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
// `--provision-with <p>` (or directly via the env var in CI).
//
// A lane = that provisioner's OWN tests + the matching arm of every shared
// scenario. Legacy single-provisioner tests live under `tests/gko/` or
// `tests/terraform/`; shared `*.scenario.ts` files emit one arm per provisioner,
// each tagged `@gko` / `@terraform`. So we select a lane by (1) IGNORING the other
// provisioner's folder and (2) DROPPING the other arm from shared scenarios via a
// case-sensitive grepInvert (case-sensitive on purpose: Playwright's --grep CLI
// flag is case-insensitive, so a bare `@gko` there would also match every
// `@GKO-1234` Xray tag). This makes `--provision-with gko` run the FULL GKO suite,
// not just the migrated scenarios.
const provisioner = process.env["E2E_PROVISIONER"]?.trim().toLowerCase();
let laneTestIgnore: RegExp | undefined;
let laneGrepInvert: RegExp | undefined;
if (provisioner === "gko") {
  laneTestIgnore = /[/\\]tests[/\\]terraform[/\\]/;
  laneGrepInvert = /@terraform\b/;
} else if (provisioner === "terraform") {
  laneTestIgnore = /[/\\]tests[/\\]gko[/\\]/;
  laneGrepInvert = /@gko\b/;
}
if (provisioner) {
  console.log(`[e2e] provisioner lane: ${provisioner}`);
}

export default defineConfig({
  globalSetup: "./global-setup.ts",
  testDir: "./tests",
  // `*.test.ts` are plain test files; `*.scenario.ts` are provisioner-matrix
  // files that expand into one test per provisioner via forEachProvisioner.
  testMatch: ["**/*.test.ts", "**/*.scenario.ts"],
  testIgnore: laneTestIgnore,
  grepInvert: laneGrepInvert,
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
