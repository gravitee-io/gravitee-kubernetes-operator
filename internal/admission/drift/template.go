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
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CompileTemplates resolves the templates of the copies drift compares, without recording any
// reference on the Secrets and ConfigMaps it reads, so that both sides hold the values the remote
// was sent. Validation compiles the new object in place, never the old one: left raw, a templated
// old value never matches the remote and every change to it is reported as drift.
//
// It returns false when drift must not be checked. A template of the old object that no longer
// resolves, its Secret deleted for instance, leaves no way to know what the old object asked for:
// a warning says so rather than refusing the very update that repairs the reference.
func CompileTemplates(ctx context.Context, oldCopy, newCopy client.Object, errs *errors.AdmissionErrors) bool {
	if err := template.Compile(ctx, newCopy, false); err != nil {
		errs.AddSeveref("could not compile templates for 'new' CRD: %s", err.Error())
		return false
	}
	if err := template.Compile(ctx, oldCopy, false); err != nil {
		errs.AddWarningf("drift detection skipped, the templates of the previous version no longer resolve: %s",
			err.Error())
		return false
	}
	return true
}
