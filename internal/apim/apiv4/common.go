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

// Package apiv4 syncs ApiV4Definition resources with APIM through the generated APIM SDK.
package apiv4

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/mapper"
)

// gatewayDefinitionVersion is the apiDefinitionVersion written in the gateway ConfigMap of a v4 API.
const gatewayDefinitionVersion = "4.0.0"

// DTO is an ApiV4Definition as APIM sees it, plus what the calls around it need. The embedded sdk.ApiV4 is what the
// PUT sends and what drift detection compares.
type DTO struct {
	sdk.ApiV4
	// HridContainsUUID is set for an API with a spec.id (created before the Automation API, or imported):
	// its UUID is sent as the hrid.
	HridContainsUUID bool `json:"-" drift:"ignore"`
	// DeleteID and DeleteHridContainsUUID address the API on delete. They follow model.APIUsesUUID, which reads
	// the status ID too, so they can differ from the upsert identity.
	DeleteID               string `json:"-" drift:"ignore"`
	DeleteHridContainsUUID bool   `json:"-" drift:"ignore"`
	// IgnoreGroups is set by the gravitee.io/ignore-groups annotation.
	IgnoreGroups bool `json:"-" drift:"ignore"`
	// Synced is set once APIM gave the API an ID (status.id): drift has nothing to compare before that.
	Synced bool `json:"-" drift:"ignore"`
	// PlanOrder is the rank of each plan (hrid to position, from 1) in the list APIM stores, drift only: APIM
	// orders plans by position and the SDK plan has no order. It is only set, and compared, when every plan has
	// an order and is published (APIM renumbers published plans only).
	PlanOrder map[string]int `json:"-"`
	// Gateway is the ConfigMap the gateway syncs the API from.
	Gateway GatewaySync `json:"-" drift:"ignore"`
}

// GatewaySync is the gateway ConfigMap of the API.
type GatewaySync struct {
	// Enabled is set when the API syncs from Kubernetes (definitionContext.syncFrom KUBERNETES). When unset, a
	// ConfigMap the API owns is deleted.
	Enabled   bool
	Name      string
	Namespace string
	Owner     metav1.OwnerReference
	Stopped   bool
	// Definition is the JSON gateway definition of the API, IDs populated.
	Definition []byte
}

// ToDTO maps an ApiV4Definition, references already resolved, to its APIM payload.
func ToDTO(obj *v1alpha1.ApiV4Definition) (DTO, error) {
	cp := obj.DeepCopy()
	cp.PopulateIDsWithContext(cp.HasContext(), k8s.IsAutomationAPIManaged(cp))

	// ToAPIV4DTO keeps the plan send order and the page mapping GKO has always sent.
	payload := model.ToAPIV4DTO(&cp.Spec.Api)
	api, err := mapper.MapViaJSON[sdk.ApiV4](payload)
	if err != nil {
		return DTO{}, err
	}

	dto := DTO{
		ApiV4:        api,
		IgnoreGroups: k8s.HasTrueAnnotation(obj, core.IgnoreGroupsAnnotation),
		Synced:       obj.Status.ID != "",
	}
	if len(payload.Plans) > 0 && !slices.ContainsFunc(payload.Plans, orderNotComparable) {
		dto.PlanOrder = planRanks(dto.Plans)
	}
	if cp.HasContext() && cp.Spec.ID != "" {
		dto.Hrid = cp.Spec.ID
		dto.HridContainsUUID = true
	}
	dto.DeleteID, dto.DeleteHridContainsUUID = deleteIdentity(obj)

	if dto.Gateway, err = toGatewaySync(cp); err != nil {
		return DTO{}, err
	}
	return dto, nil
}

func orderNotComparable(plan *model.APIV4PlanDTO) bool {
	published := plan.Status == "" || plan.Status == base.PublishedPlanStatus
	return plan.Order == nil || !published
}

// planRanks gives each plan its position in plans, from 1: the order APIM stores for a plan list it receives,
// and lists them in on GET.
func planRanks(plans []sdk.PlanV4) map[string]int {
	ranks := make(map[string]int, len(plans))
	for i, plan := range plans {
		ranks[plan.Hrid] = i + 1
	}
	return ranks
}

func deleteIdentity(obj *v1alpha1.ApiV4Definition) (string, bool) {
	if model.APIUsesUUID(obj) {
		return obj.GetID(), true
	}
	return refs.NewNamespacedNameFromObject(obj).HRID(), false
}

func toGatewaySync(cp *v1alpha1.ApiV4Definition) (GatewaySync, error) {
	gvk := cp.GetObjectKind().GroupVersionKind()
	sync := GatewaySync{
		Enabled:   !cp.IsSyncFromManagement(),
		Name:      cp.GetName(),
		Namespace: cp.GetNamespace(),
		Owner: metav1.OwnerReference{
			Kind:       gvk.Kind,
			Name:       cp.GetName(),
			APIVersion: gvk.GroupVersion().String(),
			UID:        cp.GetUID(),
		},
		Stopped: cp.IsStopped(),
	}
	if !sync.Enabled {
		return sync, nil
	}
	definition, err := json.Marshal(cp.ToGatewayDefinition())
	if err != nil {
		return GatewaySync{}, fmt.Errorf("unable to marshal the gateway definition: %w", err)
	}
	sync.Definition = definition
	return sync, nil
}

// Client is the APIM SDK client of the API's management context, plus GKO's APIM facade for the one call the
// SDK cannot make yet (APIs.UpsertV4Legacy).
type Client struct {
	*apim.SDKClient
	APIM *apim.APIM
}

// CreateAPIMClient builds the client of the API's management context. An API without a context has no
// client: nil, nil.
func CreateAPIMClient(ctx context.Context, obj *v1alpha1.ApiV4Definition) (*Client, error) {
	if !obj.HasContext() {
		return nil, nil
	}
	mctx, err := dynamic.ResolveContext(ctx, obj.ContextRef(), obj.GetNamespace())
	if err != nil {
		return nil, err
	}
	sdkClient, err := apim.NewSDKClient(ctx, mctx)
	if err != nil {
		return nil, err
	}
	facade, err := apim.FromContext(ctx, mctx, obj.GetNamespace())
	if err != nil {
		return nil, err
	}
	return &Client{SDKClient: sdkClient, APIM: facade}, nil
}

// Response is the outcome of an upsert. State is nil for an API without a management context.
type Response struct {
	State *sdk.ApiV4State
	OrgID string
	EnvID string
	// HridContainsUUID is the flag the upsert was sent with: a legacy API does not become Automation API managed.
	HridContainsUUID bool
}

func (r Response) GetOrgID() string {
	return r.OrgID
}

func (r Response) GetEnvID() string {
	return r.EnvID
}

// PayloadChanged reports a change the controller sends to APIM: the spec, or the gravitee.io/ignore-groups
// annotation (LastSpecHashPredicate.Update reconciles on both).
func PayloadChanged(oldObj, newObj *v1alpha1.ApiV4Definition) bool {
	return hash.Calculate(&oldObj.Spec) != hash.Calculate(&newObj.Spec) ||
		oldObj.GetAnnotations()[core.IgnoreGroupsAnnotation] != newObj.GetAnnotations()[core.IgnoreGroupsAnnotation]
}
