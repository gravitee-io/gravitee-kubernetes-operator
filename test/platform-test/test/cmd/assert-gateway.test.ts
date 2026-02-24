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
import assert from "node:assert";
import { assertGatewayCommand, type AssertGatewayFlags } from "../../src/cmd/assert-gateway.js";
import type { FetchFn } from "../../src/types/http.js";
import { parseArgs } from "../../src/cmd/index.js";

// ── Mock helpers ───────────────────────────────────────────────

function mockFetch(statusSequence: number[]): FetchFn {
  let callIndex = 0;
  return (_input, _init) => {
    const status = statusSequence[Math.min(callIndex++, statusSequence.length - 1)];
    return Promise.resolve(new Response(null, { status }));
  };
}

/** Minimal fast flags — skips TLS file loading, uses injected fetchFn */
const baseFlags = (overrides: Partial<AssertGatewayFlags> = {}): AssertGatewayFlags => ({
  path: "/test",
  status: 200,
  gateway: "http://localhost:30082",
  timeoutMs: 200,
  retryIntervalMs: 10,
  ...overrides,
});

// ── Flag validation ────────────────────────────────────────────

describe("assertGatewayCommand — flag validation", () => {
  it("throws when --path is empty", async () => {
    await expect(
      assertGatewayCommand(baseFlags({ path: "" })),
    ).rejects.toThrow(/path/i);
  });

  it("throws when neither --status nor --not-status is given", async () => {
    await expect(
      assertGatewayCommand(baseFlags({ status: undefined, notStatus: undefined })),
    ).rejects.toThrow(/status/i);
  });

  it("throws when both --status and --not-status are given", async () => {
    await expect(
      assertGatewayCommand(baseFlags({ status: 200, notStatus: 401 })),
    ).rejects.toThrow(/mutually exclusive/i);
  });
});

// ── assertResponds path ────────────────────────────────────────

describe("assertGatewayCommand — --status path", () => {
  it("resolves when gateway returns expected status", async () => {
    const fetchFn = mockFetch([200]);
    await expect(
      assertGatewayCommand(baseFlags({ status: 200 }), fetchFn),
    ).resolves.toBeUndefined();
  });

  it("rejects with AssertionError when gateway never returns expected status", async () => {
    const fetchFn = mockFetch([503]);
    await expect(
      assertGatewayCommand(baseFlags({ status: 200 }), fetchFn),
    ).rejects.toSatisfy((e: unknown) => e instanceof assert.AssertionError);
  });

  it("resolves after retries when status eventually matches", async () => {
    const fetchFn = mockFetch([404, 404, 200]);
    await expect(
      assertGatewayCommand(baseFlags({ status: 200 }), fetchFn),
    ).resolves.toBeUndefined();
  });
});

// ── assertNotResponds path ─────────────────────────────────────

describe("assertGatewayCommand — --not-status path", () => {
  it("resolves when gateway returns a status other than notStatus", async () => {
    const fetchFn = mockFetch([401]);
    await expect(
      assertGatewayCommand(baseFlags({ status: undefined, notStatus: 200 }), fetchFn),
    ).resolves.toBeUndefined();
  });

  it("rejects with AssertionError when gateway always returns notStatus", async () => {
    const fetchFn = mockFetch([200]);
    await expect(
      assertGatewayCommand(baseFlags({ status: undefined, notStatus: 200 }), fetchFn),
    ).rejects.toSatisfy((e: unknown) => e instanceof assert.AssertionError);
  });
});

// ── Header building ────────────────────────────────────────────

describe("assertGatewayCommand — header building", () => {
  it("passes --authorization header to fetch", async () => {
    const sentHeaders: Record<string, string>[] = [];
    const fetchFn: FetchFn = (_input, init) => {
      sentHeaders.push((init?.headers ?? {}) as Record<string, string>);
      return Promise.resolve(new Response(null, { status: 200 }));
    };
    await assertGatewayCommand(
      baseFlags({ status: 200, authorization: "Bearer token123" }),
      fetchFn,
    );
    expect(sentHeaders[0]?.["authorization"]).toBe("Bearer token123");
  });

  it("merges multiple extra headers", async () => {
    const sentHeaders: Record<string, string>[] = [];
    const fetchFn: FetchFn = (_input, init) => {
      sentHeaders.push((init?.headers ?? {}) as Record<string, string>);
      return Promise.resolve(new Response(null, { status: 200 }));
    };
    await assertGatewayCommand(
      baseFlags({ status: 200, headers: ["X-Tenant: acme", "X-Version: 2"] }),
      fetchFn,
    );
    expect(sentHeaders[0]?.["x-tenant"]).toBe("acme");
    expect(sentHeaders[0]?.["x-version"]).toBe("2");
  });
});

// ── parseArgs integration ──────────────────────────────────────

describe("parseArgs — assert-gateway flags", () => {
  it("parses --path, --status, --gateway", () => {
    const { subcommand, flags } = parseArgs([
      "assert-gateway",
      "--path", "/api",
      "--status", "200",
      "--gateway", "http://gw:8082",
    ]);
    expect(subcommand).toBe("assert-gateway");
    expect(flags["path"]).toBe("/api");
    expect(flags["status"]).toBe("200");
    expect(flags["gateway"]).toBe("http://gw:8082");
  });

  it("parses --not-status flag (hyphen preserved)", () => {
    const { flags } = parseArgs([
      "assert-gateway",
      "--path", "/api",
      "--not-status", "200",
    ]);
    expect(flags["not-status"]).toBe("200");
  });

  it("parses --cert, --key, --cacert flags", () => {
    const { flags } = parseArgs([
      "assert-gateway",
      "--path", "/mtls",
      "--status", "200",
      "--cert", "/tmp/client.crt",
      "--key", "/tmp/client.key",
      "--cacert", "/tmp/ca.crt",
    ]);
    expect(flags["cert"]).toBe("/tmp/client.crt");
    expect(flags["key"]).toBe("/tmp/client.key");
    expect(flags["cacert"]).toBe("/tmp/ca.crt");
  });

  it("parses --timeout and --retry-interval flags", () => {
    const { flags } = parseArgs([
      "assert-gateway",
      "--path", "/api",
      "--status", "200",
      "--timeout", "60000",
      "--retry-interval", "1000",
    ]);
    expect(flags["timeout"]).toBe("60000");
    expect(flags["retry-interval"]).toBe("1000");
  });
});
