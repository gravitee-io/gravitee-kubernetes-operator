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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/gateway"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	kErrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func Accept(ctx context.Context, gwc *gateway.GatewayClass) error {
	condition := k8s.NewAcceptedConditionBuilder(gwc.Object.Generation)

	paramRef := gwc.Object.Spec.ParametersRef

	if paramRef == nil {
		condition.Accept("No parameters reference")
		k8s.SetCondition(gwc, condition.Accept("No parameters reference").Build())
		return nil
	}

	if paramRef.Group != gAPIv1.Group(v1alpha1.GroupVersion.Group) {
		k8s.SetCondition(
			gwc,
			condition.
				RejectInvalidParameters("parameters reference group must be gravitee.io").
				Build(),
		)
		return nil
	}

	if paramRef.Kind != "GatewayClassParameters" {
		k8s.SetCondition(
			gwc,
			condition.
				RejectInvalidParameters("parameters reference kind must be GatewayClassParameters").
				Build(),
		)
		return nil
	}

	if gwc.Object.Spec.ParametersRef.Namespace == nil {
		k8s.SetCondition(
			gwc,
			condition.
				RejectInvalidParameters("parameters reference must be namespaced").
				Build(),
		)
		return nil
	}

	key := client.ObjectKey{
		Name:      gwc.Object.Spec.ParametersRef.Name,
		Namespace: string(*gwc.Object.Spec.ParametersRef.Namespace),
	}

	params := new(v1alpha1.GatewayClassParameters)

	err := k8s.GetClient().Get(ctx, key, params)
	if client.IgnoreNotFound(err) != nil {
		return err
	}

	if kErrors.IsNotFound(err) {
		k8s.SetCondition(
			gwc,
			condition.
				RejectInvalidParameters("parameters reference not found").
				Build(),
		)
		return nil
	}

	if !k8s.IsAccepted(params) {
		k8s.SetCondition(
			gwc,
			condition.
				RejectInvalidParameters("parameters are not accepted").
				Build(),
		)
		return nil
	}

	k8s.SetCondition(
		gwc,
		condition.Accept("gateway class has been accepted").Build(),
	)

	return nil
}
