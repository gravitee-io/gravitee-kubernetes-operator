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
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	xErrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kErrors "k8s.io/apimachinery/pkg/util/errors"
)

var errRaw = fmt.Errorf("raw error")
var errNotFound = xErrors.NewControlPlaneError(xErrors.ServerError{StatusCode: 404})
var errBadRequest = xErrors.NewControlPlaneError(xErrors.ServerError{StatusCode: 400})
var errUnauthorized = xErrors.NewControlPlaneError(xErrors.ServerError{StatusCode: 401})

var _ = Describe("Errors", func() {
	DescribeTable("recoverable errors",
		func(given error, expected bool) {
			Expect(xErrors.IsRecoverable(given)).To(Equal(expected))
		},
		Entry("With raw error", errRaw, true),
		Entry("With not found error", errNotFound, true),
		Entry("With unauthorized error", errUnauthorized, false),
		Entry("With bad request", errBadRequest, false),
	)

	DescribeTable("context error",
		func(given error, expected bool) {
			Expect(errors.Is(given, xErrors.NewControlPlaneError(nil))).To(Equal(expected))
		},
		Entry("With raw error", errRaw, false),
		Entry("With nil error", nil, false),
		Entry("With context error", xErrors.NewControlPlaneError(nil), true),
		Entry("With aggregate containing context error",
			kErrors.NewAggregate([]error{xErrors.NewControlPlaneError(nil)}), true),
		Entry("With aggregate not containing any context error",
			kErrors.NewAggregate([]error{errRaw}), false),
		Entry("With empty aggregate", kErrors.NewAggregate([]error{}), false),
	)

	// A ReconcileError wraps the underlying cause; it must stay unwrappable so the
	// standard errors machinery (and apierrors.IsNotFound, on which the delete-path
	// finalizer guards rely) can inspect the wrapped k8s error.
	DescribeTable("apierrors.IsNotFound sees through ReconcileError",
		func(given error, expected bool) {
			Expect(apierrors.IsNotFound(given)).To(Equal(expected))
		},
		Entry("With wrapped k8s NotFound",
			xErrors.NewResolveRefError(
				apierrors.NewNotFound(schema.GroupResource{Group: "gravitee.io", Resource: "portals"}, "p"),
			), true),
		Entry("With wrapped non-NotFound error", xErrors.NewResolveRefError(errRaw), false),
	)
})
