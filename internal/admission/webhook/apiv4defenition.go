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

package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func ValidateApiV4(ctx context.Context, api *v4.Api, name, ns string,
	ctxRef *refs.NamespacedName) (admission.Warnings, error) {
	// make sure Management Context exist before creating the API Definition resource
	if ctxRef != nil { //nolint:nestif // nested if is needed
		gvr := schema.GroupVersionResource{
			Group:    "gravitee.io",
			Version:  "v1alpha1",
			Resource: "apiv4definitions",
		}
		if _, err := k8s.GetDynamicClient().
			Resource(gvr).
			Namespace(ctxRef.Namespace).
			Get(ctx, ctxRef.Name, metav1.GetOptions{}); err != nil {
			return admission.Warnings{}, fmt.Errorf("can't create API [%s] because it is using "+
				"management context [%v] that doesn't exist in the cluster", api.Name, ctxRef)
		}
	} else {
		// check for unique context path
		apis, err := getListOfExistingApis(ctx, ns)

		if err != nil {
			return admission.Warnings{}, fmt.Errorf("can't list existing APIs")
		}

		existingListeners := make([]*v4.GenericListener, 0)
		for _, item := range apis.Items {
			bytes, mErr := json.Marshal(item.Object["spec"])
			if mErr != nil {
				return admission.Warnings{}, mErr
			}

			ea := new(v4.Api)
			err = json.Unmarshal(bytes, ea)
			if err != nil {
				return admission.Warnings{}, err
			}

			if name != item.Object["metadata"].(map[string]interface{})["name"] ||
				ns != item.Object["metadata"].(map[string]interface{})["namespace"] {
				existingListeners = append(existingListeners, ea.Listeners...)
			}
		}

		if err = validateApiContextPath(existingListeners, api.Listeners); err != nil {
			return admission.Warnings{}, err
		}
	}

	return admission.Warnings{}, nil
}

func getListOfExistingApis(ctx context.Context, ns string) (*unstructured.UnstructuredList, error) {
	gvr := schema.GroupVersionResource{
		Group:    "gravitee.io",
		Version:  "v1alpha1",
		Resource: "apiv4definitions",
	}
	if !env.Config.CheckApiContextPathConflictInCluster {
		return k8s.GetDynamicClient().
			Resource(gvr).
			Namespace(ns).
			List(ctx, metav1.ListOptions{})
	} else {
		return k8s.GetDynamicClient().
			Resource(gvr).
			List(ctx, metav1.ListOptions{})
	}
}

func validateApiContextPath(existingListeners, listeners []*v4.GenericListener) error {
	apiPaths := make([]string, 0)
	for _, l := range listeners {
		for _, s := range parseListener(l) {
			p, err := url.Parse(s)
			if err != nil {
				return err
			}
			apiPaths = append(apiPaths, p.String())
		}
	}

	for _, l := range existingListeners {
		paths := parseListener(l)
		err := findDuplicatePath(paths, apiPaths)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseListener(l v4.Listener) []string {
	if l == nil {
		return []string{}
	}

	switch t := l.(type) {
	case *v4.GenericListener:
		return parseListener(t.ToListener())
	case *v4.HttpListener:
		{
			paths := make([]string, 0)
			for _, path := range t.Paths {
				p := fmt.Sprintf("%s/%s", path.Host, path.Path)
				paths = append(paths, strings.ReplaceAll(p, "//", "/"))
			}
			return paths
		}
	case *v4.TCPListener:
		return t.Hosts
	}

	return []string{}
}

func findDuplicatePath(existingPaths []string, newPaths []string) error {
	for _, ep := range existingPaths {
		for _, np := range newPaths {
			if ep == np {
				return fmt.Errorf("invalid API context path [%s]. Another API with the same path already exists", ep)
			}
		}
	}
	return nil
}
