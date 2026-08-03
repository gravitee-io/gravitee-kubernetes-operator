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
 * The single entry point for running the suite. Translates the suite's custom
 * selection flags into the env vars / config that Playwright reads, then runs
 * `playwright test` (so globalSetup always runs) and forwards every other
 * argument untouched. The custom flags are orthogonal and combine freely:
 *
 *   npm run e2e -- --provision-with gko --run-up-to-version 4.12.0 [playwright args]
 *
 *   --provision-with <id>             Run only one provisioner lane (see
 *                                     src/provisioners/registry.ts for the
 *                                     known ids). -> E2E_PROVISIONER
 *                                     (case-sensitive @tag grep).
 *   --run-up-to-version <semver>      Run only features available at that version,
 *                                     i.e. skip tests tagged @since-<newer>.
 *                                     -> E2E_MAX_VERSION (enforced in e2e/setup.ts).
 *   --suite <name>                    Which suite to run (default "e2e"). "upgrade"
 *                                     selects the survival specs, which run in two
 *                                     phases across SEPARATE processes either side
 *                                     of an in-place GKO + APIM upgrade, and so use
 *                                     their own config + longer timeout.
 *   --phase <before|after>            Required by --suite upgrade: which side of the
 *                                     upgrade to run.
 *   <anything else>                   Forwarded verbatim to `playwright test`
 *                                     (e.g. --grep @GKO-2828, --headed, a file path).
 *
 * The env vars also work directly, which is handy in CI matrices:
 *   E2E_PROVISIONER=gko npm run e2e
 */

import { spawn } from "node:child_process";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { PROVISIONER_ORDER } from "../dist/provisioners/registry.js";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const PROVISIONERS = PROVISIONER_ORDER;

/**
 * Suites differ only in which Playwright config they use and whether they take a
 * phase. Everything else (lane selection, version capping, passthrough) is shared,
 * so there is one way to launch the suite regardless of which one you want.
 */
const SUITES = {
  e2e: { config: "../e2e/playwright.config.ts" },
  upgrade: { config: "../e2e/playwright.upgrade.config.ts", phases: ["before", "after"] },
};

function die(message) {
  console.error(`[e2e] ${message}`);
  process.exit(2);
}

const VALUE_FLAGS = ["--provision-with", "--run-up-to-version", "--suite", "--phase"];

const args = process.argv.slice(2);
const env = {};
const passthrough = [];
let suiteName = "e2e";
let phase;

for (let i = 0; i < args.length; i++) {
  const arg = args[i];
  let name = arg;
  let inlineValue;
  if (arg.startsWith("--") && arg.includes("=")) {
    const eq = arg.indexOf("=");
    name = arg.slice(0, eq);
    inlineValue = arg.slice(eq + 1);
  }

  if (!VALUE_FLAGS.includes(name)) {
    passthrough.push(arg);
    continue;
  }

  let value = inlineValue;
  if (value === undefined) {
    value = args[i + 1];
    if (value === undefined || value.startsWith("-")) die(`${name} requires a value`);
    i++; // consume the separate value token
  }

  switch (name) {
    case "--provision-with": {
      const provisioner = value.toLowerCase();
      if (!PROVISIONERS.includes(provisioner)) {
        die(`unknown provisioner "${value}". Known: ${PROVISIONERS.join(", ")}`);
      }
      env.E2E_PROVISIONER = provisioner;
      break;
    }
    case "--run-up-to-version":
      env.E2E_MAX_VERSION = value;
      break;
    case "--suite":
      suiteName = value.toLowerCase();
      if (!SUITES[suiteName]) {
        die(`unknown suite "${value}". Known: ${Object.keys(SUITES).join(", ")}`);
      }
      break;
    case "--phase":
      phase = value.toLowerCase();
      break;
  }
}

const suite = SUITES[suiteName];

if (suite.phases) {
  if (!phase) die(`--suite ${suiteName} requires --phase <${suite.phases.join("|")}>`);
  if (!suite.phases.includes(phase)) {
    die(`unknown phase "${phase}" for suite "${suiteName}". Known: ${suite.phases.join(", ")}`);
  }
  // The phases are separate spec files, selected by filename.
  passthrough.push(`survival.${phase}`);
} else if (phase) {
  const phased = Object.keys(SUITES).filter((n) => SUITES[n].phases);
  die(`--phase is only meaningful with a phased suite (${phased.join(", ")})`);
}

if (env.E2E_MAX_VERSION) {
  console.log(
    `[e2e] capping at APIM version ${env.E2E_MAX_VERSION}: skipping tests tagged @since-<newer>`,
  );
}
if (suiteName !== "e2e") {
  console.log(`[e2e] suite: ${suiteName}${phase ? ` (${phase} phase)` : ""}`);
}

const config = path.resolve(__dirname, suite.config);
const child = spawn("npx", ["playwright", "test", "--config", config, ...passthrough], {
  stdio: "inherit",
  env: { ...process.env, ...env },
});
child.on("exit", (code) => process.exit(code ?? 1));
