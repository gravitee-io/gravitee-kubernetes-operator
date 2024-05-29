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
	"fmt"
	"reflect"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/types/list"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Interface interface {
	WatchContexts(index indexer.IndexField) *handler.Funcs
	WatchResources(index indexer.IndexField) *handler.Funcs
	WatchApiTemplate() *handler.Funcs
	WatchTLSSecret() *handler.Funcs
}

type UpdateFunc = func(context.Context, event.UpdateEvent, workqueue.RateLimitingInterface)
type CreateFunc = func(context.Context, event.CreateEvent, workqueue.RateLimitingInterface)

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
func (w *Type) WatchContexts(index indexer.IndexField) *handler.Funcs {
	return &handler.Funcs{
		CreateFunc: w.CreateFromLookup(index),
		UpdateFunc: w.UpdateFromLookup(index),
	}
}

func ContextSecrets() *handler.Funcs {
	queueSecrets := func(obj client.Object, q workqueue.RateLimitingInterface) {
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
		CreateFunc: func(_ context.Context, e event.CreateEvent, q workqueue.RateLimitingInterface) {
			queueSecrets(e.Object, q)
		},
		DeleteFunc: func(_ context.Context, e event.DeleteEvent, q workqueue.RateLimitingInterface) {
			queueSecrets(e.Object, q)
		},
	}
}

// WatchResources can be used to trigger a reconciliation when an API resource is updated
// on resources that are depending on it. Right now this is only used for API definitions.
func (w *Type) WatchResources(index indexer.IndexField) *handler.Funcs {
	return &handler.Funcs{
		CreateFunc: w.CreateFromLookup(index),
		UpdateFunc: w.UpdateFromLookup(index),
	}
}

// WatchApiTemplate can be used to trigger a reconciliation when an API template is updated
// on resources that are depending on it. Right now this is only used for Ingress resources.
func (w *Type) WatchApiTemplate() *handler.Funcs {
	return &handler.Funcs{
		UpdateFunc: w.UpdateFromLookup(indexer.ApiTemplateField),
		CreateFunc: w.CreateFromLookup(indexer.ApiTemplateField),
	}
}

// WatchTLSSecret can be used to trigger a reconciliation when an TLS secret is updated
// on resources that are depending on it. Right now this is only used for Ingress resources.
func (w *Type) WatchTLSSecret() *handler.Funcs {
	return &handler.Funcs{
		UpdateFunc: w.UpdateFromLookup(indexer.TLSSecretField),
		CreateFunc: w.CreateFromLookup(indexer.TLSSecretField),
	}
}

// UpdateFromLookup creates an updater function that will trigger an update
// on all resources that are referencing the updated object.
// The lookupField is the field that is used to lookup the resources.
// Note that this field *must* have been registered as a cache index in our main func (see main.go).
func (w *Type) UpdateFromLookup(field indexer.IndexField) UpdateFunc {
	return func(_ context.Context, e event.UpdateEvent, q workqueue.RateLimitingInterface) {
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
func (w *Type) CreateFromLookup(field indexer.IndexField) CreateFunc {
	return func(_ context.Context, e event.CreateEvent, q workqueue.RateLimitingInterface) {
		ref := refs.NewNamespacedName(e.Object.GetNamespace(), e.Object.GetName())
		w.queueByFieldReferencing(field, ref, q)
	}
}

func (w *Type) queueByFieldReferencing(
	field indexer.IndexField,
	ref refs.NamespacedName,
	q workqueue.RateLimitingInterface,
) {
	objectList, err := list.OfType(w.objectList)

	if err != nil {
		log.FromContext(w.ctx).Error(err, "unable to create list of type", "type", w.objectList)
		return
	}

	if sErr := search.FindByFieldReferencing(w.ctx, field, ref, objectList); sErr != nil {
		log.FromContext(w.ctx).Error(sErr, "error while searching for items referencing", "reference", ref.String())
		return
	}

	items, err := meta.ExtractList(objectList)
	if err != nil {
		log.FromContext(w.ctx).Error(err, "error while extracting list items of type", "type", w.objectList)
	}

	for i := range items {
		if item, ok := items[i].(client.Object); !ok {
			log.FromContext(w.ctx).Error(
				fmt.Errorf("unable to convert the item to client.Object type"),
				"type", reflect.TypeOf(items[i]))
		} else {
			q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
				Name:      item.GetName(),
				Namespace: item.GetNamespace(),
			}})
		}
	}
}
