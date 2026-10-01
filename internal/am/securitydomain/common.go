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

package securitydomain

import (
	"context"
	"fmt"

	domain "github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/mapper"
)

// ToDomainDTO maps an AMSecurityDomain to the AM payload, keyed by its namespace-name HRID.
func ToDomainDTO(obj *v1alpha1.AMSecurityDomain) (domain.Domain, error) {
	dto, err := mapper.MapViaJSON[domain.Domain](obj.Spec.Domain)
	if err != nil {
		return domain.Domain{}, err
	}
	dto.Key = domainKey(obj)
	return dto, nil
}

func domainKey(obj *v1alpha1.AMSecurityDomain) string {
	return refs.NewNamespacedNameFromObject(obj).HRID()
}

// CreateAMClient builds the AM client from the domain's AMContext, resolved with its templates compiled.
func CreateAMClient(ctx context.Context, obj *v1alpha1.AMSecurityDomain) (*am.Client, error) {
	if !obj.HasContext() {
		return nil, fmt.Errorf("contextRef empty on %s [%s/%s]", obj.Kind, obj.GetName(), obj.GetNamespace())
	}

	// resolved like the APIM contexts: templates compiled, fetched from the API server
	resolved, err := dynamic.ResolveAMContext(ctx, obj.ContextRef(), obj.GetNamespace())
	if err != nil {
		return nil, fmt.Errorf("AMContext [%s]: %w", obj.ContextRef().String(), err)
	}
	amContext, ok := resolved.(*v1alpha1.AMContext)
	if !ok {
		return nil, fmt.Errorf("AMContext [%s]: unexpected type %T", obj.ContextRef().String(), resolved)
	}

	return am.NewSDKClient(ctx, amContext)
}
