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

package clientsecret

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Upsert creates the client secret under its application in AM. The same value again is a no-op;
// a different one is rejected by AM with 400: a secret is immutable. Once the secret is accepted, a 400
// keeps AM's message and points to rotation (AC6).
func Upsert(ctx context.Context, client *am.Client, dto ClientSecret) (Response, error) {
	resp, err := client.UpsertClientSecretWithResponse(ctx, dto.DomainKey, dto.ApplicationKey, nil, dto.ClientSecret)
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		if dto.Accepted && errors.IsBadRequest(err) {
			return Response{}, fmt.Errorf("%w; a client secret is immutable: rotate it with a new AMClientSecret", err)
		}
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
		ApplicationKey: dto.ApplicationKey,
		ExpiresAt:      resp.JSON200.ExpiresAt,
	}, nil
}

// Delete removes the client secret from its application in AM. AM answers 204 for a secret already gone;
// an application or domain already gone is not an error either.
func Delete(ctx context.Context, client *am.Client, dto ClientSecret) error {
	resp, err := client.DeleteClientSecretWithResponse(ctx, dto.DomainKey, dto.ApplicationKey, dto.Key)
	err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	})
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

// UpdateStatus copies the upsert response into the CR status.
func UpdateStatus(_ context.Context, obj *v1alpha1.AMClientSecret, resp Response) error {
	obj.Status.Key = resp.Key
	obj.Status.DomainKey = resp.DomainKey
	obj.Status.ApplicationKey = resp.ApplicationKey
	obj.Status.OrgID = resp.GetOrgID()
	obj.Status.EnvID = resp.GetEnvID()
	obj.Status.ExpiresAt = nil
	if resp.ExpiresAt != nil {
		obj.Status.ExpiresAt = new(metav1.NewTime(*resp.ExpiresAt))
	}
	return nil
}
