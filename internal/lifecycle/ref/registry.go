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
	"errors"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

var (
	ErrNotImplemented = errors.New("lifecycle/ref: not implemented")
	ErrUnknownKind    = errors.New("lifecycle/ref: unknown kind")
)

type ObjectFactory func() client.Object

// ExtractFunc pulls a value out of a fetched object (e.g. secret data key).
// Nil on RefKind means the whole object is the resolved value.
type ExtractFunc func(obj client.Object, key string) (any, error)

type RefKind struct {
	New     ObjectFactory
	Extract ExtractFunc
}

var kinds = map[string]RefKind{}

func Register(name string, new ObjectFactory, extract ExtractFunc) {
	kinds[name] = RefKind{New: new, Extract: extract}
}

func Lookup(name string) (RefKind, bool) {
	k, ok := kinds[name]
	return k, ok
}

// Init registers known kinds. Call once at process start. Do not register AMSecurityDomain until that CRD exists.
func Init() {
	Register("amcontext", func() client.Object { return &v1alpha1.AMContext{} }, nil)
	Register("secret", func() client.Object { return &corev1.Secret{} }, extractSecretKey)
}

func extractSecretKey(obj client.Object, key string) (any, error) {
	return nil, ErrNotImplemented
}
