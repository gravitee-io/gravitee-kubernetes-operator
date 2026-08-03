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

import { poll } from "../utils/match/poll.js";
import type { PollOptions } from "../types/match.js";
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
 * `Provisioned.checks`, which requires `isTerraform`/... first).
 */
export interface ProvisionerView<Detail = unknown> {
  /** MUST throw if it cannot yet tell what happened — never returns an "unknown" state. */
  read(role: Role): Promise<ProvisionerViewResult<Detail>>;
}

/**
 * The provisioner's record for `role` reaches `expected`, or throws with the observed detail.
 *
 * Polls rather than reading once: a provisioner's record settles asynchronously (GKO writes
 * `.status` after the API server accepts the manifest), and `read()` is specified to THROW
 * while it cannot yet tell. Both the transient throw and a not-yet-`expected` state are
 * retried.
 *
 * Where to use it: after `remove()` (expect `"gone"`), and on failure paths (expect
 * `"failed"`, whose detail carries the provisioner's own reason).
 *
 * Where NOT to use it: as a convergence wait after `update()`. GKO stamps
 * `observedGeneration` on its conditions but that value can lag indefinitely after a
 * re-apply (GKO-2940), so a post-update read can return the PRE-update `Accepted=True`
 * and pass without asserting anything. Gating on `observedGeneration` instead would
 * reintroduce the same GKO-2940 hang that `kubectl wait --for=condition` suffers. Until
 * GKO-2940 is fixed, `mapi` is the convergence signal after an update; this answers the
 * different question of whether the provisioner still holds a record at all.
 */
export async function assertProvisioner(
  provisioned: Provisioned<unknown>,
  role: Role,
  expected: ProvisionerState,
  options: PollOptions = {},
): Promise<void> {
  await poll(
    async () => {
      const result = await provisioned.view.read(role);
      if (result.state !== expected) {
        throw new Error(
          `expected provisioner state "${expected}" for role "${role}" ` +
            `(${provisioned.provisionerId}), got "${result.state}": ${JSON.stringify(result.detail)}`,
        );
      }
    },
    {
      description: `${provisioned.provisionerId} record for role "${role}" to be "${expected}"`,
      ...options,
    },
  );
}
