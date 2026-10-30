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
	amsdk "github.com/gravitee-io/gravitee-automation-sdk/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/mapper"
)

// Application is the AM payload of an application, with the key of the domain it lives under.
type Application struct {
	amsdk.Application
	DomainKey string
}

// Response is the upsert response of an application: the sub-resource keys and the clientId AM stored
// (spec.clientId, or the one it generated).
type Response struct {
	am.DomainSubResourceResponse
	ClientID string
}

// ToApplicationDTO converts an AMApplication object into an Application DTO, applying defaults and keys.
func ToApplicationDTO(obj *v1alpha1.AMApplication) (Application, error) {
	app := obj.Spec.Application
	dto, err := mapper.MapViaJSON[amsdk.Application](app)
	if err != nil {
		return Application{}, err
	}
	dto = dto.WithDefaults()
	dto.Key = refs.NewNamespacedNameFromObject(obj).HRID()

	// a ref wins over a key, as in AMApplication.CertificateKey and IdentityKeys
	if app.CertificateRef != nil {
		dto.Certificate = obj.SiblingKey(app.CertificateRef)
	}

	for i, idp := range obj.Spec.IdentityProviders {
		if idp.IdentityRef != nil {
			dto.IdentityProviders[i].Identity = obj.SiblingKey(idp.IdentityRef)
		}
	}

	return Application{
		Application: dto,
		DomainKey:   am.DomainKey(obj),
	}, nil
}
