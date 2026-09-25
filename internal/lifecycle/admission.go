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

package lifecycle

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission"
	admissiondrift "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func (a AdmissionLifecycle[T, D, C]) ValidateCreate(ctx context.Context, obj T) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	a.templateAndRefs(ctx, obj, errs)
	if errs.IsSevere() {
		return errs
	}

	api, ok := a.resolveContext(ctx, obj, errs)
	if !ok {
		return errs
	}

	a.preCheck(ctx, obj, errs)
	if errs.IsSevere() {
		return errs
	}

	a.dryRun(ctx, api, obj, errs)
	if errs.IsSevere() {
		return errs
	}

	a.postCheck(ctx, obj, errs)
	return errs
}

func (a AdmissionLifecycle[T, D, C]) ValidateUpdate(
	ctx context.Context, oldObj, newObj T,
) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	if newObj.IsBeingDeleted() {
		return errs
	}

	a.templateAndRefs(ctx, newObj, errs)
	if errs.IsSevere() {
		return errs
	}

	apiClient, ok := a.resolveContext(ctx, newObj, errs)
	if !ok {
		return errs
	}

	a.preCheck(ctx, newObj, errs)
	if errs.IsSevere() {
		return errs
	}

	if a.ImmutableFields != nil {
		errs.MergeWith(a.ImmutableFields(oldObj, newObj))
		if errs.IsSevere() {
			return errs
		}
	}

	a.dryRun(ctx, apiClient, newObj, errs)
	if errs.IsSevere() {
		return errs
	}

	a.postCheck(ctx, newObj, errs)
	if errs.IsSevere() {
		return errs
	}

	a.detectDrift(ctx, apiClient, oldObj, newObj, errs)
	return errs
}

func (a AdmissionLifecycle[T, D, C]) ValidateDelete(ctx context.Context, obj T) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	if a.DeleteGuard != nil {
		if err := a.DeleteGuard(ctx, obj); err != nil {
			errs.AddSevere(err.Error())
		}
	}
	return errs
}

func (a AdmissionLifecycle[T, D, C]) templateAndRefs(
	ctx context.Context, obj T, errs *errors.AdmissionErrors,
) {
	errs.Add(admission.CompileAndValidateTemplate(ctx, obj))
	if errs.IsSevere() {
		return
	}
	if err := a.resolveRefs(ctx, obj); err != nil {
		errs.AddSeveref("could not resolve references: %s", err.Error())
	}
}

// resolveContext validates the context ref by building the API client.
// Returns the client and true on success. On failure adds a severe error and returns zero C, false.
// When ClientFactory is nil (no context-dependent operations), returns zero C, true.
func (a AdmissionLifecycle[T, D, C]) resolveContext(
	ctx context.Context, obj T, errs *errors.AdmissionErrors,
) (C, bool) {
	var zero C
	if a.ClientFactory == nil {
		return zero, true
	}
	api, err := a.ClientFactory(ctx, obj)
	if err != nil {
		errs.AddSeveref("could not resolve context: %s", err.Error())
		return zero, false
	}
	return api, true
}

func (a AdmissionLifecycle[T, D, C]) preCheck(
	ctx context.Context, obj T, errs *errors.AdmissionErrors,
) {
	if a.PreCheck == nil {
		return
	}
	errs.MergeWith(a.PreCheck(ctx, obj))
}

func (a AdmissionLifecycle[T, D, C]) postCheck(
	ctx context.Context, obj T, errs *errors.AdmissionErrors,
) {
	if a.PostCheck == nil {
		return
	}
	errs.MergeWith(a.PostCheck(ctx, obj))
}

func (a AdmissionLifecycle[T, D, C]) dryRun(
	ctx context.Context, apiClient C, obj T, errs *errors.AdmissionErrors,
) {
	if a.DryRun == nil {
		return
	}
	dto := a.ToDTO(obj)
	if err := a.DryRun(ctx, apiClient, dto); err != nil {
		errs.MergeWith(err)
	}
}

func (a AdmissionLifecycle[T, D, C]) detectDrift(
	ctx context.Context, apiClient C, oldObj, newObj T, errs *errors.AdmissionErrors,
) {
	if a.GetRemote == nil {
		return
	}
	if !drift.IsDriftEnabled(newObj) {
		return
	}

	oldCopy, ok := oldObj.DeepCopyObject().(T)
	if !ok {
		errs.AddSeveref("lifecycle: DeepCopyObject is not %T", oldObj)
		return
	}
	newCopy, ok := newObj.DeepCopyObject().(T)
	if !ok {
		errs.AddSeveref("lifecycle: DeepCopyObject is not %T", newObj)
		return
	}

	if err := a.resolveRefs(ctx, oldCopy); err != nil {
		errs.AddSeveref("could not resolve references for old CRD: %s", err.Error())
		return
	}
	if err := a.resolveRefs(ctx, newCopy); err != nil {
		errs.AddSeveref("could not resolve references for new CRD: %s", err.Error())
		return
	}

	newDTO := a.ToDTO(newCopy)
	remote, err := a.GetRemote(ctx, apiClient, newDTO)
	if err != nil {
		admissiondrift.ApplyRemoteFetchPolicy(newCopy, err, errs)
		return
	}

	admissiondrift.CompareWithRemote(newCopy, a.ToDTO(oldCopy), newDTO, remote, errs)
}

func (a AdmissionLifecycle[T, D, C]) resolveRefs(ctx context.Context, obj T) error {
	ns := obj.GetNamespace()
	if a.ResolveRefs == nil {
		return nil
	}
	return a.ResolveRefs(ctx, obj, ns)
}
