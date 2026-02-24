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

import { describe, it, expect, vi } from "vitest";
import { AssertionError } from "node:assert";
import { buildPartial } from "../../src/cmd/assert-api.js";
import { parseArgs } from "../../src/cmd/index.js";
import { Mapi } from "../../src/assertions/apim/index.js";

// ── parseArgs ─────────────────────────────────────────────────

describe("parseArgs", () => {
  it("parses the subcommand as first positional", () => {
    const { subcommand } = parseArgs(["assert-api", "--api-id", "123"]);
    expect(subcommand).toBe("assert-api");
  });

  it("parses --key value flags", () => {
    const { flags } = parseArgs(["assert-api", "--api-id", "abc", "--state", "STARTED"]);
    expect(flags["api-id"]).toBe("abc");
    expect(flags["state"]).toBe("STARTED");
  });

  it("parses --key=value flags", () => {
    const { flags } = parseArgs(["assert-api", "--api-id=xyz", "--state=STOPPED"]);
    expect(flags["api-id"]).toBe("xyz");
    expect(flags["state"]).toBe("STOPPED");
  });

  it("parses --config flag", () => {
    const { flags } = parseArgs(["assert-api", "--api-id", "1", "--config", "/tmp/my.yaml"]);
    expect(flags["config"]).toBe("/tmp/my.yaml");
  });

  it("returns empty subcommand for empty argv", () => {
    const { subcommand, flags } = parseArgs([]);
    expect(subcommand).toBe("");
    expect(flags).toEqual({});
  });

  it("handles boolean flags (no value)", () => {
    const { flags } = parseArgs(["sub", "--verbose", "--state", "STARTED"]);
    expect(flags["verbose"]).toBe("true");
    expect(flags["state"]).toBe("STARTED");
  });

  it("parses --status flag", () => {
    const { flags } = parseArgs(["assert-api", "--api-id", "abc", "--status", "404"]);
    expect(flags["status"]).toBe("404");
  });
});

// ── buildPartial ──────────────────────────────────────────────

describe("buildPartial", () => {
  it("returns empty partial when no flags are set", () => {
    const partial = buildPartial({ apiId: "abc" });
    expect(partial).toEqual({});
  });

  it("includes state when --state is provided", () => {
    const partial = buildPartial({ apiId: "abc", state: "STARTED" });
    expect(partial.state).toBe("STARTED");
  });

  it("includes listener path when --path is provided", () => {
    const partial = buildPartial({ apiId: "abc", path: "/petstore" }) as {
      listeners: Array<{ paths: Array<{ path: string }> }>;
    };
    expect(partial.listeners).toHaveLength(1);
    expect(partial.listeners[0].paths?.[0].path).toBe("/petstore");
  });

  it("includes both state and path when both flags are provided", () => {
    const partial = buildPartial({ apiId: "abc", state: "STOPPED", path: "/v2" }) as {
      state: string;
      listeners: Array<{ paths: Array<{ path: string }> }>;
    };
    expect(partial.state).toBe("STOPPED");
    expect(partial.listeners[0].paths?.[0].path).toBe("/v2");
  });
});

// ── buildPartial with --expect content ────────────────────────

