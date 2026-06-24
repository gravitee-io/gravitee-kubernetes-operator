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

package drift

import (
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
)

type SingleWithPtr struct {
	Value *string `drift:"empty-is-nil" json:"value,omitempty"`
}

type MultipleWithPtr struct {
	Name  *string `drift:"empty-is-nil" json:"name,omitempty"`
	Order *int    `drift:"empty-is-nil" json:"order,omitempty"`
}

type Nested struct {
	Name        *string `drift:"empty-is-nil" json:"name,omitempty"`
	Multiple    MultipleWithPtr
	Description *string `drift:"empty-is-nil" json:"description,omitempty"`
}

type NestedWithPointer struct {
	Title       *string          `drift:"empty-is-nil" json:"title,omitempty"`
	Multiple    *MultipleWithPtr `drift:"empty-is-nil" json:"multiple,omitempty"`
	Description *string          `drift:"empty-is-nil" json:"description,omitempty"`
}

type nestedDeep struct {
	Title  string
	First  Nested
	Second *Nested
}
type nestedDeepWithPtr struct {
	Tag    string
	Nested *NestedWithPointer `drift:"empty-is-nil" json:"nested,omitempty"`
}

type Embedded struct {
	MultipleWithPtr `json:",inline"`
	Description     *string `drift:"empty-is-nil" json:"description,omitempty"`
}

type embeddedWithPtr struct {
	*MultipleWithPtr `json:",inline"`
	Description      *string `drift:"empty-is-nil" json:"description,omitempty"`
}

type TitleString string

type embeddedAndNested struct {
	Title          TitleString
	Embedded       `json:",inline"`
	NestedEmbedded Embedded
}

type singleNoPtr struct {
	Value string `drift:"empty-is-nil" json:"value,omitempty"`
}

type singleNoDriftTag struct {
	Value string `json:"value,omitempty"`
}

type singleNoTags struct {
	Value string
}

type singleNilablePtr struct {
	Value *string `drift:"empty-is-nil"`
}

type withStringArray struct {
	Name string
	Tags []string
}

type withStructArray struct {
	Name   *string           `json:"name" drift:"empty-is-nil"`
	Values []MultipleWithPtr `json:"values" drift:"empty-is-nil"`
}

type withMapOfPrimitive struct {
	Name   *string         `json:"name" drift:"empty-is-nil"`
	Values map[string]*int `json:"values" drift:"empty-is-nil"`
}
type withMapOfStruct struct {
	Name   *string                    `json:"name" drift:"empty-is-nil"`
	Values map[string]MultipleWithPtr `json:"values" drift:"empty-is-nil"`
}

type withMapOfMap struct {
	Name   *string                   `json:"name" drift:"empty-is-nil"`
	Values map[string]map[string]int `json:"values" drift:"empty-is-nil"`
}

type withMapOfArray struct {
	Name   *string          `json:"name" drift:"empty-is-nil"`
	Values map[string][]int `json:"values" drift:"empty-is-nil"`
}

type withMapOfStructArray struct {
	Name   *string                      `json:"name" drift:"empty-is-nil"`
	Values map[string][]MultipleWithPtr `json:"values" drift:"empty-is-nil"`
}

type withIntMapKey struct {
	Values map[int]string `json:"values" drift:"empty-is-nil"`
}

type withNestedIntMapKey struct {
	Values map[string]map[int]string `json:"values" drift:"empty-is-nil"`
}

type withAnyField struct {
	Value any `json:"value" drift:"type-check-panic"`
}

