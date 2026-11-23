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

// Package application maps the Application CRD to the APIM Automation API through the SDK. Its functions have
// the internal/lifecycle signatures, so the Application controller and webhook can move onto the lifecycle as is.
package application

import (
	"context"
	"net/http"
	"strconv"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/mapper"
)

// defaultMetadataFormat is the format APIM gives a metadata entry created without one.
const defaultMetadataFormat = sdk.MetadataFormat("STRING")

// DTO is the application payload with what Upsert, GetRemote and Delete need to know about the CR, since they
// never see it. The flags are not sent (Upsert sends the embedded sdk.Application) and drift detection ignores them.
type DTO struct {
	sdk.Application
	// HridContainsUUID is set for an application created before the Automation API: its UUID is sent as the hrid.
	HridContainsUUID bool `json:"-" drift:"ignore"`
	// IgnoreGroups is set by the gravitee.io/ignore-groups annotation.
	IgnoreGroups bool `json:"-" drift:"ignore"`
}

// ToDTO maps an Application CRD to the Automation API payload. SDK defaults are not applied: they differ from
// the CRD's and would defeat the drift equivalences. The CRD id is not part of the payload; a legacy application
// sends it as its hrid. picture and metadata hidden are not in the payload: APIM never applied them.
func ToDTO(obj *v1alpha1.Application) (DTO, error) {
	payload, err := mapper.MapViaJSON[sdk.Application](obj.Spec.Application)
	if err != nil {
		return DTO{}, err
	}

	// TODO: remove once apim-sdk omits an empty format (see Case 12)
	for i := range payload.Metadata {
		if payload.Metadata[i].Format == "" {
			payload.Metadata[i].Format = defaultMetadataFormat
		}
	}

	dto := DTO{
		Application:  payload,
		IgnoreGroups: k8s.HasTrueAnnotation(obj, core.IgnoreGroupsAnnotation),
	}
	if id := legacyID(obj); id != "" {
		dto.Hrid = id
		dto.HridContainsUUID = true
	} else {
		dto.Hrid = refs.NewNamespacedNameFromObject(obj).HRID()
	}
	return dto, nil
}

// legacyID returns the UUID of an application created before the Automation API, or "" for an application
// addressed by its hrid.
func legacyID(obj *v1alpha1.Application) string {
	if k8s.IsAutomationAPIManaged(obj) {
		return ""
	}
	if obj.Status.ID != "" {
		return obj.Status.ID
	}
	return obj.Spec.ID
}

// withHridContainsUUID tells APIM that the hrid of the request is the UUID of a legacy application.
// TODO: remove once apim-sdk declares hridContainsUUID on applications (see Case 12)
func withHridContainsUUID(dto DTO) []sdk.RequestEditorFn {
	if !dto.HridContainsUUID {
		return nil
	}
	return []sdk.RequestEditorFn{func(_ context.Context, req *http.Request) error {
		query := req.URL.Query()
		query.Set("hridContainsUUID", strconv.FormatBool(true))
		req.URL.RawQuery = query.Encode()
		return nil
	}}
}

// CreateAPIMClient builds the SDK client of the management context the application references.
func CreateAPIMClient(ctx context.Context, obj *v1alpha1.Application) (*apim.SDKClient, error) {
	return apim.SDKFromContextRef(ctx, obj.ContextRef(), obj.GetNamespace())
}

// Response is what Upsert returns: the application state APIM answers, and the org and env of the management
// context it was sent to.
type Response struct {
	State sdk.ApplicationState
	OrgID string
	EnvID string
	// HridContainsUUID is the flag the upsert was sent with: a legacy application does not become
	// Automation API managed.
	HridContainsUUID bool
}

func (r Response) GetOrgID() string {
	return r.OrgID
}

func (r Response) GetEnvID() string {
	return r.EnvID
}
