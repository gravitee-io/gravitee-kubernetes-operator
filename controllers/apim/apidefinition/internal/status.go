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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

func (d *Delegate) UpdateStatusSuccess(api v1alpha1.CRD) error {
	if !api.GetDeletionTimestamp().IsZero() {
		return nil
	}

	api.GetStatus().SetProcessingStatus(v1alpha1.ProcessingStatusCompleted)
	api.GetStatus().SetObservedGeneration(api.GetGeneration())
	return k8s.GetClient().Status().Update(d.ctx, api)
}

func (d *Delegate) UpdateStatusFailure(api v1alpha1.CRD) error {
	api.GetStatus().SetProcessingStatus(v1alpha1.ProcessingStatusFailed)
	return k8s.GetClient().Status().Update(d.ctx, api)
}
