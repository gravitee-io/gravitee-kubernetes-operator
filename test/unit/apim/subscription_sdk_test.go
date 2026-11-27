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

package apim_test

import (
	"context"
	"errors"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimsubscription "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/subscription"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

const subscriptionState = `{
	"id": "828be2ec-78e7-3306-8f70-a4a8fc94b384",
	"hrid": "default-sub",
	"applicationHrid": "default-app",
	"planHrid": "jwt",
	"startingAt": "2026-11-23T15:50:02.969Z",
	"endingAt": "2040-12-25T09:12:28Z"
}`

var _ = Describe("Subscription SDK calls", func() {
	var (
		ctx    context.Context
		fake   *fakeAutomationAPI
		client *apim.SDKClient
	)

	BeforeEach(func() {
		ctx = context.Background()
		fake = newFakeAutomationAPI()
		fake.response = subscriptionState
		DeferCleanup(fake.Close)

		var err error
		client, err = apim.NewSDKClient(ctx, managementContext(fake.URL, bearer("token")))
		Expect(err).ToNot(HaveOccurred())
	})

	subscriptionDTO := func(legacy apimsubscription.Legacy) apimsubscription.DTO {
		return apimsubscription.DTO{
			Subscription: sdk.Subscription{Hrid: "sub-id", ApplicationHrid: "app-id", PlanHrid: "jwt"},
			ApiHrid:      "api-id",
			Legacy:       legacy,
		}
	}

	expectFlags := func(want map[string]string) {
		GinkgoHelper()
		query := fake.lastRequest().Query
		Expect(query).To(HaveLen(len(want)))
		for flag, value := range want {
			Expect(query.Get(flag)).To(Equal(value), flag)
		}
	}

	DescribeTable("Upsert puts the subscription under its API with the three identity flags",
		func(legacy apimsubscription.Legacy, want map[string]string) {
			_, err := apimsubscription.Upsert(ctx, client, subscriptionDTO(legacy))
			Expect(err).ToNot(HaveOccurred())

			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodPut))
			Expect(req.Path).To(Equal("/automation/organizations/org1/environments/env1/apis/api-id/subscriptions"))
			Expect(req.Body).To(HaveKeyWithValue("hrid", "sub-id"))
			Expect(req.Body).ToNot(HaveKey("apiHrid"))
			expectFlags(want)
		},
		Entry("hrid-addressed",
			apimsubscription.Legacy{},
			map[string]string{"hridContainsUUID": "false", "hridContainsApiUUID": "false", "hridContainsAppUUID": "false"}),
		Entry("created before 4.12, with its API and application",
			apimsubscription.Legacy{HridContainsUUID: true, HridContainsApiUUID: true, HridContainsAppUUID: true},
			map[string]string{"hridContainsUUID": "true", "hridContainsApiUUID": "true", "hridContainsAppUUID": "true"}),
		Entry("new subscription to an API created before 4.12",
			apimsubscription.Legacy{HridContainsApiUUID: true},
			map[string]string{"hridContainsUUID": "false", "hridContainsApiUUID": "true", "hridContainsAppUUID": "false"}),
	)

	DescribeTable("GetState reads the subscription with the subscription and API flags",
		func(legacy apimsubscription.Legacy, want map[string]string) {
			_, err := apimsubscription.GetState(ctx, client, subscriptionDTO(legacy))
			Expect(err).ToNot(HaveOccurred())

			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodGet))
			Expect(req.Path).To(Equal("/automation/organizations/org1/environments/env1/apis/api-id/subscriptions/sub-id"))
			expectFlags(want)
		},
		Entry("hrid-addressed",
			apimsubscription.Legacy{},
			map[string]string{"hridContainsUUID": "false", "hridContainsApiUUID": "false"}),
		Entry("created before 4.12",
			apimsubscription.Legacy{HridContainsUUID: true, HridContainsApiUUID: true, HridContainsAppUUID: true},
			map[string]string{"hridContainsUUID": "true", "hridContainsApiUUID": "true"}),
	)

	DescribeTable("Delete removes the subscription with the subscription and API flags",
		func(legacy apimsubscription.Legacy, want map[string]string) {
			fake.status = http.StatusNoContent
			Expect(apimsubscription.Delete(ctx, client, subscriptionDTO(legacy))).To(Succeed())

			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodDelete))
			Expect(req.Path).To(Equal("/automation/organizations/org1/environments/env1/apis/api-id/subscriptions/sub-id"))
			expectFlags(want)
		},
		Entry("hrid-addressed",
			apimsubscription.Legacy{},
			map[string]string{"hridContainsUUID": "false", "hridContainsApiUUID": "false"}),
		Entry("subscription created before 4.12 under a new API",
			apimsubscription.Legacy{HridContainsUUID: true},
			map[string]string{"hridContainsUUID": "true", "hridContainsApiUUID": "false"}),
	)

	It("counts a subscription already gone from APIM as deleted", func() {
		fake.status = http.StatusNotFound
		Expect(apimsubscription.Delete(ctx, client, subscriptionDTO(apimsubscription.Legacy{}))).To(Succeed())
	})
})

