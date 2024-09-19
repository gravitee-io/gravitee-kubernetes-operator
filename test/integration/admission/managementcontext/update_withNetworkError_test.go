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

package managementcontext

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/management"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/mctx"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Validate update", labels.WithContext, func() {
	interval := constants.Interval
	admissionCtrl := mctx.AdmissionCtrl{}
	ctx := context.Background()
	cli := manager.Client()

	It("should return warnings if APIM is not accessible", func() {
		mCtx := &v1alpha1.ManagementContext{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "webhook-ctx" + random.GetSuffix(),
				Namespace: "default",
			},
			Spec: v1alpha1.ManagementContextSpec{
				Context: &management.Context{
					BaseUrl: "https://gko.example.com",
					EnvID:   "DEFAULT",
					OrgID:   "DEFAULT",
					Auth: &management.Auth{
						BearerToken: "test",
					},
				},
			},
		}

		Expect(cli.Create(context.Background(), mCtx)).To(Succeed())

		Eventually(func() error {
			return cli.Get(context.Background(), types.NamespacedName{
				Namespace: mCtx.Namespace,
				Name:      mCtx.Name,
			}, mCtx)
		}, constants.EventualTimeout, interval).Should(Succeed())

		Consistently(func() error {
			warnings, _ := admissionCtrl.ValidateUpdate(ctx, mCtx, mCtx)
			return assert.SliceOfSize("warnings", warnings, 1)
		}, constants.ConsistentTimeout, interval).Should(Succeed())
	})
})
