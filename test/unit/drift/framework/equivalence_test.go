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

package framework

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
)

var _ = Describe("EmptyIsNilString", func() {
	DescribeTable("should report equivalence",
		func(crd, remote any) {
			Expect(drift.EmptyIsNilString(crd, remote, drift.DriftContext{})).To(Equal(
				drift.Equivalence{Equivalent: drift.Equivalent},
			))
		},
		Entry("empty string vs nil", "", nil),
		Entry("nil vs empty string", nil, ""),
		Entry("empty string vs empty string", "", ""),
		Entry("nil vs nil", nil, nil),
	)

	DescribeTable("should report inequivalence",
		func(crd, remote any) {
			Expect(drift.EmptyIsNilString(crd, remote, drift.DriftContext{})).To(Equal(
				drift.Equivalence{Equivalent: drift.Inequivalent},
			))
		},
		Entry("different non-empty strings", "foo", "bar"),
		Entry("non-empty string vs nil", "foo", nil),
		Entry("nil vs non-empty string", nil, "foo"),
		Entry("non-empty string vs empty string", "foo", ""),
		Entry("empty string vs non-empty string", "", "foo"),
	)
})

var _ = Describe("EmptyIsNilInt", func() {
	DescribeTable("should report equivalence",
		func(crd, remote any) {
			Expect(drift.EmptyIsNilInt(crd, remote, drift.DriftContext{})).To(Equal(
				drift.Equivalence{Equivalent: drift.Equivalent},
			))
		},
		Entry("zero vs nil", 0, nil),
		Entry("nil vs zero", nil, 0),
		Entry("zero vs zero", 0, 0),
		Entry("nil vs nil", nil, nil),
		Entry("same non-zero value", 42, 42),
	)

	DescribeTable("should report inequivalence",
		func(crd, remote any) {
			Expect(drift.EmptyIsNilInt(crd, remote, drift.DriftContext{})).To(Equal(
				drift.Equivalence{Equivalent: drift.Inequivalent},
			))
		},
		Entry("different non-zero values", 1, 2),
		Entry("non-zero vs nil", 1, nil),
		Entry("nil vs non-zero", nil, 1),
		Entry("non-zero vs zero", 1, 0),
		Entry("zero vs non-zero", 0, 1),
		Entry("int zero vs uint zero", 0, uint(0)),
		Entry("uint zero vs int zero", uint(0), 0),
		Entry("nil vs uint zero", nil, uint(0)),
		Entry("uint zero vs nil", uint(0), nil),
	)
})

var _ = Describe("EmptyIsNilUint", func() {
	DescribeTable("should report equivalence",
		func(crd, remote any) {
			Expect(drift.EmptyIsNilUint(crd, remote, drift.DriftContext{})).To(Equal(
				drift.Equivalence{Equivalent: drift.Equivalent},
			))
		},
		Entry("zero vs nil", uint(0), nil),
		Entry("nil vs zero", nil, uint(0)),
		Entry("zero vs zero", uint(0), uint(0)),
		Entry("nil vs nil", nil, nil),
		Entry("same non-zero value", uint(42), uint(42)),
	)

	DescribeTable("should report inequivalence",
		func(crd, remote any) {
			Expect(drift.EmptyIsNilUint(crd, remote, drift.DriftContext{})).To(Equal(
				drift.Equivalence{Equivalent: drift.Inequivalent},
			))
		},
		Entry("different non-zero values", uint(1), uint(2)),
		Entry("non-zero vs nil", uint(1), nil),
		Entry("nil vs non-zero", nil, uint(1)),
		Entry("non-zero vs zero", uint(1), uint(0)),
		Entry("zero vs non-zero", uint(0), uint(1)),
		Entry("int zero vs uint zero", 0, uint(0)),
		Entry("uint zero vs int zero", uint(0), 0),
		Entry("nil vs int zero", nil, 0),
		Entry("int zero vs nil", 0, nil),
	)
})

