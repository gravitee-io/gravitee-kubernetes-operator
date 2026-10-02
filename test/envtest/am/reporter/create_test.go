// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package amreporter

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	internal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/reporter"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Create", func() {
	ctx := context.Background()

	It("should create the reporter under its domain, owned by the domain", func() {
		fixtures := fixtures().Apply()
		reporter := fixtures.AMReporter

		By("expecting the reporter to be accepted")

		Expect(assert.AMReporterAccepted(reporter)).To(Succeed())
		dto, err := internal.ToReporterDTO(reporter)
		Expect(err).ToNot(HaveOccurred())
		Expect(reporter.Status.Key).To(Equal(dto.Key))
		Expect(reporter.Status.DomainKey).To(Equal(dto.DomainKey))
		Expect(reporter.Status.Name).To(Equal(*reporter.Spec.Name))
		Expect(reporter.Status.Type).To(Equal(*reporter.Spec.Type))

		By("expecting the reporter in AM, under the domain")

		status, remote := remoteReporter(ctx, reporter)
		Expect(status).To(Equal(http.StatusOK))
		Expect(remote.Name).To(Equal(reporter.Spec.Name))
		Expect(remote.Configuration.Object).To(HaveKeyWithValue("filename", "audit"))

		By("expecting the domain to own the reporter")

		Expect(reporter.OwnerReferences).To(HaveLen(1))
		owner := reporter.OwnerReferences[0]
		Expect(owner.Kind).To(Equal("AMSecurityDomain"))
		Expect(owner.Name).To(Equal(fixtures.AMSecurityDomain.Name))
		Expect(owner.UID).To(Equal(fixtures.AMSecurityDomain.UID))
	})

	It("should wait for its domain when applied before it", func() {
		fixtures := fixtures()
		domain, reporter := fixtures.AMSecurityDomain, fixtures.AMReporter
		fixtures.AMSecurityDomain, fixtures.AMReporter = nil, nil
		fixtures.Apply()

		By("creating the reporter before its domain")

		Expect(manager.Client().Create(ctx, reporter)).To(Succeed())

		Eventually(ctx, func() error {
			if err := manager.GetLatest(ctx, reporter); err != nil {
				return err
			}
			refs := k8s.MapConditions(reporter.Status.Conditions)[k8s.ConditionResolvedRefs]
			return assert.Equals("ResolvedRefs", metaV1.ConditionFalse, refs.Status)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), reporter.Name)

		By("creating the domain")

		Expect(manager.Client().Create(ctx, domain)).To(Succeed())

		By("expecting the reporter to be accepted once the domain is in AM")

		Eventually(ctx, func() error {
			if err := manager.GetLatest(ctx, reporter); err != nil {
				return err
			}
			return assert.AMReporterAccepted(reporter)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), reporter.Name)

		status, _ := remoteReporter(ctx, reporter)
		Expect(status).To(Equal(http.StatusOK))
	})
})
