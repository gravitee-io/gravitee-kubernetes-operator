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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	"k8s.io/apimachinery/pkg/types"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1beta1"
	coreV1 "k8s.io/api/core/v1"
	netV1 "k8s.io/api/networking/v1"
)

func (d *Delegate) resolveApiDefinitionTemplate(ingress *netV1.Ingress) (*v1beta1.ApiDefinition, error) {
	var apiDefinition *v1beta1.ApiDefinition

	if name, ok := ingress.Annotations[keys.IngressTemplateAnnotation]; ok {
		apiDefinition = &v1beta1.ApiDefinition{}
		if err := d.k8s.Get(
			d.ctx, types.NamespacedName{Name: name, Namespace: ingress.Namespace}, apiDefinition,
		); err != nil {
			return nil, err
		}
	} else {
		apiDefinition = defaultApiDefinitionTemplate(ingress)
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
		log.Error(d.ctx, err, "Unable to access config map, using default HTTP not found template")
		return
	}

	if err := checkData(cm.Data); err != nil {
		log.Error(d.ctx, err, "missing key in config map, using default HTTP not found template")
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

func defaultApiDefinitionTemplate(ingress *netV1.Ingress) *v1beta1.ApiDefinition {
	return &v1beta1.ApiDefinition{
		Spec: v1beta1.ApiDefinitionSpec{
			Api: v4.Api{
				Type: v4.ProxyType,
				Plans: map[string]*v4.Plan{
					"Default Keyless (UNSECURED)": v4.NewPlan(
						base.NewPlan("Keyless (UNSECURED)", "Default ingress plan").
							WithStatus(base.PublishedPlanStatus),
					).WithSecurity(v4.NewPlanSecurity("KEY_LESS")),
				},
				FlowExecution: v4.DefaultFlowExecution(),
				ApiBase: &base.ApiBase{
					Version:     "v1beta1",
					Description: generateDescription(ingress),
				},
			},
		},
	}
}

func generateDescription(ingress *netV1.Ingress) string {
	return fmt.Sprintf("This API has been generated to handle Ingress %s", ingress.Name)
}
