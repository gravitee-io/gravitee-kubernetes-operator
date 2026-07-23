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

package gateway_test

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/gateway"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func newTestGateway(generation int64) *gateway.Gateway {
	return gateway.WrapGateway(&gwAPIv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "test-gw",
			Namespace:  "default",
			Generation: generation,
		},
		Spec: gwAPIv1.GatewaySpec{
			GatewayClassName: "test-class",
			Listeners: []gwAPIv1.Listener{
				{
					Name:     "http",
					Port:     80,
					Protocol: gwAPIv1.HTTPProtocolType,
				},
			},
		},
	})
}

var _ = Describe("CarryForwardConditionGenerations", func() {
	It("updates ObservedGeneration on existing sync conditions", func() {
		gw := newTestGateway(3)

		k8s.SetCondition(gw,
			k8s.NewAutoscalingSyncConditionBuilder(1).
				Message("HorizontalPodAutoscaler synced successfully").
				Build(),
		)
		k8s.SetCondition(gw,
			k8s.NewPDBSyncConditionBuilder(1).
				Message("PodDisruptionBudget synced successfully").
				Build(),
		)

		k8s.CarryForwardConditionGenerations(
			gw,
			gw.Object.Generation,
			k8s.ConditionAutoscalingSync,
			k8s.ConditionPDBSync,
		)

		autoscaling := k8s.GetCondition(gw, k8s.ConditionAutoscalingSync)
		Expect(autoscaling).NotTo(BeNil())
		Expect(autoscaling.ObservedGeneration).To(Equal(int64(3)))
		Expect(autoscaling.Status).To(Equal(metav1.ConditionTrue))

		pdb := k8s.GetCondition(gw, k8s.ConditionPDBSync)
		Expect(pdb).NotTo(BeNil())
		Expect(pdb.ObservedGeneration).To(Equal(int64(3)))
		Expect(pdb.Status).To(Equal(metav1.ConditionTrue))
	})

	It("is a no-op when sync conditions are absent", func() {
		gw := newTestGateway(1)

		k8s.CarryForwardConditionGenerations(
			gw,
			gw.Object.Generation,
			k8s.ConditionAutoscalingSync,
			k8s.ConditionPDBSync,
		)

		Expect(k8s.GetCondition(gw, k8s.ConditionAutoscalingSync)).To(BeNil())
		Expect(k8s.GetCondition(gw, k8s.ConditionPDBSync)).To(BeNil())
	})
})
