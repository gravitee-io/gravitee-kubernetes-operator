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
import { Gateway } from "../assertions/apim/gateway.js";
import { createTlsFetch } from "../utils/http/tls.js";
import type { GatewayConfig } from "../types/gateway.js";
import type { FetchFn } from "../types/http.js";
import { loadGraviteeConfig } from "./config.js";

export interface AssertGatewayFlags {
  /** Gateway-relative path to assert against, e.g. /petstore */
  path: string;
  /** Expected HTTP status code (mutually exclusive with notStatus) */
  status?: number;
  /** HTTP status code that must NOT appear (mutually exclusive with status) */
  notStatus?: number;
  /** Gateway base URL (default: http://localhost:30082) */
  gateway?: string;
  /** Value for the Authorization request header */
  authorization?: string;
  /** Additional headers in "Key: Value" format */
  headers?: string[];
  /** Path to client certificate PEM file (requires --key) */
  cert?: string;
  /** Path to client private key PEM file (requires --cert) */
  key?: string;
  /** Path to CA certificate PEM file */
  cacert?: string;
  /** Total retry timeout in ms (default: 30_000) */
  timeoutMs?: number;
  /** Retry interval in ms (default: 500) */
  retryIntervalMs?: number;
  /** Path to config.yaml (default: CWD) */
  configPath?: string;
}

/**
 * Execute the `assert-gateway` CLI subcommand.
 *
 * Constructs a Gateway and runs assertResponds or assertNotResponds
 * depending on the flags provided. Throws AssertionError on mismatch.
 *
 * @param flags   Parsed CLI flags
 * @param inject  Optional fetch function for testing (overrides TLS auto-detection)
 */
export async function assertGatewayCommand(flags: AssertGatewayFlags, inject?: FetchFn): Promise<void> {
  // ── Validate ──────────────────────────────────────────────

  if (!flags.path) {
    throw new Error("assert-gateway: --path is required");
  }
  if (flags.status === undefined && flags.notStatus === undefined) {
    throw new Error("assert-gateway: either --status or --not-status is required");
  }
  if (flags.status !== undefined && flags.notStatus !== undefined) {
    throw new Error("assert-gateway: --status and --not-status are mutually exclusive");
  }

  // ── Build request headers ─────────────────────────────────

  const requestHeaders: Record<string, string> = {};
  if (flags.authorization) {
    requestHeaders["authorization"] = flags.authorization;
  }
  for (const h of flags.headers ?? []) {
    const idx = h.indexOf(":");
    if (idx !== -1) {
      requestHeaders[h.slice(0, idx).trim().toLowerCase()] = h.slice(idx + 1).trim();
    }
  }

  // ── Build gateway config ──────────────────────────────────
  // --gateway flag wins; if not given, try config file; fall back to default.

  let gatewayBaseUrl = flags.gateway;
  if (!gatewayBaseUrl) {
    try {
      const cfg = await loadGraviteeConfig(flags.configPath);
      gatewayBaseUrl = cfg.gateway?.baseUrl;
    } catch {
      // Config file is optional for assert-gateway — ignore load errors
    }
  }

  const gatewayConfig: GatewayConfig = {
    baseUrl: gatewayBaseUrl ?? "http://localhost:30082",
    maxRetryMs: flags.timeoutMs,
    retryIntervalMs: flags.retryIntervalMs,
  };

  // ── Resolve fetch function ────────────────────────────────
  // inject takes priority (for tests); then TLS files; then undefined (native fetch)

  let fetchFn: FetchFn | undefined = inject;
  if (!fetchFn && flags.cert && flags.key) {
    const cert = await readFile(flags.cert);
    const key = await readFile(flags.key);
    const ca = flags.cacert ? await readFile(flags.cacert) : undefined;
    fetchFn = createTlsFetch({ cert, key, ca, rejectUnauthorized: ca !== undefined });
  }

  // ── Assert ────────────────────────────────────────────────

  const gatewayAssert = new Gateway(gatewayConfig, fetchFn);

  if (flags.status !== undefined) {
    await gatewayAssert.assertResponds(flags.path, {
      status: flags.status,
      headers: requestHeaders,
    });
  } else {
    await gatewayAssert.assertNotResponds(flags.path, {
      notStatus: flags.notStatus!,
      headers: requestHeaders,
    });
  }
}
