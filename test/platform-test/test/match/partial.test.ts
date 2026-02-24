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
import { deepPartialMatch } from "../../src/utils/match/partial.js";

describe("deepPartialMatch", () => {
  it("passes when actual matches expected primitives", () => {
    const result = deepPartialMatch({ name: "foo", age: 42 }, { name: "foo" });
    expect(result.pass).toBe(true);
    expect(result.failures).toHaveLength(0);
  });

  it("fails when a primitive field mismatches", () => {
    const result = deepPartialMatch({ name: "foo" }, { name: "bar" });
    expect(result.pass).toBe(false);
    expect(result.failures).toHaveLength(1);
    expect(result.failures[0].path).toBe("$.name");
    expect(result.failures[0].expected).toBe("bar");
    expect(result.failures[0].actual).toBe("foo");
  });

  it("passes with nested partial match", () => {
    const actual = { a: { b: { c: 1, d: 2 }, e: 3 } };
    const expected = { a: { b: { c: 1 } } };
    expect(deepPartialMatch(actual, expected).pass).toBe(true);
  });

  it("fails with nested mismatch", () => {
    const actual = { a: { b: { c: 1 } } };
    const expected = { a: { b: { c: 99 } } };
    const result = deepPartialMatch(actual, expected);
    expect(result.pass).toBe(false);
    expect(result.failures[0].path).toBe("$.a.b.c");
  });

  it("matches arrays element-by-element", () => {
    const actual = { items: ["a", "b", "c"] };
    const expected = { items: ["a", "b"] };
    expect(deepPartialMatch(actual, expected).pass).toBe(true);
  });

  it("fails when array element mismatches", () => {
    const result = deepPartialMatch({ items: ["a", "x"] }, { items: ["a", "b"] });
    expect(result.pass).toBe(false);
    expect(result.failures[0].path).toBe("$.items[1]");
  });

  it("fails when expected array is longer than actual", () => {
    const result = deepPartialMatch({ items: ["a"] }, { items: ["a", "b"] });
    expect(result.pass).toBe(false);
    expect(result.failures[0].path).toBe("$.items[1]");
  });

  it("fails when a property is missing from actual", () => {
    const result = deepPartialMatch({ a: 1 }, { a: 1, b: 2 });
    expect(result.pass).toBe(false);
    expect(result.failures[0].path).toBe("$.b");
  });

  it("matches null explicitly", () => {
    expect(deepPartialMatch({ a: null }, { a: null }).pass).toBe(true);
    expect(deepPartialMatch({ a: "x" }, { a: null }).pass).toBe(false);
  });

  it("fails when expected is object but actual is primitive", () => {
    const result = deepPartialMatch({ a: "string" }, { a: { nested: true } });
    expect(result.pass).toBe(false);
  });

  it("collects multiple failures", () => {
    const actual = { name: "wrong", state: "STOPPED", version: "1.0" };
    const expected = { name: "right", state: "STARTED" };
    const result = deepPartialMatch(actual, expected);
    expect(result.pass).toBe(false);
    expect(result.failures).toHaveLength(2);
  });

  it("handles nested objects in arrays", () => {
    const actual = { listeners: [{ type: "HTTP", entrypoints: [{ type: "proxy", host: "a.com" }] }] };
    const expected = { listeners: [{ type: "HTTP", entrypoints: [{ host: "a.com" }] }] };
    expect(deepPartialMatch(actual, expected).pass).toBe(true);
  });
});
