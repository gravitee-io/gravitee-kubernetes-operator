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

// PortalSyncedPredicate passes the Portal update a listing, page or link waiting on its portal needs: the
// one that changes the APIM id, which is what a first sync looks like. That update touches the status
// only, so LastSpecHashPredicate drops it. Any other event is dropped.
type PortalSyncedPredicate struct {
	predicate.Funcs
}

func (PortalSyncedPredicate) Create(event.CreateEvent) bool {
	return false
}

func (PortalSyncedPredicate) Update(e event.UpdateEvent) bool {
	oldPrtl, oldOK := e.ObjectOld.(*v1alpha1.Portal)
	newPrtl, newOK := e.ObjectNew.(*v1alpha1.Portal)
	if !oldOK || !newOK {
		return false
	}
	return oldPrtl.Status.ID != newPrtl.Status.ID
}

func (PortalSyncedPredicate) Delete(event.DeleteEvent) bool {
	return false
}

func (PortalSyncedPredicate) Generic(event.GenericEvent) bool {
	return false
}

// PortalChildPredicate filters the Portal events of a listing, page or link controller: a spec change, as
// for any watch, or the portal becoming synced.
func PortalChildPredicate() predicate.Predicate {
	return predicate.Or(LastSpecHashPredicate{}, PortalSyncedPredicate{})
}

// PortalThemeSyncedPredicate passes the PortalTheme update a portal waiting on its theme needs: the one that
// changes the APIM id. A portal synced before its theme is in APIM is refused, and that update touches the
// status only, so LastSpecHashPredicate drops it. Any other event is dropped.
type PortalThemeSyncedPredicate struct {
	predicate.Funcs
}

func (PortalThemeSyncedPredicate) Create(event.CreateEvent) bool {
	return false
}

func (PortalThemeSyncedPredicate) Update(e event.UpdateEvent) bool {
	oldThm, oldOK := e.ObjectOld.(*v1alpha1.PortalTheme)
	newThm, newOK := e.ObjectNew.(*v1alpha1.PortalTheme)
	if !oldOK || !newOK {
		return false
	}
	return oldThm.Status.ID != newThm.Status.ID
}

func (PortalThemeSyncedPredicate) Delete(event.DeleteEvent) bool {
	return false
}

func (PortalThemeSyncedPredicate) Generic(event.GenericEvent) bool {
	return false
}

// PortalThemeChildPredicate filters the PortalTheme events of the portal controller: a spec change, as for
// any watch, or the theme becoming synced.
func PortalThemeChildPredicate() predicate.Predicate {
	return predicate.Or(LastSpecHashPredicate{}, PortalThemeSyncedPredicate{})
}
