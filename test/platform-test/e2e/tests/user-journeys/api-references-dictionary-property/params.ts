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
 * Parameter binding for the MANUAL dictionary update scenario. ONE shared param
 * type drives both provisioners: the GKO arm renders the MANUAL Dictionary CR
 * from the params (applied via applyParams so an update() re-applies a changed
 * CR), and the Terraform arm maps the params to the `env_value` tfvar the
 * co-located main.tf reads. The resolve scenario keeps using the static
 * gko/dictionary.yaml (default value), so this only backs the update scenario.
 */

/** The GKO Dictionary CR name; its HRID is `default-<name>` (namespace + name). */
export const GKO_MANUAL_DICT_NAME = "manual-dictionary";

/** What the MANUAL dictionary update scenario is parameterized by. */
export interface ManualDictParams {
  /** The value of the dictionary's `env` property, resolved at the gateway. */
  envValue: string;
}

/** Initial params: the same "test" value the resolve scenario asserts. */
export const MANUAL_INITIAL: ManualDictParams = { envValue: "test" };

/**
 * Render the MANUAL Dictionary CR for the given params. Re-applying this same
 * document (same name) is how the GKO arm's update() propagates a changed
 * property value.
 */
export function manualDictionaryYaml(params: ManualDictParams): string {
  return `apiVersion: gravitee.io/v1alpha1
kind: Dictionary
metadata:
  name: ${GKO_MANUAL_DICT_NAME}
  labels:
    gravitee.io/e2e: "true"
spec:
  contextRef:
    name: "dev-ctx"
    namespace: "default"
  name: "${GKO_MANUAL_DICT_NAME}"
  type: MANUAL
  deployed: true
  manual:
    properties:
      env: "${params.envValue}"
`;
}

/** Map the shared params to the tfvars the terraform/main.tf fixture reads. */
export function tfManualDictVars(params: ManualDictParams): Record<string, unknown> {
  return { env_value: params.envValue };
}
