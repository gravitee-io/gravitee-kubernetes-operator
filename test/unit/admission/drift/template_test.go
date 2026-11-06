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

package drift

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/dictionary"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	admdrift "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const (
	templatedNamespace = "default"
	currentSecret      = "[[ secret `client/current` ]]"
	rotatedSecret      = "[[ secret `client/rotated` ]]"
	deletedSecret      = "[[ secret `deleted/value` ]]"
)

var _ = Describe("Drift detection on templated values", func() {
	var (
		originalDrift      env.DriftDetection
		originalTemplating bool
	)

	BeforeEach(func() {
		originalDrift = env.Config.DriftDetection
		originalTemplating = env.Config.EnableTemplating
		env.Config.DriftDetection.Enabled = true
		env.Config.DriftDetection.Policy = env.DriftPolicyDeny
		env.Config.EnableTemplating = true

		scheme := runtime.NewScheme()
		Expect(clientgoscheme.AddToScheme(scheme)).To(Succeed())
		k8s.RegisterClient(fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(&corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "client", Namespace: templatedNamespace},
				Data: map[string][]byte{
					"current": []byte("current-value"),
					"rotated": []byte("rotated-value"),
				},
			}).
			Build())
	})

	AfterEach(func() {
		env.Config.DriftDetection = originalDrift
		env.Config.EnableTemplating = originalTemplating
	})

	DescribeTable("compares the old object with its templates resolved",
		func(oldValue, newValue, remoteValue string, expectDrift bool) {
			errs := validateTemplated(oldValue, newValue, remoteValue)

			if expectDrift {
				Expect(errs.Severe).To(ContainElement(HaveField("Message", ContainSubstring("drift detected"))))
			} else {
				Expect(errs.Severe).To(BeEmpty())
				Expect(errs.Warning).To(BeEmpty())
			}
		},
		Entry("a template pointed at another secret key is a CRD-only change",
			currentSecret, rotatedSecret, "current-value", false),
		Entry("a templated value replaced by a plain one is a CRD-only change",
			currentSecret, "plain", "current-value", false),
		Entry("an unchanged template is no change at all",
			currentSecret, currentSecret, "current-value", false),
		Entry("a remote change behind an unchanged template is still drift",
			currentSecret, currentSecret, "changed-in-apim", true),
		Entry("a remote change behind a changed template is still drift",
			currentSecret, rotatedSecret, "changed-in-apim", true),
	)

	It("skips drift detection with a warning when the old templates no longer resolve", func() {
		errs := validateTemplated(deletedSecret, currentSecret, "changed-in-apim")

		Expect(errs.Severe).To(BeEmpty())
		Expect(errs.Warning).To(ConsistOf(HaveField("Message",
			ContainSubstring("drift detection skipped, the templates of the previous version no longer resolve"))))
	})

	It("refuses an update whose new templates do not resolve", func() {
		errs := validateTemplated(currentSecret, deletedSecret, "current-value")

		Expect(errs.Severe).To(ConsistOf(HaveField("Message",
			ContainSubstring("could not compile templates for 'new' CRD"))))
	})
})

func validateTemplated(oldValue, newValue, remoteValue string) *errors.AdmissionErrors {
	stored := templatedDictionary(oldValue)
	updated := templatedDictionary(newValue)
	remote := templatedDictionary(remoteValue)

	toDTO := func(crd *v1alpha1.Dictionary) any {
		return model.ToDictionaryDTO(crd.Spec.Type, "default-templated")
	}
	return admdrift.ValidateDriftWithContext(
		context.Background(),
		stored,
		updated,
		func(context.Context) (*apim.APIM, error) {
			return &apim.APIM{}, nil
		},
		func(context.Context, *v1alpha1.Dictionary) error {
			return nil
		},
		func(*apim.APIM, *v1alpha1.Dictionary) (any, error) {
			return toDTO(remote), nil
		},
		toDTO,
	)
}

func templatedDictionary(value string) *v1alpha1.Dictionary {
	return &v1alpha1.Dictionary{
		TypeMeta:   metav1.TypeMeta{Kind: "Dictionary", APIVersion: "gravitee.io/v1alpha1"},
		ObjectMeta: metav1.ObjectMeta{Namespace: templatedNamespace, Name: "templated"},
		Spec: v1alpha1.DictionarySpec{
			Type: dictionary.Type{
				Name:           "templated",
				DictionaryType: dictionary.ManualType,
				Manual: &dictionary.ManualSpec{
					Properties: map[string]string{"clientSecret": value},
				},
			},
		},
	}
}
