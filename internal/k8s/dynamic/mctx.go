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

package dynamic

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/management"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

func ExpectResolvedContext(ctx context.Context, ref core.ObjectRef, parentNs string) error {
	if _, err := ResolveContext(ctx, ref, parentNs); err != nil {
		return err
	}
	return nil
}

func ResolveContext(ctx context.Context, ref core.ObjectRef, parentNs string) (*management.Context, error) {
	context, err := resolveRefSpec(ctx, ref, parentNs, ManagementContextGVR, new(management.Context))
	if err != nil {
		return nil, err
	}

	return context, err
}

func InjectSecretIfAny(ctx context.Context, mCtx core.ContextModel, parentNs string) (*core.ContextModel, error) {
	if mCtx.HasSecretRef() || (mCtx.HasCloud() && mCtx.GetCloud().HasSecretRef()) { //nolint:nestif // normal complexity
		var name string
		if mCtx.HasSecretRef() {
			name = mCtx.GetSecretRef().GetName()
		} else {
			name = mCtx.GetCloud().GetSecretRef().GetName()
		}

		secret, err := ResolveSecret(ctx, &refs.NamespacedName{Name: name, Namespace: parentNs}, parentNs)
		if err != nil {
			return nil, err
		}
		var bearerToken string
		if mCtx.HasCloud() {
			bearerToken = string(secret.Data[core.CloudTokenSecretKey])
		} else {
			bearerToken = string(secret.Data[core.BearerTokenSecretKey])
		}
		username := string(secret.Data[core.UsernameSecretKey])
		password := string(secret.Data[core.PasswordSecretKey])

		if mCtx.GetAuth() != nil {
			mCtx.GetAuth().SetToken(bearerToken)
			mCtx.GetAuth().SetCredentials(username, password)
		}
	}

	return &mCtx, nil
}
