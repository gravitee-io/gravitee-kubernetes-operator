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

/** Public surface of the provisioner layer. */

// ── Registry ──────────────────────────────────────────────────
export { PROVISIONER_ORDER, PROVISIONER_LANES } from "./registry.js";
export type { ProvisionerLane } from "./registry.js";

// ── Core abstraction ──────────────────────────────────────────
export type {
  ProvisionerId,
  Role,
  Provisioned,
  Provisioner,
  ProvisionerChecks,
} from "./types.js";
export { assertProvisioner } from "./view.js";
export type { ProvisionerState, ProvisionerViewResult, ProvisionerView } from "./view.js";

// ── GKO ───────────────────────────────────────────────────────
export { GkoProvisioner } from "./gko/gko-provisioner.js";
export type {
  GkoScenarioSpec,
  GkoRoles,
  GkoRoleBinding,
  KubectlEngine,
} from "./gko/gko-provisioner.js";
export type { GkoView, GkoViewDetail } from "./gko/view.js";
export { subscriptionYaml, apiKeySecretYaml } from "./gko/subscription-yaml.js";
export type { SubscriptionYamlOptions, ApiKeyEntry } from "./gko/subscription-yaml.js";

// ── Terraform ─────────────────────────────────────────────────
export { TerraformProvisioner } from "./terraform/terraform-provisioner.js";
export type { TfScenarioSpec } from "./terraform/terraform-provisioner.js";
export { isTerraform } from "./terraform/checks.js";
export type { TerraformChecks } from "./terraform/checks.js";
export type { TerraformView, TerraformViewDetail } from "./terraform/view.js";
