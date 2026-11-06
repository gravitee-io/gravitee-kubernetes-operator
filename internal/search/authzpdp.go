// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package search

import (
	"context"
	"fmt"
	"sort"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"k8s.io/apimachinery/pkg/api/meta"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AuthzPdpDependent is a kind whose resources target an AuthzPdp through a pdpRef, found through
// the index its controller registers on the AuthzPdp's namespaced name.
type AuthzPdpDependent struct {
	// Kind as a message names it, e.g. AuthzSchema.
	Kind string
	// Field is the index on the pdpRef namespaced name, its namespace defaulted.
	Field IndexField
	// NewList returns an empty list of the kind.
	NewList func() client.ObjectList
	// Targets returns the targets of a resource of the kind. Set for a kind the platform does not
	// leave without a target: a PDP delete that would is refused with a 409.
	Targets func(client.Object) []authzpdp.Target
}

// AuthzPdpDependents returns the kinds that target an AuthzPdp, which deleting one warns about.
// Each authz story adds its kind to the list: AuthzSchema (AUTHZ-75), then AuthzEntity (AUTHZ-81)
// and AuthzPolicy (AUTHZ-84).
func AuthzPdpDependents() []AuthzPdpDependent {
	return []AuthzPdpDependent{
		{
			Kind:    "AuthzSchema",
			Field:   AuthzSchemaAuthzPdpField,
			NewList: func() client.ObjectList { return &v1alpha1.AuthzSchemaList{} },
		},
		{
			Kind:    "AuthzEntity",
			Field:   AuthzEntityAuthzPdpField,
			NewList: func() client.ObjectList { return &v1alpha1.AuthzEntityList{} },
			Targets: func(obj client.Object) []authzpdp.Target { return obj.(*v1alpha1.AuthzEntity).Spec.Targets },
		},
		{
			Kind:        "AuthzPolicy",
			Field:       AuthzPolicyAuthzPdpField,
			NewList:     func() client.ObjectList { return &v1alpha1.AuthzPolicyList{} },
			NeedsTarget: true,
		},
	}
}

// Dependents are the resources of one kind that target an AuthzPdp.
type Dependents struct {
	Kind string
	// Names are namespace/name, sorted.
	Names []string
	// Sole are the Names whose every target is the AuthzPdp, for a kind with Targets: the platform
	// refuses the AuthzPdp's delete while one of them does. Sorted, nil when none.
	Sole []string
}

// FindAuthzPdpDependents returns, for each kind in order, the resources whose pdpRef names pdp and
// that sync through pdp's ManagementContext: a pdpRef to a PDP of another environment does not
// resolve (see dynamic.ResolveAuthzPdpRoutingScope), so that resource does not target it. A kind
// with none is left out. A kind whose resources carry no contextRef is an error: its registration
// is a mistake, and dropping them would hide it.
func FindAuthzPdpDependents(
	ctx context.Context,
	pdp *v1alpha1.AuthzPdp,
	kinds []AuthzPdpDependent,
) ([]Dependents, error) {
	ref := *refs.NewNamespacedNameFromObject(pdp)
	pdpContext := dynamic.ContextOf(pdp)

	found := make([]Dependents, 0)
	for _, kind := range kinds {
		list := kind.NewList()
		if err := FindByFieldReferencing(ctx, kind.Field, ref, list); err != nil {
			return nil, fmt.Errorf("unable to list %s: %w", kind.Kind, err)
		}
		items, err := meta.ExtractList(list)
		if err != nil {
			return nil, fmt.Errorf("unable to read %s: %w", kind.Kind, err)
		}
		names := make([]string, 0, len(items))
		var sole []string
		for _, item := range items {
			obj, ok := item.(core.ContextAwareObject)
			if !ok {
				return nil, fmt.Errorf("unable to read the contextRef of %s: %T is not context aware", kind.Kind, item)
			}
			if dynamic.ContextOf(obj) != pdpContext {
				continue
			}
			name := refs.NewNamespacedNameFromObject(obj).String()
			names = append(names, name)
			if kind.Targets != nil && targetsOnly(obj, kind.Targets(obj), ref.String()) {
				sole = append(sole, name)
			}
		}
		if len(names) == 0 {
			continue
		}
		sort.Strings(names)
		sort.Strings(sole)
		found = append(found, Dependents{Kind: kind.Kind, Names: names, Sole: sole})
	}
	return found, nil
}

// targetsOnly tells whether every target of obj is a pdpRef to pdp, a namespace/name: obj then has
// no other gateway once pdp is deleted.
func targetsOnly(obj client.Object, targets []authzpdp.Target, pdp string) bool {
	for _, target := range targets {
		if target.PdpRef == nil || ensureNamespacedRef(obj, target.PdpRef) != pdp {
			return false
		}
	}
	return len(targets) > 0
}
