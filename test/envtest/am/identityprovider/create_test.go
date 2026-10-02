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

package amidentityprovider

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	internal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/identityprovider"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Create", func() {
	ctx := context.Background()

	It("should create the identity provider under its domain, owned by the domain", func() {
		fixtures := fixtures().Apply()
		idp := fixtures.AMIdentityProvider

		By("expecting the identity provider to be accepted")

		Expect(assert.AMIdentityProviderAccepted(idp)).To(Succeed())
		dto, err := internal.ToIdentityProviderDTO(idp)
		Expect(err).ToNot(HaveOccurred())
		Expect(idp.Status.Key).To(Equal(dto.Key))
		Expect(idp.Status.DomainKey).To(Equal(dto.DomainKey))
		Expect(idp.Status.Name).To(Equal(*idp.Spec.Name))
		Expect(idp.Status.Type).To(Equal(*idp.Spec.Type))

		By("expecting the identity provider in AM, under the domain")

		status, remote := remoteIdentityProvider(ctx, idp)
		Expect(status).To(Equal(http.StatusOK))
		Expect(remote.Name).To(Equal(idp.Spec.Name))

		By("expecting the domain to own the identity provider")

		Expect(idp.OwnerReferences).To(HaveLen(1))
		owner := idp.OwnerReferences[0]
		Expect(owner.Kind).To(Equal("AMSecurityDomain"))
		Expect(owner.Name).To(Equal(fixtures.AMSecurityDomain.Name))
		Expect(owner.UID).To(Equal(fixtures.AMSecurityDomain.UID))
	})

	It("should wait for its domain when applied before it", func() {
		fixtures := fixtures()
		domain, idp := fixtures.AMSecurityDomain, fixtures.AMIdentityProvider
		fixtures.AMSecurityDomain, fixtures.AMIdentityProvider = nil, nil
		fixtures.Apply()

		By("creating the identity provider before its domain")

		Expect(manager.Client().Create(ctx, idp)).To(Succeed())

		Eventually(ctx, func() error {
			if err := manager.GetLatest(ctx, idp); err != nil {
				return err
			}
			refs := k8s.MapConditions(idp.Status.Conditions)[k8s.ConditionResolvedRefs]
			return assert.Equals("ResolvedRefs", metaV1.ConditionFalse, refs.Status)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), idp.Name)

		By("creating the domain")

		Expect(manager.Client().Create(ctx, domain)).To(Succeed())

		By("expecting the identity provider to be accepted once the domain is in AM")

		Eventually(ctx, func() error {
			if err := manager.GetLatest(ctx, idp); err != nil {
				return err
			}
			return assert.AMIdentityProviderAccepted(idp)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), idp.Name)

		status, _ := remoteIdentityProvider(ctx, idp)
		Expect(status).To(Equal(http.StatusOK))
	})
})
