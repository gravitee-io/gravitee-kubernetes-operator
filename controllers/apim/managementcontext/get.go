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

package managementcontext

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Get(
	ctx context.Context,
	k8sClient client.Client,
	log logr.Logger,
	contextRef *model.ContextRef,
) (*gio.ManagementContext, error) {
	apimContext := new(gio.ManagementContext)
	ns := types.NamespacedName{Name: contextRef.Name, Namespace: contextRef.Namespace}

	log.Info("Looking for context from", "namespace", contextRef.Namespace, "name", contextRef.Name)

	err := k8sClient.Get(ctx, ns, apimContext)

	if err != nil {
		return nil, err
	}

	return apimContext, nil
}
