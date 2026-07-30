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
// A lane = that provisioner's OWN tests (if any) + the matching arm of every
// shared scenario. Legacy single-provisioner tests live under `tests/<segment>/`
// per `PROVISIONER_LANES[].testDirSegment`; shared `*.scenario.ts` files emit
// one arm per provisioner, each tagged with `PROVISIONER_LANES[].tag`. So we
// select a lane by (1) IGNORING every OTHER lane's legacy folder and (2)
// DROPPING every other lane's arm from shared scenarios via a case-sensitive
// grepInvert (case-sensitive on purpose: Playwright's --grep CLI flag is
// case-insensitive, so a bare `@gko` there would also match every `@GKO-1234`
// Xray tag). This makes `--provision-with gko` run the FULL GKO suite, not
// just the migrated scenarios. A lane with no `testDirSegment` (e.g. a
// provisioner that only ever appears via shared scenarios) simply contributes
// no ignore pattern for its own folder.
const provisioner = process.env["E2E_PROVISIONER"]?.trim().toLowerCase();
const lane = PROVISIONER_LANES.find((l) => l.id === provisioner);
const otherLanes = lane ? PROVISIONER_LANES.filter((l) => l.id !== lane.id) : [];
const otherTestDirSegments = otherLanes
  .map((l) => l.testDirSegment)
  .filter((segment): segment is string => Boolean(segment));
const laneTestIgnore: RegExp | undefined = otherTestDirSegments.length
  ? new RegExp(String.raw`[/\\]tests[/\\](${otherTestDirSegments.join("|")})[/\\]`)
  : undefined;
const otherTags = otherLanes.map((l) => l.tag).join("|");
const laneGrepInvert: RegExp | undefined = otherTags
  ? new RegExp(otherTags + String.raw`\b`)
  : undefined;
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
