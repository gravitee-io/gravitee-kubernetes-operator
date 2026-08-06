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

import { LOG, isEmptyString, toggleVerbosity } from "./lib/index.mjs";

import { Version } from "./lib/version.mjs";

const VERSION = argv.version;
const VERBOSE = argv.verbose;
const DRY_RUN = argv["dry-run"] !== false && argv["dry-run"] !== "false";

toggleVerbosity(VERBOSE);

if (isEmptyString(VERSION)) {
  LOG.red("You must specify a version using the --version flag");
  process.exit(1);
}

// The suite runs from the release branch, not the tag: the pipeline needs the
// conformance tooling from that branch to run at all, and impl.yaml on the
// branch is what stamps the version into the report. The job cross checks that
// impl.yaml really does declare the requested version before anything lands.
const branch = new Version(VERSION).branch();

const parameters = {
  trigger: "conformance-report",
  "conformance-version": VERSION,
  "dry-run": DRY_RUN,
};

LOG.blue(`
  Triggering conformance report pipeline
    Version    | ${VERSION}
    Branch     | ${branch}
    Dry run    | ${DRY_RUN}
`);

const pipelineURL = await triggerPipeline(parameters, branch);

LOG.blue(`Pipeline is running at ${pipelineURL}`);
