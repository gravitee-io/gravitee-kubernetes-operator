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

import type { ProvisionerChecks, Role } from "../types.js";

/**
 * GKO-specific assertions, reachable from a provisioned handle via
 * `provision.checks` once narrowed with {@link isGko}. These wrap the kubectl
 * primitives, bound to the handle's roles, so a test does not re-pass
 * (kind, name): they check Kubernetes control-plane state (conditions, events,
 * `.status`) that only the GKO path exposes.
 */
export interface GkoChecks extends ProvisionerChecks {
  readonly provisionerId: "gko";

  /** Wait for a condition (default "Accepted") on the CR playing `role`. */
  waitForCondition(role: Role, condition?: string, timeoutSeconds?: number): Promise<void>;

  /** Assert a Kubernetes event for the role's CR contains a message substring. */
  assertEventContains(role: Role, message: string): Promise<void>;

  /** Read the `.status` of the role's CR. */
  status<T = unknown>(role: Role): Promise<T>;

  /** Apply a manifest string expecting admission rejection; returns stderr. */
  applyExpectFailure(yaml: string): Promise<string>;
}

/** Narrow a provisioner-checks surface to the GKO one. */
export function isGko(checks: ProvisionerChecks): checks is GkoChecks {
  return checks.provisionerId === "gko";
}
