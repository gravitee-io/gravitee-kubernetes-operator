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

import type { ProvisionerChecks } from "../types.js";

/**
 * Terraform-specific assertions, reachable from a provisioned handle via
 * `provision.checks` once narrowed with {@link isTerraform}. They check
 * provider-level behaviour with no shared-layer home: declarative-diff drift,
 * idempotency, and Sensitive-value redaction. There is no GKO equivalent
 * (k8s reconciliation is push-based, not declarative-diff-based).
 */
export interface TerraformChecks extends ProvisionerChecks {
  readonly provisionerId: "terraform";

  /** Assert `terraform plan` reports no drift immediately after apply. */
  assertNoDrift(): Promise<void>;

  /** Assert a re-apply of the unchanged spec is a 0-added/0-changed/0-destroyed no-op. */
  assertReapplyNoop(): Promise<void>;

  /** Assert a value never appears in `terraform plan` stdout (Sensitive redaction). */
  assertRedactedInPlan(value: string): Promise<void>;

  /** Raw `terraform plan` result, for bespoke drift/idempotency assertions. */
  plan(): Promise<{ stdout: string; hasChanges: boolean }>;
}

/** Narrow a provisioner-checks surface to the Terraform one. */
export function isTerraform(checks: ProvisionerChecks): checks is TerraformChecks {
  return checks.provisionerId === "terraform";
}
