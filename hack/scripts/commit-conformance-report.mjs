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

import {
  REPORT_DIR,
  describeReport,
  loadReport,
  loadSubmissionMeta,
  renderReadme,
  reportPath,
} from "./lib/conformance.mjs";

import {
  LOG,
  PROJECT_DIR,
  isEmptyString,
  toggleVerbosity,
} from "./lib/index.mjs";

import { Version } from "./lib/version.mjs";

const VERSION = argv.version;
const VERBOSE = argv.verbose;
const DRY_RUN = argv["dry-run"] === "true" || argv["dry-run"] === true;

toggleVerbosity(VERBOSE);

if (isEmptyString(VERSION)) {
  LOG.red("You must specify a version using the --version flag");
  process.exit(1);
}

const meta = await loadSubmissionMeta();
const report = await loadReport(VERSION);
const desc = describeReport(report, meta);

// The suite exiting zero is not the bar. gotestsum is happy with a run that
// skipped tests, and a skipped test is the exact failure mode that produced the
// stale partial reports this pipeline replaces. Only a report with nothing
// failed and nothing skipped is allowed to land.
if (!desc.conformant) {
  LOG.red(`
  Not committing: the ${desc.version} run is not fully conformant.

    Skipped: ${desc.skipped.join(", ")}

  The report is still attached to the job as an artifact. Check the run before
  treating these as real limitations:
    - GATEWAY_API_MATCH_ACROSS_ROUTES=true must be set on the test binary
    - CONFORMANCE_SKIP_STARVED_TESTS must be unset
    - CONFORMANCE_RERUN=0, so nothing passed only on a retry
`);
  process.exit(1);
}

// A report claims a version. If the branch we are about to commit to declares a
// different one, the claim would be false, so stop rather than publish it.
if (desc.version !== VERSION) {
  LOG.red(
    `Report declares version ${desc.version} but ${VERSION} was requested. ` +
      `Check test/conformance/kubernetes.io/gateway-api/impl/impl.yaml on this branch.`,
  );
  process.exit(1);
}

const branch = new Version(VERSION).branch();
const reportFile = path.basename(reportPath(VERSION));
const readmeFile = path.join(REPORT_DIR, "README.md");

LOG.magenta(`
  Committing conformance report ...
    Version        | ${desc.version}
    Gateway API    | ${desc.gatewayAPIVersion} (${desc.channel} channel, ${desc.mode} mode)
    Profile        | ${desc.profileName}
    Conformance    | full — ${desc.passed} passed, 0 failed, 0 skipped
    Branch         | ${branch}
`);

// Hold the generated report in memory across the branch switch rather than
// relying on it surviving as an untracked file.
const generated = await fs.readFile(reportPath(VERSION), "utf8");

cd(PROJECT_DIR);

// A report only ever belongs on its own release branch. If that branch does not
// exist we are being asked to publish a version that was never released, so say
// so rather than letting git fail with "couldn't find remote ref".
try {
  await $`git fetch origin ${branch}`;
} catch {
  LOG.red(
    `No release branch ${branch} on origin, so there is nowhere to commit the ` +
      `${VERSION} report. Check the version you asked for.`,
  );
  process.exit(1);
}

await $`git switch ${branch} 2>/dev/null || git switch -c ${branch} origin/${branch}`;

await fs.writeFile(reportPath(VERSION), generated);
await fs.writeFile(readmeFile, renderReadme(desc, meta));

await $`git add ${reportPath(VERSION)} ${readmeFile}`;

const staged = await $`git diff --cached --name-only`;
if (isEmptyString(staged.stdout)) {
  LOG.yellow(
    `  Report for ${VERSION} is already committed and unchanged — nothing to do.`,
  );
  process.exit(0);
}

await $`git commit -m ${`ci: add conformance report for ${VERSION} [skip ci]`}`;

const shown = await $`git show --stat --oneline HEAD`;
LOG.blue(`\n${shown.stdout}`);

if (DRY_RUN) {
  LOG.yellow(`\n  Dry run — commit built on ${branch} but not pushed.\n`);
} else {
  await $`git push origin ${branch}`;
  LOG.green(`\n  Pushed ${reportFile} to ${branch}\n`);
}
