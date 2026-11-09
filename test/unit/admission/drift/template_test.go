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
	"encoding/json"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/dictionary"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	admdrift "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/template"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	templatedNamespace = "default"
	currentSecret      = "[[ secret `client/current` ]]"
	rotatedSecret      = "[[ secret `client/rotated` ]]"
	bracketsSecret     = "[[ secret `client/brackets` ]]"
	deletedSecret      = "[[ secret `deleted/value` ]]"
	removedKeySecret   = "[[ secret `client/removed` ]]"
	emptiedKeySecret   = "[[ secret `client/emptied` ]]"
	unreadableSecret   = "[[ secret `unreadable/value` ]]"
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
					"current":  []byte("current-value"),
					"rotated":  []byte("rotated-value"),
					"brackets": []byte("ab[[cd"),
					"emptied":  {},
				},
			}).
			WithInterceptorFuncs(interceptor.Funcs{Get: failToReadUnreadable}).
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
		Entry("a resolved value holding the template delimiters is not compiled again",
			currentSecret, bracketsSecret, "ab[[cd", false),
		Entry("a remote change behind an unchanged template is still drift",
			currentSecret, currentSecret, "changed-in-apim", true),
		Entry("a remote change behind a changed template is still drift",
			currentSecret, rotatedSecret, "changed-in-apim", true),
	)

	DescribeTable("skips drift detection with a warning when a secret of the old version no longer exists",
		func(oldValue string) {
			errs := validateTemplated(oldValue, currentSecret, "changed-in-apim")

			Expect(errs.Severe).To(BeEmpty())
			Expect(errs.Warning).To(ConsistOf(HaveField("Message",
				ContainSubstring("drift detection skipped, the previous version no longer resolves"))))
		},
		Entry("the Secret was deleted", deletedSecret),
		Entry("the key was removed from the Secret", removedKeySecret),
		Entry("the key was emptied in the Secret", emptiedKeySecret),
	)

	It("refuses the update when a Secret of the old version cannot be read", func() {
		errs := validateTemplated(unreadableSecret, currentSecret, "changed-in-apim")

		Expect(errs.Warning).To(BeEmpty())
		Expect(errs.Severe).To(ConsistOf(HaveField("Message",
			ContainSubstring("could not resolve the previous version for drift detection"))))
	})

	DescribeTable("handles a reference of the old version that does not resolve",
		func(failure error, expectSkipped bool) {
			failOnOld := func(_ context.Context, dict *v1alpha1.Dictionary) error {
				if dict.Spec.Manual.Properties["clientSecret"] == "old-ref" {
					return failure
				}
				return nil
			}
			updated := templatedDictionary("new-ref")

			errs := validate(admissionRequest(updated), templatedDictionary("old-ref"), updated,
				"changed-in-apim", failOnOld)

			if expectSkipped {
				Expect(errs.Severe).To(BeEmpty())
				Expect(errs.Warning).To(ConsistOf(HaveField("Message",
					ContainSubstring("drift detection skipped, the previous version no longer resolves"))))
			} else {
				Expect(errs.Warning).To(BeEmpty())
				Expect(errs.Severe).To(ConsistOf(HaveField("Message",
					ContainSubstring("could not resolve the previous version for drift detection"))))
			}
		},
		Entry("a referenced object deleted since skips drift with a warning",
			errors.NewResolveRefError(apierrors.NewNotFound(schema.GroupResource{Resource: "apiresources"}, "gone")), true),
		Entry("a referenced object that cannot be read refuses the update",
			errors.NewResolveRefError(apierrors.NewTimeoutError("the API server did not answer in time", 1)), false),
	)

	It("never prints a value a secret template resolved to", func() {
		errs := validateTemplated(currentSecret, rotatedSecret, "changed-in-apim")

		Expect(errs.Severe).To(HaveLen(1))
		Expect(errs.Severe[0].Message).To(ContainSubstring("clientSecret: (secret value, not shown)"))
		Expect(errs.Severe[0].Message).NotTo(ContainSubstring("current-value"))
		Expect(errs.Severe[0].Message).NotTo(ContainSubstring("rotated-value"))
		Expect(errs.Severe[0].Message).NotTo(ContainSubstring("changed-in-apim"))
	})

	It("never prints a value a secret template of a referenced object resolved to", func() {
		inlineReferenced := func(ctx context.Context, dict *v1alpha1.Dictionary) error {
			referenced := templatedDictionary(currentSecret)
			if err := template.Compile(ctx, referenced, false); err != nil {
				return err
			}
			dict.Spec.Manual = referenced.Spec.Manual
			return nil
		}
		updated := templatedDictionary("from-ref")

		errs := validate(admissionRequest(updated), templatedDictionary("from-ref"), updated,
			"changed-in-apim", inlineReferenced)

		Expect(errs.Severe).To(HaveLen(1))
		Expect(errs.Severe[0].Message).To(ContainSubstring("clientSecret: (secret value, not shown)"))
		Expect(errs.Severe[0].Message).NotTo(ContainSubstring("current-value"))
	})

	It("refuses an update whose new templates do not resolve", func() {
		updated := templatedDictionary(deletedSecret)

		errs := validate(admissionRequest(updated), templatedDictionary(currentSecret), updated,
			"current-value", noRefs)

		Expect(errs.Severe).To(ConsistOf(HaveField("Message",
			ContainSubstring("could not compile templates for 'new' CRD"))))
	})
})

func failToReadUnreadable(
	ctx context.Context, cli client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption,
) error {
	if key.Name == "unreadable" {
		return apierrors.NewTimeoutError("the API server did not answer in time", 1)
	}
	return cli.Get(ctx, key, obj, opts...)
}

// validateTemplated runs drift as the webhook does: the request carries the submitted object, and
// validation compiled the new object in place before drift runs.
func validateTemplated(oldValue, newValue, remoteValue string) *errors.AdmissionErrors {
	updated := templatedDictionary(newValue)
	ctx := admissionRequest(updated)
	Expect(template.Compile(ctx, updated, false)).To(Succeed())
	return validate(ctx, templatedDictionary(oldValue), updated, remoteValue, noRefs)
}

func admissionRequest(submitted *v1alpha1.Dictionary) context.Context {
	raw, err := json.Marshal(submitted)
	Expect(err).NotTo(HaveOccurred())
	return admission.NewContextWithRequest(context.Background(), admission.Request{
		AdmissionRequest: admissionv1.AdmissionRequest{Object: runtime.RawExtension{Raw: raw}},
	})
}

func noRefs(context.Context, *v1alpha1.Dictionary) error {
	return nil
}

func validate(
	ctx context.Context,
	stored, updated *v1alpha1.Dictionary,
	remoteValue string,
	resolveRefs admdrift.RefResolver[*v1alpha1.Dictionary],
) *errors.AdmissionErrors {
	remote := templatedDictionary(remoteValue)

	toDTO := func(crd *v1alpha1.Dictionary) any {
		return model.ToDictionaryDTO(crd.Spec.Type, "default-templated")
	}
	return admdrift.ValidateDriftWithContext(
		ctx,
		stored,
		updated,
		func(context.Context) (*apim.APIM, error) {
			return &apim.APIM{}, nil
		},
		resolveRefs,
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
