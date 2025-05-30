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

package watch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/types/list"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type Interface interface {
	WatchContexts(index search.IndexField) *handler.Funcs
	WatchResources(index search.IndexField) *handler.Funcs
	WatchApiTemplate() *handler.Funcs
	WatchTLSSecret() *handler.Funcs
	WatchSharedPolicyGroups(index search.IndexField) *handler.Funcs
	WatchNotifications(index search.IndexField) *handler.Funcs
	WatchTemplatingSource(objKind string) *handler.Funcs
}

type UpdateFunc = func(context.Context, event.UpdateEvent, workqueue.TypedRateLimitingInterface[reconcile.Request])
type CreateFunc = func(context.Context, event.CreateEvent, workqueue.TypedRateLimitingInterface[reconcile.Request])

var NoopCreateFunc = func(context.Context, event.CreateEvent, workqueue.TypedRateLimitingInterface[reconcile.Request]) {
	// no op
}

type Type struct {
	ctx        context.Context
	k8s        client.Client
	objectList client.ObjectList
}

// New creates a new watch instance that can be used to trigger a reconciliation
// when a resource of interest is updated or created. The objectList parameter is used to
// determine the type of resources that should be reconciled on update or create events.
func New(ctx context.Context, k8s client.Client, objectList client.ObjectList) *Type {
	return &Type{
		ctx:        ctx,
		k8s:        k8s,
		objectList: objectList,
	}
}

// WatchContexts can be used to trigger a reconciliation when a management context is updated
// on resources that should be synced with this context.
func (w *Type) WatchContexts(index search.IndexField) *handler.Funcs {
	return &handler.Funcs{
		CreateFunc: w.CreateFromLookup(index),
		UpdateFunc: w.UpdateFromLookup(index),
	}
}

func ContextSecrets() *handler.Funcs {
	queueSecrets := func(obj client.Object, q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
		ctx, ok := obj.(*v1alpha1.ManagementContext)
		if !ok {
			return
		}

		if ctx.Spec.HasSecretRef() {
			ns := ctx.Spec.Auth.SecretRef.Namespace
			if ns == "" {
				ns = ctx.Namespace
			}
			q.Add(reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      ctx.Spec.Auth.SecretRef.Name,
					Namespace: ns,
				},
			})
		}
	}

	return &handler.Funcs{
		CreateFunc: func(_ context.Context, e event.CreateEvent, q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
			queueSecrets(e.Object, q)
		},
		DeleteFunc: func(_ context.Context, e event.DeleteEvent, q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
			queueSecrets(e.Object, q)
		},
	}
}

// WatchResources can be used to trigger a reconciliation when an API resource is updated
// on resources that are depending on it. Right now this is only used for API definitions.
func (w *Type) WatchResources(index search.IndexField) *handler.Funcs {
	return &handler.Funcs{
		CreateFunc: w.CreateFromLookup(index),
		UpdateFunc: w.UpdateFromLookup(index),
	}
}

// WatchApiTemplate can be used to trigger a reconciliation when an API template is updated
// on resources that are depending on it. Right now this is only used for Ingress resources.
func (w *Type) WatchApiTemplate() *handler.Funcs {
	return &handler.Funcs{
		UpdateFunc: w.UpdateFromLookup(search.ApiTemplateField),
		CreateFunc: w.CreateFromLookup(search.ApiTemplateField),
	}
}

// WatchTLSSecret can be used to trigger a reconciliation when an TLS secret is updated
// on resources that are depending on it. Right now this is only used for Ingress resources.
func (w *Type) WatchTLSSecret() *handler.Funcs {
	return &handler.Funcs{
		UpdateFunc: w.UpdateFromLookup(search.TLSSecretField),
		CreateFunc: w.CreateFromLookup(search.TLSSecretField),
	}
}

// WatchSharedPolicyGroups can be used to trigger a reconciliation when a SPG is updated
// on resources that should be synced with this SOG.
func (w *Type) WatchSharedPolicyGroups(index search.IndexField) *handler.Funcs {
	return &handler.Funcs{
		CreateFunc: w.CreateFromLookup(index),
		UpdateFunc: w.UpdateFromLookup(index),
	}
}

