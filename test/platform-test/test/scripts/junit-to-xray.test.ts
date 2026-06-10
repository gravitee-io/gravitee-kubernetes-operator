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

import { describe, it, expect } from "vitest";
import {
  parseTestcases,
  toTests,
  buildPayload,
  tally,
  defaultSummary,
  defaultDescription,
} from "../../scripts/xray/junit-to-xray.mjs";

const tc = (name: string, body = "") =>
  `<testcase name="${name}" classname="suite.test.ts" time="1.0">${body}</testcase>`;

describe("junit-to-xray transform", () => {
  it("maps a clean testcase to PASSED keyed by its @GKO id", () => {
    const xml = tc("does a thing @GKO-100 @regression");
    expect(toTests(parseTestcases(xml))).toEqual([{ testKey: "GKO-100", status: "PASSED" }]);
  });

  it("maps <failure> and <error> to FAILED", () => {
    const xml =
      tc("fails @GKO-1", "<failure>boom</failure>") + tc("errors @GKO-2", "<error>nope</error>");
    expect(toTests(parseTestcases(xml))).toEqual([
      { testKey: "GKO-1", status: "FAILED" },
      { testKey: "GKO-2", status: "FAILED" },
    ]);
  });

  it("omits skipped / fixme'd cases so a Test's prior result is not clobbered", () => {
    const xml =
      tc("ran @GKO-1") +
      tc("skipped @GKO-2", "<properties><property name=\"fixme\" value=\"\"></property></properties><skipped></skipped>");
    expect(toTests(parseTestcases(xml))).toEqual([{ testKey: "GKO-1", status: "PASSED" }]);
  });

  it("drops testcases without a @GKO key", () => {
    const xml = tc("untagged test") + tc("tagged @GKO-9");
    expect(toTests(parseTestcases(xml))).toEqual([{ testKey: "GKO-9", status: "PASSED" }]);
  });

  it("emits one entry per key when a title carries several", () => {
    const xml = tc("covers two @GKO-1 @GKO-2");
    expect(toTests(parseTestcases(xml))).toEqual([
      { testKey: "GKO-1", status: "PASSED" },
      { testKey: "GKO-2", status: "PASSED" },
    ]);
  });

  it("does not mistake classname for the test name and survives XML entities", () => {
    const xml = tc("A &amp; B &gt; C @GKO-42");
    expect(toTests(parseTestcases(xml))).toEqual([{ testKey: "GKO-42", status: "PASSED" }]);
  });

  it("last status wins when the same key appears twice in a run", () => {
    const xml = tc("first @GKO-7") + tc("retry @GKO-7", "<failure>flaky</failure>");
    expect(toTests(parseTestcases(xml))).toEqual([{ testKey: "GKO-7", status: "FAILED" }]);
  });
});

describe("tally", () => {
  it("counts passed / failed / skipped from parsed cases", () => {
    const xml =
      tc("a @GKO-1") +
      tc("b @GKO-2", "<failure>x</failure>") +
      tc("c @GKO-3", "<skipped/>") +
      tc("untagged");
    expect(tally(parseTestcases(xml))).toEqual({ passed: 2, failed: 1, skipped: 1 });
  });
});

describe("defaultSummary", () => {
  it("includes branch, totals, and CircleCI build number", () => {
    expect(
      defaultSummary({ passed: 355, failed: 0, skipped: 0 }, { branch: "master", buildNum: "85230" }),
    ).toBe("GKO Playwright e2e on master — 355 passed, 0 failed (CircleCI #85230)");
  });

  it("surfaces skipped count when non-zero", () => {
    expect(
      defaultSummary({ passed: 12, failed: 2, skipped: 3 }, { branch: "feature/x", buildNum: "1" }),
    ).toBe("GKO Playwright e2e on feature/x — 12 passed, 2 failed, 3 skipped (CircleCI #1)");
  });

  it("falls back to '(local)' when no build number is supplied", () => {
    expect(defaultSummary({ passed: 1, failed: 0, skipped: 0 }, { branch: "master" })).toBe(
      "GKO Playwright e2e on master — 1 passed, 0 failed (local)",
    );
  });

  it("omits the branch segment when no branch is supplied", () => {
    expect(defaultSummary({ passed: 1, failed: 0, skipped: 0 }, { buildNum: "1" })).toBe(
      "GKO Playwright e2e — 1 passed, 0 failed (CircleCI #1)",
    );
  });
});

describe("defaultDescription", () => {
  it("returns just the CircleCI import line when only buildUrl is supplied", () => {
    expect(defaultDescription({ buildUrl: "https://x.test/1" })).toBe(
      "Automated import from CircleCI: https://x.test/1",
    );
  });

  it("prepends 'Commit: ...' when a commitSubject is supplied", () => {
    expect(
      defaultDescription({
        buildUrl: "https://x.test/1",
        commitSubject: "test: readable Test Execution title (GKO-2906)",
      }),
    ).toBe(
      "Commit: test: readable Test Execution title (GKO-2906)\nAutomated import from CircleCI: https://x.test/1",
    );
  });

  it("falls back to a generic line when no buildUrl is supplied", () => {
    expect(defaultDescription({})).toBe("Automated import of Playwright e2e results.");
  });
});

describe("buildPayload", () => {
  it("wraps tests with auto-templated summary/description when none supplied", () => {
    const payload = buildPayload(tc("t @GKO-1"), {
      branch: "master",
      buildNum: "42",
      buildUrl: "https://x.test/42",
    });
    expect(payload.tests).toEqual([{ testKey: "GKO-1", status: "PASSED" }]);
    expect(payload.info.summary).toBe(
      "GKO Playwright e2e on master — 1 passed, 0 failed (CircleCI #42)",
    );
    expect(payload.info.description).toBe("Automated import from CircleCI: https://x.test/42");
    expect(payload.info).not.toHaveProperty("testPlanKey");
  });

  it("includes the commit subject in the auto-templated description when supplied", () => {
    const payload = buildPayload(tc("t @GKO-1"), {
      branch: "master",
      buildNum: "42",
      buildUrl: "https://x.test/42",
      commitSubject: "fix: handle empty result file (GKO-2999)",
    });
    expect(payload.info.description).toBe(
      "Commit: fix: handle empty result file (GKO-2999)\nAutomated import from CircleCI: https://x.test/42",
    );
  });

  it("uses provided summary/description and adds testPlanKey only when set", () => {
    const payload = buildPayload(tc("t @GKO-1"), {
      summary: "run 42",
      description: "from CI",
      testPlanKey: "GKO-500",
    });
    expect(payload.info).toEqual({
      summary: "run 42",
      description: "from CI",
      testPlanKey: "GKO-500",
    });
  });
});
