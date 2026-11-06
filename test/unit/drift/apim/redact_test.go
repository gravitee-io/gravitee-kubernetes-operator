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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Redacting secret values from a drift report", func() {
	DescribeTable("hides the values of a drifted field holding a secret, on either side",
		func(crdSecret, remoteSecret string, secrets []string, hidden, shown []string) {
			merged := mergedAPIDrift(apiWithClientSecret(crdSecret), apiWithClientSecret(remoteSecret))
			merged.Redact(secrets)
			report := merged.String()

			Expect(merged.DriftDetected()).To(BeTrue(), "redaction never hides a drift")
			Expect(report).To(ContainSubstring("clientSecret"))
			for _, value := range hidden {
				Expect(report).NotTo(ContainSubstring(value))
			}
			for _, value := range shown {
				Expect(report).To(ContainSubstring(value))
			}
		},
		Entry("a secret on the CRD side hides both values",
			"rotated-secret", "stored-in-apim", []string{"rotated-secret"},
			[]string{"rotated-secret", "stored-in-apim"}, []string{"(secret value, not shown)"}),
		Entry("a secret on the remote side hides both values",
			"plain-now", "previous-secret", []string{"previous-secret"},
			[]string{"plain-now", "previous-secret"}, []string{"(secret value, not shown)"}),
		Entry("a secret embedded in a longer value hides it",
			"Bearer rotated-secret", "Bearer stored", []string{"rotated-secret"},
			[]string{"rotated-secret", "Bearer stored"}, []string{"(secret value, not shown)"}),
		Entry("a value holding no secret is printed",
			"plain-crd", "plain-remote", []string{"unrelated-secret"},
			[]string{"(secret value, not shown)"}, []string{`"plain-crd" != "plain-remote"`}),
		Entry("no secret at all leaves the report as it was",
			"plain-crd", "plain-remote", nil,
			[]string{"(secret value, not shown)"}, []string{`"plain-crd" != "plain-remote"`}),
	)

	It("hides only the fields holding a secret", func() {
		crd := apiWithClientSecret("rotated-secret")
		crd.Description = new("described in the CRD")
		remote := apiWithClientSecret("stored-in-apim")
		remote.Description = new("described in APIM")

		merged := mergedAPIDrift(crd, remote)
		merged.Redact([]string{"rotated-secret"})
		report := merged.String()

		Expect(report).To(ContainSubstring("clientSecret: (secret value, not shown)"))
		Expect(report).To(ContainSubstring(`"described in the CRD" != "described in APIM"`))
	})

	It("does not turn equivalent values into a drift", func() {
		merged := mergedAPIDrift(apiWithClientSecret("same-secret"), apiWithClientSecret("same-secret"))
		merged.Redact([]string{"same-secret"})

		Expect(merged.DriftDetected()).To(BeFalse())
	})
})

// mergedAPIDrift runs the comparison the admission path runs: old and new CRD against the remote.
func mergedAPIDrift(crd, remote model.APIV4DTO) drift.Result {
	return drift.Merge(
		drift.DetectWithNamespace(crd, remote, ""),
		drift.DetectWithNamespace(crd, remote, ""),
	)
}

func apiWithClientSecret(secret string) model.APIV4DTO {
	return model.APIV4DTO{
		Resources: []*model.APIV4ResourceDTO{{
			Enabled: true,
			Name:    new("introspection"),
			Type:    new("oauth2-am-resource"),
			Configuration: utils.NewGenericStringMap().
				Put("clientId", "client").
				Put("clientSecret", secret),
		}},
	}
}
