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
	"fmt"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	managementapierror "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi/clienterror"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	"k8s.io/client-go/util/retry"
	k8sUtil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type NonRecoverableError struct {
	cause error
}

func (e NonRecoverableError) Error() string {
	return fmt.Sprintf("Non recoverable error: %s", e.cause.Error())
}

func IsRecoverableError(err error) bool {
	return !errors.As(err, new(NonRecoverableError))
}

func (d *Delegate) UpdateStatusAndReturnError(apiDefinition *gio.ApiDefinition, reconcileErr error) error {
	reconcileErr = wrapError(reconcileErr)

	processingStatus := gio.ProcessingStatusReconciling
	if !IsRecoverableError(reconcileErr) {
		processingStatus = gio.ProcessingStatusFailed

		// Remove finalizer when API definition is failed. To allow the user to remove it.
		k8sUtil.RemoveFinalizer(apiDefinition, keys.ApiDefinitionDeletionFinalizer)
	}

	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		apiDefinition.Status.ProcessingStatus = processingStatus
		apiDefinition.Status.ObservedGeneration = apiDefinition.ObjectMeta.Generation
		return d.k8sClient.Status().Update(d.ctx, apiDefinition.DeepCopy())
	})
	if err != nil {
		d.log.Info("Unexpected error while updating API definition status.", "err", err)
		return err
	}
	return reconcileErr
}

// Wraps the error in a NonRecoverableError if it's not recoverable.
func wrapError(err error) error {
	switch {
	case managementapierror.IsUnauthorized(err):
		return NonRecoverableError{cause: err}
	default:
		return err
	}
}
