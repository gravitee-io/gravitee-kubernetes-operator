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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func createOrUpdateApiDefinition(ctx context.Context, ingress *v1.Ingress) (util.OperationResult, error) {
	apiDefinition, err := resolveApiDefinitionTemplate(ctx, ingress)
	if err != nil {
		log.Error(
			ctx,
			err,
			"Unable to resolve API definition template",
			log.KeyValues(
				ingress,
				"ingress-template-name", ingress.Annotations[core.IngressTemplateAnnotation],
			)...,
		)
		return util.OperationResultNone, err
	}

	var existingApiDefinition *v1alpha1.ApiDefinition
	nsm := types.NamespacedName{Namespace: ingress.Namespace, Name: ingress.Name}
	existingApiDefinition, err = getApiDefinition(ctx, nsm)
	if errors.IsNotFound(err) {
		return createApiDefinition(ctx, ingress, apiDefinition)
	}

	if err != nil {
		log.Error(ctx, err, "unable to create api definition from template", log.KeyValues(apiDefinition)...)
		return util.OperationResultNone, err
	}

	if !equality.Semantic.DeepEqual(existingApiDefinition, apiDefinition) {
		apiDefinition.Spec.DeepCopyInto(&existingApiDefinition.Spec)
		return updateApiDefinition(ctx, ingress, existingApiDefinition)
	}

	log.Debug(
		ctx,
		"Skipping API definition update for ingress as no changes were detected",
		log.KeyValues(ingress)...,
	)
	return util.OperationResultNone, nil
}

func createApiDefinition(
	ctx context.Context,
	ingress *v1.Ingress, apiDefinition *v1alpha1.ApiDefinition,
) (util.OperationResult, error) {
	log.Debug(
		ctx,
		"Creating API definition for ingress",
		log.KeyValues(
			ingress,
			"ingress-api-name", apiDefinition.GetName(),
			"ingress-api-namespace", apiDefinition.GetNamespace(),
		)...,
	)

	cli := k8s.GetClient()
	if err := util.SetOwnerReference(ingress, apiDefinition, cli.Scheme()); err != nil {
		return util.OperationResultNone, err
	}

	return util.OperationResultCreated, cli.Create(ctx, apiDefinition)
}

func updateApiDefinition(
	ctx context.Context,
	ingress *v1.Ingress, apiDefinition *v1alpha1.ApiDefinition,
) (util.OperationResult, error) {
	log.Debug(
		ctx,
		"Updating API definition for ingress",
		log.KeyValues(ingress)...,
	)
	cli := k8s.GetClient()
	err := util.SetOwnerReference(ingress, apiDefinition, cli.Scheme())
	if err != nil {
		return util.OperationResultNone, err
	}

	var existingApiDefinition *v1alpha1.ApiDefinition
	nsm := types.NamespacedName{Namespace: ingress.Namespace, Name: ingress.Name}
	existingApiDefinition, err = getApiDefinition(ctx, nsm)
	if err != nil {
		return util.OperationResultNone, err
	}

	apiDefinition.Spec.DeepCopyInto(&existingApiDefinition.Spec)
	return util.OperationResultUpdated, cli.Update(ctx, existingApiDefinition)
}

func getApiDefinition(ctx context.Context, key client.ObjectKey) (*v1alpha1.ApiDefinition, error) {
	api := &v1alpha1.ApiDefinition{}
	cli := k8s.GetClient()
	err := cli.Get(ctx, key, api)
	if err != nil {
		return nil, err
	}
	return api, err
}
