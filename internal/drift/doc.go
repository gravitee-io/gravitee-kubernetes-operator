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
// functions and initialize enable/disable predicates. Tests call it from a BeforeSuite hook.
// [Init] calls [InitRegistry] and [InitEnableCheck].
//
// # Comparison model
//
// Comparison is driven by struct field tags on DTO types in api/model:
//
//	drift:"<equivalence-name>[:arg1,arg2,...]"
//
// The optional arguments after the colon are passed to the equivalence function via
// [DriftContext.FuncArgs]. Property names in the output tree come from the json struct
// tag (or the lower-cased field name when no json tag is present). Embedded structs are
// flattened into their parent. Pointer fields are dereferenced before the equivalence
// kind is resolved.
//
// [DetectWithNamespace] walks two values of the same struct type recursively and builds a
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
// # Drift context
//
// [DriftContext] carries namespace information and function arguments from drift tags
// through the comparison tree. It is passed to all [EquivalenceFunc] calls.
//
// # Equivalence functions
//
// An [EquivalenceFunc] receives the CRD value, the remote value, and a [DriftContext],
// and returns an [Equivalence] with:
//   - [Equivalence.Equivalent]: the equivalence status (Equivalent, Inequivalent, CannotCompare)
//   - [Equivalence.Skip]: if true, children of this node are not compared
//   - [Equivalence.PostFunc]: optional hook called after children are processed
//   - [Equivalence.RemoteItemsFilterFunc]: optional function to filter remote slice items before comparison
//
// Registered names (see [InitRegistry]):
//
//   - empty-is-nil (string): nil or "" are equivalent; otherwise deep-equal.
//   - empty-is-nil (bool): nil or false are equivalent; otherwise deep-equal.
//   - empty-is-nil (int, int32): nil or 0 are equivalent; otherwise deep-equal.
//   - empty-is-nil (uint): nil or 0 are equivalent; otherwise deep-equal.
//   - empty-is-nil (slice, map): nil or len==0 are equivalent and skip children.
//   - empty-is-nil (struct): nil vs zero-value struct are compared; if equivalent, skip children.
//   - empty-is-true (bool): nil or true are equivalent; otherwise deep-equal.
//   - ignore (string, slice, struct): always returns CannotCompare (skips for struct via ignore-skip pattern).
//   - ignore-skip (struct): same as ignore but also sets Skip=true.
//   - trimmed (string): compares strings after [strings.TrimSpace].
//   - rfc3339 (string): compares instants; accepts RFC3339 and RFC3339Nano inputs.
//   - case-insensitive (string): compares strings case-insensitively.
//   - ignore-remote (string): ignores difference if remote value matches any of the tag arguments.
//   - ignore-namespace-prefix (string): strips namespace prefix before comparing.
//   - ignore-remote-only-metadata (slice): filters out remote-only Metadata items before comparison.
//   - ignore-unknown-crd-groups (slice): removes CRD-only strings, then compares with namespace prefix ignored.
//   - unstructured (struct): for unstructured types; hoists "object" child fields to root via PostFunc.
//
// # Drift Tag Function Arguments
//
// Some drift equivalence functions accept arguments to customize their behavior.
// Arguments are specified after a colon in the drift tag:
//
//	drift:"<equivalence-name>:arg1,arg2,..."
//
// The arguments are passed to the equivalence function via [DriftContext.FuncArgs].
//
// ## Functions with Arguments
//
// ### ignore-remote
//
// Syntax: `drift:"ignore-remote:value1,value2,..."`
//
// Ignores differences if the remote value matches any of the specified arguments.
// Useful for fields that have a known default value in APIM that may differ from
// the CRD representation.
//
// Example:
//
//	// Ignore if remote is "DEFAULT"
//	FlowMode v4.FlowMode `json:"mode,omitempty" drift:"ignore-remote:DEFAULT"`
//
//	// Ignore if remote is "AUTO" or "DEFAULT"
//	QOS v4.QOS `json:"qos,omitempty" drift:"ignore-remote:AUTO,DEFAULT"`
//
// ### ignore-namespace-prefix
//
// Syntax: `drift:"ignore-namespace-prefix"`
//
// Strips the namespace prefix from both CRD and remote values before comparing.
// The namespace is obtained from the [DriftContext]. This is useful for IDs
// that include the namespace as a prefix.
//
// Example:
//
//	// Strip "my-namespace-" prefix before comparing
//	ID string `json:"id,omitempty" drift:"ignore-namespace-prefix"`
//
// # Defaults without a drift tag
//
// When no drift tag is set on a field, the registry falls back to:
//
//   - slices, arrays: [CannotCompare] at container level; items are still compared.
//   - structs: [CannotCompare] at container level; children are still compared.
//   - other kinds: [DefaultEquivalence] (reflect.DeepEqual).
//
// Unknown drift tag names panic at runtime. Registered functions are keyed by name
// and reflect.Kind; register concrete kinds, not pointers ([RegisterEquivalenceFunc]).
//
// # Enabling drift detection
//
// Drift detection is disabled globally by default via [env.Config.DriftDetection.Enabled] (set by
// DRIFT_DETECTION_ENABLED environment variable). It can be overridden per resource using the
// gravitee.io/drift-detection annotation with values "true" or "false".
//
// Additionally, some resource types are unsupported or disabled by default:
//   - Legacy Group resources (non-Automation API) are unsupported
//   - Portal, Documentation, and PortalListing resources are disabled by default
//
// See [InitEnableCheck] and [IsDriftEnabled] for the predicate system.
//
// # Extending
//
// Register additional equivalence functions with [RegisterEquivalenceFunc] and call
// them from [InitRegistry]. Annotate api/model fields with the matching drift tag name.
//
// Reference fixtures and behaviour tables live in test/unit/drift/.
package drift
