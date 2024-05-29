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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env/template"
	netv1 "k8s.io/api/networking/v1"

	"github.com/go-logr/logr"
)

type Delegate struct {
}

func NewDelegate(log logr.Logger) *Delegate {
	return &Delegate{}
}

func (d *Delegate) ResolveTemplate(ctx context.Context, ingress *netv1.Ingress) error {
	return template.NewResolver(ctx, ingress).Resolve()
}
