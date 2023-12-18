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

	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"

	netV1 "k8s.io/api/networking/v1"

	"k8s.io/client-go/kubernetes/scheme"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	uuid "github.com/satori/go.uuid" //nolint:gomodguard // to replace with google implementation
)

var decode = scheme.Codecs.UniversalDecoder().Decode

type Fixtures struct {
	Api         *v1alpha1.ApiDefinition
	Context     *v1alpha1.ManagementContext
	Resource    *v1alpha1.ApiResource
	Ingress     *netV1.Ingress
	Application *v1alpha1.Application
}

type FixtureFiles struct {
	Api         string
	Context     string
	Resource    string
	Ingress     string
	Application string
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

	if fixtures.Context != nil && fixtures.Api != nil {
		fixtures.Api.Spec.Context = fixtures.Context.GetNamespacedName()
	}

	if fixtures.Resource != nil {
		fixtures.Api.Spec.Resources = []*base.ResourceOrRef{
			{
				Ref: &refs.NamespacedName{
					Name:      fixtures.Resource.Name,
					Namespace: fixtures.Resource.Namespace,
				},
			},
		}
	}

	if files.Ingress != "" {
		ingress, err := f.NewIngress(files.Ingress, ingressHttpPathTransformer(f))

		if err != nil {
			return nil, err
		}
		fixtures.Ingress = ingress
	}

	err := f.addApplication(files, fixtures)
	if err != nil {
		return nil, err
	}

	for _, transform := range transforms {
		transform(fixtures)
	}

	return fixtures, nil
}

func (f *FixtureGenerator) addApplication(files FixtureFiles, fixtures *Fixtures) error {
	if files.Application != "" {
		application, err := f.NewApplication(files.Application)
		if err != nil {
			return err
		}
		fixtures.Application = application
	}

	if fixtures.Context != nil && fixtures.Application != nil {
		fixtures.Application.Spec.Context = fixtures.Context.GetNamespacedName()
	}

	return nil
}

func ingressHttpPathTransformer(f *FixtureGenerator) func(ingress *netV1.Ingress) {
	return func(ingress *netV1.Ingress) {
		for i := range ingress.Spec.Rules {
			for j := range ingress.Spec.Rules[i].HTTP.Paths {
				ingress.Spec.Rules[i].HTTP.Paths[j].Path += f.Suffix
			}
		}
	}
}

func (f *FixtureGenerator) NewApiDefinition(
	path string, transforms ...func(*v1alpha1.ApiDefinition),
) (*v1alpha1.ApiDefinition, error) {
	api, err := newApiDefinition(path, transforms...)
	if err != nil {
		return nil, err
	}

	api.Name += f.Suffix
	api.Namespace = Namespace
	api.Spec.Name += f.Suffix

	if !isTemplate(api) {
		api.Spec.Proxy.VirtualHosts[0].Path += f.Suffix
	}

	return api, nil
}
func (f *FixtureGenerator) NewManagementContext(
	path string, transforms ...func(*v1alpha1.ManagementContext),
) (*v1alpha1.ManagementContext, error) {
	ctx, err := newManagementContext(path, transforms...)
	if err != nil {
		return nil, err
	}

	ctx.Name += f.Suffix
	ctx.Namespace = Namespace

	return ctx, nil
}

func (f *FixtureGenerator) NewApiResource(
	path string, transforms ...func(*v1alpha1.ApiResource),
) (*v1alpha1.ApiResource, error) {
	resource, err := newApiResource(path, transforms...)
	if err != nil {
		return nil, err
	}
	resource.Name += f.Suffix
	resource.Namespace = Namespace

	return resource, nil
}

func newApiDefinition(path string, transforms ...func(*v1alpha1.ApiDefinition)) (*v1alpha1.ApiDefinition, error) {
	crd, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	gvk := v1alpha1.GroupVersion.WithKind("ApiDefinition")
	decoded, _, err := decode(crd, &gvk, new(v1alpha1.ApiDefinition))
	if err != nil {
		return nil, err
	}

	api, ok := decoded.(*v1alpha1.ApiDefinition)
	if !ok {
		return nil, fmt.Errorf("failed to assert type of API CRD")
	}

	for _, transform := range transforms {
		transform(api)
	}

	return api, nil
}

func newApiResource(path string, transforms ...func(*v1alpha1.ApiResource)) (*v1alpha1.ApiResource, error) {
	crd, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	gvk := v1alpha1.GroupVersion.WithKind("ApiResource")
	decoded, _, err := decode(crd, &gvk, new(v1alpha1.ApiResource))
	if err != nil {
		return nil, err
	}

	resource, ok := decoded.(*v1alpha1.ApiResource)
	if !ok {
		return nil, fmt.Errorf("failed to assert type of API Resource CRD")
	}

	for _, transform := range transforms {
		transform(resource)
	}

	return resource, nil
}

func newManagementContext(
	path string, transforms ...func(*v1alpha1.ManagementContext),
) (*v1alpha1.ManagementContext, error) {
	crd, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	gvk := v1alpha1.GroupVersion.WithKind("ManagementContext")
	decoded, _, err := decode(crd, &gvk, new(v1alpha1.ManagementContext))
	if err != nil {
		return nil, err
	}

	ctx, ok := decoded.(*v1alpha1.ManagementContext)
	if !ok {
		return nil, fmt.Errorf("failed to assert type of ManagementContext CRD")
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
	ingress.Namespace = Namespace

	return ingress, nil
}

func newIngress(path string, transforms ...func(*netV1.Ingress)) (*netV1.Ingress, error) {
	crd, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	gvk := v1alpha1.GroupVersion.WithKind("Ingress")
	decoded, _, err := decode(crd, &gvk, new(netV1.Ingress))
	if err != nil {
		return nil, err
	}

	resource, ok := decoded.(*netV1.Ingress)
	if !ok {
		return nil, fmt.Errorf("failed to assert type of Ingress CRD")
	}

	for _, transform := range transforms {
		transform(resource)
	}

	return resource, nil
}

func (f *FixtureGenerator) NewApplication(path string,
	transforms ...func(application *v1alpha1.Application)) (*v1alpha1.Application, error) {
	application, err := newApplication(path, transforms...)
	if err != nil {
		return nil, err
	}

	application.Name += f.Suffix
	application.Namespace = Namespace

	return application, nil
}

func newApplication(path string, transforms ...func(application *v1alpha1.Application)) (*v1alpha1.Application, error) {
	crd, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	gvk := v1alpha1.GroupVersion.WithKind("Application")
	decoded, _, err := decode(crd, &gvk, new(v1alpha1.Application))
	if err != nil {
		return nil, err
	}

	application, ok := decoded.(*v1alpha1.Application)
	if !ok {
		return nil, fmt.Errorf("failed to assert type of Application CRD")
	}

	for _, transform := range transforms {
		transform(application)
	}

	return application, nil
}

func randomSuffix() string {
	return "-" + uuid.NewV4().String()[:7]
}

func isTemplate(api *v1alpha1.ApiDefinition) bool {
	return api.Annotations[keys.IngressTemplateAnnotation] == "true"
}
