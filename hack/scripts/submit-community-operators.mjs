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
  LOG,
  PROJECT_DIR,
  isEmptyString,
  toggleVerbosity,
  time,
} from "./lib/index.mjs";

const VERSION = argv.version;
const VERBOSE = argv.verbose;
const DRY_RUN = argv["dry-run"];
const OUTPUT_FILE = argv.output;
const GITHUB_TOKEN = process.env.GITHUB_TOKEN;

const WORKING_DIR = path.join(os.tmpdir(), "community-operators");

const FORK_REPO = "gravitee-io/community-operators";
const UPSTREAM_REPO = "k8s-operatorhub/community-operators";
const OPERATOR_NAME = "gravitee-kubernetes-operator";
const BUNDLE_DIR = path.join(PROJECT_DIR, "olm", "bundle");

toggleVerbosity(VERBOSE);

if (isEmptyString(VERSION)) {
  LOG.red("You must specify a version using the --version flag");
  process.exit(1);
}

if (isEmptyString(GITHUB_TOKEN) && !DRY_RUN) {
  LOG.red("A github token is required to submit your pull request");
  process.exit(1);
}

if (isEmptyString(OUTPUT_FILE) && !DRY_RUN) {
  LOG.red("You must specify an output file using the --output flag");
  process.exit(1);
}

if (!fs.pathExistsSync(BUNDLE_DIR)) {
  LOG.red(
    `Bundle directory ${BUNDLE_DIR} not found. Run 'make olm-bundle' first.`,
  );
  process.exit(1);
}

const operatorDir = `operators/${OPERATOR_NAME}`;
const versionDir = `${operatorDir}/${VERSION}`;
const prBranch = `gko-${VERSION}`;
const prTitle = `operator ${OPERATOR_NAME} (${VERSION})`;
const prBody = `
### New submission

Submitting \`${OPERATOR_NAME}\` version \`${VERSION}\` to OperatorHub.

**Operator name:** ${OPERATOR_NAME}
**Operator version:** ${VERSION}
**Channel:** alpha

This PR was generated automatically by the GKO release pipeline.
`;

LOG.magenta(`
  Submitting ${OPERATOR_NAME} v${VERSION} to community-operators ...
    Project dir  | ${PROJECT_DIR}
    Working dir  | ${WORKING_DIR}
    Fork repo    | ${FORK_REPO}
    Upstream     | ${UPSTREAM_REPO}
`);

async function checkoutFork() {
  await $`git clone git@github.com:${FORK_REPO}.git \
      --branch main \
      --single-branch \
      --depth 1 ${WORKING_DIR}`;
}

async function copyBundle() {
  cd(WORKING_DIR);
  await $`git switch -c ${prBranch}`;

  const targetDir = path.join(WORKING_DIR, versionDir);
  await fs.ensureDir(targetDir);

  for (const sub of ["manifests", "metadata", "tests"]) {
    const src = path.join(BUNDLE_DIR, sub);
    if (fs.pathExistsSync(src)) {
      await fs.copy(src, path.join(targetDir, sub));
    }
  }

  const ciFile = path.join(WORKING_DIR, operatorDir, "ci.yaml");
  if (!fs.pathExistsSync(ciFile)) {
    await fs.writeFile(
      ciFile,
      YAML.stringify({ updateGraph: "replaces-mode" }),
    );
  }

  await $`git add .`;
  await $`git commit -m "operator ${OPERATOR_NAME} (${VERSION})"`;
  cd(PROJECT_DIR);
}

async function submitPR() {
  cd(WORKING_DIR);
  try {
    await $`gh pr close --repo ${UPSTREAM_REPO} --delete-branch ${prBranch}`;
  } catch (_) {}
  await $`git push --set-upstream origin ${prBranch}`;
  const prURL =
    await $`gh pr create --repo ${UPSTREAM_REPO} --title ${prTitle} --body ${prBody} --base main --head gravitee-io:${prBranch}`;
  fs.writeFileSync(OUTPUT_FILE, `${prURL}`);
  cd(PROJECT_DIR);
}

LOG.blue(`
    Checking out ${FORK_REPO} ...
`);

await time(checkoutFork);

LOG.blue(`
    Copying bundle for version ${VERSION} ...
`);

await time(copyBundle);

if (!DRY_RUN) {
  LOG.blue(`
    Submitting PR to ${UPSTREAM_REPO} ...
  `);

  await time(submitPR);
} else {
  LOG.yellow(`
    Dry run — skipping push and PR submission
  `);
}