describe("buildPartial with expectContent", () => {
  it("merges expectContent into the partial", () => {
    const partial = buildPartial({
      apiId: "abc",
      expectContent: { name: "My API", visibility: "PUBLIC" },
    });
    expect(partial).toEqual({ name: "My API", visibility: "PUBLIC" });
  });

  it("--match overrides expectContent fields", () => {
    const partial = buildPartial({
      apiId: "abc",
      expectContent: { name: "From File", visibility: "PUBLIC" },
      match: '{"name":"From Match"}',
    });
    expect(partial).toEqual({ name: "From Match", visibility: "PUBLIC" });
  });

  it("--state overrides expectContent state", () => {
    const partial = buildPartial({
      apiId: "abc",
      expectContent: { state: "STOPPED", name: "My API" },
      state: "STARTED",
    });
    expect(partial.state).toBe("STARTED");
    expect((partial as { name?: string }).name).toBe("My API");
  });

  it("--path overrides expectContent listeners", () => {
    const partial = buildPartial({
      apiId: "abc",
      expectContent: { listeners: [{ type: "HTTP", paths: [{ path: "/from-file" }] }] },
      path: "/from-flag",
    }) as { listeners: Array<{ paths: Array<{ path: string }> }> };
    expect(partial.listeners[0].paths[0].path).toBe("/from-flag");
  });

  it("combines expectContent with all flag types", () => {
    const partial = buildPartial({
      apiId: "abc",
      expectContent: { categories: ["finance"], tags: ["internal"] },
      match: '{"visibility":"PUBLIC"}',
      state: "STARTED",
      path: "/petstore",
    });
    expect(partial).toMatchObject({
      categories: ["finance"],
      tags: ["internal"],
      visibility: "PUBLIC",
      state: "STARTED",
    });
  });

  it("treats absent expectContent the same as before", () => {
    const partial = buildPartial({ apiId: "abc", state: "STARTED" });
    expect(partial).toEqual({ state: "STARTED" });
  });
});

// ── assertApiCommand (mocked fetch) ──────────────────────────

describe("assertApiCommand via Mapi", () => {
  const fakeApi = {
    id: "api-1",
    name: "My API",
    definitionVersion: "V4",
    type: "PROXY",
    state: "STARTED",
    visibility: "PUBLIC",
    lifecycleState: "PUBLISHED",
    apiVersion: "1.0",
    createdAt: "2025-01-01T00:00:00Z",
    updatedAt: "2025-01-01T00:00:00Z",
    owner: { id: "u1", displayName: "Admin", type: "USER" },
    listeners: [
      { type: "HTTP", entrypoints: [{ type: "http-proxy" }], paths: [{ path: "/petstore" }] },
    ],
    endpointGroups: [{ name: "Default", type: "http-proxy" }],
  };

  function createMockApim(body: unknown, status = 200) {
    const mockFetch = vi.fn<typeof fetch>().mockResolvedValue(
      new Response(JSON.stringify(body), {
        status,
        statusText: status === 200 ? "OK" : "Not Found",
        headers: { "Content-Type": "application/json" },
      }),
    );
    const apim = new Mapi(
      { baseUrl: "http://localhost:8083", auth: { type: "basic", username: "admin", password: "admin" } },
      mockFetch,
    );
    return { apim, mockFetch };
  }

  it("passes when state matches", async () => {
    const { apim } = createMockApim(fakeApi);
    const partial = buildPartial({ apiId: "api-1", state: "STARTED" });
    await expect(apim.assertApiMatches("api-1", partial)).resolves.toBeUndefined();
  });

  it("throws AssertionError when state does not match", async () => {
    const { apim } = createMockApim(fakeApi);
    const partial = buildPartial({ apiId: "api-1", state: "STOPPED" });
    await expect(apim.assertApiMatches("api-1", partial)).rejects.toThrow(AssertionError);
  });

  it("passes when listener path matches", async () => {
    const { apim } = createMockApim(fakeApi);
    const partial = buildPartial({ apiId: "api-1", path: "/petstore" });
    await expect(apim.assertApiMatches("api-1", partial)).resolves.toBeUndefined();
  });

  it("throws AssertionError when listener path does not match", async () => {
    const { apim } = createMockApim(fakeApi);
    const partial = buildPartial({ apiId: "api-1", path: "/wrongpath" });
    await expect(apim.assertApiMatches("api-1", partial)).rejects.toThrow(AssertionError);
  });

  it("passes when --status matches a 404 response", async () => {
    const { apim } = createMockApim({ message: "Not found" }, 404);
    await expect(apim.assertApiHttpStatus("api-1", 404)).resolves.toBeUndefined();
  });

  it("throws AssertionError when --status does not match actual HTTP status", async () => {
    const { apim } = createMockApim(fakeApi, 200);
    await expect(apim.assertApiHttpStatus("api-1", 404)).rejects.toThrow(AssertionError);
  });
});
