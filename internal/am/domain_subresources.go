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

package am

import (
	"context"
	"errors"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

var ErrDomainNotReady = errors.New("not yet created in AM")

// ValidateDomainRefNamespace rejects a resource whose domainRef.namespace does not match its own namespace.
func ValidateDomainRefNamespace(obj core.AMDomainSubResource) *gerrors.AdmissionErrors {
	errs := gerrors.NewAdmissionErrors()
	domainRef := obj.GetDomainRef()
	ref := obj.GetRef()
	if domainRef.GetNamespace() != "" && domainRef.GetNamespace() != ref.GetNamespace() {
		errs.AddSeveref("domainRef.namespace [%s] must match the namespace of [%s], [%s]: "+
			"cross-namespace domainRef is not supported",
			domainRef.GetNamespace(),
			ref.GetName(),
			ref.GetNamespace())
	}
	return errs
}

// WarnMissingDomain warns if the domain is not found or not yet created in AM.
func WarnMissingDomain(ctx context.Context, obj core.AMDomainSubResource, beingDeleted bool) *gerrors.AdmissionErrors {
	errs := gerrors.NewAdmissionErrors()
	err := ResolveDomain(ctx, obj, beingDeleted)
	if apierrors.IsNotFound(err) || errors.Is(err, ErrDomainNotReady) {
		errs.AddWarningf("domain [%s] not found or not yet created in AM; [%s] will be created once the domain is",
			DomainRef(obj).String(),
			obj.GetRef().GetName())
	}
	return errs
}

// ResolveDomain checks that the parent domain exists and, on apply, that AM has created it: until then the
// identity provider cannot be created under it, so ResolvedRefs=False and the reconcile is retried.
// A missing domain on delete releases the finalizer: the domain took the identity provider with it.
func ResolveDomain(ctx context.Context, obj core.AMDomainSubResource, beingDeleted bool) error {
	domain, err := GetDomain(ctx, obj)
	if err != nil {
		return err
	}

	if !beingDeleted && domain.Status.Key == "" {
		return fmt.Errorf("AMSecurityDomain [%s]: %w", domain.GetName(), ErrDomainNotReady)
	}
	return nil
}
