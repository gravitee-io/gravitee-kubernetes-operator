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
	"errors"

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

	It("should add group reference to API", func() {
		fixtures := fixture.
			Builder().
			WithContext(constants.ContextWithCredentialsFile).
			WithGroup(constants.GroupFile).
			WithAPIv4(constants.ApiV4WithContextFile).
			Build()

		By("initializing a service account in current organization")

		apim := apim.NewClient(ctx)
		saName := random.GetName()
		Expect(apim.Org.CreateUser(model.NewServiceAccount(saName))).To(Succeed())

		By("adding the sa to the Group")

		groupMember := base.NewGraviteeGroupMember(saName, "OWNER")
		name := random.GetName()
		group := v1alpha1.Group{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: fixtures.APIv4.Namespace,
			},
			Spec: v1alpha1.GroupSpec{
				Type: &group.Type{
					Name:          name,
					NotifyMembers: false,
					Members:       []group.Member{groupMember},
				},
				Context: &refs.NamespacedName{
					Name:      fixtures.Context.Name,
					Namespace: fixtures.Context.Namespace,
				},
			},
		}

		fixtures.APIv4.Spec.GroupRefs = []refs.NamespacedName{refs.NewNamespacedName(
			group.Namespace,
			group.Name,
		)}

		fixtures = fixtures.Apply()

		Eventually(func() error {
			api := &v1alpha1.ApiV4Definition{}
			err := manager.Client().Get(ctx,
				types.NamespacedName{
					Name:      fixtures.APIv4.Name,
					Namespace: fixtures.APIv4.Namespace,
				},
				api)
			if err != nil {
				return err
			}

			if len(api.Status.Conditions) != 2 {
				return errors.New("expected exactly two conditions")
			}

			for _, condition := range api.Status.Conditions {
				if condition.Type == "ResolvedRefs" && condition.Status == "False" {
					return nil
				}
			}

			return errors.New("expected ResolvedRefs condition to be false")
		}, timeout, interval).Should(Succeed(), fixtures.APIv4.Name)

		By("creating the missing group, ResolveRefs condition must be set true")
		Expect(manager.Client().Create(ctx, &group)).To(Succeed())

		Eventually(func() error {
			api := &v1alpha1.ApiV4Definition{}
			err := manager.Client().Get(ctx,
				types.NamespacedName{
					Name:      fixtures.APIv4.Name,
					Namespace: fixtures.APIv4.Namespace,
				},
				api)
			if err != nil {
				return err
			}

			if len(api.Status.Conditions) != 2 {
				return errors.New("expected exactly two conditions")
			}

			for _, condition := range api.Status.Conditions {
				if condition.Type == "ResolvedRefs" && condition.Status == "True" {
					return nil
				}
			}

			return errors.New("expected ResolvedRefs condition to be true")
		}, timeout, interval).Should(Succeed(), fixtures.APIv4.Name)

	})
})
