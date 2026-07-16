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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hrid"
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
})