var _ = Describe("EmptyIsNilBool", func() {
	DescribeTable("should report equivalence",
		func(crd, remote any) {
			Expect(drift.EmptyIsNilBool(crd, remote, drift.DriftContext{})).To(Equal(
				drift.Equivalence{Equivalent: drift.Equivalent},
			))
		},
		Entry("false vs nil", false, nil),
		Entry("nil vs false", nil, false),
		Entry("false vs false", false, false),
		Entry("nil vs nil", nil, nil),
		Entry("true vs true", true, true),
	)

	DescribeTable("should report inequivalence",
		func(crd, remote any) {
			Expect(drift.EmptyIsNilBool(crd, remote, drift.DriftContext{})).To(Equal(
				drift.Equivalence{Equivalent: drift.Inequivalent},
			))
		},
		Entry("true vs nil", true, nil),
		Entry("nil vs true", nil, true),
		Entry("true vs false", true, false),
		Entry("false vs true", false, true),
	)
})

var _ = Describe("EmptyIsTrue", func() {
	DescribeTable("should report equivalence",
		func(crd, remote any) {
			Expect(drift.EmptyIsTrue(crd, remote, drift.DriftContext{})).To(Equal(
				drift.Equivalence{Equivalent: drift.Equivalent},
			))
		},
		Entry("nil vs nil", nil, nil),
		Entry("nil crd vs true remote", nil, true),
		Entry("true vs true", true, true),
		Entry("false vs false", false, false),
	)

	DescribeTable("should report inequivalence",
		func(crd, remote any) {
			Expect(drift.EmptyIsTrue(crd, remote, drift.DriftContext{})).To(Equal(
				drift.Equivalence{Equivalent: drift.Inequivalent},
			))
		},
		Entry("nil crd vs false remote", nil, false),
		Entry("false vs nil remote", false, nil),
		Entry("true vs false", true, false),
		Entry("false vs true", false, true),
		Entry("true vs nil remote", true, nil),
	)
})

var _ = Describe("IgnoreNamespacePrefix", func() {
	DescribeTable("should report equivalence when remote starts with namespace",
		func(crd, remote any) {
			Expect(drift.IgnoreNamespacePrefix(crd, remote, drift.DriftContext{Namespace: "foo"})).To(Equal(
				drift.Equivalence{Equivalent: drift.Equivalent},
			))
		},
		Entry("remote has ns", "api", "foo-api"),
		Entry("local has ns", "foo-api", "api"),
		Entry("both with ns", "foo-api", "foo-api"),
		Entry("none with ns", "api", "api"),
		Entry("both empty", "", ""),
		Entry("both nil", nil, nil),
		Entry("empty nil", "", nil),
		Entry("nil empty", nil, ""),
	)

	DescribeTable("should report inequivalence",
		func(crd, remote any) {
			Expect(drift.IgnoreNamespacePrefix(crd, remote, drift.DriftContext{Namespace: "foo"})).To(Equal(
				drift.Equivalence{Equivalent: drift.Inequivalent},
			))
		},
		Entry("not sharing anything", "x", "y"),
		Entry("same ns different values", "foo-x", "foo-y"),
		Entry("empty crd remote ns", "", "foo"),
		Entry("empty ns remote empty", "foo", ""),
		Entry("ns remote nil", "foo", nil),
		Entry("nil vs remote ns", nil, "foo"),
	)
})

