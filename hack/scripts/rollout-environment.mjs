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

import { ENVIRONMENTS } from "./lib/env.mjs";
import {
  LOG,
  PROJECT_DIR,
  toggleVerbosity,
  time,
  isEmptyString,
} from "./lib/index.mjs";

const WORKING_DIR = path.join(os.tmpdir(), ENVIRONMENTS.configRepoName);

const VERBOSE = argv.verbose;
const ENV = argv.env;
const COMMIT_HASH = $.env.CIRCLE_SHA1 || $.env.COMMIT_HASH;

toggleVerbosity(VERBOSE);

LOG.magenta(`
🚀 Rolling out deployments in environement ${ENV} ...
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

  if (isEmptyString(ENV)) {
    LOG.yellow(`
  ⚠️ Setting an environment via the --env flag is required. This flag must match one of the root directories of ${ENVIRONMENTS.configRepo}.
`);
    process.exit(1);
  }
}

LOG.blue(`
    🚧 Checking out ${ENVIRONMENTS.configRepo}:${ENVIRONMENTS.configBranch} ...
`);

await time(async () => {
  await $`git clone -q git@github.com:${ENVIRONMENTS.configRepo}.git \
    --branch ${ENVIRONMENTS.configBranch} \
    --single-branch \
    --depth 1 ${WORKING_DIR}`;
});

LOG.blue(`
    🚧 Annotating config values with commit hash ${COMMIT_HASH} ...
`);

await time(async () => {
  cd(WORKING_DIR);
  const gkoValuesFilePath = path.join(ENV, ENVIRONMENTS.gkoValues);
  const gkoValuesFile = await fs.readFile(gkoValuesFilePath, "utf8");
  const gkoValuesYAML = await YAML.parse(gkoValuesFile);
  const annotationKey = ENVIRONMENTS.commitHashAnnotationKey;
  gkoValuesYAML.gko.manager.annotations[annotationKey] = COMMIT_HASH;
  await fs.writeFile(gkoValuesFilePath, YAML.stringify(gkoValuesYAML));
  cd(PROJECT_DIR);
});

LOG.blue(`
    🚧 Committing config ...
`);

await time(async () => {
  cd(WORKING_DIR);
  const gkoValuesFile = path.join(ENV, ENVIRONMENTS.gkoValues);
  await $`git add ${gkoValuesFile}`;
  await $`git commit -m "ci(${ENV}): rollout config"`;
  await $`git push origin ${ENVIRONMENTS.configBranch}`;
  LOG.log();
  cd(PROJECT_DIR);
});
