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
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/types/k8s/custom"
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
	APIv4       string
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
		b.setupAPI(obj, api, suffix)
	}

	if apiV4 := decodeIfDefined(f.APIv4, &v1alpha1.ApiV4Definition{}, apiV4Kind); apiV4 != nil {
		b.setupAPIv4(obj, apiV4, suffix)
	}

	if app := decodeIfDefined(f.Application, &v1alpha1.Application{}, appKind); app != nil {
		b.setupApplication(obj, app, suffix)
	}

	if ctx := decodeIfDefined(f.Context, &v1alpha1.ManagementContext{}, ctxKind); ctx != nil {
		b.setupMgmtContext(obj, ctx, suffix)
	}

	if rsc := decodeIfDefined(f.Resource, &v1alpha1.ApiResource{}, rscKind); rsc != nil {
		b.setupAPIResource(obj, rsc, suffix)
	}

	if ing := decodeIfDefined(f.Ingress, &netV1.Ingress{}, ingKind); ing != nil {
		b.setupIngress(obj, ing, suffix)
	}

	return obj
}

func (b *FSBuilder) setupIngress(obj *Objects, ing **netV1.Ingress, suffix string) {
	obj.Ingress = *ing
	obj.Ingress.Name += suffix
	obj.Ingress.Namespace = constants.Namespace
	if obj.API != nil && isTemplate(obj.API) {
		obj.Ingress.Annotations[keys.IngressTemplateAnnotation] = obj.API.Name
	}
	if obj.APIv4 != nil && isTemplate(obj.APIv4) {
		obj.Ingress.Annotations[keys.IngressTemplateAnnotation] = obj.APIv4.Name
	}

	randomizeIngressRules(obj.Ingress, suffix)
}

func (b *FSBuilder) setupAPIResource(obj *Objects, rsc **v1alpha1.ApiResource, suffix string) {
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
	if obj.APIv4 != nil {
		obj.APIv4.Spec.Resources = []*base.ResourceOrRef{
			{
				Ref: &refs.NamespacedName{
					Name:      obj.Resource.Name,
					Namespace: constants.Namespace,
				},
			},
		}
	}
}

func (b *FSBuilder) setupMgmtContext(obj *Objects, ctx **v1alpha1.ManagementContext, suffix string) {
	obj.Context = *ctx
	obj.Context.Name += suffix
	obj.Context.Namespace = constants.Namespace
	if obj.API != nil {
		obj.API.Spec.Context = obj.Context.GetNamespacedName()
	}
	if obj.APIv4 != nil {
		obj.APIv4.Spec.Context = obj.Context.GetNamespacedName()
	}
	if obj.Application != nil {
		obj.Application.Spec.Context = obj.Context.GetNamespacedName()
	}
}

func (b *FSBuilder) setupApplication(obj *Objects, app **v1alpha1.Application, suffix string) {
	obj.Application = *app
	obj.Application.Name += suffix
	obj.Application.Namespace = constants.Namespace
}

func (b *FSBuilder) setupAPIv4(obj *Objects, apiV4 **v1alpha1.ApiV4Definition, suffix string) {
	obj.APIv4 = *apiV4
	obj.APIv4.Name += suffix
	obj.APIv4.Namespace = constants.Namespace
	obj.APIv4.Spec.Name = obj.APIv4.Name

	randomizeAPIv4Paths(obj.APIv4, suffix)
}

func (b *FSBuilder) setupAPI(obj *Objects, api **v1alpha1.ApiDefinition, suffix string) {
	obj.API = *api
	obj.API.Name += suffix
	obj.API.Namespace = constants.Namespace
	obj.API.Spec.Name = obj.API.Name

	randomizeAPIPaths(obj.API, suffix)
}

func randomizeAPIPaths(api *v1alpha1.ApiDefinition, suffix string) {
	if !isTemplate(api) {
		for _, vh := range api.Spec.Proxy.VirtualHosts {
			vh.Path = "/" + suffix[1:]
		}
	}
}

func randomizeAPIv4Paths(api *v1alpha1.ApiV4Definition, suffix string) {
	if !isTemplate(api) {
		for i, v := range api.Spec.Listeners {
			api.Spec.Listeners[i] = setPath(v, suffix)
		}
	}
}

func setPath(l v4.Listener, suffix string) *v4.GenericListener {
	switch t := l.(type) {
	case *v4.GenericListener:
		return setPath(t.ToListener(), suffix)
	case *v4.HttpListener:
		t.Paths[0].Path = "/" + suffix[1:]
		return v4.ToGenericListener(t)
	case *v4.TCPListener:
		t.Hosts[0] = constants.GatewayHost
		return v4.ToGenericListener(t)
	}
	return nil
}

func randomizeIngressRules(ing *netV1.Ingress, suffix string) {
	for i := range ing.Spec.Rules {
		for j := range ing.Spec.Rules[i].HTTP.Paths {
			ing.Spec.Rules[i].HTTP.Paths[j].Path += suffix
		}
	}
}

func isTemplate(api custom.ApiDefinitionResource) bool {
	return api.GetAnnotations()[keys.IngressTemplateAnnotation] == env.TrueString
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

func (b *FSBuilder) WithAPIv4(file string) *FSBuilder {
	b.files.APIv4 = file
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
