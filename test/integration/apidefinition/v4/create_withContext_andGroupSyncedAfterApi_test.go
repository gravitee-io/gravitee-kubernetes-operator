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

package v4

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/notification"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Create", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	// The group references a management context that does not exist yet, so the group cannot sync
	// to APIM until syncGroup creates that context.
	createUnsyncedGroup := func(groupContext *v1alpha1.ManagementContext) *v1alpha1.Group {
		name := random.GetName()
		apiGroup := &v1alpha1.Group{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: constants.Namespace},
			Spec: v1alpha1.GroupSpec{
				Type:    &group.Type{Name: name, Members: []group.Member{}},
				Context: groupContext.GetNamespacedName(),
			},
		}
		Expect(manager.Client().Create(ctx, apiGroup)).To(Succeed())
		return apiGroup
	}

	syncGroup := func(groupContext *v1alpha1.ManagementContext, apiGroup *v1alpha1.Group) {
		Expect(manager.Client().Create(ctx, groupContext)).To(Succeed())
		Eventually(func() error {
			if err := manager.GetLatest(ctx, apiGroup); err != nil {
				return err
			}
			return assert.GroupCompleted(apiGroup)
		}, timeout, interval).Should(Succeed(), apiGroup.Name)
	}

	It("should add the group to the API in APIM once the group is synced after the API", func() {
		fixtures := fixture.
			Builder().
			WithContext(constants.ContextWithCredentialsFile).
			WithAPIv4(constants.ApiV4WithContextFile).
			Build()

		By("referencing a group whose management context does not exist yet")

		groupContext := fixtures.Context.DeepCopy()
		groupContext.Name += "-group"
		apiGroup := createUnsyncedGroup(groupContext)

		fixtures.APIv4.Spec.GroupRefs = []refs.NamespacedName{refs.NewNamespacedName(apiGroup.Namespace, apiGroup.Name)}
		fixtures = fixtures.Apply()

		apim := apim.NewClient(ctx)
		apiHRID := refs.NewNamespacedNameFromObject(fixtures.APIv4).HRID()
		groupInAPIM := BeElementOf(
			model.APIGroup(refs.NewNamespacedNameFromObject(apiGroup).HRID()),
			model.APIGroup(apiGroup.Spec.Name),
		)
		groupsOfAPIInAPIM := func() ([]model.APIGroup, error) {
			api, err := apim.APIs.GetV4ByHRID(apiHRID)
			if err != nil {
				return nil, err
			}
			return api.Groups, nil
		}

		By("expecting APIM to have dropped the group it does not know yet")

		Expect(groupsOfAPIInAPIM()).NotTo(ContainElement(groupInAPIM))

		By("creating the group management context, so the group syncs to APIM")

		syncGroup(groupContext, apiGroup)

		By("expecting the API to be synced again, with the group")

		Eventually(groupsOfAPIInAPIM, timeout, interval).Should(ContainElement(groupInAPIM), fixtures.APIv4.Name)
	})

	It("should add the group to the API console notification in APIM once the group is synced after the API", func() {
		fixtures := fixture.
			Builder().
			WithContext(constants.ContextWithCredentialsFile).
			WithAPIv4(constants.ApiV4WithContextFile).
			Build()

		By("referencing a group whose management context does not exist yet, from the API and its console notification")

		groupContext := fixtures.Context.DeepCopy()
		groupContext.Name += "-group"
		apiGroup := createUnsyncedGroup(groupContext)
		groupRef := refs.NewNamespacedName(apiGroup.Namespace, apiGroup.Name)

		consoleNotification := &v1alpha1.Notification{
			ObjectMeta: metav1.ObjectMeta{Name: random.GetName(), Namespace: constants.Namespace},
			Spec: v1alpha1.NotificationSpec{
				Type: &notification.Type{
					Target:    notification.TargetConsole,
					EventType: notification.EventTypeAPI,
					Console: notification.Console{
						APIEvents: []notification.ApiEvent{"API_STARTED"},
						GroupRefs: []refs.NamespacedName{groupRef},
					},
				},
			},
		}
		Expect(manager.Client().Create(ctx, consoleNotification)).To(Succeed())

		fixtures.APIv4.Spec.GroupRefs = []refs.NamespacedName{groupRef}
		fixtures.APIv4.Spec.NotificationsRefs = []refs.NamespacedName{
			refs.NewNamespacedName(consoleNotification.Namespace, consoleNotification.Name),
		}
		fixtures = fixtures.Apply()

		apim := apim.NewClient(ctx)
		apiHRID := refs.NewNamespacedNameFromObject(fixtures.APIv4).HRID()
		consoleNotificationGroupsInAPIM := func() ([]string, error) {
			api, err := apim.APIs.GetV4ByHRID(apiHRID)
			if err != nil {
				return nil, err
			}
			console, err := apim.Notification.GetConsoleNotificationConfiguration(api.ID)
			if err != nil || console == nil {
				return nil, err
			}
			return console.Groups, nil
		}

		By("expecting APIM to have dropped the group it does not know yet")

		Expect(consoleNotificationGroupsInAPIM()).To(BeEmpty())

		By("creating the group management context, so the group syncs to APIM")

		syncGroup(groupContext, apiGroup)

		By("expecting the API to be synced again, with the group in its console notification")

		Eventually(consoleNotificationGroupsInAPIM, timeout, interval).
			Should(ContainElement(apiGroup.Status.ID), fixtures.APIv4.Name)
	})
})
