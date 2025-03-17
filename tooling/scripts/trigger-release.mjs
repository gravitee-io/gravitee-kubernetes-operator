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

import { toggleVerbosity, isEmptyString, LOG } from "./lib/index.mjs";

const VERSION = argv.version;
const VERBOSE = argv.verbose;
const DRY_RUN = argv["dry-run"] === "true" || argv["dry-run"] === true;
const LATEST = argv["latest"] === "true" || argv["latest"] === true;
const PIPELINE_BRANCH = argv["pipeline-branch"];
const TRIGGER = "release";

toggleVerbosity(VERBOSE);

if (isEmptyString(VERSION)) {
  LOG.red("You must specify a version using the --version flag");
  process.exit(1);
}

LOG.blue(`Triggering release pipeline for version ${VERSION}`);

const parameters = {
  trigger: TRIGGER,
  "release-version": VERSION,
  "dry-run": DRY_RUN,
  latest: LATEST,
};

LOG.blue(`Parameters: ${JSON.stringify(parameters)}`);

const pipelineURL = await triggerPipeline(parameters, PIPELINE_BRANCH);

LOG.blue(`Pipeline is running at ${pipelineURL}`);
