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

import * as kubectl from "../engines/kubectl.js";
import type { Provisioned, Provisioner, Role } from "../types.js";
import { BaseProvisioned } from "../base.js";
import type { GkoChecks } from "./checks.js";

/** The kubectl engine surface a parameterized apply step receives. */
export type KubectlEngine = typeof kubectl;

/** How a logical role maps to a Kubernetes CR. */
export interface GkoRoleBinding {
  /** kubectl kind, e.g. "apiv4definition" | "application" | "subscription". */
  kind: string;
  /** metadata.name of the CR. */
  name: string;
  /** condition to wait for at provision time. Default "Accepted". */
  readyCondition?: string;
}

/** Convention mapping a role to a default kubectl kind (shorthand role maps). */
const DEFAULT_KIND_BY_ROLE: Record<string, string> = {
  api: "apiv4definition",
  application: "application",
  subscription: "subscription",
  group: "group",
};

/**
 * Role map: either the shorthand `role -> CR name` (kind derived by convention)
 * or the full `role -> {kind, name}` form for multi-resource scenarios.
 */
export type GkoRoles = Record<string, string | GkoRoleBinding>;

export interface GkoScenarioSpec<P = unknown> {
  /** Static manifests (absolute paths) applied at provision(), in order. */
  manifests: string[];
  /** Logical roles owned by this scenario (resolve ids, wait conditions, clean up). */
  roles: GkoRoles;
  /**
   * Roles created by the parameterized {@link applyParams} step (e.g.
   * "subscription"). They are awaited AFTER applyParams; all other roles are
   * awaited right after the static manifests.
   */
  dynamicRoles?: Role[];
  /** Gateway context path (GKO has no output to read it from). Omit when the scenario never hits the gateway. */
  contextPath?: string;
  /**
   * Optional parameterized apply (e.g. a Subscription with a given key set).
   * Runs at provision() after the static manifests, and again on every update().
   */
  applyParams?: (k: KubectlEngine, params: P) => Promise<void>;
  /** Kubernetes namespace (defaults to the kubectl engine default). */
  namespace?: string;
}

interface ResolvedGkoSpec<P> {
  manifests: string[];
  roles: Record<string, GkoRoleBinding>;
  dynamicRoles: Role[];
  contextPath?: string;
  applyParams?: (k: KubectlEngine, params: P) => Promise<void>;
  namespace?: string;
}

function normalizeRoles(roles: GkoRoles): Record<string, GkoRoleBinding> {
  const out: Record<string, GkoRoleBinding> = {};
  for (const [role, value] of Object.entries(roles)) {
    if (typeof value === "string") {
      const kind = DEFAULT_KIND_BY_ROLE[role];
      if (!kind) {
        throw new Error(
          `GKO role "${role}" has no kind convention; use the full { kind, name } form`,
        );
      }
      out[role] = { kind, name: value };
    } else {
      out[role] = value;
    }
  }
  return out;
}

function resolveBinding<P>(spec: ResolvedGkoSpec<P>, role: Role): GkoRoleBinding {
  const binding = spec.roles[role];
  if (!binding) {
    throw new Error(
      `GKO scenario has no role "${role}". Known roles: ${Object.keys(spec.roles).join(", ")}`,
    );
  }
  return binding;
}

function buildGkoChecks<P>(spec: ResolvedGkoSpec<P>): GkoChecks {
  return {
    provisionerId: "gko",
    async waitForCondition(role: Role, condition = "Accepted", timeoutSeconds?: number) {
      const b = resolveBinding(spec, role);
      await kubectl.waitForCondition(b.kind, b.name, condition, timeoutSeconds, spec.namespace);
    },
    async assertEventContains(role: Role, message: string) {
      const b = resolveBinding(spec, role);
      await kubectl.assertEventContains(b.kind, b.name, message, spec.namespace);
    },
    async status<T = unknown>(role: Role): Promise<T> {
      const b = resolveBinding(spec, role);
      return kubectl.getStatus<T>(b.kind, b.name, spec.namespace);
    },
    async applyExpectFailure(yaml: string): Promise<string> {
      return kubectl.applyStringExpectFailure(yaml, spec.namespace);
    },
  };
}

/**
 * Tear down a GKO scenario's resources in reverse dependency order: dynamic
 * roles (subscriptions) first so GKO finalizers release, then the static
 * manifests (application, api). Tolerant: any resource may already be gone.
 * Shared by the handle's destroy() and the provisioner-level cleanup().
 */
async function teardownGko<P>(spec: ResolvedGkoSpec<P>): Promise<void> {
  for (const role of spec.dynamicRoles) {
    const b = spec.roles[role];
    if (b) await kubectl.deleteResource(b.kind, b.name, spec.namespace).catch(() => {});
  }
  for (const manifest of [...spec.manifests].reverse()) {
    await kubectl.del(manifest, spec.namespace).catch(() => {});
  }
}

class GkoProvisioned<P> extends BaseProvisioned<P> {
  readonly provisionerId = "gko" as const;
  readonly checks: GkoChecks;
  private readonly idCache = new Map<string, string>();

  constructor(private readonly spec: ResolvedGkoSpec<P>) {
    super();
    this.checks = buildGkoChecks(spec);
  }

