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

import { LOG, PROJECT_DIR, isEmptyString, Docs } from "./index.mjs";
import { Version } from "./version.mjs";

const GITHUB_TOKEN = process.env.GITHUB_TOKEN;

export async function prepareDocs(version) {
  const releasedVersion = new Version(version);
  const nextVersion = releasedVersion.nextMinor();

  if (isEmptyString(GITHUB_TOKEN)) {
    LOG.red("A github token is required to submit your pull request");
    process.exit(1);
  }

  const workingDir = path.join(os.tmpdir(), "gravitee-platform-docs");

  LOG.magenta(`
  🚀 Preparing documentation for version ${nextVersion} ...
      📦 Project dir    | ${PROJECT_DIR}
      📦 Working dir    | ${workingDir}
  `);

  await checkoutDocs(workingDir);

  await copyDocs(workingDir, releasedVersion, nextVersion);

  await prepareNewRelease(workingDir, nextVersion);

  await prepareReleaseNote(workingDir, releasedVersion, nextVersion);

  await prepareChangelog(workingDir, releasedVersion, nextVersion);

  await prepareAPIReference(workingDir, nextVersion);

  await prepareSummary(workingDir, releasedVersion, nextVersion);

  await commitChanges(workingDir, nextVersion);

  await submitChanges(workingDir, nextVersion);
}

async function checkoutDocs(workingDir) {
  await $`git clone git@github.com:${Docs.repo}.git \
      --branch main \
      --single-branch \
      --depth 1 ${workingDir}`;
}

async function copyDocs(workingDir, releasedVersion, nextVersion) {
  const releasedDocs = path.join(
    workingDir,
    Docs.baseFolder,
    releasedVersion.minor(),
  );
  const nextDocs = path.join(workingDir, Docs.baseFolder, nextVersion.minor());
  await $`cp -r ${releasedDocs} ${nextDocs}`;
}

async function prepareReleaseNote(workingDir, releasedVersion, nextVersion) {
  const releaseNoteDir = path.join(
    workingDir,
    Docs.baseFolder,
    nextVersion.minor(),
    "releases-and-changelog",
    "release-notes",
  );
  const releaseNote = path.join(
    releaseNoteDir,
    `gko-${releasedVersion.minor()}.md`,
  );
  const nextReleaseNote = path.join(
    releaseNoteDir,
    `gko-${nextVersion.minor()}.md`,
  );

  await $`mv ${releaseNote} ${nextReleaseNote}`;
  await $`echo --- > ${nextReleaseNote}`;
  await $`echo 'Gravitee Kubernetes Operator ${nextVersion.minor()} Release Notes.' >> ${nextReleaseNote}`;
  await $`echo --- >> ${nextReleaseNote}`;
  await $`echo >> ${nextReleaseNote}`;
  await $`echo '# GKO ${nextVersion.minor()}' >> ${nextReleaseNote}`;
}

async function prepareChangelog(workingDir, releasedVersion, nextVersion) {
  const changelogDir = path.join(
    workingDir,
    Docs.baseFolder,
    nextVersion.minor(),
    "releases-and-changelog",
    "changelog",
  );
  const changelog = path.join(
    changelogDir,
    `gko-${releasedVersion.minor()}.x.md`,
  );
  const nextChangelog = path.join(
    changelogDir,
    `gko-${nextVersion.minor()}.x.md`,
  );

  await $`mv ${changelog} ${nextChangelog}`;
  await $`echo '# GKO ${nextVersion.minor()}.x' > ${nextChangelog}`;
}

async function prepareNewRelease(workingDir, nextVersion) {
  const newReleaseFile = path.join(
    workingDir,
    Docs.baseFolder,
    nextVersion.minor(),
    "new-release.md",
  );
  await $`echo 'placeholder for the ${nextVersion.minor()} release' > ${newReleaseFile}`;
}

async function prepareAPIReference(workingDir, nextVersion) {
  const referenceFile = path.join(
    workingDir,
    Docs.baseFolder,
    nextVersion.minor(),
    "reference",
    "api-reference.md",
  );
  const refLinkMessage = `The Gravitee Kubernetes Operator (GKO) API reference documentation can be found [in the GKO Github repository](https://github.com/gravitee-io/gravitee-kubernetes-operator/blob/${nextVersion.minor()}.x/docs/api/reference.md).`;
  const crdLinkMessage = `The GKO CRDs can be found [on GitHub](https://github.com/gravitee-io/gravitee-kubernetes-operator/tree/${nextVersion.minor()}.x/helm/gko/crds).`;
  await $`echo '# API Reference' > ${referenceFile}`;
  await $`echo >> ${referenceFile}`;
  await $`echo ${refLinkMessage} >> ${referenceFile}`;
  await $`echo >> ${referenceFile}`;
  await $`echo ${crdLinkMessage} >> ${referenceFile}`;
}

