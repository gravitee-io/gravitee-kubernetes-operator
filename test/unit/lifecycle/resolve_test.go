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

package lifecycle_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle/ref"
)

const parentNs = "apps"

type secretHolder struct {
	Content   []byte `ref:"secret,SecretRef,tls.crt"`
	SecretRef *refs.NamespacedName
}

type pluralSecretHolder struct {
	Content   []byte `ref:"secrets,SecretRef,tls.crt"`
	SecretRef *refs.NamespacedName
}

type conventionHolder struct {
	Content    []byte `ref:"secret,ContentRef,tls.crt"`
	ContentRef *refs.NamespacedName
}

type conventionKindOnly struct {
	DomainKey    string `ref:"amsecuritydomain"`
	DomainKeyRef *refs.NamespacedName
}

type nestedHolder struct {
	Inner secretHolder
}

type listHolder struct {
	Items []secretHolder
}

type hridHolder struct {
	DomainKey string `ref:"amsecuritydomain,DomainRef"`
	DomainRef *refs.NamespacedName
}

type valueRefHolder struct {
	Content   []byte `ref:"secret,SecretRef,tls.crt"`
	SecretRef refs.NamespacedName
}

type ptrNestedHolder struct {
	Inner *secretHolder
}

type ptrListHolder struct {
	Items []*secretHolder
}

type stringSecretHolder struct {
	Content   string `ref:"secret,SecretRef,tls.crt"`
	SecretRef *refs.NamespacedName
}

type omittedRefFieldHolder struct {
	Content    []byte `ref:"secret,,tls.crt"`
	ContentRef *refs.NamespacedName
}

func registerClient(objects ...client.Object) {
	scheme := runtime.NewScheme()
	Expect(clientgoscheme.AddToScheme(scheme)).To(Succeed())
	Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
	k8s.RegisterClient(fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(objects...).
		Build())
}

func tlsSecret(ns, name string, cert []byte) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Namespace: ns, Name: name},
		Data:       map[string][]byte{"tls.crt": cert},
	}
}

func expectHRID(ns, name string) string {
	n := refs.NewNamespacedName(ns, name)
	return n.HRID()
}

func amContext(ns, name string) *v1alpha1.AMContext {
	return &v1alpha1.AMContext{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1alpha1.GroupVersion.String(),
			Kind:       "AMContext",
		},
		ObjectMeta: metav1.ObjectMeta{Namespace: ns, Name: name},
	}
}

