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
  PROJECT_DIR,
  isEmptyString,
  Docs,
  toggleVerbosity,
  time,
} from "./lib/index.mjs";
import { Version } from "./lib/version.mjs";

const VERSION = argv.version;
const VERBOSE = argv.verbose;
const RELEASE_CHANGELOG_FILE = argv.file;
const OUTPUT_FILE = argv.output;
const GITHUB_TOKEN = process.env.GITHUB_TOKEN;
const EOL = "\n";
const WORKING_DIR = path.join(os.tmpdir(), "changelog");

toggleVerbosity(VERBOSE);

if (isEmptyString(VERSION)) {
  LOG.red("You must specify a version using the --version flag");
  process.exit(1);
}

if (isEmptyString(RELEASE_CHANGELOG_FILE)) {
  LOG.red("You must specify a file using the --file flag");
  process.exit(1);
}

if (!fs.pathExistsSync(RELEASE_CHANGELOG_FILE)) {
  LOG.red(`File ${RELEASE_CHANGELOG_FILE} could not be found`);
  process.exit(1);
}

if (isEmptyString(OUTPUT_FILE)) {
  LOG.red("You must specify an output file using the --output flag");
  process.exit(1);
}

if (isEmptyString(GITHUB_TOKEN)) {
  LOG.red("A github token is required to submit your pull request");
  process.exit(1);
}

const version = new Version(VERSION);

if (version.isNoPatch()) {
  LOG.yellow(
    `No changelog to generate (version ${VERSION} is a new minor version)`
  );
  process.exit(0);
}

if (version.isPreRelease()) {
  LOG.yellow(
    `No changelog to generate (version ${VERSION} is a pre-release version)`
  );
  process.exit(0);
}

const changelogFile = `${Docs.baseFolder}/${version.minor()}/${
  Docs.changelogFolder
}/gko-${version.branch()}.md`;

LOG.blue(`Writing changelog to ${changelogFile}`);

const releaseChangelog = fs.readFileSync(RELEASE_CHANGELOG_FILE, "utf8").trim();
const changelogHeader = `# GKO ${version.branch()}`;
const prBranch = `release-gko-${VERSION}`;
const prTitle = `[GKO] Changelog for version ${VERSION}`;
const prBody = `
# GKO ${VERSION} has been released

🧐 Please review and merge this pull request to add the changelog to the documentation.
`;

LOG.magenta(`
🚀 Submitting changelog for version ${VERSION} ...
    📦 Project dir    | ${PROJECT_DIR}
    📦 Working dir    | ${WORKING_DIR}
`);

async function checkoutDocs() {
  await $`git clone git@github.com:${Docs.repo}.git \
      --branch main \
      --single-branch \
      --depth 1 ${WORKING_DIR}`;
}

async function commitChangelog() {
  cd(WORKING_DIR);
  await $`git switch -c ${prBranch}`;
  if (fs.pathExistsSync(changelogFile)) {
    LOG.blue("append changelog");
    console.log(appendChangelog());
    fs.writeFileSync(changelogFile, appendChangelog());
  } else {
    fs.writeFileSync(changelogFile, writeNewChangelog());
  }
  await $`git add .`;
  await $`git commit -m "docs: add changelog for gko-${VERSION}"`;
  cd(PROJECT_DIR);
}

async function submitChangelog() {
  cd(WORKING_DIR);
  try {
    await $`gh pr close --delete-branch ${prBranch}`;
  } catch (_) {}
  await $`git push --set-upstream origin ${prBranch}`;
  const prURL =
    await $`gh pr create --title ${prTitle} --body ${prBody} --base main --head ${prBranch}`;
  fs.writeFileSync(OUTPUT_FILE, `${prURL}`);
  cd(PROJECT_DIR);
}

function appendChangelog() {
  const changelog = fs.readFileSync(changelogFile, "utf8");
  return changelog.replace(
    changelogHeader,
    `$&${EOL}${EOL}${releaseChangelog}${EOL}`
  );
}

function writeNewChangelog() {
  return `${changelogHeader}${EOL}${EOL}${releaseChangelog}${EOL}`;
}

LOG.blue(`
    ⎈ Checking out ${Docs.repo} ...
`);

await time(checkoutDocs);

LOG.blue(`
    ⎈ Committing changelog for version ${VERSION} ...
`);

await time(commitChangelog);

LOG.blue(`
    ⎈ Submitting changelog to ${Docs.repo} ...
`);

await time(submitChangelog);
