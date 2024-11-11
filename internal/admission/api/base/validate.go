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

package base

import (
	"context"
	"reflect"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/ctxref"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func ValidateCreate(ctx context.Context, obj runtime.Object) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	// Should be the first validation, it will also compile the templates internally
	errs.Add(admission.CompileAndValidateTemplate(ctx, obj))

	if errs.IsSevere() {
		return errs
	}

	errs.Add(ctxref.Validate(ctx, obj))

	if api, ok := obj.(core.ApiDefinitionObject); ok {
		errs.Add(validatePlans(api))
		errs.Add(validateNoConflictingPath(ctx, api))
		errs.MergeWith(validateResourceOrRefs(ctx, api))
	}

	return errs
}

func ValidateUpdate(ctx context.Context, obj runtime.Object) *errors.AdmissionErrors {
	return ValidateCreate(ctx, obj)
}

func ValidateDelete(_ context.Context, obj runtime.Object) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	if api, ok := obj.(core.ApiDefinitionObject); ok {
		errs.Add(validateSubscriptionCount(api))
	}
	return errs
}

func validateSubscriptionCount(api core.ApiDefinitionObject) *errors.AdmissionError {
	st, _ := api.GetStatus().(core.SubscribableStatus)
	sc := st.GetSubscriptionCount()
	if sc > 0 {
		return errors.NewSeveref(
			"cannot delete [%s] because it is referenced in %d subscriptions. "+
				"Subscriptions must be deleted before the API definition. "+
				"You can review the subscriptions using the following command: "+
				"kubectl get subscriptions.gravitee.io -A "+
				"-o jsonpath='{.items[?(@.spec.api.name==\"%s\")].metadata.name}'",
			api.GetRef(), sc, api.GetName(),
		)
	}
	return nil
}

func ValidateSubscribedPlans(
	ctx context.Context,
	oldApi core.ApiDefinitionObject,
	newApi core.ApiDefinitionObject,
	searchIndexField indexer.IndexField,
) *errors.AdmissionError {
	st, _ := oldApi.GetStatus().(core.SubscribableStatus)
	if st.GetSubscriptionCount() == 0 {
		return nil
	}

	subs := &v1alpha1.SubscriptionList{}
	if err := search.FindByFieldReferencing(
		ctx,
		searchIndexField,
		refs.NamespacedName{
			Name:      oldApi.GetName(),
			Namespace: oldApi.GetNamespace(),
		},
		subs,
	); err != nil {
		return errors.NewSevere(err.Error())
	}

	for _, sub := range subs.Items {
		plan := newApi.GetPlan(sub.Spec.Plan)
		if plan == nil || reflect.ValueOf(plan).IsNil() {
			return errors.NewSeveref(
				"Plan [%s] could not be found in API [%s] "+
					"but there is a subscription referencing it. "+
					"You can review the depending subscriptions using the following command: "+
					"kubectl get subscriptions.gravitee.io -A "+
					"-o jsonpath='{.items[?(@.spec.api.name==\"%s\")].metadata.name}'",
				sub.Spec.Plan, newApi.GetRef(), newApi.GetName(),
			)
		}
	}
	return nil
}
