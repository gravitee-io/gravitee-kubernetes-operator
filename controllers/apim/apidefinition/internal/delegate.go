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
	"net/http"
	"time"

	"github.com/go-logr/logr"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	managementapi "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi"
	k8s "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	requestTimeoutSeconds = 5
	separator             = "/"
	defaultPlanSecurity   = "KEY_LESS"
	defaultPlanStatus     = "PUBLISHED"
	defaultPlanName       = "G.K.O. Default"
	origin                = "kubernetes"
	mode                  = "fully_managed"
)

type Delegate struct {
	ctx               context.Context
	managementContext *gio.ManagementContext
	apimClient        *managementapi.Client
	k8sClient         k8s.Client
	log               logr.Logger
}

func NewDelegate(ctx context.Context, client k8s.Client, log logr.Logger) *Delegate {
	return &Delegate{
		ctx, nil, nil, client, log,
	}
}

func (d *Delegate) SetManagementContext(managementContext *gio.ManagementContext) {
	if managementContext == nil {
		return
	}

	d.managementContext = managementContext

	httpClient := http.Client{Timeout: requestTimeoutSeconds * time.Second}
	d.apimClient = managementapi.NewClient(d.ctx, d.managementContext, httpClient)
}

func (d *Delegate) IsConnectedToManagementApi() bool {
	return d.apimClient != nil
}
