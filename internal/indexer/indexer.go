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
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	v1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type IndexField string

const (
	ContextField     IndexField = "context"
	SecretRefField   IndexField = "secretRef"
	ResourceField    IndexField = "resource"
	ApiTemplateField IndexField = "api-template"
	TLSSecretField   IndexField = "tls-secret"
	AppContextField  IndexField = "app-context"
)

func (f IndexField) String() string {
	return string(f)
}

type Indexer struct {
	Field string
	Func  Func
}

type Func = func(obj client.Object) []string

func CreateIndexerFunc[T client.Object](doIndex func(T, *[]string)) Func {
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
		Func:  CreateIndexerFunc(doIndex),
	}
}

func IndexManagementContexts(api *gio.ApiDefinition, fields *[]string) {
	if api.Spec.Context == nil {
		return
	}

	*fields = append(*fields, api.Spec.Context.String())
}

func IndexManagementContextSecrets(context *gio.ManagementContext, fields *[]string) {
	if context.Spec.HasSecretRef() {
		*fields = append(*fields, context.Spec.SecretRef().String())
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

func IndexApiTemplate(ing *v1.Ingress, fields *[]string) {
	if ing.Annotations[keys.IngressTemplateAnnotation] == "" {
		return
	}

	*fields = append(*fields, ing.Namespace+"/"+ing.Annotations[keys.IngressTemplateAnnotation])
}

func IndexTLSSecret(ing *v1.Ingress, fields *[]string) {
	if ing.Annotations[keys.IngressClassAnnotation] != keys.IngressClassAnnotationValue {
		return
	}

	if ing.Spec.TLS == nil || len(ing.Spec.TLS) == 0 {
		return
	}

	for i := range ing.Spec.TLS {
		*fields = append(*fields, ing.Namespace+"/"+ing.Spec.TLS[i].SecretName)
	}
}

func IndexApplicationManagementContexts(application *gio.Application, fields *[]string) {
	if application.Spec.Context == nil {
		return
	}

	*fields = append(*fields, application.Spec.Context.String())
}
