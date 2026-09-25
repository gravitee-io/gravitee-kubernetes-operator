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

package apim_test

import (
	"encoding/json"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Application payload", func() {
	const templateID = "3d5a2f0e-1c4b-4c8a-9f7e-2b6d8e0a1c34"

	newSpec := func() v1alpha1.ApplicationSpec {
		return v1alpha1.ApplicationSpec{Application: application.Application{
			Settings: &application.Setting{Oauth: &application.OAuthClientSettings{
				ApplicationType:          "BACKEND_TO_BACKEND",
				GrantTypes:               []application.GrantType{application.GrantTypeClientCredentials},
				AdditionalClientMetadata: map[string]string{"software_id": templateID},
			}},
		}}
	}

	It("sends the OAuth additional client metadata to APIM", func() {
		payload, err := json.Marshal(newSpec())

		Expect(err).ToNot(HaveOccurred())
		var sent struct {
			Settings struct {
				Oauth struct {
					AdditionalClientMetadata map[string]string `json:"additionalClientMetadata"`
				} `json:"oauth"`
			} `json:"settings"`
		}
		Expect(json.Unmarshal(payload, &sent)).To(Succeed())
		Expect(sent.Settings.Oauth.AdditionalClientMetadata).To(Equal(map[string]string{"software_id": templateID}))
	})

	It("compares the OAuth additional client metadata for drift", func() {
		dto := model.ToApplicationDTO(newSpec())

		Expect(dto.Settings.Oauth.AdditionalClientMetadata).To(Equal(map[string]string{"software_id": templateID}))
	})
})
