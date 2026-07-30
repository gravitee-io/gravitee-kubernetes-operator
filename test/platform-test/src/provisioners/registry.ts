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
 * Single source of truth for the set of provisioners the suite knows about.
 * Everything that needs to enumerate provisioners (the `ProvisionerId` type,
 * `forEachProvisioner`'s test-generation order, the `--provision-with` CLI
 * flag, and Playwright's per-lane file/tag filtering) derives from this
 * module instead of repeating its own hardcoded list.
 */

/** The provisioners the suite can run a scenario through, in generation order. */
export const PROVISIONER_ORDER = ["gko", "terraform"] as const;

/** Which provisioner created a resource. */
export type ProvisionerId = (typeof PROVISIONER_ORDER)[number];

/** Per-provisioner metadata needed to select a Playwright "lane". */
export interface ProvisionerLane {
  readonly id: ProvisionerId;
  /**
   * Legacy single-provisioner test folder (`e2e/tests/<segment>/`) to ignore
   * when running a DIFFERENT provisioner's lane. Absent when the provisioner
   * has no such legacy folder (e.g. a provisioner that only ever appears via
   * shared `*.scenario.ts` files).
   */
  readonly testDirSegment?: string;
  /** Case-sensitive title tag appended to generated test titles, e.g. "@gko". */
  readonly tag: string;
}

export const PROVISIONER_LANES: readonly ProvisionerLane[] = [
  { id: "gko", testDirSegment: "gko", tag: "@gko" },
  { id: "terraform", testDirSegment: "terraform", tag: "@terraform" },
];
