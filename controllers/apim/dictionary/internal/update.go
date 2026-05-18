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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func CreateOrUpdate(ctx context.Context, dict *v1alpha1.Dictionary) error {
	ns := dict.Namespace

	apimClient, err := apim.FromContextRef(ctx, dict.ContextRef(), ns)
	if err != nil {
		return err
	}
	status, err := apimClient.Dictionaries.CreateOrUpdate(dict)
	if err != nil {
		return gerrors.NewControlPlaneError(err)
	}

	// Setting fields by fields to keep the rest intact
	dict.Status.ID = status.ID
	dict.Status.OrgID = status.OrgID
	dict.Status.EnvID = status.EnvID

	return nil
}
