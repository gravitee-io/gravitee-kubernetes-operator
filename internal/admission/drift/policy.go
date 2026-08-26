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

package drift

import (
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func applyRemoteFetchPolicy(obj client.Object, err error, errs *errors.AdmissionErrors) {
	ref := client.ObjectKeyFromObject(obj)
	kind := obj.GetObjectKind().GroupVersionKind().Kind
	if errors.IsNotFound(err) {
		applyPolicy(
			env.Config.DriftDetection.OnRemoteMissing,
			func() string { return fmt.Sprintf("remote [%s] [%s] not found during drift detection", kind, ref) },
			errs,
		)
		return
	}
	applyPolicy(
		env.Config.DriftDetection.OnFetchFailure,
		func() string {
			return fmt.Sprintf("failed to fetch remote [%s] [%s] during drift detection: %s", kind, ref, err.Error())
		},
		errs,
	)
}

func applyPolicy(policy env.DriftPolicy, message func() string, errs *errors.AdmissionErrors) {
	switch policy {
	case env.DriftPolicyWarn:
		errs.AddWarning(message())
	case env.DriftPolicyAllow:
		log.Global.Warn(message())
	default:
		errs.AddSevere(message())
	}
}
