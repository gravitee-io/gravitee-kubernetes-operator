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

const VERBOSE = argv.verbose;
const PIPELINE_BRANCH = argv["pipeline-branch"];

toggleVerbosity(VERBOSE);

LOG.blue(`Triggering E2E test pipeline`);

const parameters = {
  trigger: "e2e",
};

LOG.blue(`
  Parameters: ${JSON.stringify(parameters)},
`);

const pipelineURL = await triggerPipeline(parameters, PIPELINE_BRANCH);

LOG.blue(`Pipeline is running at ${pipelineURL}`);
