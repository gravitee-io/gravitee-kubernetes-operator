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
import type { PollOptions } from "../../types/match.js";

/**
 * Retry an async assertion until it passes or times out.
 * Works with any void+throw assertion function.
 *
 * @example
 * await poll(() => apim.assertApiStarted(apiId), { timeoutMs: 15_000 });
 */
export async function poll(
  fn: () => Promise<void>,
  options: PollOptions = {},
): Promise<void> {
  const { timeoutMs = 30_000, intervalMs = 1_000, description } = options;
  const deadline = Date.now() + timeoutMs;
  let lastError: Error | undefined;

  while (Date.now() < deadline) {
    try {
      await fn();
      return;
    } catch (err) {
      lastError = err instanceof Error ? err : new Error(String(err));
      await new Promise((r) => setTimeout(r, intervalMs));
    }
  }

  const msg = description
    ? `Timed out after ${timeoutMs}ms waiting for: ${description}`
    : `Timed out after ${timeoutMs}ms`;

  // If the last attempt threw an AssertionError, re-throw it directly so
  // callers see the structured operator/actual/expected fields, not a
  // generic timeout wrapper.
  if (lastError instanceof AssertionError) {
    throw lastError;
  }

  const timeoutErr = new Error(msg);
  timeoutErr.cause = lastError;
  throw timeoutErr;
}
