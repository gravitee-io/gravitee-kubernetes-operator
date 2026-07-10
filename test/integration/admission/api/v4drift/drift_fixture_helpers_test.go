package v4

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	admissionv4 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	. "github.com/onsi/gomega"
)

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

const (
	updatedDescription     = "updated description"
	remoteDescription      = "remote updated description"
	localCRDDescription    = "local CRD description"
	driftDescriptionAssert = `description: "local CRD description" != "remote updated description"`
)

func setDescription(api *v1alpha1.ApiV4Definition, description string) {
	api.Spec.Description = utils.ToReference(description)
}

func validateDescriptionDrift(
	ctx context.Context,
	admissionCtrl admissionv4.AdmissionCtrl,
	oldAPI *v1alpha1.ApiV4Definition,
	newAPI *v1alpha1.ApiV4Definition,
	mgmtContext core.ContextModel,
) {
	apimClient := apim.NewClient(ctx)

	newAPI.PopulateIDs(mgmtContext, true)
	setDescription(newAPI, remoteDescription)

	_, err := apimClient.APIs.ImportV4(newAPI)
	Expect(err).ToNot(HaveOccurred())

	setDescription(newAPI, localCRDDescription)

	Eventually(func() error {
		_, err := admissionCtrl.ValidateUpdate(ctx, oldAPI, newAPI)
		return assert.DriftDetected(driftDescriptionAssert, err)
	}, constants.EventualTimeout, constants.Interval).Should(Succeed())
}