var _ = Describe("IgnoreRemoteArgs", func() {
	type namedString string

	DescribeTable("should report equivalence",
		func(crd, remote any, ctx drift.DriftContext) {
			Expect(drift.IgnoreRemoteArgs(crd, remote, ctx)).To(Equal(
				drift.Equivalence{Equivalent: drift.Equivalent},
			))
		},
		Entry("equal strings without args", "foo", "foo", drift.DriftContext{}),
		Entry("equal strings with unused args", "foo", "foo", drift.DriftContext{FuncArgs: []string{"bar"}}),
		Entry("both nil", nil, nil, drift.DriftContext{}),
		Entry("both empty", "", "", drift.DriftContext{}),
		Entry("remote listed in FuncArgs", "foo", "bar", drift.DriftContext{FuncArgs: []string{"bar"}}),
		Entry("remote listed among several FuncArgs", "foo", "published", drift.DriftContext{FuncArgs: []string{"draft", "published"}}),
		Entry("nil crd vs remote listed in FuncArgs", nil, "bar", drift.DriftContext{FuncArgs: []string{"bar"}}),
		Entry("empty crd vs remote listed in FuncArgs", "", "bar", drift.DriftContext{FuncArgs: []string{"bar"}}),
		Entry("named string type remote listed in FuncArgs", namedString(""), namedString("bar"), drift.DriftContext{FuncArgs: []string{"bar"}}),
	)

	DescribeTable("should report inequivalence",
		func(crd, remote any, ctx drift.DriftContext) {
			Expect(drift.IgnoreRemoteArgs(crd, remote, ctx)).To(Equal(
				drift.Equivalence{Equivalent: drift.Inequivalent},
			))
		},
		Entry("different strings without args", "foo", "bar", drift.DriftContext{}),
		Entry("different strings with nil FuncArgs", "foo", "bar", drift.DriftContext{FuncArgs: nil}),
		Entry("different strings with empty FuncArgs", "foo", "bar", drift.DriftContext{FuncArgs: []string{}}),
		Entry("remote not listed in FuncArgs", "foo", "bar", drift.DriftContext{FuncArgs: []string{"baz"}}),
		Entry("crd listed in FuncArgs is not ignored", "foo", "bar", drift.DriftContext{FuncArgs: []string{"foo"}}),
		Entry("nil vs remote not listed", nil, "bar", drift.DriftContext{FuncArgs: []string{"foo"}}),
		Entry("non-empty vs nil remote", "foo", nil, drift.DriftContext{FuncArgs: []string{"foo"}}),
	)

	It("treats a remote listed in the ignore-remote tag as equivalent through Detect", func() {
		type withIgnoreRemote struct {
			Status string `json:"status" drift:"ignore-remote:PUBLISHED"`
		}
		crd := withIgnoreRemote{Status: "DRAFT"}
		remote := withIgnoreRemote{Status: "PUBLISHED"}
		expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
	})

	It("detects drift when the remote value is not listed in the ignore-remote tag", func() {
		type withIgnoreRemote struct {
			Status string `json:"status" drift:"ignore-remote:PUBLISHED"`
		}
		crd := withIgnoreRemote{Status: "DRAFT"}
		remote := withIgnoreRemote{Status: "ARCHIVED"}
		result := drift.DetectWithNamespace(crd, remote, "")
		Expect(result.DriftDetected()).To(BeTrue())
		Expect(result.String()).To(ContainSubstring(`status: "DRAFT" != "ARCHIVED"`))
	})
})

var _ = Describe("CaseInsensitive drift tags", func() {
	DescribeTable("should report equivalence", func(a, b any) {
		e := drift.CaseInsensitive(a, b, drift.DriftContext{})
		Expect(e.Equivalent).To(Equal(drift.Equivalent))
	},
		Entry("LC UC", "foo", "FOO"),
		Entry("UC LC", "FOO", "foo"),
		Entry("equal UC", "FOO", "FOO"),
		Entry("equal LC", "foo", "foo"),
		Entry("empty", "", ""),
		Entry("nil", nil, nil),
	)
	DescribeTable("should report inequivalence", func(a, b any) {
		e := drift.CaseInsensitive(a, b, drift.DriftContext{})
		Expect(e.Equivalent).To(Equal(drift.Inequivalent))
	},
		Entry("different UC", "FOO", "bar"),
		Entry("different LC", "foo", "BAR"),
		Entry("empty", "", "bar"),
		Entry("nil", nil, ""),
		Entry("nil", "", nil),
	)
})

