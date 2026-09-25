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
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/am/domain"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	internal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/securitydomain"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newAMSecurityDomain(name, path string) *v1alpha1.AMSecurityDomain {
	enabled := true
	return &v1alpha1.AMSecurityDomain{
		ObjectMeta: metav1.ObjectMeta{Name: "test-sd"},
		Spec: v1alpha1.AMSecurityDomainSpec{
			Domain: domain.Domain{
				Name:    name,
				Path:    path,
				Enabled: &enabled,
			},
			Context: &refs.NamespacedName{Namespace: "gravitee", Name: "am-ctx"},
		},
	}
}

var _ = Describe("AMSecurityDomain DTO mapping", func() {
	It("maps CRD fields to the SDK domain type", func() {
		obj := newAMSecurityDomain("my-domain", "/my-domain")
		dto := internal.ToDomainDTO(obj)

		Expect(dto.Name).To(Equal("my-domain"))
		Expect(dto.Path).To(Equal("/my-domain"))
		Expect(dto.Enabled).ToNot(BeNil())
		Expect(*dto.Enabled).To(BeTrue())
	})

	It("maps optional fields as nil when absent", func() {
		obj := newAMSecurityDomain("minimal", "/minimal")
		dto := internal.ToDomainDTO(obj)

		Expect(dto.AccountSettings).To(BeNil())
		Expect(dto.CorsSettings).To(BeNil())
		Expect(dto.LoginSettings).To(BeNil())
		Expect(dto.PasswordSettings).To(BeNil())
		Expect(dto.WebAuthnSettings).To(BeNil())
		Expect(dto.Tags).To(BeNil())
	})

	It("maps tags when present", func() {
		obj := newAMSecurityDomain("tagged", "/tagged")
		tags := []string{"eu", "production"}
		obj.Spec.Tags = tags
		dto := internal.ToDomainDTO(obj)

		Expect(dto.Tags).ToNot(BeNil())
		Expect(dto.Tags).To(Equal([]string{"eu", "production"}))
	})

	It("keys the DTO on the resource HRID, not the contextRef", func() {
		obj := newAMSecurityDomain("d1", "/d1")
		dto := internal.ToDomainDTO(obj)

		Expect(dto.Key).To(Equal(refs.NewNamespacedNameFromObject(obj).HRID()))
		Expect(dto.Key).ToNot(ContainSubstring("am-ctx"))
	})
})

var _ = Describe("AMSecurityDomain UpdateStatus", func() {
	It("populates status fields from DomainResponse", func() {
		obj := newAMSecurityDomain("d1", "/d1")
		resp := internal.DomainResponse{
			Key:    "domain-key-123",
			OrgEnv: am.OrgEnv{OrgID: "org-1", EnvID: "env-1"},
		}

		err := internal.UpdateStatus(context.Background(), obj, resp)
		Expect(err).ToNot(HaveOccurred())
		Expect(obj.Status.Status.ID).To(Equal("domain-key-123"))
		Expect(obj.Status.Status.OrgID).To(Equal("org-1"))
		Expect(obj.Status.Status.EnvID).To(Equal("env-1"))
	})
})

var _ = Describe("AMSecurityDomain type behavior", func() {
	It("reports not being deleted when DeletionTimestamp is zero", func() {
		obj := newAMSecurityDomain("d1", "/d1")
		Expect(obj.IsBeingDeleted()).To(BeFalse())
	})

	It("reports being deleted when DeletionTimestamp is set", func() {
		obj := newAMSecurityDomain("d1", "/d1")
		now := metav1.Now()
		obj.DeletionTimestamp = &now
		Expect(obj.IsBeingDeleted()).To(BeTrue())
	})

	It("returns the context ref", func() {
		obj := newAMSecurityDomain("d1", "/d1")
		Expect(obj.ContextRef().GetName()).To(Equal("am-ctx"))
		Expect(obj.ContextRef().GetNamespace()).To(Equal("gravitee"))
	})

	It("reports HasContext when contextRef is set", func() {
		obj := newAMSecurityDomain("d1", "/d1")
		Expect(obj.HasContext()).To(BeTrue())
	})

	It("reports no context when contextRef is nil", func() {
		obj := newAMSecurityDomain("d1", "/d1")
		obj.Spec.Context = nil
		Expect(obj.HasContext()).To(BeFalse())
	})

	It("IsFailed returns false with no conditions", func() {
		obj := newAMSecurityDomain("d1", "/d1")
		Expect(obj.Status.IsFailed()).To(BeFalse())
	})

	It("IsFailed returns true when a condition is false", func() {
		obj := newAMSecurityDomain("d1", "/d1")
		obj.Status.Conditions = []metav1.Condition{
			{Type: "Accepted", Status: metav1.ConditionFalse, Reason: "Error"},
		}
		Expect(obj.Status.IsFailed()).To(BeTrue())
	})

	It("IsFailed returns false when all conditions are true", func() {
		obj := newAMSecurityDomain("d1", "/d1")
		obj.Status.Conditions = []metav1.Condition{
			{Type: "Accepted", Status: metav1.ConditionTrue, Reason: "OK"},
		}
		Expect(obj.Status.IsFailed()).To(BeFalse())
	})

	It("DeepCopyFrom copies status from another AMSecurityDomain", func() {
		src := newAMSecurityDomain("d1", "/d1")
		src.Status.Status.ID = "key-abc"
		src.Status.Status.OrgID = "org-1"

		dst := &v1alpha1.AMSecurityDomainStatus{}
		err := dst.DeepCopyFrom(src)
		Expect(err).ToNot(HaveOccurred())
		Expect(dst.Status.ID).To(Equal("key-abc"))
		Expect(dst.Status.OrgID).To(Equal("org-1"))
	})

	It("DeepCopyTo copies status to another AMSecurityDomain", func() {
		src := &v1alpha1.AMSecurityDomainStatus{}
		src.Status.ID = "key-xyz"
		src.Status.EnvID = "env-2"

		dst := newAMSecurityDomain("d1", "/d1")
		err := src.DeepCopyTo(dst)
		Expect(err).ToNot(HaveOccurred())
		Expect(dst.Status.Status.ID).To(Equal("key-xyz"))
		Expect(dst.Status.Status.EnvID).To(Equal("env-2"))
	})
})

var _ = Describe("AMSecurityDomain key validation", func() {
	withName := func(namespace, name string) *v1alpha1.AMSecurityDomain {
		obj := newAMSecurityDomain("d1", "/d1")
		obj.Namespace = namespace
		obj.Name = name
		return obj
	}

	DescribeTable("accepts names that give a valid AM key",
		func(namespace, name string) {
			Expect(internal.ValidateKey(context.Background(), withName(namespace, name)).IsSevere()).To(BeFalse())
		},
		Entry("simple name", "default", "my-domain"),
		Entry("digits", "team1", "domain2"),
		Entry("key of exactly 255 characters", "ns", strings.Repeat("a", 252)),
	)

	DescribeTable("rejects names that give an invalid AM key",
		func(namespace, name string) {
			Expect(internal.ValidateKey(context.Background(), withName(namespace, name)).IsSevere()).To(BeTrue())
		},
		Entry("dot in name", "default", "my.domain"),
		Entry("key longer than 255 characters", "ns", strings.Repeat("a", 253)),
	)
})
