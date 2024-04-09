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
	"fmt"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"k8s.io/client-go/tools/record"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"

	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/apidefinition/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const requeueAfterTime = time.Second * 5

func Reconcile(ctx context.Context,
	apiDefinition client.Object,
	c client.Client,
	r record.EventRecorder) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	d := internal.NewDelegate(ctx, c, logger)
	events := event.NewRecorder(r)

	if apiDefinition.GetAnnotations()[keys.IngressTemplateAnnotation] == "true" {
		logger.Info("syncing template", "template", apiDefinition.GetName())
		if err := d.ResolveTemplate(apiDefinition); err != nil {
			return ctrl.Result{}, err
		}
		if err := d.SyncApiDefinitionTemplate(apiDefinition, apiDefinition.GetNamespace()); err != nil {
			logger.Error(err, "Failed to sync API definition template")
			return ctrl.Result{RequeueAfter: requeueAfterTime}, err
		}

		logger.Info("template synced successfully.", "template:", apiDefinition.GetName())
		return ctrl.Result{}, nil
	}

	switch t := apiDefinition.(type) {
	case *v1alpha1.ApiDefinition:
		return reconcileApiDefinition(ctx, t, t.DeepCopy(), &t.Spec, &t.Status, t.Spec.Context, c, d, events)
	case *v1alpha1.ApiDefinitionV4:
		return reconcileApiDefinition(ctx, t, t.DeepCopy(), &t.Spec, &t.Status, t.Spec.Context, c, d, events)
	default:
		return ctrl.Result{}, fmt.Errorf("unknown type %T", t)
	}
}

func reconcileApiDefinition(ctx context.Context,
	apiDefinition, dc client.Object,
	spec any, status v1alpha1.Status, mCtx *refs.NamespacedName,
	c client.Client, d *internal.Delegate, events *event.Recorder) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	_, reconcileErr := util.CreateOrUpdate(ctx, c, dc, func() error {
		util.AddFinalizer(apiDefinition, keys.ApiDefinitionFinalizer)
		k8s.AddAnnotation(apiDefinition, keys.LastSpecHash, hash.Calculate(spec))

		if err := d.ResolveTemplate(apiDefinition); err != nil {
			status.SetProcessingStatus(v1alpha1.ProcessingStatusFailed)
			return err
		}

		if mCtx != nil {
			if err := d.ResolveContext(mCtx); err != nil {
				status.SetProcessingStatus(v1alpha1.ProcessingStatusFailed)
				logger.Info("Unable to resolve context, no attempt will be made to sync with APIM")
				return err
			}
		}

		var err error
		if !apiDefinition.GetDeletionTimestamp().IsZero() {
			err = events.Record(event.Delete, apiDefinition, func() error {
				return d.Delete(apiDefinition)
			})
		} else {
			err = events.Record(event.Update, apiDefinition, func() error {
				return d.CreateOrUpdate(apiDefinition)
			})
		}

		if err != nil {
			return err
		}
		err = status.DeepCopyFrom(apiDefinition)
		dc.SetFinalizers(apiDefinition.GetFinalizers())
		dc.SetAnnotations(apiDefinition.GetAnnotations())
		return err
	})

	if err := status.DeepCopyTo(dc); err != nil {
		return ctrl.Result{}, err
	}

	if reconcileErr == nil {
		logger.Info("API definition has been reconciled")
		return ctrl.Result{}, d.UpdateStatusSuccess(dc)
	}

	if err := d.UpdateStatusFailure(dc); err != nil {
		return ctrl.Result{}, err
	}

	if apim.IsRecoverable(reconcileErr) {
		logger.Error(reconcileErr, "Requeuing reconcile")
		return ctrl.Result{RequeueAfter: requeueAfterTime}, reconcileErr
	}

	logger.Error(reconcileErr, "Aborting reconcile")
	return ctrl.Result{}, nil
}
