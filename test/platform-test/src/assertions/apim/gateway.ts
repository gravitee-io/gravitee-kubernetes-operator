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

import { AssertionError } from "node:assert";
import { poll } from "../../utils/match/poll.js";
import type { FetchFn } from "../../types/http.js";
import type { PollOptions } from "../../types/match.js";
import type { GatewayConfig, GatewayRespondOptions, GatewayNotRespondOptions } from "../../types/gateway.js";

/**
 * Gateway (data-plane) assertion client.
 *
 * Assertions retry until the condition is met or the timeout is reached,
 * accommodating the eventual consistency of operator reconciliation and
 * gateway sync.
 *
 * Instantiate via `mapi.gateway(config)` rather than directly.
 *
 * @example
 * const gateway = mapi.gateway({ baseUrl: "http://localhost:8082" });
 *
 * // Keyless API — expect 200
 * await gateway.assertResponds("/my-api", { status: 200 });
 *
 * // JWT-secured API — expect 401 without a token
 * await gateway.assertResponds("/jwt-demo", { status: 401 });
 *
 * // JWT-secured API — expect 200 with a valid token
 * await gateway.assertResponds("/jwt-demo", {
 *   status: 200,
 *   headers: { Authorization: "Bearer <token>" },
 * });
 *
 * // After subscription removed — expect anything except 200
 * await gateway.assertNotResponds("/jwt-demo", { notStatus: 200 });
 *
 * // mTLS: supply a fetchFn built with createTlsFetch({ cert, key, ca })
 * import { createTlsFetch } from "@gravitee/platform-test/utils/http";
 * const mtlsFetch = createTlsFetch({ cert, key, ca });
 * const gateway = mapi.gateway({ baseUrl: "https://localhost:8443" }, mtlsFetch);
 * await gateway.assertResponds("/mtls-demo", { status: 200 });
 */
export class Gateway {
  private readonly baseUrl: string;
  private readonly pollOpts: PollOptions;
  private readonly fetchImpl: FetchFn;

  constructor(config: GatewayConfig, fetchFn?: FetchFn) {
    this.baseUrl = config.baseUrl.replace(/\/+$/, "");
    this.pollOpts = {
      intervalMs: config.retryIntervalMs ?? 500,
      timeoutMs: config.maxRetryMs ?? 30_000,
    };
    this.fetchImpl = fetchFn ?? globalThis.fetch;
  }

  /**
   * Assert that GET <path> returns exactly the expected HTTP status.
   * Retries until the status matches or the timeout expires.
   * @throws AssertionError on final mismatch
   */
  async assertResponds(
    path: string,
    options: GatewayRespondOptions,
    description?: string,
  ): Promise<void> {
    const url = this.buildUrl(path);
    const { status: expected, headers = {} } = options;
    const label = description ?? `GET ${url} → ${expected}`;

    await poll(
      async () => {
        const actual = await this.get(url, headers);
        if (actual !== expected) {
          throw new AssertionError({
            message: `${label}: actual status ${actual} !== expected ${expected}`,
            actual,
            expected,
            operator: "assertGatewayResponds",
          });
        }
      },
      {
        ...this.pollOpts,
        description: label,
      },
    );
  }

  /**
   * Assert that GET <path> does NOT return a specific HTTP status.
   * Retries until the status differs or the timeout expires.
   *
   * Useful when a subscription is removed and the gateway should
   * stop returning 200, but the exact rejection code is unknown.
   *
   * @throws AssertionError if notStatus is returned throughout the window
   */
  async assertNotResponds(
    path: string,
    options: GatewayNotRespondOptions,
    description?: string,
  ): Promise<void> {
    const url = this.buildUrl(path);
    const { notStatus, headers = {} } = options;

    await poll(
      async () => {
        let actual: number;
        try {
          actual = await this.get(url, headers);
        } catch {
          // Connection / TLS errors count as "not responding with notStatus"
          // (e.g. mTLS handshake rejected → no HTTP status at all).
          return;
        }
        if (actual === notStatus) {
          throw new AssertionError({
            message: `Gateway assertion failed: GET ${url} → ${actual}, expected anything except ${notStatus}`,
            actual,
            expected: `!= ${notStatus}`,
            operator: "assertGatewayNotResponds",
          });
        }
      },
      {
        ...this.pollOpts,
        description: description ?? `GET ${url} → not ${notStatus}`,
      },
    );
  }

  private buildUrl(path: string): string {
    return `${this.baseUrl}${path.startsWith("/") ? path : `/${path}`}`;
  }

  private async get(url: string, headers: Record<string, string>): Promise<number> {
    const response = await this.fetchImpl(url, { method: "GET", headers });
    return response.status;
  }
}
