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

package amclientsecret

import (
	"context"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/amclientsecret"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/am"
)

var _ = Describe("Validate create", func() {
	ctx := context.Background()
	admissionCtrl := amclientsecret.NewAdmissionCtrl()

	It("should pass dry-run validation", func() {
		fixtures := applyAllButClientSecret()
		cs := fixtures.AMClientSecret
		createSecret(ctx, cs, value)
		Expect(fixtures.AMApplication.Status.Key).ToNot(BeEmpty(), "the application must be in AM for the dry-run to run")

		warnings, err := admissionCtrl.ValidateCreate(ctx, cs)
		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(BeEmpty())
	})

	It("should fail dry-run validation when AM rejects it (AC8)", func() {
		By("restarting the mock AM with dryRunReject=true: its dry-run does not check the value rules")
		mockServer.Stop()
		rejectServer := am.Start(true)
		defer func() {
			rejectServer.Stop()
			mockServer = am.Start(false)
		}()

		cs := applyAllButClientSecret().AMClientSecret
		createSecret(ctx, cs, value)

		_, err := admissionCtrl.ValidateCreate(ctx, cs)
		Expect(err).To(MatchError(ContainSubstring("dry-run-reject mode")), "AM's own message reaches the user")
		Expect(err.Error()).ToNot(ContainSubstring(value))
	})

	It("should admit with a warning, without dry-run, when the application does not exist (AC7)", func() {
		fixtures := withApplication()
		cs := fixtures.AMClientSecret
		fixtures.AMApplication, fixtures.AMClientSecret = nil, nil
		fixtures.Apply()

		warnings, err := admissionCtrl.ValidateCreate(ctx, cs)
		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(ContainElement(ContainSubstring("will be created in AM once it exists")))
	})

	It("should admit with a warning, without dry-run, when the Secret does not exist (AC7)", func() {
		cs := applyAllButClientSecret().AMClientSecret

		warnings, err := admissionCtrl.ValidateCreate(ctx, cs)
		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(ContainElement(ContainSubstring("Secret [" + cs.Spec.SecretRef.Name + "] not found")))
	})

	It("should reject a key longer than the 64 characters AM allows for a client secret", func() {
		cs := withApplication().AMClientSecret
		cs.Name = strings.Repeat("a", 64)

		_, err := admissionCtrl.ValidateCreate(ctx, cs)
		Expect(err).To(MatchError(ContainSubstring("longer than the 64 characters")))
	})
})
