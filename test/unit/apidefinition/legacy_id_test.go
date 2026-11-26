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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// The expected ids are the ones GKO 4.11 pickID and pickCrossID return for the same input.
var _ = Describe("LegacyID", func() {
	mCtx := &management.Context{OrgID: "DEFAULT", EnvID: "DEFAULT"}

	newAPI := func() *v1alpha1.ApiV4Definition {
		return &v1alpha1.ApiV4Definition{
			ObjectMeta: metav1.ObjectMeta{Name: "my-api", Namespace: "default"},
			Spec: v1alpha1.ApiV4DefinitionSpec{
				Api: v4.Api{V4BaseApi: &v4.V4BaseApi{ApiBase: &base.ApiBase{}}},
			},
		}
	}

	It("derives the id GKO 4.11 gave an API from its namespaced name", func() {
		api := newAPI()

		Expect(api.LegacyCrossID()).To(Equal("3d7306da-8a12-9e6b-9749-deaab5bcf711"))
		Expect(api.LegacyID(mCtx)).To(Equal("964236d2-7b02-9a17-a40a-ab686e76433a"))
	})

	It("derives the id GKO 4.11 gave an API in another organization and environment", func() {
		api := newAPI()

		Expect(api.LegacyID(&management.Context{OrgID: "ORG-A", EnvID: "ENV-B"})).
			To(Equal("b299b6f8-a0ae-eca1-c5e1-b3ec7d0e0b50"))
	})

	It("derives the plan ids GKO 4.11 gave an API once its id is back in the status", func() {
		api := newAPI()
		api.Spec.Plans = &map[string]*v4.Plan{"API_KEY": {Plan: &base.Plan{}}}
		api.Status.ID = api.LegacyID(mCtx)

		api.PopulateIDs(mCtx, false)

		Expect(api.Spec.ID).To(Equal("964236d2-7b02-9a17-a40a-ab686e76433a"))
		Expect((*api.Spec.Plans)["API_KEY"].ID).To(Equal("d0aedf31-859f-47dd-ad18-df485d2709f4"))
	})

	It("derives the id GKO 4.11 gave an API from its spec cross id", func() {
		api := newAPI()
		api.Spec.CrossID = "exported-cross-id"

		Expect(api.LegacyCrossID()).To(Equal("exported-cross-id"))
		Expect(api.LegacyID(mCtx)).To(Equal("faf11dd3-fe78-e7e7-5bf3-5ddba9c43f79"))
	})
})
