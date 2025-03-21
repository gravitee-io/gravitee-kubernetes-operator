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

package search

import (
	"context"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

const kubectlCommand = "-A -o jsonpath='{.items[?(@.spec.contextRef.name==\"%s\")].metadata.name}'"
const reviewMessage = "You can review those by running the following command: "

func AssertNoContextRef(ctx context.Context, mCtx core.ContextObject) error {
	ctxRef := refs.NewNamespacedName(mCtx.GetNamespace(), mCtx.GetName())

	apis := &v1alpha1.ApiDefinitionList{}
	if err := FindByFieldReferencing(
		ctx,
		ApiContextField,
		ctxRef,
		apis,
	); err != nil {
		return err
	}
	if len(apis.Items) > 0 {
		return fmt.Errorf(
			"[%s] cannot be deleted because %d APIs are relying on this context. "+
				reviewMessage+
				"kubectl get apidefinitions.gravitee.io "+
				kubectlCommand,
			mCtx.GetName(), len(apis.Items), mCtx.GetName(),
		)
	}

	apisV4 := &v1alpha1.ApiV4DefinitionList{}
	if err := FindByFieldReferencing(
		ctx,
		ApiV4ContextField,
		ctxRef,
		apisV4,
	); err != nil {
		return err
	}

	if len(apisV4.Items) > 0 {
		return fmt.Errorf(
			"[%s] cannot be deleted because %d APIs are relying on this context. "+
				reviewMessage+
				"kubectl get apiv4definitions.gravitee.io "+
				kubectlCommand,
			mCtx.GetName(), len(apisV4.Items), mCtx.GetName(),
		)
	}

	apps := &v1alpha1.ApplicationList{}
	if err := FindByFieldReferencing(
		ctx,
		AppContextField,
		ctxRef,
		apps,
	); err != nil {
		return err
	}

	if len(apps.Items) > 0 {
		return fmt.Errorf(
			"[%s] cannot be deleted because %d applications are relying on this context. "+
				reviewMessage+
				"kubectl get applications.gravitee.io "+
				kubectlCommand,
			mCtx.GetName(), len(apps.Items), mCtx.GetName(),
		)
	}

	spg := &v1alpha1.SharedPolicyGroupList{}
	if err := FindByFieldReferencing(
		ctx,
		SPGContextField,
		ctxRef,
		spg,
	); err != nil {
		return err
	}

	if len(spg.Items) > 0 {
		return fmt.Errorf(
			"[%s] cannot be deleted because %d SharedPolicyGroups are relying on this context. "+
				reviewMessage+
				"kubectl get sharedpolicygroups.gravitee.io "+
				kubectlCommand,
			mCtx.GetName(), len(spg.Items), mCtx.GetName(),
		)
	}

	return nil
}

func AssertNoSharedPolicyGroupRef(ctx context.Context, spg *v1alpha1.SharedPolicyGroup) error {
	nsn := refs.NewNamespacedName(spg.Namespace, spg.Name)

	apisV4 := &v1alpha1.ApiV4DefinitionList{}
	if err := FindByFieldReferencing(
		ctx,
		ApiV4SharedPolicyGroupsField,
		nsn,
		apisV4,
	); err != nil {
		return err
	}

	if len(apisV4.Items) > 0 {
		return fmt.Errorf(
			"[%s] cannot be deleted because %d APIs are relying on this Shared Policy Group. "+
				reviewMessage+
				"kubectl get sharedpolicygroups.gravitee.io "+
				kubectlCommand,
			spg.Name, len(apisV4.Items), spg.Name,
		)
	}

	return nil
}
