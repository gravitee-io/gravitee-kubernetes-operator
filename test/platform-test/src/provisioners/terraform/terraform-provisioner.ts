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

import * as tfCore from "../engines/terraform-core.js";
import type { TfWorkspace } from "../engines/terraform-core.js";
import type { Provisioned, Provisioner, Role } from "../types.js";
import { BaseProvisioned } from "../base.js";
import type { TerraformChecks } from "./checks.js";
import type { TerraformView } from "./view.js";

/** Convention mapping a role to a default `terraform output` name. */
const DEFAULT_OUTPUT_BY_ROLE: Record<string, string> = {
  api: "api_id",
  subscription: "sub_id",
  application: "app_id",
  group: "group_id",
};
const DEFAULT_CONTEXT_PATH_OUTPUT = "api_context_path";

export interface TfScenarioSpec<P = unknown> {
  /** Absolute path to the fixture directory containing main.tf. */
  fixtureDir: string;
  /** Base environment with APIM_* auth/server vars (built by the e2e adapter). */
  env: Record<string, string>;
  /** role -> terraform output name. Defaults: api->api_id, subscription->sub_id, application->app_id. */
  outputs?: Partial<Record<Role, string>>;
  /**
   * role -> terraform resource address (e.g. "apim_apiv4.test"), for
   * `view.read()`. No default: unlike outputs, resource LOCAL NAMES are not
   * standardized across fixtures (confirmed: existing fixtures use "test",
   * "app", "dict", "group", "spg" ...), so guessing one would silently point
   * `view.read()` at the wrong resource instead of failing loudly. Declare it
   * per scenario for any role that needs `view`/`assertProvisioner`.
   */
  addresses?: Partial<Record<Role, string>>;
  /** Output name for the gateway context path. Default "api_context_path". */
  contextPathOutput?: string;
  /** Map params -> tfvars written before each apply (provision + update). */
  toVars?: (params: P) => Record<string, unknown>;
  /**
   * How to remove a role's resource via a realistic re-apply: tfvars to MERGE
   * over the current vars so the resource drops out of the desired state (e.g.
   * { subscription: { create_subscription: false } } against a count-gated
   * resource). Consumed by Provisioned.remove(role).
   */
  removeVars?: Partial<Record<Role, Record<string, unknown>>>;
}

function outputNameFor<P>(spec: TfScenarioSpec<P>, role: Role): string {
  const name = spec.outputs?.[role] ?? DEFAULT_OUTPUT_BY_ROLE[role];
  if (!name) {
    throw new Error(
      `Terraform scenario has no output mapped for role "${role}". ` +
        `Add it to the scenario's \`outputs\` map (e.g. { ${role}: "${role}_id" }).`,
    );
  }
  return name;
}

function addressFor<P>(spec: TfScenarioSpec<P>, role: Role): string {
  const address = spec.addresses?.[role];
  if (!address) {
    throw new Error(
      `Terraform scenario has no resource address mapped for role "${role}"; ` +
        `cannot determine view state. Add it to the scenario's \`addresses\` map ` +
        `(e.g. { ${role}: "apim_apiv4.test" }).`,
    );
  }
  return address;
}

/**
 * Builds the Terraform `ProvisionerView`: "did MY layer land this role?",
 * answered purely from the workspace's OWN state (`terraform show -json`),
 * never from mAPI. Deliberately reads STATE, not a fresh `plan` — a plan
 * re-reads the platform via the provider, which would make this a platform
 * check (that's `mapi`'s job). MUST throw when it cannot tell (no address
 * declared for the role) rather than return an ambiguous state.
 */
/** @internal exported only for unit testing; not part of the package's public surface. */
export function buildTerraformView<P>(ws: TfWorkspace, spec: TfScenarioSpec<P>): TerraformView {
  return {
    async read(role: Role) {
      const address = addressFor(spec, role);
      const { resources } = await tfCore.showState(ws);
      const resource = resources.find((r) => r.address === address);
      if (!resource) {
        return { state: "gone", detail: { reason: "not in state", address } };
      }
      if (resource.tainted) {
        return { state: "failed", detail: { reason: "tainted", address } };
      }
      return { state: "applied", detail: { reason: "in state, not tainted", address } };
    },
  };
}

