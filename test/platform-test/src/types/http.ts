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
 * Configuration for the HTTP client.
 */
export interface HttpClientConfig {
  /** Base URL of the Gravitee management API */
  baseUrl: string;
  /** Environment ID (defaults to "DEFAULT") */
  envId?: string;
  /** Authentication configuration */
  auth:
    | { type: "basic"; username: string; password: string }
    | { type: "bearer"; token: string }
    | { type: "cookie"; cookieName: string; cookieValue: string };
  /** Request timeout in milliseconds (defaults to 10_000) */
  timeoutMs?: number;
  /** Custom headers for every request */
  headers?: Record<string, string>;
}

/**
 * Typed HTTP response wrapper.
 */
export interface HttpResponse<T = unknown> {
  status: number;
  statusText: string;
  headers: Headers;
  body: T;
  elapsedMs: number;
  raw: Response;
}

/**
 * Abstraction over the fetch function signature.
 * Allows swapping native fetch for undici (e.g. for TLS/cert scenarios).
 */
export type FetchFn = typeof globalThis.fetch;

export interface TlsOptions {
  /** Client certificate PEM (for mTLS) */
  cert?: Buffer | string;
  /** Client private key PEM (for mTLS) */
  key?: Buffer | string;
  /** CA certificate PEM for server verification */
  ca?: Buffer | string;
  /**
   * Whether to reject self-signed / unverified server certificates.
   * Defaults to false (permissive) so self-signed CI certs work out of the box.
   * Automatically set to true when a CA cert is provided.
   */
  rejectUnauthorized?: boolean;
}
