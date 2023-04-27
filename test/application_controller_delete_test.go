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

package test

import (
	"fmt"
	"net/http"
	"time"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	ihttp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Deleting an Application", func() {
	httpClient := http.Client{Timeout: 5 * time.Second}

	Context("Deleting an Application", func() {

		var applicationFixture *gio.Application
		var managementContextFixture *gio.ManagementContext
		var contextLookupKey types.NamespacedName
		var appLookupKey types.NamespacedName
		It("Should delete Application", func() {
			By("Initializing the Application fixture")
			fixtureGenerator := internal.NewFixtureGenerator()
			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Application: internal.BasicApplication,
				Context:     internal.ContextWithCredentialsFile,
			})
			Expect(err).ToNot(HaveOccurred())
			managementContextFixture = fixtures.Context
			applicationFixture = fixtures.Application
			contextLookupKey = types.NamespacedName{Name: managementContextFixture.Name, Namespace: namespace}
			appLookupKey = types.NamespacedName{Name: applicationFixture.Name, Namespace: namespace}

			By("Create a management context to synchronize with the REST API")
			Expect(k8sClient.Create(ctx, managementContextFixture)).Should(Succeed())

			managementContext := new(gio.ManagementContext)
			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, managementContext)
			}, timeout, interval).Should(Succeed())

			By("Create application")
			applicationFixture.Spec.Name += fixtureGenerator.Suffix
			Expect(k8sClient.Create(ctx, applicationFixture)).Should(Succeed())

			By("Getting created application and expect to find it")
			createdApplication := &gio.Application{}
			Eventually(func() error {
				return k8sClient.Get(ctx, appLookupKey, createdApplication)
			}, timeout, interval).ShouldNot(HaveOccurred())

			expectedApiName := applicationFixture.Name
			Expect(createdApplication.Name).Should(Equal(expectedApiName))

			Expect(createdApplication.Spec.Name).Should(Equal(applicationFixture.Spec.Name))
			Expect(len(*createdApplication.Spec.ApplicationMetaData)).Should(Equal(2))

			By("Call Management API and expect the Application to be available")
			var endpoint = fmt.Sprintf("%s/organizations/%s/environments/%s/applications",
				internal.ManagementUrl, managementContext.Spec.OrgId, managementContext.Spec.EnvId)

			Eventually(func() error {
				req, _ := http.NewRequest("GET", endpoint, nil)
				req.Header.Add("Authorization", "Basic "+basicAuth(managementContext))

				resp, callErr := httpClient.Do(req)
				if callErr != nil {
					return err
				}

				applications := new([]*model.Application)
				callErr = ihttp.WriteJSON(resp, applications)
				if callErr != nil {
					return err
				}

				if applications == nil || len(*applications) < 2 {
					return fmt.Errorf("application hasn't been created")
				}

				for _, app := range *applications {
					if app.Name == applicationFixture.Spec.Name {
						return nil
					}
				}

				return fmt.Errorf("can't find any application with the given name %s", applicationFixture.Spec.Name)
			}, timeout, interval).ShouldNot(HaveOccurred())

			Eventually(func() error {
				By("Deleting the application")

				savedApplication := new(gio.Application)
				if err = k8sClient.Get(ctx, appLookupKey, savedApplication); err != nil {
					return err
				}

				return k8sClient.Delete(ctx, savedApplication)
			}).Should(Succeed())

			By("Getting the deleted application")
			deletedApplication := &gio.Application{}
			Eventually(func() error {
				return k8sClient.Get(ctx, appLookupKey, deletedApplication)
			}, timeout, interval).Should(HaveOccurred())

		})
	})

})
