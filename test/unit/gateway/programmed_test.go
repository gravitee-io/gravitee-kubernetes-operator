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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// The gateway reconciler requeues for as long as IsProgrammed reports false, so that a gateway
// waiting on its deployment is polled instead of relying only on the deployment watch. A gateway
// that reports programmed while its deployment is still pending would stop being requeued and
// stay stuck, so each condition state below is checked explicitly.
var _ = Describe("IsProgrammed", func() {
	It("is false when the programmed condition is missing", func() {
		gw := newTestGateway(1)

		Expect(k8s.IsProgrammed(gw)).To(BeFalse())
	})

	It("is false while the deployment is not ready", func() {
		gw := newTestGateway(1)

		k8s.SetCondition(gw,
			k8s.NewGatewayProgrammedConditionBuilder(1).
				Pending("waiting for gateway deployment to become ready").
				Build(),
		)

		Expect(k8s.IsProgrammed(gw)).To(BeFalse())
	})

	It("is true once all listeners have been programmed", func() {
		gw := newTestGateway(1)

		k8s.SetCondition(gw,
			k8s.NewGatewayProgrammedConditionBuilder(1).
				Program("all listeners have been programmed").
				Build(),
		)

		Expect(k8s.IsProgrammed(gw)).To(BeTrue())
	})

	It("goes back to false when the gateway stops being programmed", func() {
		gw := newTestGateway(1)

		k8s.SetCondition(gw,
			k8s.NewGatewayProgrammedConditionBuilder(1).
				Program("all listeners have been programmed").
				Build(),
		)
		Expect(k8s.IsProgrammed(gw)).To(BeTrue())

		k8s.SetCondition(gw,
			k8s.NewGatewayProgrammedConditionBuilder(2).
				Pending("waiting for gateway deployment to become ready").
				Build(),
		)

		Expect(k8s.IsProgrammed(gw)).To(BeFalse())
	})
})
