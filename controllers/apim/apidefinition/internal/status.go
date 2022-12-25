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
	"errors"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	kErrors "k8s.io/apimachinery/pkg/util/errors"
)

func (d *Delegate) UpdateStatusSuccess(api *gio.ApiDefinition) error {
	if api.IsBeingDeleted() {
		return nil
	}
	api.Status.ObservedGeneration = api.ObjectMeta.Generation
	return d.k8s.Status().Update(d.ctx, api)
}

func (d *Delegate) UpdateStatusFailure(api *gio.ApiDefinition, err error) error {
	api.Status.ObservedGeneration = api.ObjectMeta.Generation
	updateStatusContexts(api, err)
	return d.k8s.Status().Update(d.ctx, api)
}

func updateStatusContexts(api *gio.ApiDefinition, err error) {
	//nolint:errorlint // cannot use As on interface
	agg, ok := err.(kErrors.Aggregate)

	if !ok {
		updateStatusContext(api, err)
		return
	}

	agg = kErrors.Flatten(agg)

	for _, err := range agg.Errors() {
		updateStatusContext(api, err)
	}
}

func updateStatusContext(api *gio.ApiDefinition, err error) {
	contextError := &ContextError{}
	if !errors.As(err, contextError) {
		return
	}

	context, ok := api.Status.Contexts[contextError.Context]

	if !ok {
		context = gio.StatusContext{}
	}

	context.Status = gio.ProcessingStatusFailed
	api.Status.Contexts[contextError.Context] = context
}
