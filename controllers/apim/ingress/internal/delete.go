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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	v1 "k8s.io/api/networking/v1"
)

func Delete(
	ctx context.Context,
	ingress *v1.Ingress) error {
	if err := deleteIngressTLSReference(ctx, ingress); err != nil {
		log.Error(ctx, err, "An error occurred while updating the TLS secrets", log.KeyValues(ingress)...)
		return err
	}

	return nil
}
