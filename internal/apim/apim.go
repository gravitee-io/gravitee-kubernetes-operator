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

package apim

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/service"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

// APIM wraps services needed to sync resources with a given environment on a Gravitee.io APIM instance.
type APIM struct {
	APIs         *service.APIs
	Applications *service.Applications
	Env          *service.Env

	Context core.ContextModel
}

// EnvID returns the environment ID of the current managed APIM instance.
func (apim *APIM) EnvID() string {
	return apim.Context.GetEnvID()
}

// OrgID returns the organization ID of the current managed APIM instance.
func (apim *APIM) OrgID() string {
	return apim.Context.GetOrgID()
}

// FromContext returns a new APIM instance from a given reconcile context and management context.
func FromContext(ctx context.Context, context core.ContextModel, parentNs string) (*APIM, error) {
	orgID, envID := context.GetOrgID(), context.GetEnvID()
	urls, err := client.NewURLs(context.GetURL(), orgID, envID)
	if err != nil {
		return nil, err
	}

	client := &client.Client{
		HTTP: http.NewClient(ctx, toHttpAuth(context)),
		URLs: urls,
	}

	return &APIM{
		APIs:         service.NewAPIs(client),
		Applications: service.NewApplications(client),
		Env:          service.NewEnv(client),
		Context:      context,
	}, nil
}

func FromContextRef(ctx context.Context, ref core.ObjectRef, parentNs string) (*APIM, error) {
	context, err := dynamic.ResolveContext(ctx, ref, parentNs)
	if err != nil {
		return nil, err
	}
	return FromContext(ctx, context, parentNs)
}
