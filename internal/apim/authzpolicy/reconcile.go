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

package authzpolicy

import (
	"context"
	"fmt"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

// Response is what Upsert returns: the state the platform answers.
type Response struct {
	State sdk.AuthzPolicyState
}

// GetOrgID is the organization the state reports, as the status has always shown it.
func (r Response) GetOrgID() string {
	return valueOf(r.State.OrganizationId)
}

// GetEnvID is the environment the state reports, as the status has always shown it.
func (r Response) GetEnvID() string {
	return valueOf(r.State.EnvironmentId)
}

// ResolveRefs resolves, outside a delete, each pdpRef of the targets onto the routing scope the platform receives,
// then the management context. obj is the lifecycle's working copy, never saved, so the stored policy keeps its
// pdpRefs. Missing, deleting or unsynced AuthzPdps are waited out (ResolvedRefs=False names every one, requeued):
// the AuthzPdp watch requeues the policy once they sync. An AuthzPdp of another context is an IllegalStateError,
// reported on Accepted and not requeued: only a change to the pdpRef fixes it. A delete goes by hrid, so a pending
// AuthzPdp does not block it.
func ResolveRefs(ctx context.Context, obj *v1alpha1.AuthzPolicy, ns string) error {
	if !obj.IsBeingDeleted() {
		if err := dynamic.ResolveAuthzPolicyRefs(ctx, obj); err != nil {
			return err
		}
	}
	return apim.ResolveContextRef(ctx, obj, ns)
}

// Upsert creates or updates the policy and deploys it to its targets: every apply deploys. The platform answers the
// state: ids, findings and the policy status. A refusal is a 400 whose body carries the findings, see
// model.AutomationRefusal.
func Upsert(ctx context.Context, client *apim.SDKClient, dto sdk.AuthzPolicy) (Response, error) {
	state, err := upsert(ctx, client, dto, false)
	if err != nil {
		return Response{}, err
	}
	return Response{State: state}, nil
}

func upsert(ctx context.Context, client *apim.SDKClient, dto sdk.AuthzPolicy, dryRun bool) (
	sdk.AuthzPolicyState, error,
) {
	resp, err := client.UpsertAuthzPolicyWithResponse(ctx, &sdk.UpsertAuthzPolicyParams{DryRun: &dryRun}, dto)
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return sdk.AuthzPolicyState{}, err
	}
	if resp.JSON200 == nil {
		return sdk.AuthzPolicyState{}, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}
	return *resp.JSON200, nil
}

// OnUpsertError records the findings of a refused apply in the status: an entityId another policy holds, a text
// that does not parse, a violation of a strict gateway's schema. Admission skips the dry run while a pdpRef is
// pending, and another policy may take the entityId since: it stays recoverable (backoff) rather than terminal
// until the next spec change. The policy status is left as it was. Any other failure keeps the lifecycle's
// classification.
func OnUpsertError(_ context.Context, obj *v1alpha1.AuthzPolicy, err error) error {
	if !errors.IsBadRequest(err) {
		return nil
	}
	if findings, refused := model.AutomationRefusal(err); refused {
		obj.Status.Errors = findings
	}
	return fmt.Errorf("authz policy [%s] was refused by the automation api: %w", obj.GetName(), err)
}

// Delete withdraws the policy from the PDP gateways it targets and deletes it by hrid; the entities it names are
// kept. A policy already gone counts as deleted. Withdrawing a deployed policy needs the ENVIRONMENT_AUTHZ_PDP
// UPDATE permission: the platform answers 403 without it, and the lifecycle keeps the finalizer and retries, with
// Accepted naming the reason, until the context's token is granted it. dryRun is never sent: APIM 4.13 ignores it
// on a DELETE and deletes.
func Delete(ctx context.Context, client *apim.SDKClient, dto sdk.AuthzPolicy) error {
	resp, err := client.DeleteAuthzPolicyWithResponse(ctx, dto.Hrid, nil)
	return errors.IgnoreNotFound(errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}))
}

// UpdateStatus copies the upsert state into the CR status, field by field so that the conditions are kept, and
// marks the policy Automation API managed.
func UpdateStatus(_ context.Context, obj *v1alpha1.AuthzPolicy, resp Response) error {
	obj.Status.ID = valueOf(resp.State.Id)
	obj.Status.OrgID = resp.GetOrgID()
	obj.Status.EnvID = resp.GetEnvID()
	obj.Status.HRID = refs.NewNamespacedNameFromObject(obj).HRID()
	obj.Status.PolicyStatus = policyStatusOf(resp.State.PolicyStatus)
	obj.Status.Errors = ToStatusErrors(resp.State.Errors)

	k8s.AddAutomationAPIManagedCondition(obj)
	return nil
}

// ToStatusErrors maps the findings of a state onto the CR status.
func ToStatusErrors(errs *sdk.AuthzErrors) status.Errors {
	if errs == nil {
		return status.Errors{}
	}
	return status.Errors{Severe: errs.Severe, Warning: errs.Warning}
}

func policyStatusOf(s *sdk.AuthzPolicyStatePolicyStatus) string {
	if s == nil {
		return ""
	}
	return string(*s)
}

func valueOf(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
