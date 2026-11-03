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
	"fmt"
	"net/http"
	"slices"

	amsdk "github.com/gravitee-io/gravitee-automation-sdk/am-sdk/v2/pkg/sdk"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	appmodel "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/am/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Update", func() {
	ctx := context.Background()

	It("should update the application in AM, keeping its clientId", func() {
		app := fixtures(constants.AMApplicationMinimalFile).Apply().AMApplication
		clientID := app.Status.ClientID

		By("renaming the application, changing its type, its redirect URIs and adding metadata")

		Expect(manager.GetLatest(ctx, app)).To(Succeed())
		app.Spec.Name = "Renamed"
		app.Spec.Type = "web"
		app.Spec.Metadata = utils.NewGenericStringMap().Put("owner", "web-team")
		app.Spec.Settings = &appmodel.AutomationApplicationSettings{
			Oauth: &appmodel.AutomationApplicationOAuthSettings{RedirectUris: []string{"https://example.com/callback"}},
		}
		Expect(manager.Client().Update(ctx, app)).To(Succeed())

		By("expecting the changes in AM and the same clientId")

		expectRemote(ctx, app, func(remote *amsdk.Application) error {
			switch {
			case remote.Name != "Renamed":
				return fmt.Errorf("name %s", remote.Name)
			case remote.Type != amsdk.ApplicationType("web"):
				return fmt.Errorf("type %s", remote.Type)
			case remote.Metadata["owner"] != "web-team":
				return fmt.Errorf("metadata %v", remote.Metadata)
			case remote.Settings == nil || remote.Settings.Oauth == nil ||
				!slices.Equal(remote.Settings.Oauth.RedirectUris, []string{"https://example.com/callback"}):
				return fmt.Errorf("settings %v", remote.Settings)
			case *remote.ClientId != clientID:
				return fmt.Errorf("clientId %s", *remote.ClientId)
			}
			return nil
		})
		Expect(manager.GetLatest(ctx, app)).To(Succeed())
		Expect(app.Status.ClientID).To(Equal(clientID))
	})

	It("should be rejected by AM when the type changes to agent (AC4)", func() {
		app := fixtures(constants.AMApplicationMinimalFile).Apply().AMApplication

		Expect(manager.GetLatest(ctx, app)).To(Succeed())
		app.Spec.Type = "agent"
		Expect(manager.Client().Update(ctx, app)).To(Succeed())

		expectRejected(ctx, app, "type cannot be changed to or from agent")
		expectRemote(ctx, app, func(remote *amsdk.Application) error {
			if remote.Type != amsdk.ApplicationType("native") {
				return fmt.Errorf("type %s", remote.Type)
			}
			return nil
		})
	})

	It("should send an inherited block when declared and drop it when removed", func() {
		app := fixtures(constants.AMApplicationMinimalFile).Apply().AMApplication

		By("declaring a login block")

		Expect(manager.GetLatest(ctx, app)).To(Succeed())
		app.Spec.Settings = &appmodel.AutomationApplicationSettings{
			Login: &appmodel.LoginSettings{Inherited: new(false), ForgotPasswordEnabled: new(true)},
		}
		Expect(manager.Client().Update(ctx, app)).To(Succeed())

		expectRemote(ctx, app, func(remote *amsdk.Application) error {
			if remote.Settings == nil || remote.Settings.Login == nil ||
				remote.Settings.Login.ForgotPasswordEnabled == nil || !*remote.Settings.Login.ForgotPasswordEnabled {
				return fmt.Errorf("settings %v", remote.Settings)
			}
			return nil
		})

		By("removing it")

		Expect(manager.GetLatest(ctx, app)).To(Succeed())
		app.Spec.Settings = nil
		Expect(manager.Client().Update(ctx, app)).To(Succeed())

		expectRemote(ctx, app, func(remote *amsdk.Application) error {
			if remote.Settings != nil && remote.Settings.Login != nil {
				return fmt.Errorf("login %v", remote.Settings.Login)
			}
			return nil
		})
	})
})

// expectRemote waits until the application AM holds satisfies check.
func expectRemote(ctx context.Context, app *v1alpha1.AMApplication, check func(*amsdk.Application) error) {
	GinkgoHelper()
	Eventually(ctx, func() error {
		status, remote := remoteApplication(ctx, app)
		if status != http.StatusOK {
			return fmt.Errorf("status %d", status)
		}
		return check(remote)
	}, constants.EventualTimeout, constants.Interval).Should(Succeed(), app.Name)
}
