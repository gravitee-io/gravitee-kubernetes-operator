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

package lifecycle

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/template"
)

func (l ResourceLifecycle[T, D, C, R]) Reconcile(
	ctx context.Context,
	cli client.Client,
	recorder record.EventRecorder,
	req ctrl.Request,
	obj T,
) (ctrl.Result, error) {
	if err := cli.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	obj.SetNamespace(req.Namespace)

	k8s.ResetConditionsExceptAutomationAPI(asConditionAware(obj))
	dc, ok := obj.DeepCopyObject().(T)
	if !ok {
		return ctrl.Result{}, fmt.Errorf("lifecycle: DeepCopyObject is not %T", obj)
	}

	events := event.NewRecorder(recorder)
	_, err := util.CreateOrUpdate(ctx, cli, obj, func() error {
		return l.mutate(ctx, events, obj, dc)
	})

	if copyErr := dc.GetStatus().DeepCopyTo(obj); copyErr != nil {
		return ctrl.Result{}, copyErr
	}

	if err == nil {
		log.InfoEndReconcile(ctx, obj)
		return ctrl.Result{}, l.updateStatusSuccess(ctx, cli, obj)
	}

	if stErr := l.updateStatusFailure(ctx, cli, obj, err); stErr != nil {
		return ctrl.Result{}, stErr
	}

	if gerrors.IsRecoverable(err) {
		log.ErrorRequeuingReconcile(ctx, err, obj)
		return ctrl.Result{}, err
	}

	log.ErrorAbortingReconcile(ctx, err, obj)
	return ctrl.Result{}, nil
}

func (l ResourceLifecycle[T, D, C, R]) mutate(
	ctx context.Context,
	events *event.Recorder,
	obj, dc T,
) error {
	if l.Finalizer != "" && !obj.IsBeingDeleted() {
		util.AddFinalizer(obj, l.Finalizer)
	}
	k8s.AddAnnotation(obj, core.LastSpecHashAnnotation, hash.Calculate(obj.GetSpec()))

	if obj.IsBeingDeleted() {
		if err := template.ReleaseReferences(ctx, obj); err != nil {
			return err
		}
		return events.Record(event.Delete, obj, func() error {
			return l.delete(ctx, obj, dc)
		})
	}

	if err := template.Compile(ctx, dc, true); err != nil {
		return gerrors.NewCompileTemplateError(err)
	}
	return events.Record(event.Update, obj, func() error {
		return l.upsert(ctx, dc)
	})
}

func (l ResourceLifecycle[T, D, C, R]) upsert(ctx context.Context, dc T) error {
	if err := l.resolveRefs(ctx, dc); err != nil {
		return err
	}

	if l.ResolveRefs != nil {
		k8s.SetCondition(asConditionAware(dc), k8s.NewResolvedRefsConditionBuilder(dc.GetGeneration()).
			ResolveRefs("All References successfully resolved").Build())
	}

	api, dto, err := l.clientAndDTO(ctx, dc)
	if err != nil {
		return err
	}
	resp, err := l.Upsert(ctx, api, dto)
	if err := wrapUnexpectedAsControlPlane(err); err != nil {
		return err
	}
	if l.PostUpsert == nil {
		return nil
	}
	return l.PostUpsert(ctx, dc, resp)
}

func (l ResourceLifecycle[T, D, C, R]) delete(ctx context.Context, obj, dc T) error {
	if err := l.resolveRefs(ctx, dc); err != nil {
		return err
	}
	api, dto, err := l.clientAndDTO(ctx, dc)
	if err != nil {
		return err
	}
	if l.DeleteGuard != nil {
		if err := l.DeleteGuard(ctx, dc); err != nil {
			return err
		}
	}
	if err := wrapUnexpectedAsControlPlane(l.Delete(ctx, api, dto)); err != nil {
		return err
	}
	if l.Finalizer != "" {
		util.RemoveFinalizer(obj, l.Finalizer)
	}
	return nil
}

func (l ResourceLifecycle[T, D, C, R]) clientAndDTO(ctx context.Context, dc T) (C, D, error) {
	var zeroC C
	var zeroD D
	api, err := l.ClientFactory(ctx, dc)
	if err != nil {
		return zeroC, zeroD, gerrors.NewResolveRefError(err)
	}
	return api, l.ToDTO(dc), nil
}

func (l ResourceLifecycle[T, D, C, R]) resolveRefs(ctx context.Context, obj T) error {
	ns := obj.GetNamespace()
	var err error
	if l.ResolveRefs != nil {
		err = l.ResolveRefs(ctx, obj, ns)
	}
	if err == nil {
		return nil
	}
	return gerrors.NewResolveRefError(err)
}

func (l ResourceLifecycle[T, D, C, R]) updateStatusSuccess(ctx context.Context, cli client.Client, obj T) error {
	if obj.IsBeingDeleted() {
		return nil
	}
	k8s.AddSuccessfulConditions(asConditionAware(obj))
	l.dropResolvedRefsWhenUnused(obj)
	return cli.Status().Update(ctx, obj)
}

func (l ResourceLifecycle[T, D, C, R]) updateStatusFailure(ctx context.Context, cli client.Client, obj T, err error) error {
	k8s.ErrorToCondition(obj, err)
	l.dropResolvedRefsWhenUnused(obj)
	return cli.Status().Update(ctx, obj)
}

// dropResolvedRefsWhenUnused removes the ResolvedRefs condition the shared status helpers
// always add: a resource without ResolveRefs has no references to report on.
func (l ResourceLifecycle[T, D, C, R]) dropResolvedRefsWhenUnused(obj T) {
	if l.ResolveRefs != nil {
		return
	}
	ca := asConditionAware(obj)
	conditions := ca.GetConditions()
	delete(conditions, k8s.ConditionResolvedRefs)
	ca.SetConditions(slices.SortedFunc(maps.Values(conditions), func(a, b metav1.Condition) int {
		return strings.Compare(a.Type, b.Type)
	}))
}

func asConditionAware[T core.ContextAwareObject](obj T) core.ConditionAwareObject {
	ca, ok := any(obj).(core.ConditionAwareObject)
	if !ok {
		panic(fmt.Sprintf("lifecycle: %T must implement ConditionAwareObject", obj))
	}
	return ca
}

func wrapUnexpectedAsControlPlane(err error) error {
	if err == nil {
		return nil
	}
	var re gerrors.ReconcileError
	if errors.As(err, &re) {
		return err
	}
	return gerrors.NewControlPlaneError(err)
}
