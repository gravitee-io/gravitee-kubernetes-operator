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

var _ = Describe("Validate create", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := admission.AdmissionCtrl{}

	It("should resolve client certificate refs from secret and configmap", func() {
		fixtures := fixture.
			Builder().
			WithApplication(constants.ApplicationWithClientCertsRefs).
			WithContext(constants.ContextWithCredentialsFile).
			AddSecret(constants.ClientCertSecretFile).
			AddConfigMap(constants.ClientCertConfigMapFile).
			Build().
			Apply()

		By("checking that application validation succeeds (refs are resolved)")

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, fixtures.Application)
			return assert.Equals("error", nil, err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return error when certificate ref points to a CA certificate", func() {
		fixtures := fixture.
			Builder().
			WithApplication(constants.Application).
			WithContext(constants.ContextWithCredentialsFile).
			AddSecret(constants.CACertSecretFile).
			Build().
			Apply()

		By("setting a clientCertificates entry referencing the CA secret")

		fixtures.Application.Spec.Settings.TLS = &application.TLSSettings{
			ClientCertificates: []application.ClientCertificate{
				{
					Name: "ca-cert",
					Ref: &application.CertificateRef{
						Kind: "secrets",
						Name: "ca-cert-secret",
						Key:  "tls.crt",
					},
				},
			},
		}

		By("checking that application validation returns CA certificate error")

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, fixtures.Application)
			return assert.Equals(
				"error",
				errors.NewSevere(
					"certificate is a CA certificate",
				),
				err,
			)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
