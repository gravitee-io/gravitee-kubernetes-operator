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
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

func ValidateDrift(
	ctx context.Context,
	oldCRD core.ContextAwareObject,
	newCRD core.ContextAwareObject,
	resolveRefs func(context.Context, core.ContextAwareObject, *errors.AdmissionErrors),
	getRemoteObject func(*apim.APIM, core.ContextAwareObject, *errors.AdmissionErrors) (any, bool),
	dtoMapper func(any) any) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if !isDriftEnabled(newCRD) {
		return errs
	}

	oldCopy, _ := oldCRD.DeepCopyObject().(core.ContextAwareObject)
	newCopy, _ := newCRD.DeepCopyObject().(core.ContextAwareObject)

	// We need template to be compiled
	if err := admission.CompileAndValidateTemplate(ctx, oldCopy); err != nil {
		errs.AddWarningf("could not compile templates of existing CRD, drift might be detected: %s", err.Error())
	}
	if err := admission.CompileAndValidateTemplate(ctx, newCopy); err != nil {
		errs.AddWarningf("could not compile templates of updated CRD, drift might be detected: %s", err.Error())
	}

	// We must have a context to perform drift detection
	apimClient, err := apim.FromContextRef(ctx, newCopy.ContextRef(), newCopy.GetNamespace())
	if err != nil {
		errs.AddSeveref("Cannot perform drift detection without context: %s", err.Error())
		return errs
	}

	// We need to populate IDs so that remote object can be fetched
	oldAware, _ := oldCopy.(core.ConditionAware)
	newAware, _ := newCopy.(core.ConditionAware)
	oldCopy.PopulateIDs(apimClient.Context, k8s.IsAutomationAPIManaged(oldAware))
	newCopy.PopulateIDs(apimClient.Context, k8s.IsAutomationAPIManaged(newAware))

	// We need to resolve all references to compare the content that sent to the remote API
	resolveRefs(ctx, oldCopy, errs)
	resolveRefs(ctx, newCopy, errs)

	remoteObject, ok := getRemoteObject(apimClient, newCopy, errs)
	if !ok {
		return errs
	}

	oldDTO := dtoMapper(oldCopy)
	newDTO := dtoMapper(newCopy)

	oldVsRemoteResult := drift.Detect(oldDTO, remoteObject)
	newVsRemoteResult := drift.Detect(newDTO, remoteObject)

	if result := drift.Merge(oldVsRemoteResult, newVsRemoteResult); result.DriftDetected() {
		errs.AddSeveref("\ndrift detected:\n%s", result.String())
	}

	return errs
}

func isDriftEnabled(newCRD core.Object) bool {
	driftAnnot, ok := newCRD.GetAnnotations()[core.DriftDetectionAnnotation]
	if ok && driftAnnot == env.TrueString {
		return true
	} else if ok && driftAnnot == env.FalseString {
		return false
	}
	return env.Config.DriftDetection
}
