// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"encoding/json"
	goerrors "errors"
	"net/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// AutomationRefusal reads the findings of a refused apply. A gamma module refuses with a 400
// whose body is the resource state, credentials removed, with the findings in errors.severe; a
// 400 the host answers before the module sees the request has another shape and is not a refusal.
func AutomationRefusal(err error) (status.Errors, bool) {
	serverError := &errors.ServerError{}
	if !goerrors.As(err, serverError) || serverError.StatusCode != http.StatusBadRequest {
		return status.Errors{}, false
	}

	state := new(struct {
		Errors status.Errors `json:"errors"`
	})
	if json.Unmarshal([]byte(serverError.Body), state) != nil || len(state.Errors.Severe) == 0 {
		return status.Errors{}, false
	}

	return state.Errors, true
}
