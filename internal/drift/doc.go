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
//   - [Equivalence.CRDItemsFilterFunc]: optional function to filter CRD slice items before comparison
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
//   - time (struct): compares [time.Time] values as instants, ignoring location and
//     monotonic clock, and skips children; any other struct gets the default struct equivalence.
//   - case-insensitive (string): compares strings case-insensitively.
//   - ignore-remote-default (string): ignores a difference when the CRD value is unset;
//     with no arguments any remote value is accepted, with arguments only a listed
//     remote value (a server default) is.
//   - ignore-namespace-prefix (string): strips namespace prefix before comparing.
//   - ignore-only (slice): filters items present only on the side given by the tag argument (remote or crd).
//     With strip-ns, remaining keys are compared as a set after stripping the namespace prefix.
//     With expired and/or scheduled, items implementing [Expiring] / [Schedulable] are dropped first
//     (APIM omits those client certificates from GET responses).
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
// ### ignore-remote-default
//
// Syntax: `drift:"ignore-remote-default"` or `drift:"ignore-remote-default:value1,value2,..."`
//
// Ignores a difference when the CRD leaves the field unset. With no arguments, any
// remote value is accepted — use this when the CRD carries no information to predict
// what APIM resolves the field to. With arguments, only a remote value listed among
// them is accepted, and any other remote value is drift — use this when the CRD's
// silence should only cover specific, known server defaults. Either way, a CRD value
// that is set is always compared.
//
// Example:
//
//	// APIM applies DEFAULT when the payload omits the flow mode
//	Mode v4.FlowMode `json:"mode,omitempty" drift:"ignore-remote-default:DEFAULT"`
//
//	// APIM resolves an omitted next-gen portal visibility from the parent folder, which
//	// lives in another resource here — any value APIM may resolve to is accepted
//	Visibility nav.Visibility `json:"visibility,omitempty" drift:"ignore-remote-default"`
//
// Where the operator can see the whole ancestor chain — Portal.structure.topNavbar,
// ApiV4Definition.portalNavigation — the expected visibility is resolved instead, by
// [github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model.PortalDTO.WithResolvedVisibility],
// and the field is compared exactly with no tag at all.
//
// ### ignore-only
//
// Syntax: `drift:"ignore-only:remote"` or `drift:"ignore-only:crd"` or `drift:"ignore-only:crd,strip-ns"`
// or `drift:"ignore-only:crd,expired,scheduled"`
//
// Filters slice items that exist only on the named side before item comparison.
// Items must implement [Keyed], or be strings (including named string types), which are
// their own key.
// With `strip-ns`, keys are compared after stripping the namespace prefix
// and remaining membership is compared as a set (Skip).
// With `expired`, items implementing [Expiring] whose [Expiring.Expired] is true are
// dropped on both sides (GKO-2221: APIM omits expired client certificates).
// With `scheduled`, items implementing [Schedulable] whose [Schedulable.Scheduled] is true
// are dropped on both sides (GKO-2228: APIM omits future-startsAt certificates).
//
// Example:
//
//	Metadata []*APIV4MetadataEntryDTO `json:"metadata" drift:"ignore-only:remote"`
//	Groups   []APIGroup               `json:"groups" drift:"ignore-only:crd,strip-ns"`
//	ClientCertificates []ApplicationClientCertificateDTO `json:"clientCertificates,omitempty" drift:"ignore-only:crd,expired,scheduled"`
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
//     A [time.Time] has only unexported fields, so an untagged one panics: tag it
//     `drift:"time"` (or `ignore`).
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
//   - Subscriptions bound to a definition-v2 ApiDefinition are unsupported
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
