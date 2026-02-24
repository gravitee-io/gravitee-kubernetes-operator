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

export interface GatewayConfig {
  /** Base URL of the Gravitee gateway, e.g. http://localhost:8082 */
  baseUrl: string;
  /** Interval between retry attempts in ms (default: 500) */
  retryIntervalMs?: number;
  /** Max total retry window in ms (default: 30_000) */
  maxRetryMs?: number;
}

export interface GatewayRespondOptions {
  /** Expected HTTP status code */
  status: number;
  /** Optional request headers, e.g. { Authorization: "Bearer <token>" } */
  headers?: Record<string, string>;
}

export interface GatewayNotRespondOptions {
  /** HTTP status code that must NOT be returned */
  notStatus: number;
  /** Optional request headers */
  headers?: Record<string, string>;
}
