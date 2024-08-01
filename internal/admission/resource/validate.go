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

package resource

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func validateCreate(ctx context.Context, obj runtime.Object) *errors.AdmissionErrors {
	if res, ok := obj.(core.ResourceModel); ok {
		return ValidateModel(ctx, res)
	}
	return errors.NewAdmissionErrors()
}

func ValidateModel(ctx context.Context, res core.ResourceModel) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if res.GetType() == "" {
		errs.AddSevere("missing required value in property [type]")
	}

	if res.GetResourceName() == "" {
		errs.AddSevere("missing required value in property [name]")
	}

	if res.GetConfig() == nil {
		errs.AddSevere("missing required value in property [configuration]")
	}

	return errs
}
