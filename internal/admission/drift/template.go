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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/template"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// CompileTemplates resolves the templates of the old copy drift compares, without recording any
// reference on the Secrets and ConfigMaps it reads, so that it holds the values the remote was sent.
// Validation compiles the new object in place before drift runs, never the old one: left raw, a
// templated old value never matches the remote and every change to it is reported as drift.
//
// ctx carries the secret recorder of the drift check (template.WithSecretRecorder). The old copy's
// secret values are recorded, and so are those of the new object, read from the admission request
// since the object drift receives is already compiled.
//
// It returns false when drift must not be checked (see ReportUnresolvedOldVersion).
func CompileTemplates(ctx context.Context, oldCopy client.Object, errs *errors.AdmissionErrors) bool {
	if err := compileSubmitted(ctx); err != nil {
		errs.AddSeveref("could not compile templates for 'new' CRD: %s", err.Error())
		return false
	}
	if err := template.Compile(ctx, oldCopy, false); err != nil {
		ReportUnresolvedOldVersion(err, errs)
		return false
	}
	return true
}

// ReportUnresolvedOldVersion adds why the old version could not be prepared for comparison. Something
// it references that no longer exists, a Secret or an ApiResource deleted for instance, leaves no way
// to know what the old version asked for: a warning says so rather than refusing the very update that
// repairs the reference. Any other failure, the API server not answering for instance, refuses it.
func ReportUnresolvedOldVersion(err error, errs *errors.AdmissionErrors) {
	if apierrors.IsNotFound(err) {
		errs.AddWarningf("drift detection skipped, the previous version no longer resolves: %s", err.Error())
		return
	}
	errs.AddSeveref("could not resolve the previous version for drift detection: %s", err.Error())
}

// compileSubmitted compiles the object under admission as the request carried it, for the secret
// values it records only. There is none outside an admission request.
func compileSubmitted(ctx context.Context) error {
	req, err := admission.RequestFromContext(ctx)
	if err != nil || len(req.Object.Raw) == 0 {
		return nil
	}
	submitted := &unstructured.Unstructured{}
	if err := submitted.UnmarshalJSON(req.Object.Raw); err != nil {
		return err
	}
	return template.Compile(ctx, submitted, false)
}
