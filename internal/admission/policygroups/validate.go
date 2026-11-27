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

package policygroups

import (
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func validateImmutableFields(oldSpg, newSpg *v1alpha1.SharedPolicyGroup) *errors.AdmissionError {
	if utils.ToStringValue(oldSpg.Spec.CrossID) != utils.ToStringValue(newSpg.Spec.CrossID) {
		return &errors.AdmissionError{
			Severity: errors.Severe,
			Message: fmt.Sprintf("can not change Shared Policy Group [%s] CrossID, "+
				"once it is created", oldSpg.Status.CrossID),
		}
	}

	if oldSpg.Spec.ApiType != newSpg.Spec.ApiType {
		return &errors.AdmissionError{
			Severity: errors.Severe,
			Message: fmt.Sprintf("can not change Shared Policy Group [%s] ApiType [%s], "+
				"once it is created", oldSpg.Status.CrossID, oldSpg.Spec.ApiType),
		}
	}

	if *oldSpg.Spec.Phase != *newSpg.Spec.Phase {
		return &errors.AdmissionError{
			Severity: errors.Severe,
			Message: fmt.Sprintf("can not change Shared Policy Group [%s] Phase [%s], "+
				"once it is created", oldSpg.Status.CrossID, *oldSpg.Spec.Phase),
		}
	}

	return nil
}
