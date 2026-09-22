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
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle/ref"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
		applyRemoteFetchPolicy(newCopy, err, errs)
		return
	}

	oldDTO := a.ToDTO(oldCopy)
	ns := newCopy.GetNamespace()
	oldVsRemote := drift.DetectWithNamespace(oldDTO, remote, ns)
	newVsRemote := drift.DetectWithNamespace(newDTO, remote, ns)

	if result := drift.Merge(oldVsRemote, newVsRemote); result.DriftDetected() {
		applyDriftPolicy(env.Config.DriftDetection.Policy, func() string {
			if env.Config.DriftDetection.Policy == env.DriftPolicyAllow {
				objRef := client.ObjectKeyFromObject(newObj)
				kind := newObj.GetObjectKind().GroupVersionKind().Kind
				return fmt.Sprintf(
					"drift detected for resource [%s] [%s], drift policy is 'allow': drift is ignored",
					kind, objRef,
				)
			}
			return fmt.Sprintf("\ndrift detected:\n%s", result.String())
		}, errs)
	}
}

func (a AdmissionLifecycle[T, D, C]) resolveRefs(ctx context.Context, obj T) error {
	ns := obj.GetNamespace()
	if a.ResolveRefs != nil {
		return a.ResolveRefs(ctx, obj, ns)
	}
	return ref.GenericRefResolver(ctx, obj, ns)
}

func applyRemoteFetchPolicy(obj client.Object, err error, errs *errors.AdmissionErrors) {
	objRef := client.ObjectKeyFromObject(obj)
	kind := obj.GetObjectKind().GroupVersionKind().Kind
	if errors.IsNotFound(err) {
		applyDriftPolicy(
			env.Config.DriftDetection.OnRemoteMissing,
			func() string { return fmt.Sprintf("remote [%s] [%s] not found during drift detection", kind, objRef) },
			errs,
		)
		return
	}
	applyDriftPolicy(
		env.Config.DriftDetection.OnFetchFailure,
		func() string {
			return fmt.Sprintf("failed to fetch remote [%s] [%s] during drift detection: %s", kind, objRef, err.Error())
		},
		errs,
	)
}

func applyDriftPolicy(policy env.DriftPolicy, message func() string, errs *errors.AdmissionErrors) {
	switch policy {
	case env.DriftPolicyWarn:
		errs.AddWarning(message())
	case env.DriftPolicyAllow:
		log.Global.Warn(message())
	default:
		errs.AddSevere(message())
	}
}
