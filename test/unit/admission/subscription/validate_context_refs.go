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

package subscription

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/management"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/subscription"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func contextRef(namespace, name string) *refs.NamespacedName {
	return &refs.NamespacedName{Namespace: namespace, Name: name}
}

func managementContext(namespace, name string, spec management.Context) *v1alpha1.ManagementContext {
	return &v1alpha1.ManagementContext{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: name},
		Spec:       v1alpha1.ManagementContextSpec{Context: &spec},
	}
}

func localAPIM(mutate func(*management.Context)) management.Context {
	spec := management.Context{BaseUrl: "http://localhost:30083", OrgID: "DEFAULT", EnvID: "DEFAULT"}
	if mutate != nil {
		mutate(&spec)
	}
	return spec
}

// A contextRef is optional on API definitions, so it reaches validation as an interface
// wrapping a nil pointer. Every combination must yield a verdict rather than a panic.
var _ = Describe("Validating context refs", func() {
	DescribeTable("subscribing an application to an API without context",
		func(apiContext, appContext *refs.NamespacedName) {
			api := &v1alpha1.ApiV4Definition{Spec: v1alpha1.ApiV4DefinitionSpec{Context: apiContext}}
			app := &v1alpha1.Application{Spec: v1alpha1.ApplicationSpec{Context: appContext}}

			Expect(subscription.ValidateContextRefs(context.Background(), api, app)).ToNot(BeNil())
		},
		Entry("is rejected when the API has no context",
			nil, contextRef("default", "dev-ctx")),
		Entry("is rejected when the application has no context",
			contextRef("default", "dev-ctx"), nil),
		Entry("is rejected when neither has a context",
			nil, nil),
	)
})

var _ = Describe("Comparing management context targets", func() {
	path := "/proxied/automation"

	DescribeTable("two contexts",
		func(left, right *v1alpha1.ManagementContext, expected bool) {
			same, err := subscription.SameTarget(left, right)
			Expect(err).ToNot(HaveOccurred())
			Expect(same).To(Equal(expected))
		},
		Entry("match when they are the same object",
			managementContext("default", "dev-ctx", localAPIM(nil)),
			managementContext("default", "dev-ctx", localAPIM(nil)),
			true),
		Entry("match when distinct contexts target the same environment",
			managementContext("default", "team-a", localAPIM(nil)),
			managementContext("default", "team-b", localAPIM(nil)),
			true),
		Entry("match when distinct contexts in distinct namespaces target the same environment",
			managementContext("team-a", "ctx", localAPIM(nil)),
			managementContext("team-b", "ctx", localAPIM(nil)),
			true),
		Entry("match when only a trailing slash on the base URL differs",
			managementContext("default", "team-a", localAPIM(nil)),
			managementContext("default", "team-b", localAPIM(func(c *management.Context) {
				c.BaseUrl = "http://localhost:30083/"
			})),
			true),
		Entry("differ when the environment differs",
			managementContext("default", "dev", localAPIM(nil)),
			managementContext("default", "prod", localAPIM(func(c *management.Context) { c.EnvID = "PROD" })),
			false),
		Entry("differ when the organization differs",
			managementContext("default", "dev", localAPIM(nil)),
			managementContext("default", "other", localAPIM(func(c *management.Context) { c.OrgID = "OTHER" })),
			false),
		Entry("differ when the base URL differs",
			managementContext("default", "dev", localAPIM(nil)),
			managementContext("default", "other", localAPIM(func(c *management.Context) {
				c.BaseUrl = "http://apim.example.com"
			})),
			false),
		Entry("differ when the same name lives in two namespaces with distinct targets",
			managementContext("team-a", "ctx", localAPIM(nil)),
			managementContext("team-b", "ctx", localAPIM(func(c *management.Context) { c.EnvID = "PROD" })),
			false),
		Entry("differ when one context overrides the path",
			managementContext("default", "direct", localAPIM(nil)),
			managementContext("default", "proxied", localAPIM(func(c *management.Context) { c.Path = &path })),
			false),
		Entry("differ when one context targets Gravitee Cloud",
			managementContext("default", "self-hosted", localAPIM(nil)),
			managementContext("default", "cloud", localAPIM(func(c *management.Context) {
				c.Cloud = &management.Cloud{Token: "token"}
			})),
			false),
	)

	It("fails on a base URL that does not parse", func() {
		_, err := subscription.SameTarget(
			managementContext("default", "dev", localAPIM(nil)),
			managementContext("default", "broken", localAPIM(func(c *management.Context) { c.BaseUrl = "http://[::1" })),
		)
		Expect(err).To(HaveOccurred())
	})
})
