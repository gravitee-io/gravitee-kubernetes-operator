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

package template_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/dictionary"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/template"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	testNamespace            = "default"
	totalReferenceAnnotation = "gravitee.io/references"
)

var _ = Describe("ReleaseReferences", func() {
	var (
		ctx             context.Context
		scheme          *runtime.Scheme
		prevEnableTempl bool
		dictName        string
		dictObjectID    string
		deletingDict    *v1alpha1.Dictionary
		activeDict      *v1alpha1.Dictionary
	)

	BeforeEach(func() {
		ctx = context.Background()
		prevEnableTempl = env.Config.EnableTemplating
		env.Config.EnableTemplating = true

		scheme = runtime.NewScheme()
		Expect(clientgoscheme.AddToScheme(scheme)).To(Succeed())
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())

		dictName = "test-dict"
		dictObjectID = testNamespace + "/" + dictName
		deletingDict = deletingDictionary(dictName, testNamespace)
		activeDict = &v1alpha1.Dictionary{
			TypeMeta: metav1.TypeMeta{
				APIVersion: v1alpha1.GroupVersion.String(),
				Kind:       "Dictionary",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      dictName,
				Namespace: testNamespace,
			},
		}
	})

	AfterEach(func() {
		env.Config.EnableTemplating = prevEnableTempl
	})

	registerClient := func(objects ...client.Object) {
		k8s.RegisterClient(fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(objects...).
			Build())
	}

	getSecret := func(name string) *corev1.Secret {
		sec := &corev1.Secret{}
		Expect(k8s.GetClient().Get(ctx, types.NamespacedName{
			Namespace: testNamespace,
			Name:      name,
		}, sec)).To(Succeed())
		return sec
	}

	getConfigMap := func(name string) *corev1.ConfigMap {
		cm := &corev1.ConfigMap{}
		Expect(k8s.GetClient().Get(ctx, types.NamespacedName{
			Namespace: testNamespace,
			Name:      name,
		}, cm)).To(Succeed())
		return cm
	}

	annotationKey := func() string {
		return templatingAnnotationKey(deletingDict)
	}

	It("is a no-op when templating is disabled", func() {
		env.Config.EnableTemplating = false
		sec := templatingSecret("tpl-secret", dictObjectID, 1, annotationKey())
		registerClient(sec, deletingDict)

		Expect(template.ReleaseReferences(ctx, deletingDict)).To(Succeed())

		updated := getSecret("tpl-secret")
		Expect(util.ContainsFinalizer(updated, core.TemplatingFinalizer)).To(BeTrue())
		Expect(updated.Annotations[annotationKey()]).To(Equal(mustMarshal([]string{dictObjectID})))
	})

	It("is a no-op when the parent is not being deleted", func() {
		sec := templatingSecret("tpl-secret", dictObjectID, 1, templatingAnnotationKey(activeDict))
		registerClient(sec, activeDict)

		Expect(template.ReleaseReferences(ctx, activeDict)).To(Succeed())

		updated := getSecret("tpl-secret")
		Expect(util.ContainsFinalizer(updated, core.TemplatingFinalizer)).To(BeTrue())
		Expect(updated.Annotations[totalReferenceAnnotation]).To(Equal("1"))
	})

	It("removes the finalizer when releasing the last reference on a Secret", func() {
		sec := templatingSecret("tpl-secret", dictObjectID, 1, annotationKey())
		registerClient(sec, deletingDict)

		Expect(template.ReleaseReferences(ctx, deletingDict)).To(Succeed())

		updated := getSecret("tpl-secret")
		Expect(util.ContainsFinalizer(updated, core.TemplatingFinalizer)).To(BeFalse())
		Expect(updated.Annotations[annotationKey()]).To(Equal(mustMarshal([]string{})))
		Expect(updated.Annotations[totalReferenceAnnotation]).To(Equal("0"))
	})

	It("keeps the finalizer when other parents still reference the Secret", func() {
		otherParentID := testNamespace + "/other-dict"
		sec := templatingSecret("tpl-secret", dictObjectID, 2, annotationKey(), otherParentID)
		registerClient(sec, deletingDict)

		Expect(template.ReleaseReferences(ctx, deletingDict)).To(Succeed())

		updated := getSecret("tpl-secret")
		Expect(util.ContainsFinalizer(updated, core.TemplatingFinalizer)).To(BeTrue())
		Expect(updated.Annotations[annotationKey()]).To(Equal(mustMarshal([]string{otherParentID})))
		Expect(updated.Annotations[totalReferenceAnnotation]).To(Equal("1"))
	})

	It("releases templating references on a ConfigMap", func() {
		cm := templatingConfigMap("tpl-configmap", dictObjectID, 1, annotationKey())
		registerClient(cm, deletingDict)

		Expect(template.ReleaseReferences(ctx, deletingDict)).To(Succeed())

		updated := getConfigMap("tpl-configmap")
		Expect(util.ContainsFinalizer(updated, core.TemplatingFinalizer)).To(BeFalse())
		Expect(updated.Annotations[annotationKey()]).To(Equal(mustMarshal([]string{})))
		Expect(updated.Annotations[totalReferenceAnnotation]).To(Equal("0"))
	})

	It("ignores Secrets without the templating finalizer", func() {
		sec := templatingSecret("tpl-secret", dictObjectID, 1, annotationKey())
		util.RemoveFinalizer(sec, core.TemplatingFinalizer)
		registerClient(sec, deletingDict)

		Expect(template.ReleaseReferences(ctx, deletingDict)).To(Succeed())

		updated := getSecret("tpl-secret")
		Expect(util.ContainsFinalizer(updated, core.TemplatingFinalizer)).To(BeFalse())
		Expect(updated.Annotations[annotationKey()]).To(Equal(mustMarshal([]string{dictObjectID})))
		Expect(updated.Annotations[totalReferenceAnnotation]).To(Equal("1"))
	})

	It("ignores Secrets referenced by a different parent", func() {
		otherParentID := testNamespace + "/other-dict"
		sec := templatingSecret("tpl-secret", otherParentID, 1, annotationKey())
		registerClient(sec, deletingDict)

		Expect(template.ReleaseReferences(ctx, deletingDict)).To(Succeed())

		updated := getSecret("tpl-secret")
		Expect(util.ContainsFinalizer(updated, core.TemplatingFinalizer)).To(BeTrue())
		Expect(updated.Annotations[annotationKey()]).To(Equal(mustMarshal([]string{otherParentID})))
		Expect(updated.Annotations[totalReferenceAnnotation]).To(Equal("1"))
	})

	It("only scans Secrets in the parent namespace", func() {
		otherNSSec := templatingSecret("other-ns-secret", dictObjectID, 1, annotationKey())
		otherNSSec.Namespace = "other"
		registerClient(otherNSSec, deletingDict)

		Expect(template.ReleaseReferences(ctx, deletingDict)).To(Succeed())

		updated := &corev1.Secret{}
		Expect(k8s.GetClient().Get(ctx, types.NamespacedName{
			Namespace: "other",
			Name:      "other-ns-secret",
		}, updated)).To(Succeed())
		Expect(util.ContainsFinalizer(updated, core.TemplatingFinalizer)).To(BeTrue())
		Expect(updated.Annotations[annotationKey()]).To(Equal(mustMarshal([]string{dictObjectID})))
	})

	It("releases references created by Compile", func() {
		sec := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "provider-secret",
				Namespace: testNamespace,
			},
			Data: map[string][]byte{
				"url": []byte("https://example.com"),
			},
		}
		dict := &v1alpha1.Dictionary{
			TypeMeta: metav1.TypeMeta{
				APIVersion: v1alpha1.GroupVersion.String(),
				Kind:       "Dictionary",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      dictName,
				Namespace: testNamespace,
			},
		}
		dict.Spec.Name = dictName
		dict.Spec.DictionaryType = dictionary.DynamicType
		dict.Spec.Dynamic = &dictionary.DynamicSpec{
			Provider: &dictionary.Provider{
				ProviderType: "HTTP",
				URL:          "[[ secret `provider-secret/url` ]]",
				Method:       "GET",
			},
			Trigger: &dictionary.Trigger{
				Rate: 1,
				Unit: dictionary.SecondsUnit,
			},
		}
		registerClient(sec, dict)

		Expect(template.Compile(ctx, dict, true)).To(Succeed())

		annotationKey := templatingAnnotationKey(dict)
		compiled := getSecret("provider-secret")
		Expect(util.ContainsFinalizer(compiled, core.TemplatingFinalizer)).To(BeTrue())
		Expect(compiled.Annotations[annotationKey]).To(Equal(mustMarshal([]string{dictObjectID})))

		now := metav1.Now()
		dict.DeletionTimestamp = &now
		Expect(template.ReleaseReferences(ctx, dict)).To(Succeed())

		released := getSecret("provider-secret")
		Expect(util.ContainsFinalizer(released, core.TemplatingFinalizer)).To(BeFalse())
		Expect(released.Annotations[annotationKey]).To(Equal(mustMarshal([]string{})))
		Expect(released.Annotations[totalReferenceAnnotation]).To(Equal("0"))
	})

	It("releases only matching Secrets when several exist in the namespace", func() {
		matching := templatingSecret("matching-secret", dictObjectID, 1, annotationKey())
		other := templatingSecret("other-secret", testNamespace+"/other-dict", 1, annotationKey())
		registerClient(matching, other, deletingDict)

		Expect(template.ReleaseReferences(ctx, deletingDict)).To(Succeed())

		released := getSecret("matching-secret")
		Expect(util.ContainsFinalizer(released, core.TemplatingFinalizer)).To(BeFalse())

		untouched := getSecret("other-secret")
		Expect(util.ContainsFinalizer(untouched, core.TemplatingFinalizer)).To(BeTrue())
	})
})

