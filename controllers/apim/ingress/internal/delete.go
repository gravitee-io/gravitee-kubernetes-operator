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
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (d *Delegate) Delete(ingress *v1.Ingress) error {
	if err := d.deleteIngressTLSReference(ingress); err != nil {
		d.log.Error(err, "An error occurred while updating the TLS secrets")
		return err
	}

	if util.ContainsFinalizer(ingress, keys.IngressFinalizer) {
		util.RemoveFinalizer(ingress, keys.IngressFinalizer)
	}

	// because we set SetOwnerReference during phase, we don't need to delete the
	// api definition manually because it will be automatically deleted once the
	// parent is deleted so we just need to remove the finalizer
	return d.k8s.Update(d.ctx, ingress)
}
