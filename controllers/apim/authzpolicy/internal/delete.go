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

package internal

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Delete withdraws the policy from its PDP gateways and deletes it by hrid: the targets do not need
// resolving. Withdrawing a deployed policy needs the ENVIRONMENT_AUTHZ_PDP UPDATE permission: a
// 403 is a control plane error, retried with backoff, so the finalizer stays and Accepted names the
// reason until the context's token is granted it.
func Delete(ctx context.Context, policy *v1alpha1.AuthzPolicy) error {
	if !util.ContainsFinalizer(policy, core.AuthzPolicyFinalizer) {
		return nil
	}

	apimClient, err := apim.FromContextRef(ctx, policy.ContextRef(), policy.GetNamespace())
	if err != nil {
		// ManagementContext already gone: APIM is unreachable, let the finalizer be removed.
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	hrid := refs.NewNamespacedNameFromObject(policy).HRID()
	if err := gerrors.IgnoreNotFound(apimClient.AuthzPolicies.Delete(hrid)); err != nil {
		return gerrors.NewControlPlaneError(err)
	}

	return nil
}
