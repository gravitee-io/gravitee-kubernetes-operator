// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package predicate

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// AuthzPdpRoutedPredicate passes the AuthzPdp events a resource targeting the PDP through pdpRef
// needs. An update that changes status.routingScope is a first sync (or a new scope): it touches
// the status only, so LastSpecHashPredicate drops it. The update that marks the PDP for deletion
// passes, and so does the delete: the resolver reports a PDP being deleted as pending from that
// update on, and the platform removes the deleted scope from every target. A create is dropped: a
// new PDP is not synced yet, and the cache replays every PDP as a create when the operator starts.
// contextRef is immutable, so the context the resolver compares never changes on an update.
type AuthzPdpRoutedPredicate struct {
	predicate.Funcs
}

func (AuthzPdpRoutedPredicate) Create(event.CreateEvent) bool {
	return false
}

func (AuthzPdpRoutedPredicate) Update(e event.UpdateEvent) bool {
	oldPdp, oldOK := e.ObjectOld.(*v1alpha1.AuthzPdp)
	newPdp, newOK := e.ObjectNew.(*v1alpha1.AuthzPdp)
	if !oldOK || !newOK {
		return false
	}
	return oldPdp.Status.RoutingScope != newPdp.Status.RoutingScope ||
		oldPdp.IsBeingDeleted() != newPdp.IsBeingDeleted()
}

func (AuthzPdpRoutedPredicate) Delete(e event.DeleteEvent) bool {
	_, ok := e.Object.(*v1alpha1.AuthzPdp)
	return ok
}

func (AuthzPdpRoutedPredicate) Generic(event.GenericEvent) bool {
	return false
}
