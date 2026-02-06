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
	"errors"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/gateway"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/gateway-api/gateway/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	autoscalingV2 "k8s.io/api/autoscaling/v2"
	coreV1 "k8s.io/api/core/v1"
	policyV1 "k8s.io/api/policy/v1"
	kErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwAPIv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

var errSkipObject = errors.New("object should be skipped and this error should not be returned to the user")

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

	gwc, err := r.validateGatewayClass(ctx, gw)
	if err != nil {
		if errors.Is(err, errSkipObject) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	params, err := r.validateGatewayClassParameters(ctx, gwc)
	if err != nil {
		if errors.Is(err, errSkipObject) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	var dc *gateway.Gateway

	err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		freshGw := &gwAPIv1.Gateway{}
		if err := k8s.GetClient().Get(ctx, req.NamespacedName, freshGw); err != nil {
			if kErrors.IsNotFound(err) {
				return nil
			}
			return err
		}

		dc = gateway.WrapGateway(freshGw)

		gwcKey := client.ObjectKey{Name: gwc.Object.Name}
		if err := k8s.GetClient().Get(ctx, gwcKey, gwc.Object); client.IgnoreNotFound(err) != nil {
			return err
		} else if kErrors.IsNotFound(err) {
			log.Debug(ctx, "ignoring gateway as gateway class was deleted")
			return nil
		}

		return k8s.CreateOrUpdate(ctx, gw.Object, func() error {
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

				if err := internal.DetectConflicts(dc); err != nil {
					return err
				}

				if err := internal.Accept(dc); err != nil {
					return err
				}

				if !k8s.IsAccepted(dc) {
					return nil
				}

				return internal.Program(ctx, dc, params)
			})
		})
	})

	if err != nil {
		log.ErrorRequeuingReconcile(ctx, err, gw.Object)
		return ctrl.Result{}, err
	}

	if dc == nil {
		return ctrl.Result{}, nil
	}

	if err := k8s.GetClient().Get(ctx, req.NamespacedName, gw.Object); client.IgnoreNotFound(err) != nil {
		log.ErrorRequeuingReconcile(ctx, err, gw.Object)
		return ctrl.Result{}, err
	} else if kErrors.IsNotFound(err) {
		log.Debug(ctx, "Looks like the Gateway was deleted during reconciliation, no need to update status")
		return ctrl.Result{}, nil
	}

	dc.Object.Status.DeepCopyInto(&gw.Object.Status)
	if err := k8s.UpdateStatus(ctx, gw.Object); client.IgnoreNotFound(err) != nil {
		log.ErrorRequeuingReconcile(ctx, err, gw.Object)
		return ctrl.Result{}, err
	}

	log.Debug(ctx, "Re-checking and updating addresses after status update")
	if err := updateGatewayAddressesIfNeeded(ctx, gw); err != nil {
		log.ErrorRequeuingReconcile(ctx, err, gw.Object)
		return ctrl.Result{}, err
	}

	log.Debug(ctx, "Looking for service address ...")
	if result, shouldRequeue := checkLoadBalancerAddress(ctx, gw); shouldRequeue {
		return result, nil
	}

	log.InfoEndReconcile(ctx, gw.Object)
	return ctrl.Result{}, nil
}

func checkLoadBalancerAddress(ctx context.Context, gw *gateway.Gateway) (ctrl.Result, bool) {
	programmed := k8s.GetCondition(gw, k8s.ConditionProgrammed)
	if programmed == nil {
		return ctrl.Result{}, false
	}
	if programmed.Status != k8s.ConditionStatusTrue {
		return ctrl.Result{}, false
	}
	if len(gw.Object.Status.Addresses) > 0 {
		return ctrl.Result{}, false
	}

	svcList := &coreV1.ServiceList{}
	if err := k8s.GetClient().List(
		ctx,
		svcList,
		&client.ListOptions{
			Namespace:     gw.Object.Namespace,
			LabelSelector: labels.SelectorFromSet(k8s.GwAPIv1GatewayLabels(gw.Object.Name)),
		},
	); err != nil {
		return ctrl.Result{}, false
	}

	for i := range svcList.Items {
		svc := &svcList.Items[i]
		if !k8s.IsGatewayDependent(gw, svc) {
			continue
		}

		if svc.Spec.Type == coreV1.ServiceTypeLoadBalancer {
			if len(svc.Status.LoadBalancer.Ingress) == 0 {
				log.Debug(ctx, "LoadBalancer service has no IP assigned yet, requeuing gateway")
				return ctrl.Result{RequeueAfter: 5 * time.Second}, true
			}
		}
		break
	}

	return ctrl.Result{}, false
}

