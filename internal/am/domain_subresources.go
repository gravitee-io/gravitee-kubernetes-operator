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
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
// sub-resource cannot be created under it, so ResolvedRefs=False and the reconcile is retried.
// A missing domain on delete releases the finalizer: the domain took the sub-resource with it.
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

// ResolveSubResourceDomain is the lifecycle ResolveRefs hook of a domain sub-resource, see ResolveDomain.
func ResolveSubResourceDomain[T core.AMDomainSubResource](ctx context.Context, obj T, _ string) error {
	return ResolveDomain(ctx, obj, obj.IsBeingDeleted())
}

// GetSubResourceOwner returns the AMSecurityDomain that owns the sub-resource (a plain owner, not blocking deletion).
func GetSubResourceOwner[T core.AMDomainSubResource](ctx context.Context, obj T) (client.Object, bool, error) {
	d, e := GetDomain(ctx, obj)
	return d, false, e
}

// CreateSubResourceClient builds the AM client from the AMContext of the parent domain.
// A missing domain returns the NotFound as-is: on delete, the lifecycle then releases the finalizer.
// The domain is read again here, as in Owner and ResolveRefs: the hooks cannot share it, and the
// reads hit the manager's informer cache, not the API server.
func CreateSubResourceClient[T core.AMDomainSubResource](ctx context.Context, obj T) (*Client, error) {
	domain, err := GetDomain(ctx, obj)
	if err != nil {
		return nil, err
	}
	return ClientForDomain(ctx, domain)
}

// AdmissionSubResourceClient builds the AM client for admission. A domain that is missing or not yet created
// in AM gives no client and no error: the sub-resource is admitted (apply in any order) and the AM calls
// are skipped. A missing AMContext still fails.
func AdmissionSubResourceClient[T core.AMDomainSubResource](ctx context.Context, obj T) (*Client, error) {
	err := ResolveDomain(ctx, obj, obj.IsBeingDeleted())
	if apierrors.IsNotFound(err) || errors.Is(err, ErrDomainNotReady) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return CreateSubResourceClient(ctx, obj)
}

// ClientForDomain builds the AM client from the domain's AMContext, resolved with its templates compiled.
func ClientForDomain(ctx context.Context, domain *v1alpha1.AMSecurityDomain) (*Client, error) {
	if !domain.HasContext() {
		return nil, fmt.Errorf("contextRef empty on AMSecurityDomain [%s/%s]", domain.GetName(), domain.GetNamespace())
	}

	// resolved like the APIM contexts: templates compiled, fetched from the API server
	resolved, err := dynamic.ResolveAMContext(ctx, domain.ContextRef(), domain.GetNamespace())
	if err != nil {
		return nil, fmt.Errorf("AMContext [%s]: %w", domain.ContextRef().String(), err)
	}
	amContext, ok := resolved.(*v1alpha1.AMContext)
	if !ok {
		return nil, fmt.Errorf("AMContext [%s]: unexpected type %T", domain.ContextRef().String(), resolved)
	}

	return NewSDKClient(ctx, amContext)
}

// WarnSystemIgnoredFields warns when a system sub-resource sets name, type or configuration: AM ignores them.
func WarnSystemIgnoredFields(
	system *bool, name, typ *string, configuration *utils.GenericStringMap,
) *gerrors.AdmissionErrors {
	errs := gerrors.NewAdmissionErrors()
	if system == nil || !*system {
		return errs
	}
	fields := make([]string, 0)
	if configuration != nil && len(configuration.Object) > 0 {
		fields = append(fields, "configuration")
	}
	if name != nil && len(*name) > 0 {
		fields = append(fields, "name")
	}
	if typ != nil && len(*typ) > 0 {
		fields = append(fields, "type")
	}
	if len(fields) > 0 {
		errs.AddWarningf("'%s' will be ignored when 'system' is 'true'.", strings.Join(fields, "', '"))
	}
	return errs
}
