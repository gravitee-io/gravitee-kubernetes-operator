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

package securitydomain

import (
	"context"
	"fmt"
	"net/http"

	amsdk "github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

func Delete(ctx context.Context, amClient *am.Client, dto amsdk.Domain) error {
	resp, err := amClient.DeleteDomainWithResponse(ctx, dto.Identity())
	err = am.HasErrors(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	})
	if errors.IsNotFound(err) {
		// already gone from AM (deleted outside the operator): nothing left to delete
		return nil
	}
	return err
}

func Upsert(ctx context.Context, client *am.Client, dto amsdk.Domain) (DomainResponse, error) {
	resp, err := client.UpsertDomainWithResponse(ctx, nil, dto)

	if err = am.HasErrors(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return DomainResponse{}, err
	}
	if resp.JSON200 == nil {
		return DomainResponse{}, unexpectedResponse(resp.HTTPResponse)
	}

	return DomainResponse{
		Key: resp.JSON200.Key,
		OrgEnv: am.OrgEnv{
			OrgID: client.GetOrgID(),
			EnvID: client.GetEnvID(),
		},
	}, nil
}

func UpdateStatus(_ context.Context, obj *v1alpha1.AMSecurityDomain, resp DomainResponse) error {
	obj.Status.Status.ID = resp.Key
	obj.Status.Status.OrgID = resp.GetOrgID()
	obj.Status.Status.EnvID = resp.GetEnvID()
	return nil
}

func CreateAMClient(ctx context.Context, obj *v1alpha1.AMSecurityDomain) (*am.Client, error) {
	if !obj.HasContext() {
		return nil, fmt.Errorf("contextRef empty on %s [%s/%s]", obj.Kind, obj.GetName(), obj.GetNamespace())
	}

	// resolved like the APIM contexts: templates compiled, fetched from the API server
	resolved, err := dynamic.ResolveAMContext(ctx, obj.ContextRef(), obj.GetNamespace())
	if err != nil {
		return nil, fmt.Errorf("AMContext [%s]: %w", obj.ContextRef().String(), err)
	}
	amContext, ok := resolved.(*v1alpha1.AMContext)
	if !ok {
		return nil, fmt.Errorf("AMContext [%s]: unexpected type %T", obj.ContextRef().String(), resolved)
	}

	return am.NewSDKClient(ctx, amContext)
}

func DryRun(ctx context.Context, client *am.Client, dto amsdk.Domain) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	resp, err := client.UpsertDomainWithResponse(ctx, new(amsdk.UpsertDomainParams{
		DryRun: new(true),
	}), dto)

	if err = am.HasErrors(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		errs.AddSevere(err.Error())
		return errs
	}
	if resp.JSON200 == nil {
		errs.AddSevere(unexpectedResponse(resp.HTTPResponse).Error())
		return errs
	}
	return am.ToAdmissionErrors(resp.JSON200.DryRunErrors)
}

func GetRemote(ctx context.Context, client *am.Client, dto amsdk.Domain) (amsdk.Domain, error) {
	resp, err := client.GetDomainWithResponse(ctx, dto.Key)
	if err = am.HasErrors(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return amsdk.Domain{}, err
	}
	if resp.JSON200 == nil {
		return amsdk.Domain{}, unexpectedResponse(resp.HTTPResponse)
	}
	return *resp.JSON200, nil
}

func unexpectedResponse(resp *http.Response) error {
	return fmt.Errorf("unexpected AM response: status %d, content type %q",
		resp.StatusCode, resp.Header.Get("Content-Type"))
}
