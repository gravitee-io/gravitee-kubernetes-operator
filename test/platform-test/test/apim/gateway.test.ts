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
import { Gateway, createMapi } from "../../src/assertions/apim/index.js";
import type { GatewayConfig } from "../../src/types/gateway.js";
import type { FetchFn } from "../../src/types/http.js";

// ── Mock helpers ───────────────────────────────────────────────

/**
 * Build a fetch mock that returns each status code in sequence,
 * then repeats the final one indefinitely.
 */
function mockFetch(statusSequence: number[]): { fetchFn: FetchFn; calls: string[] } {
  const calls: string[] = [];
  let callIndex = 0;

  const fetchFn: FetchFn = (input, _init) => {
    const url = typeof input === "string" ? input : input instanceof URL ? input.toString() : input.url;
    calls.push(url);
    const status = statusSequence[Math.min(callIndex++, statusSequence.length - 1)];
    return Promise.resolve(new Response(null, { status }));
  };

  return { fetchFn, calls };
}

/** Default fast config for tests that expect success after a few retries */
const fastConfig = (baseUrl = "http://gw.test"): GatewayConfig => ({
  baseUrl,
  retryIntervalMs: 10,
  maxRetryMs: 200,
});

/** Config with essentially no retries — for tests that expect immediate failure */
const noRetryConfig = (baseUrl = "http://gw.test"): GatewayConfig => ({
  baseUrl,
  retryIntervalMs: 10,
  maxRetryMs: 50,
});

// ── assertResponds ─────────────────────────────────────────────

describe("Gateway.assertResponds", () => {
  it("resolves immediately when status matches on first call", async () => {
    const { fetchFn, calls } = mockFetch([200]);
    const gw = new Gateway(fastConfig(), fetchFn);
    await expect(gw.assertResponds("/echo", { status: 200 })).resolves.toBeUndefined();
    expect(calls).toHaveLength(1);
    expect(calls[0]).toBe("http://gw.test/echo");
  });

  it("retries and resolves when status eventually matches", async () => {
    const { fetchFn, calls } = mockFetch([404, 404, 200]);
    const gw = new Gateway(fastConfig(), fetchFn);
    await expect(gw.assertResponds("/slow", { status: 200 })).resolves.toBeUndefined();
    expect(calls.length).toBeGreaterThanOrEqual(3);
  });

  it("throws AssertionError when status never matches within timeout", async () => {
    const { fetchFn } = mockFetch([503]);
    const gw = new Gateway(noRetryConfig(), fetchFn);
    await expect(gw.assertResponds("/bad", { status: 200 })).rejects.toSatisfy(
      (e: unknown) => e instanceof assert.AssertionError,
    );
  });

  it("AssertionError has correct operator, actual, and expected", async () => {
    const { fetchFn } = mockFetch([503]);
    const gw = new Gateway(noRetryConfig(), fetchFn);
    try {
      await gw.assertResponds("/bad", { status: 200 });
      expect.fail("expected AssertionError");
    } catch (err) {
      expect(err).toBeInstanceOf(assert.AssertionError);
      const ae = err as assert.AssertionError;
      expect(ae.operator).toBe("assertGatewayResponds");
      expect(ae.actual).toBe(503);
      expect(ae.expected).toBe(200);
    }
  });

  it("normalizes path — adds leading slash if missing", async () => {
    const { fetchFn, calls } = mockFetch([200]);
    const gw = new Gateway(fastConfig(), fetchFn);
    await gw.assertResponds("echo", { status: 200 });
    expect(calls[0]).toBe("http://gw.test/echo");
  });

  it("strips trailing slash from baseUrl", async () => {
    const { fetchFn, calls } = mockFetch([200]);
    const gw = new Gateway(fastConfig("http://gw.test/"), fetchFn);
    await gw.assertResponds("/echo", { status: 200 });
    expect(calls[0]).toBe("http://gw.test/echo");
  });

  it("passes custom headers to each fetch call", async () => {
    const sentHeaders: string[] = [];
    const fetchFn: FetchFn = (_input, init) => {
      const h = init?.headers as Record<string, string> | undefined;
      sentHeaders.push(h?.["authorization"] ?? "none");
      return Promise.resolve(new Response(null, { status: 200 }));
    };
    const gw = new Gateway(fastConfig(), fetchFn);
    await gw.assertResponds("/secure", { status: 200, headers: { authorization: "Bearer abc" } });
    expect(sentHeaders[0]).toBe("Bearer abc");
  });

  it("forwards description when present", async () => {
    const { fetchFn } = mockFetch([404]);
    const gw = new Gateway(noRetryConfig(), fetchFn);
    try {
      await gw.assertResponds("/x", { status: 200 }, "my custom description");
      expect.fail("expected AssertionError");
    } catch (err) {
      expect(err).toBeInstanceOf(assert.AssertionError);
      const ae = err as assert.AssertionError;
      expect(ae.message).toContain("my custom description");
    }
  });
});

// ── assertNotResponds ──────────────────────────────────────────

describe("Gateway.assertNotResponds", () => {
  it("resolves immediately when status is already different", async () => {
    const { fetchFn, calls } = mockFetch([401]);
    const gw = new Gateway(fastConfig(), fetchFn);
    await expect(gw.assertNotResponds("/locked", { notStatus: 200 })).resolves.toBeUndefined();
    expect(calls).toHaveLength(1);
  });

  it("retries and resolves once status changes away from notStatus", async () => {
    const { fetchFn, calls } = mockFetch([200, 200, 401]);
    const gw = new Gateway(fastConfig(), fetchFn);
    await expect(gw.assertNotResponds("/locked", { notStatus: 200 })).resolves.toBeUndefined();
    expect(calls.length).toBeGreaterThanOrEqual(3);
  });

  it("throws AssertionError when status always equals notStatus", async () => {
    const { fetchFn } = mockFetch([200]);
    const gw = new Gateway(noRetryConfig(), fetchFn);
    await expect(gw.assertNotResponds("/open", { notStatus: 200 })).rejects.toSatisfy(
      (e: unknown) => e instanceof assert.AssertionError,
    );
  });

  it("AssertionError has correct operator, actual, and expected", async () => {
    const { fetchFn } = mockFetch([200]);
    const gw = new Gateway(noRetryConfig(), fetchFn);
    try {
      await gw.assertNotResponds("/open", { notStatus: 200 });
      expect.fail("expected AssertionError");
    } catch (err) {
      expect(err).toBeInstanceOf(assert.AssertionError);
      const ae = err as assert.AssertionError;
      expect(ae.operator).toBe("assertGatewayNotResponds");
      expect(ae.actual).toBe(200);
      expect(ae.expected).not.toBe(200);
    }
  });
});

// ── Mapi.gateway() factory ─────────────────────────────────────

describe("Mapi.gateway()", () => {
  it("returns a Gateway instance", () => {
    const apim = createMapi({ baseUrl: "http://mapi.test", auth: { type: "basic", username: "user", password: "pass" } });
    const gw = apim.gateway({ baseUrl: "http://gw.test" });
    expect(gw).toBeInstanceOf(Gateway);
  });

  it("passes custom fetchFn through to Gateway", async () => {
    const apim = createMapi({ baseUrl: "http://mapi.test", auth: { type: "basic", username: "user", password: "pass" } });
    const { fetchFn, calls } = mockFetch([200]);
    const gw = apim.gateway({ baseUrl: "http://gw.test" }, fetchFn);
    await gw.assertResponds("/ping", { status: 200 });
    expect(calls).toHaveLength(1);
    expect(calls[0]).toContain("gw.test");
  });
});
