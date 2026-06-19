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
  parseVersion,
  compareVersions,
  sinceFromTitle,
  versionSkipReason,
} from "../../src/utils/version/index.js";

describe("parseVersion", () => {
  it("parses major.minor with a default patch", () => {
    expect(parseVersion("4.12")).toEqual([4, 12, 0]);
  });
  it("parses full semver and ignores any pre-release suffix", () => {
    expect(parseVersion("4.12.3")).toEqual([4, 12, 3]);
    expect(parseVersion("4.12.0-milestone.2")).toEqual([4, 12, 0]);
  });
  it("tolerates a leading v", () => {
    expect(parseVersion("v4.11")).toEqual([4, 11, 0]);
  });
  it("returns undefined for non-versions", () => {
    expect(parseVersion("master-latest")).toBeUndefined();
    expect(parseVersion("")).toBeUndefined();
  });
});

describe("compareVersions", () => {
  it("orders by major, then minor, then patch", () => {
    expect(compareVersions("4.11", "4.12")).toBe(-1);
    expect(compareVersions("4.12", "4.11")).toBe(1);
    expect(compareVersions("4.12", "4.12.0")).toBe(0);
    expect(compareVersions("4.12.1", "4.12.0")).toBe(1);
    expect(compareVersions("5.0", "4.99")).toBe(1);
  });
  it("treats unparseable input as equal (fail open)", () => {
    expect(compareVersions("master-latest", "4.12")).toBe(0);
  });
});

describe("sinceFromTitle", () => {
  it("extracts the @since tag value from a title", () => {
    expect(sinceFromTitle("Deploy something @GKO-123 @since-4.12 @regression")).toBe("4.12");
  });
  it("returns undefined when there is no @since tag", () => {
    expect(sinceFromTitle("Deploy something @GKO-123 @regression")).toBeUndefined();
  });
});

describe("versionSkipReason", () => {
  it("runs everything when no cap is set", () => {
    expect(versionSkipReason("x @since-4.12", undefined)).toBeUndefined();
    expect(versionSkipReason("x @since-4.12", null)).toBeUndefined();
  });
  it("runs baseline (untagged) tests under any cap", () => {
    expect(versionSkipReason("plain test @regression", "4.11")).toBeUndefined();
  });
  it("skips a test that needs a newer version than the cap", () => {
    expect(versionSkipReason("new feature @since-4.12", "4.11")).toMatch(/requires APIM 4\.12/);
  });
  it("runs a test at or below the cap", () => {
    expect(versionSkipReason("feature @since-4.12", "4.12")).toBeUndefined();
    expect(versionSkipReason("feature @since-4.11", "4.12")).toBeUndefined();
  });
});
