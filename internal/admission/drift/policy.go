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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ApplyRemoteFetchPolicy applies OnRemoteMissing (HTTP 404) or OnFetchFailure to a failed remote fetch.
func ApplyRemoteFetchPolicy(obj client.Object, err error, errs *errors.AdmissionErrors) {
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

// CompareWithRemote compares the old and new DTOs of obj with the remote one, merges both results
// (see drift.Merge) and applies the drift policy when drift remains.
func CompareWithRemote(obj client.Object, oldDTO, newDTO, remote any, errs *errors.AdmissionErrors) {
	ns := obj.GetNamespace()
	oldVsRemote := drift.DetectWithNamespace(oldDTO, remote, ns)
	newVsRemote := drift.DetectWithNamespace(newDTO, remote, ns)

	result := drift.Merge(oldVsRemote, newVsRemote)
	if !result.DriftDetected() {
		return
	}
	applyPolicy(env.Config.DriftDetection.Policy, func() string {
		if env.Config.DriftDetection.Policy == env.DriftPolicyAllow {
			ref := client.ObjectKeyFromObject(obj)
			kind := obj.GetObjectKind().GroupVersionKind().Kind
			return fmt.Sprintf("drift detected for resource [%s] [%s], drift policy is 'allow': drift is ignored", kind, ref)
		}
		return fmt.Sprintf("\ndrift detected:\n%s", result.String())
	}, errs)
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
