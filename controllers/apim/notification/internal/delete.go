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
	"fmt"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"

	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

func Delete(
	ctx context.Context,
	notification *v1alpha1.Notification,
) error {
	if !util.ContainsFinalizer(notification, core.NotificationFinalizer) {
		return nil
	}

	if err := checkAPIRefs(ctx, notification); err != nil {
		return err
	}

	if err := checkAPIv4Refs(ctx, notification); err != nil {
		return err
	}

	return nil
}

func checkAPIRefs(ctx context.Context, notification *v1alpha1.Notification) error {
	apis := &v1alpha1.ApiDefinitionList{}
	err := checkRefs(ctx, notification, search.ApiNotificationRefsField, apis)
	if err != nil {
		return err
	}

	if len(apis.Items) > 0 {
		return formatApiError(apis.Items)
	}
	return nil
}

func checkAPIv4Refs(ctx context.Context, notification *v1alpha1.Notification) error {
	apis := &v1alpha1.ApiV4DefinitionList{}
	err := checkRefs(ctx, notification, search.ApiV4NotificationRefsField, apis)
	if err != nil {
		return err
	}

	if len(apis.Items) > 0 {
		return formatApV4iError(apis.Items)
	}
	return nil
}

func formatApiError(apis []v1alpha1.ApiDefinition) error {
	apiRefs := make([]string, len(apis))
	for _, api := range apis {
		apiRefs = append(apiRefs, api.GetRef().NamespacedName().String())
	}
	return fmt.Errorf("notification is referenced by APIs: %v and will remain until those are removed", apiRefs)
}
func formatApV4iError(apis []v1alpha1.ApiV4Definition) error {
	apiRefs := make([]string, len(apis))
	for _, api := range apis {
		apiRefs = append(apiRefs, api.GetRef().NamespacedName().String())
	}
	return fmt.Errorf("notification is referenced by V4 APIs: %v and will remain until those are removed", apiRefs)
}

func checkRefs(
	ctx context.Context,
	resource *v1alpha1.Notification,
	field search.IndexField,
	apis client.ObjectList) error {
	if err := search.FindByFieldReferencing(
		ctx,
		field,
		refs.NewNamespacedName(resource.Namespace, resource.Name),
		apis,
	); err != nil {
		err = fmt.Errorf("an error occurred while checking if the notification is linked to an api definition: %w", err)
		return err
	}

	return nil
}