async function prepareSummary(workingDir, releasedVersion, nextVersion) {
  const summaryFile = path.join(
    workingDir,
    Docs.baseFolder,
    nextVersion.minor(),
    "summary.md",
  );

  const summary = await $`cat ${summaryFile}`;

  const withReleaseNote = prepareReleaseNotesSummary(
    `${summary}`,
    releasedVersion,
    nextVersion,
  );
  const withReleaseNoteAndChangelog = prepareChangelogSummary(
    withReleaseNote,
    releasedVersion,
    nextVersion,
  );

  await fs.writeFile(summaryFile, withReleaseNoteAndChangelog);
}

function prepareReleaseNotesSummary(summary, releasedVersion, nextVersion) {
  const startMarker = "<!-- start_release_notes_summary -->";
  const endMarker = "<!-- end_release_notes_summary -->";
  const start = summary.indexOf(startMarker) + startMarker.length;
  const afterStart = start + startMarker.length;
  const end = summary.indexOf(endMarker);
  const nextSummary = summary.substring(afterStart, end);
  const summaryLines = nextSummary
    .split("\n")
    .filter((line) => !isEmptyString(line));
  const releasedNoteLink = `* [GKO ${releasedVersion.minor()}](https://documentation.gravitee.io/gravitee-kubernetes-operator-gko/${releasedVersion.minor()}/releases-and-changelog/release-notes/gko-${releasedVersion.minor()})`;
  const nextNoteLink = `* [GKO ${nextVersion.minor()}](releases-and-changelog/release-notes/gko-${nextVersion.minor()}.md)`;
  summaryLines.shift();
  summaryLines.unshift(releasedNoteLink);
  summaryLines.unshift(nextNoteLink);
  const nextSumarry = summaryLines.join("\n");

  return summary
    .slice(0, start)
    .concat("\n")
    .concat(nextSumarry)
    .concat("\n")
    .concat(summary.slice(end));
}

function prepareChangelogSummary(summary, releasedVersion, nextVersion) {
  const startMarker = "<!-- start_changelogs_summary -->";
  const endMarker = "<!-- end_changelogs_summary -->";
  const start = summary.indexOf(startMarker) + startMarker.length;
  const afterStart = start + startMarker.length;
  const end = summary.indexOf(endMarker);
  const nextSummary = summary.substring(afterStart, end);
  const summaryLines = nextSummary
    .split("\n")
    .filter((line) => !isEmptyString(line));
  const changelogLink = `* [GKO ${releasedVersion.minor()}.x](https://documentation.gravitee.io/gravitee-kubernetes-operator-gko/${releasedVersion.minor()}/releases-and-changelog/changelog/gko-${releasedVersion.minor()}.x)`;
  const nextChangelogLink = `* [GKO ${nextVersion.minor()}.x](releases-and-changelog/changelog/gko-${nextVersion.minor()}.md)`;
  summaryLines.shift();
  summaryLines.unshift(changelogLink);
  summaryLines.unshift(nextChangelogLink);
  const nextSumarry = summaryLines.join("\n");

  return summary
    .slice(0, start)
    .concat("\n")
    .concat(nextSumarry)
    .concat("\n")
    .concat(summary.slice(end));
}

async function commitChanges(workingDir, nextVersion) {
  const branch = `gko-prepare-${nextVersion}-release`;
  cd(workingDir);
  await $`git switch -c ${branch}`;
  await $`git add .`;
  await $`git commit -m "docs(gko): prepare documentation for version ${nextVersion}"`;
  cd(PROJECT_DIR);
}

async function submitChanges(workingDir, nextVersion) {
  const prBranch = `gko-prepare-${nextVersion}-release`;
  const prTitle = `[GKO] Bootstrap documentation for version ${nextVersion}`;
  const prBody = `
  # Prepare documentation for ${nextVersion}

  This pull request bootstraps documentation directories for version ${nextVersion}

  🧐 Please review and merge this pull request to prepare for the next release 🚀
  `;
  cd(workingDir);
  try {
    await $`gh pr close --delete-branch ${prBranch}`;
  } catch (_) {}
  await $`git push --set-upstream origin ${prBranch}`;
  const prURL =
    await $`gh pr create --title ${prTitle} --body ${prBody} --base main --head ${prBranch}`;
  cd(PROJECT_DIR);
}
