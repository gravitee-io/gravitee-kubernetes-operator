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

package apim

import (
	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	apimsubscription "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/subscription"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// legacySubscriptionState is what APIM answers to a GET by UUID (hridContainsUUID=true): it sets no hrid,
// and the application and plan hrids are their UUIDs, set by the APIM 4.12 hrid upgraders.
const legacySubscriptionState = `{
	"id": "sub-uuid",
	"hrid": null,
	"apiHrid": "api-uuid",
	"applicationHrid": "app-uuid",
	"planHrid": "plan-uuid",
	"startingAt": "2026-01-01T00:00:00Z",
	"endingAt": "2070-01-01T00:00:00Z"
}`

var _ = Describe("Subscription created before the Automation API", func() {
	It("admits an endingAt change when nothing changed in APIM", func() {
		oldDTO := legacySubscriptionDTO("2070-01-01T00:00:00Z")
		newDTO := legacySubscriptionDTO("2071-01-01T00:00:00Z")
		remote := apimsubscription.FromState(parseJSON[sdk.SubscriptionState](legacySubscriptionState), newDTO)

		merged := drift.Merge(
			drift.DetectWithNamespace(oldDTO, remote, "default"),
			drift.DetectWithNamespace(newDTO, remote, "default"),
		)

		expectNoDrift(merged)
	})
})

// legacySubscriptionDTO maps a subscription, its API and its application, all three created before the
// Automation API (a status id and no AutomationAPIManaged condition), as the reconcile and admission do.
func legacySubscriptionDTO(endingAt string) apimsubscription.DTO {
	GinkgoHelper()

	sub := &v1alpha1.Subscription{ObjectMeta: metav1.ObjectMeta{Name: "sub", Namespace: "default"}}
	sub.Spec.Plan = "jwt"
	sub.Spec.EndingAt = &endingAt
	sub.Status.ID = "sub-uuid"

	api := &v1alpha1.ApiV4Definition{ObjectMeta: metav1.ObjectMeta{Name: "api", Namespace: "default"}}
	api.Spec.V4BaseApi = &v4.V4BaseApi{ApiBase: &base.ApiBase{}}
	api.Status.ID = "api-uuid"
	plan := &v4.Plan{Plan: &base.Plan{ID: "plan-uuid"}}
	api.Spec.Plans = &map[string]*v4.Plan{"jwt": plan}

	app := &v1alpha1.Application{ObjectMeta: metav1.ObjectMeta{Name: "app", Namespace: "default"}}
	app.Status.ID = "app-uuid"

	apimsubscription.RecordRefs(sub, api, app)
	dto, err := apimsubscription.ToDTO(sub)
	Expect(err).ToNot(HaveOccurred())
	return dto
}
