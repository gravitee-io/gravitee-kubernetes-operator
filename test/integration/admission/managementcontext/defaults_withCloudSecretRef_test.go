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

package managementcontext

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/mctx"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
)

var _ = Describe("Default create", labels.WithContext, func() {
	interval := constants.Interval
	timeout := constants.EventualTimeout
	ctx := context.Background()
	admissionCtrl := mctx.AdmissionCtrl{}

	It("should configure context from secret", func() {

		By("setting a secret and a context")

		fixtures := fixture.Builder().
			AddSecret(constants.ContextCloudTokenSecretFile).
			WithContext(constants.ContextCloudWithSecretRefFile).
			Build().
			Apply()

		secret := &v1.Secret{}
		Eventually(func() error {
			return manager.Client().Get(ctx, types.NamespacedName{
				Name:      fixtures.Secrets[0].Name,
				Namespace: fixtures.Context.Namespace,
			}, secret)
		}, timeout, interval).Should(Succeed())

		By("defaulting the context")

		Consistently(func() error {
			err := admissionCtrl.Default(ctx, fixtures.Context)
			// check the token was read
			if err := assert.Equals("baseUrl",
				"https://eu.cloudgate.gravitee.io/apim/rest",
				fixtures.Context.Spec.BaseUrl); err != nil {
				return err
			}
			if err := assert.Nil("error", err); err != nil {
				return err
			}
			_, err = admissionCtrl.ValidateCreate(ctx, fixtures.Context)
			if err := assert.Nil("error", err); err != nil {
				return err
			}
			return nil
		}, constants.ConsistentTimeout, interval).Should(Succeed())
	})

})
