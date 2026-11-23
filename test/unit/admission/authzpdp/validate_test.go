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

package authzpdp_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/authzpdp"
)

var _ = Describe("Validating an AuthzPdp", func() {
	ctx := context.Background()
	ctrl := admission.NewAdmissionCtrl()
	hridFinding := "hrid [travel-orders.eu] derived from the resource namespace and name must match"

	// A Kubernetes name may hold a dot, which the automation API refuses in a hrid.
	dotted := func() *v1alpha1.AuthzPdp {
		return &v1alpha1.AuthzPdp{
			ObjectMeta: metav1.ObjectMeta{Name: "orders.eu", Namespace: "travel"},
			Spec: v1alpha1.AuthzPdpSpec{
				Context: &refs.NamespacedName{Name: "dev-ctx"},
				Type:    authzpdp.Type{Name: "Orders PDP (EU)", TargetPdpID: "orders", Tag: "eu"},
			},
		}
	}

	BeforeEach(func() {
		platformCalls.Store(0)
	})

	It("refuses a create whose namespace and name give an invalid hrid, before the dry run", func() {
		_, err := ctrl.ValidateCreate(ctx, dotted())

		Expect(err).To(MatchError(ContainSubstring(hridFinding)))
		Expect(platformCalls.Load()).To(BeZero())
	})

	It("refuses an update whose namespace and name give an invalid hrid, before the dry run", func() {
		oldObj := dotted()
		newObj := oldObj.DeepCopy()
		newObj.Spec.Name = "Orders PDP (Europe)"

		_, err := ctrl.ValidateUpdate(ctx, oldObj, newObj)

		Expect(err).To(MatchError(ContainSubstring(hridFinding)))
		Expect(platformCalls.Load()).To(BeZero())
	})
})