var _ = Describe("EmptyIsNilLen", func() {
	DescribeTable("should report equivalence",
		func(crd, remote any) {
			Expect(drift.EmptyIsNilLen(crd, remote, drift.DriftContext{})).To(Equal(
				drift.Equivalence{Equivalent: drift.Equivalent, Skip: true},
			))
		},
		Entry("empty slice vs empty slice", []string{}, []string{}),
		Entry("nil slice vs empty slice", nil, []string{}),
		Entry("empty slice vs nil", []string{}, nil),
		Entry("empty map vs empty map", map[string]string{}, map[string]string{}),
		Entry("nil map vs empty map", nil, map[string]string{}),
		Entry("empty map vs nil map", map[string]string{}, nil),
	)

	DescribeTable("should report cannot compare",
		func(crd, remote any) {
			Expect(drift.EmptyIsNilLen(crd, remote, drift.DriftContext{})).To(Equal(
				drift.Equivalence{Equivalent: drift.CannotCompare},
			))
		},
		Entry("map with same length, different content", map[string]string{"a": "1"}, map[string]string{"b": "2"}),
		Entry("slice with same length, different content", []string{"foo"}, []string{"bar"}),
		Entry("slice with different lengths", []string{"foo"}, []string{}),
		Entry("slice with different lengths, longer crd", []string{"foo", "bar"}, []string{"baz"}),
		Entry("map with different lengths", map[string]string{"a": "1"}, map[string]string{}),
		Entry("map with different lengths, longer crd", map[string]string{"a": "1", "b": "2"}, map[string]string{"c": "3"}),
	)

	It("treats nil crd as empty when remote slice is populated", func() {
		Expect(drift.EmptyIsNilLen(nil, []string{"published"}, drift.DriftContext{})).To(Equal(
			drift.Equivalence{Equivalent: drift.CannotCompare},
		))
	})

	It("treats nil crd as empty when remote map is populated", func() {
		Expect(drift.EmptyIsNilLen(nil, map[string]string{"env": "prod"}, drift.DriftContext{})).To(Equal(
			drift.Equivalence{Equivalent: drift.CannotCompare},
		))
	})
})

type emptyIsNilInnerWithTags struct {
	Tags []string `drift:"empty-is-nil" json:"tags"`
}

type emptyIsNilOuterWithInner struct {
	Inner *emptyIsNilInnerWithTags `drift:"empty-is-nil" json:"inner"`
}

var _ = Describe("EmptyIsNilLen through Detect", func() {
	It("detects drift when nil inner struct is compared to a populated remote", func() {
		result := drift.DetectWithNamespace(emptyIsNilOuterWithInner{Inner: nil}, emptyIsNilOuterWithInner{Inner: &emptyIsNilInnerWithTags{Tags: []string{"published"}}}, "")
		Expect(result.DriftDetected()).To(BeTrue())
	})
})

type emptyIsNilStructFixture struct {
	Name  string
	Count int
}

var _ = Describe("EmptyIsNilStruct", func() {
	DescribeTable("should report equivalence and skip",
		func(crd, remote any) {
			Expect(drift.EmptyIsNilStruct(crd, remote, drift.DriftContext{})).To(Equal(
				drift.Equivalence{Equivalent: drift.Equivalent, Skip: true},
			))
		},
		Entry("nil vs zero struct", nil, emptyIsNilStructFixture{}),
		Entry("zero struct vs nil", emptyIsNilStructFixture{}, nil),
	)

	DescribeTable("should report cannot compare",
		func(crd, remote any) {
			Expect(drift.EmptyIsNilStruct(crd, remote, drift.DriftContext{})).To(Equal(
				drift.Equivalence{Equivalent: drift.CannotCompare},
			))
		},
		Entry("nil vs non-zero struct", nil, emptyIsNilStructFixture{Name: "foo"}),
		Entry("non-zero struct vs nil", emptyIsNilStructFixture{Name: "foo"}, nil),
		Entry("nil vs nil", nil, nil),
		Entry("equal non-nil structs", emptyIsNilStructFixture{Name: "foo"}, emptyIsNilStructFixture{Name: "foo"}),
		Entry("different non-nil structs", emptyIsNilStructFixture{Name: "foo"}, emptyIsNilStructFixture{Name: "bar"}),
	)
})

var _ = Describe("FromDeepEqual", func() {
	DescribeTable("should report equivalence",
		func(crd, remote any) {
			Expect(drift.DefaultEquivalence(crd, remote, drift.DriftContext{})).To(Equal(
				drift.Equivalence{Equivalent: drift.Equivalent},
			))
		},
		Entry("equal strings", "foo", "foo"),
		Entry("equal ints", 42, 42),
		Entry("both nil", nil, nil),
		Entry("equal slices", []string{"a", "b"}, []string{"a", "b"}),
	)

	DescribeTable("should report inequivalence",
		func(crd, remote any) {
			Expect(drift.DefaultEquivalence(crd, remote, drift.DriftContext{})).To(Equal(
				drift.Equivalence{Equivalent: drift.Inequivalent},
			))
		},
		Entry("different strings", "foo", "bar"),
		Entry("different ints", 1, 2),
		Entry("string vs nil", "foo", nil),
		Entry("nil vs string", nil, "foo"),
		Entry("different slices", []string{"a"}, []string{"b"}),
	)
})

