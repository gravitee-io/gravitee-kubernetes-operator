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

package am

import (
	"time"

	"github.com/gravitee-io/gravitee-automation-sdk/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/clientsecret"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Client secret drift detection", func() {
	It("does not drift on the value AM never returns, the expiration it sets, or the accepted flag", func() {
		expectNoDrift(drift.DetectWithNamespace(clientSecretForDrift(), remoteClientSecret(), ""))
	})

	It("detects drift on another key", func() {
		remote := remoteClientSecret()
		remote.Key = "ns-other"

		result := drift.DetectWithNamespace(clientSecretForDrift(), remote, "")
		Expect(result.DriftDetected()).To(BeTrue())
	})
})

// clientSecretForDrift is the DTO of an accepted client secret, its value resolved.
func clientSecretForDrift() clientsecret.ClientSecret {
	GinkgoHelper()
	cs := &v1alpha1.AMClientSecret{ObjectMeta: metav1.ObjectMeta{Name: "mobile-2026", Namespace: "ns"}}
	cs.Spec.ApplicationRef = refs.NamespacedName{Name: "mobile"}
	cs.Status.Key = "ns-mobile-2026"
	cs.Resolved = v1alpha1.ClientSecretResolved{DomainKey: "ns-domain", Value: "jFXuh6wb1_dhEu5h8Yi4Dp0fC3ZboI0I-SDoaGE22eI"}
	dto, err := clientsecret.ToClientSecretDTO(cs)
	Expect(err).ToNot(HaveOccurred())
	return dto
}

// remoteClientSecret is the client secret as GetRemote returns it: no value, an expiration.
func remoteClientSecret() clientsecret.ClientSecret {
	expiresAt := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	return clientsecret.ClientSecret{
		ClientSecret:   sdk.ClientSecret{Key: "ns-mobile-2026", ExpiresAt: &expiresAt},
		DomainKey:      "ns-domain",
		ApplicationKey: "ns-mobile",
	}
}
