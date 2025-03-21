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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

func Accept(params *v1alpha1.GatewayClassParameters) {
	condition := k8s.NewAcceptedConditionBuilder(params.Generation)
	if k8s.IsResolved(params) {
		k8s.SetCondition(params, condition.Accept("Parameters have been accepted").Build())
		return
	}
	resolved := k8s.GetCondition(params, k8s.ConditionResolvedRefs)
	k8s.SetCondition(params, condition.Reason(resolved.Reason).Message(resolved.Message).Build())
}