var _ = Describe("Ignore", func() {
	It("always reports equivalence regardless of values", func() {
		Expect(drift.Ignore("export-id", "", drift.DriftContext{})).To(Equal(
			drift.Equivalence{Equivalent: drift.CannotCompare},
		))
		Expect(drift.Ignore(nil, "foo", drift.DriftContext{})).To(Equal(
			drift.Equivalence{Equivalent: drift.CannotCompare},
		))
	})
})

var _ = Describe("IgnoreSkip", func() {
	It("always reports equivalence and skips children", func() {
		Expect(drift.IgnoreSkip(struct{}{}, struct{ Name string }{Name: "foo"}, drift.DriftContext{})).To(Equal(
			drift.Equivalence{Equivalent: drift.CannotCompare, Skip: true},
		))
		Expect(drift.IgnoreSkip([]string{}, []string{"foo"}, drift.DriftContext{})).To(Equal(
			drift.Equivalence{Equivalent: drift.CannotCompare, Skip: true},
		))
	})
})

var _ = Describe("Trimmed", func() {
	DescribeTable("should report equivalence",
		func(crd, remote string) {
			Expect(drift.Trimmed(crd, remote, drift.DriftContext{})).To(Equal(
				drift.Equivalence{Equivalent: drift.Equivalent},
			))
		},
		Entry("same value", "my-api", "my-api"),
		Entry("leading and trailing whitespace on crd", "  my-api  ", "my-api"),
		Entry("leading and trailing whitespace on remote", "my-api", "  my-api  "),
		Entry("whitespace on both sides", "  my-api  ", "  my-api  "),
		Entry("empty strings", "", ""),
	)

	DescribeTable("should report inequivalence",
		func(crd, remote string) {
			Expect(drift.Trimmed(crd, remote, drift.DriftContext{})).To(Equal(
				drift.Equivalence{Equivalent: drift.Inequivalent},
			))
		},
		Entry("different values", "my-api", "other-api"),
		Entry("whitespace-only crd vs non-empty remote", "   ", "my-api"),
	)
})

var _ = Describe("RFC3339", func() {
	const (
		certStartsAt = "2024-06-15T10:30:00Z"
		certEndsAt   = "2025-06-15T10:30:00+00:00"
	)

	DescribeTable("should report equivalence",
		func(crd, remote string) {
			Expect(drift.RFC3339(crd, remote, drift.DriftContext{})).To(Equal(
				drift.Equivalence{Equivalent: drift.Equivalent},
			))
		},
		Entry("same RFC3339 timestamp", certStartsAt, certStartsAt),
		Entry("same instant with different offset notation", certEndsAt, "2025-06-15T10:30:00Z"),
		Entry("both empty certificate dates", "", ""),
	)

	DescribeTable("should report inequivalence",
		func(crd, remote string) {
			Expect(drift.RFC3339(crd, remote, drift.DriftContext{})).To(Equal(
				drift.Equivalence{Equivalent: drift.Inequivalent},
			))
		},
		Entry("different certificate validity dates", certStartsAt, certEndsAt),
		Entry("set startsAt vs empty", certStartsAt, ""),
		Entry("empty vs set endsAt", "", certEndsAt),
	)

	It("reports a parse error for invalid crd timestamp", func() {
		e := drift.RFC3339("not-a-date", certStartsAt, drift.DriftContext{})
		Expect(e.Equivalent).To(Equal(drift.Inequivalent))
		Expect(e.Reason).To(BeAssignableToTypeOf(&time.ParseError{}))
	})

	It("reports a parse error for invalid remote timestamp", func() {
		e := drift.RFC3339(certStartsAt, "not-a-date", drift.DriftContext{})
		Expect(e.Equivalent).To(Equal(drift.Inequivalent))
		Expect(e.Reason).To(BeAssignableToTypeOf(&time.ParseError{}))
	})
})

