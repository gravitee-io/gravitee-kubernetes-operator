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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func Delete(ctx context.Context, api core.ApiDefinitionObject) error {
	if !util.ContainsFinalizer(api, core.ApiDefinitionFinalizer) {
		return nil
	}

	if api.HasContext() {
		if err := deleteWithContext(ctx, api); err != nil {
			return err
		}
	}

	util.RemoveFinalizer(api, core.ApiDefinitionFinalizer)

	return nil
}

func deleteWithContext(ctx context.Context, api core.ApiDefinitionObject) error {
	apim, err := apim.FromContextRef(ctx, api.ContextRef(), api.GetNamespace())
	if err != nil {
		return err
	}
	switch {
	case api.GetDefinitionVersion() == core.ApiV2:
		return errors.IgnoreNotFound(apim.APIs.DeleteV2(api.GetID()))
	case api.GetDefinitionVersion() == core.ApiV4:
		return errors.IgnoreNotFound(apim.APIs.DeleteV4(api.GetID()))
	default:
		return fmt.Errorf("unknown version %s", api.GetDefinitionVersion())
	}
}
