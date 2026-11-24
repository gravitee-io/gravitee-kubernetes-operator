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

package service

import (
	"context"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

type Env struct {
	*client.Client
}

func NewEnv(client *client.Client) *Env {
	return &Env{Client: client}
}

// CreateGroup For tests purposes only.
func (svc *Env) CreateGroup(group *model.Group) error {
	url := svc.EnvV1Target("configuration").WithPath("groups")
	return svc.HTTP.Post(url.String(), group, group)
}

// CreateCategory For tests purposes only.
func (svc *Env) CreateCategory(category *model.Category) error {
	url := svc.EnvV1Target("configuration").WithPath("categories")
	return svc.HTTP.Post(url.String(), category, category)
}

func (svc *Env) DryRunImportGroup(ctx context.Context, grp *v1alpha1.Group) (*group.Status, error) {
	return svc.importGroup(ctx, grp, true)
}

func (svc *Env) ImportGroup(ctx context.Context, grp *v1alpha1.Group) (*group.Status, error) {
	return svc.importGroup(ctx, grp, false)
}

func (svc *Env) importGroup(ctx context.Context, grp *v1alpha1.Group, dryRun bool) (*group.Status, error) {
	body, err := model.ToGroup(grp)
	if err != nil {
		return nil, err
	}

	params := &sdk.UpsertGroupParams{DryRun: &dryRun}

	// A group created before the Automation API is addressed by its UUID, sent as the hrid.
	setHridWithUUID := grp.Spec.ID != "" && !k8s.IsAutomationAPIManaged(grp)
	if setHridWithUUID {
		body.Hrid = grp.Spec.ID
		params.HridContainsUUID = new(true)
	}

	if k8s.HasTrueAnnotation(grp, core.IgnoreMembersAnnotation) {
		params.IgnoreMembers = new(true)
	}

	client, err := svc.SDK()
	if err != nil {
		return nil, err
	}

	resp, err := client.UpsertGroupWithResponse(ctx, params, body)
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}

	if !setHridWithUUID {
		k8s.AddAutomationAPIManagedCondition(grp)
	}

	return model.ToGroupStatus(*resp.JSON200), nil
}

func (svc *Env) DeleteGroup(ctx context.Context, grp *v1alpha1.Group) error {
	id, hridContainsUUID := GetGroupID(grp)

	client, err := svc.SDK()
	if err != nil {
		return err
	}

	resp, err := client.DeleteGroupWithResponse(ctx, id, &sdk.DeleteGroupParams{HridContainsUUID: &hridContainsUUID})
	return errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	})
}

func GetGroupID(grp *v1alpha1.Group) (string, bool) {
	if k8s.IsAutomationAPIManaged(grp) {
		return refs.NewNamespacedNameFromObject(grp).HRID(), false
	}
	return grp.GetID(), true
}

func (svc *Env) Get() (*model.Env, error) {
	env := new(model.Env)
	if err := svc.HTTP.Get(svc.URLs.EnvV2.String(), env); err != nil {
		return nil, err
	}
	return env, nil
}

// GetGroupByHRID gets a group by HRID from the Automation API.
func (svc *Env) GetGroupByHRID(ctx context.Context, hrid string) (sdk.GroupState, error) {
	return svc.getGroup(ctx, hrid, false)
}

// GetGroupWithUUID gets a pre-Automation group by UUID from the Automation API.
func (svc *Env) GetGroupWithUUID(ctx context.Context, id string) (sdk.GroupState, error) {
	return svc.getGroup(ctx, id, true)
}

func (svc *Env) getGroup(ctx context.Context, hrid string, hridContainsUUID bool) (sdk.GroupState, error) {
	client, err := svc.SDK()
	if err != nil {
		return sdk.GroupState{}, err
	}

	resp, err := client.GetGroupWithResponse(ctx, hrid, &sdk.GetGroupParams{HridContainsUUID: &hridContainsUUID})
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return sdk.GroupState{}, err
	}
	if resp.JSON200 == nil {
		return sdk.GroupState{}, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}
	return *resp.JSON200, nil
}
