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

package drift

import (
	"log"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"k8s.io/apimachinery/pkg/runtime"
)

// Predicate is a function that returns true if the object should be considered for drift detection.
type Predicate func(object runtime.Object) bool

// unsupported predicate returns true if the object dpes not support drift detection.
var unsupported []Predicate

// disabled predicate returns true if the object does is disabled for drift detection (when no annotation is present).
var disabled []Predicate

func InitEnableCheck() {
	unsupported = append(unsupported, isLegacyGroup)
	disabled = append(disabled, portal, documentation, portalListing)
}

func IsDriftEnabled(crd runtime.Object) bool {
	// check if the CRD is supported
	for _, isUnsupported := range unsupported {
		if isUnsupported(crd) {
			return false
		}
	}
	if coreObj, ok := crd.(core.Object); ok {
		driftAnnot, hasAnnot := coreObj.GetAnnotations()[core.DriftDetectionAnnotation]
		if hasAnnot && driftAnnot == env.TrueString {
			return true
		} else if hasAnnot && driftAnnot == env.FalseString {
			return false
		}
	} else {
		log.Panicf("CRD Does not implement core.Object interface: %T", crd)
	}
	// no annotations: check if the CRD is disabled
	for _, isDisabled := range disabled {
		if isDisabled(crd) {
			return false
		}
	}
	return env.Config.DriftDetection
}

func isLegacyGroup(obj runtime.Object) bool {
	if g, ok := obj.(*v1alpha1.Group); ok {
		if k8s.IsAutomationAPIManaged(g) {
			return false
		}
		// we don't the drift detection for legacy groups
		return true
	}
	// not a group, pass
	return false
}

func portal(obj runtime.Object) bool {
	_, ok := obj.(*v1alpha1.Portal)
	return ok
}

func documentation(obj runtime.Object) bool {
	_, ok := obj.(*v1alpha1.Documentation)
	return ok
}

func portalListing(obj runtime.Object) bool {
	_, ok := obj.(*v1alpha1.PortalListing)
	return ok
}
