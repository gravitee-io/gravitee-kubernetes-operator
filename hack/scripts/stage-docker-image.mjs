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

import { LOG, HELM, toggleVerbosity, time, PROJECT_DIR } from "./lib/index.mjs";

const IMG = "graviteeio.azurecr.io/kubernetes-operator";
const VERSION = argv.version || (await HELM.getChartVersion());
const VERBOSE = argv.verbose;

toggleVerbosity(VERBOSE);

LOG.magenta(`
🚀 Staging version ${VERSION} ...

    📦 Project dir    | ${PROJECT_DIR}
    📦 Working dir    | ${PROJECT_DIR}`);

await checkRequirements();

async function checkRequirements() {
  if (!$.env.CIRCLECI) {
    LOG.yellow(`
  🤔 It looks like you are trying to run this script locally, while it is meant to be ran in a CI environment.

  If you are sure about what you are doing, please set the CIRCLECI environment variable to true.

`);
    process.exit(1);
  }
}

cd(PROJECT_DIR);

LOG.blue(`
🐳  Building image ...
`);

await time(async () => $`docker build -t ${IMG}:${VERSION} .`);

LOG.blue(`
🐳  Staging image ...
`);

await time(async () => $`docker push ${IMG}:${VERSION}`);

LOG.magenta(`
  🎉 Image ${IMG}:${VERSION} has been staged !`);
