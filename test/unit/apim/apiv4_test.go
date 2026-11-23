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
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apiv4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

const legacyAPIID = "6b1d5f0a-0000-0000-0000-000000000001"

// apiV4 builds an ApiV4Definition named orders in default from a spec given as JSON.
func apiV4(spec string) *v1alpha1.ApiV4Definition {
	GinkgoHelper()
	api := &v1alpha1.ApiV4Definition{
		TypeMeta:   metav1.TypeMeta{Kind: "ApiV4Definition", APIVersion: "gravitee.io/v1alpha1"},
		ObjectMeta: metav1.ObjectMeta{Name: "orders", Namespace: "default", UID: "orders-uid"},
	}
	Expect(json.Unmarshal([]byte(spec), &api.Spec)).To(Succeed())
	return api
}

// apiV4WithContext is apiV4 with a management context and a keyless plan, syncing from Kubernetes.
func apiV4WithContext() *v1alpha1.ApiV4Definition {
	GinkgoHelper()
	return apiV4(`{
		"name": "orders", "version": "1", "type": "PROXY",
		"contextRef": {"name": "dev", "namespace": "default"},
		"definitionContext": {"origin": "KUBERNETES", "syncFrom": "KUBERNETES"},
		"plans": {"KeyLess": {"name": "KeyLess", "status": "PUBLISHED", "security": {"type": "KEY_LESS"}}}
	}`)
}

// registerFakeK8s registers a fake Kubernetes client holding objects, for the gateway ConfigMap sync.
func registerFakeK8s(objects ...runtime.Object) {
	GinkgoHelper()
	scheme := runtime.NewScheme()
	Expect(clientgoscheme.AddToScheme(scheme)).To(Succeed())
	k8s.RegisterClient(fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(objects...).Build())
}

func toAPIV4DTO(api *v1alpha1.ApiV4Definition) apiv4.DTO {
	GinkgoHelper()
	dto, err := apiv4.ToDTO(api)
	Expect(err).ToNot(HaveOccurred())
	return dto
}

func ordersHRID() string {
	return (&refs.NamespacedName{Namespace: "default", Name: "orders"}).HRID()
}

