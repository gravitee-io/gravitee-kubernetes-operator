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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (d *Delegate) createOrUpdateApiDefinition(ingress *v1.Ingress) (util.OperationResult, error) {
	apiDefinition, err := d.resolveApiDefinitionTemplate(ingress)
	if err != nil {
		d.log.Error(err, "ResolveApiDefinition error")
		return util.OperationResultNone, err
	}

	var existingApiDefinition *v1alpha1.ApiDefinition
	existingApiDefinition, err = d.getApiDefinition(types.NamespacedName{Namespace: ingress.Namespace, Name: ingress.Name})
	if errors.IsNotFound(err) {
		return d.createApiDefinition(ingress, apiDefinition)
	}

	if err != nil {
		d.log.Error(err, "unable to create api definition from template")
		return util.OperationResultNone, err
	}

	if !equality.Semantic.DeepEqual(existingApiDefinition, apiDefinition) {
		apiDefinition.Spec.DeepCopyInto(&existingApiDefinition.Spec)
		return d.updateApiDefinition(ingress, existingApiDefinition)
	}

	d.log.Info(
		"No change detected on ApiDefinition. Skipped.",
		"name", apiDefinition.Name,
		"namespace", apiDefinition.Namespace,
	)
	return util.OperationResultNone, nil
}

func (d *Delegate) createApiDefinition(
	ingress *v1.Ingress, apiDefinition *v1alpha1.ApiDefinition,
) (util.OperationResult, error) {
	d.log.Info("Creating ApiDefinition", "name", apiDefinition.Name, "namespace", apiDefinition.Namespace)

	if err := util.SetOwnerReference(ingress, apiDefinition, d.k8s.Scheme()); err != nil {
		return util.OperationResultNone, err
	}

	return util.OperationResultCreated, d.k8s.Create(d.ctx, apiDefinition)
}

func (d *Delegate) updateApiDefinition(
	ingress *v1.Ingress, apiDefinition *v1alpha1.ApiDefinition,
) (util.OperationResult, error) {
	d.log.Info("Updating ApiDefinition", "name", apiDefinition.Name, "namespace", apiDefinition.Namespace)

	err := util.SetOwnerReference(ingress, apiDefinition, d.k8s.Scheme())
	if err != nil {
		return util.OperationResultNone, err
	}

	var existingApiDefinition *v1alpha1.ApiDefinition
	existingApiDefinition, err = d.getApiDefinition(types.NamespacedName{Namespace: ingress.Namespace, Name: ingress.Name})
	if err != nil {
		return util.OperationResultNone, err
	}

	apiDefinition.Spec.DeepCopyInto(&existingApiDefinition.Spec)
	return util.OperationResultUpdated, d.k8s.Update(d.ctx, existingApiDefinition)
}

func (d *Delegate) getApiDefinition(key client.ObjectKey) (*v1alpha1.ApiDefinition, error) {
	api := &v1alpha1.ApiDefinition{}
	err := d.k8s.Get(d.ctx, key, api)
	if err != nil {
		return nil, err
	}
	return api, err
}
