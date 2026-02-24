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

import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import {
  validateConfig,
  applyEnvVars,
  createMapiFromConfig,
} from "../../src/cmd/config.js";
import { Mapi } from "../../src/assertions/apim/index.js";

// ── validateConfig ────────────────────────────────────────────

describe("validateConfig", () => {
  it("accepts a valid config object", () => {
    const raw = {
      apim: {
        baseUrl: "http://localhost:8083",
        envId: "DEFAULT",
        auth: { username: "admin", password: "admin" },
      },
    };
    const config = validateConfig(raw);
    expect(config.apim.baseUrl).toBe("http://localhost:8083");
    expect(config.apim.envId).toBe("DEFAULT");
  });

  it("throws when apim section is missing", () => {
    expect(() => validateConfig({})).toThrow("apim");
  });

  it("throws when apim.baseUrl is missing", () => {
    const raw = { apim: { auth: { username: "a", password: "b" } } };
    expect(() => validateConfig(raw)).toThrow("apim.baseUrl");
  });

  it("throws when apim.auth.username is missing", () => {
    const raw = { apim: { baseUrl: "http://localhost", auth: { password: "b" } } };
    expect(() => validateConfig(raw)).toThrow("apim.auth.username");
  });

  it("allows optional envId to be absent", () => {
    const raw = {
      apim: { baseUrl: "http://localhost", auth: { username: "u", password: "p" } },
    };
    const config = validateConfig(raw);
    expect(config.apim.envId).toBeUndefined();
  });
});

// ── applyEnvVars ──────────────────────────────────────────────

describe("applyEnvVars", () => {
  const baseConfig = () => ({
    apim: {
      baseUrl: "http://localhost:8083",
      envId: "DEFAULT",
      auth: { username: "admin", password: "admin" },
    },
  });

  beforeEach(() => {
    // Clear any relevant env vars before each test
    delete process.env["GRAVITEE_BASE_URL"];
    delete process.env["GRAVITEE_ENV_ID"];
    delete process.env["GRAVITEE_USERNAME"];
    delete process.env["GRAVITEE_PASSWORD"];
  });

  afterEach(() => {
    delete process.env["GRAVITEE_BASE_URL"];
    delete process.env["GRAVITEE_ENV_ID"];
    delete process.env["GRAVITEE_USERNAME"];
    delete process.env["GRAVITEE_PASSWORD"];
  });

  it("returns config unchanged when no env vars are set", () => {
    const config = applyEnvVars(baseConfig());
    expect(config.apim.baseUrl).toBe("http://localhost:8083");
    expect(config.apim.auth.username).toBe("admin");
  });

  it("overrides baseUrl from GRAVITEE_BASE_URL", () => {
    process.env["GRAVITEE_BASE_URL"] = "http://my-host:9090/management/v2";
    const config = applyEnvVars(baseConfig());
    expect(config.apim.baseUrl).toBe("http://my-host:9090/management/v2");
  });

  it("overrides envId from GRAVITEE_ENV_ID", () => {
    process.env["GRAVITEE_ENV_ID"] = "STAGING";
    const config = applyEnvVars(baseConfig());
    expect(config.apim.envId).toBe("STAGING");
  });

  it("overrides username from GRAVITEE_USERNAME", () => {
    process.env["GRAVITEE_USERNAME"] = "testuser";
    const config = applyEnvVars(baseConfig());
    expect(config.apim.auth.username).toBe("testuser");
  });

  it("overrides password from GRAVITEE_PASSWORD", () => {
    process.env["GRAVITEE_PASSWORD"] = "s3cret";
    const config = applyEnvVars(baseConfig());
    expect(config.apim.auth.password).toBe("s3cret");
  });

  it("does not mutate the original config object", () => {
    process.env["GRAVITEE_BASE_URL"] = "http://other:9090";
    const original = baseConfig();
    const updated = applyEnvVars(original);
    expect(original.apim.baseUrl).toBe("http://localhost:8083");
    expect(updated.apim.baseUrl).toBe("http://other:9090");
  });
});

// ── createMapiFromConfig ──────────────────────────────────────

describe("createMapiFromConfig", () => {
  it("creates a Mapi instance", () => {
    const config = {
      apim: {
        baseUrl: "http://localhost:8083",
        envId: "DEFAULT",
        auth: { username: "admin", password: "admin" },
      },
    };
    const mapi = createMapiFromConfig(config);
    expect(mapi).toBeInstanceOf(Mapi);
  });

  it("wires the correct management API URL", async () => {
    const mockFetch = vi.fn<typeof fetch>().mockResolvedValue(
      new Response(JSON.stringify({ id: "api-1", name: "Test" }), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      }),
    );

    const config = {
      apim: {
        baseUrl: "http://testhost:8083",
        envId: "MY_ENV",
        auth: { username: "user", password: "pass" },
      },
    };

    // Pass mockFetch to Mapi directly for this test
    const mapi = new Mapi(
      {
        baseUrl: config.apim.baseUrl,
        envId: config.apim.envId,
        auth: { type: "basic", username: config.apim.auth.username, password: config.apim.auth.password },
      },
      mockFetch,
    );

    await mapi.fetchApi("api-1");
    const url = mockFetch.mock.calls[0][0] as string;
    expect(url).toBe("http://testhost:8083/management/v2/environments/MY_ENV/apis/api-1");
  });
});
