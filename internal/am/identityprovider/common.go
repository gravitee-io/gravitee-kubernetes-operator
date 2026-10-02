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
	amsdk "github.com/gravitee-io/gravitee-automation-sdk/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/mapper"
)

// IdentityProvider is the AM payload of an identity provider, with the key of the domain it lives under.
type IdentityProvider struct {
	amsdk.IdentityProvider
	DomainKey string
}

// Response is the upsert response of an identity provider: the sub-resource keys, and the name and type AM
// stored (its own for the system identity provider).
type Response struct {
	am.DomainSubResourceResponse
	Name string
	Type string
}

// ToIdentityProviderDTO converts an AMIdentityProvider object into an IdentityProvider DTO, applying defaults and keys.
func ToIdentityProviderDTO(obj *v1alpha1.AMIdentityProvider) (IdentityProvider, error) {
	dto, err := mapper.MapViaJSON[amsdk.IdentityProvider](obj.Spec.IdentityProvider)
	if err != nil {
		return IdentityProvider{}, err
	}
	dto = dto.WithDefaults()
	dto.Key = refs.NewNamespacedNameFromObject(obj).HRID()
	return IdentityProvider{
		IdentityProvider: dto,
		DomainKey:        am.DomainKey(obj),
	}, nil
}
