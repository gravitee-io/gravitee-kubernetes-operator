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
 * Assert that the provisioner lanes exactly PARTITION the suite: every test runs
 * in exactly one lane, none in two, none in zero.
 *
 * Lane selection is by title tag alone (see e2e/playwright.config.ts), and every
 * way of getting a tag wrong is invisible in a full-suite run:
 *
 *   - no provisioner tag  -> survives every lane's grepInvert, so it runs in ALL
 *                            lanes (the same test executed once per lane)
 *   - the wrong tag       -> runs only in some other provisioner's lane
 *   - two provisioner tags-> dropped by every lane, so it runs in NONE
 *
 * All three still pass the untagged full run, so nothing catches them until a
 * lane-scoped CI job quietly over- or under-covers. This check makes them loud.
 * It only lists tests, so it needs no cluster.
 *
 *   node scripts/check-lanes.mjs
 */

import { execFile } from "node:child_process";
import path from "node:path";
import { promisify } from "node:util";
import { fileURLToPath } from "node:url";
import { PROVISIONER_ORDER } from "../dist/provisioners/registry.js";

const execFileAsync = promisify(execFile);
const __dirname = path.dirname(fileURLToPath(import.meta.url));
const CONFIG = path.resolve(__dirname, "../e2e/playwright.config.ts");

/** Every test's fully-qualified title, via `playwright test --list --reporter=json`. */
async function listTests(provisioner) {
  const { stdout } = await execFileAsync(
    "npx",
    ["playwright", "test", "--config", CONFIG, "--list", "--reporter=json"],
    {
      env: provisioner ? { ...process.env, E2E_PROVISIONER: provisioner } : process.env,
      maxBuffer: 64 * 1024 * 1024,
    },
  );

  const titles = new Set();
  const walk = (suite, prefix) => {
    const here = `${prefix}>${suite.title ?? ""}`;
    for (const child of suite.suites ?? []) walk(child, here);
    for (const spec of suite.specs ?? []) titles.add(`${here}>${spec.title}`);
  };
  for (const suite of JSON.parse(stdout).suites ?? []) walk(suite, "");
  return titles;
}

const all = await listTests(undefined);
const lanes = new Map();
for (const id of PROVISIONER_ORDER) lanes.set(id, await listTests(id));

const laneOf = new Map();
const inMultipleLanes = [];
for (const [id, titles] of lanes) {
  for (const title of titles) {
    if (laneOf.has(title)) inMultipleLanes.push({ title, lanes: [laneOf.get(title), id] });
    laneOf.set(title, id);
  }
}
const inNoLane = [...all].filter((title) => !laneOf.has(title));

const summary = [...lanes].map(([id, t]) => `${id}=${t.size}`).join("  ");
console.log(`[check-lanes] total=${all.size}  ${summary}`);

if (inNoLane.length || inMultipleLanes.length) {
  for (const title of inNoLane) {
    console.error(`[check-lanes] in NO lane (missing a provisioner tag?): ${title}`);
  }
  for (const { title, lanes: ids } of inMultipleLanes) {
    console.error(`[check-lanes] in lanes ${ids.join(" + ")}: ${title}`);
  }
  process.exit(1);
}

console.log("[check-lanes] OK: the lanes partition the suite exactly");
