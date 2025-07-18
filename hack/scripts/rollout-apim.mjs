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

import { CONFIG } from "./lib/config.mjs";
import {
  LOG,
  PROJECT_DIR,
  toggleVerbosity,
  time,
  isEmptyString,
} from "./lib/index.mjs";

const WORKING_DIR = path.join(os.tmpdir(), CONFIG.repoName);
const VERBOSE = argv.verbose;
const COMMIT_HASH = argv.srcSha;
const SOURCE_BRANCH = argv.srcBRanch;

toggleVerbosity(VERBOSE);

const env = await getEnv();

LOG.magenta(`
🚀 Rolling out deployments in environement ${env} ...
    📦 Project dir    | ${PROJECT_DIR}
    📦 Working dir    | ${WORKING_DIR}`);

await checkRequirements();

async function checkRequirements() {
  if (isEmptyString($.env.CIRCLECI)) {
    LOG.yellow(`
  🤔 It looks like you are trying to run this script locally, while it is meant to be ran in a CI environment.

  If you are sure about what you are doing, set the CIRCLECI environment variable to true.
`);
    process.exit(1);
  }

  if (isEmptyString(COMMIT_HASH)) {
    LOG.yellow(`
  Git commit hash must be set either in CIRCLE_SHA1 or COMMIT_HASH environement variable.
`);
    process.exit(1);
  }

  if (isEmptyString(env)) {
    LOG.yellow(`
  🧐 It looks like the origin branch does not require to rollout any component.
`);
    process.exit(0);
  }
}

LOG.blue(`
    🚧 Checking out ${CONFIG.repo}:${CONFIG.branch} ...
`);

await time(async () => {
  await $`git clone -q git@github.com:${CONFIG.repo}.git \
    --branch ${CONFIG.branch} \
    --single-branch \
    --depth 1 ${WORKING_DIR}`;
});

LOG.blue(`
    🚧 Annotating config values with commit hash ${COMMIT_HASH} ...
`);

await time(async () => {
  cd(WORKING_DIR);
  const apimValuesFilePath = path.join(env, CONFIG.apimValues);
  const apimValuesFile = await fs.readFile(apimValuesFilePath, "utf8");
  const apimValuesYAML = await YAML.parse(apimValuesFile);
  const annotationKey = CONFIG.apimCommitHashAnnotationKey;
  apimValuesYAML.apim.common.annotations[annotationKey] = COMMIT_HASH;
  await fs.writeFile(apimValuesFilePath, YAML.stringify(apimValuesYAML));
  cd(PROJECT_DIR);
});

LOG.blue(`
    🚧 Committing config ...
`);

await time(async () => {
  cd(WORKING_DIR);
  const apimValuesFile = path.join(env, CONFIG.apimValues);
  await $`git add ${apimValuesFile}`;
  await $`git commit -m "ci(${env}): rollout APIM config (${COMMIT_HASH})"`;
  await $`git push origin ${CONFIG.branch}`;
  LOG.log();
  cd(PROJECT_DIR);
});

async function getEnv() {
  if (SOURCE_BRANCH == "master") {
    return "dev";
  }
  if (SOURCE_BRANCH == STABLE.getBranch()) {
    return "stable;";
  }
}
