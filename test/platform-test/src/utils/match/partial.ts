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

import type { AssertionFailure, AssertionReport } from "../../types/match.js";

/**
 * Deep partial match: checks that every field in `expected` exists in `actual`
 * with the same value. Fields NOT in `expected` are ignored.
 */
export function deepPartialMatch(
  actual: unknown,
  expected: unknown,
  path: string = "$",
): AssertionReport {
  const failures: AssertionFailure[] = [];
  collect(actual, expected, path, failures);
  return { pass: failures.length === 0, failures, actual };
}

function collect(
  actual: unknown,
  expected: unknown,
  path: string,
  failures: AssertionFailure[],
): void {
  // Primitives and null: strict equality
  if (expected === null || expected === undefined || typeof expected !== "object") {
    if (!Object.is(actual, expected)) {
      failures.push({
        path,
        message: `Expected ${JSON.stringify(expected)} but got ${JSON.stringify(actual)}`,
        expected,
        actual,
      });
    }
    return;
  }

  // Expected is object but actual isn't
  if (actual === null || actual === undefined || typeof actual !== "object") {
    failures.push({
      path,
      message: `Expected an object but got ${JSON.stringify(actual)}`,
      expected,
      actual,
    });
    return;
  }

  // Array matching
  if (Array.isArray(expected)) {
    if (!Array.isArray(actual)) {
      failures.push({
        path,
        message: `Expected an array but got ${typeof actual}`,
        expected,
        actual,
      });
      return;
    }
    for (let i = 0; i < expected.length; i++) {
      if (i >= actual.length) {
        failures.push({
          path: `${path}[${i}]`,
          message: `Array index ${i} missing (array has ${actual.length} elements)`,
          expected: expected[i],
          actual: undefined,
        });
        continue;
      }
      collect(actual[i], expected[i], `${path}[${i}]`, failures);
    }
    return;
  }

  // Object matching (partial)
  const expectedObj = expected as Record<string, unknown>;
  const actualObj = actual as Record<string, unknown>;

  for (const key of Object.keys(expectedObj)) {
    const childPath = path === "$" ? `$.${key}` : `${path}.${key}`;
    if (!(key in actualObj)) {
      failures.push({
        path: childPath,
        message: `Property "${key}" missing from actual object`,
        expected: expectedObj[key],
        actual: undefined,
      });
      continue;
    }
    collect(actualObj[key], expectedObj[key], childPath, failures);
  }
}
