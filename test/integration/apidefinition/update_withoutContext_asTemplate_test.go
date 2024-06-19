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

package apidefinition

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	xhttp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/endpoint"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Delete", labels.WithoutContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	httpCli := xhttp.NewNoAuthClient(ctx)

	It("should update template and ingress", func() {
		fixtures := fixture.Builder().
			WithAPI(constants.ApiWithTemplateAnnotation).
			WithIngress(constants.IngressWithTemplateFile).
			Build().
			Apply()

		By("expecting to find an API definition created for ingress")

		apiDef := new(v1alpha1.ApiDefinition)

		Eventually(func() error {
			return manager.Client().Get(ctx, types.NamespacedName{
				Namespace: fixtures.Ingress.Namespace,
				Name:      fixtures.Ingress.Name,
			}, apiDef)
		}, timeout, interval).Should(Succeed())

		By("checking access to backend service, expecting status 401")

		url := endpoint.ForV2(apiDef)

		Eventually(func() error {
			err := httpCli.Get(url.WithPath("hostname"), new(Host), xhttp.WithHost("httpbin.example.com"))
			if !errors.IsUnauthorized(err) {
				return assert.Equals("error", "[UNAUTHORIZED]", err)
			}
			return nil
		}, timeout, interval).Should(Succeed())

		By("updating the template")

		updated := fixtures.API.DeepCopy()
		updated.Spec.Plans = append(updated.Spec.Plans,
			v2.NewPlan(
				base.NewPlan("Default keyless plan").
					WithStatus(base.PublishedPlanStatus),
			).WithSecurity("KEY_LESS").WithName("Key Less"),
		)

		Eventually(func() error {
			return manager.UpdateSafely(ctx, updated)
		}, timeout, interval).Should(Succeed())

		By("checking access to backend service, expecting status 200")

		Eventually(func() error {
			return httpCli.Get(url.WithPath("hostname"), new(Host), xhttp.WithHost("httpbin.example.com"))
		}, timeout, interval).Should(Succeed())
	})
})
