// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dynamic

import (
	"context"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

// AssertCatalogMcpServersSynced returns an error naming the first CatalogMcpServer a studio
// references that does not exist or that the platform has not registered yet. The platform
// refuses a studio selecting such a server, so callers wait for it instead of calling the
// platform: the error is a state to wait out, never a terminal one.
func AssertCatalogMcpServersSynced(ctx context.Context, proxy *v1alpha1.McpProxy) error {
	for _, ref := range proxy.Spec.ServerRefs(proxy.GetNamespace()) {
		srv := &v1alpha1.CatalogMcpServer{}
		key := types.NamespacedName{Namespace: ref.Namespace, Name: ref.Name}
		if err := k8s.GetClient().Get(ctx, key, srv); err != nil {
			if apierrors.IsNotFound(err) {
				return fmt.Errorf("catalog mcp server [%s] not found", ref.String())
			}
			return fmt.Errorf("unable to read catalog mcp server [%s]: %w", ref.String(), err)
		}
		if srv.Status.ID == "" {
			return fmt.Errorf("catalog mcp server [%s] is not synced yet", ref.String())
		}
	}
	return nil
}
