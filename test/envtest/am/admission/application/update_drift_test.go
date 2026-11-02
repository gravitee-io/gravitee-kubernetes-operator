// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package amapplication

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	amsdk "github.com/gravitee-io/gravitee-automation-sdk/am-sdk/v2/pkg/sdk"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/amapplication"
	internal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
)

var _ = Describe("Validate drift", func() {
	ctx := context.Background()
	admissionCtrl := amapplication.NewAdmissionCtrl()

	// inAM applies everything but the application, and stores the application in the AM mock as AM returns it:
	// the spec, with the fields remote overrides.
	inAM := func(remote func(*internal.Application)) *v1alpha1.AMApplication {
		app := applyAllButApplication(constants.AMApplicationFullFile).AMApplication

		dto, err := internal.ToApplicationDTO(app)
		Expect(err).ToNot(HaveOccurred())
		remote(&dto)
		resp, err := am.NewSDKClient().UpsertApplicationWithResponse(ctx, dto.DomainKey, nil, dto.Application)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.JSON200).ToNot(BeNil(), string(resp.Body))
		// synced as the controller leaves it: drift only runs for a resource AM accepted
		app.Status.Key = resp.JSON200.Key
		return app
	}
	validateUpdate := func(app *v1alpha1.AMApplication) func() error {
		return func() error {
			updated := app.DeepCopy()
			updated.Annotations = map[string]string{"updated": "true"}
			_, err := admissionCtrl.ValidateUpdate(ctx, app, updated)
			return err
		}
	}

	It("should not drift on the clientId and OAuth settings AM computed", func() {
		app := inAM(func(remote *internal.Application) {
			remote.ClientId = new("generated-by-am")
			remote.Settings.Oauth.ResponseTypes = []string{"code"}
			remote.Settings.Oauth.TokenEndpointAuthMethod = new("none")
			remote.Settings.Oauth.ForcePKCE = new(true)
			remote.Settings.Login = &amsdk.LoginSettings{Inherited: new(true)}
		})

		Consistently(validateUpdate(app), constants.ConsistentTimeout, constants.Interval).Should(Succeed())
	})

	It("should detect drift on the name", func() {
		app := inAM(func(remote *internal.Application) {
			remote.Name = "Renamed in AM"
		})

		Eventually(validateUpdate(app), constants.EventualTimeout, constants.Interval).
			Should(MatchError(And(ContainSubstring("drift detected"), ContainSubstring("Renamed in AM"))))
	})
})
