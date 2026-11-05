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

package dictionary

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/dictionary"
	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/dictionary"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate create", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := adm.AdmissionCtrl{}

	It("should return severe error when MANUAL type has dynamic field", func() {
		fixtures := fixture.
			Builder().
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		dict := fixture.
			Builder().
			WithDictionary(constants.DictionaryManualFile).
			Build()

		dict.Dictionary.Spec.Context = fixtures.Context.GetNamespacedName()
		dict.Dictionary.Spec.Dynamic = &dictionary.DynamicSpec{
			Provider: &dictionary.Provider{
				ProviderType:  "HTTP",
				URL:           "https://example.com",
				Method:        "GET",
				Specification: "[]",
			},
			Trigger: &dictionary.Trigger{
				Rate: 30,
				Unit: dictionary.SecondsUnit,
			},
		}

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, dict.Dictionary)
			return assert.Equals(
				"error",
				errors.NewSevere("dictionary type is MANUAL but 'dynamic' field is set, use 'manual' instead"),
				err,
			)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return severe error when DYNAMIC type has manual field", func() {
		fixtures := fixture.
			Builder().
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		dict := fixture.
			Builder().
			WithDictionary(constants.DictionaryDynamicFile).
			Build()

		dict.Dictionary.Spec.Context = fixtures.Context.GetNamespacedName()
		dict.Dictionary.Spec.Manual = &dictionary.ManualSpec{
			Properties: map[string]string{"key": "value"},
		}

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, dict.Dictionary)
			return assert.Equals(
				"error",
				errors.NewSevere("dictionary type is DYNAMIC but 'manual' field is set, use 'dynamic' instead"),
				err,
			)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return severe error when MANUAL type has no manual field", func() {
		fixtures := fixture.
			Builder().
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		dict := fixture.
			Builder().
			WithDictionary(constants.DictionaryManualFile).
			Build()

		dict.Dictionary.Spec.Context = fixtures.Context.GetNamespacedName()
		dict.Dictionary.Spec.Manual = nil

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, dict.Dictionary)
			return assert.Equals(
				"error",
				errors.NewSevere("dictionary type is MANUAL but 'manual' field is not set"),
				err,
			)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return severe error when DYNAMIC type has no dynamic field", func() {
		fixtures := fixture.
			Builder().
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		dict := fixture.
			Builder().
			WithDictionary(constants.DictionaryDynamicFile).
			Build()

		dict.Dictionary.Spec.Context = fixtures.Context.GetNamespacedName()
		dict.Dictionary.Spec.Dynamic = nil

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, dict.Dictionary)
			return assert.Equals(
				"error",
				errors.NewSevere("dictionary type is DYNAMIC but 'dynamic' field is not set"),
				err,
			)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
