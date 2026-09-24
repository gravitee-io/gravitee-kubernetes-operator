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
	"encoding/json"
	"fmt"
	"net/http"

	amsdk "github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
)

func Delete(ctx context.Context, amClient *am.Client, dto amsdk.Domain) error {
	resp, err := amClient.DeleteDomainWithResponse(ctx, dto.Identity())
	return am.HasErrors(err, func() *http.Response {
		return resp.HTTPResponse
	})
}

func Upsert(ctx context.Context, client *am.Client, dto amsdk.Domain) (DomainResponse, error) {
	resp, err := client.UpsertDomainWithResponse(ctx, nil, dto)

	if err = am.HasErrors(err, func() *http.Response {
		return resp.HTTPResponse
	}); err != nil {
		return DomainResponse{}, err
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

	amContext := &v1alpha1.AMContext{}
	ref := obj.ContextRef()
	if ref.GetNamespace() == "" {
		ref.SetNamespace(obj.Namespace)
	}
	err := k8s.GetClient().Get(ctx, ref.NamespacedName(), amContext)
	if err != nil {
		return nil, fmt.Errorf("AMContext [%s] not found", ref.String())
	}

	return am.NewSDKClient(ctx, amContext)
}

func DryRun(ctx context.Context, client *am.Client, dto amsdk.Domain) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	resp, err := client.UpsertDomainWithResponse(ctx, new(amsdk.UpsertDomainParams{
		DryRun: new(true),
	}), dto, func(ctx context.Context, req *http.Request) error {
		js, _ := json.Marshal(dto)
		log.Info(ctx, "DryRun", "body", string(js))
		return nil
	})

	if err = am.HasErrors(err, func() *http.Response {
		return resp.HTTPResponse
	}); err != nil {
		errs.AddSevere(err.Error())
		return errs
	}
	return am.ToAdmissionErrors(resp.JSON200.DryRunErrors)
}

func GetRemote(ctx context.Context, client *am.Client, dto amsdk.Domain) (amsdk.Domain, error) {
	resp, err := client.GetDomainWithResponse(ctx, dto.Key)
	if err = am.HasErrors(err, func() *http.Response {
		return resp.HTTPResponse
	}); err != nil {
		return amsdk.Domain{}, err
	}
	return *resp.JSON200, nil
}

func DeleteGuard(_ context.Context, obj *v1alpha1.AMSecurityDomain) error {
	// if there is a status with data, that means it has been created at least once
	// there it cannot be deleted if the CRD has no contextRef
	if obj.Status.ID != "" && !obj.HasContext() {
		return fmt.Errorf("cannot delete %s [%s] without contextRef", obj.Kind, obj.GetRef().String())
	}
	return nil
}
