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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/resource"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

func validateResourceOrRefs(ctx context.Context, api core.ApiDefinitionObject) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	for _, res := range api.GetResources() {
		if res.IsRef() {
			if r, err := dynamic.ResolveResource(ctx, res.GetRef(), api.GetNamespace()); err != nil {
				errs.AddSeveref(
					"api references resource [%s] that does not exist in the cluster",
					res.GetRef(),
				)
			} else {
				res.SetObject(toResourceOrRef(r))
			}
		} else {
			errs.MergeWith(resource.ValidateModel(ctx, res.GetObject()))
		}
	}

	return errs
}

func toResourceOrRef(r core.ResourceModel) core.ResourceModel {
	rn := r.GetResourceName()
	t := r.GetType()
	return &base.Resource{
		Enabled:       true,
		Name:          &rn,
		Type:          &t,
		Configuration: r.GetConfig(),
	}
}
