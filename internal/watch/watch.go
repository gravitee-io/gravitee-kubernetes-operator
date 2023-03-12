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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/list"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Interface interface {
	WatchContexts() *handler.Funcs
	WatchResources() *handler.Funcs
}

type UpdateFunc = func(event.UpdateEvent, workqueue.RateLimitingInterface)
type CreateFunc = func(event.CreateEvent, workqueue.RateLimitingInterface)

type Type struct {
	ctx  context.Context
	k8s  client.Client
	kind schema.GroupVersionKind
}

func New(ctx context.Context, k8s client.Client, gvk schema.GroupVersionKind) *Type {
	return &Type{
		ctx:  ctx,
		k8s:  k8s,
		kind: gvk,
	}
}

// Watch context can be used to trigger a reconciliation when a context is updated
// on resources that should be synced with this context.
func (w *Type) WatchContexts() *handler.Funcs {
	return &handler.Funcs{
		UpdateFunc: w.UpdateFromLookup(indexer.ContextField),
		CreateFunc: w.CreateFromLookup(indexer.ContextField),
	}
}

// Watch resources can be used to trigger a reconciliation when an API resource is updated
// on resources that are referencing this resource.
func (w *Type) WatchResources() *handler.Funcs {
	return &handler.Funcs{
		UpdateFunc: w.UpdateFromLookup(indexer.ResourceField),
	}
}

// UpdateFromLookup creates an updater function that will trigger an update
// on all resources that are referencing the updated object
// The lookupField is the field that is used to lookup the resources
// Note that this field *must* have been registered as a cache index in our main func (see main.go).
func (w *Type) UpdateFromLookup(field indexer.IndexField) UpdateFunc {
	return func(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
		ref := model.NewNamespacedName(e.ObjectNew.GetNamespace(), e.ObjectNew.GetName())
		w.queueByFieldReferencing(field, ref, q)
	}
}

// CreateFromLookup creates an updater function that will trigger an update
// on all resources that are referencing the created object
// The lookupField is the field that is used to lookup the resources
// Note that this field *must* have been registered as a cache index in our main func (see main.go).
// This can be used to reconcile resources when have been created before the referenced object (e.g. an API context).
func (w *Type) CreateFromLookup(field indexer.IndexField) CreateFunc {
	return func(e event.CreateEvent, q workqueue.RateLimitingInterface) {
		ref := model.NewNamespacedName(e.Object.GetNamespace(), e.Object.GetName())
		w.queueByFieldReferencing(field, ref, q)
	}
}

func (w *Type) queueByFieldReferencing(
	field indexer.IndexField,
	ref model.NamespacedName,
	q workqueue.RateLimitingInterface,
) {
	list, err := list.OfKind(w.kind)

	if err != nil {
		log.FromContext(w.ctx).Error(err, "unable to initialize list from watcher kind", "kind", w.kind.String())
		return
	}

	if sErr := search.New(w.ctx, w.k8s).FindByFieldReferencing(field, ref, list); sErr != nil {
		log.FromContext(w.ctx).Error(sErr, "error while searching for items referencing", "reference", ref.String())
		return
	}

	for _, item := range list.GetItems() {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      item.GetName(),
			Namespace: item.GetNamespace(),
		}})
	}
}
