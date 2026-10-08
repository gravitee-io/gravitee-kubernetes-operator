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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type setOnly struct {
	Tags []string `json:"tags,omitempty" drift:"set"`
}

type setOfPointers struct {
	Tags []*string `json:"tags,omitempty" drift:"set"`
}

type pointerToSet struct {
	Tags *[]string `json:"tags,omitempty" drift:"set"`
}

type setOfInts struct {
	Codes []int `json:"codes,omitempty" drift:"set"`
}

type label string

type setOfNamedStrings struct {
	Labels []label `json:"labels,omitempty" drift:"set"`
}

// mirrors sdk.Domain.CorsSettings: an untagged pointer to a struct holding sets
type setInPointerStruct struct {
	Cors *setOnly `json:"cors,omitempty"`
}

type setInStruct struct {
	Cors setOnly `json:"cors"`
}

type setInListItems struct {
	Items []setOnly `json:"items"`
}

type setInMapValues struct {
	Values map[string]setOnly `json:"values"`
}

type setOfStructs struct {
	Values []MultipleWithPtr `json:"values" drift:"set"`
}

type setOnString struct {
	Value string `json:"value" drift:"set"`
}

func detect(crd, remote any) drift.Result {
	return drift.DetectWithNamespace(crd, remote, "")
}

