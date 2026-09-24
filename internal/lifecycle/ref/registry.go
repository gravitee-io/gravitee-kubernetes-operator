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

// Package ref resolves `ref` struct tags onto a CR using k8s.GetClient().
package ref

import (
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

type ObjectFactory func() client.Object

// ExtractFunc pulls a value out of a fetched object (e.g. secret data key).
type ExtractFunc func(obj client.Object, key string) (any, error)

type Kind struct {
	New     ObjectFactory
	Extract ExtractFunc
}

var kinds = map[string]Kind{}

func Register(name string, factory ObjectFactory, extract ExtractFunc) {
	if name == "" {
		panic("lifecycle/ref: empty kind name")
	}
	if factory == nil {
		panic("lifecycle/ref: nil ObjectFactory")
	}
	if extract == nil {
		panic("lifecycle/ref: nil ExtractFunc")
	}
	kinds[name] = Kind{New: factory, Extract: extract}
}

func Lookup(name string) (Kind, bool) {
	k, ok := kinds[name]
	if !ok && !strings.HasSuffix(name, "s") {
		k, ok = kinds[name+"s"]
	} else if !ok {
		k, ok = kinds[strings.TrimSuffix(name, "s")]
	}
	return k, ok
}

// Init registers known kinds. Call once at process start. Do not register AMSecurityDomain until that CRD exists.
func Init() {
	Register(core.CRDAMContextResource, func() client.Object { return &v1alpha1.AMContext{} }, noop)
	Register(core.CRDAMSecurityDomainResource, func() client.Object { return &v1alpha1.AMSecurityDomain{} }, noop)
	Register("secret", func() client.Object { return &corev1.Secret{} }, extractSecretKey)
}

func noop(client.Object, string) (any, error) {
	return nil, nil
}

func extractSecretKey(obj client.Object, key string) (any, error) {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		return nil, NewWrappedError("extract", "secret", key, ErrNotASecret)
	}
	data, ok := secret.Data[key]
	if !ok {
		return nil, NewWrappedError("extract", "secret", key, ErrSecretKeyMissing)
	}
	return data, nil
}