var _ = Describe("DefaultEquivalencePostPullUpObjectChildren", func() {
	DescribeTable("delegates struct equivalence to defaultStructEquivalence",
		func(crd, remote any, expected drift.Equivalence) {
			e := drift.DefaultEquivalencePostPullUpObjectChildren(crd, remote, drift.DriftContext{})
			Expect(e.Equivalent).To(Equal(expected.Equivalent))
			Expect(e.Skip).To(Equal(expected.Skip))
			Expect(e.PostFunc).NotTo(BeNil())
		},
		Entry("nil config vs populated GenericStringMap",
			nil,
			unstructured.Unstructured{Object: map[string]any{"type": "HTTP"}},
			drift.Equivalence{Equivalent: drift.CannotCompare, Skip: false}),
		Entry("populated config vs nil",
			unstructured.Unstructured{Object: map[string]any{"type": "HTTP"}},
			nil, drift.Equivalence{Equivalent: drift.CannotCompare, Skip: false}),
		Entry("both populated configs",
			unstructured.Unstructured{Object: map[string]any{"type": "HTTP"}},
			unstructured.Unstructured{Object: map[string]any{"type": "HTTP"}},
			drift.Equivalence{Equivalent: drift.CannotCompare}),
	)

	It("hoists object children to the root via PostFunc", func() {
		e := drift.DefaultEquivalencePostPullUpObjectChildren(unstructured.Unstructured{Object: map[string]any{"type": "HTTP"}}, unstructured.Unstructured{Object: map[string]any{"type": "HTTP"}}, drift.DriftContext{})

		child := &drift.Result{}
		child.AppendChild(&drift.Result{Property: "other"}, false)
		o := child.AppendChild(&drift.Result{Property: "object"}, false)
		o.AppendChild(&drift.Result{Property: "type", CRDValue: "HTTP"}, false)
		o.AppendChild(&drift.Result{Property: "path", CRDValue: "/api"}, false)
		e.PostFunc(child)

		Expect(child.Children()).To(ConsistOf(
			&drift.Result{Property: "other"},
			&drift.Result{Property: "path", CRDValue: "/api"},
			&drift.Result{Property: "type", CRDValue: "HTTP"},
		))
	})

	It("ignores an empty object child", func() {
		e := drift.DefaultEquivalencePostPullUpObjectChildren(nil, nil, drift.DriftContext{})

		child := &drift.Result{}
		child.AppendChild(&drift.Result{Property: "object"}, false)
		child.AppendChild(&drift.Result{Property: "remaining"}, false)
		e.PostFunc(child)

		Expect(child.Children()).To(ConsistOf(
			&drift.Result{Property: "remaining"},
		))
	})

	Describe("RFC3339 drift tags", func() {

		// withRFC3339Dates mimics application.ClientCertificate date fields.
		type withRFC3339Dates struct {
			StartsAt string `json:"startsAt" drift:"rfc3339"`
			EndsAt   string `json:"endsAt" drift:"rfc3339"`
		}

		const validDate = "2024-06-15T10:30:00Z"

		It("reports parse error through Detect and Result.String for invalid crd timestamp", func() {
			crd := withRFC3339Dates{StartsAt: "not-a-date", EndsAt: validDate}
			remote := withRFC3339Dates{StartsAt: validDate, EndsAt: validDate}
			result := drift.DetectWithNamespace(crd, remote, "")
			Expect(result.DriftDetected()).To(BeTrue())
			Expect(result.String()).To(MatchRegexp(
				`startsAt: "not-a-date" != "` + validDate + `" \(error: parsing time`,
			))
		})

		It("reports parse error through Detect and Result.String for invalid remote timestamp", func() {
			crd := withRFC3339Dates{StartsAt: validDate, EndsAt: validDate}
			remote := withRFC3339Dates{StartsAt: validDate, EndsAt: "not-a-date"}
			result := drift.DetectWithNamespace(crd, remote, "")
			Expect(result.DriftDetected()).To(BeTrue())
			Expect(result.String()).To(MatchRegexp(
				`endsAt: "` + validDate + `" != "not-a-date" \(error: parsing time`,
			))
		})

		It("detects no drift for equivalent instants with different RFC3339 notations", func() {
			crd := withRFC3339Dates{StartsAt: validDate, EndsAt: "2025-06-15T10:30:00+00:00"}
			remote := withRFC3339Dates{StartsAt: validDate, EndsAt: "2025-06-15T10:30:00Z"}
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		})
	})
})

type namedItem struct{ ID string }

func (n namedItem) MatchKey() string { return n.ID }

type withIgnoreOnlyRemote struct {
	Items []namedItem `json:"items" drift:"ignore-only:remote"`
}

