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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// GroupSyncedPredicate passes the Group events an API referencing the group needs. On top of
// LastSpecHashPredicate's create (a group not synced yet, so the API resolves the ref), it passes
// the update that records a successful APIM import (last-spec-hash annotation changed), which
// leaves the spec unchanged and is dropped by LastSpecHashPredicate, and the updates of a group
// being deleted.
type GroupSyncedPredicate struct {
	predicate.Funcs
}

func (GroupSyncedPredicate) Create(e event.CreateEvent) bool {
	return LastSpecHashPredicate{}.Create(e)
}

func (GroupSyncedPredicate) Update(e event.UpdateEvent) bool {
	oldGroup, oldOK := e.ObjectOld.(*v1alpha1.Group)
	newGroup, newOK := e.ObjectNew.(*v1alpha1.Group)
	if !oldOK || !newOK {
		return false
	}
	return newGroup.GetDeletionTimestamp() != nil ||
		oldGroup.Annotations[core.LastSpecHashAnnotation] != newGroup.Annotations[core.LastSpecHashAnnotation]
}

func (GroupSyncedPredicate) Delete(event.DeleteEvent) bool {
	return false
}

func (GroupSyncedPredicate) Generic(event.GenericEvent) bool {
	return false
}
