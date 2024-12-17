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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/mctx"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Validate delete", labels.WithContext, func() {
	ctx := context.Background()
	timeout := constants.EventualTimeout
	interval := constants.Interval

	admissionCtrl := adm.AdmissionCtrl{}

	fixtures := fixture.Builder().
		WithContext(constants.ContextWithCredentialsFile).
		WithApplication(constants.Application).
		WithAPIv4(constants.ApiV4WithContextFile).
		WithAPI(constants.ApiWithContextFile).
		Build().
		Apply()

	DescribeTable(
		"should fail",
		func(mctx *v1alpha1.ManagementContext, refObj client.Object, expected error) {

			Eventually(func() error {
				_, err := admissionCtrl.ValidateDelete(ctx, mctx)
				return assert.Equals("error", expected, err)
			}, timeout, interval).Should(Succeed(), mctx.GetName())

			Expect(manager.Delete(ctx, refObj)).To(Succeed())
		},
		Entry(
			"with API ref",
			fixtures.Context,
			fixtures.API,
			errors.NewSeveref(
				"[%s] cannot be deleted because %d APIs are relying on this context. "+
					"You can and review this APIs using the following command: "+
					"kubectl get apidefinitions.gravitee.io -A "+
					"-o jsonpath='{.items[?(@.spec.contextRef.name==\"%s\")].metadata.name}'",
				fixtures.Context.Name, 1, fixtures.Context.Name,
			),
		),
		Entry(
			"with APIv4 ref",
			fixtures.Context,
			fixtures.APIv4,
			errors.NewSeveref(
				"[%s] cannot be deleted because %d APIs are relying on this context. "+
					"You can and review this APIs using the following command: "+
					"kubectl get apiv4definitions.gravitee.io -A "+
					"-o jsonpath='{.items[?(@.spec.contextRef.name==\"%s\")].metadata.name}'",
				fixtures.Context.Name, 1, fixtures.Context.Name,
			),
		),
		Entry(
			"with application ref",
			fixtures.Context,
			fixtures.Application,
			errors.NewSeveref(
				"[%s] cannot be deleted because %d applications are relying on this context. "+
					"You can and review this applications using the following command: "+
					"kubectl get applications.gravitee.io -A "+
					"-o jsonpath='{.items[?(@.spec.contextRef.name==\"%s\")].metadata.name}'",
				fixtures.Context.Name, 1, fixtures.Context.Name,
			),
		),
	)
})
