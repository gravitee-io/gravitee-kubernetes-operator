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

	It("should return error when clientCertificate and clientCertificates are both set", func() {
		fixtures := fixture.
			Builder().
			WithApplication(constants.Application).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("setting both clientCertificate and clientCertificates")

		fixtures.Application.Spec.Settings.TLS = &application.TLSSettings{
			ClientCertificate: "some-cert",
			ClientCertificates: []application.ClientCertificate{
				{
					Name:    "cert1",
					Content: "some-content",
				},
			},
		}

		By("checking that application validation returns error")

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, fixtures.Application)
			return assert.Equals(
				"error",
				errors.NewSevere(
					"clientCertificate and clientCertificates cannot be used at the same time",
				),
				err,
			)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return error when clientCertificates entry has both content and ref", func() {
		fixtures := fixture.
			Builder().
			WithApplication(constants.Application).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("setting a clientCertificates entry with both content and ref")

		fixtures.Application.Spec.Settings.TLS = &application.TLSSettings{
			ClientCertificates: []application.ClientCertificate{
				{
					Name:    "cert1",
					Content: "some-content",
					Ref: &application.CertificateRef{
						Kind: "secrets",
						Name: "some-secret",
					},
				},
			},
		}

		By("checking that application validation returns error")

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, fixtures.Application)
			return assert.Equals(
				"error",
				errors.NewSevere(
					"clientCertificates[0]: content and ref cannot both be set",
				),
				err,
			)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return error when clientCertificates entry has neither content nor ref", func() {
		fixtures := fixture.
			Builder().
			WithApplication(constants.Application).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("setting a clientCertificates entry with no content and no ref")

		fixtures.Application.Spec.Settings.TLS = &application.TLSSettings{
			ClientCertificates: []application.ClientCertificate{
				{
					Name: "cert1",
				},
			},
		}

		By("checking that application validation returns error")

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, fixtures.Application)
			return assert.Equals(
				"error",
				errors.NewSevere(
					"clientCertificates[0]: either content or ref must be set",
				),
				err,
			)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return error when clientCertificates entry has invalid content", func() {
		fixtures := fixture.
			Builder().
			WithApplication(constants.Application).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("setting a clientCertificates entry with invalid content")

		fixtures.Application.Spec.Settings.TLS = &application.TLSSettings{
			ClientCertificates: []application.ClientCertificate{
				{
					Name:    "cert1",
					Content: "invalid-content",
				},
			},
		}

		By("checking that application dryrun validation")

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, fixtures.Application)
			return assert.Equals(
				"error",
				errors.NewSevere(
					"certificate is empty",
				),
				err,
			)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return error when clientCertificates entry has invalid inverted start and end dates ", func() {
		fixtures := fixture.
			Builder().
			WithApplication(constants.Application).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("setting a clientCertificates entry with end date before start date")

		fixtures.Application.Spec.Settings.TLS = &application.TLSSettings{
			ClientCertificates: []application.ClientCertificate{
				{
					Name:     "cert1",
					StartsAt: "2027-01-01T09:01:50Z",
					EndsAt:   "2020-01-01T09:01:50Z", // end before start
					Content: `
-----BEGIN CERTIFICATE-----
MIIC8DCCAdigAwIBAgIUZm+8jaHYuScdmU/PWGUpE44XkYcwDQYJKoZIhvcNAQEL
BQAwDjEMMAoGA1UEAwwDR0tPMB4XDTI2MDIyMDA5MDE1MVoXDTM2MDIxODA5MDE1
MVowEjEQMA4GA1UEAwwHY2xpZW50MTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
AQoCggEBAKsPx0rL8snrF7OrxDMi6t6GH6XaeyFWrvrZ4vO0DkmN+8aUZjsVK8+6
LAZ5jj5tGBGGc2jjC7/45nvjhFvT4pUSOTqJx4b6lJUV5MLXGihw+xYCAPF+hAcB
eyb4GZJXW1KBsMF51Fkx5v/W8pr/NkwkhRcI064DOslyMBcJepu2GjkecLnUQ5Ow
XBu52xp5ZSbasOpNR8uO/ZsNkviS1zFSo0J9ckXq+iJkgyx9D7QXIj4vpxfN/Wy+
tx02WlZjRdiEdQqubN2hv4c1/BHRJAqo1T1JzllY+KgCohN9uKVkgF5FV5/G0BSj
AXBovRTWpFGAsBpz1Mt2x0PSN2hGl4UCAwEAAaNCMEAwHQYDVR0OBBYEFP65h0BM
HvXgVQYSAE026kKhdOJWMB8GA1UdIwQYMBaAFNkFt/7UZZDMVZ55lcNDVDCR5rvo
MA0GCSqGSIb3DQEBCwUAA4IBAQCVXwc+BFOZU1sn6jhGQ/ADS1fRJ4KS1Tr5wKO5
aTrFIFo5JJ/akNlg92k8rOXMFactQfMWpb57MIeEgBEtmQ5cAbWFqyyuMlDZEyOf
OKQ1hfLII5je1p9Wy+ec2SoVIkbxyBIer4hAOXqnoMKr2lp/u6JoknJZ5DcvRgTb
BOwMRlejW7NI7DAgpgX+vGeNp3CFoeFCwzuLfRYOv35f7KDDTLn9UjoVNWznQcF9
8EyG7naoQcvzo4zcXCXb/l+rHUKOC3ecGSlJXbeGp1jgvpYsNScj7tRkkG0pp8Ig
H35MbSlE4AQ7desE0rDTg0KytcXVNZ2pTWgEcs43R/Vim4jj
-----END CERTIFICATE-----
`,
				},
			},
		}

		By("checking that application dry run validation returns error")

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, fixtures.Application)
			return assert.Equals(
				"error",
				errors.NewSeveref(
					"certificate [%s]: startsAt must be before endsAt", "cert1",
				),
				err,
			)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	// TODO test default (kind and name)
})
