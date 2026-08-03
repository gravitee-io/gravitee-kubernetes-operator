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
import { PROVISIONER_LANES } from "../src/provisioners/registry.js";

const __dirname = path.dirname(fileURLToPath(import.meta.url));

// Optional provisioner-lane filter, set by `scripts/e2e.mjs` from
// `--provision-with <p>` (or directly via the env var in CI).
//
// A lane = every test that runs through that provisioner: its own tests under
// `tests/<id>/` plus the matching arm of each shared `*.scenario.ts`. Selection
// is purely by TITLE TAG — every provisioner-specific test carries its
// `PROVISIONER_LANES[].tag` and `forEachProvisioner` appends it to each
// generated arm — so a lane is chosen by dropping every OTHER lane's tag.
// Nothing keys off the folder a test lives in, which is what lets the tree be
// reorganised without touching lane logic.
//
// The grep is case-SENSITIVE on purpose: Playwright's `--grep` CLI flag is
// case-insensitive, so a bare `@gko` there would also match every `@GKO-1234`
// Xray tag and select the whole suite.
const provisioner = process.env["E2E_PROVISIONER"]?.trim().toLowerCase();
const lane = PROVISIONER_LANES.find((l) => l.id === provisioner);
const otherTags = lane
  ? PROVISIONER_LANES.filter((l) => l.id !== lane.id)
      .map((l) => l.tag)
      .join("|")
  : "";
const laneGrepInvert: RegExp | undefined = otherTags
  ? new RegExp(otherTags + String.raw`\b`)
  : undefined;
if (provisioner) {
  // stderr, not stdout: the JSON reporter writes its report to stdout, and a
  // stray log line there makes the output unparseable.
  console.error(`[e2e] provisioner lane: ${provisioner}`);
}

export default defineConfig({
  globalSetup: "./global-setup.ts",
  testDir: "./tests",
  // `*.test.ts` are plain test files; `*.scenario.ts` are provisioner-matrix
  // files that expand into one test per provisioner via forEachProvisioner.
  testMatch: ["**/*.test.ts", "**/*.scenario.ts"],
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
