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

// ── Core abstraction ──────────────────────────────────────────
export type {
  ProvisionerId,
  Role,
  ResourceRef,
  Provisioned,
  Provisioner,
  ProvisionerChecks,
} from "./types.js";

// ── GKO ───────────────────────────────────────────────────────
export { GkoProvisioner } from "./gko/gko-provisioner.js";
export type {
  GkoScenarioSpec,
  GkoRoles,
  GkoRoleBinding,
  KubectlEngine,
} from "./gko/gko-provisioner.js";
export { isGko } from "./gko/checks.js";
export type { GkoChecks } from "./gko/checks.js";
export { subscriptionYaml, apiKeySecretYaml } from "./gko/subscription-yaml.js";
export type { SubscriptionYamlOptions, ApiKeyEntry } from "./gko/subscription-yaml.js";

// ── Terraform ─────────────────────────────────────────────────
export { TerraformProvisioner } from "./terraform/terraform-provisioner.js";
export type { TfScenarioSpec } from "./terraform/terraform-provisioner.js";
export { isTerraform } from "./terraform/checks.js";
export type { TerraformChecks } from "./terraform/checks.js";