function buildTerraformChecks<P>(ws: TfWorkspace, spec: TfScenarioSpec<P>): TerraformChecks {
  return {
    provisionerId: "terraform",
    async taint(role: Role) {
      await tfCore.taint(ws, addressFor(spec, role));
    },
    async untaint(role: Role) {
      await tfCore.untaint(ws, addressFor(spec, role));
    },
    async assertNoDrift() {
      const { hasChanges, stdout } = await tfCore.plan(ws);
      if (hasChanges) {
        throw new Error(`terraform plan reported drift immediately after apply:\n${stdout}`);
      }
    },
    async assertReapplyNoop() {
      const out = await tfCore.apply(ws);
      if (!/0 added.*0 changed.*0 destroyed/.test(out)) {
        throw new Error(`expected a no-op re-apply (0/0/0), got:\n${out}`);
      }
    },
    async assertRedactedInPlan(value: string) {
      const { stdout } = await tfCore.plan(ws);
      if (stdout.includes(value)) {
        throw new Error("terraform plan output leaked a value expected to be redacted as Sensitive");
      }
    },
    async plan() {
      return tfCore.plan(ws);
    },
  };
}

class TerraformProvisioned<P> extends BaseProvisioned<P> {
  readonly provisionerId = "terraform" as const;
  readonly checks: TerraformChecks;
  readonly view: TerraformView;
  private readonly idCache = new Map<string, string>();
  private contextPathCache?: string;

  constructor(
    private readonly ws: TfWorkspace,
    private readonly spec: TfScenarioSpec<P>,
    private lastVars: Record<string, unknown>,
  ) {
    super();
    this.checks = buildTerraformChecks(ws, spec);
    this.view = buildTerraformView(ws, spec);
  }

  protected async resolveId(role: Role): Promise<string> {
    const cached = this.idCache.get(role);
    if (cached) return cached;
    const value = await tfCore.output(this.ws, outputNameFor(this.spec, role));
    this.idCache.set(role, value);
    return value;
  }

  async contextPath(): Promise<string> {
    if (this.contextPathCache !== undefined) return this.contextPathCache;
    const output = this.spec.contextPathOutput ?? DEFAULT_CONTEXT_PATH_OUTPUT;
    this.contextPathCache = await tfCore.output(this.ws, output);
    return this.contextPathCache;
  }

  async update(params: P): Promise<void> {
    if (this.spec.toVars) {
      this.lastVars = this.spec.toVars(params);
      await tfCore.writeVars(this.ws, this.lastVars);
    }
    await tfCore.apply(this.ws);
  }

  async remove(role: Role): Promise<void> {
    const vars = this.spec.removeVars?.[role];
    if (!vars) {
      throw new Error(
        `Terraform scenario has no removeVars for role "${role}". Declare how the ` +
          `role drops out of the desired state, e.g. removeVars: { ${role}: { create_${role}: false } }.`,
      );
    }
    // Re-apply with the resource dropped from the desired state, the way a user
    // would (delete the resource block + apply), leaving the rest standing.
    await tfCore.writeVars(this.ws, { ...this.lastVars, ...vars });
    await tfCore.apply(this.ws);
    this.idCache.delete(role);
  }

  async destroy(): Promise<void> {
    // destroyWorkspace runs `destroy` (swallowing errors) then removes the temp
    // dir, so it is always safe to call, including after an inline destroy.
    await tfCore.destroyWorkspace(this.ws);
  }
}

/**
 * Provisions a scenario through the Terraform APIM provider: copies the fixture
 * to a temp workspace, writes any parameter tfvars, applies, and resolves
 * logical roles to APIM ids from `terraform output`. The TfWorkspace and its
 * tfstate-lock/timeout semantics stay entirely inside this implementation.
 */
export class TerraformProvisioner<P = unknown> implements Provisioner<P> {
  readonly provisionerId = "terraform" as const;

  constructor(private readonly spec: TfScenarioSpec<P>) {}

  async provision(params: P): Promise<Provisioned<P>> {
    const ws = await tfCore.initWorkspace(this.spec.fixtureDir, this.spec.env);
    const vars = this.spec.toVars ? this.spec.toVars(params) : {};
    try {
      if (this.spec.toVars) {
        await tfCore.writeVars(ws, vars);
      }
      await tfCore.apply(ws);
    } catch (err) {
      // Apply failed: tear the workspace down so we don't orphan the tfstate
      // lock or leave half-created APIM resources behind, then rethrow.
      await tfCore.destroyWorkspace(ws).catch(() => {});
      throw err;
    }
    return new TerraformProvisioned<P>(ws, this.spec, vars);
  }
}
