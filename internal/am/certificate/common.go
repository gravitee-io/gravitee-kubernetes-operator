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

package certificate

import (
	"time"

	amsdk "github.com/gravitee-io/gravitee-automation-sdk/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/mapper"
)

// Certificate is the AM payload of a certificate, with the key of the domain it lives under.
type Certificate struct {
	amsdk.Certificate
	DomainKey string
}

// Response is the upsert response of a certificate: the sub-resource keys and its expiry, when AM knows it.
type Response struct {
	am.DomainSubResourceResponse
	ExpiresAt *time.Time
}

// ToCertificateDTO converts an AMCertificate object into a Certificate DTO, applying defaults and keys.
func ToCertificateDTO(obj *v1alpha1.AMCertificate) (Certificate, error) {
	dto, err := mapper.MapViaJSON[amsdk.Certificate](obj.Spec.Certificate)
	if err != nil {
		return Certificate{}, err
	}
	dto = dto.WithDefaults()
	dto.Key = refs.NewNamespacedNameFromObject(obj).HRID()
	return Certificate{
		Certificate: dto,
		DomainKey:   am.DomainKey(obj),
	}, nil
}
