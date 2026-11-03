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

package clientsecret

import (
	"context"
	"net/http"

	amsdk "github.com/gravitee-io/gravitee-automation-sdk/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// keyMaxLength is the longest key AM accepts for a client secret, shorter than for the other AM resources.
const keyMaxLength = 64

// PreCheck validates the client secret:
// key pattern and length;
// applicationRef.namespace must match the resource's own namespace;
// If the application, the Secret or its entry is not found, a warning is issued.
func PreCheck(ctx context.Context, obj *v1alpha1.AMClientSecret) *gerrors.AdmissionErrors {
	errs := gerrors.NewAdmissionErrors()

	key := refs.NewNamespacedNameFromObject(obj).HRID()
	errs.MergeWith(am.ValidateKey(key))
	if errs.IsSevere() {
		return errs
	}
	if len(key) > keyMaxLength {
		errs.AddSeveref("key [%s] (namespace-name) is longer than the %d characters AM allows for a client secret; "+
			"shorten the resource name", key, keyMaxLength)
		return errs
	}

	if ns := obj.Spec.ApplicationRef.Namespace; ns != "" && ns != obj.Namespace {
		errs.AddSeveref("AMClientSecret namespace is [%s] applicationRef is [%s]: cross-namespace is not supported",
			obj.Namespace, ns)
		return errs
	}

	errs.MergeWith(warnMissingRefs(ctx, obj))
	return errs
}

// warnMissingRefs warns about a dependency that does not exist yet. The messages come from ResolveRefs:
// they name the Secret and the key, never the value.
func warnMissingRefs(ctx context.Context, obj *v1alpha1.AMClientSecret) *gerrors.AdmissionErrors {
	errs := gerrors.NewAdmissionErrors()
	if err := ResolveRefs(ctx, obj.DeepCopy(), obj.Namespace); isMissingRef(err) {
		errs.AddWarningf("%s; the client secret will be created in AM once it exists", err)
	}
	return errs
}

// DryRun validates the client secret against AM without persisting it. It is skipped without a client
// (application or domain not in AM yet) or without a value (Secret or entry missing). Any AM error, a
// different value for an existing secret or a value too short included, rejects it with AM's message.
func DryRun(ctx context.Context, client *am.Client, dto ClientSecret) *gerrors.AdmissionErrors {
	errs := gerrors.NewAdmissionErrors()
	if client == nil || dto.Value == nil {
		return errs
	}

	resp, err := client.UpsertClientSecretWithResponse(ctx, dto.DomainKey, dto.ApplicationKey,
		&amsdk.UpsertClientSecretParams{DryRun: new(true)}, dto.ClientSecret)

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

// GetRemote fetches the client secret for drift detection. Without a client (application or domain not
// in AM yet) nothing can have drifted: the resource is its own remote, so an update is admitted.
func GetRemote(ctx context.Context, client *am.Client, dto ClientSecret) (ClientSecret, error) {
	if client == nil {
		return dto, nil
	}
	resp, err := client.GetClientSecretWithResponse(ctx, dto.DomainKey, dto.ApplicationKey, dto.Key)
	if err = am.HasErrors(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return ClientSecret{}, err
	}
	if resp.JSON200 == nil {
		return ClientSecret{}, am.UnexpectedResponse(resp.HTTPResponse)
	}
	return ClientSecret{ClientSecret: *resp.JSON200, DomainKey: dto.DomainKey, ApplicationKey: dto.ApplicationKey}, nil
}
