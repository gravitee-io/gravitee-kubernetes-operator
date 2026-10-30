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
	"net/http"

	amsdk "github.com/gravitee-io/gravitee-automation-sdk/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

const refNotFoundMessage = "%s [%s] not found; [%s] will be created in AM once it exists"
const refError = "%s [%s] could not be read: %s"

// PreCheck validates the application
// key pattern and length;
// domainRef.namespace must match the resource's own namespace;
// If the domain is not found or ready, a warning is issued.
func PreCheck(ctx context.Context, obj *v1alpha1.AMApplication) *gerrors.AdmissionErrors {
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

	if !obj.IsBeingDeleted() {
		errs.MergeWith(validateCertificateRef(ctx, obj))
		errs.MergeWith(validateIdentityProvidersRef(ctx, obj))
	}

	return errs
}

func validateCertificateRef(ctx context.Context, obj *v1alpha1.AMApplication) *gerrors.AdmissionErrors {
	errs := gerrors.NewAdmissionErrors()
	if obj.Spec.CertificateRef != nil {
		if obj.Spec.CertificateRef.Namespace != "" && obj.Spec.CertificateRef.Namespace != obj.Namespace {
			errs.AddSeveref("AMApplication namespace is [%s] certificateRef is [%s]: cross-namespace is not supported", obj.Namespace, obj.Spec.CertificateRef.Namespace)
			return errs
		}
		if cert, err := resolveCertificate(ctx, obj.Spec.CertificateRef, obj.Namespace); apierrors.IsNotFound(err) {
			errs.AddWarningf(refNotFoundMessage, "AMCertificate", obj.Spec.CertificateRef.Name, obj.Name)
		} else if err != nil {
			errs.AddWarningf(refError, "AMCertificate", obj.Spec.CertificateRef.Name, err)
		} else if err := validateIsInSameDomain("AMCertificate", cert, obj); err != nil {
			errs.AddSevere(err.Error())
		}
	}
	return errs
}

func validateIdentityProvidersRef(ctx context.Context, obj *v1alpha1.AMApplication) *gerrors.AdmissionErrors {
	errs := gerrors.NewAdmissionErrors()
	for i, provider := range obj.Spec.IdentityProviders {
		if provider.IdentityRef != nil {
			if provider.IdentityRef.Namespace != "" && provider.IdentityRef.Namespace != obj.Namespace {
				errs.AddSeveref("AMApplication namespace is [%s] identityRef[%d] is [%s]: cross-namespace is not supported", obj.Namespace, i, provider.IdentityRef.Namespace)
				return errs
			}
			if idp, err := resolveIdentityProvider(ctx, provider.IdentityRef, obj.Namespace); apierrors.IsNotFound(err) {
				errs.AddWarningf(refNotFoundMessage, "AMIdentityProvider", provider.IdentityRef.Name, obj.Name)
			} else if err != nil {
				errs.AddWarningf(refError, "AMIdentityProvider", provider.IdentityRef.Name, err)
			} else if err := validateIsInSameDomain("AMIdentityProvider", idp, obj); err != nil {
				errs.AddSevere(err.Error())
			}
		}
	}
	return errs
}

// DryRun validates the application against AM without persisting it. It is skipped without a client
// (domain not in AM yet). Any AM error, an unreachable AM, a duplicate clientId or a type change to or from
// agent included, rejects it with AM's message.
func DryRun(ctx context.Context, client *am.Client, dto Application) *gerrors.AdmissionErrors {
	errs := gerrors.NewAdmissionErrors()
	if client == nil {
		return errs
	}

	resp, err := client.UpsertApplicationWithResponse(ctx, dto.DomainKey,
		&amsdk.UpsertApplicationParams{DryRun: new(true)}, dto.Application)

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

// GetRemote fetches the application for drift detection. Without a client (domain not in AM
// yet) nothing can have drifted: the resource is its own remote, so an update is admitted, as a create is.
func GetRemote(ctx context.Context, client *am.Client, dto Application) (Application, error) {
	if client == nil {
		return dto, nil
	}
	resp, err := client.GetApplicationWithResponse(ctx, dto.DomainKey, dto.Key)
	if err = am.HasErrors(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return Application{}, err
	}
	if resp.JSON200 == nil {
		return Application{}, am.UnexpectedResponse(resp.HTTPResponse)
	}
	return Application{Application: *resp.JSON200, DomainKey: dto.DomainKey}, nil
}
