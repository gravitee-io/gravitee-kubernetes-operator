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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/mctx"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Validate update", labels.WithContext, func() {
	interval := constants.Interval
	admissionCtrl := mctx.AdmissionCtrl{}
	ctx := context.Background()
	cli := manager.Client()

	It("should return error if secret is missing", func() {
		fixtures := fixture.Builder().
			WithContext(constants.ContextWithSecretFile).
			Build()

		fixtures.Context.Spec.SecretRef().Name = "unknown-secret"
		fixtures.Apply()

		mCtx := new(v1alpha1.ManagementContext)
		Eventually(func() error {
			return cli.Get(context.Background(), types.NamespacedName{
				Namespace: fixtures.Context.Namespace,
				Name:      fixtures.Context.Name,
			}, mCtx)
		}, constants.EventualTimeout, interval).Should(Succeed())

		Consistently(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, mCtx)
			return err
		}, constants.ConsistentTimeout, interval).ShouldNot(Succeed())
	})
})
