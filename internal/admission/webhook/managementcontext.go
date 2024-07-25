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

package webhook

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/uuid"
)

func CheckAPIMAvailability(ctxRef *refs.NamespacedName) error {
	cli, err := apim.FromContextRef(context.Background(), ctxRef)
	if err != nil {
		return fmt.Errorf("can't create apim client for this management context [%s]", ctxRef.Name)
	}

	_, err = cli.APIs.GetV4ByID(uuid.NewV4String())

	var opError *net.OpError
	if errors.As(err, &opError) {
		return fmt.Errorf("unable to reach APIM, [%s] is not available", cli.Context.BaseUrl)
	}

	return nil
}
