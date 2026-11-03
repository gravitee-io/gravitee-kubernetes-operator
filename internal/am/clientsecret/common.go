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
	"time"

	amsdk "github.com/gravitee-io/gravitee-automation-sdk/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
)

// ClientSecret is the AM payload of a client secret, with the keys of the domain and application it lives under.
type ClientSecret struct {
	amsdk.ClientSecret
	DomainKey      string
	ApplicationKey string
	// Accepted: AM has accepted the secret once (status.key set). Never sent to AM.
	Accepted bool `drift:"ignore"`
}

// Response is the upsert response of a client secret: the sub-resource keys and the expiration AM set.
type Response struct {
	am.DomainSubResourceResponse
	ApplicationKey string
	ExpiresAt      *time.Time
}

// ToClientSecretDTO converts an AMClientSecret object into a ClientSecret DTO. The value and the domain key
// come from obj.Resolved, filled by the reference resolvers: the value is left out when it is not resolved.
func ToClientSecretDTO(obj *v1alpha1.AMClientSecret) (ClientSecret, error) {
	dto := amsdk.ClientSecret{Key: refs.NewNamespacedNameFromObject(obj).HRID()}
	if value := obj.Resolved.Value; value != "" {
		dto.Value = &value
	}
	application := refs.NewNamespacedName(obj.GetNamespace(), obj.Spec.ApplicationRef.Name)
	return ClientSecret{
		ClientSecret:   dto,
		DomainKey:      obj.Resolved.DomainKey,
		ApplicationKey: application.HRID(),
		Accepted:       obj.Status.Key != "",
	}, nil
}
