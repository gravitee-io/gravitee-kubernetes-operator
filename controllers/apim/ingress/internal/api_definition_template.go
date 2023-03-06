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
	"net/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/ingress/internal/mapper"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	"k8s.io/apimachinery/pkg/types"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	coreV1 "k8s.io/api/core/v1"
	netV1 "k8s.io/api/networking/v1"
)

func (d *Delegate) resolveApiDefinitionTemplate(ingress *netV1.Ingress) (*v1alpha1.ApiDefinition, error) {
	var apiDefinition *v1alpha1.ApiDefinition

	if name, ok := ingress.Annotations[keys.IngressTemplateAnnotation]; ok {
		apiDefinition = &v1alpha1.ApiDefinition{}
		if err := d.k8s.Get(
			d.ctx, types.NamespacedName{Name: name, Namespace: ingress.Namespace}, apiDefinition,
		); err != nil {
			return nil, err
		}
	} else {
		apiDefinition = defaultApiDefinitionTemplate()
	}

	return mapper.New(d.getMapperOpts()).Map(apiDefinition, ingress), nil
}

func (d *Delegate) getMapperOpts() mapper.Opts {
	opts := mapper.NewOpts()
	d.setNotFoundTemplate(&opts)
	return opts
}

func (d *Delegate) setNotFoundTemplate(opts *mapper.Opts) {
	ns, name := env.Config.CMTemplate404NS, env.Config.CMTemplate404Name

	if name == "" {
		return
	}

	cm := coreV1.ConfigMap{}
	if err := d.k8s.Get(d.ctx, types.NamespacedName{Namespace: ns, Name: name}, &cm); err != nil {
		d.log.Error(err, "unable to access config map, using default HTTP not found template")
		return
	}

	if err := checkData(cm.Data); err != nil {
		d.log.Error(err, "missing key in config map, using default HTTP not found template")
		return
	}

	opts.Templates[http.StatusNotFound] = mapper.ResponseTemplate{
		Content:     cm.Data["content"],
		ContentType: cm.Data["contentType"],
	}
}

func checkData(template map[string]string) error {
	if _, ok := template["content"]; !ok {
		return fmt.Errorf("missing content in template")
	}

	if _, ok := template["contentType"]; !ok {
		return fmt.Errorf("missing contentType in template")
	}

	return nil
}

func defaultApiDefinitionTemplate() *v1alpha1.ApiDefinition {
	return &v1alpha1.ApiDefinition{
		Spec: v1alpha1.ApiDefinitionSpec{
			Api: model.Api{
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
