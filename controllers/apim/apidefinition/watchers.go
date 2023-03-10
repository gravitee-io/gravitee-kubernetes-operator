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

package apidefinition

import (
	"context"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type UpdateFunc = func(event.UpdateEvent, workqueue.RateLimitingInterface)
type CreateFunc = func(event.CreateEvent, workqueue.RateLimitingInterface)

// ApiUpdateFilter filters out update event that are coming from internal updates such as adding finalizers.
type ApiUpdateFilter struct {
	predicate.Funcs
}

func (ApiUpdateFilter) Update(e event.UpdateEvent) bool {
	if e.ObjectOld == nil || e.ObjectNew == nil {
		return false
	}
	if len(e.ObjectOld.GetFinalizers()) != len(e.ObjectNew.GetFinalizers()) {
		return false
	}
	return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
}

// NewUpdateFromLookup creates an updater function that will trigger an update
// on all API definitions that are referencing the updated object
// The lookupField is the field that is used to lookup the API definitions
// Note that this field *must* have been registered as a cache index in our main func (see main.go).
func (r *Reconciler) NewUpdateFromLookup(field indexer.IndexField) UpdateFunc {
	return func(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
		ref := model.NewNamespacedName(e.ObjectNew.GetNamespace(), e.ObjectNew.GetName())
		r.queueRefs(field, ref, q)
	}
}

// NewCreateFromLookup creates an updater function that will trigger an update
// on all API definitions that are referencing the created object
// The lookupField is the field that is used to lookup the API definitions
// Note that this field *must* have been registered as a cache index in our main func (see main.go).
// This can be used to reconcile API definitions when have been created before the referenced object (e.g. an API context).
func (r *Reconciler) NewCreateFromLookup(field indexer.IndexField) CreateFunc {
	return func(e event.CreateEvent, q workqueue.RateLimitingInterface) {
		ref := model.NewNamespacedName(e.Object.GetNamespace(), e.Object.GetName())
		r.queueRefs(field, ref, q)
	}
}

// ContextWatcher creates a watcher that will trigger an update on all API definitions
// that are referencing the updated or created context.
// API can thus be created before referencing
// a context, and will be reconciled when the context is later created.
func (r *Reconciler) ContextWatcher(field indexer.IndexField) *handler.Funcs {
	return &handler.Funcs{
		UpdateFunc: r.NewUpdateFromLookup(field),
		CreateFunc: r.NewCreateFromLookup(field),
	}
}

func (r *Reconciler) ResourceWatcher(field indexer.IndexField) *handler.Funcs {
	return &handler.Funcs{
		UpdateFunc: r.NewUpdateFromLookup(field),
	}
}

func (r *Reconciler) queueRefs(
	indexer indexer.IndexField,
	ref model.NamespacedName,
	q workqueue.RateLimitingInterface,
) {
	ctx := context.Background()
	apis := &gio.ApiDefinitionList{}
	if err := search.New(ctx, r.Client).FindByFieldReferencing(indexer, ref, apis); err != nil {
		log.FromContext(ctx).WithValues("reference", ref.String()).Error(
			err,
			"unable to list APIs referencing resource, skipping update",
		)
		return
	}

	for _, api := range apis.Items {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      api.Name,
			Namespace: api.Namespace,
		}})
	}
}
