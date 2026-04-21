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
import { AssertionError } from "node:assert";
import { throwIfFailed } from "../../src/utils/match/result.js";

describe("throwIfFailed", () => {
  it("returns void when report passes", () => {
    const result = throwIfFailed({ pass: true, failures: [], actual: {} });
    expect(result).toBeUndefined();
  });

  it("throws AssertionError when report has failures", () => {
    expect(() =>
      throwIfFailed({
        pass: false,
        failures: [
          { jsonPath: "$.name", message: "mismatch", expected: "foo", actual: "bar" },
        ],
        actual: { name: "bar" },
      }),
    ).toThrow(AssertionError);
  });

  it("includes all failure paths in the error message", () => {
    try {
      throwIfFailed({
        pass: false,
        failures: [
          { jsonPath: "$.name", message: "m1", expected: "a", actual: "b" },
          { jsonPath: "$.state", message: "m2", expected: "STARTED", actual: "STOPPED" },
        ],
        actual: { name: "b", state: "STOPPED" },
      });
      expect.unreachable("should have thrown");
    } catch (err) {
      expect(err).toBeInstanceOf(AssertionError);
      const ae = err as AssertionError;
      expect(ae.message).toContain("2 mismatches");
      expect(ae.message).toContain("$.name");
      expect(ae.message).toContain("$.state");
      expect(ae.operator).toBe("deepPartialMatch");
    }
  });
});
