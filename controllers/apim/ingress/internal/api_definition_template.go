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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
)

func DefaultApiDefinitionTemplate() *v1alpha1.ApiDefinition {
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

func (d *Delegate) ResolveApiDefinition(ingress *v1.Ingress, namespace string) (*v1alpha1.ApiDefinition, error) {
	var apiDefinition *v1alpha1.ApiDefinition
	name, hasIngressTemplateAnnotation := ingress.Annotations[keys.IngressTemplateAnnotation]

	if hasIngressTemplateAnnotation {
		apiDefinition = &v1alpha1.ApiDefinition{}
		err := d.k8s.Get(d.ctx, types.NamespacedName{Name: name, Namespace: namespace}, apiDefinition)
		if err != nil {
			return nil, err
		}
	} else {
		apiDefinition = DefaultApiDefinitionTemplate()
	}

	return MergeApiDefinition(d.log, apiDefinition, ingress), nil
}
