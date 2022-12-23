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
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type UpdaterFunc = func(event.UpdateEvent, workqueue.RateLimitingInterface)

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

// CreateUpdaterFromLookup creates an updater function that will trigger an update
// on all API definitions that are referencing the updated object
// The lookupField is the field that is used to lookup the API definitions
// Note that this field *must* have been registered as a cache index in our main func (see main.go).
func (r *Reconciler) CreateUpdaterFromLookup(lookupField string) UpdaterFunc {
	return func(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
		ctx := context.Background()
		ref := model.NewNamespacedName(e.ObjectNew.GetNamespace(), e.ObjectNew.GetName())
		log := log.FromContext(ctx).WithValues("reference", ref.String())

		apis, err := r.listForRef(ctx, ref, lookupField)
		if err != nil {
			log.Error(err, "unable to list APIs for context, skipping update")
			return
		}

		for _, api := range apis {
			q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
				Name:      api.Name,
				Namespace: api.Namespace,
			}})
		}
	}
}

func (r *Reconciler) listForRef(
	ctx context.Context, ref model.NamespacedName, field string,
) ([]gio.ApiDefinition, error) {
	log := log.FromContext(ctx).WithValues("field", field)
	apiDefinitionList := &gio.ApiDefinitionList{}

	filter := &client.ListOptions{
		FieldSelector: fields.SelectorFromSet(fields.Set{field: ref.String()}),
	}

	if err := r.Client.List(ctx, apiDefinitionList, filter); err != nil {
		log.Error(err, "unable to list API definitions from reference")
		return nil, err
	}

	return apiDefinitionList.Items, nil
}