var _ = Describe("Subscription UpdateStatus", func() {
	state := func(startingAt, endingAt *string) sdk.SubscriptionState {
		return sdk.SubscriptionState{
			Hrid:       "default-sub",
			Id:         new("828be2ec-78e7-3306-8f70-a4a8fc94b384"),
			StartingAt: startingAt,
			EndingAt:   endingAt,
		}
	}

	It("copies the id and dates APIM answers, and marks the subscription Automation API managed", func() {
		sub := newSub("sub", "", false)
		resp := apimsubscription.Response{State: state(new("2026-11-23T15:50:02.969Z"), new("2040-12-25T09:12:28Z"))}

		Expect(apimsubscription.UpdateStatus(sub, resp)).To(Succeed())

		Expect(sub.Status.ID).To(Equal("828be2ec-78e7-3306-8f70-a4a8fc94b384"))
		Expect(sub.Status.StartedAt).To(Equal("2026-11-23T15:50:02Z"))
		Expect(sub.Status.EndingAt).To(Equal("2040-12-25T09:12:28Z"))
		Expect(k8s.IsAutomationAPIManaged(sub)).To(BeTrue())
	})

	It("clears the end date when APIM answers none", func() {
		sub := newSub("sub", "828be2ec-78e7-3306-8f70-a4a8fc94b384", true)
		sub.Status.EndingAt = "2040-12-25T09:12:28Z"

		Expect(apimsubscription.UpdateStatus(sub, apimsubscription.Response{State: state(new("2026-11-23T15:50:02Z"), nil)})).
			To(Succeed())

		Expect(sub.Status.EndingAt).To(BeEmpty())
	})

	It("keeps the id it has when APIM answers none", func() {
		sub := newSub("sub", "828be2ec-78e7-3306-8f70-a4a8fc94b384", false)
		resp := apimsubscription.Response{
			State:  sdk.SubscriptionState{Hrid: "default-sub", StartingAt: new("2026-11-23T15:50:02Z")},
			Legacy: apimsubscription.Legacy{HridContainsUUID: true},
		}

		Expect(apimsubscription.UpdateStatus(sub, resp)).To(Succeed())

		Expect(sub.Status.ID).To(Equal("828be2ec-78e7-3306-8f70-a4a8fc94b384"))
	})

	It("leaves a subscription upserted by its UUID addressed by its UUID", func() {
		sub := newSub("sub", "828be2ec-78e7-3306-8f70-a4a8fc94b384", false)
		resp := apimsubscription.Response{
			State:  state(new("2026-11-23T15:50:02Z"), nil),
			Legacy: apimsubscription.Legacy{HridContainsUUID: true},
		}

		Expect(apimsubscription.UpdateStatus(sub, resp)).To(Succeed())

		Expect(k8s.IsAutomationAPIManaged(sub)).To(BeFalse())
	})

	DescribeTable("refuses a response without a usable start date, leaving the status as it was",
		func(startingAt *string, message string) {
			sub := newSub("sub", "", false)

			err := apimsubscription.UpdateStatus(sub, apimsubscription.Response{State: state(startingAt, nil)})

			Expect(err).To(MatchError(ContainSubstring(message)))
			Expect(sub.Status.ID).To(BeEmpty())
			Expect(k8s.IsAutomationAPIManaged(sub)).To(BeFalse())
		},
		Entry("missing", nil, "subscription [default-sub] has no start date in APIM"),
		Entry("not a date", new("yesterday"), "parse the start date of subscription [default-sub]"),
	)
})

var _ = Describe("Subscription FromState", func() {
	addressed := func(legacy bool) apimsubscription.DTO {
		return apimsubscription.DTO{ApiHrid: "default-api", Legacy: apimsubscription.Legacy{HridContainsUUID: legacy}}
	}

	It("keeps the hrid APIM answers for a subscription addressed by its hrid", func() {
		state := sdk.SubscriptionState{Id: new("sub-uuid"), Hrid: "default-sub"}

		Expect(apimsubscription.FromState(state, addressed(false)).Hrid).To(Equal("default-sub"))
	})

	It("leaves the hrid empty when APIM answers neither hrid nor id", func() {
		Expect(apimsubscription.FromState(sdk.SubscriptionState{}, addressed(true)).Hrid).To(BeEmpty())
	})
})

var _ = Describe("Subscription SetProcessingStatus", func() {
	It("writes Completed after a successful reconcile", func() {
		sub := newSub("sub", "", true)
		sub.Status.ProcessingStatus = core.ProcessingStatusFailed

		apimsubscription.SetProcessingStatus(sub, nil)

		Expect(sub.Status.ProcessingStatus).To(Equal(core.ProcessingStatusCompleted))
	})

	It("writes Failed after a failed reconcile", func() {
		sub := newSub("sub", "", true)
		sub.Status.ProcessingStatus = core.ProcessingStatusCompleted

		apimsubscription.SetProcessingStatus(sub, errors.New("upsert refused"))

		Expect(sub.Status.ProcessingStatus).To(Equal(core.ProcessingStatusFailed))
	})
})
