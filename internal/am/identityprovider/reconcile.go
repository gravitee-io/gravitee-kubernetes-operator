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

package identityprovider

import (
	"context"
	"net/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// Upsert creates or updates the identity provider under its domain in AM.
func Upsert(ctx context.Context, client *am.Client, dto IdentityProvider) (Response, error) {
	resp, err := client.UpsertIdentityProviderWithResponse(ctx, dto.DomainKey, nil, dto.IdentityProvider)
	if err = am.HasErrors(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return Response{}, err
	}
	if resp.JSON200 == nil {
		return Response{}, am.UnexpectedResponse(resp.HTTPResponse)
	}

	return Response{
		DomainSubResourceResponse: am.DomainSubResourceResponse{
			DomainKey: dto.DomainKey,
			BaseResponse: am.BaseResponse{
				Key: resp.JSON200.Key,
				OrgEnv: am.OrgEnv{
					OrgID: client.GetOrgID(),
					EnvID: client.GetEnvID(),
				},
			},
		},
		Name: utils.SafeDereference(resp.JSON200.Name),
		Type: utils.SafeDereference(resp.JSON200.Type),
	}, nil
}

// Delete removes the identity provider from its domain in AM. An identity provider already gone,
// deleted outside the operator or with its domain, is not an error.
func Delete(ctx context.Context, client *am.Client, dto IdentityProvider) error {
	resp, err := client.DeleteIdentityProviderWithResponse(ctx, dto.DomainKey, dto.Key)
	err = am.HasErrors(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	})
	if errors.IsNotFound(err) {
		// already gone from AM (deleted outside the operator, or with its domain): nothing left to delete
		return nil
	}
	return err
}

// UpdateStatus copies the upsert response into the CR status.
func UpdateStatus(_ context.Context, obj *v1alpha1.AMIdentityProvider, resp Response) error {
	obj.Status.Key = resp.Key
	obj.Status.DomainKey = resp.DomainKey
	obj.Status.OrgID = resp.GetOrgID()
	obj.Status.EnvID = resp.GetEnvID()
	obj.Status.Name = resp.Name
	obj.Status.Type = resp.Type
	return nil
}
