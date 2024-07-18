/*
 * Copyright (C) 2015 The Gravitee team (http://gravitee.io)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	runtime "k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ webhook.Defaulter = &ApiV4Definition{}
var _ webhook.Validator = &ApiV4Definition{}

func (api *ApiV4Definition) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(api).
		Complete()
}

func (api *ApiV4Definition) Default() {}

func (api *ApiV4Definition) ValidateCreate() (admission.Warnings, error) {
	return validateApi(api)
}

func (api *ApiV4Definition) ValidateUpdate(_ runtime.Object) (admission.Warnings, error) {
	return validateApi(api)
}

func (*ApiV4Definition) ValidateDelete() (admission.Warnings, error) {
	return admission.Warnings{}, nil
}

func validateApi(api *ApiV4Definition) (admission.Warnings, error) {
	// make sure Management Context exist before creating the API Definition resource
	if api.HasContext() {
		mCtx := new(ManagementContext)
		if err := k8s.GetClient().Get(context.Background(), api.ContextRef().NamespacedName(), mCtx); err != nil {
			return admission.Warnings{}, fmt.Errorf("can't create api [%s] because it is using "+
				"management context [%v] that doesn't exist in the cluster", api.Name, api.ContextRef().NamespacedName())
		}
	} else {
		// check for unique context path
		apis := new(ApiV4DefinitionList)
		if err := k8s.GetClient().List(context.Background(), apis); err != nil {
			return admission.Warnings{}, err
		}

		existingListeners := make([]*v4.GenericListener, 0)
		for _, item := range apis.Items {
			existingListeners = append(existingListeners, item.Spec.Listeners...)
		}

		if err := validateApiContextPath(existingListeners, api.Spec.Listeners); err != nil {
			return admission.Warnings{}, err
		}
	}

	return admission.Warnings{}, nil
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
				return fmt.Errorf("invalid api context path. the same path already exist [%s]", ep)
			}
		}
	}
	return nil
}
