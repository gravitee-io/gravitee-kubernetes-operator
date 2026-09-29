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
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
)

var _ = Describe("SetServiceType", func() {
	lbClass := "example.com/lb"

	// The CRD defaults externalTrafficPolicy to Cluster, so every params object carries it.
	params := func(serviceType coreV1.ServiceType) *gateway.Service {
		return &gateway.Service{
			Type:                  ptr.To(serviceType),
			ExternalTrafficPolicy: coreV1.ServiceExternalTrafficPolicyCluster,
			LoadBalancerClass:     ptr.To(lbClass),
		}
	}

	It("keeps externalTrafficPolicy and loadBalancerClass on a LoadBalancer service", func() {
		svc := k8s.DefaultService.DeepCopy()
		k8s.SetServiceType(svc, params(coreV1.ServiceTypeLoadBalancer))

		Expect(svc.Spec.Type).To(Equal(coreV1.ServiceTypeLoadBalancer))
		Expect(svc.Spec.ExternalTrafficPolicy).To(Equal(coreV1.ServiceExternalTrafficPolicyCluster))
		Expect(svc.Spec.LoadBalancerClass).To(Equal(ptr.To(lbClass)))
	})

	It("keeps externalTrafficPolicy and drops loadBalancerClass on a NodePort service", func() {
		svc := k8s.DefaultService.DeepCopy()
		k8s.SetServiceType(svc, params(coreV1.ServiceTypeNodePort))

		Expect(svc.Spec.Type).To(Equal(coreV1.ServiceTypeNodePort))
		Expect(svc.Spec.ExternalTrafficPolicy).To(Equal(coreV1.ServiceExternalTrafficPolicyCluster))
		Expect(svc.Spec.LoadBalancerClass).To(BeNil())
	})

	It("drops externalTrafficPolicy and loadBalancerClass on a ClusterIP service", func() {
		svc := k8s.DefaultService.DeepCopy()
		k8s.SetServiceType(svc, params(coreV1.ServiceTypeClusterIP))

		Expect(svc.Spec.Type).To(Equal(coreV1.ServiceTypeClusterIP))
		Expect(svc.Spec.ExternalTrafficPolicy).To(BeEmpty())
		Expect(svc.Spec.LoadBalancerClass).To(BeNil())
	})

	It("clears the values an existing LoadBalancer service held when it moves to ClusterIP", func() {
		svc := k8s.DefaultService.DeepCopy()
		k8s.SetServiceType(svc, params(coreV1.ServiceTypeLoadBalancer))
		k8s.SetServiceType(svc, params(coreV1.ServiceTypeClusterIP))

		Expect(svc.Spec.Type).To(Equal(coreV1.ServiceTypeClusterIP))
		Expect(svc.Spec.ExternalTrafficPolicy).To(BeEmpty())
		Expect(svc.Spec.LoadBalancerClass).To(BeNil())
	})
})