type withIgnoreOnlyCRD struct {
	Items []namedItem `json:"items" drift:"ignore-only:crd"`
}

type withIgnoreOnlyCRDStripNS struct {
	Items []namedItem `json:"items" drift:"ignore-only:crd,strip-ns"`
}

var _ = Describe("IgnoreOnlyArgs", func() {
	remoteCtx := drift.DriftContext{FuncArgs: []string{"remote"}}
	crdCtx := drift.DriftContext{FuncArgs: []string{"crd"}}

	DescribeTable("should report equivalence and skip for empty slices",
		func(crd, remote any) {
			Expect(drift.IgnoreOnlyArgs(crd, remote, remoteCtx)).To(Equal(
				drift.Equivalence{Equivalent: drift.Equivalent, Skip: true},
			))
		},
		Entry("nil vs nil", nil, nil),
		Entry("nil vs empty slice", nil, []namedItem{}),
		Entry("empty slice vs nil", []namedItem{}, nil),
		Entry("empty slice vs empty slice", []namedItem{}, []namedItem{}),
	)

	DescribeTable("should report cannot compare without filter",
		func(crd, remote any, ctx drift.DriftContext) {
			e := drift.IgnoreOnlyArgs(crd, remote, ctx)
			Expect(e.Equivalent).To(Equal(drift.CannotCompare))
			Expect(e.RemoteItemsFilterFunc).To(BeNil())
			Expect(e.CRDItemsFilterFunc).To(BeNil())
		},
		Entry("same names", []namedItem{{ID: "owner"}}, []namedItem{{ID: "owner"}}, remoteCtx),
		Entry("non-keyed slice items", []string{"foo"}, []string{"foo"}, remoteCtx),
		Entry("crd item missing from remote", []namedItem{{ID: "owner"}}, []namedItem{}, remoteCtx),
		Entry("empty context does not filter",
			[]namedItem{{ID: "owner"}},
			[]namedItem{{ID: "owner"}, {ID: "sync-id"}},
			drift.DriftContext{},
		),
	)

	It("provides a filter that removes remote-only items by identifier", func() {
		crd := []namedItem{{ID: "owner"}}
		remote := []namedItem{{ID: "owner"}, {ID: "sync-id"}}

		e := drift.IgnoreOnlyArgs(crd, remote, remoteCtx)
		Expect(e.Equivalent).To(Equal(drift.CannotCompare))
		Expect(e.RemoteItemsFilterFunc).NotTo(BeNil())
		Expect(e.CRDItemsFilterFunc).To(BeNil())

		filtered := e.RemoteItemsFilterFunc(remote)
		Expect(filtered).To(ConsistOf(namedItem{ID: "owner"}))
	})

	It("filters remote-only items when crd is empty", func() {
		remote := []namedItem{{ID: "sync-id"}}

		e := drift.IgnoreOnlyArgs(nil, remote, remoteCtx)
		Expect(e.Equivalent).To(Equal(drift.CannotCompare))
		Expect(e.RemoteItemsFilterFunc).NotTo(BeNil())

		filtered := e.RemoteItemsFilterFunc(remote)
		Expect(filtered).To(BeEmpty())
	})

	It("keeps non-keyed items when filtering remote-only items", func() {
		crd := []namedItem{{ID: "owner"}}
		remote := []namedItem{{ID: "owner"}, {ID: "sync-id"}}

		e := drift.IgnoreOnlyArgs(crd, remote, remoteCtx)
		Expect(e.RemoteItemsFilterFunc).NotTo(BeNil())

		mixed := []any{
			namedItem{ID: "owner"},
			namedItem{ID: "sync-id"},
			"keep-me",
		}
		filtered := e.RemoteItemsFilterFunc(mixed)
		Expect(filtered).To(ConsistOf(namedItem{ID: "owner"}, "keep-me"))
	})

	It("provides a filter that removes crd-only items by identifier", func() {
		crd := []namedItem{{ID: "owner"}, {ID: "local-only"}}
		remote := []namedItem{{ID: "owner"}}

		e := drift.IgnoreOnlyArgs(crd, remote, crdCtx)
		Expect(e.Equivalent).To(Equal(drift.CannotCompare))
		Expect(e.CRDItemsFilterFunc).NotTo(BeNil())
		Expect(e.RemoteItemsFilterFunc).To(BeNil())

		filtered := e.CRDItemsFilterFunc(crd)
		Expect(filtered).To(ConsistOf(namedItem{ID: "owner"}))
	})

	It("filters crd-only items when remote is empty", func() {
		crd := []namedItem{{ID: "local-only"}}

		e := drift.IgnoreOnlyArgs(crd, nil, crdCtx)
		Expect(e.Equivalent).To(Equal(drift.CannotCompare))
		Expect(e.CRDItemsFilterFunc).NotTo(BeNil())

		filtered := e.CRDItemsFilterFunc(crd)
		Expect(filtered).To(BeEmpty())
	})

	It("Detect ignores remote-only items", func() {
		crd := withIgnoreOnlyRemote{Items: []namedItem{{ID: "owner"}}}
		remote := withIgnoreOnlyRemote{Items: []namedItem{{ID: "owner"}, {ID: "sync-id"}}}
		expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
	})

	It("Detect ignores crd-only items", func() {
		crd := withIgnoreOnlyCRD{Items: []namedItem{{ID: "owner"}, {ID: "local-only"}}}
		remote := withIgnoreOnlyCRD{Items: []namedItem{{ID: "owner"}}}
		expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
	})

	It("Detect ignores namespace prefix with strip-ns", func() {
		crd := withIgnoreOnlyCRDStripNS{Items: []namedItem{{ID: "group1"}, {ID: "crd-only"}}}
		remote := withIgnoreOnlyCRDStripNS{Items: []namedItem{{ID: "test-namespace-group1"}}}
		expectNoDrift(drift.DetectWithNamespace(crd, remote, "test-namespace"))
	})
})

