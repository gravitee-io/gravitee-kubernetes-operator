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

import type { Provisioned, Role } from "./types.js";

/**
 * What a provisioner's OWN record says about a role. Deliberately three
 * outcomes and no fourth: there is no "unknown" state. A provisioner that
 * cannot yet tell what happened MUST throw rather than return an ambiguous
 * result — an "unknown" outcome would let a half-implemented `read()` pass
 * silently instead of failing loudly.
 */
export type ProvisionerState = "applied" | "failed" | "gone";

export interface ProvisionerViewResult<Detail = unknown> {
  readonly state: ProvisionerState;
  /** Provisioner-specific evidence for `state` (a condition, a plan diff, an HTTP error body, ...). */
  readonly detail: Detail;
}

/**
 * Provisioner-internal, agnostic readout: "did MY layer land this role?" — as
 * opposed to `mapi`, which answers "is the resource actually right?" and is
 * identical across provisioners. Reached via `Provisioned.view`, usable from
 * shared scenario bodies without narrowing to a concrete provisioner (unlike
 * `Provisioned.checks`, which requires `isGko`/`isTerraform`/... first).
 */
export interface ProvisionerView<Detail = unknown> {
  /** MUST throw if it cannot yet tell what happened — never returns an "unknown" state. */
  read(role: Role): Promise<ProvisionerViewResult<Detail>>;
}

/** The provisioner's record for `role` matches `expected`, or throws with the observed detail. */
export async function assertProvisioner(
  provisioned: Provisioned<unknown>,
  role: Role,
  expected: ProvisionerState,
): Promise<void> {
  const result = await provisioned.view.read(role);
  if (result.state !== expected) {
    throw new Error(
      `expected provisioner state "${expected}" for role "${role}" ` +
        `(${provisioned.provisionerId}), got "${result.state}": ${JSON.stringify(result.detail)}`,
    );
  }
}
