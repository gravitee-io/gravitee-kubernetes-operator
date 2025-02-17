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

package v4

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func validateSharedPolicyGroups(ctx context.Context, coreApi core.ApiDefinitionObject) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	api, ok := coreApi.(*v1alpha1.ApiV4Definition)
	if !ok {
		errs.AddSevere("unable to convert to ApiV4Definition")
	}

	sharedPolicyGroups := api.Spec.GetAllSharedPolicyGroups()
	for i := range sharedPolicyGroups {
		spg := sharedPolicyGroups[i]
		if spg.Namespace == "" {
			spg.Namespace = api.Namespace
		}

		obj := &v1alpha1.SharedPolicyGroup{}
		key := client.ObjectKey{Namespace: spg.Namespace, Name: spg.Name}
		if err := k8s.GetClient().Get(ctx, key, obj); err != nil {
			errs.AddSeveref("unable to get Shared Policy Group [%s] in namespace [%s]", spg.Name, spg.Namespace)
		}
	}

	return errs
}
