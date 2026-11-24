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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	xErrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ = Describe("NotFoundRef", func() {
	notFound := apierrors.NewNotFound(schema.GroupResource{Group: "gravitee.io", Resource: "portals"}, "main")

	It("words a NotFound and keeps it a NotFound", func() {
		err := xErrors.NotFoundRef(notFound, "portal [%s] can't be resolved", "main")

		Expect(err).To(MatchError("portal [main] can't be resolved"))
		Expect(apierrors.IsNotFound(err)).To(BeTrue())
		Expect(apierrors.IsNotFound(xErrors.NewResolveRefError(err))).To(BeTrue())
	})

	It("returns any other error as is", func() {
		Expect(xErrors.NotFoundRef(errRaw, "portal [%s] can't be resolved", "main")).To(BeIdenticalTo(errRaw))
	})
})
