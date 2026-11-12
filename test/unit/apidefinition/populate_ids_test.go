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

package apidefinition_test

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/management"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hrid"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("PopulateIDs", func() {
	legacyAPI := func(plans *map[string]*v4.Plan, pages *map[string]*v4.Page) *v1alpha1.ApiV4Definition {
		return &v1alpha1.ApiV4Definition{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-api",
				Namespace: "default",
			},
			Spec: v1alpha1.ApiV4DefinitionSpec{
				Api: v4.Api{
					V4BaseApi: &v4.V4BaseApi{
						ApiBase: &base.ApiBase{Name: "Test API", Version: "1.0"},
					},
					Plans: plans,
					Pages: pages,
				},
			},
			Status: v1alpha1.ApiV4DefinitionStatus{
				Status: base.Status{
					ApiStatus: base.ApiStatus{
						ID: "existing-uuid",
					},
					Plans: map[string]string{
						"api-key":    "plan-uuid-1",
						"My Plan.v1": "plan-uuid-2",
					},
				},
			},
		}
	}

	DescribeTable("sets HRID on plans in legacy path",
		func(key string, expectedHRID string) {
			plans := map[string]*v4.Plan{
				key: {Plan: &base.Plan{}, Name: "Test Plan"},
			}
			api := legacyAPI(&plans, nil)
			api.Status.Plans = map[string]string{key: "plan-uuid"}

			api.PopulateIDs(nil, false)

			Expect(api.Spec.Plans).NotTo(BeNil())
			plan, ok := (*api.Spec.Plans)[key]
			Expect(ok).To(BeTrue(), "plan key %q should exist in result", key)
			Expect(plan.HRID).To(Equal(expectedHRID))
			Expect(plan.HRID).To(Equal(hrid.NameToValidHRID(key)))
		},
		Entry("simple key", "api-key", "api-key"),
		Entry("key with spaces and dots", "My Plan.v1", "My-Plan-v1"),
	)

	DescribeTable("sets HRID on pages in legacy path",
		func(key string, expectedHRID string) {
			pages := map[string]*v4.Page{
				key: {Page: &base.Page{Name: "Test Page", Type: "MARKDOWN"}},
			}
			api := legacyAPI(nil, &pages)

			api.PopulateIDs(nil, false)

			Expect(api.Spec.Pages).NotTo(BeNil())
			page, ok := (*api.Spec.Pages)[key]
			Expect(ok).To(BeTrue(), "page key %q should exist in result", key)
			Expect(page.HRID).To(Equal(expectedHRID))
			Expect(page.HRID).To(Equal(hrid.NameToValidHRID(key)))
		},
		Entry("simple key", "getting-started", "getting-started"),
		Entry("key with spaces and dots", "My Page.v1", "My-Page-v1"),
	)

	// Without spec.id, a plan id comes from status.plans, never from the spec.
	DescribeTable("sends the status plan id in HRID mode with a management context",
		func(statusPlans map[string]string, automationAPIManaged bool, expectedID string) {
			plans := map[string]*v4.Plan{
				"k": {Plan: &base.Plan{ID: "X"}, Name: "Test Plan"},
			}
			api := legacyAPI(&plans, nil)
			api.Status = v1alpha1.ApiV4DefinitionStatus{}
			api.Status.Plans = statusPlans
			if automationAPIManaged {
				api.Status.ID = "existing-uuid"
				k8s.AddAutomationAPIManagedCondition(api)
			}

			api.PopulateIDs(&management.Context{}, k8s.IsAutomationAPIManaged(api))

			Expect(api.Spec.ID).To(BeEmpty())
			Expect(api.Spec.Plans).NotTo(BeNil())
			plan, ok := (*api.Spec.Plans)["k"]
			Expect(ok).To(BeTrue())
			Expect(plan.HRID).To(Equal("k"))
			Expect(plan.ID).To(Equal(expectedID))
			Expect(plans["k"].ID).To(Equal("X"), "the CR spec itself is not modified")
		},
		Entry("status id differs from the spec id", map[string]string{"k": "S"}, false, "S"),
		Entry("status id differs from the spec id, AutomationAPIManaged",
			map[string]string{"k": "S"}, true, "S"),
		Entry("no status entry for the plan", map[string]string{"other": "S"}, false, ""),
		Entry("no status plans, AutomationAPIManaged", nil, true, ""),
	)

	It("keeps a plan without base fields in HRID mode with a management context", func() {
		plans := map[string]*v4.Plan{
			"k": {Name: "Test Plan"},
		}
		api := legacyAPI(&plans, nil)
		api.Status = v1alpha1.ApiV4DefinitionStatus{}

		api.PopulateIDs(&management.Context{}, k8s.IsAutomationAPIManaged(api))

		Expect((*api.Spec.Plans)["k"].HRID).To(Equal("k"))
	})

	It("keeps the spec plan and page ids in HRID mode when spec.id is set", func() {
		plans := map[string]*v4.Plan{
			"k": {Plan: &base.Plan{ID: "X"}, Name: "Test Plan"},
		}
		pages := map[string]*v4.Page{
			"doc": {Page: &base.Page{ID: "P", Name: "doc"}},
		}
		api := legacyAPI(&plans, &pages)
		api.Status = v1alpha1.ApiV4DefinitionStatus{}
		api.Status.Plans = map[string]string{"k": "S"}
		api.Spec.ID = "exported-uuid"
		k8s.AddAutomationAPIManagedCondition(api)

		api.PopulateIDs(&management.Context{}, k8s.IsAutomationAPIManaged(api))

		Expect(api.Spec.ID).To(Equal("exported-uuid"))
		Expect((*api.Spec.Plans)["k"].ID).To(Equal("X"))
		Expect((*api.Spec.Pages)["doc"].ID).To(Equal("P"))
	})
})
