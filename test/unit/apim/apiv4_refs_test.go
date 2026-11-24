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

package apim_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apidefinition"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apiv4"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// registerGraviteeK8s registers a fake Kubernetes client that knows the gravitee.io kinds, holding objects.
func registerGraviteeK8s(objects ...client.Object) {
	GinkgoHelper()
	scheme := runtime.NewScheme()
	Expect(clientgoscheme.AddToScheme(scheme)).To(Succeed())
	Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
	k8s.RegisterClient(fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).Build())
}

func developersGroup() *v1alpha1.Group {
	return &v1alpha1.Group{
		ObjectMeta: metav1.ObjectMeta{Name: "developers", Namespace: "default"},
		Spec:       v1alpha1.GroupSpec{Type: &group.Type{Name: "Developers"}},
	}
}

// withoutContext builds an API without a management context (no console notification resolution needed).
func withoutContext(spec string) *v1alpha1.ApiV4Definition {
	GinkgoHelper()
	api := apiV4(spec)
	api.Spec.Context = nil
	return api
}

var _ = Describe("ApiV4Definition reference resolution", func() {
	ctx := context.Background()

	DescribeTable("words a missing reference as the admission webhook does, and stays a NotFound",
		func(spec, message string) {
			registerGraviteeK8s()
			err := apiv4.ResolveRefs(ctx, withoutContext(spec), "default")

			Expect(err).To(MatchError(message))
			Expect(apierrors.IsNotFound(err)).To(BeTrue())
			Expect(apierrors.IsNotFound(gerrors.NewResolveRefError(err))).To(BeTrue())
		},
		Entry("ApiResource",
			`{"name": "orders", "version": "1", "type": "PROXY", "resources": [{"ref": {"name": "cache", "namespace": "default"}}]}`,
			"api references resource [default/cache] that does not exist in the cluster",
		),
		Entry("Shared Policy Group",
			`{"name": "orders", "version": "1", "type": "PROXY", "flows": [{"name": "all", "enabled": true, `+
				`"request": [{"enabled": true, "sharedPolicyGroupRef": {"name": "spg"}}]}]}`,
			"unable to get Shared Policy Group [spg] in namespace [default]",
		),
	)

	It("words a missing notification as the admission webhook does, and stays a NotFound", func() {
		registerGraviteeK8s()
		api := apiV4WithContext()
		api.Spec.NotificationsRefs = withoutContext(`{"notificationsRefs": [{"name": "alerts", "namespace": "default"}]}`).
			Spec.NotificationsRefs

		err := apiv4.ResolveRefsForAdmission(ctx, api, "default")

		Expect(err).To(MatchError("api references notification [default/alerts] that does not exist in the cluster"))
		Expect(apierrors.IsNotFound(err)).To(BeTrue())
	})

	It("leaves the v2 resolvers' errors as they were", func() {
		registerGraviteeK8s()
		api := withoutContext(`{"resources": [{"ref": {"name": "cache", "namespace": "default"}}]}`)

		err := apidefinition.ResolveResources(ctx, api.Spec.Resources, false)

		Expect(err).To(MatchError(ContainSubstring(`"cache" not found`)))
		Expect(err.Error()).ToNot(ContainSubstring("that does not exist in the cluster"))
	})

	It("lists a group once however many times the API is resolved", func() {
		registerGraviteeK8s(developersGroup())
		api := withoutContext(`{"name": "orders", "version": "1", "type": "PROXY", "groups": ["Developers"], ` +
			`"groupRefs": [{"name": "developers"}]}`)

		Expect(apiv4.ResolveRefsForAdmission(ctx, api, "default")).To(Succeed())
		Expect(apiv4.ResolveRefsForAdmission(ctx, api, "default")).To(Succeed())

		Expect(api.Spec.Groups).To(Equal([]string{"Developers"}))
	})

	It("reports a missing group on ResolvedRefs without failing", func() {
		registerGraviteeK8s()
		api := withoutContext(`{"name": "orders", "version": "1", "type": "PROXY", "groupRefs": [{"name": "developers"}]}`)

		Expect(apiv4.ResolveRefs(ctx, api, "default")).To(Succeed())

		Expect(k8s.HasUnresolvedRefs(api)).To(BeTrue())
	})

	It("resolves nothing for an API being deleted", func() {
		registerGraviteeK8s()
		api := withoutContext(`{"name": "orders", "version": "1", "type": "PROXY", ` +
			`"resources": [{"ref": {"name": "cache", "namespace": "default"}}]}`)
		api.DeletionTimestamp = &metav1.Time{Time: time.Now()}

		Expect(apiv4.ResolveRefs(ctx, api, "default")).To(Succeed())
		Expect(api.Spec.Resources[0].Ref).ToNot(BeNil())
	})
})
