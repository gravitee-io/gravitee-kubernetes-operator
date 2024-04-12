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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"
	coreV1 "k8s.io/api/core/v1"
	netV1 "k8s.io/api/networking/v1"
)

type Files struct {
	Secrets     []string
	ConfigMaps  []string
	Context     string
	Resource    string
	API         string
	Application string
	Ingress     string
}

type FSBuilder struct {
	files *Files
}

func Builder() *FSBuilder {
	return &FSBuilder{
		files: &Files{
			Secrets:    []string{},
			ConfigMaps: []string{},
		},
	}
}

func (b *FSBuilder) Build() *Objects {
	f := b.files
	obj := &Objects{}

	obj.Secrets = decodeList(f.Secrets, &coreV1.Secret{}, secKind)
	obj.ConfigMaps = decodeList(f.ConfigMaps, &coreV1.ConfigMap{}, cmKind)

	suffix := random.GetSuffix()
	obj.randomSuffix = suffix

	if api := decodeIfDefined(f.API, &v1alpha1.ApiDefinition{}, apiKind); api != nil {
		obj.API = *api
		obj.API.Name += suffix
		obj.API.Namespace = constants.Namespace
		obj.API.Spec.Name = obj.API.Name

		randomizeAPIPaths(obj.API, suffix)
	}

	if app := decodeIfDefined(f.Application, &v1alpha1.Application{}, appKind); app != nil {
		obj.Application = *app
		obj.Application.Name += suffix
		obj.Application.Namespace = constants.Namespace
	}

	if ctx := decodeIfDefined(f.Context, &v1alpha1.ManagementContext{}, ctxKind); ctx != nil {
		obj.Context = *ctx
		obj.Context.Name += suffix
		obj.Context.Namespace = constants.Namespace
		if obj.API != nil {
			obj.API.Spec.Context = obj.Context.GetNamespacedName()
		}
		if obj.Application != nil {
			obj.Application.Spec.Context = obj.Context.GetNamespacedName()
		}
	}

	if rsc := decodeIfDefined(f.Resource, &v1alpha1.ApiResource{}, rscKind); rsc != nil {
		obj.Resource = *rsc
		obj.Resource.Name += suffix
		obj.Resource.Namespace = constants.Namespace
		if obj.API != nil {
			obj.API.Spec.Resources = []*base.ResourceOrRef{
				{
					Ref: &refs.NamespacedName{
						Name:      obj.Resource.Name,
						Namespace: constants.Namespace,
					},
				},
			}
		}
	}

	if ing := decodeIfDefined(f.Ingress, &netV1.Ingress{}, ingKind); ing != nil {
		obj.Ingress = *ing
		obj.Ingress.Name += suffix
		obj.Ingress.Namespace = constants.Namespace
		if obj.API != nil && isTemplate(obj.API) {
			obj.Ingress.Annotations[keys.IngressTemplateAnnotation] = obj.API.Name
		}

		randomizeIngressRules(obj.Ingress, suffix)
	}

	return obj
}

func randomizeAPIPaths(api *v1alpha1.ApiDefinition, suffix string) {
	if !isTemplate(api) {
		for _, vh := range api.Spec.Proxy.VirtualHosts {
			vh.Path = "/" + suffix[1:]
		}
	}
}

func randomizeIngressRules(ing *netV1.Ingress, suffix string) {
	for i := range ing.Spec.Rules {
		for j := range ing.Spec.Rules[i].HTTP.Paths {
			ing.Spec.Rules[i].HTTP.Paths[j].Path += suffix
		}
	}
}

func isTemplate(api *v1alpha1.ApiDefinition) bool {
	return api.Annotations[keys.IngressTemplateAnnotation] == env.TrueString
}

func (b *FSBuilder) AddSecret(file string) *FSBuilder {
	b.files.Secrets = append(b.files.Secrets, file)
	return b
}

func (b *FSBuilder) AddConfigMap(file string) *FSBuilder {
	b.files.ConfigMaps = append(b.files.ConfigMaps, file)
	return b
}

func (b *FSBuilder) WithContext(file string) *FSBuilder {
	b.files.Context = file
	return b
}

func (b *FSBuilder) WithResource(file string) *FSBuilder {
	b.files.Resource = file
	return b
}

func (b *FSBuilder) WithAPI(file string) *FSBuilder {
	b.files.API = file
	return b
}

func (b *FSBuilder) WithApplication(file string) *FSBuilder {
	b.files.Application = file
	return b
}

func (b *FSBuilder) WithIngress(file string) *FSBuilder {
	b.files.Ingress = file
	return b
}
