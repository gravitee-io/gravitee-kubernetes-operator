// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package catalogmcpserver

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/ctxref"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// validateCreate runs the checks in order: the context, then the platform dry run, which
// discovers the upstream and refuses an unreachable server or a taken entityId. Schema and CEL
// own the spec-only rules (entityId grammar and immutability, auth variant pairing).
func validateCreate(ctx context.Context, srv *v1alpha1.CatalogMcpServer) *errors.AdmissionErrors {
	errs := validateContext(ctx, srv)
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateDryRun(ctx, srv))
	return errs
}

// validateContext requires a context that resolves. Resolving also compiles templates, so a dry
// run that follows carries resolved credentials.
func validateContext(ctx context.Context, srv *v1alpha1.CatalogMcpServer) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if !srv.HasContext() {
		errs.AddSevere("a management context reference (spec.contextRef) is required")
		return errs
	}

	errs.Add(ctxref.Validate(ctx, srv))
	return errs
}

func validateDryRun(ctx context.Context, srv *v1alpha1.CatalogMcpServer) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	cp := srv.DeepCopy()

	apimClient, err := apim.FromContextRef(ctx, cp.ContextRef(), cp.GetNamespace())
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	state, err := apimClient.CatalogMcpServers.DryRunCreateOrUpdate(cp)
	if err != nil {
		errs.MergeWith(refusalToAdmissionErrors(err))
		return errs
	}

	errs.MergeWith(errors.NewAdmissionErrorsFromStatus(state.Errors))

	return errs
}

// refusalToAdmissionErrors reports a platform refusal finding by finding; anything else is
// reported as the request error.
func refusalToAdmissionErrors(err error) *errors.AdmissionErrors {
	if findings, refused := model.CatalogMcpServerRefusal(err); refused {
		return errors.NewAdmissionErrorsFromStatus(findings)
	}

	errs := errors.NewAdmissionErrors()
	errs.AddSevere(err.Error())
	return errs
}
