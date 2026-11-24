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

package application

import (
	"context"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// PayloadChanged reports an update the controller will send to APIM: a spec change, or a change of the
// ignore-groups annotation, which decides whether the groups are applied.
func PayloadChanged(oldObj, newObj *v1alpha1.Application) bool {
	return oldObj.GetSpec().Hash() != newObj.GetSpec().Hash() ||
		oldObj.GetAnnotations()[core.IgnoreGroupsAnnotation] != newObj.GetAnnotations()[core.IgnoreGroupsAnnotation]
}

// DryRun validates the application against APIM without persisting it. It returns APIM's validation errors as
// findings, and APIM's error when the call fails. Warnings are only reported when there is no severe error.
func DryRun(ctx context.Context, client *apim.SDKClient, dto DTO) (*errors.AdmissionErrors, error) {
	state, err := upsert(ctx, client, dto, true)
	if err != nil {
		return nil, err
	}
	errs := errors.NewAdmissionErrors()
	if state.Errors == nil {
		return errs, nil
	}
	for _, severe := range state.Errors.Severe {
		errs.AddSevere(severe)
	}
	if errs.IsSevere() {
		return errs, nil
	}
	for _, warning := range state.Errors.Warning {
		errs.AddWarning(warning)
	}
	return errs, nil
}

// GetRemote fetches the application from APIM as the payload a PUT would send, for drift detection.
// With ignore-groups, APIM owns the groups: they are left out, as the drift mapper leaves out the CR's.
func GetRemote(ctx context.Context, client *apim.SDKClient, dto DTO) (DTO, error) {
	state, err := GetState(ctx, client, dto)
	if err != nil {
		return DTO{}, err
	}
	// Spec is a shallow copy of the state: only reassign its fields, never mutate shared slices or maps.
	remote := DTO{Application: state.Spec(), Legacy: dto.Legacy, IgnoreGroups: dto.IgnoreGroups}
	if dto.IgnoreGroups {
		remote.Groups = nil
	}
	return remote, nil
}

// GetState fetches the application state from APIM, including the server-managed fields (id, status).
func GetState(ctx context.Context, client *apim.SDKClient, dto DTO) (sdk.ApplicationState, error) {
	resp, err := client.GetApplicationWithResponse(ctx, dto.Identity(), dto.Legacy.editors()...)
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return sdk.ApplicationState{}, err
	}
	if resp.JSON200 == nil {
		return sdk.ApplicationState{}, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}
	return *resp.JSON200, nil
}
