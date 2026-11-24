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

package group

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/ctxref"
	apimgroup "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

func validateCreate(ctx context.Context, group *v1alpha1.Group) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	errs.Add(ctxref.Validate(ctx, group))
	if errs.IsSevere() {
		return errs
	}
	if k8s.HasTrueAnnotation(group, core.IgnoreMembersAnnotation) && len(group.Spec.Members) > 0 {
		errs.AddWarningf("spec.members is not applied while %s is true", core.IgnoreMembersAnnotation)
	}
	errs.MergeWith(validateDryRun(ctx, group))
	return errs
}

func validateDryRun(ctx context.Context, group *v1alpha1.Group) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	client, err := apimgroup.CreateAPIMClient(ctx, group)
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	dto, err := apimgroup.ToDTO(group)
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	findings, err := apimgroup.DryRun(ctx, client, dto)
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	errs.MergeWith(findings)
	return errs
}
