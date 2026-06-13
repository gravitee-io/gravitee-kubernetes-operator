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

	if err := assertNoApiDefinitions(ctx, ctxRef, mCtx.GetName()); err != nil {
		return err
	}

	if err := assertNoApiV4Definitions(ctx, ctxRef, mCtx.GetName()); err != nil {
		return err
	}

	if err := assertNoApplications(ctx, ctxRef, mCtx.GetName()); err != nil {
		return err
	}

	if err := assertNoSharedPolicyGroups(ctx, ctxRef, mCtx.GetName()); err != nil {
		return err
	}

	if err := assertNoGroups(ctx, ctxRef, mCtx.GetName()); err != nil {
		return err
	}

	if err := assertNoDictionaries(ctx, ctxRef, mCtx.GetName()); err != nil {
		return err
	}

	if err := assertNoPortals(ctx, ctxRef, mCtx.GetName()); err != nil {
		return err
	}

	return nil
}

func assertNoApiDefinitions(ctx context.Context, ctxRef refs.NamespacedName, contextName string) error {
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
			contextName, len(apis.Items), contextName,
		)
	}
	return nil
}

func assertNoApiV4Definitions(ctx context.Context, ctxRef refs.NamespacedName, contextName string) error {
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
			contextName, len(apisV4.Items), contextName,
		)
	}
	return nil
}

func assertNoApplications(ctx context.Context, ctxRef refs.NamespacedName, contextName string) error {
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
			contextName, len(apps.Items), contextName,
		)
	}
	return nil
}

func assertNoSharedPolicyGroups(ctx context.Context, ctxRef refs.NamespacedName, contextName string) error {
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
			contextName, len(spg.Items), contextName,
		)
	}
	return nil
}

func assertNoGroups(ctx context.Context, ctxRef refs.NamespacedName, contextName string) error {
	groups := &v1alpha1.GroupList{}
	if err := FindByFieldReferencing(
		ctx,
		GroupContextField,
		ctxRef,
		groups,
	); err != nil {
		return err
	}

	if len(groups.Items) > 0 {
		return fmt.Errorf(
			"[%s] cannot be deleted because %d groups are relying on this context. "+
				reviewMessage+
				"kubectl get groups.gravitee.io "+
				kubectlCommand,
			contextName, len(groups.Items), contextName,
		)
	}
	return nil
}

func assertNoDictionaries(ctx context.Context, ctxRef refs.NamespacedName, contextName string) error {
	dictionaries := &v1alpha1.DictionaryList{}
	if err := FindByFieldReferencing(
		ctx,
		DictionaryContextField,
		ctxRef,
		dictionaries,
	); err != nil {
		return err
	}

	if len(dictionaries.Items) > 0 {
		return fmt.Errorf(
			"[%s] cannot be deleted because %d dictionaries are relying on this context. "+
				reviewMessage+
				"kubectl get dictionaries.gravitee.io "+
				kubectlCommand,
			contextName, len(dictionaries.Items), contextName,
		)
	}
	return nil
}

func assertNoPortals(ctx context.Context, ctxRef refs.NamespacedName, contextName string) error {
	portals := &v1alpha1.PortalList{}
	if err := FindByFieldReferencing(
		ctx,
		PortalContextField,
		ctxRef,
		portals,
	); err != nil {
		return err
	}

	if len(portals.Items) > 0 {
		return fmt.Errorf(
			"[%s] cannot be deleted because %d portals are relying on this context. "+
				reviewMessage+
				"kubectl get portals.gravitee.io "+
				kubectlCommand,
			contextName, len(portals.Items), contextName,
		)
	}
	return nil
}

func AssertNoPortalListingRef(ctx context.Context, prtl *v1alpha1.Portal) error {
	nsn := refs.NewNamespacedName(prtl.Namespace, prtl.Name)

	listings := &v1alpha1.PortalListingList{}
	if err := FindByFieldReferencing(
		ctx,
		PortalListingPortalField,
		nsn,
		listings,
	); err != nil {
		return err
	}

	if len(listings.Items) > 0 {
		return fmt.Errorf(
			"[%s] cannot be deleted because %d portal listings are relying on this portal. "+
				reviewMessage+
				"kubectl get portallistings.gravitee.io "+
				"-A -o jsonpath='{.items[?(@.spec.portalRef.name==\"%s\")].metadata.name}'",
			prtl.Name, len(listings.Items), prtl.Name,
		)
	}

	return nil
}

// AssertNoPortalDocumentationRef blocks deletion of a Portal while any
// Documentation page is still attached to it.
func AssertNoPortalDocumentationRef(ctx context.Context, prtl *v1alpha1.Portal) error {
	nsn := refs.NewNamespacedName(prtl.Namespace, prtl.Name)
	return assertNoDocumentations(ctx, DocumentationPortalField, nsn, "portal", prtl.Name, "portalRef")
}

// AssertNoApiDocumentationRef blocks deletion of an API while any Documentation
// page is still attached to it.
func AssertNoApiDocumentationRef(ctx context.Context, api core.ApiDefinitionObject) error {
	nsn := refs.NewNamespacedName(api.GetNamespace(), api.GetName())
	return assertNoDocumentations(ctx, DocumentationApiField, nsn, "API", api.GetName(), "apiRef")
}

func assertNoDocumentations(
	ctx context.Context,
	field IndexField,
	nsn refs.NamespacedName,
	ownerKind, ownerName, jsonPathRef string,
) error {
	docs := &v1alpha1.DocumentationList{}
	if err := FindByFieldReferencing(
		ctx,
		field,
		nsn,
		docs,
	); err != nil {
		return err
	}

	if len(docs.Items) > 0 {
		return fmt.Errorf(
			"[%s] cannot be deleted because %d documentation pages are relying on this %s. "+
				reviewMessage+
				"kubectl get documentations.gravitee.io "+
				"-A -o jsonpath='{.items[?(@.spec.%s.name==\"%s\")].metadata.name}'",
			ownerName, len(docs.Items), ownerKind, jsonPathRef, ownerName,
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
