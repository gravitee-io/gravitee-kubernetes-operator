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

	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/kube/custom"

	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/apidefinition/internal"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const requeueAfterTime = time.Second * 5

func Reconcile(
	ctx context.Context,
	apiDefinition custom.ApiDefinition,
	r record.EventRecorder,
) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	events := event.NewRecorder(r)

	if apiDefinition.GetAnnotations()[keys.IngressTemplateAnnotation] == "true" {
		logger.Info("syncing template", "template", apiDefinition.GetName())
		if err := template.Compile(ctx, apiDefinition); err != nil {
			return ctrl.Result{}, err
		}
		if err := internal.SyncApiDefinitionTemplate(ctx, apiDefinition, apiDefinition.GetNamespace()); err != nil {
			logger.Error(err, "Failed to sync API definition template")
			return ctrl.Result{RequeueAfter: requeueAfterTime}, err
		}

		logger.Info("template synced successfully.", "template:", apiDefinition.GetName())
		return ctrl.Result{}, nil
	}

	return reconcileApiDefinition(ctx, apiDefinition, events)
}

func reconcileApiDefinition(
	ctx context.Context,
	apiDefinition custom.ApiDefinition,
	events *event.Recorder,
) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	dc := apiDefinition.DeepCopyResource()
	status := apiDefinition.GetStatus()
	_, reconcileErr := util.CreateOrUpdate(ctx, k8s.GetClient(), dc, func() error {
		util.AddFinalizer(apiDefinition, keys.ApiDefinitionFinalizer)
		k8s.AddAnnotation(apiDefinition, keys.LastSpecHash, apiDefinition.GetSpec().Hash())

		if err := template.Compile(ctx, apiDefinition); err != nil {
			status.SetProcessingStatus(custom.ProcessingStatusFailed)
			return err
		}

		var err error
		if !apiDefinition.GetDeletionTimestamp().IsZero() {
			err = events.Record(event.Delete, apiDefinition, func() error {
				return internal.Delete(ctx, apiDefinition)
			})
		} else {
			err = events.Record(event.Update, apiDefinition, func() error {
				return internal.CreateOrUpdate(ctx, apiDefinition)
			})
		}

		if err != nil {
			return err
		}

		dc.SetFinalizers(apiDefinition.GetFinalizers())
		dc.SetAnnotations(apiDefinition.GetAnnotations())
		return err
	})

	if err := status.DeepCopyTo(dc); err != nil {
		return ctrl.Result{}, err
	}

	if reconcileErr == nil {
		logger.Info("API definition has been reconciled")
		return ctrl.Result{}, internal.UpdateStatusSuccess(ctx, dc)
	}

	if err := internal.UpdateStatusFailure(ctx, dc); err != nil {
		return ctrl.Result{}, err
	}

	if errors.IsRecoverable(reconcileErr) {
		logger.Error(reconcileErr, "Requeuing reconcile")
		return ctrl.Result{RequeueAfter: requeueAfterTime}, reconcileErr
	}

	logger.Error(reconcileErr, "Aborting reconcile")
	return ctrl.Result{}, nil
}
