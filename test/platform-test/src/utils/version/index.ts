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
 * Version helpers for the `--run-up-to-version` filter.
 *
 * A test declares the oldest APIM version it needs with a `@since-<version>` tag
 * in its title (see `since()` in `e2e/helpers/tags.ts`). When the suite is run
 * with `--run-up-to-version <v>` (which sets `E2E_MAX_VERSION`), any test whose
 * `@since` version is newer than `<v>` is skipped. A test with no `@since` tag is
 * treated as baseline and always runs. The skip itself is wired in
 * `e2e/setup.ts`; everything here is pure and unit-tested.
 */

/** A parsed (major, minor, patch) tuple. */
export type Version = [number, number, number];

/**
 * Parse a dotted version string into (major, minor, patch). Tolerates a leading
 * "v", a missing patch (defaults to 0), and any pre-release / build suffix:
 *   "4.12" -> [4,12,0]   "v4.12.3" -> [4,12,3]   "4.12.0-milestone.2" -> [4,12,0]
 * Returns undefined when there is no leading `major.minor` (e.g. "master-latest").
 */
export function parseVersion(input: string): Version | undefined {
  const match = /^v?(\d+)\.(\d+)(?:\.(\d+))?/.exec(input.trim());
  if (!match) return undefined;
  return [Number(match[1]), Number(match[2]), match[3] === undefined ? 0 : Number(match[3])];
}

/**
 * Compare two version strings by (major, minor, patch). Returns -1, 0, or 1.
 * If either side is unparseable, returns 0 so callers fail open (never gate by
 * mistake on something like "master-latest").
 */
export function compareVersions(a: string, b: string): number {
  const pa = parseVersion(a);
  const pb = parseVersion(b);
  if (!pa || !pb) return 0;
  for (let i = 0; i < 3; i++) {
    if (pa[i] !== pb[i]) return pa[i] < pb[i] ? -1 : 1;
  }
  return 0;
}

const SINCE_RE = /@since-(\S+)/;

/** Extract the `@since-<version>` value from a test title, if present. */
export function sinceFromTitle(title: string): string | undefined {
  const match = SINCE_RE.exec(title);
  return match ? match[1] : undefined;
}

/**
 * Decide whether a test should be skipped when the run is capped at `maxVersion`
 * (the `--run-up-to-version` / `E2E_MAX_VERSION` value). Returns a human-readable
 * reason to skip, or undefined to run.
 *
 * - No cap set -> run everything.
 * - No `@since` tag -> baseline -> always run.
 * - `@since-X` with X strictly newer than the cap -> skip.
 */
export function versionSkipReason(
  title: string,
  maxVersion: string | undefined | null,
): string | undefined {
  if (!maxVersion) return undefined;
  const since = sinceFromTitle(title);
  if (since === undefined) return undefined;
  if (compareVersions(since, maxVersion) > 0) {
    return `requires APIM ${since}, run capped at ${maxVersion} (--run-up-to-version)`;
  }
  return undefined;
}
