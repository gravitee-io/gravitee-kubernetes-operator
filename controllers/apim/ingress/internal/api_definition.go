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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (d *Delegate) createOrUpdateApiDefinition(ctx context.Context, ingress *v1.Ingress) (util.OperationResult, error) {
	apiDefinition, err := d.resolveApiDefinitionTemplate(ctx, ingress)
	if err != nil {
		log.FromContext(ctx).Error(err, "ResolveApiDefinition error")
		return util.OperationResultNone, err
	}

	var existingApiDefinition *v1alpha1.ApiDefinition
	nsm := types.NamespacedName{Namespace: ingress.Namespace, Name: ingress.Name}
	existingApiDefinition, err = d.getApiDefinition(ctx, nsm)
	if errors.IsNotFound(err) {
		return d.createApiDefinition(ctx, ingress, apiDefinition)
	}

	if err != nil {
		log.FromContext(ctx).Error(err, "unable to create api definition from template")
		return util.OperationResultNone, err
	}

	if !equality.Semantic.DeepEqual(existingApiDefinition, apiDefinition) {
		apiDefinition.Spec.DeepCopyInto(&existingApiDefinition.Spec)
		return d.updateApiDefinition(ctx, ingress, existingApiDefinition)
	}

	log.FromContext(ctx).Info(
		"No change detected on ApiDefinition. Skipped.",
		"name", apiDefinition.Name,
		"namespace", apiDefinition.Namespace,
	)
	return util.OperationResultNone, nil
}

func (d *Delegate) createApiDefinition(
	ctx context.Context,
	ingress *v1.Ingress, apiDefinition *v1alpha1.ApiDefinition,
) (util.OperationResult, error) {
	log.FromContext(ctx).Info("Creating ApiDefinition", "name", apiDefinition.Name, "namespace", apiDefinition.Namespace)

	cli := k8s.GetClient()
	if err := util.SetOwnerReference(ingress, apiDefinition, cli.Scheme()); err != nil {
		return util.OperationResultNone, err
	}

	return util.OperationResultCreated, cli.Create(ctx, apiDefinition)
}

func (d *Delegate) updateApiDefinition(
	ctx context.Context,
	ingress *v1.Ingress, apiDefinition *v1alpha1.ApiDefinition,
) (util.OperationResult, error) {
	log.FromContext(ctx).Info("Updating ApiDefinition", "name", apiDefinition.Name, "namespace", apiDefinition.Namespace)
	cli := k8s.GetClient()
	err := util.SetOwnerReference(ingress, apiDefinition, cli.Scheme())
	if err != nil {
		return util.OperationResultNone, err
	}

	var existingApiDefinition *v1alpha1.ApiDefinition
	nsm := types.NamespacedName{Namespace: ingress.Namespace, Name: ingress.Name}
	existingApiDefinition, err = d.getApiDefinition(ctx, nsm)
	if err != nil {
		return util.OperationResultNone, err
	}

	apiDefinition.Spec.DeepCopyInto(&existingApiDefinition.Spec)
	return util.OperationResultUpdated, cli.Update(ctx, existingApiDefinition)
}

func (d *Delegate) getApiDefinition(ctx context.Context, key client.ObjectKey) (*v1alpha1.ApiDefinition, error) {
	api := &v1alpha1.ApiDefinition{}
	cli := k8s.GetClient()
	err := cli.Get(ctx, key, api)
	if err != nil {
		return nil, err
	}
	return api, err
}
