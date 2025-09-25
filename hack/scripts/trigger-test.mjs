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

import { triggerPipeline } from "./lib/circleci.mjs";

import { toggleVerbosity, LOG } from "./lib/index.mjs";

import { Version } from "./lib/version.mjs";

const VERBOSE = argv.verbose;
const BRANCH = argv["branch"];
const NOTIFY = argv["notify"] === "true" || argv["notify"] === true;
const APIM_VERSION = argv["apimVersion"];

const apimVersion = new Version(APIM_VERSION);
const gkoVersion = new Version(BRANCH);

if (apimVersion !== "master-latest") {
  if (apimVersion.majorDigit < gkoVersion.majorDigit) {
    LOG.blue(
      `Skipping test because we don't support forward compatibility between major versions`,
    );
    process.exit(0);
  }
  
  if (apimVersion.minorDigit < gkoVersion.minorDigit) {
    LOG.blue(
      `Skipping test because we don't support forward compatibility between minor versions`,
    );
    process.exit(0);
  }
}


toggleVerbosity(VERBOSE);

LOG.blue(`Triggering test pipeline`);

const parameters = {
  trigger: "test",
  "apim-version": APIM_VERSION,
  notify: NOTIFY,
};

LOG.blue(`
  Parameters: ${JSON.stringify(parameters)},
  Branch: ${BRANCH}
`);

const pipelineURL = await triggerPipeline(parameters, BRANCH);

LOG.blue(`Pipeline is running at ${pipelineURL}`);
