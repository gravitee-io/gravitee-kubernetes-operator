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

package reporter

import (
	amsdk "github.com/gravitee-io/gravitee-automation-sdk/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/mapper"
)

// Reporter is the AM payload of a reporter, with the key of the domain it lives under.
type Reporter struct {
	amsdk.Reporter
	DomainKey string
}

// Response is the upsert response of a reporter: the sub-resource keys and the name and type AM stored
// (its own for the system reporter).
type Response struct {
	am.DomainSubResourceResponse
	Name string
	Type string
}

// ToReporterDTO converts an AMReporter object into a Reporter DTO, applying defaults and keys.
func ToReporterDTO(obj *v1alpha1.AMReporter) (Reporter, error) {
	dto, err := mapper.MapViaJSON[amsdk.Reporter](obj.Spec.Reporter)
	if err != nil {
		return Reporter{}, err
	}
	dto = dto.WithDefaults()
	dto.Key = refs.NewNamespacedNameFromObject(obj).HRID()
	return Reporter{
		Reporter:  dto,
		DomainKey: am.DomainKey(obj),
	}, nil
}
