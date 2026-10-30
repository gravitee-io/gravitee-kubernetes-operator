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
)

func ResolveRefs(ctx context.Context, obj *v1alpha1.AMApplication, ns string) error {
	err := am.ResolveSubResourceDomain(ctx, obj, ns)

	if err != nil || obj.IsBeingDeleted() {
		return err
	}

	if certificateRef := obj.Spec.CertificateRef; certificateRef != nil {
		cert, err := resolveCertificate(ctx, certificateRef, ns)
		if err != nil {
			if apierrors.IsNotFound(err) {
				return fmt.Errorf("AMCertificate [%s] not found: %w", certificateRef.Name, err)
			}
			return err
		}
		if err := validateIsInSameDomain("AMCertificate", cert, obj); err != nil {
			return err
		}
	}

	for _, appIdP := range obj.Spec.IdentityProviders {
		if identityRef := appIdP.IdentityRef; identityRef != nil {
			idp, err := resolveIdentityProvider(ctx, identityRef, ns)
			if err != nil {
				if apierrors.IsNotFound(err) {
					return fmt.Errorf("AMIdentityProvider [%s] not found: %w", identityRef.Name, err)
				}
				return err
			}
			if err := validateIsInSameDomain("AMIdentityProvider", idp, obj); err != nil {
				return err
			}
		}
	}

	return nil

}

func resolveIdentityProvider(ctx context.Context, idpRef *refs.NamespacedName, ns string) (*v1alpha1.AMIdentityProvider, error) {
	idp := &v1alpha1.AMIdentityProvider{}
	err := k8s.GetClient().Get(ctx, namespaceName(*idpRef, ns), idp)
	if err != nil {
		return nil, err
	}
	return idp, nil
}

func resolveCertificate(ctx context.Context, certificateRef *refs.NamespacedName, ns string) (*v1alpha1.AMCertificate, error) {
	cert := &v1alpha1.AMCertificate{}
	err := k8s.GetClient().Get(ctx, namespaceName(*certificateRef, ns), cert)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func namespaceName(ref refs.NamespacedName, ns string) types.NamespacedName {
	if ref.Namespace == "" {
		ref.Namespace = ns
	}
	return ref.NamespacedName()
}

func validateIsInSameDomain(refKind string, ref core.AMDomainSubResource, recipient core.AMDomainSubResource) error {
	if refDomain, recipientDomain := am.DomainRef(ref), am.DomainRef(recipient); *refDomain != *recipientDomain {
		return fmt.Errorf("%s [%s] belongs to domain [%s], it should belong to [%s] or AM will not find it", refKind, ref.GetRef().GetName(), refDomain.Name, recipientDomain.Name)
	}
	return nil
}
