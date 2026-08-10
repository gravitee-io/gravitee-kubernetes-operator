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

package docs

import (
	"context"

	documentation "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/docs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

// The area enum is declared on the CRD, so an unknown value is refused by the
// API server on create rather than by the admission webhook. No management
// context is needed: the request never reaches a controller.
var _ = Describe("Create", labels.WithoutContext, func() {
	ctx := context.Background()

	It("should reject a documentation whose area is not a known value", func() {
		doc := fixture.Builder().
			WithDocumentation(constants.DocumentationPortalFile).
			Build().
			Documentation

		doc.Spec.Area = utils.ToReference(documentation.PageArea("INVALID"))

		err := manager.Client().Create(ctx, doc)

		Expect(apierrors.IsInvalid(err)).To(BeTrue(), "expected an invalid object error, got %v", err)
		Expect(err.Error()).To(ContainSubstring("spec.area"))
		Expect(err.Error()).To(ContainSubstring(string(documentation.TopNavbar)))
		Expect(err.Error()).To(ContainSubstring(string(documentation.Homepage)))
	})
})
