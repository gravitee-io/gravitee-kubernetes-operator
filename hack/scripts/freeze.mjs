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

import { HELM, LOG, PROJECT_DIR, isEmptyString } from "./lib/index.mjs";
import { Version } from "./lib/version.mjs";
import { rolloutMergify } from "./lib/rollout-mergify.mjs";
import { rolloutTestScheduler } from "./lib/rollout-test-scheduler.mjs";
import { prepareDocs } from "./lib/prepare-docs.mjs";

const VERSION = argv.version;

if (isEmptyString(VERSION)) {
  LOG.red(
    "You must specify a version using the --version flag (e.g. --version 4.11.0)",
  );
  process.exit(1);
}

const frozenVersion = new Version(VERSION);

if (frozenVersion.isPatch()) {
  LOG.red(
    "Freeze is only supported for minor releases (patch digit must be 0)",
  );
  process.exit(1);
}

const releaseBranch = frozenVersion.branch();
const nextMinorVersion = frozenVersion.nextMinor();
const nextMinorRC = nextMinorVersion.rc();

LOG.blue(`
    ❄️  Freezing version ${frozenVersion}
        Release branch: ${releaseBranch}
        Next minor:     ${nextMinorRC}`);

// -- Release branch ----------------------------------------------------------

LOG.blue(`
    🚧 Switching to master and creating release branch ${releaseBranch}`);

await $`git switch master`;
await $`git switch -c ${releaseBranch}`;

const branchLabel = `apply-on-${releaseBranch.replaceAll(".", "-")}`;

LOG.blue(`
    🏷️  Creating GitHub label ${branchLabel}`);

await $`gh label create ${branchLabel} --description ${"Backport to " + releaseBranch} --force`;

LOG.blue(`
    🚧 Rolling out mergify config on ${releaseBranch}`);

await rolloutMergify(frozenVersion.toString());

LOG.blue(`
    📝 Updating hack/stable.yaml → branch: ${releaseBranch}`);

await updateStableYaml(releaseBranch);

LOG.blue(`
    📝 Updating impl.yaml → version: ${frozenVersion}`);

await updateImplYaml(frozenVersion.toString());

LOG.blue(`
    📝 Updating hack/apim.yaml → image: ${releaseBranch}-latest, chart: ${frozenVersion.minor()}.*`);

await updateApimYaml({
  imageVersion: `${releaseBranch}-latest`,
  chartVersion: `${frozenVersion.minor()}.*`,
});

await $`make generate manifests reference helm-reference`;
await $`make add-license`;

await $`git add .`;
await $`git commit -m "chore: freeze ${frozenVersion}"`;
await $`git push -u origin ${releaseBranch}`;

// -- Master branch -----------------------------------------------------------

LOG.blue(`
    🚧 Switching to master`);

await $`git switch master`;

LOG.blue(`
    ⎈ Moving helm chart version to ${nextMinorRC}`);

await HELM.setChartVersion(nextMinorRC.toString());

await $`git add helm/gko/Chart.yaml`;

LOG.blue(`
    📝 Updating hack/apim.yaml → chart version: ${nextMinorVersion.minor()}.*`);

await updateApimYaml({ chartVersion: `${nextMinorVersion.minor()}.*` });

await $`git add hack/apim.yaml`;

LOG.blue(`
    🚧 Rolling out mergify config on master`);

await rolloutMergify(frozenVersion.toString());

await $`git add .mergify.yml`;

LOG.blue(`
    🚧 Rolling out test scheduler on master`);

await rolloutTestScheduler(frozenVersion.toString());

await $`git add .github/workflows/schedule-test.yml`;

await $`make generate manifests reference helm-reference`;
await $`make add-license`;

await $`git add .`;
await $`git commit -m "chore: prepare for ${nextMinorVersion}"`;
await $`git push -u origin master`;

// -- Documentation -----------------------------------------------------------

LOG.blue(`
    📚 Preparing documentation for ${nextMinorVersion}`);

await prepareDocs(frozenVersion.toString());

LOG.green(`
    ✅ Freeze complete for ${frozenVersion}`);

// -- Helpers -----------------------------------------------------------------

async function updateStableYaml(branch) {
  const filePath = path.join(PROJECT_DIR, "hack", "stable.yaml");
  const content = await fs.readFile(filePath, "utf8");
  const yaml = YAML.parse(content);
  yaml.branch = branch;
  await fs.writeFile(filePath, YAML.stringify(yaml));
}

async function updateApimYaml({ imageVersion, chartVersion }) {
  const filePath = path.join(PROJECT_DIR, "hack", "apim.yaml");
  const content = await fs.readFile(filePath, "utf8");
  const yaml = YAML.parse(content);
  if (imageVersion) {
    yaml.apim.image.version = imageVersion;
  }
  if (chartVersion) {
    yaml.apim.chart.version = chartVersion;
  }
  await fs.writeFile(filePath, YAML.stringify(yaml));
}

async function updateImplYaml(version) {
  const filePath = path.join(
    PROJECT_DIR,
    "test",
    "conformance",
    "kubernetes.io",
    "gateway-api",
    "impl",
    "impl.yaml",
  );
  const content = await fs.readFile(filePath, "utf8");
  const yaml = YAML.parse(content);
  yaml.version = version;
  await fs.writeFile(filePath, YAML.stringify(yaml));
}
