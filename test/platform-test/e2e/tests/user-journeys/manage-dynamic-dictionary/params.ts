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
 * Parameter binding for the manage-dynamic-dictionary journey. ONE shared param
 * type (`DynamicDictParams`) drives both provisioners: the GKO arm renders the
 * DYNAMIC Dictionary CR from the params (applied via applyParams so an update()
 * re-applies a changed CR), and the Terraform arm maps the params to the tfvars
 * the co-located main.tf reads. This is the scenario-specific seam the
 * provisioner core deliberately does not own (mirrors subscribe-and-call/params.ts).
 */

/** The GKO Dictionary CR name; its HRID is `default-<name>` (namespace + name). */
export const GKO_DICT_NAME = "dyn-dictionary";

/**
 * JOLT shift spec mapping every echoed response header into a key/value
 * dictionary property. The echo endpoint reflects the provider's own
 * `X-Test-Specific` request header, so this yields a property of that name whose
 * value is the request-header value. Kept as a single-line JSON string so it
 * embeds cleanly in the generated YAML.
 */
const JOLT_SPEC = JSON.stringify([
  { operation: "shift", spec: { headers: { "*": { $: "[#2].key", "@": "[#2].value" } } } },
]);

/** What every manage-dynamic-dictionary scenario can be parameterized by. */
export interface DynamicDictParams {
  /** The provider request-header value the JOLT spec surfaces as `X-Test-Specific`. */
  headerValue: string;
  /** Deploy (start) the dictionary at the gateway, or stop/undeploy it. */
  deployed: boolean;
}

/** The initial params every scenario provisions with: a deployed dict serving "ABCDEF". */
export const INITIAL: DynamicDictParams = { headerValue: "ABCDEF", deployed: true };

/**
 * Render the DYNAMIC Dictionary CR for the given params. Re-applying this same
 * document (same name) is how the GKO arm's update() propagates a changed
 * provider header value or flips the deployed flag.
 */
export function dictionaryYaml(params: DynamicDictParams): string {
  return `apiVersion: gravitee.io/v1alpha1
kind: Dictionary
metadata:
  name: ${GKO_DICT_NAME}
  labels:
    gravitee.io/e2e: "true"
spec:
  contextRef:
    name: "dev-ctx"
    namespace: "default"
  name: "${GKO_DICT_NAME}"
  description: "DYNAMIC dictionary exposing echo headers as properties"
  type: DYNAMIC
  deployed: ${params.deployed}
  dynamic:
    provider:
      type: HTTP
      url: "https://api.gravitee.io/echo"
      method: GET
      headers:
        - name: "X-Test-Specific"
          value: "${params.headerValue}"
      specification: '${JOLT_SPEC}'
    trigger:
      rate: 5
      unit: SECONDS
`;
}

/** Map the shared params to the tfvars the terraform/main.tf fixture reads. */
export function tfDynamicDictVars(params: DynamicDictParams): Record<string, unknown> {
  return {
    header_value: params.headerValue,
    deployed: params.deployed,
    create_dictionary: true,
  };
}
