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

// Package catalogmcpserver maps the CatalogMcpServer CRD to the AIM catalog of the APIM Automation API
// through the SDK.
package catalogmcpserver

import (
	"context"
	"fmt"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/catalogmcpserver"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
)

// ToDTO maps the CRD onto the payload a PUT sends. It is the single mapping used by the sync path, the dry
// run and drift detection. SDK defaults are not applied. The CRD nests each auth variant under its own key
// while the payload carries it flat with its type, so the auth is mapped field by field, not through JSON.
func ToDTO(obj *v1alpha1.CatalogMcpServer) (sdk.AimCatalogMcpServer, error) {
	spec := obj.Spec.Type

	auth, err := toAuth(spec.Connection.Auth)
	if err != nil {
		return sdk.AimCatalogMcpServer{}, err
	}

	transport := sdk.AimCatalogMcpServerTransport(spec.Connection.Transport)
	if transport == "" {
		transport = sdk.AimCatalogMcpServerTransportHTTP
	}

	return sdk.AimCatalogMcpServer{
		Hrid:        refs.NewNamespacedNameFromObject(obj).HRID(),
		EntityId:    spec.EntityID,
		Description: spec.Description,
		Connection: sdk.AimCatalogMcpServerConnection{
			Endpoint:  spec.Connection.Endpoint,
			Transport: transport,
			Auth:      auth,
		},
	}, nil
}

// toAuth always returns an auth, NONE when the spec has none: the platform reads an absent auth back as
// NONE, and drift would otherwise report it. Each variant carries its type, as the platform's copy does
// once decoded, so drift compares it too.
func toAuth(auth *catalogmcpserver.Auth) (*sdk.AimCatalogMcpServerAuth, error) {
	none := &sdk.AimCatalogMcpServerAuth{
		None: &sdk.AimCatalogMcpServerNoAuth{Type: sdk.AimCatalogMcpServerNoAuthTypeNONE},
	}
	if auth == nil || auth.Type == "" {
		return none, nil
	}

	switch auth.Type {
	case catalogmcpserver.AuthTypeNone:
		return none, nil
	case catalogmcpserver.AuthTypeHeader:
		if auth.Header == nil {
			return nil, fmt.Errorf("auth type %s has no header", auth.Type)
		}
		return &sdk.AimCatalogMcpServerAuth{Header: &sdk.AimCatalogMcpServerHeaderAuth{
			Type:  sdk.AimCatalogMcpServerHeaderAuthTypeHEADER,
			Name:  auth.Header.Name,
			Value: new(auth.Header.Value),
		}}, nil
	case catalogmcpserver.AuthTypeOAuth2:
		if auth.OAuth2 == nil {
			return nil, fmt.Errorf("auth type %s has no oauth2", auth.Type)
		}
		oauth2 := &sdk.AimCatalogMcpServerOAuth2Auth{
			Type:         sdk.AimCatalogMcpServerOAuth2AuthTypeOAUTH2,
			ClientId:     auth.OAuth2.ClientID,
			ClientSecret: new(auth.OAuth2.ClientSecret),
			TokenUrl:     auth.OAuth2.TokenURL,
		}
		if auth.OAuth2.Scope != nil {
			oauth2.Scope = new(*auth.OAuth2.Scope)
		}
		return &sdk.AimCatalogMcpServerAuth{OAuth2: oauth2}, nil
	}

	return nil, fmt.Errorf("unknown auth type %s", auth.Type)
}

// CreateAPIMClient builds the SDK client of the management context the server references.
func CreateAPIMClient(ctx context.Context, obj *v1alpha1.CatalogMcpServer) (*apim.SDKClient, error) {
	return apim.SDKFromContextRef(ctx, obj.ContextRef(), obj.GetNamespace())
}
