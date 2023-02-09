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

package internal

import (
	"fmt"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	"k8s.io/apimachinery/pkg/types"

	"github.com/go-logr/logr"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (d *Delegate) ResolveApiDefinitionTemplate(ingress *v1.Ingress) (*v1alpha1.ApiDefinition, error) {
	var apiDefinition *v1alpha1.ApiDefinition
	if name, ok := ingress.Annotations[keys.IngressTemplateAnnotation]; ok {
		apiDefinition = &v1alpha1.ApiDefinition{}
		if err := d.k8s.Get(d.ctx, types.NamespacedName{Name: name, Namespace: ingress.Namespace}, apiDefinition); err != nil {
			return nil, err
		}
	} else {
		apiDefinition = defaultApiDefinitionTemplate()
	}

	return mergeApiDefinition(d.log, apiDefinition, ingress), nil
}

func defaultApiDefinitionTemplate() *v1alpha1.ApiDefinition {
	return &v1alpha1.ApiDefinition{
		Spec: v1alpha1.ApiDefinitionSpec{
			Api: model.Api{
				Name: "default-keyless",
				Plans: []*model.Plan{
					{
						Name:     "Default keyless plan",
						Security: "KEY_LESS",
						Status:   "PUBLISHED",
					},
				},
			},
		},
	}
}

// MergeApiDefinition
// Transform the ingress as an API Definition as per https://kubernetes.io/docs/concepts/services-networking/ingress/#the-ingress-resource
func mergeApiDefinition(
	log logr.Logger,
	apiDefinition *v1alpha1.ApiDefinition,
	ingress *v1.Ingress,
) *v1alpha1.ApiDefinition {
	api := &v1alpha1.ApiDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingress.Name,
			Namespace: ingress.Namespace,
		},
		Spec: *apiDefinition.Spec.DeepCopy(),
	}

	log.Info("Merge Ingress with API Definition")

	for _, rule := range ingress.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			service := path.Backend.Service

			//TODO: How-to dedal with PathType ?
			api.Spec.Proxy = &model.Proxy{
				VirtualHosts: []*model.VirtualHost{
					{
						Path: path.Path,
					},
				},
				Groups: []*model.EndpointGroup{
					{
						Name: "default",
						Endpoints: []*model.HttpEndpoint{
							{
								Name:   service.Name,
								Target: fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", service.Name, ingress.Namespace, service.Port.Number),
							},
						},
					},
				},
			}
		}
	}
	return api
}
