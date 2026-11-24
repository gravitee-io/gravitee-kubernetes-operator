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

package portallink

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func verifyImmutableFields(oldLink, newLink *v1alpha1.PortalLink) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	// A link is parented by exactly one of a portal or an API; switching between
	// the two is a different (and clearer) error than moving it to a different
	// parent of the same kind.
	switch {
	case oldLink.IsPortalLink() && newLink.IsApiLink():
		errs.AddSevere("a portal link cannot be reassigned from a portal to an API")
		return errs
	case oldLink.IsApiLink() && newLink.IsPortalLink():
		errs.AddSevere("a portal link cannot be reassigned from an API to a portal")
		return errs
	}

	if refString(newLink.Spec.Portal) != refString(oldLink.Spec.Portal) {
		errs.AddSeveref(
			"portalRef is immutable; the link cannot be moved to a different portal "+
				"(from [%s] to [%s])",
			refString(oldLink.Spec.Portal), refString(newLink.Spec.Portal),
		)
		return errs
	}
	if refString(newLink.Spec.API) != refString(oldLink.Spec.API) {
		errs.AddSeveref(
			"apiRef is immutable; the link cannot be moved to a different API "+
				"(from [%s] to [%s])",
			refString(oldLink.Spec.API), refString(newLink.Spec.API),
		)
		return errs
	}

	return errs
}

func refString(ref *refs.NamespacedName) string {
	if ref == nil {
		return ""
	}
	return ref.String()
}