type timedItem struct {
	ID          string `json:"id"`
	IsExpired   bool   `json:"-"`
	IsScheduled bool   `json:"-"`
}

func (t timedItem) MatchKey() string { return t.ID }
func (t timedItem) Expired() bool    { return t.IsExpired }
func (t timedItem) Scheduled() bool  { return t.IsScheduled }

type withIgnoreOnlyCRDHidden struct {
	Items []timedItem `json:"items" drift:"ignore-only:crd,expired,scheduled"`
}

var _ = Describe("IgnoreOnlyArgs expired and scheduled", func() {
	hiddenCtx := drift.DriftContext{FuncArgs: []string{"crd", "expired", "scheduled"}}

	It("skips comparison when every item is expired or scheduled", func() {
		crd := []timedItem{{ID: "old", IsExpired: true}, {ID: "later", IsScheduled: true}}
		e := drift.IgnoreOnlyArgs(crd, []timedItem{}, hiddenCtx)
		Expect(e.Equivalent).To(Equal(drift.Equivalent))
		Expect(e.Skip).To(BeTrue())
	})

	It("Detect ignores an expired CRD item omitted by remote", func() {
		crd := withIgnoreOnlyCRDHidden{Items: []timedItem{
			{ID: "expired", IsExpired: true},
			{ID: "active"},
		}}
		remote := withIgnoreOnlyCRDHidden{Items: []timedItem{{ID: "active"}}}
		expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
	})

	It("Detect ignores a scheduled CRD item omitted by remote", func() {
		crd := withIgnoreOnlyCRDHidden{Items: []timedItem{{ID: "scheduled", IsScheduled: true}}}
		remote := withIgnoreOnlyCRDHidden{Items: []timedItem{}}
		expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
	})

	It("Detect ignores expired and scheduled items together", func() {
		crd := withIgnoreOnlyCRDHidden{Items: []timedItem{
			{ID: "expired", IsExpired: true},
			{ID: "scheduled", IsScheduled: true},
			{ID: "active"},
		}}
		remote := withIgnoreOnlyCRDHidden{Items: []timedItem{{ID: "active"}}}
		expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
	})

	It("Detect still reports drift on remaining active items", func() {
		crd := withIgnoreOnlyCRDHidden{Items: []timedItem{
			{ID: "expired", IsExpired: true},
			{ID: "scheduled", IsScheduled: true},
			{ID: "active"},
		}}
		remote := withIgnoreOnlyCRDHidden{Items: []timedItem{{ID: "other"}}}
		got := drift.DetectWithNamespace(crd, remote, "")
		Expect(got.DriftDetected()).To(BeTrue())
	})
})
