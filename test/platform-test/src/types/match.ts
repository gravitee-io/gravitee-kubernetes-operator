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

/**
 * Utility: recursively makes all properties optional.
 */
export type DeepPartial<T> = T extends object
  ? { [P in keyof T]?: DeepPartial<T[P]> }
  : T;

/**
 * A single field-level assertion failure.
 */
export interface AssertionFailure {
  jsonPath: string;
  message: string;
  expected: unknown;
  actual: unknown;
}

/**
 * Aggregated result of a partial match assertion.
 */
export interface AssertionReport {
  pass: boolean;
  failures: AssertionFailure[];
  actual: unknown;
}

export interface PollOptions {
  /** Maximum time to wait in ms (default: 30_000) */
  timeoutMs?: number;
  /** Interval between retries in ms (default: 1_000) */
  intervalMs?: number;
  /** Human-readable description for timeout error */
  description?: string;
}
