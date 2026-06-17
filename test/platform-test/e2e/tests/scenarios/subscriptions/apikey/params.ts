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
 * Parameter binding for the `subscriptions/apikey` pilot. ONE shared param type
 * (`ApiKeyParams`) drives every provisioner; the small per-provisioner `apply`
 * closures translate it into the GKO Subscription YAML and the Terraform tfvars,
 * reusing the builders the original tests already had. This is the
 * scenario-specific seam the provisioner core deliberately does not own.
 */

import {
  subscriptionYaml,
  type ApiKeyEntry,
  type KubectlEngine,
} from "../../../../../src/provisioners/index.js";

/** What every apikey scenario can be parameterized by. */
export interface ApiKeyParams {
  /** Desired key set. Empty/undefined lets APIM auto-generate a single key. */
  keys?: ApiKeyEntry[];
}

/**
 * Per-process suffix so re-runs pick fresh api-key values. APIM enforces api-key
 * value uniqueness per API across active AND revoked states, and the local
 * MongoDB persists across cluster lifecycle, so hardcoded values would collide
 * with "API Key already exists" on a second run. Date.now alone collides if two
 * processes start in the same millisecond, so mix in a short random suffix.
 */
export const RUN_ID = `${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 6)}`;

/** Generate a unique api-key value padded to at least 32 chars. */
export function uniqueKey(prefix: string): string {
  return `${prefix}-${RUN_ID}`.padEnd(32, "0");
}

/** Standard GKO base resources for the single-api-key-plan apikey scenarios. */
export const APIKEY_GKO = {
  apiName: "e2e-v4-apikey-plan",
  appName: "e2e-app-simple",
  plan: "ApiKey",
  contextPath: "/e2e-v4-apikey-plan",
  apiManifest: "plans/v4-apikey/crd.yaml",
  appManifest: "applications/application-simple/crd.yaml",
} as const;

/**
 * GKO apply closure for a given Subscription name: (re)builds and applies the
 * Subscription manifest for the requested key set. Empty keys omit the apiKeys
 * block so APIM auto-generates one.
 */
export function gkoApplyApiKeys(subName: string) {
  return async (kubectl: KubectlEngine, params: ApiKeyParams): Promise<void> => {
    await kubectl.applyString(
      subscriptionYaml({
        name: subName,
        apiName: APIKEY_GKO.apiName,
        plan: APIKEY_GKO.plan,
        applicationName: APIKEY_GKO.appName,
        keys: params.keys ?? [],
      }),
    );
  };
}

/**
 * Terraform tfvars closure for a given hrid suffix: maps the shared params to
 * the `keys` + `hrid_suffix` variables the apikey-custom fixture expects.
 * `expire_at` is emitted as null (not omitted) so the JSON shape stays stable
 * across applies, which the provider's optional(string) attribute treats as
 * "no expireAt".
 */
export function tfApiKeyVars(hridSuffix: string) {
  return (params: ApiKeyParams): Record<string, unknown> => ({
    hrid_suffix: hridSuffix,
    keys: (params.keys ?? []).map((k) => ({ key: k.key, expire_at: k.expireAt ?? null })),
  });
}
