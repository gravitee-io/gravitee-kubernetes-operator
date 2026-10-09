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

	"github.com/gravitee-io/gravitee-automation-sdk/common/pkg/unstructured"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/amreporter"
	internal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/reporter"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
)

var _ = Describe("Validate drift", func() {
	ctx := context.Background()
	admissionCtrl := amreporter.NewAdmissionCtrl()

	// inAM applies the domain only, edits the reporter spec with local, and stores the reporter in the
	// AM mock as AM returns it: the spec, with the fields remote overrides.
	inAM := func(local func(*v1alpha1.AMReporterSpec), remote func(*internal.Reporter)) *v1alpha1.AMReporter {
		fixtures := withDomain()
		reporter := fixtures.AMReporter
		local(&reporter.Spec)
		fixtures.AMReporter = nil
		fixtures.Apply()

		dto, err := internal.ToReporterDTO(reporter)
		Expect(err).ToNot(HaveOccurred())
		remote(&dto)
		resp, err := am.NewSDKClient().UpsertReporterWithResponse(ctx, dto.DomainKey, nil, dto.Reporter)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.JSON200).ToNot(BeNil(), string(resp.Body))
		return reporter
	}
	validateUpdate := func(reporter *v1alpha1.AMReporter) func() error {
		return func() error {
			updated := reporter.DeepCopy()
			updated.Annotations = map[string]string{"updated": "true"}
			_, err := admissionCtrl.ValidateUpdate(ctx, reporter, updated)
			return err
		}
	}

	// Set, name, type, configuration and the attribute mappings are ignored by AM for the system reporter,
	// which returns its own.
	It("should not drift on the fields AM replaces for the system reporter", func() {
		reporter := inAM(func(spec *v1alpha1.AMReporterSpec) {
			spec.System = new(true)
		}, func(remote *internal.Reporter) {
			remote.Name = new("Default")
			remote.Type = new("mongodb")
			remote.Configuration = unstructured.StringifiedFrom(map[string]any{"uri": "mongodb://localhost"})
		})

		Consistently(validateUpdate(reporter), constants.ConsistentTimeout, constants.Interval).Should(Succeed())
	})

	It("should detect drift on the name of a regular reporter", func() {
		reporter := inAM(func(*v1alpha1.AMReporterSpec) {}, func(remote *internal.Reporter) {
			remote.Name = new("Renamed in AM")
		})

		Eventually(validateUpdate(reporter), constants.EventualTimeout, constants.Interval).
			Should(MatchError(And(ContainSubstring("drift detected"), ContainSubstring("Renamed in AM"))))
	})
})
