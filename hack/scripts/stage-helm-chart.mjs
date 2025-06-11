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
import { Version } from "./lib/version.mjs";

const WORKING_DIR = path.join(os.tmpdir(), "helm-charts");
const PROJECT_DIR = path.join(__dirname, "..");

const VERSION = await HELM.getChartVersion();
const VERBOSE = argv.verbose;

const ACR_USER_NAME = $.env.ACR_USER_NAME;
const ACR_PASSWORD = $.env.ACR_PASSWORD;

toggleVerbosity(VERBOSE);

LOG.magenta(`
🚀 Staging version ${VERSION} ...
    📦 Project dir    | ${PROJECT_DIR}
    📦 Working dir    | ${WORKING_DIR}`);

await checkRequirements();

async function checkRequirements() {
  // if (isEmptyString(ACR_USER_NAME)) {
  //   LOG.red(
  //     "Azure container registry username is mandatory. Please set the ACR_USER_NAME environment variable.",
  //   );
  //   process.exit(1);
  // }

  // if (isEmptyString(ACR_PASSWORD)) {
  //   LOG.red(
  //     "Azure container registry password is mandatory. Please set the ACR_PASSWORD environment variable.",
  //   );
  //   process.exit(1);
  // }

  if (!$.env.CIRCLECI) {
    LOG.yellow(`
  ⚠️ it looks like you are trying to run this script locally, while it is meant to be ran in a CI environment.

  If you are sure you want to continue, please set the CIRCLECI environment variable to true.

`);
    process.exit(1);
  }
}

LOG.blue(`
⎈ Packaging chart ...
`);

await time(packageChart);

async function packageChart() {
  await $`helm package -d ${HELM.chartDir} ${HELM.chartDir} --app-version ${VERSION} --version ${VERSION}`;
}

LOG.blue(`
⎈ Staging chart ...
`);

await time(stageChart);

async function stageChart() {
  // await $`helm registry login graviteeio.azurecr.io --username $ACR_USER_NAME --password $ACR_PASSWORD`;
  await $`helm push ${HELM.chartDir}/gko-${VERSION}.tgz oci://graviteeio.azurecr.io/helm/`;
}

LOG.magenta(`
  🎉 version ${VERSION} has been staged !`);
