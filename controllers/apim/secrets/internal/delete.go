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
	v1 "k8s.io/api/core/v1"
)

func Delete(ctx context.Context, secret *v1.Secret) error {
	return checkContextFinalizer(ctx, secret)
}

func checkContextFinalizer(ctx context.Context, secret *v1.Secret) error {
	contextRefs, err := getReferences(ctx, secret, new(v1alpha1.ManagementContextList))

	if err != nil {
		return err
	}

	refCount := len(contextRefs)

	if refCount >= 1 {
		return fmt.Errorf("secret is used by %d management context, cannot be deleted", refCount)
	}

	return nil
}
