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

package group

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Create", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("should create group with default roles mapping", func() {
		fixtures := fixture.
			Builder().
			WithContext(constants.ContextWithCredentialsFile).
			Build()

		By("initializing a service account in current organization")

		apim := apim.NewClient(ctx)
		saName := random.GetName()
		Expect(apim.Org.CreateUser(model.NewServiceAccount(saName))).To(Succeed())

		By("creating a group with roles mapping")

		groupMember := base.NewGraviteeGroupMember(saName, "OWNER")
		name := random.GetName()
		group := v1alpha1.Group{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: constants.Namespace,
			},
			Spec: v1alpha1.GroupSpec{
				Type: &group.Type{
					Name:          name,
					NotifyMembers: false,
					Members:       []group.Member{groupMember},
					Roles: &group.Roles{
						API:         "OWNER",
						Application: "USER",
					},
				},
				Context: &refs.NamespacedName{
					Name:      fixtures.Context.Name,
					Namespace: fixtures.Context.Namespace,
				},
			},
		}

		fixtures.Apply()
		Expect(manager.Client().Create(ctx, &group)).To(Succeed())

		Eventually(func() error {
			g := &v1alpha1.Group{}
			err := manager.Client().Get(ctx,
				types.NamespacedName{
					Name:      group.Name,
					Namespace: group.Namespace,
				},
				g)
			if err != nil {
				return err
			}

			// Verify that group was created and status is populated
			if g.Status.ID == "" {
				return ErrStatusNotPopulated
			}

			return nil
		}, timeout, interval).Should(Succeed(), group.Name)

		By("cleaning up the created group")
		Expect(manager.Delete(ctx, &group)).To(Succeed())
	})

	It("should create group with partial roles mapping", func() {
		fixtures := fixture.
			Builder().
			WithContext(constants.ContextWithCredentialsFile).
			Build()

		By("initializing a service account in current organization")

		apim := apim.NewClient(ctx)
		saName := random.GetName()
		Expect(apim.Org.CreateUser(model.NewServiceAccount(saName))).To(Succeed())

		By("creating a group with only API role mapping")

		groupMember := base.NewGraviteeGroupMember(saName, "OWNER")
		name := random.GetName()
		group := v1alpha1.Group{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: constants.Namespace,
			},
			Spec: v1alpha1.GroupSpec{
				Type: &group.Type{
					Name:          name,
					NotifyMembers: false,
					Members:       []group.Member{groupMember},
					Roles: &group.Roles{
						API: "OWNER",
					},
				},
				Context: &refs.NamespacedName{
					Name:      fixtures.Context.Name,
					Namespace: fixtures.Context.Namespace,
				},
			},
		}

		fixtures.Apply()
		Expect(manager.Client().Create(ctx, &group)).To(Succeed())

		Eventually(func() error {
			g := &v1alpha1.Group{}
			err := manager.Client().Get(ctx,
				types.NamespacedName{
					Name:      group.Name,
					Namespace: group.Namespace,
				},
				g)
			if err != nil {
				return err
			}

			// Verify that group was created and status is populated
			if g.Status.ID == "" {
				return ErrStatusNotPopulated
			}

			return nil
		}, timeout, interval).Should(Succeed(), group.Name)

		By("cleaning up the created group")
		Expect(manager.Delete(ctx, &group)).To(Succeed())
	})

	It("should create group without roles mapping", func() {
		fixtures := fixture.
			Builder().
			WithContext(constants.ContextWithCredentialsFile).
			WithGroup(constants.GroupFile).
			Build()

		By("initializing a service account in current organization")

		apim := apim.NewClient(ctx)
		saName := random.GetName()
		Expect(apim.Org.CreateUser(model.NewServiceAccount(saName))).To(Succeed())

		By("adding the sa to the Group")

		groupMember := base.NewGraviteeGroupMember(saName, "OWNER")
		fixtures.Group.Spec.Members = []group.Member{groupMember}

		fixtures.Apply()

		Eventually(func() error {
			g := &v1alpha1.Group{}
			err := manager.Client().Get(ctx,
				types.NamespacedName{
					Name:      fixtures.Group.Name,
					Namespace: fixtures.Group.Namespace,
				},
				g)
			if err != nil {
				return err
			}

			// Verify that group was created and status is populated
			if g.Status.ID == "" {
				return ErrStatusNotPopulated
			}

			return nil
		}, timeout, interval).Should(Succeed(), fixtures.Group.Name)

		By("cleaning up the created group")
		Expect(manager.Delete(ctx, fixtures.Group)).To(Succeed())
	})
})

var ErrStatusNotPopulated = StatusNotPopulatedError{Message: "status not populated"}

type StatusNotPopulatedError struct {
	Message string
}

func (e StatusNotPopulatedError) Error() string {
	return e.Message
}
