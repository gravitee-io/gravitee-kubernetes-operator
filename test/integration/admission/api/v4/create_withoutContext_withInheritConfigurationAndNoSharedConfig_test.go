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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"

	apiV4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate create", labels.WithoutContext, func() {
	ctx := context.Background()
	admissionCtrl := v4.AdmissionCtrl{}

	It("should return error when an endpoint inherits configuration from an EndpointGroup "+
		"with no sharedConfiguration", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4).
			Build()

		group := apiV4.NewHttpEndpointGroup("default-group")
		endpoint := apiV4.NewHttpEndpoint("default")
		endpoint.Inherit = true
		group.Endpoints = []*apiV4.Endpoint{endpoint}

		By("setting an EndpointGroup with no sharedConfiguration")

		fixtures.APIv4.Spec.EndpointGroups = []*apiV4.EndpointGroup{group}

		By("checking that api does not pass validation")

		_, err := admissionCtrl.ValidateCreate(ctx, fixtures.APIv4)

		Expect(err).To(Equal(
			errors.NewSeveref(
				"Endpoint [%s] in EndpointGroup [%s] of API [%s] has inheritConfiguration set to true "+
					"but the EndpointGroup does not define a sharedConfiguration to inherit from",
				"default", "default-group", fixtures.APIv4.Name,
			),
		))
	})

	It("should pass validation when the EndpointGroup defines a sharedConfiguration", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4).
			Build()

		group := apiV4.NewHttpEndpointGroup("default-group")
		group.SharedConfig = utils.ToGenericStringMap(map[string]interface{}{
			"target": "https://api.gravitee.io/echo",
		})
		endpoint := apiV4.NewHttpEndpoint("default")
		endpoint.Inherit = true
		group.Endpoints = []*apiV4.Endpoint{endpoint}

		By("setting an EndpointGroup with a sharedConfiguration")

		fixtures.APIv4.Spec.EndpointGroups = []*apiV4.EndpointGroup{group}

		By("checking that api passes validation")

		_, err := admissionCtrl.ValidateCreate(ctx, fixtures.APIv4)

		Expect(err).ToNot(HaveOccurred())
	})
})
