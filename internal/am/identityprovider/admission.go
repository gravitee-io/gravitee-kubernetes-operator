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

package identityprovider

import (
	"context"
	"errors"
	"net/http"
	"strings"

	amsdk "github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// PreCheck validates the identity provider
// key pattern and length;
// domainRef.namespace must match the resource's own namespace;
// If the domain is not found or ready, a warning is issued.
// if 'system' is true: setting any of name, type, and configuration them gives a warning.
func PreCheck(ctx context.Context, obj *v1alpha1.AMIdentityProvider) *gerrors.AdmissionErrors {
	errs := gerrors.NewAdmissionErrors()

	key := refs.NewNamespacedNameFromObject(obj).HRID()
	errs.MergeWith(am.ValidateKey(key))
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(am.ValidateDomainRefNamespace(obj))
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(am.WarnMissingDomain(ctx, obj, obj.IsBeingDeleted()))

	if obj.Spec.System != nil && *obj.Spec.System {
		fields := make([]string, 0)
		if obj.Spec.Configuration != nil && len(obj.Spec.Configuration.Object) > 0 {
			fields = append(fields, "configuration")
		}
		if obj.Spec.Name != nil && len(*obj.Spec.Name) > 0 {
			fields = append(fields, "name")
		}
		if obj.Spec.Type != nil && len(*obj.Spec.Type) > 0 {
			fields = append(fields, "type")
		}
		if len(fields) > 0 {
			errs.AddWarningf("'%s' will be ignored when 'system' is 'true'.", strings.Join(fields, "', '"))
		}
	}

	return errs
}

// AdmissionClient builds the AM client for admission. A domain that is missing or not yet created in AM
// gives no client and no error: the identity provider is admitted (apply in any order) and the AM calls
// are skipped. A missing AMContext still fails.
func AdmissionClient(ctx context.Context, obj *v1alpha1.AMIdentityProvider) (*am.Client, error) {
	err := ResolveDomain(ctx, obj, obj.GetNamespace())
	if apierrors.IsNotFound(err) || errors.Is(err, am.ErrDomainNotReady) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return CreateAMClient(ctx, obj)
}

// DryRun validates the identity provider against AM without persisting it. It is skipped without a client
// (domain not in AM yet). Any AM error, an unreachable AM included, rejects it, as for the domain and the
// APIM resources.
func DryRun(ctx context.Context, client *am.Client, dto IdentityProvider) *gerrors.AdmissionErrors {
	errs := gerrors.NewAdmissionErrors()
	if client == nil {
		return errs
	}

	resp, err := client.UpsertIdentityProviderWithResponse(ctx, dto.DomainKey,
		&amsdk.UpsertIdentityProviderParams{DryRun: new(true)}, dto.IdentityProvider)

	if err = am.HasErrors(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		errs.AddSevere(err.Error())
		return errs
	}
	if resp.JSON200 == nil {
		errs.AddSevere(am.UnexpectedResponse(resp.HTTPResponse).Error())
		return errs
	}
	return am.ToAdmissionErrors(resp.JSON200.DryRunErrors)
}

// GetRemote fetches the identity provider for drift detection. Without a client (domain not in AM
// yet) nothing can have drifted: the resource is its own remote, so an update is admitted, as a create is.
func GetRemote(ctx context.Context, client *am.Client, dto IdentityProvider) (IdentityProvider, error) {
	if client == nil {
		return dto, nil
	}
	resp, err := client.GetIdentityProviderWithResponse(ctx, dto.DomainKey, dto.Key)
	if err = am.HasErrors(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return IdentityProvider{}, err
	}
	if resp.JSON200 == nil {
		return IdentityProvider{}, am.UnexpectedResponse(resp.HTTPResponse)
	}
	return IdentityProvider{IdentityProvider: *resp.JSON200, DomainKey: dto.DomainKey}, nil
}

// ToIdentityProviderDTOForDrift wraps ToIdentityProviderDTO and blanks "name", "configuration" and "type" for the
// system identity provider: AM ignores them and returns its own.
func ToIdentityProviderDTOForDrift(obj *v1alpha1.AMIdentityProvider) (IdentityProvider, error) {
	dto, err := ToIdentityProviderDTO(obj)
	if err != nil {
		return IdentityProvider{}, err
	}
	if utils.SafeDereference(dto.System) {
		dto.Name = nil
		dto.Configuration = nil
		dto.Type = nil
	}
	return dto, nil
}
