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

package reporter

import (
	"context"
	"net/http"

	amsdk "github.com/gravitee-io/gravitee-automation-sdk/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// PreCheck validates the reporter
// key pattern and length;
// domainRef.namespace must match the resource's own namespace;
// If the domain is not found or ready, a warning is issued.
// if 'system' is true: setting any of name, type, and configuration gives a warning.
func PreCheck(ctx context.Context, obj *v1alpha1.AMReporter) *gerrors.AdmissionErrors {
	errs := gerrors.NewAdmissionErrors()

	key := refs.NewNamespacedNameFromObject(obj).HRID()
	errs.MergeWith(am.ValidateKey(key))
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(am.ValidateDomainRefNamespace(obj))
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(am.WarnMissingDomain(ctx, obj, obj.IsBeingDeleted()))
	errs.MergeWith(am.WarnSystemIgnoredFields(obj.Spec.System, obj.Spec.Name, obj.Spec.Type, obj.Spec.Configuration))

	return errs
}

// DryRun validates the reporter against AM without persisting it. It is skipped without a client
// (domain not in AM yet). Any AM error, an unreachable AM or a license refusal included, rejects it,
// as for the other AM kinds.
func DryRun(ctx context.Context, client *am.Client, dto Reporter) *gerrors.AdmissionErrors {
	errs := gerrors.NewAdmissionErrors()
	if client == nil {
		return errs
	}

	resp, err := client.UpsertReporterWithResponse(ctx, dto.DomainKey,
		&amsdk.UpsertReporterParams{DryRun: new(true)}, dto.Reporter)

	if err = am.HasErrors(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		errs.AddSevere(err.Error())
		return errs
	}
	if resp.JSON200 == nil {
		errs.AddSevere(am.UnexpectedResponse(resp.HTTPResponse).Error())
		return errs
	}
	return am.ToAdmissionErrors(resp.JSON200.DryRunErrors)
}

// GetRemote fetches the reporter for drift detection. Without a client (domain not in AM
// yet) nothing can have drifted: the resource is its own remote, so an update is admitted, as a create is.
func GetRemote(ctx context.Context, client *am.Client, dto Reporter) (Reporter, error) {
	if client == nil {
		return dto, nil
	}
	resp, err := client.GetReporterWithResponse(ctx, dto.DomainKey, dto.Key)
	if err = am.HasErrors(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return Reporter{}, err
	}
	if resp.JSON200 == nil {
		return Reporter{}, am.UnexpectedResponse(resp.HTTPResponse)
	}
	return Reporter{Reporter: *resp.JSON200, DomainKey: dto.DomainKey}, nil
}

// ToReporterDTOForDrift wraps ToReporterDTO and blanks the fields AM ignores for the system reporter:
// "name", "type", "configuration", "attributeMappings" and "attributeMappingEventTypes".
func ToReporterDTOForDrift(obj *v1alpha1.AMReporter) (Reporter, error) {
	dto, err := ToReporterDTO(obj)
	if err != nil {
		return Reporter{}, err
	}
	if utils.SafeDereference(dto.System) {
		dto.Name = nil
		dto.Type = nil
		dto.Configuration = nil
		dto.AttributeMappings = nil
		dto.AttributeMappingEventTypes = nil
	}
	return dto, nil
}
