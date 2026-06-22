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
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
)

const (
	registryTestAlwaysEqualName = "registry-test-always-equal"
	registryTestTypeCheckName   = "registry-test-type-check"
	registryTestCustomSliceName = "registry-test-custom-slice"
)

type registryCustomString struct {
	Name string `drift:"registry-test-always-equal" json:"name"`
}

type registryInitBoolField struct {
	Flag *bool `drift:"empty-is-nil" json:"flag"`
}

type registryIntFallback struct {
	Count int `drift:"registry-test-unknown-func" json:"count"`
}

type registryFloatFallback struct {
	Score float64 `json:"score"`
}

type registryArrayFallback struct {
	Tags [2]string `json:"tags"`
}

type registryUnknownSliceTag struct {
	Items []string `drift:"registry-test-unknown-slice" json:"items"`
}

type registryCustomSliceTag struct {
	Items []string `drift:"registry-test-custom-slice" json:"items"`
}

type registryTypeCheckAny struct {
	Value any `drift:"registry-test-type-check" json:"value"`
}

var _ = Describe("Registry", func() {
	Describe("Register", func() {
		It("panics when registering a pointer kind", func() {
			Expect(func() {
				drift.Register("registry-test-pointer-kind", reflect.Pointer, drift.FromDeepEqual)
			}).To(PanicWith("cannot register a pointer to a struct, use a concrete type or an interface"))
		})

		It("stores a custom equivalence func retrievable via Detect", func() {
			drift.Register(registryTestAlwaysEqualName, reflect.String, func(_, _ any) drift.Equivalence {
				return drift.Equivalence{Equivalent: drift.Equivalent}
			})

			expectNoDrift(drift.Detect(
				registryCustomString{Name: "foo"},
				registryCustomString{Name: "bar"},
			))
		})

		It("wraps registered funcs with assertTypes", func() {
			drift.Register(registryTestTypeCheckName, reflect.Interface, drift.FromDeepEqual)

			Expect(func() {
				drift.Detect(
					registryTypeCheckAny{Value: 1},
					registryTypeCheckAny{Value: "not-an-int"},
				)
			}).To(PanicWith(MatchRegexp(`drift detection only work comparing values of same type`)))
		})
	})

	Describe("Get", func() {
		It("returns Init-registered empty-is-nil for bool kind", func() {
			expectNoDrift(drift.Detect(
				registryInitBoolField{Flag: nil},
				registryInitBoolField{Flag: ptr(false)},
			))
		})

		It("falls back to defaultEquivalence when drift tag is absent", func() {
			expectDrift(drift.Detect(
				registryFloatFallback{Score: 1.5},
				registryFloatFallback{Score: 2.5},
			), "score: 1.5 != 2.5")
		})

		It("falls back to defaultSliceArrayEquivalence for unregistered array kind", func() {
			result := drift.Detect(
				registryArrayFallback{Tags: [2]string{"a", "b"}},
				registryArrayFallback{Tags: [2]string{"x", "y"}},
			)

			Expect(result.Children).To(HaveLen(1))
			Expect(result.Children[0].Equivalent).To(Equal(drift.CannotCompare))
			Expect(result.DriftDetected()).To(BeFalse())
		})

		It("panics if the drift function is unknown", func() {
			Expect(func() {
				drift.Detect(
					registryUnknownSliceTag{Items: []string{"a"}},
					registryUnknownSliceTag{Items: []string{"b"}},
				)
			}).To(PanicWith("drift function 'registry-test-unknown-slice' not found for kind 'slice'"))
		})

		It("uses a registered slice equivalence func when name is known", func() {
			drift.Register(registryTestCustomSliceName, reflect.Slice, func(_, _ any) drift.Equivalence {
				return drift.Equivalence{Equivalent: drift.Equivalent, Skip: true}
			})

			expectNoDrift(drift.Detect(
				registryCustomSliceTag{Items: []string{"a"}},
				registryCustomSliceTag{Items: []string{"b"}},
			))
		})
	})

	Describe("defaultEquivalence", func() {
		It("reports equivalent values via deep equal", func() {
			expectNoDrift(drift.Detect(
				registryFloatFallback{Score: 3.14},
				registryFloatFallback{Score: 3.14},
			))
		})

		It("reports inequivalent values via deep equal", func() {
			expectDrift(drift.Detect(
				registryFloatFallback{Score: 1},
				registryFloatFallback{Score: 2},
			), "score: 1 != 2")
		})
	})

	Describe("defaultSliceArrayEquivalence", func() {
		It("returns CannotCompare without reporting drift", func() {
			result := drift.Detect(
				registryArrayFallback{Tags: [2]string{"same", "values"}},
				registryArrayFallback{Tags: [2]string{"same", "values"}},
			)

			Expect(result.Children).To(HaveLen(1))
			Expect(result.Children[0].Equivalent).To(Equal(drift.CannotCompare))
			Expect(result.DriftDetected()).To(BeFalse())
		})
	})
})
