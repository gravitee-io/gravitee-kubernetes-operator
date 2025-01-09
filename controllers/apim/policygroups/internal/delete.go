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

package internal

import (
	"context"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func Delete(
	ctx context.Context,
	spg *v1alpha1.SharedPolicyGroup,
) error {
	if !util.ContainsFinalizer(spg, core.SharedPolicyGroupFinalizer) {
		return nil
	}

	if err := search.AssertNoSharedPolicyGroupRef(ctx, spg); err != nil {
		return err
	}

	apim, apimErr := apim.FromContextRef(ctx, spg.Spec.Context, spg.GetNamespace())
	if apimErr != nil {
		return apimErr
	}

	if spg.Status.ID == "" {
		return fmt.Errorf("can not delete a CRD that hasn't been successfuly created in APIM")
	}

	if err := apim.SharedPolicyGroup.Delete(spg.Status.ID); errors.IgnoreNotFound(err) != nil {
		return err
	}

	return nil
}
