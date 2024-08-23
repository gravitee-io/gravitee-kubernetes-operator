package application

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/application"
	admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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

var _ = Describe("Validate create", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := admission.AdmissionCtrl{}

	It("should return error on create with wrong settings", func() {
		fixtures := fixture.
			Builder().
			WithApplication(constants.Application).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("adding oauth settings while simple settings are set")

		fixtures.Application.Spec.Settings.Oauth = &application.OAuthClientSettings{
			ApplicationType: "WEB",
		}

		By("checking that application validation returns error")

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Application, fixtures.Application)
			return assert.Equals(
				"error",
				errors.NewSevere(
					"configuring both OAuth and simple settings is not allowed",
				),
				err,
			)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
