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
)

func UpdateCondition(ctx context.Context, mgtCtx *v1alpha1.ManagementContext, err error) error {
	if err != nil {
		k8s.ErrorToCondition(mgtCtx, err)
	} else {
		k8s.AddCondition(mgtCtx, k8s.NewAcceptedConditionBuilder(mgtCtx.GetGeneration()).
			Accept("Successfully reconciled").Build())
	}

	return k8s.GetClient().Status().Update(ctx, mgtCtx)
}
