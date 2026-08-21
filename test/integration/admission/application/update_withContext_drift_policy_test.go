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

package application

import (
	"context"
	"fmt"
	"strings"

	admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate drift policies", labels.WithContext, func() {
	ctx := context.Background()
	admissionCtrl := admission.AdmissionCtrl{}

	var original env.DriftDetection
	var originalTimeout int

	BeforeEach(func() {
		original = env.Config.DriftDetection
		originalTimeout = env.Config.HTTPClientTimeoutSeconds
		env.Config.DriftDetection.Enabled = true
		env.Config.DriftDetection.Policy = env.DriftPolicyDeny
		env.Config.DriftDetection.OnRemoteMissing = env.DriftPolicyDeny
		env.Config.DriftDetection.FetchFailurePolicy = env.DriftPolicyDeny
	})

	AfterEach(func() {
		env.Config.DriftDetection = original
		env.Config.HTTPClientTimeoutSeconds = originalTimeout
	})

	It("should warn and allow the update when drift policy is warn", func() {
		fixtures := applyMinimalApplication()
		newApp := fixtures.Application.DeepCopy()
		validateDescriptionDrift(ctx, admissionCtrl, fixtures.Application, newApp, fixtures.Context)

		env.Config.DriftDetection.Policy = env.DriftPolicyWarn

		warnings, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Application, newApp)
		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(HaveLen(1))
		Expect(warnings[0]).To(ContainSubstring("drift detected"))
		Expect(warnings[0]).To(ContainSubstring(driftDescriptionAssert))
	})

	It("should allow the update with no admission output when drift policy is allow", func() {
		fixtures := applyMinimalApplication()
		newApp := fixtures.Application.DeepCopy()
		validateDescriptionDrift(ctx, admissionCtrl, fixtures.Application, newApp, fixtures.Context)

		env.Config.DriftDetection.Policy = env.DriftPolicyAllow

		warnings, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Application, newApp)
		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(BeEmpty())
	})

	It("should warn and allow the update when the remote is missing and policy is warn", func() {
		fixtures := applyMinimalApplication()
		env.Config.DriftDetection.OnRemoteMissing = env.DriftPolicyWarn

		newApp := fixtures.Application.DeepCopy()
		setDescription(newApp, updatedDescription)
		newApp.Name = random.GetName() // this triggers 404 on the GET during drift

		Eventually(func() error {
			warnings, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Application, newApp)
			if err != nil {
				return err
			}
			if err = assert.SliceOfSize("warnings", warnings, 1); err != nil {
				return err
			}
			return contains(warnings[0], "not found during drift detection")
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())
	})

	It("should allow the update with no admission output when the remote is missing and policy is allow", func() {
		fixtures := applyMinimalApplication()
		env.Config.DriftDetection.OnRemoteMissing = env.DriftPolicyAllow

		newApp := fixtures.Application.DeepCopy()
		setDescription(newApp, updatedDescription)
		newApp.Name = random.GetName() // this triggers 404 on the GET during drift

		warnings, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Application, newApp)
		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(BeEmpty())
	})

})

func applyMinimalApplication() *fixture.Objects {
	GinkgoHelper()
	return fixture.
		Builder().
		WithApplication(constants.Application).
		WithContext(constants.ContextWithCredentialsFile).
		Build().
		Apply()
}

func contains(actual, substring string) error {
	if !strings.Contains(actual, substring) {
		return fmt.Errorf("expected %q to contain %q", actual, substring)
	}
	return nil
}
