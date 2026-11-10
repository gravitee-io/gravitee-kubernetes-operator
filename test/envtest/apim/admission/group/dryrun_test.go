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

package apimgroupadmission

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
)

var _ = Describe("Group dry run", func() {
	ctx := context.Background()

	It("rejects a group whose APIM dry run reports a severe error, and persists nothing", func() {
		objects := groupFixtures(constants.GroupFile)
		grp := objects.Group
		objects.Group = nil
		objects.Apply()
		grp.Spec.Context = objects.Context.GetNamespacedName()

		_, err := admission.AdmissionCtrl{}.ValidateCreate(ctx, grp)
		Expect(err).To(HaveOccurred())

		status, _ := remoteGroup(ctx, grp)
		Expect(status).To(Equal(http.StatusNotFound))
	})
})
