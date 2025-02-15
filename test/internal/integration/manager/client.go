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

package manager

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Client() client.Client {
	return k8s.GetClient()
}

func GetLatest[T client.Object](ctx context.Context, obj T) error {
	return k8s.GetLatest(ctx, obj)
}

func Delete[T client.Object](ctx context.Context, obj T) error {
	return k8s.Delete(ctx, obj)
}

func UpdateSafely[T client.Object](ctx context.Context, objNew T) error {
	return k8s.UpdateSafely(ctx, objNew)
}
