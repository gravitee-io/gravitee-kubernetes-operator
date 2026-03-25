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
import { HttpClient } from "../../utils/http/http.js";
import { deepPartialMatch } from "../../utils/match/partial.js";
import { throwIfFailed } from "../../utils/match/result.js";
import type { FetchFn } from "../../types/http.js";
import type { DeepPartial, AssertionReport } from "../../types/match.js";
import type { MapiConfig } from "../../types/mapi.js";
import type { GatewayConfig } from "../../types/gateway.js";
import type { Api, Application, Plan, PaginatedResult, Subscription } from "../../types/apim.js";
import { Gateway } from "./gateway.js";

/**
 * mAPI assertion client.
 *
 * Fetches resources from the Gravitee v2 Management API and asserts
 * they match an expected partial shape.
 *
 * @example
 * const mapi = createMapi({ baseUrl: "http://localhost:8083", auth: { type: "basic", username: "admin", password: "admin" } });
 * await mapi.assertApiMatches(apiId, { name: "My API", state: "STARTED" });
 */
export class Mapi {
  /** @internal */
  readonly http: HttpClient;

  constructor(config: MapiConfig, fetchFn?: FetchFn) {
    this.http = new HttpClient(config, fetchFn);
  }

  // ── API Assertions ──────────────────────────────────────────

  /**
   * Fetch an API by ID and assert it partially matches the expected shape.
   * @throws AssertionError if any specified field doesn't match
   */
  async assertApiMatches(apiId: string, expected: DeepPartial<Api>): Promise<void> {
    const api = await this.fetchApi(apiId);
    const report = deepPartialMatch(api, expected);
    throwIfFailed(report);
  }

  /** Non-throwing variant — returns the report for soft assertions. */
  async checkApiMatches(apiId: string, expected: DeepPartial<Api>): Promise<AssertionReport> {
    const api = await this.fetchApi(apiId);
    return deepPartialMatch(api, expected);
  }

  async assertApiState(apiId: string, state: Api["state"]): Promise<void> {
    return this.assertApiMatches(apiId, { state } as DeepPartial<Api>);
  }

  async assertApiStarted(apiId: string): Promise<void> {
    return this.assertApiState(apiId, "STARTED");
  }

  async assertApiStopped(apiId: string): Promise<void> {
    return this.assertApiState(apiId, "STOPPED");
  }

  /**
   * Assert that the management API returns a specific HTTP status code for the given API ID.
   *
   * Useful for asserting absence (e.g. 404 after deletion) without fetching the resource body.
   * @throws AssertionError if the actual status does not match expectedStatus
   */
  async assertApiHttpStatus(apiId: string, expectedStatus: number): Promise<void> {
    const path = this.http.managementV2Path(`/apis/${apiId}`);
    const res = await this.http.get<unknown>(path);
    if (res.status !== expectedStatus) {
      throw new AssertionError({
        message: `Expected HTTP ${expectedStatus} but got ${res.status} for API ${apiId}`,
        expected: expectedStatus,
        actual: res.status,
        operator: "assertApiStatus",
      });
    }
  }

  // ── Plan Assertions ─────────────────────────────────────────

  async assertPlanMatches(apiId: string, planId: string, expected: DeepPartial<Plan>): Promise<void> {
    const plan = await this.fetchPlan(apiId, planId);
    throwIfFailed(deepPartialMatch(plan, expected));
  }

  async assertPlanStatus(apiId: string, planId: string, status: Plan["status"]): Promise<void> {
    return this.assertPlanMatches(apiId, planId, { status } as DeepPartial<Plan>);
  }

  async assertPlanPublished(apiId: string, planId: string): Promise<void> {
    return this.assertPlanStatus(apiId, planId, "PUBLISHED");
  }

  // ── Subscription Assertions ─────────────────────────────────

  async assertSubscriptionMatches(
    apiId: string,
    subscriptionId: string,
    expected: DeepPartial<Subscription>,
  ): Promise<void> {
    const sub = await this.fetchSubscription(apiId, subscriptionId);
    throwIfFailed(deepPartialMatch(sub, expected));
  }

