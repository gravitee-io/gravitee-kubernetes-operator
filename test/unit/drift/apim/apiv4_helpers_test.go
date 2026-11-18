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
	return apiv4.DriftDTO(apiv4.RemoteDTO(crd, withEchoedBooleans(crd, parseJSON[sdk.ApiV4State](remote))))
}

// sdkV4 maps a DTO of GKO's former v4 payload model to the SDK payload, normalised for drift.
func sdkV4(dto model.APIV4DTO) apiv4.DTO {
	GinkgoHelper()
	api, err := mapper.MapViaJSON[sdk.ApiV4](dto)
	Expect(err).ToNot(HaveOccurred())
	return apiv4.DriftDTO(apiv4.DTO{ApiV4: api})
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
	return apiv4.DriftDTO(apiv4.RemoteDTO(crd, withEchoedBooleans(crd, state)))
}

// withEchoedBooleans fills the booleans a test GET response leaves out with the value the CR sends. GKO
// always sends them and APIM's GET always returns them (the v4drift integration suite syncs full APIs
// without drift on them), so an omitted one stands for "what GKO sent", not for the SDK default.
func withEchoedBooleans(crd apiv4.DTO, state sdk.ApiV4State) sdk.ApiV4State {
	if state.NotifyMembers == nil {
		state.NotifyMembers = crd.NotifyMembers
	}
	if state.FlowExecution != nil && state.FlowExecution.MatchRequired == nil && crd.FlowExecution != nil {
		flowExecution := *state.FlowExecution
		flowExecution.MatchRequired = crd.FlowExecution.MatchRequired
		state.FlowExecution = &flowExecution
	}
	sent := make(map[string]sdk.PageV4, len(crd.Pages))
	for _, page := range crd.Pages {
		sent[page.Hrid] = page
	}
	pages := make([]sdk.PageV4, len(state.Pages))
	for i, page := range state.Pages {
		if crdPage, ok := sent[page.Hrid]; ok {
			if page.Homepage == nil {
				page.Homepage = crdPage.Homepage
			}
			if page.Published == nil {
				page.Published = crdPage.Published
			}
		}
		pages[i] = page
	}
	if state.Pages != nil {
		state.Pages = pages
	}
	return state
}

// completeAPIV4Payload is the complete fixture as sent to APIM, before DriftDTO drops the fields
// drift leaves out (listener path overrideAccess), so that every SDK property is walked.
func completeAPIV4Payload() apiv4.DTO {
	GinkgoHelper()
	api, err := mapper.MapViaJSON[sdk.ApiV4](loadFixture[model.APIV4DTO]("api_v4_dto.json"))
	Expect(err).ToNot(HaveOccurred())
	return apiv4.DTO{ApiV4: api, PlanOrder: map[string]int{"fixture": 1}}
}
