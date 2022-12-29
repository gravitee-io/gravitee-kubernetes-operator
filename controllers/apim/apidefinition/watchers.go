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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
func (r *Reconciler) NewUpdateFromLookup(lookupField string) UpdateFunc {
	return func(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
		ref := model.NewNamespacedName(e.ObjectNew.GetNamespace(), e.ObjectNew.GetName())
		r.queueRefs(lookupField, ref, q)
	}
}

// NewCreateFromLookup creates an updater function that will trigger an update
// on all API definitions that are referencing the created object
// The lookupField is the field that is used to lookup the API definitions
// Note that this field *must* have been registered as a cache index in our main func (see main.go).
// This can be used to reconcile API definitions when have been created before the referenced object (e.g. an API context).
func (r *Reconciler) NewCreateFromLookup(lookupField string) CreateFunc {
	return func(e event.CreateEvent, q workqueue.RateLimitingInterface) {
		ref := model.NewNamespacedName(e.Object.GetNamespace(), e.Object.GetName())
		r.queueRefs(lookupField, ref, q)
	}
}

// ContextWatcher creates a watcher that will trigger an update on all API definitions
// that are referencing the updated or created context.
// API can thus be created before referencing
// a context, and will be reconciled when the context is later created.
func (r *Reconciler) ContextWatcher(lookupField string) *handler.Funcs {
	return &handler.Funcs{
		UpdateFunc: r.NewUpdateFromLookup(indexer.ContextField.String()),
		CreateFunc: r.NewCreateFromLookup(indexer.ContextField.String()),
	}
}

func (r *Reconciler) ResourceWatcher(lookupField string) *handler.Funcs {
	return &handler.Funcs{
		UpdateFunc: r.NewUpdateFromLookup(indexer.ResourceField.String()),
	}
}

func (r *Reconciler) queueRefs(lookupField string, ref model.NamespacedName, q workqueue.RateLimitingInterface) {
	ctx := context.Background()
	log := log.FromContext(ctx).WithValues("reference", ref.String())

	apis, err := r.listForRef(ctx, ref, lookupField)
	if err != nil {
		log.Error(err, "unable to list APIs referencing resource, skipping update")
		return
	}

	for _, api := range apis {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      api.Name,
			Namespace: api.Namespace,
		}})
	}
}

func (r *Reconciler) listForRef(
	ctx context.Context, ref model.NamespacedName, field string,
) ([]gio.ApiDefinition, error) {
	apiDefinitionList := &gio.ApiDefinitionList{}

	filter := &client.ListOptions{
		FieldSelector: fields.SelectorFromSet(fields.Set{field: ref.String()}),
	}

	if err := r.Client.List(ctx, apiDefinitionList, filter); err != nil {
		return nil, err
	}

	return apiDefinitionList.Items, nil
}