  async assertSubscriptionStatus(
    apiId: string,
    subscriptionId: string,
    status: Subscription["status"],
  ): Promise<void> {
    return this.assertSubscriptionMatches(apiId, subscriptionId, { status } as DeepPartial<Subscription>);
  }

  async assertSubscriptionAccepted(apiId: string, subscriptionId: string): Promise<void> {
    return this.assertSubscriptionStatus(apiId, subscriptionId, "ACCEPTED");
  }

  // ── Gateway (data plane) ────────────────────────────────────

  /**
   * Create a Gateway for the mAPI data-plane.
   *
   * The gateway URL is separate from the management API URL — typically
   * localhost:8082 (HTTP) or localhost:8443 (HTTPS/mTLS).
   *
   * For mTLS, pass a `fetchFn` created with `createTlsFetch({ cert, key, ca })`.
   *
   * @example
   * const gateway = mapi.gateway({ baseUrl: "http://localhost:8082" });
   * await gateway.assertResponds("/my-api", { status: 200 });
   */
  gateway(config: GatewayConfig, fetchFn?: FetchFn): Gateway {
    return new Gateway(config, fetchFn);
  }

  // ── Application Assertions ──────────────────────────────────

  async assertApplicationMatches(appId: string, expected: DeepPartial<Application>): Promise<void> {
    const app = await this.fetchApplication(appId);
    throwIfFailed(deepPartialMatch(app, expected));
  }

  async assertApplicationHttpStatus(appId: string, expectedStatus: number): Promise<void> {
    const path = this.http.managementV1Path(`/applications/${appId}`);
    const res = await this.http.get<unknown>(path);
    if (res.status !== expectedStatus) {
      throw new AssertionError({
        message: `Expected HTTP ${expectedStatus} but got ${res.status} for Application ${appId}`,
        expected: expectedStatus,
        actual: res.status,
        operator: "assertApplicationHttpStatus",
      });
    }
  }

  // ── Plan List Helpers ──────────────────────────────────────

  async listApiPlans(apiId: string): Promise<Plan[]> {
    const path = this.http.managementV2Path(`/apis/${apiId}/plans?page=1&perPage=100`);
    const res = await this.http.get<PaginatedResult<Plan>>(path);
    if (res.status !== 200) {
      throw new Error(`Failed to list plans for API ${apiId}: ${res.status}`);
    }
    return res.body.data;
  }

  // ── Fetch Helpers ───────────────────────────────────────────

  async fetchApplication(appId: string): Promise<Application> {
    const path = this.http.managementV1Path(`/applications/${appId}`);
    const res = await this.http.get<Application>(path);
    if (res.status !== 200) {
      throw new Error(`Failed to fetch Application ${appId}: ${res.status} ${res.statusText}\n${JSON.stringify(res.body, null, 2)}`);
    }
    return res.body;
  }

  async fetchApi(apiId: string): Promise<Api> {
    const path = this.http.managementV2Path(`/apis/${apiId}`);
    const res = await this.http.get<Api>(path);
    if (res.status !== 200) {
      throw new Error(`Failed to fetch API ${apiId}: ${res.status} ${res.statusText}\n${JSON.stringify(res.body, null, 2)}`);
    }
    return res.body;
  }

  async fetchPlan(apiId: string, planId: string): Promise<Plan> {
    const path = this.http.managementV2Path(`/apis/${apiId}/plans/${planId}`);
    const res = await this.http.get<Plan>(path);
    if (res.status !== 200) {
      throw new Error(`Failed to fetch Plan ${planId}: ${res.status} ${res.statusText}\n${JSON.stringify(res.body, null, 2)}`);
    }
    return res.body;
  }

  async fetchSubscription(apiId: string, subscriptionId: string): Promise<Subscription> {
    const path = this.http.managementV2Path(`/apis/${apiId}/subscriptions/${subscriptionId}`);
    const res = await this.http.get<Subscription>(path);
    if (res.status !== 200) {
      throw new Error(`Failed to fetch Subscription ${subscriptionId}: ${res.status} ${res.statusText}\n${JSON.stringify(res.body, null, 2)}`);
    }
    return res.body;
  }
}

/** Convenience factory */
export function createMapi(config: MapiConfig, fetchFn?: FetchFn): Mapi {
  return new Mapi(config, fetchFn);
}
