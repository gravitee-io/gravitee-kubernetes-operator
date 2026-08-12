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

package dynamic

import (
	"context"
	"testing"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDynamic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "package dynamic")
}

// References are not always required, e.g. contextRef is optional on API definitions.
// An unset reference reaches the resolvers as an interface wrapping a nil pointer, which
// must be reported as an error rather than dereferenced. These resolve before any client
// call, so they need no cluster.
var _ = Describe("Resolving a nil reference", func() {
	ctx := context.Background()

	DescribeTable("returns an error instead of panicking",
		func(resolve func(core.ObjectRef) error) {
			// An unset reference reaches a resolver as an interface wrapping a nil
			// pointer, which is not equal to nil on its own.
			var unset *refs.NamespacedName
			Expect(resolve(unset)).To(HaveOccurred())

			// A caller may also pass no reference at all.
			Expect(resolve(nil)).To(HaveOccurred())
		},
		Entry("management context", func(ref core.ObjectRef) error {
			_, err := dynamic.ResolveContext(ctx, ref, "default")
			return err
		}),
		Entry("API definition", func(ref core.ObjectRef) error {
			_, err := dynamic.ResolveAPI(ctx, ref, "default")
			return err
		}),
		Entry("application", func(ref core.ObjectRef) error {
			_, err := dynamic.ResolveApplication(ctx, ref, "default")
			return err
		}),
		Entry("group", func(ref core.ObjectRef) error {
			_, err := dynamic.ResolveGroup(ctx, ref, "default")
			return err
		}),
		Entry("API resource", func(ref core.ObjectRef) error {
			_, err := dynamic.ResolveResource(ctx, ref, "default")
			return err
		}),
		Entry("portal", func(ref core.ObjectRef) error {
			_, err := dynamic.ResolvePortal(ctx, ref, "default")
			return err
		}),
		Entry("notification", func(ref core.ObjectRef) error {
			_, err := dynamic.ResolveNotification(ctx, ref, "default")
			return err
		}),
		Entry("secret", func(ref core.ObjectRef) error {
			_, err := dynamic.ResolveSecret(ctx, ref, "default")
			return err
		}),
	)
})
