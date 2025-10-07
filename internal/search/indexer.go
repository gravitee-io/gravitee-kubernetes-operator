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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type IndexField string

const (
	ApiContextField              IndexField = "context"
	ApiV4ContextField            IndexField = "api-v4-context"
	SecretRefField               IndexField = "secretRef"
	ApiResourceField             IndexField = "resource"
	ApiNotificationRefsField     IndexField = "api-notificationRefs"
	ApiGroupField                IndexField = "api-group"
	ApiV4NotificationRefsField   IndexField = "api-v4-notificationRefs"
	ApiV4GroupField              IndexField = "api-v4-group"
	ApiV4ResourceField           IndexField = "api-v4-resource"
	ApiV4SharedPolicyGroupsField IndexField = "api-v4-spg"
	ApiTemplateField             IndexField = "api-template"
	TLSSecretField               IndexField = "tls-secret"
	AppContextField              IndexField = "app-context"
	ApiV2SubsField               IndexField = "api-v2-subscription"
	ApiV4SubsField               IndexField = "api-v4-subscription"
	SPGContextField              IndexField = "spg-context"
)

func (f IndexField) String() string {
	return string(f)
}

type Indexer struct {
	Field string
	Func  client.IndexerFunc
}

func InitCache(ctx context.Context, cache cache.Cache) error {
	errs := make([]error, 0)

	err := newIndexer(ctx, cache, &v1alpha1.ApiDefinition{}, ApiContextField, indexManagementContexts)
	if err != nil {
		errs = append(errs, err)
	}

	err = newIndexer(ctx, cache, &v1alpha1.ApiV4Definition{}, ApiV4ContextField, indexApiV4ManagementContexts)
	if err != nil {
		errs = append(errs, err)
	}

	err = newIndexer(ctx, cache, &v1alpha1.ApiDefinition{}, ApiResourceField, indexApiResourceRefs)
	if err != nil {
		errs = append(errs, err)
	}

	err = newIndexer(ctx, cache, &v1alpha1.ApiV4Definition{}, ApiV4ResourceField, indexApiV4ResourceRefs)
	if err != nil {
		errs = append(errs, err)
	}

	err = newIndexer(ctx, cache, &v1alpha1.ApiDefinition{}, ApiNotificationRefsField, indexNotificationRefs)
	if err != nil {
		errs = append(errs, err)
	}

	err = newIndexer(ctx, cache, &v1alpha1.ApiV4Definition{}, ApiV4NotificationRefsField, indexApiV4NotificationRefs)
	if err != nil {
		errs = append(errs, err)
	}

	err = newIndexer(ctx, cache, &v1alpha1.ApiDefinition{}, ApiGroupField, indexApiGroupsRefs)
	if err != nil {
		errs = append(errs, err)
	}

	err = newIndexer(ctx, cache, &v1alpha1.ApiV4Definition{}, ApiV4GroupField, indexApiV4GroupsRefs)
	if err != nil {
		errs = append(errs, err)
	}

	err = newIndexer(ctx, cache, &v1alpha1.ApiV4Definition{}, ApiV4SharedPolicyGroupsField,
		indexApiV4FlowsSharedPolicyGroupsRefs)
	if err != nil {
		errs = append(errs, err)
	}

	err = newIndexer(ctx, cache, &v1alpha1.ManagementContext{}, SecretRefField, indexManagementContextSecrets)
	if err != nil {
		errs = append(errs, err)
	}

	err = newIndexer(ctx, cache, &v1.Ingress{}, ApiTemplateField, indexApiTemplate)
	if err != nil {
		errs = append(errs, err)
	}

	err = newIndexer(ctx, cache, &v1.Ingress{}, TLSSecretField, indexTLSSecret)
	if err != nil {
		errs = append(errs, err)
	}

	err = newIndexer(ctx, cache, &v1alpha1.Application{}, AppContextField, indexApplicationManagementContexts)
	if err != nil {
		errs = append(errs, err)
	}

	err = newIndexer(ctx, cache, &v1alpha1.Subscription{}, ApiV2SubsField, indexAPIv2Subscriptions)
	if err != nil {
		errs = append(errs, err)
	}

	err = newIndexer(ctx, cache, &v1alpha1.Subscription{}, ApiV4SubsField, indexAPIv4Subscriptions)
	if err != nil {
		errs = append(errs, err)
	}

	err = newIndexer(ctx, cache, &v1alpha1.SharedPolicyGroup{}, SPGContextField, indexSharedPolicyGroupManagementContexts)
	if err != nil {
		errs = append(errs, err)
	}

	return errors.NewAggregate(errs)
}

func createIndexerFunc[T client.Object](doIndex func(T, *[]string)) client.IndexerFunc {
	return func(obj client.Object) []string {
		fields := []string{}
		o, ok := obj.(T)

		if !ok {
			return fields
		}

		doIndex(o, &fields)

		return fields
	}
}

func newIndexer[T client.Object](ctx context.Context, cache cache.Cache, obj client.Object,
	indexField IndexField, doIndex func(T, *[]string)) error {
	indexer := Indexer{
		Field: string(indexField),
		Func:  createIndexerFunc(doIndex),
	}

	if err := cache.IndexField(ctx,
		obj,
		indexer.Field,
		indexer.Func); err != nil {
		return err
	}

	return nil
}

func indexManagementContexts(api *v1alpha1.ApiDefinition, fields *[]string) {
	if api.Spec.Context == nil {
		return
	}

	ctxRef := api.Spec.Context.DeepCopy()
	if ctxRef.Namespace == "" {
		ctxRef.Namespace = api.Namespace
	}

	*fields = append(*fields, ensureNamespacedRef(api, api.Spec.Context))
}