var _ = Describe("Detect", func() {

	Describe("when drift is detected", func() {

		DescribeTable("simple properties",
			func(crd, api any, expectedString string) {
				expectDrift(drift.Detect(crd, api), expectedString)
			},
			Entry("single pointer property with different values",
				SingleWithPtr{Value: ptr("foo")},
				SingleWithPtr{Value: ptr("bar")},
				`value: "foo" != "bar"`,
			),
			Entry("single non-pointer property with different values",
				singleNoPtr{Value: "foo"},
				singleNoPtr{Value: "bar"},
				`value: "foo" != "bar"`,
			),
			Entry("single property without drift tag falls back to default equal",
				singleNoDriftTag{Value: "foo"},
				singleNoDriftTag{Value: "bar"},
				`value: "foo" != "bar"`,
			),
			Entry("single property without any tag falls back to default equal",
				singleNoTags{Value: "foo"},
				singleNoTags{Value: "bar"},
				`value: "foo" != "bar"`,
			),
			Entry("double pointer property with one difference",
				MultipleWithPtr{Name: ptr("foo"), Order: ptr(1)},
				MultipleWithPtr{Name: ptr("foo"), Order: ptr(2)},
				"order: 1 != 2",
			),
			Entry("double pointer property with all different",
				MultipleWithPtr{Name: ptr("foo"), Order: ptr(1)},
				MultipleWithPtr{Name: ptr("bar"), Order: ptr(2)},
				`name: "foo" != "bar"
order: 1 != 2`,
			),
		)

		Describe("nested structures", func() {

			It("detects partial drift", func() {
				crd := Nested{Name: ptr("foo"), Multiple: MultipleWithPtr{Name: ptr("foo"), Order: ptr(1)}, Description: ptr("bar")}
				api := Nested{Name: ptr("foo"), Multiple: MultipleWithPtr{Name: ptr("foo"), Order: ptr(2)}, Description: ptr("bar")}
				expectDrift(drift.Detect(crd, api), `multiple:
  order: 1 != 2`)
			})

			It("detects drift on all props", func() {
				crd := Nested{Name: ptr("foo"), Multiple: MultipleWithPtr{Name: ptr("foo"), Order: ptr(1)}, Description: ptr("foo")}
				api := Nested{Name: ptr("bar"), Multiple: MultipleWithPtr{Name: ptr("bar"), Order: ptr(2)}, Description: ptr("bar")}
				expectDrift(drift.Detect(crd, api), `name: "foo" != "bar"
multiple:
  name: "foo" != "bar"
  order: 1 != 2
description: "foo" != "bar"`)
			})

			It("detects drift on complex structure", func() {
				crd := nestedDeep{Title: "CRD",
					First:  Nested{Name: ptr("First CRD"), Multiple: MultipleWithPtr{Name: ptr("foo"), Order: ptr(1)}},
					Second: &Nested{Name: ptr("Second CRD"), Multiple: MultipleWithPtr{Name: ptr("bar"), Order: ptr(1)}, Description: ptr("Second CRD")},
				}
				api := nestedDeep{Title: "API",
					First:  Nested{Name: ptr("First API"), Multiple: MultipleWithPtr{Name: ptr("baz"), Order: ptr(1)}},
					Second: &Nested{Name: ptr("Second API"), Multiple: MultipleWithPtr{Name: ptr("bar"), Order: ptr(2)}, Description: ptr("Second API")},
				}
				expectDrift(drift.Detect(crd, api), `title: "CRD" != "API"
first:
  name: "First CRD" != "First API"
  multiple:
    name: "foo" != "baz"
second:
  name: "Second CRD" != "Second API"
  multiple:
    order: 1 != 2
  description: "Second CRD" != "Second API"`)
			})

			It("detects drifts in nested with pointers", func() {
				crd := NestedWithPointer{Title: ptr("test"), Multiple: &MultipleWithPtr{Name: ptr("foo")}, Description: ptr("Test")}
				api := NestedWithPointer{Title: ptr("test"), Multiple: nil, Description: ptr("Test")}
				expectDrift(drift.Detect(crd, api), `multiple:
  name: "foo" != <nil>`)
			})

			It("detects drifts in deep nested with structs with pointer", func() {
				crd := nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{Title: ptr(""), Multiple: nil, Description: nil}}
				api := nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{Title: nil, Multiple: &MultipleWithPtr{
					Name:  ptr("x"),
					Order: ptr(1),
				}, Description: ptr("")}}
				expectDrift(drift.Detect(crd, api), `nested:
  multiple:
    name: <nil> != "x"
    order: <nil> != 1`)
			})

			It("detects drifts in deep nested with structs with pointer", func() {
				crd := nestedDeepWithPtr{Tag: "test", Nested: nil}
				api := nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{Title: nil, Multiple: &MultipleWithPtr{
					Name:  ptr("x"),
					Order: ptr(1),
				}, Description: ptr("")}}
				expectDrift(drift.Detect(crd, api), `nested:
  multiple:
    name: <nil> != "x"
    order: <nil> != 1`)
			})

			It("detects drifts in deep nested with structs with pointer", func() {
				crd := nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{Title: nil, Multiple: &MultipleWithPtr{
					Name:  ptr("x"),
					Order: ptr(1),
				}, Description: ptr("")}}
				api := nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{Title: ptr(""), Multiple: nil, Description: nil}}
				expectDrift(drift.Detect(crd, api), `nested:
  multiple:
    name: "x" != <nil>
    order: 1 != <nil>`)
			})

			It("detects drifts in deep nested with structs with pointer", func() {
				crd := nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{Title: nil, Multiple: &MultipleWithPtr{
					Name:  ptr("x"),
					Order: ptr(1),
				}, Description: ptr("")}}
				api := nestedDeepWithPtr{Tag: "test", Nested: nil}
				expectDrift(drift.Detect(crd, api), `nested:
  multiple:
    name: "x" != <nil>
    order: 1 != <nil>`)
			})

		})

		Describe("embedded structures", func() {

			It("detects partial drift", func() {
				crd := Embedded{MultipleWithPtr: MultipleWithPtr{Name: ptr("foo"), Order: ptr(1)}, Description: ptr("bar")}
				api := Embedded{MultipleWithPtr: MultipleWithPtr{Name: ptr("foo"), Order: ptr(2)}, Description: ptr("bar")}
				expectDrift(drift.Detect(crd, api), "order: 1 != 2")
			})

			It("detects drift on all props", func() {
				crd := Embedded{MultipleWithPtr: MultipleWithPtr{Name: ptr("foo"), Order: ptr(1)}, Description: ptr("foo")}
				api := Embedded{MultipleWithPtr: MultipleWithPtr{Name: ptr("bar"), Order: ptr(2)}, Description: ptr("bar")}
				expectDrift(drift.Detect(crd, api), `name: "foo" != "bar"
order: 1 != 2
description: "foo" != "bar"`)
			})
			It("detects partial drift", func() {
				crd := embeddedWithPtr{MultipleWithPtr: &MultipleWithPtr{Name: ptr("foo"), Order: ptr(1)}, Description: ptr("bar")}
				api := embeddedWithPtr{MultipleWithPtr: &MultipleWithPtr{Name: ptr("foo"), Order: ptr(2)}, Description: ptr("bar")}
				expectDrift(drift.Detect(crd, api), "order: 1 != 2")
			})

			It("detects drift on all props", func() {
				crd := embeddedWithPtr{MultipleWithPtr: &MultipleWithPtr{Name: ptr("foo"), Order: ptr(1)}, Description: ptr("foo")}
				api := embeddedWithPtr{MultipleWithPtr: &MultipleWithPtr{Name: ptr("bar"), Order: ptr(2)}, Description: ptr("bar")}
				expectDrift(drift.Detect(crd, api), `name: "foo" != "bar"
order: 1 != 2
description: "foo" != "bar"`)
			})

			It("detects drift with nested and embedded", func() {
				crd := embeddedAndNested{Title: "Test", Embedded: Embedded{
					MultipleWithPtr: MultipleWithPtr{Name: ptr("x"), Order: ptr(10)}, Description: ptr("coming soon"),
				},
					NestedEmbedded: Embedded{MultipleWithPtr: MultipleWithPtr{Name: ptr("foo")}, Description: nil}}
				api := embeddedAndNested{Title: "Test", Embedded: Embedded{
					MultipleWithPtr: MultipleWithPtr{Name: ptr("y"), Order: ptr(42)}, Description: ptr("coming soon"),
				},
					NestedEmbedded: Embedded{MultipleWithPtr: MultipleWithPtr{Name: ptr("bar")}, Description: ptr("baz")}}
				expectDrift(drift.Detect(crd, api), `name: "x" != "y"
order: 10 != 42
nestedEmbedded:
  name: "foo" != "bar"
  description: <nil> != "baz"`)
			})

		})

		Describe("struct with slices", func() {

			It("detects array drift", func() {
				crd := withStringArray{Name: "test", Tags: []string{"bar", "bar"}}
				api := withStringArray{Name: "test", Tags: []string{"foo", "bar"}}
				expectDrift(drift.Detect(crd, api), `tags[0]: "bar" != "foo"`)
			})

			It("detects array drift", func() {
				crd := withStringArray{Name: "test", Tags: []string{"bar", "bar"}}
				api := withStringArray{Name: "test", Tags: []string{"foo"}}
				expectDrift(drift.Detect(crd, api), `tags[0]: "bar" != "foo"
tags[1]: "bar" != ""`)
			})

			It("detects array drift", func() {
				crd := withStringArray{Name: "test", Tags: []string{"bar"}}
				api := withStringArray{Name: "test", Tags: []string{"foo", "bar"}}
				expectDrift(drift.Detect(crd, api), `tags[0]: "bar" != "foo"
tags[1]: "" != "bar"`)
			})

			It("detects array drift", func() {
				crd := withStringArray{Name: "test"}
				api := withStringArray{Name: "test", Tags: []string{"foo", "bar"}}
				expectDrift(drift.Detect(crd, api), `tags[0]: "" != "foo"
tags[1]: "" != "bar"`)
			})

			It("detects array of struct drift against nil", func() {
				crd := withStructArray{Name: ptr("test"), Values: []MultipleWithPtr{{Name: ptr("foo")}}}
				api := withStructArray{Name: ptr("test"), Values: nil}
				expectDrift(drift.Detect(crd, api), `values[0]:
  name: "foo" != <nil>`)
			})

			It("detects array of struct drift against nil", func() {
				crd := withStructArray{Name: ptr("test"), Values: nil}
				api := withStructArray{Name: ptr("test"), Values: []MultipleWithPtr{{Name: ptr("foo")}}}
				expectDrift(drift.Detect(crd, api), `values[0]:
  name: <nil> != "foo"`)
			})

			It("detects array of struct drift complex setup", func() {
				crd := withStructArray{
					Name: ptr("test"),
					Values: []MultipleWithPtr{
						{Name: ptr("bar")},
						{Name: ptr("foo"), Order: ptr(1)},
					},
				}
				api := withStructArray{
					Name: nil,
					Values: []MultipleWithPtr{
						{Name: ptr("foo"), Order: ptr(0)},
						{Name: ptr("bar")},
					},
				}
				expectDrift(drift.Detect(crd, api), `name: "test" != <nil>
values[0]:
  name: "bar" != "foo"
values[1]:
  name: "foo" != "bar"
  order: 1 != <nil>`)
			})
		})

		Describe("struct with int maps", func() {
			It("detects map drift, partial", func() {
				crd := withMapOfPrimitive{Name: ptr("test"), Values: map[string]*int{"x": ptr(42)}}
				api := withMapOfPrimitive{Name: ptr("test"), Values: map[string]*int{"x": ptr(0)}}
				expectDrift(drift.Detect(crd, api), `values:
  x: 42 != 0`)
			})

			It("detects int map drift, complete", func() {
				crd := withMapOfPrimitive{Name: ptr("foo"), Values: map[string]*int{"x": ptr(42)}}
				api := withMapOfPrimitive{Name: ptr("bar"), Values: map[string]*int{"x": ptr(0)}}
				expectDrift(drift.Detect(crd, api), `name: "foo" != "bar"
values:
  x: 42 != 0`)
			})

			It("detects map drift, partial api missing", func() {
				crd := withMapOfPrimitive{Name: ptr("test"), Values: map[string]*int{"x": ptr(42)}}
				api := withMapOfPrimitive{Name: ptr("test"), Values: map[string]*int{}}
				expectDrift(drift.Detect(crd, api), `values:
  x: 42 != <nil>`)
			})

			It("detects map drift, partial crd missing", func() {
				crd := withMapOfPrimitive{Name: ptr("test"), Values: map[string]*int{}}
				api := withMapOfPrimitive{Name: ptr("test"), Values: map[string]*int{"x": ptr(42)}}
				expectDrift(drift.Detect(crd, api), `values:
  x: <nil> != 42`)
			})
			It("detects map drift, partial crd nil", func() {
				crd := withMapOfPrimitive{Name: ptr("test")}
				api := withMapOfPrimitive{Name: ptr("test"), Values: map[string]*int{"x": ptr(42)}}
				expectDrift(drift.Detect(crd, api), `values:
  x: <nil> != 42`)
			})
			It("detects map drift, partial api nil", func() {
				crd := withMapOfPrimitive{Name: ptr("test"), Values: map[string]*int{"x": ptr(42)}}
				api := withMapOfPrimitive{Name: ptr("test")}
				expectDrift(drift.Detect(crd, api), `values:
  x: 42 != <nil>`)
			})
			It("detects map drift, partial api nil", func() {
				crd := withMapOfPrimitive{Name: ptr("test"), Values: map[string]*int{"x": ptr(42)}}
				api := withMapOfPrimitive{Name: ptr("test"), Values: map[string]*int{"y": ptr(42)}}
				expectDrift(drift.Detect(crd, api), `values:
  x: 42 != <nil>
  y: <nil> != 42`)
			})

			It(" detect map of struct drift, partial", func() {
				crd := withMapOfStruct{Name: ptr("test"), Values: map[string]MultipleWithPtr{"x": {Name: ptr("foo"), Order: ptr(66)}}}
				api := withMapOfStruct{Name: ptr("test"), Values: map[string]MultipleWithPtr{"x": {Name: ptr("bar"), Order: ptr(66)}}}
				expectDrift(drift.Detect(crd, api), `values:
  x:
    name: "foo" != "bar"`)
			})

			It(" detect map of struct drift, complete", func() {
				crd := withMapOfStruct{Name: ptr("test"), Values: map[string]MultipleWithPtr{"x": {Name: ptr("foo"), Order: ptr(66)}}}
				api := withMapOfStruct{Name: ptr("prod"), Values: map[string]MultipleWithPtr{"x": {Name: ptr("bar"), Order: ptr(66)}}}
				expectDrift(drift.Detect(crd, api), `name: "test" != "prod"
values:
  x:
    name: "foo" != "bar"`)
			})

			It(" detect map of struct drift, empty", func() {
				crd := withMapOfStruct{Name: ptr("test"), Values: map[string]MultipleWithPtr{"x": {Order: ptr(66)}}}
				api := withMapOfStruct{Name: ptr("test"), Values: map[string]MultipleWithPtr{"x": {}}}
				expectDrift(drift.Detect(crd, api), `values:
  x:
    order: 66 != <nil>`)
			})

			It(" detect map of struct drift, missing api", func() {
				crd := withMapOfStruct{Name: ptr("test"), Values: map[string]MultipleWithPtr{"x": {Order: ptr(66)}}}
				api := withMapOfStruct{Name: ptr("test"), Values: map[string]MultipleWithPtr{}}
				expectDrift(drift.Detect(crd, api), `values:
  x:
    order: 66 != <nil>`)
			})
			It(" detect map of struct drift, nil api", func() {
				crd := withMapOfStruct{Name: ptr("test"), Values: map[string]MultipleWithPtr{"x": {Order: ptr(66)}}}
				api := withMapOfStruct{Name: ptr("test")}
				expectDrift(drift.Detect(crd, api), `values:
  x:
    order: 66 != <nil>`)
			})

			It(" detect map of struct drift, missing crd", func() {
				crd := withMapOfStruct{Name: ptr("test"), Values: map[string]MultipleWithPtr{}}
				api := withMapOfStruct{Name: ptr("test"), Values: map[string]MultipleWithPtr{"x": {Order: ptr(66)}}}
				expectDrift(drift.Detect(crd, api), `values:
  x:
    order: <nil> != 66`)
			})
			It(" detect map of struct drift, nil crd", func() {
				crd := withMapOfStruct{Name: ptr("test")}
				api := withMapOfStruct{Name: ptr("test"), Values: map[string]MultipleWithPtr{"x": {Order: ptr(66)}}}
				expectDrift(drift.Detect(crd, api), `values:
  x:
    order: <nil> != 66`)
			})

			It(" detect map of arrays drift, missing crd", func() {
				crd := withMapOfArray{Name: ptr("test"), Values: map[string][]int{"x": {1, 2}}}
				api := withMapOfArray{Name: ptr("test"), Values: map[string][]int{"x": {1, 4, 5}}}
				expectDrift(drift.Detect(crd, api), `values:
  x[1]: 2 != 4
  x[2]: 0 != 5`)
			})

			It(" detect map of arrays drift, partial missing api", func() {
				crd := withMapOfArray{Name: ptr("test"), Values: map[string][]int{"x": {1, 2, 3}}}
				api := withMapOfArray{Name: ptr("test"), Values: map[string][]int{"x": {1, 4}}}
				expectDrift(drift.Detect(crd, api), `values:
  x[1]: 2 != 4
  x[2]: 3 != 0`)
			})

			It(" detect map of struct arrays drift, not equal", func() {
				crd := withMapOfStructArray{Name: ptr("test"), Values: map[string][]MultipleWithPtr{"x": {{Name: ptr("foo"), Order: ptr(66)}}}}
				api := withMapOfStructArray{Name: ptr("test"), Values: map[string][]MultipleWithPtr{"x": {{Name: ptr("bar"), Order: ptr(67)}}}}
				expectDrift(drift.Detect(crd, api), `values:
  x[0]:
    name: "foo" != "bar"
    order: 66 != 67`)
			})

			It(" detect map of struct arrays drift, missing api", func() {
				crd := withMapOfStructArray{Name: ptr("test"), Values: map[string][]MultipleWithPtr{"x": {{Name: ptr("foo"), Order: ptr(66)}}}}
				api := withMapOfStructArray{Name: ptr("test"), Values: map[string][]MultipleWithPtr{}}
				expectDrift(drift.Detect(crd, api), `values:
  x[0]:
    name: "foo" != <nil>
    order: 66 != <nil>`)
			})

			It(" detect map of struct arrays drift, missing crd", func() {
				crd := withMapOfStructArray{Name: ptr("test"), Values: map[string][]MultipleWithPtr{}}
				api := withMapOfStructArray{Name: ptr("test"), Values: map[string][]MultipleWithPtr{"x": {{Name: ptr("foo"), Order: ptr(66)}}}}
				expectDrift(drift.Detect(crd, api), `values:
  x[0]:
    name: <nil> != "foo"
    order: <nil> != 66`)
			})
			It(" detect map of struct arrays drift, missing crd item", func() {
				crd := withMapOfStructArray{Name: ptr("test"), Values: map[string][]MultipleWithPtr{"x": {
					{Name: ptr("foo"), Order: ptr(66)},
				}}}
				api := withMapOfStructArray{Name: ptr("test"), Values: map[string][]MultipleWithPtr{"x": {
					{Name: ptr("foo"), Order: ptr(66)},
					{Name: ptr("foo"), Order: ptr(66)},
				}}}
				expectDrift(drift.Detect(crd, api), `values:
  x[1]:
    name: <nil> != "foo"
    order: <nil> != 66`)
			})
			It(" detect map of struct arrays drift, missing crd entry", func() {
				crd := withMapOfStructArray{Name: ptr("test"), Values: map[string][]MultipleWithPtr{
					"x": {{Name: ptr("foo"), Order: ptr(66)}}}}
				api := withMapOfStructArray{Name: ptr("test"), Values: map[string][]MultipleWithPtr{
					"x": {{Name: ptr("foo"), Order: ptr(66)}},
					"y": {{Name: ptr("foo"), Order: ptr(66)}}}}
				expectDrift(drift.Detect(crd, api), `values:
  y[0]:
    name: <nil> != "foo"
    order: <nil> != 66`)
			})

			It(" detect map of struct arrays drift, partial equivalence equal", func() {
				crd := withMapOfStructArray{Name: ptr("test"), Values: map[string][]MultipleWithPtr{"x": {{Name: ptr(""), Order: ptr(66)}}}}
				api := withMapOfStructArray{Name: ptr("test"), Values: map[string][]MultipleWithPtr{"x": {{Name: nil, Order: ptr(67)}}}}
				expectDrift(drift.Detect(crd, api), `values:
  x[0]:
    order: 66 != 67`)
			})

			It(" detect map of maps drift, partial", func() {
				crd := withMapOfMap{Name: ptr("test"), Values: map[string]map[string]int{"x": {"y": 1}}}
				api := withMapOfMap{Name: ptr("test"), Values: map[string]map[string]int{"x": {"y": 0}}}
				expectDrift(drift.Detect(crd, api), `values:
  x:
    y: 1 != 0`)
			})
			It(" detect map of maps drift, complete", func() {
				crd := withMapOfMap{Name: ptr("test"), Values: map[string]map[string]int{"x": {"y": 1}}}
				api := withMapOfMap{Name: ptr("prod"), Values: map[string]map[string]int{"x": {"y": 0}}}
				expectDrift(drift.Detect(crd, api), `name: "test" != "prod"
values:
  x:
    y: 1 != 0`)
			})

			It(" detect map of maps drift, partial missing api", func() {
				crd := withMapOfMap{Name: ptr("test"), Values: map[string]map[string]int{"x": {"y": 1}}}
				api := withMapOfMap{Name: ptr("test"), Values: map[string]map[string]int{"x": {}}}
				expectDrift(drift.Detect(crd, api), `values:
  x:
    y: 1 != <nil>`)
			})

			It(" detect map of maps drift, partial missing crd", func() {
				crd := withMapOfMap{Name: ptr("test"), Values: map[string]map[string]int{"x": {}}}
				api := withMapOfMap{Name: ptr("test"), Values: map[string]map[string]int{"x": {"y": 1}}}
				expectDrift(drift.Detect(crd, api), `values:
  x:
    y: <nil> != 1`)
			})

			It(" detect map of maps of struct drift, partial", func() {
				crd := withMapOfMap{Name: ptr("test"), Values: map[string]map[string]int{"x": {}}}
				api := withMapOfMap{Name: ptr("test"), Values: map[string]map[string]int{"x": {"y": 1}}}
				expectDrift(drift.Detect(crd, api), `values:
  x:
    y: <nil> != 1`)
			})
		})
	})

	Describe("when no drift is detected", func() {

		DescribeTable("equivalent values",
			func(crd, api any) {
				expectNoDrift(drift.Detect(crd, api))
			},
			Entry("empty-is-nil: no values",
				singleNilablePtr{},
				singleNilablePtr{},
			),
			Entry("empty-is-nil: empty string vs nil",
				singleNilablePtr{Value: ptr("")},
				singleNilablePtr{Value: nil},
			),
			Entry("empty-is-nil: nil vs empty string",
				singleNilablePtr{Value: nil},
				singleNilablePtr{Value: ptr("")},
			),
			Entry("empty-is-nil: nil vs nil",
				singleNilablePtr{Value: nil},
				singleNilablePtr{Value: nil},
			),
			Entry("empty-is-nil: empty string vs empty string",
				singleNilablePtr{Value: ptr("")},
				singleNilablePtr{Value: ptr("")},
			),
			Entry("double pointer property with no values",
				MultipleWithPtr{},
				MultipleWithPtr{},
			),
			Entry("double pointer property with same values",
				MultipleWithPtr{Name: ptr("foo"), Order: ptr(1)},
				MultipleWithPtr{Name: ptr("foo"), Order: ptr(1)},
			),
			Entry("nested structs, equal",
				Nested{},
				Nested{},
			),
			Entry("nested structs, equal",
				Nested{Name: ptr("foo"), Multiple: MultipleWithPtr{Name: ptr("foo"), Order: ptr(2)}, Description: ptr("bar")},
				Nested{Name: ptr("foo"), Multiple: MultipleWithPtr{Name: ptr("foo"), Order: ptr(2)}, Description: ptr("bar")},
			), Entry("nested structs, equivalent",
				Nested{Name: ptr("foo"), Multiple: MultipleWithPtr{Name: ptr("foo"), Order: ptr(0)}, Description: ptr("bar")},
				Nested{Name: ptr("foo"), Multiple: MultipleWithPtr{Name: ptr("foo"), Order: nil}, Description: ptr("bar")},
			),
			Entry("struct with nested pointer structs, empty",
				NestedWithPointer{},
				NestedWithPointer{},
			),
			Entry("struct with nested pointer structs, equivalents",
				NestedWithPointer{Title: ptr(""), Multiple: &MultipleWithPtr{}, Description: ptr("")},
				NestedWithPointer{Title: ptr(""), Multiple: nil, Description: ptr("")},
			),
			Entry("struct with deep nested pointer structs, empty",
				nestedDeepWithPtr{},
				nestedDeepWithPtr{},
			),
			Entry("struct with deep nested pointer structs, equals",
				nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{
					Title:       ptr("foo"),
					Multiple:    &MultipleWithPtr{Name: ptr("bar"), Order: ptr(1)},
					Description: ptr("Test!")}},
				nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{
					Title:       ptr("foo"),
					Multiple:    &MultipleWithPtr{Name: ptr("bar"), Order: ptr(1)},
					Description: ptr("Test!")}},
			),
			Entry("struct with deep nested pointer structs, equivalents at field level",
				nestedDeepWithPtr{Tag: "test", Nested: nil},
				nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{
					Title:       nil,
					Multiple:    &MultipleWithPtr{Name: nil, Order: nil},
					Description: nil}},
			),
			Entry("struct with deep nested pointer structs, equivalents at field & struct level",
				nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{Title: ptr(""), Multiple: nil, Description: nil}},
				nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{Title: nil, Multiple: nil, Description: ptr("")}}),
			Entry("struct with deep nested pointer structs, equivalents at field level",
				nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{Title: ptr(""),
					Multiple: &MultipleWithPtr{Name: nil, Order: ptr(0)}, Description: nil}},
				nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{Title: nil,
					Multiple: &MultipleWithPtr{Name: ptr(""), Order: nil}, Description: ptr("")}},
			),
			Entry("embedded structs",
				Embedded{},
				Embedded{},
			),
			Entry("embedded structs with pointer",
				Embedded{MultipleWithPtr: MultipleWithPtr{Name: ptr("foo"), Order: ptr(2)}, Description: ptr("bar")},
				Embedded{MultipleWithPtr: MultipleWithPtr{Name: ptr("foo"), Order: ptr(2)}, Description: ptr("bar")},
			),
			Entry("embedded structs, equivalent",
				Embedded{MultipleWithPtr: MultipleWithPtr{Name: ptr("foo"), Order: ptr(2)}, Description: ptr("")},
				Embedded{MultipleWithPtr: MultipleWithPtr{Name: ptr("foo"), Order: ptr(2)}, Description: nil},
			),
			Entry("embedded structs",
				embeddedWithPtr{},
				embeddedWithPtr{},
			),
			Entry("embedded structs",
				embeddedWithPtr{MultipleWithPtr: &MultipleWithPtr{Name: ptr("foo"), Order: ptr(2)}, Description: ptr("bar")},
				embeddedWithPtr{MultipleWithPtr: &MultipleWithPtr{Name: ptr("foo"), Order: ptr(2)}, Description: ptr("bar")},
			),
			Entry("embedded structs with pointer, equivalent",
				embeddedWithPtr{MultipleWithPtr: &MultipleWithPtr{Name: ptr("foo"), Order: ptr(2)}, Description: ptr("")},
				embeddedWithPtr{MultipleWithPtr: &MultipleWithPtr{Name: ptr("foo"), Order: ptr(2)}, Description: nil},
			),
			Entry("struct with string array",
				withStringArray{},
				withStringArray{},
			),
			Entry("struct with string array",
				withStringArray{Name: "test", Tags: []string{"foo", "bar"}},
				withStringArray{Name: "test", Tags: []string{"foo", "bar"}},
			),
			Entry("struct with struct array, equals",
				withStructArray{Name: ptr("test"), Values: []MultipleWithPtr{{Name: ptr("foo"), Order: ptr(1)}, {Name: ptr("bar")}}},
				withStructArray{Name: ptr("test"), Values: []MultipleWithPtr{{Name: ptr("foo"), Order: ptr(1)}, {Name: ptr("bar")}}},
			),
			Entry("struct with struct array, equivalents",
				withStructArray{Name: nil, Values: []MultipleWithPtr{{Name: ptr(""), Order: ptr(0)}, {Name: nil, Order: nil}}},
				withStructArray{Name: ptr(""), Values: []MultipleWithPtr{{Name: nil}, {Name: ptr(""), Order: ptr(0)}}},
			),
			Entry("struct with equivalent equivalents",
				withStructArray{Name: ptr("test"), Values: []MultipleWithPtr{}},
				withStructArray{Name: ptr("test"), Values: nil},
			),
			Entry("struct with complex structure, equivalents",
				nestedDeep{Title: "Test",
					First:  Nested{Name: ptr("First"), Multiple: MultipleWithPtr{Name: ptr("foo"), Order: ptr(1)}},
					Second: &Nested{Name: ptr(""), Multiple: MultipleWithPtr{Name: ptr("foo"), Order: ptr(0)}, Description: nil}},
				nestedDeep{Title: "Test",
					First:  Nested{Name: ptr("First"), Multiple: MultipleWithPtr{Name: ptr("foo"), Order: ptr(1)}},
					Second: &Nested{Name: nil, Multiple: MultipleWithPtr{Name: ptr("foo"), Order: nil}, Description: ptr("")}},
			),
			Entry("struct with complex structure, equivalents",
				nestedDeep{},
				nestedDeep{},
			),
			Entry("struct with complex nested and embedded, equivalents",
				embeddedAndNested{Title: "Test",
					Embedded:       Embedded{MultipleWithPtr: MultipleWithPtr{Name: ptr("test"), Order: nil}, Description: ptr("coming soon")},
					NestedEmbedded: Embedded{MultipleWithPtr: MultipleWithPtr{Name: ptr("foo")}, Description: ptr("bar")}},
				embeddedAndNested{Title: "Test",
					Embedded:       Embedded{MultipleWithPtr: MultipleWithPtr{Name: ptr("test"), Order: nil}, Description: ptr("coming soon")},
					NestedEmbedded: Embedded{MultipleWithPtr: MultipleWithPtr{Name: ptr("foo")}, Description: ptr("bar")}},
			),
			Entry("struct with complex nested and embedded, equivalents",
				embeddedAndNested{},
				embeddedAndNested{},
			),
			Entry("struct with map of a simple type, empty",
				withMapOfPrimitive{},
				withMapOfPrimitive{},
			),
			Entry("struct with map of a simple type, equal",
				withMapOfPrimitive{Name: ptr("test"), Values: map[string]*int{"foo": ptr(5), "bar": ptr(7), "baz": ptr(11)}},
				withMapOfPrimitive{Name: ptr("test"), Values: map[string]*int{"foo": ptr(5), "bar": ptr(7), "baz": ptr(11)}},
			),
			Entry("struct with map of a simple type, equivalent",
				withMapOfPrimitive{Name: ptr("test"), Values: map[string]*int{"x": nil, "y": nil, "z": nil}},
				withMapOfPrimitive{Name: ptr("test"), Values: map[string]*int{"x": ptr(0), "y": ptr(0), "z": ptr(0)}},
			),
			Entry("struct with map of a simple type, equivalent crd missing",
				withMapOfPrimitive{Name: ptr("test"), Values: map[string]*int{}},
				withMapOfPrimitive{Name: ptr("test"), Values: map[string]*int{"x": ptr(0), "y": ptr(0), "z": ptr(0)}},
			),
			Entry("struct with map of a simple type, equivalent api missing",
				withMapOfPrimitive{Name: ptr("test"), Values: map[string]*int{"x": ptr(0), "y": ptr(0), "z": ptr(0)}},
				withMapOfPrimitive{Name: ptr("test"), Values: map[string]*int{}},
			),
			Entry("struct with map of a simple type, no map",
				withMapOfPrimitive{Name: ptr("test"), Values: map[string]*int{"x": ptr(0), "y": ptr(0), "z": ptr(0)}},
				withMapOfPrimitive{Name: ptr("test"), Values: nil},
			),
			Entry("struct with map of a simple type, no map",
				withMapOfPrimitive{Name: ptr("test"), Values: nil},
				withMapOfPrimitive{Name: ptr("test"), Values: map[string]*int{"x": ptr(0), "y": ptr(0), "z": ptr(0)}},
			),
			Entry("struct with map of struct, equal",
				withMapOfStruct{},
				withMapOfStruct{},
			),
			Entry("struct with map of struct, equal",
				withMapOfStruct{Name: ptr("test"), Values: map[string]MultipleWithPtr{
					"x": {Name: ptr("foo"), Order: ptr(66)},
					"y": {Name: ptr("bar"), Order: ptr(22)}}},
				withMapOfStruct{Name: ptr("test"), Values: map[string]MultipleWithPtr{
					"x": {Name: ptr("foo"), Order: ptr(66)},
					"y": {Name: ptr("bar"), Order: ptr(22)}}},
			),
			Entry("struct with map of struct, equivalent",
				withMapOfStruct{Name: ptr("test"), Values: map[string]MultipleWithPtr{"x": {Name: ptr(""), Order: nil}}},
				withMapOfStruct{Name: ptr("test"), Values: map[string]MultipleWithPtr{"x": {Name: nil, Order: ptr(0)}}},
			),
			Entry("struct with map of struct, equivalent missing api",
				withMapOfStruct{Name: ptr("test"), Values: map[string]MultipleWithPtr{"x": {}}},
				withMapOfStruct{Name: ptr("test"), Values: map[string]MultipleWithPtr{}},
			),
			Entry("struct with map of struct, equivalent missing crd",
				withMapOfStruct{Name: ptr("test"), Values: map[string]MultipleWithPtr{}},
				withMapOfStruct{Name: ptr("test"), Values: map[string]MultipleWithPtr{"x": {}}},
			),
			Entry("struct with map of array, equal",
				withMapOfArray{},
				withMapOfArray{},
			),
			Entry("struct with map of array, equal",
				withMapOfArray{Name: ptr("test"), Values: map[string][]int{"x": {42}, "y": {66}}},
				withMapOfArray{Name: ptr("test"), Values: map[string][]int{"x": {42}, "y": {66}}},
			), Entry("struct with map of array, equivalent",
				withMapOfArray{Name: ptr("test"), Values: map[string][]int{"x": {}, "y": {66}}},
				withMapOfArray{Name: ptr("test"), Values: map[string][]int{"x": nil, "y": {66}}},
			),
			Entry("struct with map of array, equivalent missing api",
				withMapOfArray{Name: ptr("test"), Values: map[string][]int{"x": {}, "y": {66}}},
				withMapOfArray{Name: ptr("test"), Values: map[string][]int{"y": {66}}},
			),
			Entry("struct with map of array, equivalent missing crd",
				withMapOfArray{Name: ptr("test"), Values: map[string][]int{"y": {66}}},
				withMapOfArray{Name: ptr("test"), Values: map[string][]int{"x": {}, "y": {66}}},
			),
			Entry("struct with map of struct array, equal",
				withMapOfStructArray{},
				withMapOfStructArray{},
			),
			Entry("struct with map of struct array, equal",
				withMapOfStructArray{Name: ptr("test"), Values: map[string][]MultipleWithPtr{
					"x": {MultipleWithPtr{Name: ptr("foo"), Order: ptr(66)}},
					"y": {MultipleWithPtr{Name: ptr("bar"), Order: ptr(22)}}}},
				withMapOfStructArray{Name: ptr("test"), Values: map[string][]MultipleWithPtr{
					"x": {MultipleWithPtr{Name: ptr("foo"), Order: ptr(66)}},
					"y": {MultipleWithPtr{Name: ptr("bar"), Order: ptr(22)}}}},
			),
			Entry("struct with map of struct array, equivalent",
				withMapOfStructArray{Name: ptr("test"), Values: map[string][]MultipleWithPtr{
					"x": {MultipleWithPtr{Name: ptr(""), Order: ptr(0)}},
					"y": {MultipleWithPtr{Name: ptr(""), Order: ptr(0)}}}},
				withMapOfStructArray{Name: ptr("test"), Values: map[string][]MultipleWithPtr{
					"x": {MultipleWithPtr{Name: nil, Order: nil}},
					"y": {MultipleWithPtr{Name: nil, Order: nil}}}},
			),
			Entry("struct with map of map, equal",
				withMapOfMap{},
				withMapOfMap{},
			),
			Entry("struct with map of map, equal",
				withMapOfMap{Name: ptr("test"), Values: map[string]map[string]int{"x": {"y": 42}}},
				withMapOfMap{Name: ptr("test"), Values: map[string]map[string]int{"x": {"y": 42}}},
			),
			Entry("struct with map of map, equivalent",
				withMapOfMap{Name: ptr("test"), Values: map[string]map[string]int{"x": {}}},
				withMapOfMap{Name: ptr("test"), Values: map[string]map[string]int{"x": nil}},
			),
		)
	})

	Describe("when invalid inputs cause panic", func() {
		DescribeTable("non-struct root values",
			func(crd, api any) {
				Expect(func() {
					drift.Detect(crd, api)
				}).To(PanicWith(MatchRegexp(`detect drift only supports structs`)))
			},
			Entry("crd is a string", "not a struct", SingleWithPtr{}),
			Entry("api is a string", SingleWithPtr{}, "not a struct"),
			Entry("crd is a pointer to struct", &SingleWithPtr{}, SingleWithPtr{}),
			Entry("api is a pointer to struct", SingleWithPtr{}, &SingleWithPtr{}),
		)

		It("panics when a map key is not a string", func() {
			crd := withIntMapKey{Values: map[int]string{1: "a"}}
			api := withIntMapKey{Values: map[int]string{1: "b"}}
			Expect(func() {
				drift.Detect(crd, api)
			}).To(PanicWith(MatchRegexp(`map key must be of type string`)))
		})

		It("panics when a nested map key is not a string", func() {
			crd := withNestedIntMapKey{Values: map[string]map[int]string{"outer": {1: "a"}}}
			api := withNestedIntMapKey{Values: map[string]map[int]string{"outer": {1: "b"}}}
			Expect(func() {
				drift.Detect(crd, api)
			}).To(PanicWith(MatchRegexp(`map key must be of type string`)))
		})

		It("panics when registering an equivalence for pointer kind", func() {
			Expect(func() {
				drift.Register("pointer-kind", reflect.Pointer, drift.FromDeepEqual)
			}).To(PanicWith("cannot register a pointer to a struct, use a concrete type or an interface"))
		})

		It("panics when comparing values of different types through a registered equivalence", func() {
			drift.Register("type-check-panic", reflect.Interface, drift.FromDeepEqual)
			crd := withAnyField{Value: "foo"}
			api := withAnyField{Value: 42}
			Expect(func() {
				drift.Detect(crd, api)
			}).To(PanicWith(MatchRegexp(`drift detection only work comparing values of same type`)))
		})
	})

})

