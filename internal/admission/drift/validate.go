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
	"k8s.io/apimachinery/pkg/runtime"
)

// RefResolver is a function that resolves all references in the CRD;
// it used on old and new CRD to compare the content that will be sent to the remote API.
type RefResolver func(ctx context.Context, obj runtime.Object, errs *errors.AdmissionErrors)

// RemoteObjectGetter is a function that returns the remote object that will be compared with the local ones;
// it must make sure to use either Automation API or Management API.
// it returns the remote object and a boolean indicating if the object was found.
type RemoteObjectGetter func(*apim.APIM, runtime.Object, *errors.AdmissionErrors) (any, bool)

// DTOMapper is a function that converts the CRD into a DTO that can be compared with the remote object.
type DTOMapper func(any) any

type ContextResolver func(ctx context.Context) (*apim.APIM, error)

func ValidateDrift(
	ctx context.Context,
	oldCRD core.ContextAwareObject,
	newCRD core.ContextAwareObject,
	resolveRefs RefResolver,
	getRemoteObject RemoteObjectGetter,
	dtoMapper DTOMapper) *errors.AdmissionErrors {
	return ValidateDriftWithContext(ctx, oldCRD, newCRD, func(ctx context.Context) (*apim.APIM, error) {
		return apim.FromContextRef(ctx, newCRD.ContextRef(), newCRD.GetNamespace())
	}, resolveRefs, getRemoteObject, dtoMapper)
}

func ValidateDriftWithContext(
	ctx context.Context,
	oldCRD runtime.Object,
	newCRD runtime.Object,
	resolveContext ContextResolver,
	resolveRefs RefResolver,
	getRemoteObject RemoteObjectGetter,
	mapDTO DTOMapper) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if isDriftEnabled(newCRD) {
		return errs
	}

	oldCopy := oldCRD.DeepCopyObject()
	newCopy := newCRD.DeepCopyObject()

	// We need template to be compiled
	if err := admission.CompileAndValidateTemplate(ctx, oldCopy); err != nil {
		errs.AddWarningf("could not compile templates of existing CRD, drift might be detected: %s", err.Error())
	}
	if err := admission.CompileAndValidateTemplate(ctx, newCopy); err != nil {
		errs.AddWarningf("could not compile templates of updated CRD, drift might be detected: %s", err.Error())
	}

	// We need to resolve all references to compare the content that sent to the remote API
	resolveRefs(ctx, oldCopy, errs)
	resolveRefs(ctx, newCopy, errs)

	apimClient, err := resolveContext(ctx)
	if err != nil {
		errs.AddSeveref("could not resolve context for CRD: %s", err.Error())
		return errs
	}

	remoteObject, ok := getRemoteObject(apimClient, newCopy, errs)
	if !ok {
		return errs
	}

	oldDTO := mapDTO(oldCopy)
	newDTO := mapDTO(newCopy)

	oldVsRemoteResult := drift.Detect(oldDTO, remoteObject)
	newVsRemoteResult := drift.Detect(newDTO, remoteObject)

	if result := drift.Merge(oldVsRemoteResult, newVsRemoteResult); result.DriftDetected() {
		errs.AddSeveref("\ndrift detected:\n%s", result.String())
	}

	return errs
}

func isDriftEnabled(newCRD runtime.Object) bool {
	coreObj, _ := newCRD.(core.Object)
	driftAnnot, ok := coreObj.GetAnnotations()[core.DriftDetectionAnnotation]
	if ok && driftAnnot == env.TrueString {
		return true
	} else if ok && driftAnnot == env.FalseString {
		return false
	}
	return env.Config.DriftDetection
}
