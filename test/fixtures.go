package test

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes/scheme"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	uuid "github.com/satori/go.uuid" //nolint:gomodguard // to replace with google implementation
)

var decode = scheme.Codecs.UniversalDecoder().Decode

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

func NewManagementContext(path string, transforms ...func(*gio.ManagementContext)) (*gio.ManagementContext, error) {
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
