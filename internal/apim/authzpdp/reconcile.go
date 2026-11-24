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

package authzpdp

import (
	"context"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/authz"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"k8s.io/utils/ptr"
)

// Response is what Upsert returns: the state the platform answers.
type Response struct {
	State sdk.AuthzPdpState
}

// GetOrgID is the organization the state reports, as the status has always shown it.
func (r Response) GetOrgID() string {
	return ptr.Deref(r.State.OrganizationId, "")
}

// GetEnvID is the environment the state reports, as the status has always shown it.
func (r Response) GetEnvID() string {
	return ptr.Deref(r.State.EnvironmentId, "")
}

// Upsert creates or updates the PDP gateway and its endpoint. The platform answers the state: ids, findings,
// the routing scope and the AuthZEN discovery path. A refusal is a 400 whose body carries the findings, see
// model.AutomationRefusal.
func Upsert(ctx context.Context, client *apim.SDKClient, dto sdk.AuthzPdp) (Response, error) {
	state, err := upsert(ctx, client, dto, false)
	if err != nil {
		return Response{}, err
	}
	return Response{State: state}, nil
}

func upsert(ctx context.Context, client *apim.SDKClient, dto sdk.AuthzPdp, dryRun bool) (sdk.AuthzPdpState, error) {
	resp, err := client.UpsertAuthzPdpWithResponse(ctx, &sdk.UpsertAuthzPdpParams{DryRun: &dryRun}, dto)
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return sdk.AuthzPdpState{}, err
	}
	if resp.JSON200 == nil {
		return sdk.AuthzPdpState{}, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}
	return *resp.JSON200, nil
}

// OnUpsertError records the findings of a refused apply in the status. Admission ran the same checks as a dry
// run, so a refusal here comes from a state that changed since (an endpoint path another API took, a sibling's
// strict): it stays recoverable (backoff) rather than terminal until the next spec change. The routing scope is
// left as it was: a PDP synced before still exists. Any other failure keeps the lifecycle's classification.
func OnUpsertError(_ context.Context, obj *v1alpha1.AuthzPdp, err error) error {
	return authz.OnRefusal("authz pdp", obj, &obj.Status.Errors, err)
}

// Delete removes the PDP gateway and its endpoint. A PDP already gone counts as deleted. The platform answers
// 409, deleting nothing, while the delete would leave a policy or an entity with no target: the lifecycle keeps
// the finalizer and retries, on backoff and whenever one of them changes or goes (watch.WatchAuthzPdpReleases).
// dryRun is never sent: APIM 4.13 ignores it on a DELETE and deletes.
func Delete(ctx context.Context, client *apim.SDKClient, dto sdk.AuthzPdp) error {
	resp, err := client.DeleteAuthzPdpWithResponse(ctx, dto.Hrid, nil)
	return errors.IgnoreNotFound(errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}))
}

// UpdateStatus copies the upsert state into the CR status, field by field so that the conditions are kept,
// and marks the PDP Automation API managed.
func UpdateStatus(_ context.Context, obj *v1alpha1.AuthzPdp, resp Response) error {
	obj.Status.ID = ptr.Deref(resp.State.Id, "")
	obj.Status.OrgID = resp.GetOrgID()
	obj.Status.EnvID = resp.GetEnvID()
	obj.Status.HRID = refs.NewNamespacedNameFromObject(obj).HRID()
	obj.Status.Errors = authz.ToStatusErrors(resp.State.Errors)
	obj.Status.RoutingScope = ptr.Deref(resp.State.RoutingScope, "")
	obj.Status.WellKnownPath = ptr.Deref(resp.State.WellKnownPath, "")

	k8s.AddAutomationAPIManagedCondition(obj)
	return nil
}
