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

import type { Provisioned, ProvisionerChecks, ProvisionerId, Role } from "./types.js";

/** Build the internal role string from a kind + optional disambiguating label. */
function roleFor(kind: string, label?: string): Role {
  return (label ? `${kind}:${label}` : kind) as Role;
}

/**
 * Shared base for the concrete provisioned handles. It defines the plain id
 * getters (`apiId()`, `subscriptionId()`, ...) ONCE, each delegating to a single
 * `resolveId(role)` that GKO/Terraform implement. A test author calls
 * `provision.apiId()`; the "role" string stays an internal detail. Add a new
 * kind's getter here in one line.
 */
export abstract class BaseProvisioned<P = unknown> implements Provisioned<P> {
  abstract readonly provisionerId: ProvisionerId;
  abstract readonly checks: ProvisionerChecks;

  /** Resolve a logical role to its APIM id (UUID). Resolved once, then cached. */
  protected abstract resolveId(role: Role): Promise<string>;

  abstract contextPath(): Promise<string>;
  abstract update(params: P): Promise<void>;
  abstract remove(role: Role): Promise<void>;
  abstract destroy(): Promise<void>;

  apiId(label?: string): Promise<string> {
    return this.resolveId(roleFor("api", label));
  }
  subscriptionId(label?: string): Promise<string> {
    return this.resolveId(roleFor("subscription", label));
  }
  applicationId(label?: string): Promise<string> {
    return this.resolveId(roleFor("application", label));
  }
  groupId(label?: string): Promise<string> {
    return this.resolveId(roleFor("group", label));
  }
}