func updateGatewayAddressesIfNeeded(ctx context.Context, gw *gateway.Gateway) error {
	programmed := k8s.GetCondition(gw, k8s.ConditionProgrammed)
	if programmed == nil || programmed.Status != k8s.ConditionStatusTrue {
		return nil
	}

	// Re-check addresses to ensure they're up-to-date after status update
	oldAddresses := gw.Object.Status.Addresses
	if err := internal.UpdateGatewayAddresses(ctx, gw); err != nil {
		return err
	}

	if !addressesEqual(oldAddresses, gw.Object.Status.Addresses) {
		return k8s.UpdateStatus(ctx, gw.Object)
	}

	return nil
}

func addressesEqual(a1, a2 []gwAPIv1.GatewayStatusAddress) bool {
	if len(a1) != len(a2) {
		return false
	}
	for i := range a1 {
		if a1[i].Value != a2[i].Value {
			return false
		}
		if (a1[i].Type == nil) != (a2[i].Type == nil) {
			return false
		}
		if a1[i].Type != nil && a2[i].Type != nil && *a1[i].Type != *a2[i].Type {
			return false
		}
	}
	return true
}

func (r *Reconciler) validateGatewayClass(
	ctx context.Context,
	gw *gateway.Gateway,
) (*gateway.GatewayClass, error) {
	gwcName := string(gw.Object.Spec.GatewayClassName)

	if gwcName == "" {
		log.Debug(ctx, "ignoring gateway as no gateway class name is defined")
		return nil, errSkipObject
	}

	gwcKey := client.ObjectKey{Name: gwcName}
	gwc := gateway.WrapGatewayClass(&gwAPIv1.GatewayClass{})

	if err := k8s.GetClient().Get(ctx, gwcKey, gwc.Object); client.IgnoreNotFound(err) != nil {
		return nil, err
	} else if kErrors.IsNotFound(err) {
		log.Debug(ctx, "ignoring gateway as gateway class name was not found")
		return nil, errSkipObject
	}

	if gwc.Object.Spec.ControllerName != core.GraviteeGatewayClassController {
		log.Debug(ctx, "ignoring gateway as controller name does not match")
		return nil, errSkipObject
	}

	return gwc, nil
}

func (r *Reconciler) validateGatewayClassParameters(
	ctx context.Context,
	gwc *gateway.GatewayClass,
) (*v1alpha1.GatewayClassParameters, error) {
	paramRef := gwc.Object.Spec.ParametersRef

	if paramRef == nil {
		return nil, errSkipObject
	}

	if paramRef.Group != gwAPIv1.Group(v1alpha1.GroupVersion.Group) {
		return nil, errSkipObject
	}

	if paramRef.Kind != "GatewayClassParameters" {
		return nil, errSkipObject
	}

	key := client.ObjectKey{
		Name:      paramRef.Name,
		Namespace: string(*paramRef.Namespace),
	}

	params := new(v1alpha1.GatewayClassParameters)

	if err := k8s.GetClient().Get(ctx, key, params); client.IgnoreNotFound(err) != nil {
		return nil, err
	} else if kErrors.IsNotFound(err) {
		log.Debug(ctx, "ignoring gateway as gateway class parameters were not found")
		return nil, errSkipObject
	}

	return params, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gwAPIv1.Gateway{}).
		Watches(&gwAPIv1.GatewayClass{}, internal.WatchGatewayClasses()).
		Watches(&gwAPIv1.HTTPRoute{}, internal.WatchHTTPRoutes()).
		Watches(&v1alpha1.KafkaRoute{}, internal.WatchKafkaRoutes()).
		Watches(&coreV1.Service{}, internal.WatchServices()).
		Watches(&coreV1.Secret{}, internal.WatchSecrets()).
		Watches(&gwAPIv1beta1.ReferenceGrant{}, internal.WatchReferenceGrants()).
		Watches(&autoscalingV2.HorizontalPodAutoscaler{}, internal.WatchHPAs()).
		Watches(&policyV1.PodDisruptionBudget{}, internal.WatchPDBs()).
		Complete(r)
}
