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
import { join } from "node:path";
import YAML from "yaml";
import { Mapi } from "../assertions/apim/index.js";
import type { MapiConfig } from "../types/mapi.js";
import type { GraviteeTestConfig } from "../types/config.js";

/**
 * Validate and cast a raw parsed YAML object to GraviteeTestConfig.
 * Throws an informative error if required fields are missing.
 */
export function validateConfig(raw: Record<string, unknown>): GraviteeTestConfig {
  const apim = raw["apim"] as Record<string, unknown> | undefined;
  if (!apim || typeof apim !== "object") {
    throw new Error("Config missing required section: apim");
  }

  const auth = apim["auth"] as Record<string, unknown> | undefined;
  if (!auth || typeof auth !== "object") {
    throw new Error("Config missing required section: apim.auth");
  }

  const baseUrl = apim["baseUrl"] as string | undefined;
  if (!baseUrl) {
    throw new Error("Config missing required field: apim.baseUrl");
  }

  const username = auth["username"] as string | undefined;
  if (!username) {
    throw new Error("Config missing required field: apim.auth.username");
  }

  const password = auth["password"] as string | undefined;
  if (!password) {
    throw new Error("Config missing required field: apim.auth.password");
  }

  const gateway = raw["gateway"] as Record<string, unknown> | undefined;

  return {
    apim: {
      baseUrl,
      envId: apim["envId"] as string | undefined,
      auth: { username, password },
    },
    gateway: gateway
      ? {
          baseUrl: gateway["baseUrl"] as string | undefined,
          mtlsBaseUrl: gateway["mtlsBaseUrl"] as string | undefined,
        }
      : undefined,
  };
}

/**
 * Apply environment variable overrides to a loaded config.
 * Env vars take precedence over file values.
 *
 * | Env var                    | Config field               |
 * |----------------------------|----------------------------|
 * | GRAVITEE_BASE_URL          | apim.baseUrl               |
 * | GRAVITEE_ENV_ID            | apim.envId                 |
 * | GRAVITEE_USERNAME          | apim.auth.username         |
 * | GRAVITEE_PASSWORD          | apim.auth.password         |
 * | GRAVITEE_GATEWAY_URL       | gateway.baseUrl            |
 * | GRAVITEE_GATEWAY_MTLS_URL  | gateway.mtlsBaseUrl        |
 */
export function applyEnvVars(config: GraviteeTestConfig): GraviteeTestConfig {
  const gatewayUrl = process.env["GRAVITEE_GATEWAY_URL"] ?? config.gateway?.baseUrl;
  const gatewayMtlsUrl = process.env["GRAVITEE_GATEWAY_MTLS_URL"] ?? config.gateway?.mtlsBaseUrl;
  return {
    apim: {
      ...config.apim,
      baseUrl: process.env["GRAVITEE_BASE_URL"] ?? config.apim.baseUrl,
      envId: process.env["GRAVITEE_ENV_ID"] ?? config.apim.envId,
      auth: {
        username: process.env["GRAVITEE_USERNAME"] ?? config.apim.auth.username,
        password: process.env["GRAVITEE_PASSWORD"] ?? config.apim.auth.password,
      },
    },
    gateway:
      gatewayUrl !== undefined || gatewayMtlsUrl !== undefined
        ? { baseUrl: gatewayUrl, mtlsBaseUrl: gatewayMtlsUrl }
        : config.gateway,
  };
}

/**
 * Load and parse a config.yaml config file.
 *
 * @param configPath - explicit path; defaults to `config.yaml` in CWD
 */
export async function loadGraviteeConfig(configPath?: string): Promise<GraviteeTestConfig> {
  const filePath = configPath ?? join(process.cwd(), "config.yaml");
  let content: string;
  try {
    content = await readFile(filePath, "utf-8");
  } catch (err) {
    throw new Error(
      `Cannot read config file "${filePath}": ${err instanceof Error ? err.message : String(err)}`,
    );
  }
  const raw = YAML.parse(content);
  const config = validateConfig(raw);
  return applyEnvVars(config);
}

/**
 * Create a Mapi instance from a GraviteeTestConfig.
 * Maps the simple auth shape (username/password) to the full MapiConfig.
 */
export function createMapiFromConfig(config: GraviteeTestConfig): Mapi {
  const mapiConfig: MapiConfig = {
    baseUrl: config.apim.baseUrl,
    envId: config.apim.envId,
    auth: {
      type: "basic",
      username: config.apim.auth.username,
      password: config.apim.auth.password,
    },
  };
  return new Mapi(mapiConfig);
}
