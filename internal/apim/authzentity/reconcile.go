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

package authzentity

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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"k8s.io/utils/ptr"
)

// Response is what Upsert returns: the state the platform answers.
type Response struct {
	State sdk.AuthzEntityState
}

// GetOrgID is the organization the state reports, as the status has always shown it.
func (r Response) GetOrgID() string {
	return ptr.Deref(r.State.OrganizationId, "")
}

// GetEnvID is the environment the state reports, as the status has always shown it.
func (r Response) GetEnvID() string {
	return ptr.Deref(r.State.EnvironmentId, "")
}

// ResolveRefs resolves, outside a delete, each pdpRef of the targets onto the routing scope the platform receives,
// then the management context. obj is the lifecycle's working copy, never saved, so the stored entity keeps its
// pdpRefs. Missing, deleting or unsynced AuthzPdps are waited out (ResolvedRefs=False names every one, requeued):
// the AuthzPdp watch requeues the entity once they sync. An AuthzPdp of another context is an IllegalStateError,
// reported on Accepted and not requeued: only a change to the pdpRef fixes it. A delete goes by hrid: it drops the
// targets from the working copy instead, so that ToDTO maps it whatever the state of its AuthzPdps.
func ResolveRefs(ctx context.Context, obj *v1alpha1.AuthzEntity, ns string) error {
	if obj.IsBeingDeleted() {
		obj.Spec.Targets = nil
	} else if err := dynamic.ResolveAuthzEntityRefs(ctx, obj); err != nil {
		return err
	}
	return apim.ResolveContextRef(ctx, obj, ns)
}

// Upsert creates or updates the entity and publishes it to its targets. The platform answers the state: ids and
// findings. A refusal is a 400 whose body carries the findings, see model.AutomationRefusal.
func Upsert(ctx context.Context, client *apim.SDKClient, dto sdk.AuthzEntity) (Response, error) {
	state, err := upsert(ctx, client, dto, false)
	if err != nil {
		return Response{}, err
	}
	return Response{State: state}, nil
}

func upsert(ctx context.Context, client *apim.SDKClient, dto sdk.AuthzEntity, dryRun bool) (
	sdk.AuthzEntityState, error,
) {
	resp, err := client.UpsertAuthzEntityWithResponse(ctx, &sdk.UpsertAuthzEntityParams{DryRun: &dryRun}, dto)
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return sdk.AuthzEntityState{}, err
	}
	if resp.JSON200 == nil {
		return sdk.AuthzEntityState{}, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}
	return *resp.JSON200, nil
}

// OnUpsertError records the findings of a refused apply in the status: an entityId another source holds
// (AUTHZ-5), an attribute off a strict schema, a parent the platform deleted. Admission skips the dry run while a
// pdpRef is pending, and another source may take the entityId since: it stays recoverable (backoff) rather than
// terminal until the next spec change. Any other failure keeps the lifecycle's classification.
func OnUpsertError(_ context.Context, obj *v1alpha1.AuthzEntity, err error) error {
	return authz.OnRefusal("authz entity", obj, &obj.Status.Errors, err)
}

// Delete deletes the entity by hrid. The platform deletes its RESOURCE children with it and detaches its PRINCIPAL
// children. An entity already gone counts as deleted. Over the platform's cascade limit, it answers 413 and
// deletes nothing: the lifecycle keeps the finalizer and retries, and Accepted names the reason until the cascade
// fits. dryRun is never sent: APIM 4.13 ignores it on a DELETE and deletes.
func Delete(ctx context.Context, client *apim.SDKClient, dto sdk.AuthzEntity) error {
	resp, err := client.DeleteAuthzEntityWithResponse(ctx, dto.Hrid, nil)
	return errors.IgnoreNotFound(errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}))
}

// UpdateStatus copies the upsert state into the CR status, field by field so that the conditions are kept, and
// marks the entity Automation API managed.
func UpdateStatus(_ context.Context, obj *v1alpha1.AuthzEntity, resp Response) error {
	obj.Status.ID = ptr.Deref(resp.State.Id, "")
	obj.Status.OrgID = resp.GetOrgID()
	obj.Status.EnvID = resp.GetEnvID()
	obj.Status.HRID = refs.NewNamespacedNameFromObject(obj).HRID()
	obj.Status.Errors = authz.ToStatusErrors(resp.State.Errors)

	k8s.AddAutomationAPIManagedCondition(obj)
	return nil
}
