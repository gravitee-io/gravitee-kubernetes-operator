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
	"encoding/json"
	goerrors "errors"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/ctxref"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// validateCreate runs the checks in order: a context is required, it must resolve (this also
// compiles templates so the dry run carries resolved credentials), then the platform dry run,
// which discovers the upstream and refuses an unreachable server or a taken entityId. Schema
// and CEL own the spec-only rules (entityId grammar and immutability, auth variant pairing).
func validateCreate(ctx context.Context, srv *v1alpha1.CatalogMcpServer) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if !srv.HasContext() {
		errs.AddSevere("a management context reference (spec.contextRef) is required")
		return errs
	}

	errs.Add(ctxref.Validate(ctx, srv))
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateDryRun(ctx, srv))
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

// refusalToAdmissionErrors turns a platform refusal into the findings the author reads. A 400
// answers the resource state with the findings in errors.severe (credentials already removed),
// so those are reported one by one; anything else is reported as the request error.
func refusalToAdmissionErrors(err error) *errors.AdmissionErrors {
	if errors.IsBadRequest(err) {
		serverError := &errors.ServerError{}
		if goerrors.As(err, serverError) {
			state := new(model.CatalogMcpServerState)
			if jsonErr := json.Unmarshal([]byte(serverError.Body), state); jsonErr == nil &&
				len(state.Errors.Severe) > 0 {
				return errors.NewAdmissionErrorsFromStatus(state.Errors)
			}
		}
	}

	errs := errors.NewAdmissionErrors()
	errs.AddSevere(err.Error())
	return errs
}
