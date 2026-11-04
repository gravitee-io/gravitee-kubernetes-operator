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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
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
}

// AuthzPdpDependents lists the kinds that target an AuthzPdp, which deleting one warns about.
// Each authz story adds its kind: AuthzSchema (AUTHZ-75), AuthzEntity (AUTHZ-81) and
// AuthzPolicy (AUTHZ-84).
var AuthzPdpDependents = []AuthzPdpDependent{}

// Dependents are the resources of one kind that target an AuthzPdp.
type Dependents struct {
	Kind string
	// Names are namespace/name, sorted.
	Names []string
}

// FindAuthzPdpDependents returns, for each kind in order, the resources whose pdpRef names the
// AuthzPdp ref. A kind with none is left out.
func FindAuthzPdpDependents(
	ctx context.Context,
	ref refs.NamespacedName,
	kinds []AuthzPdpDependent,
) ([]Dependents, error) {
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
		for _, item := range items {
			if obj, ok := item.(client.Object); ok {
				names = append(names, refs.NewNamespacedNameFromObject(obj).String())
			}
		}
		if len(names) == 0 {
			continue
		}
		sort.Strings(names)
		found = append(found, Dependents{Kind: kind.Kind, Names: names})
	}
	return found, nil
}
