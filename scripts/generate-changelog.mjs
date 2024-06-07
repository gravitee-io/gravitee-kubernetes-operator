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

import { LOG, isEmptyString, toggleVerbosity } from "./lib/index.mjs";

const VERSION = argv.version;
const VERBOSE = argv.verbose;
const JIRA_TOKEN = process.env.JIRA_TOKEN;
const JIRA_BASE = "https://gravitee.atlassian.net/rest/api/3";
const JIRA_PROJECT = "GKO";
const DATE_OPTS = { year: "numeric", month: "long", day: "numeric" };

const JIRA_HEADERS = {
  Authorization: `Basic ${JIRA_TOKEN}`,
  Accept: "application/json",
};

const CHANGELOG_ISSUES = ["Story"];

toggleVerbosity(VERBOSE);

if (isEmptyString(VERSION)) {
  LOG.red("You must specify a version using the --version flag");
  await $`exit 1`;
}

if (isEmptyString(JIRA_TOKEN)) {
  LOG.red("JIRA_TOKEN must be defined as an environment variable");
  await $`exit 1`;
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
  const query = `jql=project=${JIRA_PROJECT} AND fixVersion=${versionId}`;

  const issues = await fetch(`${JIRA_BASE}/search?${query}`, {
    method: "GET",
    headers: JIRA_HEADERS,
  })
    .then((response) => response.json())
    .then((body) => body.issues);

  return issues
    .filter((issue) => CHANGELOG_ISSUES.includes(issue.fields.issuetype.name))
    .map((issue) => ({
      key: issue.key,
      githubIssue: issue.fields.customfield_10115,
      summary: issue.fields.summary,
      components: issue.fields.components,
      type: issue.fields.issuetype.name,
    }));
}

const jiraVersion = await getJiraVersion(VERSION);

const jiraIssues = await getJiraIssues(jiraVersion.id);

let changelog = `
## GKO ${VERSION} - ${new Date().toLocaleDateString("en-US", DATE_OPTS)}
`;

changelog += `
<details>
  <summary><b>What's new ?</summary>
`;

for (const issue of jiraIssues) {
  changelog += `
    * ${issue.summary}
    `;
}

changelog += `
</details>
`;

echo(changelog);
