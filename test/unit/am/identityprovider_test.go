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

package am_test

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	internal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/identityprovider"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

var _ = Describe("AMIdentityProvider ResolveDomain", func() {
	const ns = "ns"

	domain := func(key string) *v1alpha1.AMSecurityDomain {
		d := &v1alpha1.AMSecurityDomain{ObjectMeta: metav1.ObjectMeta{Name: "domain", Namespace: ns}}
		d.Status.Key = key
		return d
	}
	idp := func() *v1alpha1.AMIdentityProvider {
		p := &v1alpha1.AMIdentityProvider{ObjectMeta: metav1.ObjectMeta{Name: "idp", Namespace: ns}}
		p.Spec.DomainRef = refs.NamespacedName{Name: "domain"}
		return p
	}
	deleting := func(p *v1alpha1.AMIdentityProvider) *v1alpha1.AMIdentityProvider {
		p.DeletionTimestamp = &metav1.Time{Time: time.Now()}
		return p
	}
	cluster := func(objects ...client.Object) {
		scheme := runtime.NewScheme()
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
		k8s.RegisterClient(fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).Build())
	}

	It("resolves a domain AM has created", func() {
		cluster(domain("ns-domain"))
		Expect(am.ResolveSubResourceDomain(context.Background(), idp(), ns)).To(Succeed())
	})

	It("fails with ErrDomainNotReady while AM has not created the domain", func() {
		cluster(domain(""))
		err := am.ResolveSubResourceDomain(context.Background(), idp(), ns)
		Expect(errors.Is(err, am.ErrDomainNotReady)).To(BeTrue(), "got %v", err)
	})

	It("returns NotFound when the domain does not exist", func() {
		cluster()
		err := am.ResolveSubResourceDomain(context.Background(), idp(), ns)
		Expect(apierrors.IsNotFound(err)).To(BeTrue(), "got %v", err)
	})

	It("does not require a ready domain on delete", func() {
		cluster(domain(""))
		Expect(am.ResolveSubResourceDomain(context.Background(), deleting(idp()), ns)).To(Succeed())
	})

	It("returns NotFound on delete when the domain is gone", func() {
		cluster()
		err := am.ResolveSubResourceDomain(context.Background(), deleting(idp()), ns)
		Expect(apierrors.IsNotFound(err)).To(BeTrue(), "got %v", err)
	})
})

// idpWithReadyDomain registers a fake cluster holding a domain AM has created, and returns an identity
// provider under it.
func idpWithReadyDomain() *v1alpha1.AMIdentityProvider {
	d := &v1alpha1.AMSecurityDomain{ObjectMeta: metav1.ObjectMeta{Name: "domain", Namespace: "ns"}}
	d.Status.Key = "ns-domain"
	scheme := runtime.NewScheme()
	Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
	k8s.RegisterClient(fake.NewClientBuilder().WithScheme(scheme).WithObjects(d).Build())

	p := &v1alpha1.AMIdentityProvider{ObjectMeta: metav1.ObjectMeta{Name: "idp", Namespace: "ns"}}
	p.Spec.DomainRef = refs.NamespacedName{Name: "domain"}
	p.Spec.Name = new("Inline users")
	p.Spec.Type = new("inline-am-idp")
	p.Spec.Configuration = utils.NewGenericStringMap().Put("users", []any{map[string]any{"username": "jdoe"}})
	return p
}

var _ = Describe("AMIdentityProvider DTO mapping", func() {
	It("keys the identity provider and its domain by their HRIDs", func() {
		idp := idpWithReadyDomain()
		idp.Spec.DomainRef.Namespace = "ignored"

		dto, err := internal.ToIdentityProviderDTO(idp)
		Expect(err).ToNot(HaveOccurred())

		Expect(dto.Key).To(Equal("ns-idp"))
		Expect(dto.DomainKey).To(Equal("ns-domain"))
	})

	It("sends configuration as a JSON string", func() {
		dto, err := internal.ToIdentityProviderDTO(idpWithReadyDomain())
		Expect(err).ToNot(HaveOccurred())

		payload, err := json.Marshal(dto.IdentityProvider)
		Expect(err).ToNot(HaveOccurred())
		var sent map[string]any
		Expect(json.Unmarshal(payload, &sent)).To(Succeed())
		Expect(sent["configuration"]).To(BeAssignableToTypeOf(""))
		Expect(sent["configuration"]).To(MatchJSON(`{"users":[{"username":"jdoe"}]}`))
		Expect(sent["name"]).To(Equal("Inline users"))
		Expect(sent["type"]).To(Equal("inline-am-idp"))
	})

	It("defaults system to false", func() {
		dto, err := internal.ToIdentityProviderDTO(idpWithReadyDomain())
		Expect(err).ToNot(HaveOccurred())

		Expect(dto.System).ToNot(BeNil())
		Expect(*dto.System).To(BeFalse())
	})
})