var _ = Describe("ApiV4Definition APIM mapping", func() {
	Describe("ToDTO identity", func() {
		It("sends the namespaced name as hrid for an API with a context and no spec.id", func() {
			dto := toAPIV4DTO(apiV4WithContext())

			Expect(dto.Hrid).To(Equal(ordersHRID()))
			Expect(dto.HridContainsUUID).To(BeFalse())
			Expect(dto.Identity()).To(Equal(ordersHRID()))
			Expect(dto.Synced).To(BeFalse())
		})

		It("sends the spec.id as hrid, flagged as a UUID, for an API with a context and a spec.id", func() {
			api := apiV4WithContext()
			api.Spec.ID = legacyAPIID

			dto := toAPIV4DTO(api)

			Expect(dto.Hrid).To(Equal(legacyAPIID))
			Expect(dto.HridContainsUUID).To(BeTrue())
		})

		It("does not flag a UUID for an API without a context, whose spec.id is generated", func() {
			api := apiV4WithContext()
			api.Spec.Context = nil

			dto := toAPIV4DTO(api)

			Expect(dto.Hrid).To(Equal(ordersHRID()))
			Expect(dto.HridContainsUUID).To(BeFalse())
			Expect(api.Spec.ID).To(BeEmpty(), "ToDTO must not populate the IDs of the object it was given")
			Expect(string(dto.Gateway.Definition)).To(ContainSubstring(`"id":"`))
		})

		It("is synced once APIM gave the API an ID", func() {
			api := apiV4WithContext()
			api.Status.ID = legacyAPIID
			Expect(toAPIV4DTO(api).Synced).To(BeTrue())
		})
	})

	Describe("ToDTO delete identity", func() {
		It("deletes an Automation API managed API by hrid", func() {
			api := apiV4WithContext()
			api.Status.ID = legacyAPIID
			k8s.AddAutomationAPIManagedCondition(api)

			dto := toAPIV4DTO(api)

			Expect(dto.DeleteID).To(Equal(ordersHRID()))
			Expect(dto.DeleteHridContainsUUID).To(BeFalse())
		})

		It("deletes a legacy API by the UUID in its status", func() {
			api := apiV4WithContext()
			api.Status.ID = legacyAPIID

			dto := toAPIV4DTO(api)

			Expect(dto.DeleteID).To(Equal(legacyAPIID))
			Expect(dto.DeleteHridContainsUUID).To(BeTrue())
		})

		It("deletes an API never synced by hrid", func() {
			dto := toAPIV4DTO(apiV4WithContext())

			Expect(dto.DeleteID).To(Equal(ordersHRID()))
			Expect(dto.DeleteHridContainsUUID).To(BeFalse())
		})
	})

	It("reads ignore-groups from the annotation", func() {
		api := apiV4WithContext()
		Expect(toAPIV4DTO(api).IgnoreGroups).To(BeFalse())

		api.Annotations = map[string]string{core.IgnoreGroupsAnnotation: "true"}
		Expect(toAPIV4DTO(api).IgnoreGroups).To(BeTrue())
	})

	Describe("plan send order", func() {
		plans := func(dto apiv4.DTO) []string {
			hrids := make([]string, 0, len(dto.Plans))
			for _, plan := range dto.Plans {
				hrids = append(hrids, plan.Hrid)
			}
			return hrids
		}

		It("sends plans by order, then by hrid, and ranks them for drift when all are published and ordered", func() {
			dto := toAPIV4DTO(apiV4(`{
				"name": "orders", "version": "1", "type": "PROXY",
				"contextRef": {"name": "dev", "namespace": "default"},
				"plans": {
					"a": {"name": "a", "order": 20, "status": "PUBLISHED", "security": {"type": "KEY_LESS"}},
					"b": {"name": "b", "order": 10, "status": "PUBLISHED", "security": {"type": "API_KEY"}},
					"c": {"name": "c", "order": 10, "status": "PUBLISHED", "security": {"type": "JWT"}}
				}
			}`))

			Expect(plans(dto)).To(Equal([]string{"b", "c", "a"}))
			Expect(dto.PlanOrder).To(Equal(map[string]apiv4.PlanRank{"b": {Order: 1}, "c": {Order: 2}, "a": {Order: 3}}))
		})

		It("leaves the plan order out of drift when a plan has no order", func() {
			dto := toAPIV4DTO(apiV4(`{
				"name": "orders", "version": "1", "type": "PROXY",
				"contextRef": {"name": "dev", "namespace": "default"},
				"plans": {
					"a": {"name": "a", "order": 1, "status": "PUBLISHED", "security": {"type": "KEY_LESS"}},
					"b": {"name": "b", "status": "PUBLISHED", "security": {"type": "API_KEY"}}
				}
			}`))

			Expect(plans(dto)).To(Equal([]string{"a", "b"}))
			Expect(dto.PlanOrder).To(BeNil())
		})

		It("leaves the plan order out of drift when a plan is not published", func() {
			dto := toAPIV4DTO(apiV4(`{
				"name": "orders", "version": "1", "type": "PROXY",
				"contextRef": {"name": "dev", "namespace": "default"},
				"plans": {
					"a": {"name": "a", "order": 1, "status": "PUBLISHED", "security": {"type": "KEY_LESS"}},
					"b": {"name": "b", "order": 2, "status": "STAGING", "security": {"type": "API_KEY"}}
				}
			}`))

			Expect(dto.PlanOrder).To(BeNil())
		})
	})

	Describe("GatewaySync", func() {
		It("enables the ConfigMap of an API that syncs from Kubernetes, owned by the API", func() {
			gateway := toAPIV4DTO(apiV4WithContext()).Gateway

			Expect(gateway.Enabled).To(BeTrue())
			Expect(gateway.Stopped).To(BeFalse())
			Expect(gateway.Name).To(Equal("orders"))
			Expect(gateway.Namespace).To(Equal("default"))
			Expect(gateway.Owner).To(Equal(metav1.OwnerReference{
				Kind: "ApiV4Definition", Name: "orders", APIVersion: "gravitee.io/v1alpha1", UID: "orders-uid",
			}))
			Expect(string(gateway.Definition)).To(ContainSubstring(`"name":"orders"`))
		})

		It("disables the ConfigMap of an API that syncs from management", func() {
			api := apiV4WithContext()
			api.Spec.DefinitionContext = nil

			gateway := toAPIV4DTO(api).Gateway

			Expect(gateway.Enabled).To(BeFalse())
			Expect(gateway.Definition).To(BeNil())
			Expect(gateway.Owner.UID).To(Equal(types.UID("orders-uid")))
		})

		It("marks a stopped API", func() {
			api := apiV4WithContext()
			api.Spec.State = "STOPPED"
			Expect(toAPIV4DTO(api).Gateway.Stopped).To(BeTrue())
		})
	})

	Describe("Upsert gateway ConfigMap", func() {
		var ctx context.Context

		cmKey := types.NamespacedName{Name: "orders", Namespace: "default"}

		configMap := func(owner types.UID) *corev1.ConfigMap {
			return &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
				Name: "orders", Namespace: "default",
				OwnerReferences: []metav1.OwnerReference{{
					Kind: "ApiV4Definition", Name: "orders", APIVersion: "gravitee.io/v1alpha1", UID: owner,
				}},
			}}
		}

		BeforeEach(func() {
			ctx = context.Background()
		})

		It("writes the ConfigMap of an API without context that syncs from Kubernetes", func() {
			registerFakeK8s()
			api := apiV4WithContext()
			api.Spec.Context = nil

			_, err := apiv4.Upsert(ctx, nil, toAPIV4DTO(api))
			Expect(err).ToNot(HaveOccurred())

			cm := &corev1.ConfigMap{}
			Expect(k8s.GetClient().Get(ctx, cmKey, cm)).To(Succeed())
			Expect(cm.Data).To(HaveKeyWithValue("apiDefinitionVersion", "4.0.0"))
			Expect(cm.Data).To(HaveKeyWithValue("organizationId", "DEFAULT"))
			Expect(cm.Data).To(HaveKeyWithValue("environmentId", "DEFAULT"))
			Expect(cm.Data["definition"]).To(ContainSubstring(`"name":"orders"`))
			Expect(cm.OwnerReferences).To(HaveLen(1))
			Expect(cm.OwnerReferences[0].UID).To(Equal(types.UID("orders-uid")))
		})

		It("deletes the ConfigMap of a stopped API", func() {
			registerFakeK8s(configMap("orders-uid"))
			api := apiV4WithContext()
			api.Spec.Context = nil
			api.Spec.State = "STOPPED"

			_, err := apiv4.Upsert(ctx, nil, toAPIV4DTO(api))
			Expect(err).ToNot(HaveOccurred())
			Expect(kerrors.IsNotFound(k8s.GetClient().Get(ctx, cmKey, &corev1.ConfigMap{}))).To(BeTrue())
		})

		It("deletes the ConfigMap it owns when the API moves to management sync", func() {
			registerFakeK8s(configMap("orders-uid"))
			api := apiV4WithContext()
			api.Spec.Context = nil
			api.Spec.DefinitionContext = nil

			_, err := apiv4.Upsert(ctx, nil, toAPIV4DTO(api))
			Expect(err).ToNot(HaveOccurred())
			Expect(kerrors.IsNotFound(k8s.GetClient().Get(ctx, cmKey, &corev1.ConfigMap{}))).To(BeTrue())
		})

		It("keeps a ConfigMap of the same name that another object owns", func() {
			registerFakeK8s(configMap("someone-else"))
			api := apiV4WithContext()
			api.Spec.Context = nil
			api.Spec.DefinitionContext = nil

			_, err := apiv4.Upsert(ctx, nil, toAPIV4DTO(api))
			Expect(err).ToNot(HaveOccurred())
			Expect(k8s.GetClient().Get(ctx, cmKey, &corev1.ConfigMap{})).To(Succeed())
		})
	})

	Describe("PayloadChanged", func() {
		It("is false for the same spec and annotations", func() {
			Expect(apiv4.PayloadChanged(apiV4WithContext(), apiV4WithContext())).To(BeFalse())
		})

		It("is true when the spec changes", func() {
			changed := apiV4WithContext()
			changed.Spec.Description = new("changed")
			Expect(apiv4.PayloadChanged(apiV4WithContext(), changed)).To(BeTrue())
		})

		It("is true when the ignore-groups annotation changes", func() {
			changed := apiV4WithContext()
			changed.Annotations = map[string]string{core.IgnoreGroupsAnnotation: "true"}
			Expect(apiv4.PayloadChanged(apiV4WithContext(), changed)).To(BeTrue())
		})

		It("is false when another annotation changes", func() {
			changed := apiV4WithContext()
			changed.Annotations = map[string]string{"example.com/other": "x"}
			Expect(apiv4.PayloadChanged(apiV4WithContext(), changed)).To(BeFalse())
		})
	})
})
