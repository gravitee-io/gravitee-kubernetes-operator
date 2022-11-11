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
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes/scheme"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	uuid "github.com/satori/go.uuid" //nolint:gomodguard // to replace with google implementation
)

var decode = scheme.Codecs.UniversalDecoder().Decode

type ApiWithContext struct {
	Api     *gio.ApiDefinition
	Context *gio.ManagementContext
}

func NewApiDefinition(path string, transforms ...func(*gio.ApiDefinition)) (*gio.ApiDefinition, error) {
	crd, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	gvk := gio.GroupVersion.WithKind("ApiDefinition")
	decoded, _, err := decode(crd, &gvk, new(gio.ApiDefinition))
	if err != nil {
		return nil, err
	}

	api, ok := decoded.(*gio.ApiDefinition)
	if !ok {
		return nil, fmt.Errorf("failed to assert type of API CRD")
	}

	for _, transform := range transforms {
		transform(api)
	}

	addRandomSuffixes(api)

	return api, nil
}

func NewApiWithRandomContext(
	apiPath string, contextPath string, transforms ...func(*ApiWithContext),
) (*ApiWithContext, error) {
	api, err := NewApiDefinition(apiPath)
	if err != nil {
		return nil, err
	}

	ctx, err := newManagementContext(contextPath, func(ctx *gio.ManagementContext) {
		ctx.Name += "-" + uuid.NewV4().String()[:7]
		api.Spec.Context = &model.NamespacedName{
			Name:      ctx.Name,
			Namespace: ctx.Namespace,
		}
	})

	if err != nil {
		return nil, err
	}

	apiWithContext := &ApiWithContext{
		Api:     api,
		Context: ctx,
	}

	for _, transform := range transforms {
		transform(apiWithContext)
	}

	return apiWithContext, nil
}

func newManagementContext(path string, transforms ...func(*gio.ManagementContext)) (*gio.ManagementContext, error) {
	crd, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	gvk := gio.GroupVersion.WithKind("ManagementContext")
	decoded, _, err := decode(crd, &gvk, new(gio.ManagementContext))
	if err != nil {
		return nil, err
	}

	ctx, ok := decoded.(*gio.ManagementContext)
	if !ok {
		return nil, fmt.Errorf("failed to assert type of Management Context CRD")
	}

	for _, transform := range transforms {
		transform(ctx)
	}

	return ctx, nil
}

func addRandomSuffixes(api *gio.ApiDefinition) {
	suffix := "-" + uuid.NewV4().String()[:7]
	api.Name += suffix
	api.Spec.Name += suffix
	api.Spec.Proxy.VirtualHosts[0].Path += suffix
}
