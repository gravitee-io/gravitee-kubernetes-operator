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
 * Playwright binding for the provisioner layer. `forProvisioners` expands ONE
 * scenario into ONE tagged test per provisioner. The body is
 * provisioner-agnostic: it gets a `provision` handle plus the shared mapi/gateway
 * fixtures and runs the SAME assertions regardless of how the resource was
 * created. Each test's provisioner tag (e.g. `@gko`/`@terraform`) and its
 * per-provisioner Xray id are appended to the title so `--grep` selection keeps
 * working.
 */

import { test } from "../setup.js";
import type { Mapi, Gateway } from "../../src/index.js";
import type { ProvisionerId, Provisioned, Provisioner } from "../../src/provisioners/index.js";
import { TF_WORKSPACE_TIMEOUT_MS } from "../../src/provisioners/engines/terraform-core.js";

const PROVISIONER_ORDER: readonly ProvisionerId[] = ["gko", "terraform"] as const;

type ProvisionerFactory<P> = () => Provisioner<P> | Promise<Provisioner<P>>;

export interface ScenarioDef<P> {
  /** Human-readable scenario title; the provisioner tag + Xray id are appended. */
  title: string;
  /** Per-provisioner provisioner factories. A provisioner absent here is not generated. */
  provisioners: Partial<Record<ProvisionerId, ProvisionerFactory<P>>>;
  /**
   * Per-provisioner Xray id(s). GKO and TF are different Jira tickets, so each arm
   * carries its own. A list is allowed for the case where one provisioner splits
   * into several tickets what the other does in one (e.g. GKO count + gateway
   * as two ids vs a single combined TF id): all ids land in that arm's title.
   */
  xray: Partial<Record<ProvisionerId, string | string[]>>;
  /** Extra title tags appended to every generated test, e.g. [TAGS.REGRESSION]. */
  tags?: string[];
  /** Provisioners planned but not yet implemented -> a visible skipped entry, never a silent gap. */
  pending?: Partial<Record<ProvisionerId, string>>;
  /** Per-provisioner Playwright timeout override (ms). */
  timeoutMs?: Partial<Record<ProvisionerId, number>>;
}

export interface ScenarioBodyArgs<P> {
  /** Live handle to the provisioned scenario for the current provisioner. */
  provision: Provisioned<P>;
  mapi: Mapi;
  gateway: Gateway;
}

export type ScenarioBody<P> = (args: ScenarioBodyArgs<P>) => Promise<void>;

interface ActiveProvision {
  provisioner: Provisioner<unknown>;
  provision?: Provisioned<unknown>;
}

/**
 * The scenario currently being provisioned, tracked at module scope so a single
 * afterEach can tear it down even when a test TIMES OUT. An inline `finally`
 * does not run on a Playwright timeout, so it alone cannot satisfy the AGENTS.md
 * "always clean up, with a safety net" rule. Safe to keep at module scope
 * because the suite runs serially (workers: 1).
 */
let active: ActiveProvision | undefined;

async function teardownActive(): Promise<void> {
  const current = active;
  active = undefined;
  if (!current) return;
  // destroy()/cleanup() are documented to never throw, but guard anyway so a
  // teardown failure never masks the test's own result.
  if (current.provision) {
    await current.provision.destroy().catch(() => {});
  } else if (current.provisioner.cleanup) {
    await current.provisioner.cleanup().catch(() => {});
  }
}

function buildTitle<P>(scenario: ScenarioDef<P>, provisionerId: ProvisionerId): string {
  const xray = scenario.xray[provisionerId];
  let xrayIds: string[] = [];
  if (Array.isArray(xray)) xrayIds = xray;
  else if (xray) xrayIds = [xray];
  const tokens = [scenario.title, ...xrayIds, `@${provisionerId}`, ...(scenario.tags ?? [])];
  return tokens.filter((t): t is string => Boolean(t)).join(" ");
}

/**
 * Expand one scenario into one tagged Playwright test per provisioner.
 *
 * - A provisioner with a factory runs the shared `body` against a live handle.
 * - A provisioner in `pending` renders as an explicit `test.fixme` (a visible gap,
 *   never a silent skip, never red).
 * - A provisioner absent from both is N/A and emits nothing.
 */
export function forProvisioners<P>(
  scenario: ScenarioDef<P>,
  body: ScenarioBody<P>,
  initialParams: P,
): void {
  // Safety net, file-scoped (this runs during the scenario file's load).
  // Registered per call; redundant hooks are harmless no-ops once `active` is
  // cleared by the first one.
  test.afterEach(teardownActive);

  for (const provisionerId of PROVISIONER_ORDER) {
    const title = buildTitle(scenario, provisionerId);
    const factory = scenario.provisioners[provisionerId];

    if (!factory) {
      const reason = scenario.pending?.[provisionerId];
      if (reason) {
        test(title, async () => {
          test.fixme(true, `pending provisioner (${provisionerId}): ${reason}`);
        });
      }
      continue;
    }

    test(title, async ({ mapi, gateway }) => {
      const defaultTimeout = provisionerId === "terraform" ? TF_WORKSPACE_TIMEOUT_MS : undefined;
      const timeout = scenario.timeoutMs?.[provisionerId] ?? defaultTimeout;
      if (timeout) test.setTimeout(timeout);

      const provisioner = (await factory()) as Provisioner<unknown>;
      const tracked: ActiveProvision = { provisioner };
      active = tracked;
      try {
        const provision = (await provisioner.provision(initialParams)) as Provisioned<P>;
        tracked.provision = provision;
        await body({ provision, mapi, gateway });
      } finally {
        await teardownActive();
      }
    });
  }
}
