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

// This file holds the addressing of applications created before APIM 4.12, by their UUID. Removing it leaves
// these to fix: the applyLegacy call in ToDTO, the request editors of the SDK calls (none), and UpdateStatus
// (always Automation API managed).

import (
	"context"
	"net/http"
	"strconv"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// Legacy flags a DTO whose Hrid is the UUID of an application created before the Automation API.
// The field is the query parameter it sets.
type Legacy struct {
	// HridContainsUUID is set when the application Hrid is its UUID.
	HridContainsUUID bool
}

// editors tells APIM that the hrid of the request is the UUID of a legacy application.
// TODO: remove once apim-sdk declares hridContainsUUID on applications (see Case 12)
func (l Legacy) editors() []sdk.RequestEditorFn {
	if !l.HridContainsUUID {
		return nil
	}
	return []sdk.RequestEditorFn{func(_ context.Context, req *http.Request) error {
		query := req.URL.Query()
		query.Set("hridContainsUUID", strconv.FormatBool(true))
		req.URL.RawQuery = query.Encode()
		return nil
	}}
}

// applyLegacy swaps in the UUID of an application created before the Automation API, and flags it.
func applyLegacy(dto *DTO, obj *v1alpha1.Application) {
	if id := legacyID(obj); id != "" {
		dto.Hrid = id
		dto.Legacy.HridContainsUUID = true
	}
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