var _ = Describe("GenericRefResolver", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
	})

	It("fills a []byte from secret data by explicit ref field and key", func() {
		sec := tlsSecret(parentNs, "tls", []byte("cert-bytes"))
		registerClient(sec)

		obj := &secretHolder{
			SecretRef: &refs.NamespacedName{Name: "tls"},
		}
		Expect(ref.GenericRefResolver(ctx, obj, parentNs)).To(Succeed())
		Expect(obj.Content).To(Equal([]byte("cert-bytes")))
	})

	It("resolves a plural kind tag against a singular registry name", func() {
		sec := tlsSecret(parentNs, "tls", []byte("plural-tag"))
		registerClient(sec)

		obj := &pluralSecretHolder{
			SecretRef: &refs.NamespacedName{Name: "tls"},
		}
		Expect(ref.GenericRefResolver(ctx, obj, parentNs)).To(Succeed())
		Expect(obj.Content).To(Equal([]byte("plural-tag")))
	})

	It("uses an explicit ref field named <field>Ref", func() {
		sec := tlsSecret(parentNs, "tls", []byte("from-convention"))
		registerClient(sec)

		obj := &conventionHolder{
			ContentRef: &refs.NamespacedName{Name: "tls"},
		}
		Expect(ref.GenericRefResolver(ctx, obj, parentNs)).To(Succeed())
		Expect(obj.Content).To(Equal([]byte("from-convention")))
	})

	It("defaults the sibling field to <field>Ref when the tag is kind only", func() {
		am := amContext(parentNs, "ctx-b")
		registerClient(am)

		obj := &conventionKindOnly{
			DomainKeyRef: &refs.NamespacedName{Name: "ctx-b"},
		}
		Expect(ref.GenericRefResolver(ctx, obj, parentNs)).To(Succeed())
		Expect(obj.DomainKey).To(Equal(expectHRID(parentNs, "ctx-b")))
	})

	It("uses parentNs when the ref has no namespace", func() {
		sec := tlsSecret(parentNs, "tls", []byte("parent-ns"))
		registerClient(sec)

		obj := &secretHolder{
			SecretRef: &refs.NamespacedName{Name: "tls"},
		}
		Expect(ref.GenericRefResolver(ctx, obj, parentNs)).To(Succeed())
		Expect(obj.Content).To(Equal([]byte("parent-ns")))
	})

	It("uses the ref namespace when set", func() {
		sec := tlsSecret("other", "tls", []byte("other-ns"))
		registerClient(sec)

		obj := &secretHolder{
			SecretRef: &refs.NamespacedName{Namespace: "other", Name: "tls"},
		}
		Expect(ref.GenericRefResolver(ctx, obj, parentNs)).To(Succeed())
		Expect(obj.Content).To(Equal([]byte("other-ns")))
	})

	It("walks a nested struct", func() {
		sec := tlsSecret(parentNs, "tls", []byte("nested"))
		registerClient(sec)

		obj := &nestedHolder{
			Inner: secretHolder{
				SecretRef: &refs.NamespacedName{Name: "tls"},
			},
		}
		Expect(ref.GenericRefResolver(ctx, obj, parentNs)).To(Succeed())
		Expect(obj.Inner.Content).To(Equal([]byte("nested")))
	})

	It("walks a slice of structs", func() {
		sec := tlsSecret(parentNs, "tls", []byte("item"))
		registerClient(sec)

		obj := &listHolder{
			Items: []secretHolder{
				{SecretRef: &refs.NamespacedName{Name: "tls"}},
				{SecretRef: &refs.NamespacedName{Name: "tls"}},
			},
		}
		Expect(ref.GenericRefResolver(ctx, obj, parentNs)).To(Succeed())
		Expect(obj.Items[0].Content).To(Equal([]byte("item")))
		Expect(obj.Items[1].Content).To(Equal([]byte("item")))
	})

	It("extracts HRID from a referenced AMContext", func() {
		am := amContext(parentNs, "ctx-a")
		registerClient(am)

		obj := &hridHolder{
			DomainRef: &refs.NamespacedName{Name: "ctx-a"},
		}
		Expect(ref.GenericRefResolver(ctx, obj, parentNs)).To(Succeed())
		Expect(obj.DomainKey).To(Equal(expectHRID(parentNs, "ctx-a")))
	})

	It("wraps unknown kind as Error and ErrUnknownKind", func() {
		type badKind struct {
			X    string `ref:"nope"`
			XRef *refs.NamespacedName
		}
		registerClient()
		err := ref.GenericRefResolver(ctx, &badKind{XRef: &refs.NamespacedName{Name: "x"}}, parentNs)
		Expect(err).To(HaveOccurred())
		Expect(errors.Is(err, ref.ErrUnknownKind)).To(BeTrue())
		Expect(errors.Is(err, &ref.Error{Op: "resolve", Kind: "nope"})).To(BeTrue())
	})

	It("unwraps k8s NotFound from a missing secret", func() {
		registerClient()
		obj := &secretHolder{
			SecretRef: &refs.NamespacedName{Name: "missing"},
		}
		err := ref.GenericRefResolver(ctx, obj, parentNs)
		Expect(err).To(HaveOccurred())
		Expect(apierrors.IsNotFound(err)).To(BeTrue())
		Expect(errors.Is(err, &ref.Error{Op: "resolve"})).To(BeTrue())
	})

	It("rejects a non-struct root", func() {
		registerClient()
		n := 1
		err := ref.GenericRefResolver(ctx, &n, parentNs)
		Expect(errors.Is(err, ref.ErrStructKind)).To(BeTrue())
	})

	It("skips resolve when the target field is already set", func() {
		sec := tlsSecret(parentNs, "tls", []byte("from-cluster"))
		registerClient(sec)

		obj := &secretHolder{
			Content:   []byte("already"),
			SecretRef: &refs.NamespacedName{Name: "tls"},
		}
		Expect(ref.GenericRefResolver(ctx, obj, parentNs)).To(Succeed())
		Expect(obj.Content).To(Equal([]byte("already")))
	})

	It("errors when the sibling ref is nil", func() {
		registerClient(tlsSecret(parentNs, "tls", []byte("x")))
		obj := &secretHolder{}
		err := ref.GenericRefResolver(ctx, obj, parentNs)
		Expect(err).To(HaveOccurred())
		Expect(errors.Is(err, ref.ErrMissingRef)).To(BeTrue())
	})

	It("accepts a value NamespacedName sibling", func() {
		sec := tlsSecret(parentNs, "tls", []byte("from-value"))
		registerClient(sec)

		obj := &valueRefHolder{
			SecretRef: refs.NamespacedName{Name: "tls"},
		}
		Expect(ref.GenericRefResolver(ctx, obj, parentNs)).To(Succeed())
		Expect(obj.Content).To(Equal([]byte("from-value")))
	})

	It("walks a pointer nested struct", func() {
		sec := tlsSecret(parentNs, "tls", []byte("ptr-nested"))
		registerClient(sec)

		obj := &ptrNestedHolder{
			Inner: &secretHolder{
				SecretRef: &refs.NamespacedName{Name: "tls"},
			},
		}
		Expect(ref.GenericRefResolver(ctx, obj, parentNs)).To(Succeed())
		Expect(obj.Inner.Content).To(Equal([]byte("ptr-nested")))
	})

	It("walks a slice of struct pointers", func() {
		sec := tlsSecret(parentNs, "tls", []byte("ptr-item"))
		registerClient(sec)

		obj := &ptrListHolder{
			Items: []*secretHolder{
				{SecretRef: &refs.NamespacedName{Name: "tls"}},
			},
		}
		Expect(ref.GenericRefResolver(ctx, obj, parentNs)).To(Succeed())
		Expect(obj.Items[0].Content).To(Equal([]byte("ptr-item")))
	})

	It("succeeds on an empty slice", func() {
		registerClient()
		obj := &listHolder{Items: []secretHolder{}}
		Expect(ref.GenericRefResolver(ctx, obj, parentNs)).To(Succeed())
		Expect(obj.Items).To(BeEmpty())
	})

	It("errors on a type mismatch between extract and target", func() {
		sec := tlsSecret(parentNs, "tls", []byte("cert"))
		registerClient(sec)

		obj := &stringSecretHolder{
			SecretRef: &refs.NamespacedName{Name: "tls"},
		}
		err := ref.GenericRefResolver(ctx, obj, parentNs)
		Expect(errors.Is(err, ref.ErrTypeMismatch)).To(BeTrue())
	})

	It("errors when the root struct is not addressable", func() {
		registerClient(tlsSecret(parentNs, "tls", []byte("x")))
		obj := secretHolder{
			SecretRef: &refs.NamespacedName{Name: "tls"},
		}
		err := ref.GenericRefResolver(ctx, obj, parentNs)
		Expect(errors.Is(err, ref.ErrCannotSetField)).To(BeTrue())
	})

	It("errors when the secret has no such data key", func() {
		sec := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Namespace: parentNs, Name: "tls"},
			Data:       map[string][]byte{"other": []byte("x")},
		}
		registerClient(sec)

		obj := &secretHolder{
			SecretRef: &refs.NamespacedName{Name: "tls"},
		}
		err := ref.GenericRefResolver(ctx, obj, parentNs)
		Expect(errors.Is(err, ref.ErrSecretKeyMissing)).To(BeTrue())
	})

	It("resolves secret,,tls.crt using ContentRef", func() {
		sec := tlsSecret(parentNs, "tls", []byte("omitted-ref"))
		registerClient(sec)

		obj := &omittedRefFieldHolder{
			ContentRef: &refs.NamespacedName{Name: "tls"},
		}
		Expect(ref.GenericRefResolver(ctx, obj, parentNs)).To(Succeed())
		Expect(obj.Content).To(Equal([]byte("omitted-ref")))
	})
})
