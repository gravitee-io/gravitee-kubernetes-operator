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

package application

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/application"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func appWith(mutate func(*application.Application)) *v1alpha1.Application {
	app := &v1alpha1.Application{}
	app.Spec.Name = "deprecated-fields"
	mutate(&app.Spec.Application)
	return app
}

func warningsOf(app *v1alpha1.Application) []string {
	warnings, err := admission.WarnDeprecatedFields(app).Map()
	Expect(err).ToNot(HaveOccurred())
	return warnings
}

var _ = Describe("Deprecated application fields", func() {
	It("warns that picture is ignored and points to pictureUrl", func() {
		app := appWith(func(spec *application.Application) {
			spec.Picture = new("data:image/png;base64,iVBORw0KGgo=")
		})

		Expect(warningsOf(app)).To(ConsistOf(
			"[spec.picture] is deprecated and ignored by APIM, it will be removed in 2027-12; use [spec.pictureUrl]",
		))
	})

	It("warns once per metadata entry that sets hidden", func() {
		app := appWith(func(spec *application.Application) {
			spec.Metadata = &[]application.Metadata{
				{Name: "team", Value: new("payments"), Hidden: new(true)},
				{Name: "tier", Value: new("gold")},
				{Name: "region", Value: new("eu"), Hidden: new(false)},
			}
		})

		Expect(warningsOf(app)).To(ConsistOf(
			"[spec.metadata[team].hidden] is deprecated and ignored by APIM, it will be removed in 2027-12",
			"[spec.metadata[region].hidden] is deprecated and ignored by APIM, it will be removed in 2027-12",
		))
	})

	It("does not warn when neither field is set", func() {
		app := appWith(func(spec *application.Application) {
			spec.PictureURL = new("https://example.com/app.png")
			spec.Metadata = &[]application.Metadata{{Name: "tier", Value: new("gold")}}
		})

		Expect(warningsOf(app)).To(BeEmpty())
	})
})
