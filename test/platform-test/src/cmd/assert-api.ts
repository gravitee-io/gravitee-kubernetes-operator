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

import { readFile } from "node:fs/promises";
import { parse } from "yaml";
import type { DeepPartial } from "../types/match.js";
import type { Api, ApiState, HttpListener } from "../types/apim.js";
import { loadGraviteeConfig, createMapiFromConfig } from "./config.js";

export interface AssertApiFlags {
  /** API ID to fetch and assert */
  apiId: string;
  /** Assert the management API returns exactly this HTTP status code (e.g. 404 after deletion) */
  expectedStatus?: number;
  /** Expected lifecycle state (e.g. "STARTED", "STOPPED") */
  state?: string;
  /** Expected listener path (e.g. "/petstore") */
  path?: string;
  /** Raw JSON partial to merge into the assertion (e.g. '{"categories":["finance"]}') */
  match?: string;
  /** Path to a YAML file containing the expected partial API shape */
  expectFile?: string;
  /** Pre-parsed content from --expect file (for testability; set by assertApiCommand) */
  expectContent?: Record<string, unknown>;
  /** Path to config.yaml config file */
  configPath?: string;
}

/**
 * Load and parse a YAML expect file as a plain object.
 * Throws a descriptive error if the file cannot be read or parsed.
 */
export async function loadExpectFile(filePath: string): Promise<Record<string, unknown>> {
  let content: string;
  try {
    content = await readFile(filePath, "utf-8");
  } catch (err) {
    throw new Error(
      `assert-api: cannot read --expect file "${filePath}": ${err instanceof Error ? err.message : String(err)}`,
    );
  }

  let parsed: unknown;
  try {
    parsed = parse(content);
  } catch (err) {
    throw new Error(
      `assert-api: --expect file "${filePath}" is not valid YAML: ${err instanceof Error ? err.message : String(err)}`,
    );
  }

  if (parsed === null || parsed === undefined || typeof parsed !== "object" || Array.isArray(parsed)) {
    throw new Error(
      `assert-api: --expect file "${filePath}" must contain a YAML mapping (object), got ${Array.isArray(parsed) ? "array" : String(parsed)}`,
    );
  }

  return parsed as Record<string, unknown>;
}

/**
 * Build a DeepPartial<Api> from CLI flags.
 *
 * Merge order (later wins):
 *   1. `--expect <file>` content (loaded externally, passed as `expectContent`)
 *   2. `--match <JSON>` inline partial
 *   3. `--state`, `--path` individual flags
 */
export function buildPartial(flags: AssertApiFlags): DeepPartial<Api> {
  const partial: DeepPartial<Api> = {};

  if (flags.expectContent) {
    Object.assign(partial, flags.expectContent);
  }

  if (flags.match) {
    try {
      Object.assign(partial, JSON.parse(flags.match));
    } catch {
      throw new Error(`assert-api: --match value is not valid JSON: ${flags.match}`);
    }
  }

  if (flags.state) {
    partial.state = flags.state as ApiState;
  }

  if (flags.path) {
    // Map to V4 HTTP listener path shape; DeepPartial allows partial Listener arrays
    (partial as { listeners?: DeepPartial<HttpListener>[] }).listeners = [
      { paths: [{ path: flags.path }] },
    ];
  }

  return partial;
}

/**
 * Execute the `assert-api` subcommand.
 *
 * Loads config, fetches the API, runs the partial assertion.
 * Throws AssertionError on mismatch (caller should handle exit code).
 */
export async function assertApiCommand(flags: AssertApiFlags): Promise<void> {
  if (!flags.apiId) {
    throw new Error("assert-api: --api-id is required");
  }

  if (flags.expectFile) {
    flags.expectContent = await loadExpectFile(flags.expectFile);
  }

  const config = await loadGraviteeConfig(flags.configPath);
  const mapi = createMapiFromConfig(config);

  if (flags.expectedStatus !== undefined) {
    await mapi.assertApiHttpStatus(flags.apiId, flags.expectedStatus);
  } else {
    await mapi.assertApiMatches(flags.apiId, buildPartial(flags));
  }
}
