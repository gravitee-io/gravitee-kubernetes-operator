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

	amsdk "github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/v2/pkg/sdk"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// ValidateKey checks domain key
// see am.ValidateKey
func ValidateKey(_ context.Context, obj *v1alpha1.AMSecurityDomain) *errors.AdmissionErrors {
	key := domainKey(obj)
	return am.ValidateKey(key)
}

// DryRun validates the domain against AM without persisting it, and reports AM's dry-run errors.
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
		errs.AddSevere(am.UnexpectedResponse(resp.HTTPResponse).Error())
		return errs
	}
	return am.ToAdmissionErrors(resp.JSON200.DryRunErrors)
}

// GetRemote fetches the domain from AM, for drift detection.
func GetRemote(ctx context.Context, client *am.Client, dto amsdk.Domain) (amsdk.Domain, error) {
	resp, err := client.GetDomainWithResponse(ctx, dto.Key)
	if err = am.HasErrors(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return amsdk.Domain{}, err
	}
	if resp.JSON200 == nil {
		return amsdk.Domain{}, am.UnexpectedResponse(resp.HTTPResponse)
	}
	return *resp.JSON200, nil
}
