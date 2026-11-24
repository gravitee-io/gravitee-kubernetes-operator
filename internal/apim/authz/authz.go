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

// Package authz holds what the AuthZ resources share when they map onto the APIM Automation API: the findings of
// a state or a refused apply, and the targets a schema, an entity or a policy sends.
package authz

import (
	"fmt"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ToStatusErrors maps the findings of a state onto the CR status.
func ToStatusErrors(errs *sdk.AuthzErrors) status.Errors {
	if errs == nil {
		return status.Errors{}
	}
	return status.Errors{Severe: errs.Severe, Warning: errs.Warning}
}

// OnRefusal records in errs the findings of an apply the Automation API refused with a 400, and returns the
// error to report, naming the resource: it stays recoverable. It returns nil for any other failure, which keeps
// the lifecycle's classification.
func OnRefusal(resource string, obj metav1.Object, errs *status.Errors, err error) error {
	if !errors.IsBadRequest(err) {
		return nil
	}
	if findings, refused := model.AutomationRefusal(err); refused {
		*errs = findings
	}
	return fmt.Errorf("%s [%s] was refused by the automation api: %w", resource, obj.GetName(), err)
}

// ToTarget returns the scope target sends once the refs have resolved: nil without a target or a scope. A pdpRef
// left is an error: sent as no target, it would publish to the default engine.
func ToTarget(obj metav1.Object, target *authzpdp.Target) (*string, error) {
	switch {
	case target == nil:
		return nil, nil
	case target.PdpRef != nil:
		return nil, unresolved(obj, "target", *target)
	case target.Scope != "":
		return new(target.Scope), nil
	}
	return nil, nil
}

// ToTargets returns the scopes targets send once the refs have resolved, in declaration order. A pdpRef left is
// an error: it has no scope to send.
func ToTargets(obj metav1.Object, targets []authzpdp.Target) ([]string, error) {
	var scopes []string
	for i, target := range targets {
		if target.PdpRef != nil {
			return nil, unresolved(obj, fmt.Sprintf("targets[%d]", i), target)
		}
		scopes = append(scopes, target.Scope)
	}
	return scopes, nil
}

func unresolved(obj metav1.Object, field string, target authzpdp.Target) error {
	ref := dynamic.AuthzPdpRef(obj, *target.PdpRef)
	return fmt.Errorf("%s.pdpRef [%s] is not resolved", field, ref.String())
}
