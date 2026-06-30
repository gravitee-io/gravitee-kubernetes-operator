// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package drift compares a Kubernetes CRD payload with a remote Gravitee APIM object
// and reports structural differences as a tree of [Result] nodes.
//
// It is used at admission time to reject updates when the remote API was changed
// outside the operator (see [github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift]).
//
// # Initialization
//
// Call [Init] once at process startup (from main) to register built-in equivalence
// functions. Tests call it from a BeforeSuite hook.
//
// # Comparison model
//
// Comparison is driven by struct field tags on DTO types in api/model:
//
//	drift:"<equivalence-name>"
//
// Property names in the output tree come from the json struct tag (or the lower-cased
// field name when no json tag is present). Embedded structs are flattened into their
// parent. Pointer fields are dereferenced before the equivalence kind is resolved.
//
// [Detect] walks two values of the same struct type recursively and builds a
// [Result] tree. Leaf nodes that are inequivalent are formatted as pseudo-YAML by
// [Result.String]. [Result.DriftDetected] returns true when any node in the tree is
// inequivalent.
//
// At admission, the operator compares old and new CRD revisions against the same
// remote snapshot using [Merge]:
//
//   - O/R: old CRD DTO vs remote
//   - N/R: new CRD DTO vs remote
//
// Drift is reported only when both comparisons diverge from remote in a way that
// cannot be explained by the user's CRD update (see [Merge] for the five cases).
//
// # Equivalence functions
//
// An [EquivalenceFunc] receives the CRD value and the remote value and returns an
// [Equivalence] with a status, an optional [Equivalence.Skip] flag, and an optional
// [Equivalence.PostFunc] hook.
//
// Registered names (see [Init]):
//
//   - empty-is-nil (string): nil or "" are equivalent; otherwise deep-equal.
//   - empty-is-nil (bool): nil or false are equivalent; otherwise deep-equal.
//   - empty-is-nil (int): nil or 0 are equivalent; otherwise deep-equal.
//   - empty-is-nil (uint): nil or 0 are equivalent; otherwise deep-equal.
//   - empty-is-nil (slice, map): 0-len or nil is considered as equivalent.
//   - empty-is-nil (struct): nil vs zero-value struct is equivalent and skips
//     children; otherwise children are compared.
//   - ignore (string): always equivalent (e.g. APIM-export-only IDs).
//   - ignore (struct): always equivalent and skips children.
//   - trimmed (string): compares strings after [strings.TrimSpace].
//   - rfc3339 (string): compares instants; accepts RFC3339 and RFC3339Nano inputs.
//   - unstructured (struct): for [k8s.io/apimachinery/pkg/apis/meta/v1/unstructured.Unstructured]
//     and [github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils.GenericStringMap];
//     hoists the "object" child fields to the root via [Equivalence.PostFunc].
//
// # Defaults without a drift tag
//
// When no drift tag is set on a field, the registry falls back to:
//
//   - slices: [CannotCompare] at container level, items are still compared.
//   - structs: [CannotCompare] at container level, items are still compared.
//   - other kinds: [FromDeepEqual] (reflect.DeepEqual).
//
// Unknown drift tag names panic at runtime. Registered functions are keyed by name
// and reflect.Kind; register concrete kinds, not pointers ([RegisterEquivalenceFunc]).
//
// # Extending
//
// Register additional equivalence functions with [RegisterEquivalenceFunc] and call
// them from [Init]. Annotate api/model fields with the matching drift tag name.
//
// Reference fixtures and behaviour tables live in test/unit/drift/.
//
// # Enabling drift detection
//
// Globally via the ENABLE_DRIFT_DETECTION environment variable (enabled by default).
// Per resource via the gravitee.io/drift-detection annotation (true/false).
package drift
