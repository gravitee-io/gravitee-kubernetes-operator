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

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/template"
)

func resolveSharedPolicyGroups(ctx context.Context, spec *v1alpha1.ApiV4DefinitionSpec) error {
	err := resolveFLowSharedPolicyGroupsReferences(ctx, spec.Flows)
	if err != nil {
		return err
	}

	if spec.Plans == nil {
		return nil
	}

	for _, plan := range *spec.Plans {
		err := resolveFLowSharedPolicyGroupsReferences(ctx, plan.Flows)
		if err != nil {
			return err
		}
	}

	return nil
}

//nolint:gocognit // acceptable complexity
func resolveFLowSharedPolicyGroupsReferences(ctx context.Context, flows []*v4.Flow) error {
	if len(flows) == 0 {
		return nil
	}

	for _, flow := range flows {
		for _, flowStep := range flow.Request {
			err := resolveIfSharedPolicyGroupRef(ctx, flowStep)
			if err != nil {
				return err
			}
		}
		for _, flowStep := range flow.Response {
			err := resolveIfSharedPolicyGroupRef(ctx, flowStep)
			if err != nil {
				return err
			}
		}
		for _, flowStep := range flow.Connect {
			err := resolveIfSharedPolicyGroupRef(ctx, flowStep)
			if err != nil {
				return err
			}
		}
		for _, flowStep := range flow.Interact {
			err := resolveIfSharedPolicyGroupRef(ctx, flowStep)
			if err != nil {
				return err
			}
		}
		for _, flowStep := range flow.Publish {
			err := resolveIfSharedPolicyGroupRef(ctx, flowStep)
			if err != nil {
				return err
			}
		}
		for _, flowStep := range flow.Subscribe {
			err := resolveIfSharedPolicyGroupRef(ctx, flowStep)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func resolveIfSharedPolicyGroupRef(ctx context.Context, flowStep *v4.FlowStep) error {
	if flowStep.SharedPolicyGroup == nil {
		return nil
	}

	ns := flowStep.SharedPolicyGroup
	if ns.Namespace == "" {
		key, _ := ctx.Value(nsKey).(string)
		ns.Namespace = key
	}

	spg := new(v1alpha1.SharedPolicyGroup)

	if err := k8s.GetClient().Get(ctx, client.ObjectKey{Namespace: ns.Namespace, Name: ns.Name}, spg); err != nil {
		return gerrors.NewResolveRefError(err)
	}

	if err := template.Compile(ctx, spg, true); err != nil {
		return err
	}

	flowStep.Name = &spg.Name
	flowStep.Policy = utils.ToReference("shared-policy-group-policy")
	flowStep.Configuration = utils.ToGenericStringMap(map[string]interface{}{
		"sharedPolicyGroupId": spg.Status.CrossID,
	})

	return nil
}
