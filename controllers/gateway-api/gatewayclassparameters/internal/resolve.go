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

package internal

import (
	"context"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"sigs.k8s.io/controller-runtime/pkg/client"

	coreV1 "k8s.io/api/core/v1"
	kErrors "k8s.io/apimachinery/pkg/api/errors"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func Resolve(ctx context.Context, params *v1alpha1.GatewayClassParameters) error {
	condition := k8s.NewResolvedRefsConditionBuilder(params.Generation)

	if params.Spec.Gravitee == nil || params.Spec.Gravitee.LicenseRef == nil {
		k8s.SetCondition(params, condition.Message("No license to resolve").Build())
		return nil
	}

	ref := params.Spec.Gravitee.LicenseRef

	ns := ref.Namespace
	if ns == nil {
		gwNs := gwAPIv1.Namespace(params.Namespace)
		ns = &gwNs
	}

	key := client.ObjectKey{Namespace: string(*ns), Name: string(ref.Name)}
	secret := &coreV1.Secret{}

	if err := k8s.GetClient().Get(ctx, key, secret); client.IgnoreNotFound(err) != nil {
		return err
	} else if kErrors.IsNotFound(err) {
		condition.RejectLicenseNotFound(
			fmt.Sprintf("License secret [%s] could not be resolved", key.String()),
		)
		k8s.SetCondition(params, condition.Build())
		return nil
	}

	k8s.SetCondition(params, condition.Message("Gravitee license has been resolved").Build())

	return nil
}
