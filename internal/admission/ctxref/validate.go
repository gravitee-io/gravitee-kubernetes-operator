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

package ctxref

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"k8s.io/apimachinery/pkg/runtime"
)

func Validate(ctx context.Context, obj runtime.Object) *errors.AdmissionError {
	if ctxAware, ok := obj.(core.ContextAwareObject); ok {
		if ctxAware.HasContext() {
			return validateContextRefExists(ctx, ctxAware)
		}
	}
	return nil
}

func validateContextRefExists(ctx context.Context, ctxAware core.ContextAwareObject) *errors.AdmissionError {
	ctxRef := ctxAware.ContextRef()

	// Should be the first validation, it will also compile the templates internally
	tmpErr := admission.CompileAndValidateTemplate(ctx, ctxAware)
	if tmpErr != nil {
		return tmpErr
	}

	if err := dynamic.ExpectResolvedContext(ctx, ctxRef, ctxAware.GetNamespace()); err != nil {
		return errors.NewSeveref(
			"resource [%s] references management context [%v] that doesn't exist in the cluster",
			ctxAware.GetName(),
			ctxRef,
		)
	}
	return nil
}
