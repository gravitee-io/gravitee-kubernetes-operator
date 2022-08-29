package internal

import (
	"errors"
	"fmt"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	managementapierror "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi/clienterror"
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

	processingStatus := gio.ProcessingStatusFailed
	if IsRecoverableError(reconcileErr) {
		processingStatus = gio.ProcessingStatusReconciling
	}

	apiDefinition.Status.ProcessingStatus = processingStatus

	// Updated succeed, update Generation & Status
	apiDefinition.Status.Generation = apiDefinition.ObjectMeta.Generation
	err := d.k8sClient.Status().Update(d.ctx, apiDefinition.DeepCopy())
	if err != nil {
		d.log.Info("Unexpected error while updating API definition status.", "err", err)
		return err
	}
	return reconcileErr
}

// Wraps the error in a NonRecoverableError if it's not recoverable.
func wrapError(err error) error {
	switch {
	case errors.As(err, new(managementapierror.ApiUnauthorizedError)):
		return NonRecoverableError{cause: err}
	default:
		return err
	}
}
