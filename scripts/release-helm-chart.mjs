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
  HELM,
  toggleVerbosity,
  time,
  isEmptyString,
} from "./lib/index.mjs";

const WORKING_DIR = path.join(os.tmpdir(), "helm-charts");
const PROJECT_DIR = path.join(__dirname, "..");

const VERSION = argv.version;
const VERBOSE = argv.verbose;
const DRY_RUN = argv["dry-run"];

const GITHUB_TOKEN = $.env.GITHUB_TOKEN;

toggleVerbosity(VERBOSE);

LOG.magenta(`
🚀 Releasing version ${VERSION} ...
    📦 Project dir    | ${PROJECT_DIR}
    📦 Working dir    | ${WORKING_DIR}`);

await checkRequirements();

async function checkRequirements() {
  if (isEmptyString(VERSION)) {
    LOG.red("You must specify a version to release using the --version flag");
    await $`exit 1`;
  }

  if (isEmptyString(GITHUB_TOKEN) && !DRY_RUN) {
    LOG.red(
      "A github token is needed to push the release. Please set the GITHUB_TOKEN environment variable."
    );
    await $`exit 1`;
  }

  if (!$.env.CIRCLECI) {
    LOG.yellow(`
  ⚠️ it looks like you are trying to run this script locally, while it is meant to be ran in a CI environment.

  If you are sure you want to continue, please set the CIRCLECI environment variable to true.

`);
    await $`exit 1`;
  }
}

LOG.blue(`
⎈ Checking out ${HELM.chartsRepo}:${HELM.releaseBranch} ...
`);

await time(checkoutHelmCharts);

async function checkoutHelmCharts() {
  await $`git clone git@github.com:${HELM.chartsRepo}.git \
    --branch ${HELM.releaseBranch} \
    --single-branch \
    --depth 1 ${WORKING_DIR}`;
}

LOG.blue(`
⎈ Packaging chart ...
`);

await time(packageChart);

async function packageChart() {
  await $`helm package -d ${WORKING_DIR}/helm/gko ${HELM.chartDir} --app-version ${VERSION} --version ${VERSION}`;
}

LOG.blue(`
⎈ Indexing repository ...
`);

await time(indexRepo);

async function indexRepo() {
  await $`helm repo index \
      --url https://helm.gravitee.io/helm/gko \
      --merge ${WORKING_DIR}/index.yaml ${WORKING_DIR}/helm/gko`;

  await $`mv ${WORKING_DIR}/helm/gko/index.yaml ${WORKING_DIR}/index.yaml`;
}

LOG.blue(`
⎈ Committing release ...
`);

if (!DRY_RUN) {
  await time(publishRelease);
} else {
  LOG.yellow(`  ⚠️ This is a dry run, release will not be committed ..
  `);
}

async function publishRelease() {
  cd(WORKING_DIR);
  await $`git remote set-url origin "https://${GITHUB_TOKEN}@github.com/${HELM.chartsRepo}.git"`;
  await $`git add helm/gko/gko-${VERSION}.tgz index.yaml`;
  await $`git commit -m "chore(gko): release version ${VERSION}"`;
  await $`git push origin ${HELM.releaseBranch}`;
  cd(PROJECT_DIR);
}

LOG.blue(`
⎈ Setting chart version to ${VERSION} ...
`);

await time(setChartVersion);

async function setChartVersion() {
  const chartFile = await fs.readFile(`${HELM.chartDir}/Chart.yaml`, "utf8");
  const chartYaml = await YAML.parse(chartFile);
  chartYaml.version = VERSION;
  chartYaml.appVersion = VERSION;
  await fs.writeFile(`${HELM.chartDir}/Chart.yaml`, YAML.stringify(chartYaml));
}

if (!DRY_RUN) {
  LOG.magenta(`
🎉 version ${VERSION} has been released !`);
} else {
  LOG.magenta(`
🎉 dry run done for version ${VERSION}`);
}
