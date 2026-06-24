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
 * Shared definition for the upgrade SURVIVAL check, used by both phases:
 *
 *   survival.before.spec.ts  - provisions this scenario on the OLD line
 *                              (released GKO 4.11 + APIM 4.11) and leaves it.
 *   survival.after.spec.ts   - re-attaches to it on the NEW line (branch GKO +
 *                              APIM 4.12), in a SEPARATE process.
 *
 * The two phases share state only through the cluster (the CRs persist across the
 * in-place upgrade), so the resource identities MUST be fixed and identical on
 * both sides - hence the single source of truth here. Reuses the suite's `dev-ctx`
 * ManagementContext (applied by global-setup), so no context fixture is needed.
 *
 * Core resource set (a JWT-plan V4 API that works on the 4.11 line): a published
 * API, an application, and a subscription. mTLS + cross-subscription expansion
 * (runbook phases 4-5) layers on top of this once the core is proven end-to-end.
 */

import { gkoScenario } from "../../helpers/provisioner-env.js";
import { fixture } from "../../setup.js";

/** Fixed identities + expectations shared by the before/after phases. */
export const SURVIVAL = {
  apiName: "legacy-api",
  appName: "legacy-app",
  subName: "legacy-legacy-jwt",
  contextPath: "/legacy",
  updatedDescription: "Legacy API - JWT plan (updated after upgrade)",
} as const;

/**
 * Factory for the survival scenario's GKO provisioner. Both phases call this to
 * get a provisioner bound to the same CRs: the before-phase calls `provision()`,
 * the after-phase calls `attach()` (rebuilding the handle from the still-present
 * CRs). Static manifests (JWT secret, API, application) are applied in order; the
 * subscription is the dynamic role, applied after the API/app reconcile.
 */
export const survivalScenario = gkoScenario({
  manifests: ["upgrade/jwt-secret.yaml", "upgrade/api-legacy.yaml", "upgrade/app-legacy.yaml"],
  roles: {
    api: { kind: "apiv4definition", name: SURVIVAL.apiName },
    application: { kind: "application", name: SURVIVAL.appName },
    subscription: { kind: "subscription", name: SURVIVAL.subName },
  },
  dynamicRoles: ["subscription"],
  contextPath: SURVIVAL.contextPath,
  applyParams: async (kubectl) => {
    await kubectl.apply(fixture("upgrade/sub-legacy-legacy-jwt.yaml"));
  },
});

/** Fixed identity for the V2 keyless survival API (GKO-1060). */
export const SURVIVAL_V2 = {
  apiName: "legacy-v2-api",
  contextPath: "/legacy-v2",
} as const;

/**
 * GKO provisioner for the V2 keyless survival API (kind `apidefinition`). A
 * keyless plan means the gateway serves it without a token, so "the V2 API
 * survives and stays reachable" is a plain 200 before and after the upgrade.
 */
export const survivalV2Scenario = gkoScenario({
  manifests: ["upgrade/v2-legacy-api.yaml"],
  roles: {
    api: { kind: "apidefinition", name: SURVIVAL_V2.apiName },
  },
  contextPath: SURVIVAL_V2.contextPath,
});

/** Fixed identities for the 4.12-only V2-subscription check. */
export const SURVIVAL_V2_SUB = {
  apiName: "v2-sub-api",
  appName: "v2-sub-app",
  subName: "v2-sub",
} as const;

/**
 * GKO provisioner for a V2 API + app + JWT-plan subscription. V2 subscriptions go
 * through the Automation API (4.12-only), so this is provisioned only in the
 * after-phase, gated on the upgrade target. provision() waits for the Subscription
 * CR to reach Accepted, which is the signal that the 4.12 V2-subscription path works.
 */
export const survivalV2SubScenario = gkoScenario({
  manifests: ["upgrade/v2-sub-api.yaml", "upgrade/v2-sub-app.yaml"],
  roles: {
    api: { kind: "apidefinition", name: SURVIVAL_V2_SUB.apiName },
    application: { kind: "application", name: SURVIVAL_V2_SUB.appName },
    subscription: { kind: "subscription", name: SURVIVAL_V2_SUB.subName },
  },
  dynamicRoles: ["subscription"],
  applyParams: async (kubectl) => {
    await kubectl.apply(fixture("upgrade/v2-sub.yaml"));
  },
});