var _ = Describe("AMIdentityProvider drift mapping", func() {
	It("blanks name, type and configuration of the system identity provider: AM supplies them", func() {
		idp := idpWithReadyDomain()
		idp.Spec.System = new(true)

		dto, err := internal.ToIdentityProviderDTOForDrift(idp)
		Expect(err).ToNot(HaveOccurred())

		Expect(dto.Name).To(BeNil())
		Expect(dto.Type).To(BeNil())
		Expect(dto.Configuration).To(BeNil())
		Expect(dto.Key).To(Equal("ns-idp"))
	})

	It("keeps name, type and configuration of a regular identity provider", func() {
		idp := idpWithReadyDomain()

		dto, err := internal.ToIdentityProviderDTOForDrift(idp)
		Expect(err).ToNot(HaveOccurred())
		expected, err := internal.ToIdentityProviderDTO(idp)
		Expect(err).ToNot(HaveOccurred())

		Expect(dto).To(Equal(expected))
	})
})

var _ = Describe("AMIdentityProvider UpdateStatus", func() {
	It("copies the keys, and the name and type AM stored", func() {
		idp := &v1alpha1.AMIdentityProvider{}
		resp := internal.Response{
			DomainSubResourceResponse: am.DomainSubResourceResponse{
				DomainKey: "ns-domain",
				BaseResponse: am.BaseResponse{
					Key:    "ns-idp",
					OrgEnv: am.OrgEnv{OrgID: "DEFAULT", EnvID: "DEFAULT"},
				},
			},
			Name: "Default Identity Provider",
			Type: "gravitee-am-idp",
		}

		Expect(internal.UpdateStatus(context.Background(), idp, resp)).To(Succeed())

		Expect(idp.Status.Key).To(Equal("ns-idp"))
		Expect(idp.Status.DomainKey).To(Equal("ns-domain"))
		Expect(idp.Status.OrgID).To(Equal("DEFAULT"))
		Expect(idp.Status.EnvID).To(Equal("DEFAULT"))
		Expect(idp.Status.Name).To(Equal("Default Identity Provider"))
		Expect(idp.Status.Type).To(Equal("gravitee-am-idp"))
	})
})

// warnings returns the warnings as the webhook reports them.
func warnings(errs *gerrors.AdmissionErrors) []string {
	w, _ := errs.Map()
	return w
}

var _ = Describe("AMIdentityProvider PreCheck", func() {
	ctx := context.Background()

	It("admits a valid identity provider without warnings", func() {
		errs := internal.PreCheck(ctx, idpWithReadyDomain())

		Expect(errs.IsSevere()).To(BeFalse())
		Expect(warnings(errs)).To(BeEmpty())
	})

	It("rejects a key AM would refuse", func() {
		idp := idpWithReadyDomain()
		idp.Name = "idp.with.dots"

		Expect(internal.PreCheck(ctx, idp).IsSevere()).To(BeTrue())
	})

	It("rejects a domainRef in another namespace", func() {
		idp := idpWithReadyDomain()
		idp.Spec.DomainRef.Namespace = "other"

		Expect(internal.PreCheck(ctx, idp).IsSevere()).To(BeTrue())
	})

	It("warns when the domain is missing", func() {
		idp := idpWithReadyDomain()
		idp.Spec.DomainRef.Name = "missing"

		errs := internal.PreCheck(ctx, idp)

		Expect(errs.IsSevere()).To(BeFalse())
		Expect(warnings(errs)).To(ContainElement(ContainSubstring("not found or not yet created in AM")))
	})

	It("warns that name, type and configuration are ignored for the system identity provider", func() {
		idp := idpWithReadyDomain()
		idp.Spec.System = new(true)
		idp.Spec.Type = nil

		errs := internal.PreCheck(ctx, idp)

		Expect(errs.IsSevere()).To(BeFalse())
		Expect(warnings(errs)).To(ConsistOf("'configuration', 'name' will be ignored when 'system' is 'true'."))
	})
})
