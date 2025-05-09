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

	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

func Delete(
	ctx context.Context,
	resource *v1alpha1.Notification,
) error {
	if !util.ContainsFinalizer(resource, core.NotificationFinalizer) {
		return nil
	}

	apis := &v1alpha1.ApiDefinitionList{}
	err := checkRefs(ctx, resource, apis)
	if err != nil {
		return err
	}

	if len(apis.Items) > 0 {
		return fmt.Errorf("resource is referenced and will remain")
	}


	apisV4 := &v1alpha1.ApiV4DefinitionList{}
	err = checkRefs(ctx, resource, apisV4)
	if err != nil {
		return err
	}

	if len(apis.Items) > 0 {
		return fmt.Errorf("resource is referenced and will remain")
	}

	return nil
}

func checkRefs(ctx context.Context, resource *v1alpha1.Notification, apis client.ObjectList) error{
	//if err := search.FindByFieldReferencing(
	//	ctx,
	//	search.NotificationRefsField,
	//	refs.NewNamespacedName(resource.Namespace, resource.Name),
	//	apis,
	//); err != nil {
	//	err = fmt.Errorf("an error occurred while checking if the notification is linked to an api definition: %w", err)
	//	return err
	//}

	return nil
}
