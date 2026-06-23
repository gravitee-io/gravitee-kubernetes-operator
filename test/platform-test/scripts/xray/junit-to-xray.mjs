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

// Transform a Playwright JUnit XML report into the Xray Cloud JSON results
// format, keyed by the @GKO-NNNN Jira Test issue id embedded in each test
// title. Unlike the JUnit importer (which matches Tests by classname+name and
// would create duplicate Test issues), this lets Xray attach each result to the
// existing Test by its key, updating that Test's "Test Coverage" panel.
//
// Reads the JUnit XML from argv[2] (or stdin), writes Xray JSON to stdout:
//   { info: { summary, description, [testPlanKey] }, tests: [ { testKey, status } ] }
//
// Status mapping: <failure>/<error> -> FAILED, otherwise PASSED. Skipped /
// test.fixme'd cases are omitted entirely so a temporarily-disabled test does
// not overwrite a Test's previously recorded result with TODO.
//
// Env (all optional):
//   XRAY_SUMMARY         explicit Test Execution summary; overrides the auto-templated one
//   XRAY_DESCRIPTION     explicit description; overrides the auto-templated one
//   XRAY_TEST_PLAN_KEY   existing Test Plan key to link the Execution to
//   XRAY_BRANCH          branch name woven into the default summary (e.g. "master")
//   XRAY_BUILD_NUM       CircleCI build number woven into the default summary
//   XRAY_BUILD_URL       CircleCI build URL woven into the default description
//   XRAY_COMMIT_SUBJECT  git commit subject; prepended to the default description as "Commit: <subject>"

import { readFileSync, realpathSync } from "node:fs";
import { fileURLToPath } from "node:url";

const KEY_RE = /@(GKO-\d+)/g;

export function decodeEntities(s) {
  return s
    .replace(/&lt;/g, "<")
    .replace(/&gt;/g, ">")
    .replace(/&quot;/g, '"')
    .replace(/&apos;/g, "'")
    .replace(/&#39;/g, "'")
    .replace(/&amp;/g, "&");
}

export function parseTestcases(xml) {
  const cases = [];
  const re = /<testcase\b([^>]*)>([\s\S]*?)<\/testcase>/g;
  let m;
  while ((m = re.exec(xml)) !== null) {
    const attrs = m[1];
    const inner = m[2];
    const nameMatch = /(?:^|\s)name="([^"]*)"/.exec(attrs);
    if (!nameMatch) continue;
    const name = decodeEntities(nameMatch[1]);
    const failed = /<(failure|error)\b/.test(inner);
    const skipped = /<skipped\b/.test(inner);
    cases.push({ name, failed, skipped });
  }
  return cases;
}

export function toTests(cases) {
  // One Test issue maps to one result. If a title carries several @GKO keys,
  // record the same status against each. Last write wins per key within a run.
  const byKey = new Map();
  for (const c of cases) {
    if (c.skipped) continue; // not executed -> leave the Test's prior result intact
    const keys = [...c.name.matchAll(KEY_RE)].map((k) => k[1]);
    const status = c.failed ? "FAILED" : "PASSED";
    for (const key of keys) byKey.set(key, status);
  }
  return [...byKey].map(([testKey, status]) => ({ testKey, status }));
}

export function tally(cases) {
  // Count cases from the raw parsed list so the summary reflects what
  // *actually ran* (including untagged ones), independent of testKey dedup.
  let passed = 0,
    failed = 0,
    skipped = 0;
  for (const c of cases) {
    if (c.skipped) skipped++;
    else if (c.failed) failed++;
    else passed++;
  }
  return { passed, failed, skipped };
}

export function defaultSummary({ passed, failed, skipped }, { branch, buildNum } = {}) {
  // Format the Test Execution title so it's readable at a glance in the
  // Jira/Xray list view — branch and outcome up front, build number as a
  // parenthetical for click-through traceability.
  const where = branch ? ` on ${branch}` : "";
  const skippedPart = skipped > 0 ? `, ${skipped} skipped` : "";
  const buildPart = buildNum ? ` (CircleCI #${buildNum})` : " (local)";
  return `GKO Playwright e2e${where} — ${passed} passed, ${failed} failed${skippedPart}${buildPart}`;
}

export function defaultDescription({ buildUrl, commitSubject } = {}) {
  // Keep the description deliberately small — the Xray-generated "Tests" panel
  // already shows totals + per-test rundown inside the issue body, so anything
  // we add here that duplicates it is noise. Two pieces of context that aren't
  // visible anywhere else on the page:
  //   1. The CircleCI build URL (one click to logs / artifacts).
  //   2. The commit subject (what change was being tested, in plain English).
  // Falls back to the original merged-master one-liner when either is missing.
  const base = buildUrl
    ? `Automated import from CircleCI: ${buildUrl}`
    : "Automated import of Playwright e2e results.";
  return commitSubject ? `Commit: ${commitSubject}\n${base}` : base;
}

export function buildPayload(
  xml,
  { summary, description, testPlanKey, branch, buildNum, buildUrl, commitSubject } = {},
) {
  const cases = parseTestcases(xml);
  const tests = toTests(cases);
  const totals = tally(cases);
  const info = {
    summary: summary || defaultSummary(totals, { branch, buildNum }),
    description: description || defaultDescription({ buildUrl, commitSubject }),
  };
  if (testPlanKey) info.testPlanKey = testPlanKey;
  return { info, tests };
}

function main() {
  const path = process.argv[2];
  const xml = path ? readFileSync(path, "utf8") : readFileSync(0, "utf8");
  const payload = buildPayload(xml, {
    summary: process.env.XRAY_SUMMARY,
    description: process.env.XRAY_DESCRIPTION,
    testPlanKey: process.env.XRAY_TEST_PLAN_KEY,
    branch: process.env.XRAY_BRANCH,
    buildNum: process.env.XRAY_BUILD_NUM,
    buildUrl: process.env.XRAY_BUILD_URL,
    commitSubject: process.env.XRAY_COMMIT_SUBJECT,
  });
  process.stdout.write(JSON.stringify(payload, null, 2) + "\n");
}

// Run only when executed directly, not when imported (e.g. by unit tests).
const invokedDirectly =
  process.argv[1] &&
  realpathSync(process.argv[1]) === realpathSync(fileURLToPath(import.meta.url));
if (invokedDirectly) main();
