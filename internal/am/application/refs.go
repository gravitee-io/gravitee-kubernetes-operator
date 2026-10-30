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

package application

import (
	"context"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// siblingObject is a CR the application references: an AMCertificate or an AMIdentityProvider.
type siblingObject interface {
	client.Object
	core.AMDomainSubResource
}

// siblingRef is a certificateRef or an identityRef of the application, with the CR it is read into.
type siblingRef struct {
	kind   string
	field  string
	ref    *refs.NamespacedName
	target siblingObject
}

// siblingRefs lists the certificateRef and identityRef of the application.
func siblingRefs(obj *v1alpha1.AMApplication) []siblingRef {
	var siblings []siblingRef
	if ref := obj.Spec.CertificateRef; ref != nil {
		siblings = append(siblings, siblingRef{
			kind: "AMCertificate", field: "certificateRef", ref: ref, target: &v1alpha1.AMCertificate{},
		})
	}
	for i, idp := range obj.Spec.IdentityProviders {
		if ref := idp.IdentityRef; ref != nil {
			siblings = append(siblings, siblingRef{
				kind: "AMIdentityProvider", field: fmt.Sprintf("identityRef[%d]", i), ref: ref,
				target: &v1alpha1.AMIdentityProvider{},
			})
		}
	}
	return siblings
}

// get reads the referenced CR in the application's namespace, as its AM key and the indexers do
// (a cross-namespace ref is rejected by admission).
func (s siblingRef) get(ctx context.Context, ns string) error {
	return k8s.GetClient().Get(ctx, types.NamespacedName{Namespace: ns, Name: s.ref.Name}, s.target)
}

func ResolveRefs(ctx context.Context, obj *v1alpha1.AMApplication, ns string) error {
	err := am.ResolveSubResourceDomain(ctx, obj, ns)

	if err != nil || obj.IsBeingDeleted() {
		return err
	}

	for _, sibling := range siblingRefs(obj) {
		err := sibling.get(ctx, ns)
		if apierrors.IsNotFound(err) {
			return fmt.Errorf("%s [%s] not found: %w", sibling.kind, sibling.ref.Name, err)
		}
		if err != nil {
			return fmt.Errorf("%s [%s] could not be read: %w", sibling.kind, sibling.ref.Name, err)
		}
		if err := validateIsInSameDomain(sibling.kind, sibling.target, obj); err != nil {
			return err
		}
	}

	return nil
}

func validateIsInSameDomain(refKind string, ref core.AMDomainSubResource, recipient core.AMDomainSubResource) error {
	if refDomain, recipientDomain := am.DomainRef(ref), am.DomainRef(recipient); *refDomain != *recipientDomain {
		return fmt.Errorf("%s [%s] belongs to domain [%s], it should belong to [%s] or AM will not find it", refKind, ref.GetRef().GetName(), refDomain.Name, recipientDomain.Name)
	}
	return nil
}
