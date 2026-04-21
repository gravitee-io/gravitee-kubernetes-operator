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

import type { HttpClientConfig, HttpResponse } from "../../types/http.js";

/**
 * Internal HTTP client for Gravitee management APIs.
 *
 * Uses native fetch. The management API is always reached over plain HTTP
 * with basic/bearer/cookie auth in our test scenarios — mTLS only applies
 * to the data plane, which has its own injectable fetch on `Gateway`.
 */
export class HttpClient {
  private readonly baseUrl: string;
  private readonly envId: string;
  private readonly timeoutMs: number;
  private readonly defaultHeaders: Record<string, string>;
  private readonly authHeader: Record<string, string>;

  constructor(config: HttpClientConfig) {
    this.baseUrl = config.baseUrl.replace(/\/+$/, "");
    this.envId = config.envId ?? "DEFAULT";
    this.timeoutMs = config.timeoutMs ?? 10_000;
    this.defaultHeaders = config.headers ?? {};

    // Pre-compute auth header
    const { auth } = config;
    if (auth.type === "basic") {
      this.authHeader = {
        Authorization: `Basic ${btoa(`${auth.username}:${auth.password}`)}`,
      };
    } else if (auth.type === "bearer") {
      this.authHeader = { Authorization: `Bearer ${auth.token}` };
    } else {
      this.authHeader = { Cookie: `${auth.cookieName}=${auth.cookieValue}` };
    }
  }

  async get<T = unknown>(path: string): Promise<HttpResponse<T>> {
    return this.request<T>("GET", path);
  }

  async post<T = unknown>(path: string, body?: unknown): Promise<HttpResponse<T>> {
    return this.request<T>("POST", path, body);
  }

  async put<T = unknown>(path: string, body?: unknown): Promise<HttpResponse<T>> {
    return this.request<T>("PUT", path, body);
  }

  async delete<T = unknown>(path: string): Promise<HttpResponse<T>> {
    return this.request<T>("DELETE", path);
  }

  /** Build v1 management API path for the configured environment */
  managementV1Path(resource: string): string {
    return `/management/organizations/DEFAULT/environments/${this.envId}${resource}`;
  }

  /** Build v2 management API path for the configured environment */
  managementV2Path(resource: string): string {
    return `/management/v2/environments/${this.envId}${resource}`;
  }

  private async request<T>(method: string, path: string, body?: unknown): Promise<HttpResponse<T>> {
    const url = `${this.baseUrl}${path}`;

    const headers: Record<string, string> = {
      Accept: "application/json",
      ...this.defaultHeaders,
      ...this.authHeader,
    };

    if (body !== undefined) {
      headers["Content-Type"] = "application/json";
    }

    const start = performance.now();

    const response = await fetch(url, {
      method,
      headers,
      body: body !== undefined ? JSON.stringify(body) : undefined,
      signal: AbortSignal.timeout(this.timeoutMs),
    });

    const elapsedMs = Math.round(performance.now() - start);

    let responseBody: T;
    const contentType = response.headers.get("content-type") ?? "";
    if (contentType.includes("application/json")) {
      responseBody = (await response.json()) as T;
    } else {
      responseBody = (await response.text()) as unknown as T;
    }

    return {
      status: response.status,
      statusText: response.statusText,
      headers: response.headers,
      body: responseBody,
      elapsedMs,
      raw: response,
    };
  }
}
