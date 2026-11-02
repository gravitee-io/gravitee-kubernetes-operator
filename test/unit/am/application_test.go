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

package am_test

import (
	"context"
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/application"
)

// newApplication decodes spec into an application of the ns namespace.
func newApplication(spec string) *v1alpha1.AMApplication {
	app := &v1alpha1.AMApplication{ObjectMeta: metav1.ObjectMeta{Name: "mobile", Namespace: "ns"}}
	Expect(json.Unmarshal([]byte(spec), &app.Spec)).To(Succeed())
	return app
}

const minimalApplication = `{
	"domainRef": {"name": "domain"},
	"name": "Mobile",
	"type": "native",
	"certificateRef": {"name": "signing"}
}`

// sentApplication is the application payload as sent to AM.
func sentApplication(app *v1alpha1.AMApplication) map[string]any {
	dto, err := application.ToApplicationDTO(app)
	Expect(err).ToNot(HaveOccurred())
	payload, err := json.Marshal(dto.Application)
	Expect(err).ToNot(HaveOccurred())
	var sent map[string]any
	Expect(json.Unmarshal(payload, &sent)).To(Succeed())
	return sent
}

var _ = Describe("AMApplication DTO mapping", func() {
	It("keys the application and its domain by their HRIDs", func() {
		app := newApplication(minimalApplication)
		app.Spec.DomainRef.Namespace = "ignored"

		dto, err := application.ToApplicationDTO(app)
		Expect(err).ToNot(HaveOccurred())

		Expect(dto.Key).To(Equal("ns-mobile"))
		Expect(dto.DomainKey).To(Equal("ns-domain"))
	})

	It("defaults enabled to true and template to false, and leaves clientId out when not set", func() {
		sent := sentApplication(newApplication(minimalApplication))

		Expect(sent["enabled"]).To(BeTrue())
		Expect(sent["template"]).To(BeFalse())
		Expect(sent).ToNot(HaveKey("clientId"))
	})

	It("sends metadata and jwks as objects and thresholds as numbers", func() {
		sent := sentApplication(newApplication(`{
			"domainRef": {"name": "domain"},
			"name": "Mobile",
			"type": "native",
			"certificate": "signing",
			"metadata": {"owner": "payments", "tier": 1},
			"settings": {
				"oauth": {"jwks": {"keys": [{"kty": "RSA", "kid": "k1"}]}},
				"riskAssessment": {"deviceAssessment": {"enabled": true, "thresholds": {"HIGH": 0.8}}}
			}
		}`))

		Expect(sent["metadata"]).To(Equal(map[string]any{"owner": "payments", "tier": float64(1)}))
		settings := sent["settings"].(map[string]any)
		oauth := settings["oauth"].(map[string]any)
		Expect(oauth["jwks"]).To(Equal(map[string]any{
			"keys": []any{map[string]any{"kty": "RSA", "kid": "k1"}},
		}))
		risk := settings["riskAssessment"].(map[string]any)
		device := risk["deviceAssessment"].(map[string]any)
		Expect(device["thresholds"]).To(Equal(map[string]any{"HIGH": 0.8}))
	})

	It("sends a certificateRef as the certificate's key", func() {
		sent := sentApplication(newApplication(minimalApplication))

		Expect(sent["certificate"]).To(Equal("ns-signing"))
		Expect(sent).ToNot(HaveKey("certificateRef"))
	})

	It("keeps the order of mixed identity and identityRef items", func() {
		sent := sentApplication(newApplication(`{
			"domainRef": {"name": "domain"},
			"name": "Mobile",
			"type": "native",
			"certificate": "signing",
			"identityProviders": [
				{"identity": "ldap", "priority": 0},
				{"identityRef": {"name": "users"}, "priority": 1},
				{"identity": "id:6f1c3e2a-9b7d-4c1e-8f0a-2d5b7e9c4a11", "priority": 2}
			]
		}`))

		Expect(sent["identityProviders"]).To(Equal([]any{
			map[string]any{"identity": "ldap", "priority": float64(0)},
			map[string]any{"identity": "ns-users", "priority": float64(1)},
			map[string]any{"identity": "id:6f1c3e2a-9b7d-4c1e-8f0a-2d5b7e9c4a11", "priority": float64(2)},
		}))
	})

	It("sends the same application for a key and for a ref giving that key", func() {
		byRef := newApplication(`{
			"domainRef": {"name": "domain"},
			"name": "Mobile",
			"type": "native",
			"certificateRef": {"name": "signing"},
			"identityProviders": [{"identityRef": {"name": "users"}}]
		}`)
		byKey := newApplication(`{
			"domainRef": {"name": "domain"},
			"name": "Mobile",
			"type": "native",
			"certificate": "ns-signing",
			"identityProviders": [{"identity": "ns-users"}]
		}`)

		Expect(sentApplication(byRef)).To(Equal(sentApplication(byKey)))
	})
})

var _ = Describe("AMApplication reference keys", func() {
	It("gives the certificate key from certificate or certificateRef", func() {
		Expect(newApplication(minimalApplication).CertificateKey()).To(Equal("ns-signing"))

		app := newApplication(minimalApplication)
		app.Spec.CertificateRef, app.Spec.Certificate = nil, new("legacy-signing")
		Expect(app.CertificateKey()).To(Equal("legacy-signing"))
	})

	It("gives no certificate key for an id", func() {
		app := newApplication(minimalApplication)
		app.Spec.CertificateRef, app.Spec.Certificate = nil, new("id:6f1c3e2a-9b7d-4c1e-8f0a-2d5b7e9c4a11")

		Expect(app.CertificateKey()).To(BeEmpty())
	})

	It("gives the identity keys from identity or identityRef, leaving ids out", func() {
		app := newApplication(`{
			"domainRef": {"name": "domain"},
			"name": "Mobile",
			"type": "native",
			"certificate": "signing",
			"identityProviders": [
				{"identity": "ldap"},
				{"identityRef": {"name": "users"}},
				{"identity": "id:6f1c3e2a-9b7d-4c1e-8f0a-2d5b7e9c4a11"}
			]
		}`)

		Expect(app.IdentityKeys()).To(Equal([]string{"ldap", "ns-users"}))
	})
})

var _ = Describe("AMApplication UpdateStatus", func() {
	It("copies the keys and the clientId AM stored", func() {
		app := &v1alpha1.AMApplication{}
		resp := application.Response{
			DomainSubResourceResponse: am.DomainSubResourceResponse{
				DomainKey: "ns-domain",
				BaseResponse: am.BaseResponse{
					Key:    "ns-mobile",
					OrgEnv: am.OrgEnv{OrgID: "DEFAULT", EnvID: "DEFAULT"},
				},
			},
			ClientID: "generated-client-id",
		}

		Expect(application.UpdateStatus(context.Background(), app, resp)).To(Succeed())

		Expect(app.Status.Key).To(Equal("ns-mobile"))
		Expect(app.Status.DomainKey).To(Equal("ns-domain"))
		Expect(app.Status.OrgID).To(Equal("DEFAULT"))
		Expect(app.Status.EnvID).To(Equal("DEFAULT"))
		Expect(app.Status.ClientID).To(Equal("generated-client-id"))
	})
})
