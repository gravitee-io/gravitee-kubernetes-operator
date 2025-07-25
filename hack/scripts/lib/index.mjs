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

const { log } = console;

const newLoggerFn =
  (color) =>
  (...args) =>
    log(color(...args));

const $Quote = $.quote;
const $NoQuote = (unescaped) => unescaped;

// Record time taken by a fn to execute
export async function time(fn) {
  const start = Date.now();
  await fn();
  const end = Date.now();
  green(`Done in ${(end - start) / 1000}s`);
}

// Toggle zx verbosity (print commands before executing them)
export function toggleVerbosity(verbose = false) {
  $.verbose = verbose;
}

// Disable quote escaping for zx
export function setNoQuoteEscape() {
  $.quote = $NoQuote;
}

// Enable quote escaping for zx
export function setQuoteEscape() {
  $.quote = $Quote;
}

// Utility function to convert camelCase to kebab-case
// Used to generate files from resources names
// e.g ApiDefinition -> api-definition
export function toKebabCase(camelCaseString) {
  return camelCaseString.replace(/([a-z])([A-Z])/g, "$1-$2").toLowerCase();
}

// Can be used to check if a flag as been passed as a string to a script
export function isNonEmptyString(str) {
  return String(str) === str && str.trim().length > 0;
}

export function isEmptyString(str) {
  return !isNonEmptyString(str);
}

// see https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Map/groupBy
// should be replaced by the native function when available
export function groupBy(arr, func) {
  const groups = new Map();
  for (const item of arr) {
    const key = func(item);
    if (groups.has(key)) {
      groups.get(key).push(item);
    } else {
      groups.set(func(item), [item]);
    }
  }
  return groups;
}

// Color loggers
const green = newLoggerFn(chalk.green);
const blue = newLoggerFn(chalk.blue);
const magenta = newLoggerFn(chalk.magenta);
const yellow = newLoggerFn(chalk.yellow);
const red = newLoggerFn(chalk.red);

export const LOG = Object.seal({ green, blue, magenta, yellow, red, log });

// Path to the local helm chart directory
export const PROJECT_DIR = path.join(__dirname, "..", "..");

// Path to the local helm chart directory
const chartDir = path.join(PROJECT_DIR, "helm", "gko");

// Path to the helm crds directory. This resources are not templated.
const crdDir = path.join(chartDir, "crds", "gravitee.io");

// The gravitee.io official helm charts repository
const chartsRepo = "gravitee-io/helm-charts";

// The release branch for gravitee.io official helm charts repository
const releaseBranch = "gh-pages";

const releaseVersionAnnotation = "gravitee.io/operator.version";

export const HELM = {
  chartDir,
  crdDir,
  chartsRepo,
  releaseBranch,
  releaseVersionAnnotation,
  getChartVersion,
  setChartVersion,
  annotateCRDs,
};

const docsRepo = "gravitee-io/gravitee-platform-docs";
const docsRepoURL = `https://github.com/gravitee-io/${docsRepo}`;

export const Docs = {
  repo: docsRepo,
  repoURL: docsRepoURL,
  baseFolder: "docs/gko",
  changelogFolder: "releases-and-changelog/changelog",
};

async function getChartVersion() {
  const chartFile = await fs.readFile(`${chartDir}/Chart.yaml`, "utf8");
  const chartYaml = await YAML.parse(chartFile);
  return chartYaml.version;
}

async function setChartVersion(version) {
  const chartFile = await fs.readFile(`${chartDir}/Chart.yaml`, "utf8");
  const chartYaml = await YAML.parse(chartFile);
  chartYaml.version = version;
  chartYaml.appVersion = version;
  await fs.writeFile(`${chartDir}/Chart.yaml`, YAML.stringify(chartYaml));
}

async function annotateCRDs(version) {
  const fileNames = await fs.readdir(crdDir);
  for (const fileName of fileNames) {
    await annotateCRD(fileName, version);
  }
}

async function annotateCRD(fileName, version) {
  const crdFile = await fs.readFile(`${crdDir}/${fileName}`, "utf8");
  const crdYaml = await YAML.parse(crdFile);
  crdYaml.metadata.annotations[releaseVersionAnnotation] = version;
  await fs.writeFile(`${crdDir}/${fileName}`, YAML.stringify(crdYaml));
}
