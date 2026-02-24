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

import { AssertionError } from "node:assert";
import type { AssertionReport } from "../../types/match.js";

/**
 * Convert an AssertionReport to void+throw.
 * If the report has failures, throws an AssertionError with detailed info.
 */
export function throwIfFailed(report: AssertionReport): void {
  if (report.pass) return;

  const lines = report.failures.map(
    (f) =>
      `  path:     ${f.path}\n  expected: ${JSON.stringify(f.expected)}\n  actual:   ${JSON.stringify(f.actual)}`,
  );

  throw new AssertionError({
    message: `Assertion failed (${report.failures.length} mismatch${report.failures.length > 1 ? "es" : ""}):\n${lines.join("\n\n")}`,
    actual: report.actual,
    expected: "(partial match)",
    operator: "deepPartialMatch",
  });
}
