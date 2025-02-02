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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
<<<<<<< HEAD
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
=======
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
>>>>>>> 539e666 (fix: remove secret controller)
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	"golang.org/x/net/context"
	v1 "k8s.io/api/core/v1"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func Delete(
	ctx context.Context,
	instance *v1alpha1.ManagementContext,
) error {
	if !util.ContainsFinalizer(instance, keys.ManagementContextFinalizer) {
		return nil
	}

<<<<<<< HEAD
	apis := &v1alpha1.ApiDefinitionList{}
	if err := search.FindByFieldReferencing(
		ctx,
		indexer.ApiContextField,
		refs.NewNamespacedName(instance.Namespace, instance.Name),
		apis,
	); err != nil {
		err = fmt.Errorf("an error occurred while checking if the management context is linked to an api definition: %w", err)
		return err
	}

	if len(apis.Items) > 0 {
		return fmt.Errorf("can not delete %s because %d api(s) relying on this context", instance.Name, len(apis.Items))
	}

	apisV4 := &v1alpha1.ApiV4DefinitionList{}
	if err := search.FindByFieldReferencing(
		ctx,
		indexer.ApiV4ContextField,
		refs.NewNamespacedName(instance.Namespace, instance.Name),
		apisV4,
	); err != nil {
		err = fmt.Errorf("can not check if the management context is linked to an api v4 definition: %w", err)
		return err
	}

	if len(apisV4.Items) > 0 {
		return fmt.Errorf("can not delete %s because %d api(s) relying on this context", instance.Name, len(apisV4.Items))
	}

	apps := &v1alpha1.ApplicationList{}
	err := search.FindByFieldReferencing(
		ctx,
		indexer.AppContextField,
		refs.NewNamespacedName(instance.Namespace, instance.Name),
		apps,
	)

	if err != nil {
		err = fmt.Errorf("an error occurred while checking if the management context is linked to an application: %w", err)
		return err
	}

	if len(apps.Items) > 0 {
		return fmt.Errorf("can not delete %s because %d application(s) are relying on this context",
			instance.Name, len(apps.Items))
	}

	util.RemoveFinalizer(instance, keys.ManagementContextFinalizer)

	return nil
=======
	if instance.HasSecretRef() {
		secret := &v1.Secret{}

		nsn := getSecretRef(instance)
		if err := k8s.GetClient().Get(ctx, nsn, secret); err != nil {
			return err
		}

		isRef, err := isReferenced(ctx, *instance.Spec.Auth.SecretRef)
		if err != nil {
			return err
		}

		if !isRef {
			util.RemoveFinalizer(secret, core.ManagementContextSecretFinalizer)
		}

		if err := k8s.GetClient().Update(ctx, secret); err != nil {
			return err
		}
	}

	return search.AssertNoContextRef(ctx, instance)
>>>>>>> 539e666 (fix: remove secret controller)
}
