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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/policygroups"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"

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

	if api.Spec.Flows != nil {
		errs.Add(validateFlowSharedPolicyGroups(ctx, api.Spec.Flows, api.Namespace))
	}

	if api.Spec.Plans != nil {
		for _, plan := range *api.Spec.Plans {
			if plan.Flows != nil {
				errs.Add(validateFlowSharedPolicyGroups(ctx, plan.Flows, api.Namespace))
			}
		}
	}

	return errs
}

//nolint:gocognit // acceptable complexity
func validateFlowSharedPolicyGroups(ctx context.Context, flows []*v4.Flow,
	parentNS string) *errors.AdmissionError {
	for _, flow := range flows {
		for _, flowStep := range flow.Request {
			return validateSharedPolicyGroup(ctx, flowStep, parentNS, "REQUEST")
		}

		for _, flowStep := range flow.Response {
			if flowStep.SharedPolicyGroup != nil {
				return validateSharedPolicyGroup(ctx, flowStep, parentNS, "RESPONSE")
			}
		}

		for _, flowStep := range flow.Connect {
			if flowStep.SharedPolicyGroup != nil {
				return validateSharedPolicyGroup(ctx, flowStep, parentNS, "CONNECT")
			}
		}

		for _, flowStep := range flow.Interact {
			if flowStep.SharedPolicyGroup != nil {
				return validateSharedPolicyGroup(ctx, flowStep, parentNS, "INTERACT")
			}
		}

		for _, flowStep := range flow.Publish {
			if flowStep.SharedPolicyGroup != nil {
				return validateSharedPolicyGroup(ctx, flowStep, parentNS, "PUBLISH")
			}
		}

		for _, flowStep := range flow.Subscribe {
			if flowStep.SharedPolicyGroup != nil {
				return validateSharedPolicyGroup(ctx, flowStep, parentNS, "SUBSCRIBE")
			}
		}
	}

	return nil
}

func validateSharedPolicyGroup(ctx context.Context, flowStep *v4.FlowStep, parentNS string,
	phase policygroups.FlowPhase) *errors.AdmissionError {
	if flowStep.SharedPolicyGroup != nil {
		spg, err := getSharePolicyGroup(ctx, flowStep.SharedPolicyGroup, parentNS)
		if err != nil {
			return err
		}

		if *spg.Spec.Phase != phase {
			return errors.NewSeveref("Incompatible Shared Policy Group [%s] with phase [%s] for FlowStep [%s]",
				spg.Name, *spg.Spec.Phase, phase)
		}
	}

	return nil
}

func getSharePolicyGroup(ctx context.Context, spg *refs.NamespacedName, parentNs string) (*v1alpha1.SharedPolicyGroup,
	*errors.AdmissionError) {
	if spg.Namespace == "" {
		spg.Namespace = parentNs
	}

	obj := &v1alpha1.SharedPolicyGroup{}
	key := client.ObjectKey{Namespace: spg.Namespace, Name: spg.Name}
	if err := k8s.GetClient().Get(ctx, key, obj); err != nil {
		return nil, errors.NewSeveref("unable to get Shared Policy Group [%s] in namespace [%s]", spg.Name, spg.Namespace)
	}

	return obj, nil
}