func (w *Type) WatchNotifications(index search.IndexField) *handler.Funcs {
	return &handler.Funcs{
		CreateFunc: w.CreateFromLookup(index),
		UpdateFunc: w.UpdateFromLookup(index),
	}
}

func (w *Type) WatchTemplatingSource(objKind string) *handler.Funcs {
	return &handler.Funcs{
		CreateFunc: NoopCreateFunc,
		UpdateFunc: w.UpdateForTemplating(objKind),
	}
}

// UpdateFromLookup creates an updater function that will trigger an update
// on all resources that are referencing the updated object.
// The lookupField is the field that is used to lookup the resources.
// Note that this field *must* have been registered as a cache index in our main func (see main.go).
func (w *Type) UpdateFromLookup(field search.IndexField) UpdateFunc {
	return func(_ context.Context, e event.UpdateEvent, q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
		ref := refs.NewNamespacedName(e.ObjectNew.GetNamespace(), e.ObjectNew.GetName())
		w.queueByFieldReferencing(field, ref, q)
	}
}

// CreateFromLookup creates an updater function that will trigger an update
// on all resources that are referencing the created object.
// The lookupField is the field that is used to lookup the resources.
// Note that this field *must* have been registered as a cache index in our main func (see main.go).
// This can be used to reconcile resources when have been created before their dependencies.
// For example, when an API is created before the management context it references.
func (w *Type) CreateFromLookup(field search.IndexField) CreateFunc {
	return func(_ context.Context, e event.CreateEvent, q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
		ref := refs.NewNamespacedName(e.Object.GetNamespace(), e.Object.GetName())
		w.queueByFieldReferencing(field, ref, q)
	}
}

func (w *Type) queueByFieldReferencing(
	field search.IndexField,
	ref refs.NamespacedName,
	q workqueue.TypedRateLimitingInterface[reconcile.Request],
) {
	objectList, err := list.OfType(w.objectList)
	queueKind := w.objectList.GetObjectKind().GroupVersionKind().Kind

	if err != nil {
		log.Error(w.ctx, err, fmt.Sprintf("Unable to create list from kind [%s]", queueKind))
		return
	}

	if err := search.FindByFieldReferencing(w.ctx, field, ref, objectList); err != nil {
		log.Error(
			w.ctx,
			err,
			fmt.Sprintf("Error while searching for items referencing [%s]", ref),
		)
		return
	}

	items, err := meta.ExtractList(objectList)
	if err != nil {
		log.Error(
			w.ctx,
			err,
			fmt.Sprintf("Error while extracting list items of kind [%s]", queueKind),
		)
	}

	for i := range items {
		if item, ok := items[i].(client.Object); !ok {
			log.Error(
				w.ctx,
				errors.New("unsupported type"),
				"List item is not a client object and cannot be added to the queue",
			)
		} else {
			q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
				Name:      item.GetName(),
				Namespace: item.GetNamespace(),
			}})
		}
	}
}

func (w *Type) UpdateForTemplating(objKind string) UpdateFunc {
	return func(_ context.Context, e event.UpdateEvent, q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
		w.queueForTemplating(objKind, e.ObjectOld, q)
	}
}

func (w *Type) queueForTemplating(
	objKind string,
	obj client.Object,
	q workqueue.TypedRateLimitingInterface[reconcile.Request],
) {
	if !util.ContainsFinalizer(obj, core.TemplatingFinalizer) {
		return
	}

	annotationValue := obj.GetAnnotations()[getObjectAnnotationName(objKind)]
	if annotationValue != "" {
		values := make([]string, 0)
		if err := json.Unmarshal([]byte(annotationValue), &values); err != nil {
			log.Error(
				w.ctx,
				err,
				fmt.Sprintf("Error while extracting list items of kind [%s]", objKind),
			)
			return
		}

		for _, val := range values {
			q.Add(reconcile.Request{NamespacedName: getNamespacedName(val)})
		}
	}
}

func getNamespacedName(annotationValueItem string) types.NamespacedName {
	nsAndName := strings.Split(annotationValueItem, "/")
	return types.NamespacedName{
		Name:      nsAndName[1],
		Namespace: nsAndName[0],
	}
}

func getObjectAnnotationName(objKind string) string {
	return "gravitee.io/" + objKind
}