var _ = Describe("Merge", func() {

	DescribeTable("no drift",
		func(oldCRD, newCRD, api any) {
			or := drift.Detect(oldCRD, api)
			nr := drift.Detect(newCRD, api)
			result := drift.Merge(or, nr)
			expectNoDrift(result)
		},
		// Case 1 (AAA)
		Entry("case 1 simple",
			MultipleWithPtr{Name: ptr("foo"), Order: ptr(0)},
			MultipleWithPtr{Name: ptr("foo"), Order: ptr(0)},
			MultipleWithPtr{Name: ptr("foo")},
		),
		Entry("case 1 deep nested",
			nestedDeepWithPtr{Tag: "test", Nested: nil},
			nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{
				Title: ptr("")}},
			nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{
				Title:       nil,
				Multiple:    &MultipleWithPtr{Name: nil, Order: nil},
				Description: nil}},
		),
		Entry("case 1 maps",
			withMapOfArray{Name: ptr("test"), Values: map[string][]int{"x": {}, "y": {66}, "z": {12}}},
			withMapOfArray{Name: ptr("test"), Values: map[string][]int{"y": {66}, "z": {12}}},
			withMapOfArray{Name: ptr("test"), Values: map[string][]int{"x": {}, "y": {66}, "z": {12}}},
		),
		// case 2 CRD change matches remote change (ABB)
		Entry("case 2 simple",
			MultipleWithPtr{Name: ptr("Old Name"), Order: ptr(0)},
			MultipleWithPtr{Name: ptr("New Name"), Order: ptr(0)},
			MultipleWithPtr{Name: ptr("New Name")},
		),
		Entry("case 2 deep nested",
			nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{
				Title: ptr("Old Title")}},
			nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{
				Title: ptr("New Title")}},
			nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{
				Title:       ptr("New Title"),
				Multiple:    &MultipleWithPtr{Name: nil, Order: nil},
				Description: nil}},
		),
		Entry("case 2 maps",
			withMapOfArray{Name: ptr("test"), Values: map[string][]int{"x": {0}, "y": {66}, "z": {12}}},
			withMapOfArray{Name: ptr("test"), Values: map[string][]int{"x": {42}, "y": {66}, "z": {12}}},
			withMapOfArray{Name: ptr("test"), Values: map[string][]int{"x": {42}, "y": {66}, "z": {12}}},
		),
		// case 4 Only CRD changes (ABA)
		Entry("case 2 simple",
			MultipleWithPtr{Name: ptr("Old Name"), Order: ptr(0)},
			MultipleWithPtr{Name: ptr("New Name"), Order: ptr(0)},
			MultipleWithPtr{Name: ptr("Old Name")},
		),
		Entry("case 2 deep nested",
			nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{
				Title: ptr("Old Title")}},
			nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{
				Title: ptr("New Title")}},
			nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{
				Title:       ptr("Old Title"),
				Multiple:    &MultipleWithPtr{Name: nil, Order: nil},
				Description: nil}},
		),
		Entry("case 2 maps",
			withMapOfArray{Name: ptr("test"), Values: map[string][]int{"x": {0}, "y": {66}, "z": {12}}},
			withMapOfArray{Name: ptr("test"), Values: map[string][]int{"x": {42}, "y": {66}, "z": {12}}},
			withMapOfArray{Name: ptr("test"), Values: map[string][]int{"x": {0}, "y": {66}, "z": {12}}},
		),
	)
	DescribeTable("expect drift",
		func(oldCRD, newCRD, api any, output string) {
			or := drift.Detect(oldCRD, api)
			nr := drift.Detect(newCRD, api)
			result := drift.Merge(or, nr)
			expectDrift(result, output)
		},
		// case 2 on order, case 3 on name (AAB)
		Entry("case 3 simple",
			MultipleWithPtr{Name: ptr("CRD Name"), Order: ptr(0)},
			MultipleWithPtr{Name: ptr("CRD Name"), Order: ptr(1)},
			MultipleWithPtr{Name: ptr("Remote Name"), Order: ptr(0)},
			`name: "CRD Name" != "Remote Name"`,
		),
		// case 2 on nested name, case 3 on title (AAB)
		Entry("case 3 deep nested",
			nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{
				Title:    ptr("Old Title"),
				Multiple: &MultipleWithPtr{Name: ptr("CRD Name"), Order: ptr(0)},
			}},
			nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{
				Title:    ptr("New Title"),
				Multiple: &MultipleWithPtr{Name: ptr("CRD Name"), Order: ptr(0)},
			}},
			nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{
				Title:    ptr("Old Title"),
				Multiple: &MultipleWithPtr{Name: ptr("Remote Name"), Order: ptr(0)},
			}},
			`nested:
  multiple:
    name: "CRD Name" != "Remote Name"`,
		),
		// case 2 on x, case 3 on y (AAB)
		Entry("case 3 maps",
			withMapOfArray{Name: ptr("test"), Values: map[string][]int{"x": {0}, "y": {66}, "z": {12}}},
			withMapOfArray{Name: ptr("test"), Values: map[string][]int{"x": {42}, "y": {66}, "z": {12}}},
			withMapOfArray{Name: ptr("test"), Values: map[string][]int{"x": {42}, "y": {66, 67}, "z": {12}}},
			`values:
  y[1]: 0 != 67`,
		),
		// case 2 on order, case 5 on name (ABC)
		Entry("case 5 simple",
			MultipleWithPtr{Name: ptr("Old Name"), Order: ptr(0)},
			MultipleWithPtr{Name: ptr("New Name"), Order: ptr(1)},
			MultipleWithPtr{Name: ptr("Remote Name"), Order: ptr(0)},
			`name: "New Name" != "Remote Name"`,
		),
		// case 2 on nested name, case 5 on title (ABC)
		Entry("case 5 deep nested",
			nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{
				Title:    ptr("Old Title"),
				Multiple: &MultipleWithPtr{Name: ptr("Old Name"), Order: ptr(0)},
			}},
			nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{
				Title:    ptr("New Title"),
				Multiple: &MultipleWithPtr{Name: ptr("New Name"), Order: ptr(0)},
			}},
			nestedDeepWithPtr{Tag: "test", Nested: &NestedWithPointer{
				Title:    ptr("Old Title"),
				Multiple: &MultipleWithPtr{Name: ptr("Remote Name"), Order: ptr(0)},
			}},
			`nested:
  multiple:
    name: "New Name" != "Remote Name"`,
		),
		// case 2 on x, case 5 on y (ABC)
		Entry("case 5 maps",
			withMapOfArray{Name: ptr("test"), Values: map[string][]int{"a": {42}, "b": {11}, "c": {0}, "x": {0}, "y": {66}, "z": {12}}},
			withMapOfArray{Name: ptr("test"), Values: map[string][]int{"a": {0}, "b": {11}, "c": {42}, "x": {42}, "y": {67}, "z": {12}}},
			withMapOfArray{Name: ptr("test"), Values: map[string][]int{"a": {42}, "b": {11}, "c": {42}, "x": {42}, "y": {68}, "z": {12}}},
			`values:
  y[0]: 67 != 68`,
		),

		Entry("case 1AAA (tag), 2ABB (title), 3AAB (name) => drift, 4(order), 5ABC(description) => drift deep nested",
			nestedDeepWithPtr{
				Tag: "Same", Nested: &NestedWithPointer{
					Title: ptr("Old Title"),
					Multiple: &MultipleWithPtr{
						Name:  ptr("Old Name"),
						Order: ptr(2)},
					Description: ptr("Old Description"),
				}},
			nestedDeepWithPtr{
				Tag: "Same", Nested: &NestedWithPointer{
					Title: ptr("New Title"),
					Multiple: &MultipleWithPtr{
						Name:  ptr("New Name"),
						Order: ptr(30)},
					Description: ptr("New Description"),
				}},
			nestedDeepWithPtr{Tag: "Same", Nested: &NestedWithPointer{
				Title: ptr("New Title"),
				Multiple: &MultipleWithPtr{
					Name:  ptr("Remote Name"),
					Order: ptr(2)},
				Description: ptr("Remote Description"),
			}},
			`nested:
  multiple:
    name: "New Name" != "Remote Name"
  description: "New Description" != "Remote Description"`,
		),
	)

})
