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
import { rolloutMergify } from "./rollout-mergify.mjs";
import { rolloutTestScheduller } from "./rollout-test-scheduler.mjs";

const releasedVersion = new Version(await HELM.getChartVersion());
const candidateVersion = releasedVersion.nextPatch().rc().toString();

LOG.blue(`
    ⎈ Moving helm chart version from ${releasedVersion} to ${candidateVersion}`);

HELM.setChartVersion(candidateVersion.toString());

if (releasedVersion.isPatch()) {
  process.exit(0);
}

LOG.blue(`
    🚧 Rolling out mergify config`);

await rolloutMergify(releasedVersion.toString());

LOG.blue(`
    🚧 Rolling out test scheduler`);

await rolloutTestScheduller(releasedVersion.toString());
