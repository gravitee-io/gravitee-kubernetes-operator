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

package admissionwebhook

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Webhook", labels.WithContext, func() {
	timeout := constants.EventualTimeout / 10
	interval := constants.Interval

	ctx := context.Background()

	It("should get errors for API creation because of existing path", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4).
			Build().
			Apply()

		By("Check API creation validation")
		Eventually(func() error {
			api := &v1alpha1.ApiV4Definition{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fixtures.APIv4.Name + "-duplicate",
					Namespace: fixtures.APIv4.Namespace,
				},
			}

			fixtures.APIv4.Spec.DeepCopyInto(&api.Spec)

			if err := manager.Client().Create(ctx, api); err != nil {
				return err
			}

			_, err := api.ValidateCreate()
			return err
		}, timeout, interval).ShouldNot(Succeed())
	})
})
