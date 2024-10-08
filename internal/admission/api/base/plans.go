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

package base

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func validatePlans(api core.ApiDefinitionObject) *errors.AdmissionError {
	if !api.HasPlans() && api.GetState() != "STOPPED" {
		return errors.NewSeveref(
			"cannot apply API [%s]. "+
				"Its state is set to STARTED, but the API has no plans. "+
				"APIs must have at least one plan in order to be deployed.",
			api.GetName(),
		)
	}
	return nil
}
