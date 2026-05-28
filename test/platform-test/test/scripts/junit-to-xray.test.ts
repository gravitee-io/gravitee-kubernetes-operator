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
import { parseTestcases, toTests, buildPayload } from "../../scripts/xray/junit-to-xray.mjs";

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

describe("buildPayload", () => {
  it("wraps tests with default info when none supplied", () => {
    const payload = buildPayload(tc("t @GKO-1"));
    expect(payload.tests).toEqual([{ testKey: "GKO-1", status: "PASSED" }]);
    expect(payload.info.summary).toBeTruthy();
    expect(payload.info).not.toHaveProperty("testPlanKey");
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
