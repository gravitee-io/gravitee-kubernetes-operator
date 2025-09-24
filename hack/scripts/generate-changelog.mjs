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
  groupBy,
  isEmptyString,
  isNonEmptyString,
  toggleVerbosity,
} from "./lib/index.mjs";
import { Version } from "./lib/version.mjs";

const VERSION = argv.version;
const VERBOSE = argv.verbose;
const JIRA_TOKEN = process.env.JIRA_TOKEN;
const JIRA_BASE = "https://gravitee.atlassian.net/rest/api/3";
const JIRA_PROJECT = "GKO";
const GH_ISSUES_BASE = "https://github.com/gravitee-io/issues/issues";
const DATE_OPTS = { year: "numeric", month: "long", day: "numeric" };
const EOL = "\n";
const TAB = "  ";

toggleVerbosity(VERBOSE);

const JIRA_HEADERS = {
  Authorization: `Basic ${JIRA_TOKEN}`,
  Accept: "application/json",
};

const LOG_COMPONENTS = new Map([
  ["GKO", { label: "GKO", order: 0 }],
  ["APIM", { label: "APIM", order: 1 }],
  ["Others", { label: "Others", order: 2 }],
]);

const LOG_TYPES = new Map([
  ["Public Bug", { label: "Bug fixes", order: 0 }],
  ["Public Improvement", { label: "Improvements", order: 1 }],
  ["Public Security", { label: "Security", order: 2 }],
]);

if (isEmptyString(VERSION)) {
  LOG.red("You must specify a version using the --version flag");
  process.exit(1);
}

const version = new Version(VERSION);
if (version.isPreRelease()) {
  LOG.yellow(
    `No changelog to generate (version ${VERSION} is a pre-release version)`,
  );
  process.exit(0);
}

if (isEmptyString(JIRA_TOKEN)) {
  LOG.red("JIRA_TOKEN must be defined as an environment variable");
  process.exit(1);
}

async function getJiraVersion(versionName) {
  return fetch(`${JIRA_BASE}/project/GKO/versions`, {
    method: "GET",
    headers: JIRA_HEADERS,
  })
    .then((response) => response.json())
    .then((versions) => versions.find(({ name }) => name === versionName));
}

async function getJiraIssues(versionId) {
  const issues = await fetch(`${JIRA_BASE}/search/jql`, {
    method: "POST",
    headers: {
      ...JIRA_HEADERS,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      jql: `project = '${JIRA_PROJECT}' AND fixVersion = '${versionId}'`,
      fields: ["issuetype", "summary", "components", "customfield_10115"],
    }),
  })
    .then((response) => response.json())
    .then((body) => body.issues);

  return issues
    .filter((issue) => LOG_TYPES.has(issue.fields.issuetype.name))
    .map((issue) => ({
      key: issue.key,
      githubIssue: issue.fields.customfield_10115,
      summary: issue.fields.summary,
      type: issue.fields.issuetype.name,
      component: findComponent(issue),
    }));
}

function findComponent(issue) {
  return (
    [...LOG_COMPONENTS.keys()].find((component) =>
      issue.fields.components.some((c) => c.name === component),
    ) || "Others"
  );
}

function groupByTypeAndComponent(issues) {
  return groupByType(issues).map(([type, issues]) => [
    type,
    groupByComponent(issues),
  ]);
}

function groupByType(issues) {
  return [...groupBy(issues, (issue) => issue.type).entries()].sort(
    ([t1], [t2]) => LOG_TYPES.get(t1).order - LOG_TYPES.get(t2).order,
  );
}

function groupByComponent(issues) {
  return [...groupBy(issues, (issue) => issue.component).entries()].sort(
    ([c1], [c2]) => LOG_COMPONENTS.get(c1).order - LOG_COMPONENTS.get(c2).order,
  );
}

function buildLogs([type, componentToIssues]) {
  return `
<details>
<summary>${LOG_TYPES.get(type).label}</summary>

${componentToIssues.map(buildComponentLogs).join(EOL)}
</details>
`;
}

function buildComponentLogs([component, issues]) {
  return `${TAB}**${LOG_COMPONENTS.get(component).label}**

${issues.map(buildSummary).join(EOL)}
`;
}

function buildSummary(issue) {
  return isNonEmptyString(issue.githubIssue)
    ? `${TAB}* ${issue.summary} [#${issue.githubIssue}](${GH_ISSUES_BASE}/${issue.githubIssue})`
    : `${TAB}* ${issue.summary}`;
}

const releaseDate = new Date().toLocaleDateString("en-US", DATE_OPTS);
const releaseChangelogHeader = `## Gravitee Kubernetes Operator ${VERSION} - ${releaseDate}`;
const noChangeMessage = `
${releaseChangelogHeader}

There is nothing new in version ${VERSION}.

> This version was generated to keep the kubernetes operator in sync with other gravitee products.

`;

const noChangelogMessage = `
${VERSION} was not created in Jira. Assuming no changelog should be generated.
`;

const jiraVersion = await getJiraVersion(VERSION);

if (!jiraVersion) {
  echo(noChangelogMessage);
  process.exit(0);
}

const jiraIssues = await getJiraIssues(jiraVersion.id);

if (jiraIssues.length === 0) {
  echo(noChangeMessage);
} else {
  echo(`${releaseChangelogHeader}
    ${groupByTypeAndComponent(jiraIssues).map(buildLogs).join(EOL)}
`);
}
