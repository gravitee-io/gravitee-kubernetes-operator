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

package indexer

import (
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type IndexField string

const (
	ContextField  IndexField = "context"
	ResourceField IndexField = "resource"
)

func (f IndexField) String() string {
	return string(f)
}

type Indexer struct {
	Field string
	Func  Func
}

type Func = func(obj client.Object) []string

func CreateIndexerFunc[T client.Object](field IndexField, doIndex func(T, *[]string)) Func {
	return func(obj client.Object) []string {
		fields := []string{}
		o, ok := obj.(T)

		if !ok {
			return fields
		}

		doIndex(o, &fields)

		return fields
	}
}

func NewIndexer[T client.Object](field IndexField, doIndex func(T, *[]string)) Indexer {
	return Indexer{
		Field: string(field),
		Func:  CreateIndexerFunc(field, doIndex),
	}
}

func IndexApiContexts(api *gio.ApiDefinition, fields *[]string) {
	if api.Spec.Contexts == nil {
		return
	}

	for location := range api.Status.Contexts {
		*fields = append(*fields, location)
	}
}

func IndexApiResourceRefs(api *gio.ApiDefinition, fields *[]string) {
	if api.Spec.Resources == nil {
		return
	}

	for _, resource := range api.Spec.Resources {
		if resource.IsRef() {
			*fields = append(*fields, resource.Ref.String())
		}
	}
}
