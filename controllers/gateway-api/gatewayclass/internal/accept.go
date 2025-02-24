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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	kErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func Accept(ctx context.Context, gwc *gAPIv1.GatewayClass) (*metav1.Condition, error) {
	condition := k8s.NewAcceptedConditionBuilder(gwc.Generation)

	paramRef := gwc.Spec.ParametersRef

	if paramRef == nil {
		condition.Accept("No parameters reference")
		return condition.Build(), nil
	}

	if paramRef.Group != gAPIv1.Group(v1alpha1.GroupVersion.Group) {
		condition.RejectInvalidParameters("parameters reference group must be gravitee.io")
		return condition.Build(), nil
	}

	if paramRef.Kind != "GatewayClassParameters" {
		condition.RejectInvalidParameters("parameters reference kind must be GatewayClassParameters")
		return condition.Build(), nil
	}

	if gwc.Spec.ParametersRef.Namespace == nil {
		condition.RejectInvalidParameters("parameters reference must be namespaced")
		return condition.Build(), nil
	}

	key := client.ObjectKey{
		Name:      gwc.Spec.ParametersRef.Name,
		Namespace: string(*gwc.Spec.ParametersRef.Namespace),
	}

	gw := new(v1alpha1.GatewayClassParameters)

	err := k8s.GetClient().Get(ctx, key, gw)
	if client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	if kErrors.IsNotFound(err) {
		condition.RejectInvalidParameters("parameters reference not found")
		return condition.Build(), nil
	}

	return condition.Accept("gateway class has been accepted").Build(), nil
}
