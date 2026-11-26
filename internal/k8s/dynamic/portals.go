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

package dynamic

import (
	"context"
	"errors"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

// ErrPortalNotSynced is returned for a Portal GKO has not written to APIM yet. APIM refuses a
// listing, a page or a link of a portal it does not hold while the environment holds another
// portal, its default portal included, so a resource applied with its Portal waits for it.
var ErrPortalNotSynced = errors.New("portal not synced")

func ResolvePortal(ctx context.Context, ref core.ObjectRef, parentNs string) (*v1alpha1.Portal, error) {
	prtl, err := resolveRef(ctx, ref, parentNs, PortalGVR, new(v1alpha1.Portal))
	if err != nil {
		return nil, err
	}

	return prtl, nil
}

// CheckPortalSynced returns ErrPortalNotSynced while prtl has no APIM id: admission warns and
// skips the dry run, and the controller requeues.
func CheckPortalSynced(prtl *v1alpha1.Portal) error {
	if prtl.Status.ID == "" {
		return fmt.Errorf("%w: [%s/%s]", ErrPortalNotSynced, prtl.GetNamespace(), prtl.GetName())
	}
	return nil
}
