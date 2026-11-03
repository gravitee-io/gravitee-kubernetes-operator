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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"
)

// The CRD rules (schema and CEL) are evaluated by the API server, not by the webhook handler.
var _ = Describe("CRD validation", func() {
	ctx := context.Background()

	create := func() *v1alpha1.AMClientSecret {
		cs := &v1alpha1.AMClientSecret{
			ObjectMeta: metaV1.ObjectMeta{Name: random.GetName(), Namespace: constants.Namespace},
		}
		cs.Spec.ApplicationRef = refs.NamespacedName{Name: "some-application"}
		cs.Spec.SecretRef = am.SecretKeyRef{Name: "some-secret", Key: "client-secret"}
		Expect(manager.Client().Create(ctx, cs)).To(Succeed())
		DeferCleanup(func() {
			Expect(client.IgnoreNotFound(manager.Client().Delete(ctx, cs))).To(Succeed())
		})
		return cs
	}

	DescribeTable("should reject a change of an immutable reference (AC4)",
		func(change func(*v1alpha1.AMClientSecret), message string) {
			cs := create()
			Eventually(ctx, func() error {
				if err := manager.GetLatest(ctx, cs); err != nil {
					return err
				}
				change(cs)
				return manager.Client().Update(ctx, cs)
			}, constants.EventualTimeout, constants.Interval).Should(MatchError(ContainSubstring(message)))
		},
		Entry("applicationRef", func(cs *v1alpha1.AMClientSecret) {
			cs.Spec.ApplicationRef.Name = "other-application"
		}, "applicationRef is immutable"),
		Entry("secretRef.name", func(cs *v1alpha1.AMClientSecret) {
			cs.Spec.SecretRef.Name = "other-secret"
		}, "secretRef is immutable"),
		Entry("secretRef.key", func(cs *v1alpha1.AMClientSecret) {
			cs.Spec.SecretRef.Key = "other-key"
		}, "secretRef is immutable"),
	)

	It("should require a non-empty applicationRef name", func() {
		cs := &v1alpha1.AMClientSecret{
			ObjectMeta: metaV1.ObjectMeta{Name: random.GetName(), Namespace: constants.Namespace},
		}
		cs.Spec.SecretRef = am.SecretKeyRef{Name: "some-secret", Key: "client-secret"}
		Expect(manager.Client().Create(ctx, cs)).To(MatchError(ContainSubstring("applicationRef.name must not be empty")))
	})
})
