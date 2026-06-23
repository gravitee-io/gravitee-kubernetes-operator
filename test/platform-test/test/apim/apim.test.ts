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

import { describe, it, expect, vi, afterEach } from "vitest";
import { AssertionError } from "node:assert";
import { Mapi } from "../../src/assertions/apim/index.js";

/** Create a Mapi with globalThis.fetch stubbed to return the given body. */
function createMockMapi(body: unknown, status = 200) {
  const mockFetch = vi.fn<typeof fetch>().mockResolvedValue(
    new Response(JSON.stringify(body), {
      status,
      statusText: status === 200 ? "OK" : "Not Found",
      headers: { "Content-Type": "application/json" },
    }),
  );
  vi.stubGlobal("fetch", mockFetch);

  const mapi = new Mapi({
    baseUrl: "http://localhost:8083",
    auth: { type: "basic", username: "admin", password: "admin" },
  });

  return { mapi, mockFetch };
}

describe("Mapi", () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });
  const fakeApi = {
    id: "api-1",
    name: "Petstore API",
    definitionVersion: "V4",
    type: "PROXY",
    state: "STARTED",
    visibility: "PRIVATE",
    lifecycleState: "PUBLISHED",
    apiVersion: "1.0",
    createdAt: "2025-01-01T00:00:00Z",
    updatedAt: "2025-01-01T00:00:00Z",
    owner: { id: "user-1", displayName: "Admin", type: "USER" },
    listeners: [{ type: "HTTP", entrypoints: [{ type: "http-proxy" }] }],
    endpointGroups: [{ name: "Default", type: "http-proxy" }],
  };

  describe("assertApiMatches", () => {
    it("resolves when API matches expected partial", async () => {
      const { mapi } = createMockMapi(fakeApi);
      await expect(
        mapi.assertApiMatches("api-1", { name: "Petstore API", state: "STARTED" }),
      ).resolves.toBeUndefined();
    });

    it("throws AssertionError when field mismatches", async () => {
      const { mapi } = createMockMapi(fakeApi);
      await expect(
        mapi.assertApiMatches("api-1", { name: "Wrong Name" }),
      ).rejects.toThrow(AssertionError);
    });

    it("supports deep nested partial matching", async () => {
      const { mapi } = createMockMapi(fakeApi);
      await expect(
        mapi.assertApiMatches("api-1", {
          listeners: [{ type: "HTTP" }],
        }),
      ).resolves.toBeUndefined();
    });

    it("throws AssertionError on nested mismatch", async () => {
      const { mapi } = createMockMapi(fakeApi);
      await expect(
        mapi.assertApiMatches("api-1", {
          listeners: [{ type: "TCP" }],
        }),
      ).rejects.toThrow(AssertionError);
    });
  });

  describe("assertApiStarted / assertApiStopped", () => {
    it("passes when API is STARTED", async () => {
      const { mapi } = createMockMapi(fakeApi);
      await expect(mapi.assertApiStarted("api-1")).resolves.toBeUndefined();
    });

    it("fails when API is not STOPPED", async () => {
      const { mapi } = createMockMapi(fakeApi);
      await expect(mapi.assertApiStopped("api-1")).rejects.toThrow(AssertionError);
    });
  });

  describe("checkApiMatches", () => {
    it("returns report with pass=true on match", async () => {
      const { mapi } = createMockMapi(fakeApi);
      const report = await mapi.checkApiMatches("api-1", { state: "STARTED" });
      expect(report.pass).toBe(true);
      expect(report.failures).toHaveLength(0);
    });

    it("returns report with failures on mismatch (does not throw)", async () => {
      const { mapi } = createMockMapi(fakeApi);
      const report = await mapi.checkApiMatches("api-1", { state: "STOPPED" });
      expect(report.pass).toBe(false);
      expect(report.failures).toHaveLength(1);
    });
  });

  describe("fetchApi error handling", () => {
    it("throws Error when API returns non-200", async () => {
      const { mapi } = createMockMapi({ message: "Not found" }, 404);
      await expect(mapi.fetchApi("missing")).rejects.toThrow("Failed to fetch API missing: 404");
    });
  });

  describe("HTTP client wiring", () => {
    it("calls the correct v2 management API path", async () => {
      const { mapi, mockFetch } = createMockMapi(fakeApi);
      await mapi.assertApiMatches("api-1", { name: "Petstore API" });
      expect(mockFetch).toHaveBeenCalledOnce();
      const url = mockFetch.mock.calls[0][0] as string;
      expect(url).toBe(
        "http://localhost:8083/management/v2/environments/DEFAULT/apis/api-1",
      );
    });
  });

  describe("assertApiHttpStatus", () => {
    it("resolves when actual status matches expected", async () => {
      const { mapi } = createMockMapi(fakeApi, 200);
      await expect(mapi.assertApiHttpStatus("api-1", 200)).resolves.toBeUndefined();
    });

    it("resolves when expected status is 404 and API returns 404", async () => {
      const { mapi } = createMockMapi({ message: "Not found" }, 404);
      await expect(mapi.assertApiHttpStatus("api-1", 404)).resolves.toBeUndefined();
    });

    it("throws AssertionError when status does not match", async () => {
      const { mapi } = createMockMapi(fakeApi, 200);
      await expect(mapi.assertApiHttpStatus("api-1", 404)).rejects.toThrow(AssertionError);
    });

    it("AssertionError carries actual and expected status codes", async () => {
      const { mapi } = createMockMapi({ message: "Not found" }, 404);
      try {
        await mapi.assertApiHttpStatus("api-1", 200);
        expect.unreachable("should have thrown");
      } catch (err) {
        expect(err).toBeInstanceOf(AssertionError);
        const ae = err as AssertionError;
        expect(ae.actual).toBe(404);
        expect(ae.expected).toBe(200);
        expect(ae.operator).toBe("assertApiStatus");
      }
    });
  });
});