func indexApiV4ManagementContexts(api *v1alpha1.ApiV4Definition, fields *[]string) {
	if api.Spec.Context == nil {
		return
	}

	*fields = append(*fields, ensureNamespacedRef(api, api.Spec.Context))
}

func indexManagementContextSecrets(context *v1alpha1.ManagementContext, fields *[]string) {
	if context.Spec.HasSecretRef() {
		*fields = append(*fields, ensureNamespacedRef(context, context.Spec.SecretRef()))
	}
}

func indexApiResourceRefs(api *v1alpha1.ApiDefinition, fields *[]string) {
	if api.Spec.Resources == nil {
		return
	}

	for _, resource := range api.Spec.Resources {
		if resource.IsRef() {
			*fields = append(*fields, ensureNamespacedRef(api, resource.Ref))
		}
	}
}

func indexApiV4ResourceRefs(api *v1alpha1.ApiV4Definition, fields *[]string) {
	if api.Spec.Resources == nil {
		return
	}

	for _, resource := range api.Spec.Resources {
		if resource.IsRef() {
			*fields = append(*fields, ensureNamespacedRef(api, resource.Ref))
		}
	}
}

func indexNotificationRefs(api *v1alpha1.ApiDefinition, fields *[]string) {
	if api.Spec.NotificationsRefs == nil {
		return
	}

	for _, ref := range api.Spec.NotificationsRefs {
		r := ref
		*fields = append(*fields, ensureNamespacedRef(api, &r))
	}
}

func indexApiV4NotificationRefs(api *v1alpha1.ApiV4Definition, fields *[]string) {
	if api.Spec.NotificationsRefs == nil {
		return
	}

	for _, ref := range api.Spec.NotificationsRefs {
		r := ref
		*fields = append(*fields, ensureNamespacedRef(api, &r))
	}
}

func indexApiV4FlowsSharedPolicyGroupsRefs(api *v1alpha1.ApiV4Definition, fields *[]string) {
	for _, sharedPolicyGroup := range api.Spec.GetAllSharedPolicyGroups() {
		*fields = append(*fields, ensureNamespacedRef(api, sharedPolicyGroup))
	}
}

func indexApiGroupsRefs(api *v1alpha1.ApiDefinition, fields *[]string) {
	for i := 0; i < len(api.Spec.GroupRefs); i++ {
		*fields = append(*fields, ensureNamespacedRef(api, &api.Spec.GroupRefs[i]))
	}
}

func indexApiV4GroupsRefs(api *v1alpha1.ApiV4Definition, fields *[]string) {
	for i := 0; i < len(api.Spec.GroupRefs); i++ {
		*fields = append(*fields, ensureNamespacedRef(api, &api.Spec.GroupRefs[i]))
	}
}

func indexApiTemplate(ing *v1.Ingress, fields *[]string) {
	if ing.Annotations[core.IngressTemplateAnnotation] == "" {
		return
	}

	*fields = append(*fields, ing.Namespace+"/"+ing.Annotations[core.IngressTemplateAnnotation])
}

func indexTLSSecret(ing *v1.Ingress, fields *[]string) {
	if !k8s.IsGraviteeIngress(ing) {
		return
	}

	if len(ing.Spec.TLS) == 0 {
		return
	}

	for i := range ing.Spec.TLS {
		*fields = append(*fields, ing.Namespace+"/"+ing.Spec.TLS[i].SecretName)
	}
}

func indexApplicationManagementContexts(application *v1alpha1.Application, fields *[]string) {
	if application.Spec.Context == nil {
		return
	}

	*fields = append(*fields, ensureNamespacedRef(application, application.Spec.Context))
}

func indexAPIv4Subscriptions(sub *v1alpha1.Subscription, fields *[]string) {
	kind := sub.Spec.API.Kind
	ns := sub.Spec.API.Namespace
	if kind == "" {
		kind = core.CRDApiV4DefinitionResource
	}
	kind = dynamic.ResourceFromKind(kind)
	if ns == "" {
		ns = sub.GetNamespace()
	}
	nsn := refs.NamespacedName{
		Name:      sub.Spec.API.Name,
		Namespace: ns,
	}
	if kind == core.CRDApiV4DefinitionResource {
		*fields = append(*fields, nsn.String())
	}
}

func indexAPIv2Subscriptions(sub *v1alpha1.Subscription, fields *[]string) {
	kind := sub.Spec.API.Kind
	ns := sub.Spec.API.Namespace
	if kind == "" {
		kind = core.CRDApiV4DefinitionResource
	}
	kind = dynamic.ResourceFromKind(kind)
	if ns == "" {
		ns = sub.GetNamespace()
	}
	nsn := refs.NamespacedName{
		Name:      sub.Spec.API.Name,
		Namespace: ns,
	}
	if kind == core.CRDApiDefinitionResource {
		*fields = append(*fields, nsn.String())
	}
}

func indexSharedPolicyGroupManagementContexts(spg *v1alpha1.SharedPolicyGroup, fields *[]string) {
	if spg.Spec.Context == nil {
		return
	}

	*fields = append(*fields, ensureNamespacedRef(spg, spg.Spec.Context))
}

func ensureNamespacedRef(obj client.Object, ref core.ObjectRef) string {
	cp := refs.NewNamespacedName(ref.GetNamespace(), ref.GetName())
	if cp.Namespace == "" {
		cp.Namespace = obj.GetNamespace()
	}
	return cp.String()
}
