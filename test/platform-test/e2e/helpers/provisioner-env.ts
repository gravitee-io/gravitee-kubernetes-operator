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
 * e2e binding for the runner-agnostic provisioner core. This supplies the two
 * things the `src/` core must not own (it cannot import from `e2e/`): the APIM
 * auth/server environment loaded from `config.yaml`, and fixture-path
 * resolution via `fixture()`.
 *
 * Scenario authors use `gkoScenario(...)` / `tfScenario(...)` with
 * fixture-RELATIVE paths; this module resolves them to the absolute paths /
 * env the provisioner constructors expect, and returns the `() => Provisioner`
 * factories that `forProvisioners` consumes.
 */

import { fixture as resolveFixture } from "../setup.js";
import { terraformEnv } from "./terraform.js";
import {
  GkoProvisioner,
  TerraformProvisioner,
  type GkoScenarioSpec,
  type Provisioner,
  type TfScenarioSpec,
} from "../../src/provisioners/index.js";

/**
 * The TF APIM env (config.yaml + process.env) is static per run, so resolve it
 * once and memoize. Built lazily on first Terraform scenario use.
 */
let tfEnvPromise: Promise<Record<string, string>> | undefined;
function tfEnv(): Promise<Record<string, string>> {
  tfEnvPromise ??= terraformEnv();
  return tfEnvPromise;
}

/** GKO scenario as authored in tests: manifests are fixture-relative paths. */
export interface GkoScenarioInput<P = unknown>
  extends Omit<GkoScenarioSpec<P>, "manifests"> {
  /** Fixture-relative manifest paths, e.g. "subscriptions/apikey-auto/crd.yaml". */
  manifests: string[];
}

/** Build a GKO provisioner factory from a fixture-relative scenario. */
export function gkoScenario<P = unknown>(
  input: GkoScenarioInput<P>,
): () => Provisioner<P> {
  return () =>
    new GkoProvisioner<P>({
      ...input,
      manifests: input.manifests.map((m) => resolveFixture(m)),
    });
}

/** Terraform scenario as authored in tests: `fixture` is a folder name. */
export interface TfScenarioInput<P = unknown>
  extends Omit<TfScenarioSpec<P>, "fixtureDir" | "env"> {
  /** Fixture-relative folder name containing main.tf, e.g. "subscriptions/apikey-auto". */
  fixture: string;
}

/** Build a Terraform provisioner factory from a fixture-relative scenario. */
export function tfScenario<P = unknown>(
  input: TfScenarioInput<P>,
): () => Promise<Provisioner<P>> {
  const { fixture: fixtureName, ...rest } = input;
  return async () =>
    new TerraformProvisioner<P>({
      ...rest,
      fixtureDir: resolveFixture(fixtureName),
      env: await tfEnv(),
    });
}
