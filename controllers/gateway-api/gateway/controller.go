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

package gateway

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/gateway"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/gateway-api/gateway/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	kErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	coreV1 "k8s.io/api/core/v1"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type Reconciler struct {
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//nolint:gocognit,funlen // keep
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	gw := gateway.WrapGateway(&gwAPIv1.Gateway{})

	if err := k8s.GetClient().Get(ctx, req.NamespacedName, gw.Object); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	events := event.NewRecorder(r.Recorder)

	gwcName := string(gw.Object.Spec.GatewayClassName)

	if gwcName == "" {
		log.Debug(ctx, "ignoring gateway as no gateway class name is defined")
		return ctrl.Result{}, nil
	}

	gwcKey := client.ObjectKey{Name: gwcName}
	gwc := gateway.WrapGatewayClass(&gwAPIv1.GatewayClass{})

	if err := k8s.GetClient().Get(ctx, gwcKey, gwc.Object); client.IgnoreNotFound(err) != nil {
		return k8s.RequeueError(err)
	} else if kErrors.IsNotFound(err) {
		log.Debug(ctx, "ignoring gateway as gateway class name was not found")
		return ctrl.Result{}, nil
	}

	if gwc.Object.Spec.ControllerName != core.GraviteeGatewayClassController {
		log.Debug(ctx, "ignoring gateway as controller name does not match")
		return ctrl.Result{}, nil
	}

	paramRef := gwc.Object.Spec.ParametersRef

	if paramRef == nil {
		return ctrl.Result{}, nil
	}

	if paramRef.Group != gwAPIv1.Group(v1alpha1.GroupVersion.Group) {
		return ctrl.Result{}, nil
	}

	if paramRef.Kind != "GatewayClassParameters" {
		return ctrl.Result{}, nil
	}

	key := client.ObjectKey{
		Name:      paramRef.Name,
		Namespace: string(*paramRef.Namespace),
	}

	params := new(v1alpha1.GatewayClassParameters)

	if err := k8s.GetClient().Get(ctx, key, params); client.IgnoreNotFound(err) != nil {
		return k8s.RequeueError(err)
	} else if kErrors.IsNotFound(err) {
		log.Debug(ctx, "ignoring gateway as gateway class parameters were not found")
		return ctrl.Result{}, nil
	}

	dc := gw.DeepCopy()

	err := k8s.CreateOrUpdate(ctx, gw.Object, func() error {
		util.AddFinalizer(gw.Object, core.GatewayFinalizer)

		if !gw.Object.DeletionTimestamp.IsZero() {
			return events.Record(event.Delete, gw.Object, func() error {
				util.RemoveFinalizer(gw.Object, core.GatewayFinalizer)
				return nil
			})
		}

		gwcAccepted := k8s.GetCondition(gwc, k8s.ConditionAccepted)

		if gwcAccepted == nil {
			log.Debug(ctx, "ignoring gateway as gateway class accepted condition is not set")
			return nil
		}

		if gwcAccepted.Status == k8s.ConditionStatusFalse {
			log.Debug(ctx, "ignoring gateway as gateway class is not accepted")
			return nil
		}

		return events.Record(event.Update, gw.Object, func() error {
			internal.Init(dc)
			if err := internal.Resolve(ctx, dc, params); err != nil {
				return err
			}
			internal.DetectConflicts(dc)
			internal.Accept(dc)
			if !k8s.IsAccepted(dc) {
				return nil
			}
			return internal.Program(ctx, dc, params)
		})
	})

	if err != nil {
		log.ErrorRequeuingReconcile(ctx, err, gw.Object)
		return k8s.RequeueError(err)
	}

	dc.Object.Status.DeepCopyInto(&gw.Object.Status)
	if err := k8s.UpdateStatus(ctx, gw.Object); client.IgnoreNotFound(err) != nil {
		log.ErrorRequeuingReconcile(ctx, err, gw.Object)
		return k8s.RequeueError(err)
	}

	log.InfoEndReconcile(ctx, gw.Object)
	return ctrl.Result{Requeue: false}, err
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gwAPIv1.Gateway{}).
		Watches(&gwAPIv1.GatewayClass{}, internal.WatchGatewayClasses()).
		Watches(&gwAPIv1.HTTPRoute{}, internal.WatchHTTPRoutes()).
		Watches(&v1alpha1.KafkaRoute{}, internal.WatchKafkaRoutes()).
		Watches(&coreV1.Service{}, internal.WatchServices()).
		Watches(&coreV1.Secret{}, internal.WatchSecrets()).
		Complete(r)
}