func deletingDictionary(name, namespace string) *v1alpha1.Dictionary {
	now := metav1.Now()
	return &v1alpha1.Dictionary{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1alpha1.GroupVersion.String(),
			Kind:       "Dictionary",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			Namespace:         namespace,
			DeletionTimestamp: &now,
			Finalizers:        []string{core.DictionaryFinalizer},
		},
	}
}

func templatingSecret(name, objectID string, totalRefs int, annotationKey string, extraObjectIDs ...string) *corev1.Secret {
	objectIDs := append([]string{objectID}, extraObjectIDs...)
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
			Annotations: map[string]string{
				annotationKey:            mustMarshal(objectIDs),
				totalReferenceAnnotation: itoa(totalRefs),
			},
			Finalizers: []string{core.TemplatingFinalizer},
		},
	}
}

func templatingConfigMap(name, objectID string, totalRefs int, annotationKey string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
			Annotations: map[string]string{
				annotationKey:            mustMarshal([]string{objectID}),
				totalReferenceAnnotation: itoa(totalRefs),
			},
			Finalizers: []string{core.TemplatingFinalizer},
		},
	}
}

func mustMarshal(values []string) string {
	b, err := json.Marshal(values)
	Expect(err).NotTo(HaveOccurred())
	return string(b)
}

func itoa(v int) string {
	return strconv.Itoa(v)
}

func templatingAnnotationKey(obj client.Object) string {
	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	Expect(err).NotTo(HaveOccurred())
	return "gravitee.io/" + strings.ToLower(fmt.Sprintf("%v", u["kind"])) + "s"
}
