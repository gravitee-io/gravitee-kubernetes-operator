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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

func (d *Delegate) resolveResources(resources []*base.ResourceOrRef) error {
	if resources == nil {
		return nil
	}

	for _, resource := range resources {
		if err := d.resolveIfRef(resource); err != nil {
			return err
		}
	}

	return nil
}

func (d *Delegate) resolveIfRef(resourceOrRef *base.ResourceOrRef) error {
	if !resourceOrRef.IsRef() {
		return nil
	}

	namespacedName := resourceOrRef.Ref.ToK8sType()
	resource := new(v1alpha1.ApiResource)

	d.log.Info("Looking for api resource from", "namespace", namespacedName.Namespace, "name", namespacedName.Name)

	if err := d.k8s.Get(d.ctx, namespacedName, resource); err != nil {
		return err
	}

	resourceOrRef.Resource = resource.Spec.Resource

	return nil
}
