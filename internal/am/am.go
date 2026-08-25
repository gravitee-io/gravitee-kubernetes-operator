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

package am

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/service"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

const automationBasePath = "/automation"

// AM wraps the Automation API services needed to talk to one AM environment.
type AM struct {
	Domains *service.Domains
	Context core.ContextModel
}

func FromContext(ctx context.Context, amCtx core.ContextObject, parentNs string) (*AM, error) {
	urls, err := client.NewURLs(
		amCtx.GetURL(),
		getAutomationPath(amCtx),
		amCtx.GetOrgID(),
		amCtx.GetEnvID(),
	)
	if err != nil {
		return nil, err
	}

	if _, err = dynamic.InjectSecretIfAny(ctx, amCtx); err != nil {
		return nil, err
	}

	httpClient, err := http.NewClient(ctx, toHttpAuth(amCtx))
	if err != nil {
		return nil, err
	}

	c := &client.Client{
		HTTP: httpClient,
		URLs: urls,
	}

	return &AM{
		Domains: service.NewDomains(c),
		Context: amCtx,
	}, nil
}

func FromContextRef(ctx context.Context, ref core.ObjectRef, parentNs string) (*AM, error) {
	amCtx, err := dynamic.ResolveAMContext(ctx, ref, parentNs)
	if err != nil {
		return nil, err
	}
	return FromContext(ctx, amCtx, parentNs)
}

func getAutomationPath(amCtx core.ContextModel) string {
	if amCtx.GetPath() != nil {
		return *amCtx.GetPath()
	}
	return automationBasePath
}
