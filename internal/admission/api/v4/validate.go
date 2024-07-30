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

package v4

import (
	"context"
	"net/url"
	"slices"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/types/k8s/custom"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func validateCreate(ctx context.Context, obj runtime.Object) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	if api, ok := obj.(custom.ApiDefinitionResource); ok {
		errs = errs.MergeWith(base.ValidateCreate(ctx, obj))
		if errs.IsSevere() {
			return errs
		}
		errs.Add(validateNoConflictingPath(ctx, api))
		if errs.IsSevere() {
			return errs
		}
		if api.HasContext() {
			errs = errs.MergeWith(validateDryRun(ctx, api))
		}
	}
	return errs
}

// TODO this should be move to base once implemented for v2
func validateNoConflictingPath(ctx context.Context, api custom.ApiDefinitionResource) *errors.AdmissionError {
	apiPaths, err := api.GetContextPaths()
	if err != nil {
		return errors.NewSevere(err.Error())
	}
	existingPaths, err := getExistingPaths(ctx, api)
	if err != nil {
		return errors.NewSevere(err.Error())
	}
	for _, apiPath := range apiPaths {
		if _, pErr := url.Parse(apiPath); pErr != nil {
			return errors.NewSevere(
				"path [%s] is invalid",
				apiPath,
			)
		}
		if slices.Contains(existingPaths, apiPath) {
			return errors.NewSevere(
				"invalid API context path [%s]. Another API with the same path already exists",
				apiPath,
			)
		}
	}
	return nil
}

func validateDryRun(ctx context.Context, api custom.ApiDefinitionResource) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	cp, _ := api.DeepCopyResource().(custom.ApiDefinitionResource)

	apim, err := apim.FromContextRef(ctx, cp.ContextRef(), cp.GetNamespace())
	if err != nil {
		errs.AddSevere(err.Error())
	}

	cp.PopulateIDs(apim.Context)

	impl, ok := cp.GetDefinition().(*v4.Api)
	if !ok {
		errs.AddSevere("unable to call dry run import because api is not a v4 API")
	}

	status, err := apim.APIs.DryRunImportV4(impl)
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}
	for _, severe := range status.Errors.Severe {
		errs.AddSevere(severe)
	}
	if errs.IsSevere() {
		return errs
	}
	for _, warning := range status.Errors.Warning {
		errs.AddWarning(warning)
	}
	return errs
}

func getExistingPaths(ctx context.Context, api custom.ApiDefinitionResource) ([]string, error) {
	existingPaths := make([]string, 0)
	unstructuredList, err := getListOfExistingApis(ctx, api.GetNamespace())
	if err != nil {
		return existingPaths, err
	}

	for _, item := range unstructuredList.Items {
		converted, cErr := dynamic.Convert(item.Object["spec"], new(v4.Api))
		if cErr != nil {
			return existingPaths, cErr
		}
		convertedPaths, pErr := converted.GetContextPaths()
		if pErr != nil {
			return existingPaths, pErr
		}
		if !isCurrentApi(item, api) {
			existingPaths = append(existingPaths, convertedPaths...)
		}
	}
	return existingPaths, nil
}

func isCurrentApi(item unstructured.Unstructured, api custom.ApiDefinitionResource) bool {
	return api.GetName() == item.Object["metadata"].(map[string]interface{})["name"] &&
		api.GetNamespace() == item.Object["metadata"].(map[string]interface{})["namespace"]
}

func getListOfExistingApis(ctx context.Context, ns string) (*unstructured.UnstructuredList, error) {
	gvr := schema.GroupVersionResource{
		Group:    "gravitee.io",
		Version:  "v1alpha1",
		Resource: "apiv4definitions",
	}
	if !env.Config.CheckApiContextPathConflictInCluster {
		return dynamic.GetClient().
			Resource(gvr).
			Namespace(ns).
			List(ctx, metav1.ListOptions{})
	} else {
		return dynamic.GetClient().
			Resource(gvr).
			List(ctx, metav1.ListOptions{})
	}
}
