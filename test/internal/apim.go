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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/service"
)

type APIM struct {
	*apim.APIM

	Applications  *service.Applications
	Subscriptions *service.Subscriptions
}

func NewAPIM(ctx context.Context) (*APIM, error) {
	context, err := newManagementContext(ContextWithCredentialsFile)
	if err != nil {
		return nil, err
	}

	apim, err := apim.FromContext(ctx, context.Spec.Context)
	if err != nil {
		return nil, err
	}

	applications := service.NewApplications(apim.APIs.Client)
	subscriptions := service.NewSubscriptions(apim.APIs.Client)

	return &APIM{
		APIM:          apim,
		Applications:  applications,
		Subscriptions: subscriptions,
	}, nil
}