  protected async resolveId(role: Role): Promise<string> {
    const cached = this.idCache.get(role);
    if (cached) return cached;
    const b = resolveBinding(this.spec, role);
    const status = await kubectl.getStatus<{ id?: string }>(b.kind, b.name, this.spec.namespace);
    if (!status?.id) {
      throw new Error(`GKO ${b.kind}/${b.name} has no .status.id yet`);
    }
    this.idCache.set(role, status.id);
    return status.id;
  }

  async contextPath(): Promise<string> {
    if (this.spec.contextPath === undefined) {
      throw new Error("GKO scenario has no contextPath; declare it to use gateway assertions");
    }
    return this.spec.contextPath;
  }

  async update(params: P): Promise<void> {
    if (this.spec.applyParams) {
      await this.spec.applyParams(kubectl, params);
    }
    // No waitForCondition here on purpose: a changed-spec re-apply can leave
    // Accepted=true momentarily (the operator has not yet observed the change),
    // so `kubectl wait` would return before convergence. Callers poll APIM via
    // mapi for the real convergence signal, exactly as the original tests did.
  }

  async remove(role: Role): Promise<void> {
    // Delete just this role's CR (the way a user removes one resource), leaving
    // the rest of the scenario standing, and wait for it to be gone.
    const b = resolveBinding(this.spec, role);
    await kubectl.deleteResource(b.kind, b.name, this.spec.namespace);
    await kubectl.waitForDeletion(b.kind, b.name, undefined, this.spec.namespace);
    this.idCache.delete(role);
  }

  async destroy(): Promise<void> {
    await teardownGko(this.spec);
  }
}

/**
 * Provisions a scenario through GKO: applies CRD manifests via kubectl, waits
 * for the operator to reconcile, and resolves logical roles to APIM ids from
 * each CR's `.status.id`.
 */
export class GkoProvisioner<P = unknown> implements Provisioner<P> {
  readonly provisionerId = "gko" as const;
  private readonly spec: ResolvedGkoSpec<P>;

  constructor(spec: GkoScenarioSpec<P>) {
    this.spec = {
      manifests: spec.manifests,
      roles: normalizeRoles(spec.roles),
      dynamicRoles: spec.dynamicRoles ?? [],
      contextPath: spec.contextPath,
      applyParams: spec.applyParams,
      namespace: spec.namespace,
    };
  }

  async provision(params: P): Promise<Provisioned<P>> {
    const ns = this.spec.namespace;
    const dynamic = new Set<string>(this.spec.dynamicRoles);

    // 1. Apply static manifests (API + application) in order.
    for (const manifest of this.spec.manifests) {
      await kubectl.apply(manifest, ns);
    }

    // 2. Wait for non-dynamic roles to reconcile.
    for (const [role, b] of Object.entries(this.spec.roles)) {
      if (dynamic.has(role)) continue;
      await kubectl.waitForCondition(b.kind, b.name, b.readyCondition ?? "Accepted", undefined, ns);
    }

    // 3. Apply the parameterized step (e.g. the subscription with its key set).
    if (this.spec.applyParams) {
      await this.spec.applyParams(kubectl, params);
    }

    // 4. Wait for the dynamic roles to reconcile.
    for (const role of this.spec.dynamicRoles) {
      const b = this.spec.roles[role];
      if (b) {
        await kubectl.waitForCondition(b.kind, b.name, b.readyCondition ?? "Accepted", undefined, ns);
      }
    }

    return new GkoProvisioned<P>(this.spec);
  }

  /** Best-effort teardown when provision() failed partway (no live handle). */
  async cleanup(): Promise<void> {
    await teardownGko(this.spec);
  }

  /**
   * Rebuild a live handle for a scenario provisioned in an EARLIER process,
   * used by the upgrade survival flow (provision before the upgrade, attach
   * after it, in a fresh `npm run e2e`). GKO resources are Kubernetes CRs that
   * persist across an in-place APIM upgrade, so there is no in-memory state to
   * carry: we confirm each referenced role's CR is still present (a missing one
   * means it did not survive) and return a fresh handle that re-resolves
   * `.status.id` on demand. The `hrid` in `refs` is not needed here because the
   * CR kind/name from the spec already locates the resource; `refs` declares
   * which roles the caller expects to have survived.
   */
  async attach(refs: Partial<Record<Role, { hrid: string }>>): Promise<Provisioned<P>> {
    const ns = this.spec.namespace;
    // Poll for a settle window rather than a single-shot check: just after the
    // operator + CRD re-apply, a one-off `get` can transiently fail (API-server
    // settle, CRD discovery refresh) even though the CR is present. A resource
    // that genuinely did not survive never reappears, so this still fails after
    // the window - it avoids a false "did not survive" without masking a real one.
    const SETTLE_MS = 30_000;
    for (const role of Object.keys(refs) as Role[]) {
      const b = resolveBinding(this.spec, role);
      if (!(await kubectl.waitForExists(b.kind, b.name, ns, SETTLE_MS))) {
        throw new Error(
          `GKO ${b.kind}/${b.name} (role "${role}") did not survive the upgrade: ` +
            `not found in the cluster after ${SETTLE_MS / 1_000}s`,
        );
      }
    }
    return new GkoProvisioned<P>(this.spec);
  }
}
