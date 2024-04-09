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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (d *Delegate) UpdateStatusSuccess(api client.Object) error {
	if !api.GetDeletionTimestamp().IsZero() {
		return nil
	}

	switch t := api.(type) {
	case *v1alpha1.ApiDefinition:
		t.Status.ObservedGeneration = t.ObjectMeta.Generation
		t.Status.Status = v1alpha1.ProcessingStatusCompleted
		return d.k8s.Status().Update(d.ctx, t)
	case *v1alpha1.ApiDefinitionV4:
		t.Status.ObservedGeneration = t.ObjectMeta.Generation
		t.Status.Status = v1alpha1.ProcessingStatusCompleted
		return d.k8s.Status().Update(d.ctx, t)
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

func (d *Delegate) UpdateStatusFailure(api client.Object) error {
	switch t := api.(type) {
	case *v1alpha1.ApiDefinition:
		t.Status.Status = v1alpha1.ProcessingStatusFailed
		return d.k8s.Status().Update(d.ctx, t)
	case *v1alpha1.ApiDefinitionV4:
		t.Status.Status = v1alpha1.ProcessingStatusFailed
		return d.k8s.Status().Update(d.ctx, t)
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}
