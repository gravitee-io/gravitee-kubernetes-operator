/*
Copyright 2022 DAVID BRASSELY.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package apidefinition

import (
	"context"
	"time"

	"k8s.io/client-go/tools/record"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/template"

	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"

	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/apidefinition/internal"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const requeueAfterTime = time.Second * 5

func Reconcile(
	ctx context.Context,
	apiDefinition core.ApiDefinitionObject,
	r record.EventRecorder,
) (ctrl.Result, error) {
	events := event.NewRecorder(r)

	if apiDefinition.GetAnnotations()[core.IngressTemplateAnnotation] == "true" {
		return reconcileApiTemplate(ctx, apiDefinition)
	}

	return reconcileApiDefinition(ctx, apiDefinition, events)
}

func reconcileApiTemplate(ctx context.Context, apiDefinition core.ApiDefinitionObject) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	status := apiDefinition.GetStatus()
	_, err := util.CreateOrUpdate(ctx, k8s.GetClient(), apiDefinition, func() error {
		dc, _ := apiDefinition.DeepCopyObject().(core.ApiDefinitionObject)
		if err := template.Compile(ctx, dc); err != nil {
			return err
		}

		if err := internal.SyncApiDefinitionTemplate(ctx, dc, apiDefinition.GetNamespace()); err != nil {
			return err
		}

		apiDefinition.SetFinalizers(dc.GetFinalizers())
		status = dc.GetStatus()

		return nil
	})

	if err != nil {
		logger.Error(err, "Failed to sync API definition template")
		return ctrl.Result{RequeueAfter: requeueAfterTime}, err
	}

	logger.Info("template synced successfully.", "template:", apiDefinition.GetName())
	if err := status.DeepCopyTo(apiDefinition); err != nil {
		return ctrl.Result{}, err
	}

	if err := internal.UpdateStatusSuccess(ctx, apiDefinition); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
func reconcileApiDefinition(
	ctx context.Context,
	apiDefinition core.ApiDefinitionObject,
	events *event.Recorder,
) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	status := apiDefinition.GetStatus()
	_, reconcileErr := util.CreateOrUpdate(ctx, k8s.GetClient(), apiDefinition, func() error {
		util.AddFinalizer(apiDefinition, core.ApiDefinitionFinalizer)
		k8s.AddAnnotation(apiDefinition, core.LastSpecHashAnnotation, apiDefinition.GetSpec().Hash())

		dc, _ := apiDefinition.DeepCopyObject().(core.ApiDefinitionObject)

		if err := template.Compile(ctx, dc); err != nil {
			status.SetProcessingStatus(core.ProcessingStatusFailed)
			return err
		}

		var err error
		if !apiDefinition.GetDeletionTimestamp().IsZero() {
			err = events.Record(event.Delete, apiDefinition, func() error {
				return internal.Delete(ctx, dc)
			})
		} else {
			err = events.Record(event.Update, apiDefinition, func() error {
				return internal.CreateOrUpdate(ctx, dc)
			})
		}

		if err != nil {
			return err
		}

		apiDefinition.SetFinalizers(dc.GetFinalizers())
		status = dc.GetStatus()
		return err
	})

	if err := status.DeepCopyTo(apiDefinition); err != nil {
		return ctrl.Result{}, err
	}

	if reconcileErr == nil {
		logger.Info("API definition has been reconciled")
		return ctrl.Result{}, internal.UpdateStatusSuccess(ctx, apiDefinition)
	}

	if err := internal.UpdateStatusFailure(ctx, apiDefinition); err != nil {
		return ctrl.Result{}, err
	}

	if errors.IsRecoverable(reconcileErr) {
		logger.Error(reconcileErr, "Requeuing reconcile")
		return ctrl.Result{RequeueAfter: requeueAfterTime}, reconcileErr
	}

	logger.Error(reconcileErr, "Aborting reconcile")
	return ctrl.Result{}, nil
}
