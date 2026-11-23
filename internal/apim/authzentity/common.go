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

// Package authzentity maps the AuthzEntity CRD to the AuthZ entities of the APIM Automation API through the SDK.
package authzentity

import (
	"context"
	"fmt"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-automation-sdk/common/pkg/unstructured"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
)

// ToDTO maps an entity whose refs dynamic.ResolveAuthzEntityRefs has resolved onto the payload a PUT sends: each
// target is a scope, sent in declaration order. It is the single mapping used by the sync path, the dry run and
// drift detection, and reads nothing from the cluster. The platform sorts targets, parents and attribute keys and
// answers {} and [] for absent ones: the SDK's drift tags compare none of them by position. active is always
// sent, so that an unset one reads as true. An unresolved pdpRef is an error: it has no scope to send.
func ToDTO(obj *v1alpha1.AuthzEntity) (sdk.AuthzEntity, error) {
	targets, err := toTargets(obj)
	if err != nil {
		return sdk.AuthzEntity{}, err
	}

	var parents []sdk.AuthzEntityParent
	for _, parent := range obj.Spec.Parents {
		parents = append(parents, sdk.AuthzEntityParent{EntityType: parent.EntityType, EntityId: parent.EntityID})
	}

	return sdk.AuthzEntity{
		Hrid:       refs.NewNamespacedNameFromObject(obj).HRID(),
		EntityId:   obj.Spec.EntityID,
		Kind:       sdk.AuthzEntityKind(obj.Spec.Kind),
		EntityType: obj.Spec.EntityType,
		Attributes: toAttributes(obj.Spec.Attributes),
		Parents:    parents,
		Active:     new(obj.Spec.IsActive()),
		Targets:    targets,
	}, nil
}

func toTargets(obj *v1alpha1.AuthzEntity) ([]string, error) {
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

func toAttributes(attributes *utils.GenericStringMap) *unstructured.Unstructured {
	if attributes == nil {
		return nil
	}
	return unstructured.From(attributes.DeepCopy().Object)
}

// CreateAPIMClient builds the SDK client of the management context the entity references.
func CreateAPIMClient(ctx context.Context, obj *v1alpha1.AuthzEntity) (*apim.SDKClient, error) {
	return apim.SDKFromContextRef(ctx, obj.ContextRef(), obj.GetNamespace())
}
