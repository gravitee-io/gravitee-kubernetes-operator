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

package v4

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
)

type Status struct {
	base.Status `json:",inline"`
	// This field is used to store the list of plans that have been created
	// for the API definition if a management context has been defined
	// to sync the API with an APIM instance
	Plans map[string]string `json:"plans,omitempty"`
	// When API has been created regardless of errors, this field is
	// used to persist the error message encountered during admission
	Errors Errors `json:"errors,omitempty"`
}

type Errors struct {
	// warning errors do not block object reconciliation,
	// most of the time because the value is ignored or defaulted
	// when the API gets synced with APIM
	Warning []string `json:"warning,omitempty"`
	// severe errors do not pass admission and will block reconcile
	// hence, this field should always be during the admission phase
	// and is very unlikely to be persisted in the status
	Severe []string `json:"severe,omitempty"`
}
