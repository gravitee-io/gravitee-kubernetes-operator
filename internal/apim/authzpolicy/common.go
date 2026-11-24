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

// Package authzpolicy maps the AuthzPolicy CRD to the AuthZ policies of the APIM Automation API through the SDK.
package authzpolicy

import (
	"context"
	"fmt"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
)

// ToDTO maps a policy whose refs dynamic.ResolveAuthzPolicyRefs has resolved onto the payload a PUT sends: each
// target is a scope, sent in declaration order. It is the single mapping used by the sync path, the dry run and
// drift detection, and reads nothing from the cluster. The platform stores and answers policyText verbatim, so it
// compares exactly; it sorts targets and answers [] for none, which the SDK's drift tag does not compare by
// position. An unresolved pdpRef is an error: it has no scope to send.
func ToDTO(obj *v1alpha1.AuthzPolicy) (sdk.AuthzPolicy, error) {
	targets, err := toTargets(obj)
	if err != nil {
		return sdk.AuthzPolicy{}, err
	}

	return sdk.AuthzPolicy{
		Hrid:       refs.NewNamespacedNameFromObject(obj).HRID(),
		EntityId:   obj.Spec.EntityID,
		Name:       obj.Spec.Name,
		PolicyType: sdk.AuthzPolicyType(obj.Spec.PolicyType),
		PolicyText: obj.Spec.PolicyText,
		Targets:    targets,
	}, nil
}

func toTargets(obj *v1alpha1.AuthzPolicy) ([]string, error) {
	var targets []string
	for i, target := range obj.Spec.Targets {
		if target.PdpRef != nil {
			ref := *target.PdpRef
			if ref.Namespace == "" {
				ref.Namespace = obj.GetNamespace()
			}
			return nil, fmt.Errorf("targets[%d].pdpRef [%s] is not resolved", i, ref.String())
		}
		targets = append(targets, target.Scope)
	}
	return targets, nil
}

// CreateAPIMClient builds the SDK client of the management context the policy references.
func CreateAPIMClient(ctx context.Context, obj *v1alpha1.AuthzPolicy) (*apim.SDKClient, error) {
	return apim.SDKFromContextRef(ctx, obj.ContextRef(), obj.GetNamespace())
}
