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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1beta1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Update an Application", func() {
	var applicationFixture *v1beta1.Application
	var managementContextFixture *v1beta1.ManagementContext
	var contextLookupKey types.NamespacedName
	var appLookupKey types.NamespacedName

	It("Should create Application with context", func() {
		By("Initializing the Application fixture")
		fixtureGenerator := internal.NewFixtureGenerator()
		fixtures, fErr := fixtureGenerator.NewFixtures(internal.FixtureFiles{
			Application: internal.BasicApplication,
			Context:     internal.ClusterContextFile,
		})
		Expect(fErr).ToNot(HaveOccurred())
		managementContextFixture = fixtures.Context
		applicationFixture = fixtures.Application
		contextLookupKey = types.NamespacedName{Name: managementContextFixture.Name, Namespace: namespace}
		appLookupKey = types.NamespacedName{Name: applicationFixture.Name, Namespace: namespace}

		By("Creating a management context to synchronize with the REST API")
		Expect(k8sClient.Create(ctx, managementContextFixture)).Should(Succeed())

		managementContext := new(v1beta1.ManagementContext)
		Eventually(func() error {
			return k8sClient.Get(ctx, contextLookupKey, managementContext)
		}, timeout, interval).Should(Succeed())

		By("Creating an application")
		applicationFixture.Spec.Name += fixtureGenerator.Suffix
		Expect(k8sClient.Create(ctx, applicationFixture)).Should(Succeed())

		By("Getting created application and expect to find it")
		createdApplication := &v1beta1.Application{}
		Eventually(func() error {
			if fErr = k8sClient.Get(ctx, appLookupKey, createdApplication); fErr != nil {
				return fErr
			}
			return internal.AssertApplicationStatusIsSet(createdApplication)
		}, timeout, interval).ShouldNot(HaveOccurred())

		expectedApplicationName := applicationFixture.Name
		Expect(createdApplication.Name).Should(Equal(expectedApplicationName))

		Expect(createdApplication.Spec.Name).Should(Equal(applicationFixture.Spec.Name))
		Expect(len(*createdApplication.Spec.ApplicationMetaData)).Should(Equal(2))

		By("Calling Management API and expect the Application to be available")
		Eventually(func() error {
			apimClient, err := internal.NewAPIM(ctx)
			Expect(err).ToNot(HaveOccurred())
			app, err := apimClient.Applications.GetByID(createdApplication.Status.ID)
			Expect(err).ToNot(HaveOccurred())
			if err != nil {
				return err
			}
			return internal.AssertEquals("name", app.Name, createdApplication.Spec.Name)
		}, timeout, interval).ShouldNot(HaveOccurred())

		By("Updating application")
		updatedApplication := createdApplication.DeepCopy()
		updatedApplication.Spec.Name += "-updated"

		Eventually(func() error {
			update := new(v1beta1.Application)
			if fErr = k8sClient.Get(ctx, appLookupKey, update); fErr != nil {
				return fErr
			}
			updatedApplication.Spec.DeepCopyInto(&update.Spec)
			return k8sClient.Update(ctx, update)
		}, timeout, interval).ShouldNot(HaveOccurred())

		By("Calling Management API and expect the updated Application to be available")
		Eventually(func() error {
			apimClient, err := internal.NewAPIM(ctx)
			Expect(err).ToNot(HaveOccurred())
			app, err := apimClient.Applications.GetByID(updatedApplication.Status.ID)
			Expect(err).ToNot(HaveOccurred())
			if err != nil {
				return err
			}
			return internal.AssertEquals("name", app.Name, updatedApplication.Spec.Name)
		}, timeout, interval).ShouldNot(HaveOccurred())

	})
})
