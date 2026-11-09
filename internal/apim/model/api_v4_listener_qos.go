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

package model

import v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"

// ForDrift applies the normalisations the ApiV4Definition drift comparison needs, to the CR and the
// remote side alike. Drift path only: the payload sent to APIM is unchanged.
func (dto APIV4DTO) ForDrift() APIV4DTO {
	return dto.WithResolvedVisibility().WithoutOverrideAccess().WithoutDefaultQos()
}

// WithoutDefaultQos removes a qos equal to the AUTO default from every listener entrypoint, so that
// drift detection treats it as absent. APIM stores AUTO for an entrypoint without qos, and stores
// no qos at all on a native (KAFKA) entrypoint, whatever was sent: there any qos is removed.
// Drift path only: the payload sent to APIM is unchanged.
func (dto APIV4DTO) WithoutDefaultQos() APIV4DTO {
	for _, listener := range dto.Listeners {
		if listener == nil || listener.GenericStringMap == nil {
			continue
		}
		native := listener.GetString("type") == string(v4.KafkaType)
		for _, item := range listener.GetSlice("entrypoints") {
			entrypoint, ok := item.(map[string]any)
			if ok && (native || entrypoint["qos"] == string(v4.AutoQOS)) {
				delete(entrypoint, "qos")
			}
		}
	}
	return dto
}
