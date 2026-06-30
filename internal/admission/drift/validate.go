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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// RefResolver is a function that resolves all references in the CRD and can add errors. Caller is due to check them.
type RefResolver func(ctx context.Context, obj runtime.Object) error

// RemoteObjectGetter is a function that returns the remote object that will be compared with the local ones;
// it must make sure to use either Automation API or Management API.
// it returns the remote object and a boolean indicating if the object was found.
type RemoteObjectGetter func(*apim.APIM, runtime.Object, *errors.AdmissionErrors) any

// DTOMapper is a function that converts the CRD into a DTO that can be compared with the remote object.
type DTOMapper func(any) any

// MapDTO wraps a typed mapper as a DTOMapper.
func MapDTO[T any, D any](mapper func(T) D) DTOMapper {
	return func(o any) any {
		return mapper(o.(T)) //nolint:errcheck // with the generic use we are safe
	}
}

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

	if !drift.IsDriftEnabled(newCRD) {
		return errs
	}

	oldCopy := oldCRD.DeepCopyObject()
	newCopy := newCRD.DeepCopyObject()

	// We need to resolve all references to compare the content that sent to the remote API
	err := resolveRefs(ctx, oldCopy)
	if err != nil {
		errs.AddSeveref("could not resolve references for 'old' CRD: %s", err.Error())
		return errs
	}
	err = resolveRefs(ctx, newCopy)
	if err != nil {
		errs.AddSeveref("could not resolve references for 'new' CRD: %s", err.Error())
		return errs
	}

	apimClient, err := resolveContext(ctx)
	if err != nil {
		errs.AddSeveref("could not resolve context for CRD: %s", err.Error())
		return errs
	}

	remoteObject := getRemoteObject(apimClient, newCopy, errs)
	if errs.IsSevere() {
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
