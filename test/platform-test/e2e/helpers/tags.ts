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
 * Xray test ID registry.
 *
 * Use these constants in test titles for filtering with --grep:
 *   test(`should start API ${XRAY.API_LIFECYCLE.START_STOP}`, ...)
 *
 * Run a single Xray test:
 *   npm run e2e -- --grep @GKO-1464
 *
 * Run regression pack:
 *   npm run e2e:regression
 */

export const XRAY = {
  API_LIFECYCLE: {
    DEPLOY_V4_SYNC_K8S: "@GKO-69",
    START_STOP_V2_V4_NATIVE: "@GKO-1464",
  },
  PLANS: {
    KEYLESS_PLAN_V4: "@GKO-110",
  },
  TERRAFORM: {
    APPLY_COMPLEX_CONFIG: "@GKO-1926",
    IDEMPOTENT_CONFIG: "@GKO-1932",
  },
  DEPLOYMENT_RECONCILIATION: {
    RECONCILE_API_CONFIG: "@GKO-1444",
  },
  WEBHOOKS: {
    REJECT_INVALID_CRS: "@GKO-1447",
  },
} as const;

export const TAGS = {
  REGRESSION: "@regression",
} as const;
