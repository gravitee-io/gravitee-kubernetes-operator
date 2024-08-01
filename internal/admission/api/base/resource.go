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

package base

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/resource"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

func validateResourceOrRefs(ctx context.Context, api core.ApiDefinitionObject) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	for _, res := range api.GetResources() {
		if res.IsRef() && dynamic.ExpectResolvedResource(ctx, res.GetRef(), api.GetNamespace()) != nil {
			errs.AddSevere("api references resource [%s] that does not exist in the cluster", res.GetRef())
		} else {
			errs.MergeWith(resource.ValidateModel(ctx, res.GetObject()))
		}
	}
	return errs
}
