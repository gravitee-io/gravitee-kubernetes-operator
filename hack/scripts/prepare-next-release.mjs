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

import { HELM, LOG } from "./lib/index.mjs";
import { Version } from "./lib/version.mjs";
import { rolloutMergify } from "./lib/rollout-mergify.mjs";
import { rolloutTestScheduller } from "./lib/rollout-test-scheduler.mjs";
import { prepareDocs } from "./lib/prepare-docs.mjs";

const releasedVersion = new Version(await HELM.getChartVersion());
const nextPatchVersion = releasedVersion.nextPatch();
const patchCandidateVersion = nextPatchVersion.rc();

LOG.blue(`
    🚧 Switching to ${releasedVersion.branch()} branch`);

await $`git switch ${releasedVersion.branch()}`;

LOG.blue(`
    ⎈ Moving helm chart version from ${releasedVersion} to ${patchCandidateVersion}`);

await HELM.setChartVersion(patchCandidateVersion.toString());

await $`make add-license`;

await $`git add helm/gko/Chart.yaml`;

if (releasedVersion.isPatch()) {
  LOG.blue(`
    🚧 Committing changes to branch ${releasedVersion.branch()}`);

  await $`git commit -m "chore: prepare for ${nextPatchVersion}"`;
  await $`git push -u origin ${releasedVersion.branch()}`;

  process.exit(0);
}

LOG.blue(`
    🚧 Rolling out mergify config`);

await rolloutMergify(releasedVersion.toString());

LOG.blue(`
    🚧 Committing changes to branch ${releasedVersion.branch()}`);

await $`git commit -m "chore: prepare for ${nextPatchVersion}"`;
await $`git push -u origin ${releasedVersion.branch()}`;

const nextMinorVersion = releasedVersion.nextMinor();
const minorCandidateVersion = nextMinorVersion.rc();

LOG.blue(`
    🚧 Switching to master branch`);

await $`git switch master`;

LOG.blue(`
    ⎈ Moving helm chart version from to ${minorCandidateVersion}`);

await HELM.setChartVersion(minorCandidateVersion.toString());

await $`git add helm/gko/Chart.yaml`;

LOG.blue(`
    🚧 Rolling out mergify config`);

await rolloutMergify(releasedVersion.toString());

await $`git add .mergify.yml`;

LOG.blue(`
    🚧 Rolling out test scheduler`);

await rolloutTestScheduller(releasedVersion.toString());

await $`git add .github/workflows/schedule-test.yml `;

await $`git commit -m "chore: prepare for ${nextMinorVersion}"`;
await $`git push -u origin master`;

LOG.blue(`
    🚧 Preparing documentation for ${nextMinorVersion}`);

await prepareDocs(releasedVersion.toString());
