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

package oidcgroupmapping

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/idpgroupmapping/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/template"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type Reconciler struct {
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	oidcgroupmapping := &v1alpha1.IDPGroupMapping{}
	if err := k8s.GetClient().Get(ctx, req.NamespacedName, oidcgroupmapping); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	oidcgroupmapping.SetNamespace(req.Namespace)

	events := event.NewRecorder(r.Recorder)

	oidcgroupmapping.SetConditions([]metav1.Condition{})
	dc := oidcgroupmapping.DeepCopy()

	_, err := util.CreateOrUpdate(ctx, k8s.GetClient(), oidcgroupmapping, func() error {
		util.AddFinalizer(oidcgroupmapping, core.IDPGroupMappingFinalizer)
		k8s.AddAnnotation(oidcgroupmapping, core.LastSpecHashAnnotation, hash.Calculate(&dc.Spec))

		if err := template.Compile(ctx, dc, true); err != nil {
			oidcgroupmapping.Status.ProcessingStatus = core.ProcessingStatusFailed
			return err
		}

		var err error
		if oidcgroupmapping.IsBeingDeleted() {
			err = events.Record(event.Delete, oidcgroupmapping, func() error {
				if err := internal.Delete(ctx, dc); err != nil {
					return err
				}
				util.RemoveFinalizer(oidcgroupmapping, core.GroupFinalizer)
				return nil
			})
		} else {
			err = events.Record(event.Update, oidcgroupmapping, func() error {
				return internal.CreateOrUpdate(ctx, dc)
			})
		}

		return err
	})

	if err := dc.GetStatus().DeepCopyTo(oidcgroupmapping); err != nil {
		return ctrl.Result{}, err
	}

	if err == nil {
		log.InfoEndReconcile(ctx, oidcgroupmapping)
		return ctrl.Result{}, internal.UpdateStatusSuccess(ctx, oidcgroupmapping)
	}

	// An error occurred during the reconcile
	if err := internal.UpdateStatusFailure(ctx, oidcgroupmapping, err); err != nil {
		return ctrl.Result{}, err
	}

	if errors.IsRecoverable(err) {
		log.ErrorRequeuingReconcile(ctx, err, oidcgroupmapping)
		return ctrl.Result{}, err
	}

	log.ErrorAbortingReconcile(ctx, err, oidcgroupmapping)
	return ctrl.Result{}, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.IDPGroupMapping{}).
		WithEventFilter(predicate.LastSpecHashPredicate{}).
		Complete(r)
}
