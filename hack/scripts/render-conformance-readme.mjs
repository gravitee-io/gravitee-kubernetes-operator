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
  REPORT_DIR,
  describeReport,
  loadReport,
  loadSubmissionMeta,
  renderReadme,
} from "./lib/conformance.mjs";

import {
  LOG,
  PROJECT_DIR,
  isEmptyString,
  toggleVerbosity,
} from "./lib/index.mjs";

const VERBOSE = argv.verbose;
const CHECK = argv.check === true || argv.check === "true";

toggleVerbosity(VERBOSE);

// Default to the version the suite currently reports as, so a release bump of
// impl.yaml is the only place a version is written by hand.
async function currentVersion() {
  const implFile = path.join(
    PROJECT_DIR,
    "test",
    "conformance",
    "kubernetes.io",
    "gateway-api",
    "impl",
    "impl.yaml",
  );
  return YAML.parse(await fs.readFile(implFile, "utf8")).version;
}

const VERSION = isEmptyString(argv.version)
  ? await currentVersion()
  : String(argv.version);

const meta = await loadSubmissionMeta();
const report = await loadReport(VERSION);
const desc = describeReport(report, meta);
const readme = renderReadme(desc, meta);

// See submit-conformance-report.mjs: a partial report is usually a mis-run, not
// a limitation. Render it, but never let it pass unremarked.
if (!desc.conformant) {
  LOG.yellow(`
  WARNING: the ${VERSION} report is PARTIAL (skipped: ${desc.skipped.join(", ")}).

  Confirm these are real limitations and not a mis-run before committing this
  README — GATEWAY_API_MATCH_ACROSS_ROUTES unset on the test binary, or
  CIRCLECI=true, each silently skip tests GKO actually passes.
`);
}

const target = path.join(REPORT_DIR, "README.md");

if (CHECK) {
  const current = fs.pathExistsSync(target)
    ? await fs.readFile(target, "utf8")
    : "";
  if (current !== readme) {
    LOG.red(
      `${target} is out of date for version ${VERSION}.\n` +
        `Run: npx zx hack/scripts/render-conformance-readme.mjs`,
    );
    process.exit(1);
  }
  LOG.green(`README is up to date for version ${VERSION}`);
} else {
  await fs.writeFile(target, readme);
  LOG.green(`Rendered ${target} for version ${VERSION}`);
}
