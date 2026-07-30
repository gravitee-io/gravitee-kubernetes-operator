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

import type { K8sCondition } from "../engines/kubectl.js";
import type { ProvisionerView } from "../view.js";

/** Evidence GKO's `view.read()` reports: the condition it inspected, plus severe errors on failure. */
export interface GkoViewDetail {
  reason: string;
  condition?: K8sCondition;
  severeErrors?: string[];
}

/**
 * GKO's `ProvisionerView`, reachable from a provisioned handle via
 * `provision.view` — no narrowing needed (unlike `checks`), `view` is meant to
 * be called from shared, provisioner-agnostic scenario bodies directly.
 */
export type GkoView = ProvisionerView<GkoViewDetail>;
