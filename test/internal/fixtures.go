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

	netV1 "k8s.io/api/networking/v1"

	"k8s.io/client-go/kubernetes/scheme"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	uuid "github.com/satori/go.uuid" //nolint:gomodguard // to replace with google implementation
)

var decode = scheme.Codecs.UniversalDecoder().Decode

type Fixtures struct {
	Api      *gio.ApiDefinition
	Context  *gio.ManagementContext
	Resource *gio.ApiResource
	Ingress  *netV1.Ingress
}

type FixtureFiles struct {
	Api      string
	Context  string
	Resource string
	Ingress  string
}

type FixtureGenerator struct {
	Suffix string
}

func NewFixtureGenerator() *FixtureGenerator {
	return &FixtureGenerator{
		Suffix: randomSuffix(),
	}
}

func (f *FixtureGenerator) AddSuffix(property string) string {
	return property + f.Suffix
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

	if files.Context != "" {
		ctx, err := f.NewManagementContext(files.Context)
		if err != nil {
			return nil, err
		}
		fixtures.Context = ctx
	}

	if files.Resource != "" {
		resource, err := f.NewApiResource(files.Resource)
		if err != nil {
			return nil, err
		}
		fixtures.Resource = resource
	}

	if fixtures.Context != nil {
		fixtures.Api.Spec.Context = fixtures.Context.GetNamespacedName()
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

	if files.Ingress != "" {
		ingress, err := f.NewIngress(files.Ingress)
		if err != nil {
			return nil, err
		}
		fixtures.Ingress = ingress
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

	api.Name += f.Suffix
	api.Spec.Name += f.Suffix
	api.Spec.Proxy.VirtualHosts[0].Path += f.Suffix

	return api, nil
}

func (f *FixtureGenerator) NewManagementContext(
	path string, transforms ...func(*gio.ManagementContext),
) (*gio.ManagementContext, error) {
	ctx, err := newManagementContext(path, transforms...)
	if err != nil {
		return nil, err
	}

	ctx.Name += f.Suffix

	return ctx, nil
}

func (f *FixtureGenerator) NewApiResource(path string, transforms ...func(*gio.ApiResource)) (*gio.ApiResource, error) {
	resource, err := newApiResource(path, transforms...)
	if err != nil {
		return nil, err
	}
	resource.Name += f.Suffix

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
		return nil, fmt.Errorf("failed to assert type of API Context CRD")
	}

	for _, transform := range transforms {
		transform(ctx)
	}

	return ctx, nil
}

func (f *FixtureGenerator) NewIngress(path string, transforms ...func(*netV1.Ingress)) (*netV1.Ingress, error) {
	ingress, err := newIngress(path, transforms...)
	if err != nil {
		return nil, err
	}
	ingress.Name += f.Suffix

	return ingress, nil
}

func newIngress(path string, transforms ...func(*netV1.Ingress)) (*netV1.Ingress, error) {
	crd, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	gvk := gio.GroupVersion.WithKind("Ingress")
	decoded, _, err := decode(crd, &gvk, new(netV1.Ingress))
	if err != nil {
		return nil, err
	}

	resource, ok := decoded.(*netV1.Ingress)
	if !ok {
		return nil, fmt.Errorf("failed to assert type of API CRD")
	}

	for _, transform := range transforms {
		transform(resource)
	}

	return resource, nil
}

func randomSuffix() string {
	return "-" + uuid.NewV4().String()[:7]
}
