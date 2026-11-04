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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/dictionary"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/template"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("Compile encrypted dictionary properties", func() {
	var (
		ctx             context.Context
		prevEnableTempl bool
	)

	BeforeEach(func() {
		ctx = context.Background()
		prevEnableTempl = env.Config.EnableTemplating
		env.Config.EnableTemplating = true

		scheme := runtime.NewScheme()
		Expect(clientgoscheme.AddToScheme(scheme)).To(Succeed())
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())

		k8s.RegisterClient(fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(&corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "backend-creds", Namespace: testNamespace},
				Data:       map[string][]byte{"api-key": []byte("s3cret")},
			}).
			Build())
	})

	AfterEach(func() {
		env.Config.EnableTemplating = prevEnableTempl
	})

	It("resolves a secret reference inside encryptedProperties", func() {
		dict := &v1alpha1.Dictionary{
			TypeMeta: metav1.TypeMeta{
				APIVersion: v1alpha1.GroupVersion.String(),
				Kind:       "Dictionary",
			},
			ObjectMeta: metav1.ObjectMeta{Name: "secret-dict", Namespace: testNamespace},
			Spec: v1alpha1.DictionarySpec{
				Type: dictionary.Type{
					Name:           "secret-dict",
					DictionaryType: dictionary.ManualType,
					Manual: &dictionary.ManualSpec{
						Properties: map[string]string{"url": "https://backend"},
						EncryptedProperties: map[string]dictionary.EncryptableValue{
							"apiKey": {Value: "[[ secret `backend-creds/api-key` ]]"},
						},
					},
				},
			},
		}

		Expect(template.Compile(ctx, dict, true)).To(Succeed())

		Expect(dict.Spec.Manual.EncryptedProperties["apiKey"].Value).To(Equal("s3cret"))
		Expect(dict.Spec.Manual.Properties["url"]).To(Equal("https://backend"))
	})
})
