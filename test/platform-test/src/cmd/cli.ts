#!/usr/bin/env node
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
 * @gravitee/platform-test CLI
 *
 * Usage:
 *   npx @gravitee/platform-test <subcommand> [flags]
 *
 * Subcommands:
 *   assert-api   Assert that an API matches an expected partial shape
 *
 * Exit codes:
 *   0  Assertion passed
 *   1  Assertion failed or usage error
 *   2  Config / network error
 */

import { AssertionError } from "node:assert";
import { fileURLToPath } from "node:url";
import { parseArgs } from "./parse-args.js";
import { assertApiCommand } from "./assert-api.js";
import { assertGatewayCommand } from "./assert-gateway.js";

// ── Dispatch ──────────────────────────────────────────────────

async function main(): Promise<void> {
  const { subcommand, flags } = parseArgs(process.argv.slice(2));

  switch (subcommand) {
    case "assert-api": {
      await assertApiCommand({
        apiId: flags["api-id"] ?? "",
        expectedStatus: flags["status"] !== undefined ? parseInt(flags["status"], 10) : undefined,
        state: flags["state"],
        path: flags["path"],
        match: flags["match"],
        matchFile: flags["match-file"],
        configPath: flags["config"],
      });
      break;
    }

    case "assert-gateway": {
      await assertGatewayCommand({
        path: flags["path"] ?? "",
        status: flags["status"] !== undefined ? parseInt(flags["status"], 10) : undefined,
        notStatus: flags["not-status"] !== undefined ? parseInt(flags["not-status"], 10) : undefined,
        gateway: flags["gateway"],
        authorization: flags["authorization"],
        cert: flags["cert"],
        key: flags["key"],
        cacert: flags["cacert"],
        timeoutMs: flags["timeout"] !== undefined ? parseInt(flags["timeout"], 10) : undefined,
        retryIntervalMs: flags["retry-interval"] !== undefined ? parseInt(flags["retry-interval"], 10) : undefined,
        configPath: flags["config"],
      });
      break;
    }

    case "help":
    case "--help":
    case "-h":
    case "": {
      process.stdout.write(
        [
          "Usage: platform-test <subcommand> [flags]",
          "",
          "Subcommands:",
          "  assert-api      Assert that an APIM API matches expected state",
          "  assert-gateway  Assert that the APIM gateway responds with expected status",
          "",
          "Flags for assert-api:",
          "  --api-id <id>              API ID to assert (required)",
          "  --status <code>            Expected HTTP status code (e.g. 404); skips property checks",
          "  --state  <state>           Expected state (e.g. STARTED, STOPPED)",
          "  --path   <path>            Expected listener path (e.g. /petstore)",
          "  --match  <json>            Arbitrary JSON partial to assert (e.g. '{\"categories\":[\"finance\"]}')",
          "  --match-file <file>        YAML file with expected partial API shape (merged before --match/--state/--path)",
          "  --config <file>            Path to config.yaml (default: CWD)",
          "",
          "Flags for assert-gateway:",
          "  --path           <path>    Gateway path to call, e.g. /petstore (required)",
          "  --status         <code>    Expected HTTP status code (mutually exclusive with --not-status)",
          "  --not-status     <code>    Status code that must NOT appear",
          "  --gateway        <url>     Gateway base URL (default: http://localhost:30082)",
          "  --authorization  <value>   Authorization header value (e.g. Bearer <token>)",
          "  --cert           <file>    Client certificate PEM (mTLS)",
          "  --key            <file>    Client private key PEM (mTLS)",
          "  --cacert         <file>    CA certificate PEM",
          "  --timeout        <ms>      Total retry timeout in ms (default: 30000)",
          "  --retry-interval <ms>      Interval between retries in ms (default: 500)",
          "",
        ].join("\n"),
      );
      break;
    }

    default: {
      process.stderr.write(`Unknown subcommand: "${subcommand}". Run with --help for usage.\n`);
      process.exit(1);
    }
  }
}

// Entry point guard: only run when executed directly, not when imported as a module.
// In Node.js ESM, process.argv[1] is the path of the directly-executed script.
if (process.argv[1] === fileURLToPath(import.meta.url)) {
  main().catch((err: unknown) => {
    if (err instanceof AssertionError) {
      // Assertion failure: print readable diff to stderr, exit 1
      process.stderr.write(`Assertion failed:\n${err.message}\n`);
      process.exit(1);
    }

    // Config / network / unexpected error: exit 2
    const message = err instanceof Error ? err.message : String(err);
    process.stderr.write(`Error: ${message}\n`);
    process.exit(2);
  });
}