var _ = Describe("set equivalence", func() {
	DescribeTable("no drift",
		func(crd, remote any) {
			expectNoDrift(detect(crd, remote))
		},
		Entry("same order", setOnly{[]string{"a", "b"}}, setOnly{[]string{"a", "b"}}),
		Entry("other order", setOnly{[]string{"a", "b", "c"}}, setOnly{[]string{"c", "a", "b"}}),
		Entry("duplicates in the CRD", setOnly{[]string{"a", "a", "b"}}, setOnly{[]string{"b", "a"}}),
		Entry("duplicates on the remote", setOnly{[]string{"a", "b"}}, setOnly{[]string{"b", "a", "b"}}),
		Entry("nil and nil", setOnly{}, setOnly{}),
		Entry("nil and empty", setOnly{}, setOnly{[]string{}}),
		Entry("empty and nil", setOnly{[]string{}}, setOnly{}),
		Entry("ints in another order", setOfInts{[]int{1, 2}}, setOfInts{[]int{2, 1}}),
		Entry("named strings in another order",
			setOfNamedStrings{[]label{"x", "y"}}, setOfNamedStrings{[]label{"y", "x"}}),
		Entry("pointer items with equal values",
			setOfPointers{[]*string{new("a"), new("b")}}, setOfPointers{[]*string{new("b"), new("a")}}),
		Entry("pointer to a set, other order",
			pointerToSet{&[]string{"a", "b"}}, pointerToSet{&[]string{"b", "a"}}),
		Entry("pointer to a set, nil and empty", pointerToSet{}, pointerToSet{&[]string{}}),
		Entry("set in a pointer struct, other order",
			setInPointerStruct{&setOnly{[]string{"GET", "POST"}}}, setInPointerStruct{&setOnly{[]string{"POST", "GET"}}}),
		Entry("set in a pointer struct, both structs nil", setInPointerStruct{}, setInPointerStruct{}),
		Entry("set in a struct value, other order",
			setInStruct{setOnly{[]string{"a", "b"}}}, setInStruct{setOnly{[]string{"b", "a"}}}),
		Entry("set in list items, other order",
			setInListItems{[]setOnly{{[]string{"a", "b"}}}}, setInListItems{[]setOnly{{[]string{"b", "a"}}}}),
		Entry("set in map values, other order",
			setInMapValues{map[string]setOnly{"k": {[]string{"a", "b"}}}},
			setInMapValues{map[string]setOnly{"k": {[]string{"b", "a"}}}}),
	)

	DescribeTable("drift",
		func(crd, remote any, expected string) {
			expectDrift(detect(crd, remote), expected)
		},
		Entry("item added on the remote",
			setOnly{[]string{"a", "b"}}, setOnly{[]string{"b", "a", "c"}}, `tags: (2 unchanged)
  <unchanged>
              !=
                 + "c"`),
		Entry("item missing on the remote",
			setOnly{[]string{"a", "b", "c"}}, setOnly{[]string{"a", "b"}}, `tags: (2 unchanged)
  + "c"
        !=
           <unchanged>`),
		Entry("one item replaced",
			setOnly{[]string{"a", "b"}}, setOnly{[]string{"a", "c"}}, `tags: (1 unchanged)
  + "b"
        !=
           + "c"`),
		Entry("nil CRD, remote has items",
			setOnly{}, setOnly{[]string{"a"}}, `tags:
  <unchanged>
              !=
                 + "a"`),
		Entry("CRD has items, nil remote",
			setOnly{[]string{"a"}}, setOnly{}, `tags:
  + "a"
        !=
           <unchanged>`),
		Entry("case is significant",
			setOnly{[]string{"GET"}}, setOnly{[]string{"get"}}, `tags:
  + "GET"
          !=
             + "get"`),
		Entry("ints",
			setOfInts{[]int{1, 2}}, setOfInts{[]int{2, 3}}, `codes: (1 unchanged)
  + 1
      !=
         + 3`),
		Entry("pointer items print their value",
			setOfPointers{[]*string{new("a")}}, setOfPointers{[]*string{new("b")}}, `tags:
  + "a"
        !=
           + "b"`),
		Entry("pointer to a set",
			pointerToSet{&[]string{"a"}}, pointerToSet{&[]string{"b"}}, `tags:
  + "a"
        !=
           + "b"`),
		Entry("set nested in a pointer struct is indented",
			setInPointerStruct{&setOnly{[]string{"GET", "POST"}}},
			setInPointerStruct{&setOnly{[]string{"GET", "PUT"}}}, `cors:
  tags: (1 unchanged)
    + "POST"
             !=
                + "PUT"`),
		Entry("set nested in list items is indented",
			setInListItems{[]setOnly{{[]string{"a"}}}},
			setInListItems{[]setOnly{{[]string{"b"}}}}, `items[0]:
  tags:
    + "a"
          !=
             + "b"`),
	)

	Describe("a pointer struct present on one side only", func() {
		It("does not panic and drifts when only the remote has the struct", func() {
			var result drift.Result
			Expect(func() {
				result = detect(setInPointerStruct{}, setInPointerStruct{&setOnly{[]string{"a"}}})
			}).NotTo(Panic())
			Expect(result.DriftDetected()).To(BeTrue())
		})

		It("does not panic and drifts when only the CRD has the struct", func() {
			var result drift.Result
			Expect(func() {
				result = detect(setInPointerStruct{&setOnly{[]string{"a"}}}, setInPointerStruct{})
			}).NotTo(Panic())
			Expect(result.DriftDetected()).To(BeTrue())
		})

		It("does not drift when the only remote struct holds an empty set", func() {
			expectNoDrift(detect(setInPointerStruct{}, setInPointerStruct{&setOnly{[]string{}}}))
		})
	})

	It("lists the differing items in a stable order", func() {
		crd := setOnly{[]string{"a", "b", "c", "d", "e"}}
		remote := setOnly{[]string{"z", "y", "x", "w", "v"}}
		first := detect(crd, remote)
		for range 50 {
			next := detect(crd, remote)
			Expect(next.String()).To(Equal(first.String()))
		}
	})

	It("lists the differing items in the order each side declares them", func() {
		expectDrift(detect(setOnly{[]string{"c", "a", "b"}}, setOnly{[]string{"z", "x"}}), `tags:
  + "c"
  + "a"
  + "b"
        !=
           + "z"
           + "x"`)
	})

	It("compares non-scalar items by position", func() {
		crd := setOfStructs{[]MultipleWithPtr{{Name: new("a")}, {Name: new("b")}}}
		remote := setOfStructs{[]MultipleWithPtr{{Name: new("b")}, {Name: new("a")}}}
		result := detect(crd, remote)
		Expect(result.DriftDetected()).To(BeTrue())
	})

	It("panics when set tags a string", func() {
		Expect(func() {
			detect(setOnString{"a"}, setOnString{"b"})
		}).To(PanicWith(MatchRegexp(`drift function 'set' not found for kind 'string'`)))
	})

	DescribeTable("Merge",
		func(oldCRD, newCRD, remote any, expected string) {
			result := drift.Merge(detect(oldCRD, remote), detect(newCRD, remote))
			if expected == "" {
				expectNoDrift(result)
				return
			}
			expectDrift(result, expected)
		},
		Entry("unchanged CRD, remote reordered",
			setOnly{[]string{"a", "b"}}, setOnly{[]string{"a", "b"}}, setOnly{[]string{"b", "a"}}, ""),
		Entry("CRD reordered to match the remote",
			setOnly{[]string{"a", "b"}}, setOnly{[]string{"b", "a"}}, setOnly{[]string{"b", "a"}}, ""),
		Entry("CRD adds an item, remote unchanged",
			setOnly{[]string{"a", "b"}}, setOnly{[]string{"b", "a", "c"}}, setOnly{[]string{"b", "a"}}, ""),
		Entry("CRD removes an item, remote unchanged",
			setOnly{[]string{"a", "b"}}, setOnly{[]string{"a"}}, setOnly{[]string{"b", "a"}}, ""),
		Entry("unchanged CRD, remote adds an item",
			setOnly{[]string{"a", "b"}}, setOnly{[]string{"a", "b"}}, setOnly{[]string{"b", "a", "c"}},
			`tags: (2 unchanged)
  <unchanged>
              !=
                 + "c"`),
		Entry("CRD realigns with the remote, with duplicates",
			setOnly{[]string{"a", "a"}}, setOnly{[]string{"c", "c", "b"}}, setOnly{[]string{"b", "c"}}, ""),
		Entry("CRD and remote both changed, remote adds an item",
			setOnly{[]string{"a"}}, setOnly{[]string{"b", "c", "b"}}, setOnly{[]string{"c", "b", "x"}},
			`tags: (2 unchanged)
  <unchanged>
              !=
                 + "x"`),
		Entry("CRD and remote both changed, CRD keeps an extra item",
			setOnly{[]string{"a", "x"}}, setOnly{[]string{"b", "c", "b", "x"}}, setOnly{[]string{"c", "b"}},
			`tags: (2 unchanged)
  + "x"
        !=
           <unchanged>`),
		Entry("nested: unchanged CRD, remote reordered",
			setInPointerStruct{&setOnly{[]string{"GET", "POST"}}},
			setInPointerStruct{&setOnly{[]string{"GET", "POST"}}},
			setInPointerStruct{&setOnly{[]string{"POST", "GET"}}}, ""),
	)
})
