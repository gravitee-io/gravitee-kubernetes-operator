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

	v1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (d *Delegate) CreateOrUpdate(
	ctx context.Context,
	desired *v1.Ingress) error {
	if err := d.updateIngressTLSReference(ctx, desired); err != nil {
		log.FromContext(ctx).Error(err, "An error occurred while updating the PEM registry")
		return err
	}

	operation, apiDefinitionError := d.createOrUpdateApiDefinition(ctx, desired)
	if apiDefinitionError != nil {
		log.FromContext(ctx).Error(
			apiDefinitionError,
			"An error occurs while creating or updating the ApiDefinition",
			"Operation", operation,
		)
		return apiDefinitionError
	}

	return nil
}
