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
 * Enum value rejection.
 *
 * The operator only exposes the API / SharedPolicyGroup enum values the
 * Automation API actually supports; the CRD schemas were trimmed accordingly.
 * These negative tests assert the Kubernetes API server rejects the removed
 * values at admission time (CRD schema validation), and double as a regression
 * guard: if a future CRD regeneration re-introduces EDGE / NATIVE / INTERACT /
 * CONNECT / ENTRYPOINT_CONNECT, the matching test fails.
 *
 * Rejection happens at the CRD schema layer (apiserver), before the operator
 * admission webhook and before any APIM call, so the manifests are otherwise
 * valid — the only reason they are refused is the offending enum value. Each
 * assertion pins the rejection to that value AND its field path so a manifest
 * mistake cannot make the test pass for the wrong reason.
 *
 * Xray tests:
 *   GKO-2965: V4 API EDGE type is rejected by the CRD schema
 *   GKO-2966: SPG NATIVE apiType is rejected by the CRD schema
 *   GKO-2967: SPG removed flow phases are rejected by the CRD schema
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running (the fixed CRDs are installed)
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectl from "../../../helpers/kubectl.js";

const V4_NAME = "e2e-gko2934-v4-edge";
const SPG_NATIVE_NAME = "e2e-gko2934-spg-native";
const SPG_PHASE_PREFIX = "e2e-gko2934-spg-phase";

// SPG flow phases removed from the FlowPhase enum (no longer supported).
const REMOVED_SPG_PHASES = ["INTERACT", "CONNECT", "ENTRYPOINT_CONNECT"] as const;

// k8s metadata.name must be RFC 1123 DNS-safe, so map underscores to hyphens
// (ENTRYPOINT_CONNECT -> entrypoint-connect). This keeps every phase manifest
// valid for all reasons EXCEPT the offending phase enum value — the only
// intended rejection cause.
const phaseResourceName = (phase: string) =>
  `${SPG_PHASE_PREFIX}-${phase.toLowerCase().replaceAll("_", "-")}`;

/** Build a structurally valid SharedPolicyGroup manifest with a given apiType/phase. */
function spgManifest(opts: { name: string; apiType: string; phase: string }): string {
  return `apiVersion: gravitee.io/v1alpha1
kind: SharedPolicyGroup
metadata:
  name: ${opts.name}
spec:
  contextRef:
    name: "dev-ctx"
    namespace: "default"
  name: "${opts.name}"
  description: "enum rejection negative test"
  apiType: ${opts.apiType}
  phase: ${opts.phase}
  steps:
    - name: "Transform Headers"
      enabled: true
      policy: transform-headers
      configuration:
        addHeaders:
          - name: "X-Test"
            value: "x"
`;
}

/** Build a structurally valid (minimal) V4 API manifest with a given type. */
function v4Manifest(opts: { name: string; type: string }): string {
  return `apiVersion: gravitee.io/v1alpha1
kind: ApiV4Definition
metadata:
  name: ${opts.name}
spec:
  contextRef:
    name: dev-ctx
    namespace: default
  name: ${opts.name}
  version: "1.0"
  description: "enum rejection negative test"
  type: ${opts.type}
  state: STARTED
  listeners:
    - type: HTTP
      paths:
        - path: /${opts.name}
      entrypoints:
        - type: http-proxy
          qos: AUTO
  endpointGroups:
    - name: Default HTTP proxy group
      type: http-proxy
      endpoints:
        - name: Default HTTP proxy
          type: http-proxy
          inheritConfiguration: false
          configuration:
            target: https://api.gravitee.io/echo
  flowExecution:
    mode: DEFAULT
    matchRequired: false
  plans:
    KeyLess:
      name: Free plan
      security:
        type: KEY_LESS
`;
}

test.describe("Admission — Enum value rejection", () => {
  // Safety-net cleanup: schema-rejected resources are never persisted, but if a
  // regression makes admission accept one, remove it so it can't leak into
  // downstream tests. Each delete ignores errors (the resource is usually gone).
  test.afterEach(async () => {
    await kubectl.deleteResource("apiv4definition", V4_NAME).catch(() => {});
    await kubectl.deleteResource("sharedpolicygroup", SPG_NATIVE_NAME).catch(() => {});
    for (const phase of REMOVED_SPG_PHASES) {
      await kubectl.deleteResource("sharedpolicygroup", phaseResourceName(phase)).catch(() => {});
    }
  });

  // ── V4 API: EDGE type ───────────────────────────────────────

  test(`V4 API with EDGE type is rejected by the CRD schema ${XRAY.ENUM_VALIDATION.V4_EDGE_TYPE_REJECTED} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyStringExpectFailure(
      v4Manifest({ name: V4_NAME, type: "EDGE" }),
    );
    expect(stderr.toLowerCase()).toContain("unsupported value");
    expect(stderr).toContain("EDGE");
    expect(stderr).toContain("spec.type");
  });

  // ── SPG: NATIVE apiType ─────────────────────────────────────

  test(`SharedPolicyGroup with NATIVE apiType is rejected by the CRD schema ${XRAY.ENUM_VALIDATION.SPG_APITYPE_NATIVE_REJECTED} ${TAGS.REGRESSION}`, async ({
    kubectl,
  }) => {
    const stderr = await kubectl.applyStringExpectFailure(
      spgManifest({ name: SPG_NATIVE_NAME, apiType: "NATIVE", phase: "REQUEST" }),
    );
    expect(stderr.toLowerCase()).toContain("unsupported value");
    expect(stderr).toContain("NATIVE");
    expect(stderr).toContain("spec.apiType");
  });

  // ── SPG: removed flow phases (INTERACT / CONNECT / ENTRYPOINT_CONNECT) ──

  for (const phase of REMOVED_SPG_PHASES) {
    test(`SharedPolicyGroup with ${phase} phase is rejected by the CRD schema ${XRAY.ENUM_VALIDATION.SPG_REMOVED_PHASE_REJECTED} ${TAGS.REGRESSION}`, async ({
      kubectl,
    }) => {
      const stderr = await kubectl.applyStringExpectFailure(
        spgManifest({
          name: phaseResourceName(phase),
          apiType: "PROXY",
          phase,
        }),
      );
      expect(stderr.toLowerCase()).toContain("unsupported value");
      expect(stderr).toContain(phase);
      expect(stderr).toContain("spec.phase");
    });
  }
});
