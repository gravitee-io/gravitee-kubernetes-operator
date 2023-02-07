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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (d *Delegate) CreateOrUpdateIngress(ingress *v1.Ingress) (util.OperationResult, error) {
	desired := &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingress.Name,
			Namespace: ingress.Namespace,
		},
	}

	return util.CreateOrUpdate(d.ctx, d.k8s, desired, func() error {
		ingress.Spec.DeepCopyInto(&desired.Spec)

		if !desired.DeletionTimestamp.IsZero() {
			if util.ContainsFinalizer(desired, keys.IngressFinalizer) {
				util.RemoveFinalizer(desired, keys.IngressFinalizer)
			}
			return nil
		}

		if !util.ContainsFinalizer(desired, keys.IngressFinalizer) {
			util.AddFinalizer(desired, keys.IngressFinalizer)
		}

		return nil
	})
}
