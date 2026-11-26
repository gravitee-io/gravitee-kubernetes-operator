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

package navigation

import (
	"context"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// WarnPortalNotSynced warns that the dry run of a listing, link or page is skipped while its portal is not
// synced: CreateAPIMClient then gives no client, and the controller waits for the portal. A parent that
// cannot be read was reported by the client factory, which runs first.
func WarnPortalNotSynced(ctx context.Context, portalRef, apiRef *refs.NamespacedName, ns string) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	parent, err := ResolveParent(ctx, portalRef, apiRef, ns)
	if err != nil {
		return errs
	}
	if err = CheckPortalSynced(parent); err != nil {
		errs.AddWarningf("platform checks deferred (%s)", err)
	}
	return errs
}

// Findings turns the errors of a dry-run state into admission findings: severe first, then warnings.
func Findings(stateErrors *sdk.Errors) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	if stateErrors == nil {
		return errs
	}
	for _, severe := range stateErrors.Severe {
		errs.AddSevere(severe)
	}
	for _, warning := range stateErrors.Warning {
		errs.AddWarning(warning)
	}
	return errs
}

// FindingsOnError reports a dry run that fails (transport, HTTP error, undecodable answer) as a severe
// finding with the error's own message, the way the webhooks always have, rather than as an error the
// lifecycle would prefix with "dry run failed: ".
func FindingsOnError[C any, D any](
	dryRun func(ctx context.Context, client C, dto D) (*errors.AdmissionErrors, error),
) func(ctx context.Context, client C, dto D) (*errors.AdmissionErrors, error) {
	return func(ctx context.Context, client C, dto D) (*errors.AdmissionErrors, error) {
		findings, err := dryRun(ctx, client, dto)
		if err != nil {
			errs := errors.NewAdmissionErrors()
			errs.AddSevere(err.Error())
			return errs, nil
		}
		return findings, nil
	}
}
