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

type Fixtures struct {
	Api      *gio.ApiDefinition
	Contexts []gio.ApiContext
	Resource *gio.ApiResource
}

type FixtureFiles struct {
	Api      string
	Contexts []string
	Resource string
}

type FixtureGenerator struct {
	suffix string
}

func NewFixtureGenerator() *FixtureGenerator {
	return &FixtureGenerator{
		suffix: randomSuffix(),
	}
}

func (f *FixtureGenerator) NewFixtures(files FixtureFiles, transforms ...func(*Fixtures)) (*Fixtures, error) {
	fixtures := &Fixtures{}

	if files.Api != "" {
		api, err := f.NewApiDefinition(files.Api)
		if err != nil {
			return nil, err
		}
		fixtures.Api = api
	}

	if files.Contexts != nil {
		fixtures.Contexts = make([]gio.ApiContext, 0)

		for _, file := range files.Contexts {
			ctx, err := f.NewApiContext(file)
			if err != nil {
				return nil, err
			}
			fixtures.Contexts = append(fixtures.Contexts, *ctx)
		}
	}

	if files.Resource != "" {
		resource, err := f.NewApiResource(files.Resource)
		if err != nil {
			return nil, err
		}
		fixtures.Resource = resource
	}

	if fixtures.Contexts != nil {
		apiContexts := fixtures.Api.Spec.Contexts
		for _, ctx := range fixtures.Contexts {
			apiContexts = append(apiContexts, ctx.GetNamespacedName())
		}
		fixtures.Api.Spec.Contexts = apiContexts
	}

	if fixtures.Resource != nil {
		fixtures.Api.Spec.Resources = []*model.ResourceOrRef{
			{
				Ref: &model.NamespacedName{
					Name:      fixtures.Resource.Name,
					Namespace: fixtures.Resource.Namespace,
				},
			},
		}
	}

	for _, transform := range transforms {
		transform(fixtures)
	}

	return fixtures, nil
}

func (f *FixtureGenerator) NewApiDefinition(
	path string, transforms ...func(*gio.ApiDefinition),
) (*gio.ApiDefinition, error) {
	api, err := newApiDefinition(path, transforms...)
	if err != nil {
		return nil, err
	}

	api.Name += f.suffix
	api.Spec.Name += f.suffix
	api.Spec.Proxy.VirtualHosts[0].Path += f.suffix

	return api, nil
}

func (f *FixtureGenerator) NewApiContext(
	path string, transforms ...func(*gio.ApiContext),
) (*gio.ApiContext, error) {
	ctx, err := newApiContext(path, transforms...)
	if err != nil {
		return nil, err
	}

	ctx.Name += f.suffix

	return ctx, nil
}

func (f *FixtureGenerator) NewApiResource(path string, transforms ...func(*gio.ApiResource)) (*gio.ApiResource, error) {
	resource, err := newApiResource(path, transforms...)
	if err != nil {
		return nil, err
	}
	resource.Name += f.suffix

	return resource, nil
}

func newApiDefinition(path string, transforms ...func(*gio.ApiDefinition)) (*gio.ApiDefinition, error) {
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

	return api, nil
}

func newApiResource(path string, transforms ...func(*gio.ApiResource)) (*gio.ApiResource, error) {
	crd, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	gvk := gio.GroupVersion.WithKind("ApiResource")
	decoded, _, err := decode(crd, &gvk, new(gio.ApiResource))
	if err != nil {
		return nil, err
	}

	resource, ok := decoded.(*gio.ApiResource)
	if !ok {
		return nil, fmt.Errorf("failed to assert type of API CRD")
	}

	for _, transform := range transforms {
		transform(resource)
	}

	return resource, nil
}

func newApiContext(path string, transforms ...func(*gio.ApiContext)) (*gio.ApiContext, error) {
	crd, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	gvk := gio.GroupVersion.WithKind("ApiContext")
	decoded, _, err := decode(crd, &gvk, new(gio.ApiContext))
	if err != nil {
		return nil, err
	}

	ctx, ok := decoded.(*gio.ApiContext)
	if !ok {
		return nil, fmt.Errorf("failed to assert type of API Context CRD")
	}

	for _, transform := range transforms {
		transform(ctx)
	}

	return ctx, nil
}

func randomSuffix() string {
	return "-" + uuid.NewV4().String()[:7]
}
