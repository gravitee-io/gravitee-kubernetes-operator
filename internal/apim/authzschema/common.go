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

// Package authzschema maps the AuthzSchema CRD to the AuthZ schema documents of the APIM Automation API through
// the SDK.
package authzschema

import (
	"context"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/authz"
)

// ToDTO maps a schema whose refs dynamic.ResolveAuthzSchemaRefs has resolved onto the payload a PUT sends: its
// target is a scope, or absent. It is the single mapping used by the sync path, the dry run and drift detection,
// and reads nothing from the cluster. The platform stores schema and target as sent, so both compare exactly.
// An unresolved pdpRef is an error: sent as no target, it would publish the document to the default engine.
func ToDTO(obj *v1alpha1.AuthzSchema) (sdk.AuthzSchema, error) {
	target, err := authz.ToTarget(obj, obj.Spec.Target)
	if err != nil {
		return sdk.AuthzSchema{}, err
	}

	return sdk.AuthzSchema{
		Hrid:   refs.NewNamespacedNameFromObject(obj).HRID(),
		Name:   obj.Spec.Name,
		Schema: obj.Spec.Schema,
		Target: target,
	}, nil
}

// CreateAPIMClient builds the SDK client of the management context the schema references.
func CreateAPIMClient(ctx context.Context, obj *v1alpha1.AuthzSchema) (*apim.SDKClient, error) {
	return apim.SDKFromContextRef(ctx, obj.ContextRef(), obj.GetNamespace())
}
