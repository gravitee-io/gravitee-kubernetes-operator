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

// Package group maps the Group CRD to the APIM Automation API through the SDK. Its functions have the
// internal/lifecycle signatures, so the Group controller and webhook can move onto the lifecycle as is.
package group

import (
	"context"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/mapper"
)

// DTO is the group payload with what Upsert, GetRemote and Delete need to know about the CR, since they
// never see it. The flags are not sent (Upsert sends the embedded sdk.Group) and drift detection ignores them.
type DTO struct {
	sdk.Group
	// HridContainsUUID is set for a group created before the Automation API: its UUID is sent as the hrid.
	HridContainsUUID bool `json:"-" drift:"ignore"`
	// IgnoreMembers is set by the gravitee.io/ignore-members annotation.
	IgnoreMembers bool `json:"-" drift:"ignore"`
}

// ToDTO maps a Group CRD to the Automation API payload. SDK defaults are not applied: they differ from
// the CRD's and would defeat the drift equivalences. The CRD id is not part of the payload; a legacy group
// sends it as its hrid.
func ToDTO(obj *v1alpha1.Group) (DTO, error) {
	payload, err := mapper.MapViaJSON[sdk.Group](obj.Spec.Type)
	if err != nil {
		return DTO{}, err
	}

	dto := DTO{
		Group:         payload,
		IgnoreMembers: k8s.HasTrueAnnotation(obj, core.IgnoreMembersAnnotation),
	}
	if id := legacyID(obj); id != "" {
		dto.Hrid = id
		dto.HridContainsUUID = true
	} else {
		dto.Hrid = refs.NewNamespacedNameFromObject(obj).HRID()
	}
	return dto, nil
}

// legacyID returns the UUID of a group created before the Automation API, or "" for a group addressed by its hrid.
func legacyID(obj *v1alpha1.Group) string {
	if k8s.IsAutomationAPIManaged(obj) {
		return ""
	}
	if obj.Status.ID != "" {
		return obj.Status.ID
	}
	return obj.Spec.ID
}

// CreateAPIMClient builds the SDK client of the management context the group references.
func CreateAPIMClient(ctx context.Context, obj *v1alpha1.Group) (*apim.SDKClient, error) {
	return apim.SDKFromContextRef(ctx, obj.ContextRef(), obj.GetNamespace())
}

// Response is what Upsert returns: the group state APIM answers, and the org and env of the management
// context it was sent to.
type Response struct {
	State sdk.GroupState
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
