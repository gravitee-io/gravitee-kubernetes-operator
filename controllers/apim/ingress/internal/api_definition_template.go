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
	"context"
	"fmt"
	"net/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/ingress/internal/mapper"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"k8s.io/apimachinery/pkg/types"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	coreV1 "k8s.io/api/core/v1"
	netV1 "k8s.io/api/networking/v1"
)

func resolveApiDefinitionTemplate(
	ctx context.Context,
	ingress *netV1.Ingress,
) (*v1alpha1.ApiDefinition, error) {
	var apiDefinition *v1alpha1.ApiDefinition

	if name, ok := ingress.Annotations[core.IngressTemplateAnnotation]; ok {
		apiDefinition = &v1alpha1.ApiDefinition{}
		cli := k8s.GetClient()
		if err := cli.Get(
			ctx, types.NamespacedName{Name: name, Namespace: ingress.Namespace}, apiDefinition,
		); err != nil {
			return nil, err
		}
	} else {
		apiDefinition = defaultApiDefinitionTemplate()
	}

	return mapper.New(getMapperOpts(ctx, ingress)).Map(apiDefinition, ingress), nil
}

func getMapperOpts(ctx context.Context, ingress *netV1.Ingress) mapper.Opts {
	opts := mapper.NewOpts()
	setNotFoundTemplate(ctx, ingress, &opts)
	return opts
}

func setNotFoundTemplate(ctx context.Context, ingress *netV1.Ingress, opts *mapper.Opts) {
	ns, name := env.Config.CMTemplate404NS, env.Config.CMTemplate404Name

	if name == "" {
		return
	}

	cm := coreV1.ConfigMap{}
	cli := k8s.GetClient()
	if err := cli.Get(ctx, types.NamespacedName{Namespace: ns, Name: name}, &cm); err != nil {
		log.Error(
			ctx,
			err,
			"unable to access config map, using default HTTP not found template",
			log.KeyValues(ingress)...,
		)
		return
	}

	if err := checkData(cm.Data); err != nil {
		log.Error(
			ctx,
			err,
			"missing key in config map, using default HTTP not found template",
			log.KeyValues(ingress)...,
		)
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
		Spec: v1alpha1.ApiDefinitionV2Spec{
			Api: v2.Api{
				Plans: []*v2.Plan{
					v2.NewPlan(
						base.NewPlan().
							WithStatus(base.PublishedPlanStatus),
					).WithDescription("Default ingress keyless plan").
						WithSecurity("KEY_LESS").WithName("Key Less"),
				},
				ApiBase: &base.ApiBase{
					Version: "1.0.0",
				},
				Description: "This API was generated on behalf of an Kubernetes ingress resource",
			},
			IsLocal: true,
		},
	}
}
