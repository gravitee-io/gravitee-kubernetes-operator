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
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (d *Delegate) CreateOrUpdate(desired *v1.Ingress) error {
	if err := d.addFinalizer(desired); err != nil {
		d.log.Error(err, "An error occurs while adding finalizer to the Ingress", "Ingress", desired)
		return err
	}

	if err := d.updateIngressTLSReference(desired); err != nil {
		d.log.Error(err, "An error occurred while updating the PEM registry")
		return err
	}

	operation, apiDefinitionError := d.createOrUpdateApiDefinition(desired)
	if apiDefinitionError != nil {
		d.log.Error(
			apiDefinitionError,
			"An error occurs while creating or updating the ApiDefinition",
			"Operation", operation,
		)
		return apiDefinitionError
	}

	return nil
}

func (d *Delegate) addFinalizer(desired *v1.Ingress) error {
	ingress := &v1.Ingress{}
	if err := d.k8s.Get(
		d.ctx, types.NamespacedName{Namespace: desired.Namespace, Name: desired.Name}, ingress,
	); err != nil {
		return err
	}

	if !util.ContainsFinalizer(ingress, keys.IngressFinalizer) {
		util.AddFinalizer(ingress, keys.IngressFinalizer)
	}

	return d.k8s.Update(d.ctx, ingress)
}
