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

package apim

import (
	"slices"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apiv4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/mapper"
)

// crdV4 maps a CRD spec (JSON) to the DTO the reconcile sends, through apiv4.ToDTO.
func crdV4(spec string) apiv4.DTO {
	GinkgoHelper()
	return crdV4Spec(parseJSON[*v4.Api](spec))
}

func crdV4Spec(spec *v4.Api) apiv4.DTO {
	GinkgoHelper()
	api := &v1alpha1.ApiV4Definition{
		ObjectMeta: metav1.ObjectMeta{Name: "api", Namespace: "default"},
		Spec:       v1alpha1.ApiV4DefinitionSpec{Api: *spec},
	}
	dto, err := apiv4.ToDTO(api)
	Expect(err).ToNot(HaveOccurred())
	return dto
}

// crdDriftV4 maps a CRD spec (JSON) the way the admission drift check does.
func crdDriftV4(spec string) apiv4.DTO {
	GinkgoHelper()
	return apiv4.DriftDTO(crdV4(spec))
}

func crdDriftV4Spec(spec *v4.Api) apiv4.DTO {
	GinkgoHelper()
	return apiv4.DriftDTO(crdV4Spec(spec))
}

// remoteDriftV4 maps a GET response (JSON) the way the admission drift check does, against the CR's DTO.
func remoteDriftV4(crd apiv4.DTO, remote string) apiv4.DTO {
	GinkgoHelper()
	return apiv4.DriftDTO(apiv4.RemoteDTO(crd, parseJSON[sdk.ApiV4State](remote)))
}

// sdkV4 maps a DTO of GKO's former v4 payload model to the SDK payload, normalised for drift. The plan order is
// compared when every plan has one, as for a CR.
func sdkV4(dto model.APIV4DTO) apiv4.DTO {
	GinkgoHelper()
	api, err := mapper.MapViaJSON[sdk.ApiV4](dto)
	Expect(err).ToNot(HaveOccurred())
	sent := apiv4.DTO{ApiV4: api}
	if len(dto.Plans) > 0 && !slices.ContainsFunc(dto.Plans, func(plan *model.APIV4PlanDTO) bool { return plan.Order == nil }) {
		sent.PlanOrder = make(map[string]apiv4.PlanRank, len(dto.Plans))
		for _, plan := range dto.Plans {
			sent.PlanOrder[plan.HRID] = apiv4.PlanRank{Order: *plan.Order}
		}
	}
	return apiv4.DriftDTO(sent)
}

// remoteDriftV4JSON maps a GET response (JSON) the way the admission drift check does, for a CR without plans.
func remoteDriftV4JSON(remote string) apiv4.DTO {
	GinkgoHelper()
	return apiv4.DriftDTO(apiv4.RemoteDTO(apiv4.DTO{}, parseJSON[sdk.ApiV4State](remote)))
}

// sdkV4JSON parses an SDK payload, normalised for drift.
func sdkV4JSON(api string) apiv4.DTO {
	GinkgoHelper()
	return apiv4.DriftDTO(apiv4.DTO{ApiV4: parseJSON[sdk.ApiV4](api)})
}

// remoteOfCRD is the remote side of the drift check for a GET response given as GKO's former v4 payload model.
func remoteOfCRD(crd apiv4.DTO, remote model.APIV4DTO) apiv4.DTO {
	GinkgoHelper()
	state, err := mapper.MapViaJSON[sdk.ApiV4State](remote)
	Expect(err).ToNot(HaveOccurred())
	return apiv4.DriftDTO(apiv4.RemoteDTO(crd, state))
}

