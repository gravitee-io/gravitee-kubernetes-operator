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

package amcontext

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/amctx"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate AMContext", func() {
	ctx := context.Background()
	admissionCtrl := amctx.AdmissionCtrl{}

	build := func(file string) *v1alpha1.AMContext {
		return fixture.Builder().WithAMContext(file).Build().AMContext
	}

	withInlineToken := func(obj *v1alpha1.AMContext, token string) *v1alpha1.AMContext {
		obj.Spec.Auth.SecretRef = nil
		obj.Spec.Auth.BearerToken = token
		return obj
	}

	createSecret := func(data map[string][]byte) string {
		secret := &coreV1.Secret{
			ObjectMeta: metaV1.ObjectMeta{Name: random.GetName(), Namespace: constants.Namespace},
			Data:       data,
		}
		Expect(manager.Client().Create(ctx, secret)).To(Succeed())
		DeferCleanup(func() {
			Expect(client.IgnoreNotFound(manager.Client().Delete(ctx, secret))).To(Succeed())
		})
		return secret.Name
	}

	It("should accept a valid secretRef", func() {
		warnings, err := admissionCtrl.ValidateCreate(ctx, build(constants.AMContextFile))
		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(BeEmpty())
	})

	It("should accept a valid inline token", func() {
		obj := withInlineToken(build(constants.AMContextFile), "admin-token")
		_, err := admissionCtrl.ValidateCreate(ctx, obj)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should reject a missing Secret", func() {
		obj := build(constants.AMContextFile)
		name := random.GetName()
		obj.Spec.Auth.SecretRef = &refs.NamespacedName{Name: name}

		_, err := admissionCtrl.ValidateCreate(ctx, obj)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("doesn't exist"))
		Expect(err.Error()).To(ContainSubstring(name))
	})

	It("should reject a Secret without bearerToken key", func() {
		obj := build(constants.AMContextFile)
		obj.Spec.Auth.SecretRef = &refs.NamespacedName{Name: createSecret(map[string][]byte{})}

		_, err := admissionCtrl.ValidateCreate(ctx, obj)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed with status 401"))
	})

	It("should reject a wrong inline token", func() {
		_, err := admissionCtrl.ValidateCreate(ctx, build(constants.AMContextBadTokenFile))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed with status 401"))
	})

	It("should admit with a warning when AM is unreachable", func() {
		warnings, err := admissionCtrl.ValidateCreate(ctx, build(constants.AMContextUnreachableFile))
		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(ContainElement(And(
			ContainSubstring("unable to reach AM"),
			ContainSubstring("http://localhost:1"),
		)))
	})

	It("should accept a templated token", func() {
		name := createSecret(map[string][]byte{"token": []byte("admin-token")})
		obj := withInlineToken(build(constants.AMContextFile), "[[ secret `"+name+"/token` ]]")

		_, err := admissionCtrl.ValidateCreate(ctx, obj)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should reject an update to a wrong token", func() {
		oldObj := build(constants.AMContextFile)
		newObj := withInlineToken(oldObj.DeepCopy(), "invalid-token")

		_, err := admissionCtrl.ValidateUpdate(ctx, oldObj, newObj)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed with status 401"))
	})
})
