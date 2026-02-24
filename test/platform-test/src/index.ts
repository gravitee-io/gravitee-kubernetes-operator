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

// ── Types (all client-facing types from a single directory) ───
export type {
  // APIM domain models
  Api, ApiV1, ApiV4, ApiV2, ApiFederated, ApiFederatedAgent,
  ApiState, ApiType,
  Plan, PlanV4, PlanV2, PlanFederated, PlanStatus, PlanSecurity, PlanMode, PlanType,
  Subscription, SubscriptionStatus, SubscriptionConsumerStatus, ConsumerStatus,
  // mAPI / Gateway configuration
  MapiConfig,
  GatewayConfig, GatewayRespondOptions, GatewayNotRespondOptions,
  // Matching engine
  DeepPartial, AssertionReport, AssertionFailure, PollOptions,
  // HTTP
  HttpClientConfig, HttpResponse, TlsOptions,
  // Config
  GraviteeTestConfig,
} from "./types/index.js";

// ── mAPI ──────────────────────────────────────────────────────
export { Mapi, createMapi } from "./assertions/apim/index.js";

// ── Gateway ───────────────────────────────────────────────────
export { Gateway } from "./assertions/apim/index.js";

// ── Matching Engine ───────────────────────────────────────────
export { deepPartialMatch, throwIfFailed, poll } from "./utils/match/index.js";

// ── HTTP (for advanced use / custom clients) ──────────────────
export { HttpClient, createTlsFetch } from "./utils/http/index.js";

// ── Config / CLI helpers ──────────────────────────────────────
export { loadGraviteeConfig, createMapiFromConfig, applyEnvVars } from "./cmd/config.js";
