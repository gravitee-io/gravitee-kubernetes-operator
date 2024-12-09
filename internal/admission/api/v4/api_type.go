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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func validateApiType(oldApi, newAPI core.ApiDefinitionObject) *errors.AdmissionError {
	if nativeAPI == oldApi.GetType() && oldApi.GetType() != newAPI.GetType() {
		return errors.NewSeveref("it is not possible to change the API type 'NATIVE' "+
			"to something else [%s]", newAPI.GetType())
	} else if nativeAPI != oldApi.GetType() && nativeAPI == newAPI.GetType() {
		return errors.NewSevere("it is not possible to convert a none NATIVE API to a NATIVE API")
	}

	return nil
}
