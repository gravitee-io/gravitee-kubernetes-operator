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

package drift_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
)

var _ = Describe("EmptyIsNilString", func() {
	DescribeTable("should report equivalence",
		func(crd, api any) {
			Expect(drift.EmptyIsNilString(crd, api)).To(Equal(
				drift.Equivalence{Equivalent: drift.Equivalent},
			))
		},
		Entry("empty string vs nil", "", nil),
		Entry("nil vs empty string", nil, ""),
		Entry("empty string vs empty string", "", ""),
		Entry("nil vs nil", nil, nil),
	)

	DescribeTable("should report inequivalence",
		func(crd, api any) {
			Expect(drift.EmptyIsNilString(crd, api)).To(Equal(
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
		func(crd, api any) {
			Expect(drift.EmptyIsNilInt(crd, api)).To(Equal(
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
		func(crd, api any) {
			Expect(drift.EmptyIsNilInt(crd, api)).To(Equal(
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
		func(crd, api any) {
			Expect(drift.EmptyIsNilUint(crd, api)).To(Equal(
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
		func(crd, api any) {
			Expect(drift.EmptyIsNilUint(crd, api)).To(Equal(
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
		func(crd, api any) {
			Expect(drift.EmptyIsNilBool(crd, api)).To(Equal(
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
		func(crd, api any) {
			Expect(drift.EmptyIsNilBool(crd, api)).To(Equal(
				drift.Equivalence{Equivalent: drift.Inequivalent},
			))
		},
		Entry("true vs nil", true, nil),
		Entry("nil vs true", nil, true),
		Entry("true vs false", true, false),
		Entry("false vs true", false, true),
	)
})

var _ = Describe("EmptyIsNilLen", func() {
	DescribeTable("should report equivalence",
		func(crd, api any) {
			Expect(drift.EmptyIsNilLen(crd, api)).To(Equal(
				drift.Equivalence{Equivalent: drift.Equivalent},
			))
		},
		Entry("empty slice vs empty slice", []string{}, []string{}),
		Entry("nil slice vs empty slice", []string(nil), []string{}),
		Entry("slice with same length, different content", []string{"foo"}, []string{"bar"}),
		Entry("empty map vs empty map", map[string]string{}, map[string]string{}),
		Entry("nil map vs empty map", map[string]string(nil), map[string]string{}),
		Entry("map with same length, different content", map[string]string{"a": "1"}, map[string]string{"b": "2"}),
	)

	DescribeTable("should report cannot compare",
		func(crd, api any) {
			Expect(drift.EmptyIsNilLen(crd, api)).To(Equal(
				drift.Equivalence{Equivalent: drift.CannotCompare},
			))
		},
		Entry("slice with different lengths", []string{"foo"}, []string{}),
		Entry("slice with different lengths, longer crd", []string{"foo", "bar"}, []string{"baz"}),
		Entry("map with different lengths", map[string]string{"a": "1"}, map[string]string{}),
		Entry("map with different lengths, longer crd", map[string]string{"a": "1", "b": "2"}, map[string]string{"c": "3"}),
	)
})

type emptyIsNilStructFixture struct {
	Name  string
	Count int
}

var _ = Describe("EmptyIsNilStruct", func() {
	DescribeTable("should report equivalence and skip",
		func(crd, api any) {
			Expect(drift.EmptyIsNilStruct(crd, api)).To(Equal(
				drift.Equivalence{Equivalent: drift.Equivalent, Skip: true},
			))
		},
		Entry("nil vs zero struct", nil, emptyIsNilStructFixture{}),
		Entry("zero struct vs nil", emptyIsNilStructFixture{}, nil),
	)

	DescribeTable("should report cannot compare",
		func(crd, api any) {
			Expect(drift.EmptyIsNilStruct(crd, api)).To(Equal(
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
		func(crd, api any) {
			Expect(drift.FromDeepEqual(crd, api)).To(Equal(
				drift.Equivalence{Equivalent: drift.Equivalent},
			))
		},
		Entry("equal strings", "foo", "foo"),
		Entry("equal ints", 42, 42),
		Entry("both nil", nil, nil),
		Entry("equal slices", []string{"a", "b"}, []string{"a", "b"}),
	)

	DescribeTable("should report inequivalence",
		func(crd, api any) {
			Expect(drift.FromDeepEqual(crd, api)).To(Equal(
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
