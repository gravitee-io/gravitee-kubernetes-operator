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
	netv1 "k8s.io/api/networking/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1beta1"
)

// This function is applied to all ingresses which are using the ApiDefinition template
// As per Kubernetes Finalizers (https://kubernetes.io/docs/concepts/overview/working-with-objects/finalizers/)
// First return value defines if we should requeue or not.
func (d *Delegate) SyncApiDefinitionTemplate(api *v1beta1.ApiDefinition, ns string) error {
	// We are first looking if the template is in deletion phase, the Kubernetes API marks the object for
	// deletion by populating .metadata.deletionTimestamp
	if !api.DeletionTimestamp.IsZero() {
		return d.delete(api, ns)
	}

	if !util.ContainsFinalizer(api, keys.ApiDefinitionTemplateFinalizer) {
		util.AddFinalizer(api, keys.ApiDefinitionTemplateFinalizer)
		return d.k8s.Update(d.ctx, api)
	}

	return nil
}

func (d *Delegate) delete(apiDefinition *v1beta1.ApiDefinition, namespace string) error {
	if !util.ContainsFinalizer(apiDefinition, keys.ApiDefinitionTemplateFinalizer) {
		return nil
	}

	ingressList := netv1.IngressList{}

	// Retrieves the ingresses from the namespace
	err := d.k8s.List(d.ctx, &ingressList, client.InNamespace(namespace))
	if err != nil && !kerrors.IsNotFound(err) {
		return err
	}

	var ingresses []string

	for i := range ingressList.Items {
		if ingressList.Items[i].GetAnnotations()[keys.IngressTemplateAnnotation] == apiDefinition.Name {
			ingresses = append(ingresses, ingressList.Items[i].GetName())
		}
	}

	// There are existing ingresses which are still relying on this ApiDefinition template, re-schedule deletion
	if len(ingresses) > 0 {
		return fmt.Errorf("can not delete %s %v depends on it", apiDefinition.Name, ingresses)
	}

	util.RemoveFinalizer(apiDefinition, keys.ApiDefinitionTemplateFinalizer)

	return d.k8s.Update(d.ctx, apiDefinition)
}
