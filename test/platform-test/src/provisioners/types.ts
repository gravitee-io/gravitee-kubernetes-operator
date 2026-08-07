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
 * Provisioner abstraction.
 *
 * A test defines intent once; a provisioner applies it through one path (GKO,
 * Terraform, later UI, ...); the same shared assertions (`mapi`/`gateway`) run
 * against the resulting platform state. The provisioner hides two things that
 * differ per provisioner: how resources are created, and how a logical role
 * resolves to an APIM id. Add a provisioner by implementing this interface.
 */

import type { ProvisionerId } from "./registry.js";
import type { ProvisionerView } from "./view.js";

/** Which provisioner created a resource. Re-exported for existing importers. */
export type { ProvisionerId };

/**
 * A logical role within a scenario, mapped per-provisioner to a concrete resource.
 * The string union gives autocomplete for the common roles while still allowing
 * scenario-specific keys (e.g. "api:two-plans").
 */
export type Role = "api" | "application" | "subscription" | "plan" | (string & {});

/**
 * Escape hatch for assertions a provisioner alone can make, reached via
 * {@link Provisioned.checks} and narrowed by a type guard (`isTerraform`). The base only
 * carries the discriminant; a concrete shape (`TerraformChecks`) extends it.
 *
 * Deliberately rare. Anything expressible as "did MY layer land this role?" belongs in
 * {@link ProvisionerView}, which shared scenario bodies can call without narrowing. GKO
 * declares no checks: its control-plane readouts are exactly that question, and the
 * remaining Kubernetes primitives (admission rejection, events) are reached through the
 * kubectl engine in tests where nothing is provisioned.
 */
export interface ProvisionerChecks {
  readonly provisionerId: ProvisionerId;
}

/**
 * Live handle to one provisioned scenario instance. Provisioner-agnostic surface a
 * test body interacts with. Methods that hit the cluster/APIM are async.
 */
export interface Provisioned<P = unknown> {
  /** Which provisioner produced this handle. */
  readonly provisionerId: ProvisionerId;

  /**
   * The resource's APIM id (UUID), by kind. Pass an optional `label` only for
   * the rare scenario that has two of the same kind (e.g. `apiId("two-plans")`);
   * omit it for the usual single-resource case. Resolved once, then cached.
   */
  apiId(label?: string): Promise<string>;
  subscriptionId(label?: string): Promise<string>;
  applicationId(label?: string): Promise<string>;
  groupId(label?: string): Promise<string>;
  sharedPolicyGroupId(label?: string): Promise<string>;

  /** The API gateway context path (e.g. "/e2e-...") for data-plane assertions. */
  contextPath(): Promise<string>;

  /** Re-provision with changed parameters (rotation-style tests). */
  update(params: P): Promise<void>;

  /**
   * Remove the resource playing `role`, the way a user would, leaving the rest
   * of the scenario standing: GKO deletes that CR; Terraform drops the resource
   * from the desired state and re-applies. Use this to test partial-teardown
   * effects (e.g. deleting only a Subscription revokes its key while the API
   * stays up) without tearing the whole scenario down.
   */
  remove(role: Role): Promise<void>;

  /** Tear down. MUST be idempotent and safe in finally/afterAll (never throws). */
  destroy(): Promise<void>;

  /**
   * Provisioner-internal, agnostic readout ("did MY layer land this role?").
   * Usable from shared scenario bodies without narrowing, unlike `checks`.
   */
  readonly view: ProvisionerView;

  /**
   * Assertions only this provisioner can make (narrow with `isTerraform`). Absent on
   * provisioners that have none, so shared bodies cannot reach for it by habit.
   */
  readonly checks?: ProvisionerChecks;
}

/**
 * Factory: knows how to stand up ONE scenario for ONE provisioner. A scenario
 * supplies one Provisioner per provisioner it supports.
 */
export interface Provisioner<P = unknown> {
  readonly provisionerId: ProvisionerId;

  /** Stand the scenario up with initial params, returning a live handle. */
  provision(params: P): Promise<Provisioned<P>>;

  /**
   * Best-effort teardown WITHOUT a live handle, for the case where
   * `provision()` itself fails partway (e.g. GKO applied a CR but reconcile got
   * stuck). Idempotent and tolerant; safe to call in a finally block. Omit it
   * when the provisioner self-cleans on a failed provision (Terraform does).
   */
  cleanup?(): Promise<void>;

  /**
   * Seam for upgrade testing (provision -> upgrade cluster -> re-assert):
   * rebuild a handle from stable HRIDs after the original in-memory state is
   * gone. Optional because not every provisioner implements it (GKO does;
   * Terraform does not yet).
   */
  attach?(refs: Partial<Record<Role, { hrid: string }>>): Promise<Provisioned<P>>;
}
