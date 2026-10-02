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

package amreporter

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/amreporter"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
)

var _ = Describe("Validate create", func() {
	ctx := context.Background()
	admissionCtrl := amreporter.NewAdmissionCtrl()

	It("should pass dry-run validation", func() {
		fixtures := withDomain()
		reporter := fixtures.AMReporter
		fixtures.AMReporter = nil
		fixtures.Apply()

		warnings, err := admissionCtrl.ValidateCreate(ctx, reporter)
		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(BeEmpty())
	})

	It("should fail dry-run validation when AM rejects it", func() {
		By("restarting the mock AM with dryRunReject=true")
		mockServer.Stop()
		rejectServer := am.Start(true)
		defer func() {
			rejectServer.Stop()
			mockServer = am.Start(false)
		}()

		fixtures := withDomain()
		reporter := fixtures.AMReporter
		fixtures.AMReporter = nil
		fixtures.Apply()

		_, err := admissionCtrl.ValidateCreate(ctx, reporter)
		Expect(err).To(HaveOccurred())
	})

	It("should admit with a warning, without dry-run, when the domain does not exist", func() {
		fixtures := fixture.Builder().
			WithAMContext(constants.AMContextFile).
			WithAMReporter(constants.AMReporterFileFile).
			Build()
		reporter := fixtures.AMReporter
		fixtures.AMReporter = nil
		fixtures.Apply()

		// a dry-run under a missing domain would get a 404 from AM, a severe error
		warnings, err := admissionCtrl.ValidateCreate(ctx, reporter)
		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(ContainElement(ContainSubstring("not found or not yet created in AM")))
	})

	It("should reject a domainRef in another namespace", func() {
		reporter := withDomain().AMReporter
		reporter.Spec.DomainRef.Namespace = "other"

		_, err := admissionCtrl.ValidateCreate(ctx, reporter)
		Expect(err).To(MatchError(ContainSubstring("domainRef.namespace")))
	})

	It("should reject a name that does not make a valid AM key", func() {
		reporter := withDomain().AMReporter
		reporter.Name = "reporter.with.dots"

		_, err := admissionCtrl.ValidateCreate(ctx, reporter)
		Expect(err).To(HaveOccurred())
	})

	It("should admit a system reporter with a warning when name, type and configuration are set", func() {
		fixtures := withDomain()
		reporter := fixtures.AMReporter
		reporter.Spec.System = new(true)
		fixtures.AMReporter = nil
		fixtures.Apply()

		warnings, err := admissionCtrl.ValidateCreate(ctx, reporter)
		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(ConsistOf("'configuration', 'name', 'type' will be ignored when 'system' is 'true'."))
	})
})
