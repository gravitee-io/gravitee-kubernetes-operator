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
	"net/http"

	amsdk "github.com/gravitee-io/gravitee-automation-sdk/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// Delete removes the domain from AM, which also deletes its identity providers. A domain already gone is
// not an error.
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

// Upsert creates or updates the domain in AM.
func Upsert(ctx context.Context, client *am.Client, dto amsdk.Domain) (am.BaseResponse, error) {
	resp, err := client.UpsertDomainWithResponse(ctx, nil, dto)

	if err = am.HasErrors(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return am.BaseResponse{}, err
	}
	if resp.JSON200 == nil {
		return am.BaseResponse{}, am.UnexpectedResponse(resp.HTTPResponse)
	}

	return am.BaseResponse{
		Key: resp.JSON200.Key,
		OrgEnv: am.OrgEnv{
			OrgID: client.GetOrgID(),
			EnvID: client.GetEnvID(),
		},
	}, nil
}

// UpdateStatus copies the upsert response into the CR status.
func UpdateStatus(_ context.Context, obj *v1alpha1.AMSecurityDomain, resp am.BaseResponse) error {
	obj.Status.Status.Key = resp.Key
	obj.Status.Status.OrgID = resp.GetOrgID()
	obj.Status.Status.EnvID = resp.GetEnvID()
	return nil
}
