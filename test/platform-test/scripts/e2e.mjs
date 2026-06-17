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
 * Thin e2e entry point. Translates the suite's custom selection flags into env
 * vars that e2e/playwright.config.ts reads, then runs `playwright test` with the
 * config (so globalSetup always runs) and forwards every other argument
 * untouched. The custom flags are orthogonal and combine freely:
 *
 *   npm run e2e -- --provision-with gko --run-up-to-version 4.12.0 [playwright args]
 *
 *   --provision-with <gko|terraform>  Run only one provisioner lane.
 *                                     -> E2E_PROVISIONER (case-sensitive @tag grep).
 *   --run-up-to-version <semver>      Run only features available at that version.
 *                                     -> E2E_MAX_VERSION. STUBBED: accepted but not
 *                                        enforced yet (version-gating is future work).
 *   <anything else>                   Forwarded verbatim to `playwright test`
 *                                     (e.g. --grep @GKO-2828, --headed, a file path).
 *
 * The env vars also work directly, which is handy in CI matrices:
 *   E2E_PROVISIONER=gko npm run e2e
 */

import { spawn } from "node:child_process";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const CONFIG = path.resolve(__dirname, "../e2e/playwright.config.ts");
const PROVISIONERS = ["gko", "terraform"];

function die(message) {
  console.error(`[e2e] ${message}`);
  process.exit(2);
}

const args = process.argv.slice(2);
const env = {};
const passthrough = [];

for (let i = 0; i < args.length; i++) {
  const arg = args[i];
  let name = arg;
  let inlineValue;
  if (arg.startsWith("--") && arg.includes("=")) {
    const eq = arg.indexOf("=");
    name = arg.slice(0, eq);
    inlineValue = arg.slice(eq + 1);
  }

  if (name === "--provision-with" || name === "--run-up-to-version") {
    let value = inlineValue;
    if (value === undefined) {
      value = args[i + 1];
      if (value === undefined || value.startsWith("-")) die(`${name} requires a value`);
      i++; // consume the separate value token
    }
    if (name === "--provision-with") {
      const provisioner = value.toLowerCase();
      if (!PROVISIONERS.includes(provisioner)) {
        die(`unknown provisioner "${value}". Known: ${PROVISIONERS.join(", ")}`);
      }
      env.E2E_PROVISIONER = provisioner;
    } else {
      env.E2E_MAX_VERSION = value;
    }
  } else {
    passthrough.push(arg);
  }
}

if (env.E2E_MAX_VERSION) {
  console.warn(
    `[e2e] --run-up-to-version=${env.E2E_MAX_VERSION} is accepted but NOT enforced yet: ` +
      `version-gating (@since-<semver> tags + skip) is not implemented. Running the full selection.`,
  );
}

const child = spawn("npx", ["playwright", "test", "--config", CONFIG, ...passthrough], {
  stdio: "inherit",
  env: { ...process.env, ...env },
});
child.on("exit", (code) => process.exit(code ?? 1));
