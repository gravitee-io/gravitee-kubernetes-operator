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
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RefResolver is a function that resolves all references in the CRD and can add errors. Caller is due to check them.
type RefResolver[T client.Object] func(ctx context.Context, obj T) error

// RemoteObjectGetter is a function that returns the remote object that will be compared with the local ones;
// it must make sure to use either Automation API or Management API.
type RemoteObjectGetter[T client.Object] func(*apim.APIM, T, *errors.AdmissionErrors) any

// DTOMapper is a function that converts the CRD into a DTO that can be compared with the remote object.
type DTOMapper[T client.Object] func(T) any

// MapDTO wraps a typed mapper as a DTOMapper.
func MapDTO[T client.Object, D any](mapper func(T) D) DTOMapper[T] {
	return func(o T) any {
		return mapper(o)
	}
}

type ContextResolver func(ctx context.Context) (*apim.APIM, error)

func ValidateDrift[T core.ContextAwareObject](
	ctx context.Context,
	oldCRD T,
	newCRD T,
	resolveRefs RefResolver[T],
	getRemoteObject RemoteObjectGetter[T],
	dtoMapper DTOMapper[T]) *errors.AdmissionErrors {
	return ValidateDriftWithContext(ctx, oldCRD, newCRD, func(ctx context.Context) (*apim.APIM, error) {
		return apim.FromContextRef(ctx, newCRD.ContextRef(), newCRD.GetNamespace())
	}, resolveRefs, getRemoteObject, dtoMapper)
}

func ValidateDriftWithContext[T client.Object](
	ctx context.Context,
	oldCRD T,
	newCRD T,
	resolveContext ContextResolver,
	resolveRefs RefResolver[T],
	getRemoteObject RemoteObjectGetter[T],
	mapDTO DTOMapper[T]) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if !drift.IsDriftEnabled(newCRD) {
		return errs
	}

	oldCopy := oldCRD.DeepCopyObject().(T)
	newCopy := newCRD.DeepCopyObject().(T)

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

	ns := newCopy.GetNamespace()

	oldVsRemoteResult := drift.DetectWithNamespace(oldDTO, remoteObject, ns)
	newVsRemoteResult := drift.DetectWithNamespace(newDTO, remoteObject, ns)

	if result := drift.Merge(oldVsRemoteResult, newVsRemoteResult); result.DriftDetected() {
		errs.AddSeveref("\ndrift detected:\n%s", result.String())
	}

	return errs
}
