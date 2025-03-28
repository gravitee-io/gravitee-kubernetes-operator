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

package v4

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate update", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := v4.AdmissionCtrl{}

	It("should return error on API creation with invalid api resource", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4WithOauth2AmResourceFile).
			Build().
			Apply()
		ctxFixtures := fixture.
			Builder().
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("preparing API for import")
		fixtures.APIv4.Spec.Context = &refs.NamespacedName{
			Name:      ctxFixtures.Context.Name,
			Namespace: ctxFixtures.Context.Namespace,
		}

		By("adding an invalid to the API")
		invalidConfig := &utils.GenericStringMap{
			Unstructured: struct{ Object map[string]interface{} }{Object: map[string]interface{}{
				"wrong_json": "[{\"object\":object}]",
			}},
		}

		fixtures.APIv4.Spec.Resources[0].Configuration = invalidConfig
		By("checking that API validation returns errors")

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.APIv4, fixtures.APIv4)

			return assert.Equals(
				"severe",
				errors.NewSeveref(
					"Resource [%s] configuration is not valid",
					*fixtures.APIv4.Spec.Resources[0].Name,
				).Error(),
				err.Error(),
			)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
