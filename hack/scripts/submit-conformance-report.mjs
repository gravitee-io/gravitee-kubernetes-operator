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
  FORK_REPO,
  UPSTREAM_REPO,
  describeReport,
  findImplementationsList,
  loadReport,
  loadSubmissionMeta,
  patchImplementationsList,
  renderReadme,
  reportPath,
  resolveReportDirVersion,
} from "./lib/conformance.mjs";

import {
  LOG,
  PROJECT_DIR,
  isEmptyString,
  time,
  toggleVerbosity,
} from "./lib/index.mjs";

const VERSION = argv.version;
const VERBOSE = argv.verbose;
const DRY_RUN = argv["dry-run"] === "true" || argv["dry-run"] === true;
const OUTPUT_FILE = argv.output;
const ALLOW_PARTIAL =
  argv["allow-partial"] === true || argv["allow-partial"] === "true";

const WORKING_DIR = path.join(os.tmpdir(), "gateway-api-conformance");

toggleVerbosity(VERBOSE);

if (isEmptyString(VERSION)) {
  LOG.red("You must specify a version using the --version flag");
  process.exit(1);
}

// The upstream PR is opened by a human on purpose. Every conformance report PR
// in kubernetes-sigs/gateway-api is authored by a person, and EasyCLA would
// need a machine account allowlisted under the Gravitee corporate CLA before a
// bot could sign one. This script goes as far as a pushed branch and stops.
const meta = await loadSubmissionMeta();
const report = await loadReport(VERSION);
const desc = describeReport(report, meta);

// We submit for full conformance. A skipped test is far more often a mis-run
// than a real limitation: HTTPRouteMatchingAcrossRoutes skips unless
// GATEWAY_API_MATCH_ACROSS_ROUTES=true is set on the *test binary*, and
// HTTPRouteWeight and HTTPRouteRedirectPortAndScheme skip under
// CONFORMANCE_SKIP_STARVED_TESTS. Publishing such a report understates what GKO
// supports and is hard to walk back once merged upstream, so partial takes an
// explicit --allow-partial.
if (!desc.conformant && !ALLOW_PARTIAL) {
  LOG.red(`
  Refusing to submit: the ${desc.version} report is a PARTIAL conformance report.

    Skipped: ${desc.skipped.join(", ")}

  Check the run before assuming these are real limitations:
    - GATEWAY_API_MATCH_ACROSS_ROUTES=true must be set on the test binary
    - CONFORMANCE_SKIP_STARVED_TESTS must be unset
    - CONFORMANCE_RERUN=0, so no test passed only on a retry

  Re-run the suite and regenerate the report, or pass --allow-partial if the
  limitation is genuine and you intend to publish it.
`);
  process.exit(1);
}

const reportFile = `standard-${desc.version}-default-report.yaml`;
const prBranch = `gko-${desc.version}-${desc.gatewayAPIVersion}`;
const prTitle = `Gravitee Kubernetes Operator ${desc.version} conformance report for ${desc.gatewayAPIVersion.replace(/^v/, "")}`;
const prBody = `**What type of PR is this?**

/kind documentation
/area conformance-test

**What this PR does / why we need it**:

Adds ${meta.name} version ${desc.version} conformance report for gateway API ${desc.gatewayAPIVersion.replace(/^v/, "")}

**Which issue(s) this PR fixes**:

Fixes #

**Does this PR introduce a user-facing change?**:

\`\`\`release-note
NONE
\`\`\`
`;

LOG.magenta(`
  Preparing ${meta.name} ${desc.version} conformance submission ...
    Working dir    | ${WORKING_DIR}
    Fork repo      | ${FORK_REPO}
    Upstream       | ${UPSTREAM_REPO}
    Gateway API    | ${desc.gatewayAPIVersion} (${desc.channel} channel, ${desc.mode} mode)
    Profile        | ${desc.profileName}
    Conformance    | ${desc.conformant ? "full" : `partial (skipped: ${desc.skipped.join(", ")})`}
    Branch         | ${prBranch}
`);

// Resolved against the fork once it is checked out, because upstream symlinks
// patch versions at a per-minor directory. Set before anything is rendered:
// the badge has to advertise the same path the report is committed to.
let reportDir;

async function checkoutFork() {
  await fs.remove(WORKING_DIR);
  await $`git clone git@github.com:${FORK_REPO}.git \
      --single-branch \
      --depth 1 ${WORKING_DIR}`;
  cd(WORKING_DIR);
  await $`git remote add upstream https://github.com/${UPSTREAM_REPO}.git`;
  await $`git fetch upstream main --depth 1`;
  await $`git reset --hard upstream/main`;
  cd(PROJECT_DIR);

  desc.reportDirVersion = resolveReportDirVersion(
    WORKING_DIR,
    desc.gatewayAPIVersion,
  );
  reportDir = `conformance/reports/${desc.reportDirVersion}/${meta.directory}`;
  LOG.blue(`  Reports directory: ${reportDir}`);
}

async function applyChanges() {
  cd(WORKING_DIR);
  await $`git switch -c ${prBranch}`;

  const targetDir = path.join(WORKING_DIR, reportDir);
  await fs.ensureDir(targetDir);

  // The report is copied verbatim: it is the evidence, never regenerated here.
  await fs.copy(reportPath(desc.version), path.join(targetDir, reportFile));
  await fs.writeFile(
    path.join(targetDir, "README.md"),
    renderReadme(desc, meta),
  );

  const listFile = findImplementationsList(WORKING_DIR);
  const listContent = await fs.readFile(listFile, "utf8");
  await fs.writeFile(
    listFile,
    patchImplementationsList(listContent, desc, meta),
  );

  await $`git add .`;
  await $`git commit --signoff -m ${prTitle}`;

  const stat = await $`git show --stat --oneline HEAD`;
  LOG.blue(`\n${stat.stdout}`);

  cd(PROJECT_DIR);
}

async function pushBranch() {
  cd(WORKING_DIR);
  await $`git push --force --set-upstream origin ${prBranch}`;
  cd(PROJECT_DIR);
}

LOG.blue(`\n    Checking out ${FORK_REPO} ...\n`);
await time(checkoutFork);

LOG.blue(`\n    Applying report, README and implementations list edits ...\n`);
await time(applyChanges);

const [forkOrg, forkName] = FORK_REPO.split("/");
const compareURL = `https://github.com/${UPSTREAM_REPO}/compare/main...${forkOrg}:${forkName}:${prBranch}?expand=1`;

if (DRY_RUN) {
  LOG.yellow(`
    Dry run — branch built at ${WORKING_DIR} but not pushed.
    Inspect it with:

      git -C ${WORKING_DIR} show HEAD
  `);
} else {
  LOG.blue(`\n    Pushing ${prBranch} to ${FORK_REPO} ...\n`);
  await time(pushBranch);

  LOG.green(`
    Branch pushed. Open the pull request yourself:

      ${compareURL}

    Title:

      ${prTitle}
  `);
  LOG.log(`\n${prBody}`);
}

if (!isEmptyString(OUTPUT_FILE)) {
  await fs.writeFile(
    OUTPUT_FILE,
    JSON.stringify(
      {
        version: desc.version,
        gatewayAPIVersion: desc.gatewayAPIVersion,
        conformant: desc.conformant,
        skipped: desc.skipped,
        branch: prBranch,
        pushed: !DRY_RUN,
        compareURL,
        title: prTitle,
        body: prBody,
      },
      null,
      2,
    ),
  );
  LOG.green(`  Submission summary written to ${OUTPUT_FILE}`);
}
