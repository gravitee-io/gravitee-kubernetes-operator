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

import { LOG } from "./lib/index.mjs";

const FORK_REPO = "gravitee-io-labs/community-operators";
const UPSTREAM_REPO = "k8s-operatorhub/community-operators";
const OPERATOR_NAME = "gravitee-kubernetes-operator";
const WORKING_DIR = path.join(os.tmpdir(), "community-operators-rebase");

process.env.GH_TOKEN = process.env.GH_TOKEN || process.env.GITHUB_TOKEN;

LOG.blue("Searching for open GKO PRs on OperatorHub ...");

const prsRaw =
  await $`gh search prs --repo ${UPSTREAM_REPO} ${OPERATOR_NAME} \
    --author graviteeio --state open \
    --json number,title,headRefName --limit 20`;

const prs = JSON.parse(prsRaw.stdout);

if (prs.length === 0) {
  LOG.green("No open PRs found. Nothing to do.");
  process.exit(0);
}

LOG.blue(`Found ${prs.length} open PR(s). Checking CI status ...`);

const stale = [];
for (const pr of prs) {
  const checksRaw =
    await $`gh pr view ${pr.number} --repo ${UPSTREAM_REPO} \
      --json statusCheckRollup \
      --jq ${`[.statusCheckRollup[] | select(.name == "operator-ci")] | map({conclusion}) | first`}`;

  let check;
  try {
    check = JSON.parse(checksRaw.stdout);
  } catch {
    continue;
  }

  if (check && check.conclusion === "FAILURE") {
    stale.push(pr);
  }
}

if (stale.length === 0) {
  LOG.green("All PRs have passing or pending CI. Nothing to rebase.");
  process.exit(0);
}

LOG.yellow(
  `${stale.length} PR(s) need rebase: ${stale.map((p) => p.headRefName).join(", ")}`,
);

LOG.blue("Cloning fork ...");
await fs.remove(WORKING_DIR);
await $`git clone git@github.com:${FORK_REPO}.git --no-checkout --filter=blob:none ${WORKING_DIR}`;
cd(WORKING_DIR);
await $`git remote add upstream https://github.com/${UPSTREAM_REPO}.git`;
await $`git fetch upstream main`;

for (const pr of stale) {
  const branch = pr.headRefName;
  LOG.blue(`  Rebasing ${branch} (PR #${pr.number}) ...`);

  try {
    await $`git fetch origin ${branch}`;
    await $`git checkout -B ${branch} origin/${branch}`;
    await $`git rebase upstream/main`;
    await $`git push --force origin ${branch}`;
    LOG.green(`  Rebased and pushed ${branch}`);
  } catch (e) {
    try {
      await $`git rebase --abort`;
    } catch {
      // rebase may not be in progress
    }
    LOG.red(`  Failed to rebase ${branch}: ${e.message}`);
  }
}

LOG.green("Done.");
