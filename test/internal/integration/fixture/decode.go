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

package fixture

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/examples"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	coreV1 "k8s.io/api/core/v1"
	netV1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/onsi/gomega"
)

var universalDecode = scheme.Codecs.UniversalDecoder().Decode

var (
	cmKind    = coreV1.SchemeGroupVersion.WithKind("ConfigMap")
	secKind   = coreV1.SchemeGroupVersion.WithKind("Secret")
	ingKind   = netV1.SchemeGroupVersion.WithKind("Ingress")
	ctxKind   = v1alpha1.GroupVersion.WithKind("ManagementContext")
	rscKind   = v1alpha1.GroupVersion.WithKind("ApiResource")
	appKind   = v1alpha1.GroupVersion.WithKind("Application")
	apiKind   = v1alpha1.GroupVersion.WithKind("ApiDefinition")
	apiV4Kind = v1alpha1.GroupVersion.WithKind("ApiV4Definition")
)

func decodeIfDefined[T client.Object](path string, rcv T, kind schema.GroupVersionKind) *T {
	if len(path) == 0 {
		return nil
	}
	obj := decode(path, rcv, kind)
	return &obj
}

func decodeList[T client.Object](paths []string, rcv T, kind schema.GroupVersionKind) []T {
	obj := make([]T, 0)
	for _, path := range paths {
		cp, ok := rcv.DeepCopyObject().(T)
		Expect(ok).To(BeTrue())
		cp.SetNamespace(constants.Namespace)
		obj = append(obj, decode(path, cp, kind))
	}
	return obj
}

func decode[T client.Object](path string, rcv T, kind schema.GroupVersionKind) T {
	data, err := examples.FS.ReadFile(path)
	Expect(err).ToNot(HaveOccurred())

	decoded, _, err := universalDecode(data, &kind, rcv)
	Expect(err).ToNot(HaveOccurred())

	obj, ok := decoded.(T)
	Expect(ok).To(BeTrue())

	return obj
}
