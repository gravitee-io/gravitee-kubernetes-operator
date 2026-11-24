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

package lifecycle_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

var _ = Describe("AdmissionLifecycle resolve error messages", func() {
	ctx := context.Background()

	failingRefs := func(context.Context, *v1alpha1.Group, string) error { return errors.New("missing ref") }
	failingContext := func(context.Context, *v1alpha1.Group) (*testClient, error) {
		return nil, errors.New("missing context")
	}

	DescribeTable("prefixes a resolve error by default, and not with UnprefixedResolveErrors",
		func(unprefixed bool, refsFail bool, expected string) {
			a := minimalAdmission(&testClient{})
			a.UnprefixedResolveErrors = unprefixed
			if refsFail {
				a.ResolveRefs = failingRefs
			} else {
				a.ClientFactory = failingContext
			}

			errs := a.ValidateCreate(ctx, newGroup("g1", "ns"))

			Expect(errs.Severe).To(ConsistOf(MatchError(expected)))
		},
		Entry("references", false, true, "could not resolve references: missing ref"),
		Entry("context", false, false, "could not resolve context: missing context"),
		Entry("references, unprefixed", true, true, "missing ref"),
		Entry("context, unprefixed", true, false, "missing context"),
	)
})
