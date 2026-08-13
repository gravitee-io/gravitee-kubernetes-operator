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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/subscription"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func contextRef(namespace, name string) *refs.NamespacedName {
	return &refs.NamespacedName{Namespace: namespace, Name: name}
}

// A contextRef is optional on API definitions, so it reaches validation as an interface
// wrapping a nil pointer. Every combination must yield a verdict rather than a panic.
var _ = Describe("Validating context refs", func() {
	DescribeTable("subscribing an application to an API",
		func(apiContext, appContext *refs.NamespacedName, expectRejection bool) {
			api := &v1alpha1.ApiV4Definition{Spec: v1alpha1.ApiV4DefinitionSpec{Context: apiContext}}
			app := &v1alpha1.Application{Spec: v1alpha1.ApplicationSpec{Context: appContext}}

			err := subscription.ValidateContextRefs(api, app)

			if expectRejection {
				Expect(err).ToNot(BeNil())
				return
			}
			Expect(err).To(BeNil())
		},
		Entry("is accepted when both reference the same context",
			contextRef("default", "dev-ctx"), contextRef("default", "dev-ctx"), false),
		Entry("is rejected when the context names differ",
			contextRef("default", "prod-ctx"), contextRef("default", "dev-ctx"), true),
		Entry("is rejected when the context namespaces differ",
			contextRef("other", "dev-ctx"), contextRef("default", "dev-ctx"), true),
		Entry("is rejected when the API has no context",
			nil, contextRef("default", "dev-ctx"), true),
		Entry("is rejected when the application has no context",
			contextRef("default", "dev-ctx"), nil, true),
		Entry("is rejected when neither has a context",
			nil, nil, true),
	)
})
