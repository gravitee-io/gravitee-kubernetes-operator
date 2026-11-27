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

// Package sharedpolicygroup maps the SharedPolicyGroup CRD to the APIM Automation API through the SDK. Its
// functions have the internal/lifecycle signatures.
package sharedpolicygroup

import (
	"context"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/mapper"
)

// DTO is the shared policy group payload with what Upsert, GetRemote and Delete need to know about the CR,
// since they never see it. The flag is not sent (Upsert sends the embedded sdk.SharedPolicyGroup) and drift
// detection ignores it.
type DTO struct {
	sdk.SharedPolicyGroup
	// HridContainsUUID is set for a group created before the Automation API: its cross id is sent as the hrid.
	HridContainsUUID bool `json:"-" drift:"ignore"`
}

// ToDTO maps a SharedPolicyGroup CRD to the Automation API payload. SDK defaults are not applied: they would
// defeat the drift equivalences. A step without a policy is sent with an empty policy, which APIM refuses.
func ToDTO(obj *v1alpha1.SharedPolicyGroup) (DTO, error) {
	cp := obj.DeepCopy()
	cp.PopulateIDs(nil, k8s.IsAutomationAPIManaged(cp))

	payload, err := mapper.MapViaJSON[sdk.SharedPolicyGroup](cp.Spec.SharedPolicyGroup)
	if err != nil {
		return DTO{}, err
	}

	dto := DTO{SharedPolicyGroup: payload}
	if id := legacyID(cp); id != "" {
		dto.Hrid = id
		dto.HridContainsUUID = true
	}
	return dto, nil
}

// legacyID returns the cross id of a group created before the Automation API, or "" for a group addressed by
// its hrid.
func legacyID(obj *v1alpha1.SharedPolicyGroup) string {
	if k8s.IsAutomationAPIManaged(obj) {
		return ""
	}
	return obj.GetID()
}

// CreateAPIMClient builds the SDK client of the management context the group references.
func CreateAPIMClient(ctx context.Context, obj *v1alpha1.SharedPolicyGroup) (*apim.SDKClient, error) {
	return apim.SDKFromContextRef(ctx, obj.ContextRef(), obj.GetNamespace())
}

// Response is what Upsert returns: the state APIM answers, and the org and env of the management context it was
// sent to.
type Response struct {
	State sdk.SharedPolicyGroupState
	OrgID string
	EnvID string
	// HridContainsUUID is the flag the upsert was sent with: a legacy group does not become Automation API managed.
	HridContainsUUID bool
}

func (r Response) GetOrgID() string {
	return r.OrgID
}

func (r Response) GetEnvID() string {
	return r.EnvID
}
