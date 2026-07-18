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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/gateway"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/gateway-api/httproute/internal/mapper"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/uuid"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	configMapPrefix     = "httproute-"
	gatewayConfigMapPfx = "gw-route-"
	gatewayAPIPfx       = "gw-route-"
	maxK8sNameLength    = 253
	hashSuffixLength    = 8
	definitionKey       = "definition"
	apiDefVersionKey    = "apiDefinitionVersion"
	apiDefVersionValue  = "4.0.0"
	orgKey              = "organizationId"
	envKey              = "environmentId"
	defaultOrgID        = "DEFAULT"
	defaultEnvID        = "DEFAULT"
	managedByKey        = "managed-by"
	gioTypeKey          = "gio-type"
	gatewayOwnerKey     = "gateway-owner"
	routeLabelPrefix    = "httproute.gravitee.io/"
)

// ---------------------------------------------------------------------------
// Program creates one ApiV4Definition CR per HTTPRoute (default path).
// ---------------------------------------------------------------------------

func Program(ctx context.Context, route *gwAPIv1.HTTPRoute) error {
	api, err := mapper.Map(ctx, route)
	if err != nil {
		return err
	}
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		if err := k8s.GetClient().Get(ctx, client.ObjectKeyFromObject(api), api); client.IgnoreNotFound(err) != nil {
			return err
		}

		return k8s.CreateOrUpdate(ctx, api, func() error {
			api.SetOwnerReferences(getOwnerReferences(route))
			if api.Labels == nil {
				api.Labels = make(map[string]string)
			}
			for k, v := range routeLabels(route) {
				api.Labels[k] = v
			}
			spec, err := mapper.MapSpec(ctx, route)
			if err != nil {
				return err
			}
			api.Spec = spec
			return nil
		})
	})
}

// ---------------------------------------------------------------------------
// ProgramConfigMap creates one ConfigMap per HTTPRoute (skipAPIDefinition
// without matchAcrossRoutes). This is the original per-route behaviour that
// preserves backward-compatible naming and API IDs.
// ---------------------------------------------------------------------------

func ProgramConfigMap(ctx context.Context, route *gwAPIv1.HTTPRoute) error {
	api, err := mapper.Map(ctx, route)
	if err != nil {
		return err
	}

	api.Spec.ID = uuid.FromStrings(string(route.UID))
	api.Spec.CrossID = uuid.FromStrings(route.Namespace, route.Name)

	if api.Spec.Plans != nil {
		for key, plan := range *api.Spec.Plans {
			plan.ID = uuid.FromStrings(api.Spec.ID, key)
		}
	}

	gwDef := api.Spec.Api.ToGatewayDefinition()
	jsonSpec, err := json.Marshal(gwDef)
	if err != nil {
		return err
	}

	cmLabels := map[string]string{
		managedByKey: core.CRDGroup,
		gioTypeKey:   core.CRDApiDefinitionResource + "." + core.CRDGroup,
	}
	for k, v := range routeLabels(route) {
		cmLabels[k] = v
	}

	cm := &v1.ConfigMap{
		ObjectMeta: metaV1.ObjectMeta{
			Name:            buildConfigMapName(route),
			Namespace:       route.Namespace,
			OwnerReferences: getOwnerReferences(route),
			Labels:          cmLabels,
		},
		Data: map[string]string{
			definitionKey:    string(jsonSpec),
			apiDefVersionKey: apiDefVersionValue,
			orgKey:           defaultOrgID,
			envKey:           defaultEnvID,
		},
	}

	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		existing := &v1.ConfigMap{}
		err := k8s.GetClient().Get(ctx, client.ObjectKeyFromObject(cm), existing)
		if errors.IsNotFound(err) {
			return k8s.GetClient().Create(ctx, cm)
		}
		if err != nil {
			return err
		}
		existing.Labels = cm.Labels
		existing.Data = cm.Data
		existing.SetOwnerReferences(cm.GetOwnerReferences())
		return k8s.GetClient().Update(ctx, existing)
	})
}

// ---------------------------------------------------------------------------
// ProgramMergedConfigMaps merges HTTPRoutes with overlapping context paths
// into a single ConfigMap per group (skipAPIDefinition + matchAcrossRoutes).
// ---------------------------------------------------------------------------

func ProgramMergedConfigMaps(ctx context.Context, route *gwAPIv1.HTTPRoute) error {
	for _, gw := range resolveAcceptedParentGateways(ctx, route) {
		if err := programGatewayMergedConfigMaps(ctx, gw); err != nil {
			return err
		}
	}

	return nil
}

func programGatewayMergedConfigMaps(ctx context.Context, gw *gwAPIv1.Gateway) error {
	allRoutes, err := listRoutesForGateway(ctx, gw)
	if err != nil {
		return err
	}

	if len(allRoutes) == 0 {
		return deleteAllMergedConfigMaps(ctx, gw)
	}

	mapped, err := mapAllRoutes(ctx, allRoutes)
	if err != nil {
		return err
	}
	groups := groupByOverlappingPaths(mapped)

	activeNames := make(map[string]bool)
	for _, group := range groups {
		cmName, err := programGroupConfigMap(ctx, gw, group)
		if err != nil {
			return err
		}
		activeNames[cmName] = true
	}

	return cleanupStaleMergedConfigMaps(ctx, gw, activeNames)
}

func programGroupConfigMap(ctx context.Context, gw *gwAPIv1.Gateway, group []mappedRoute) (string, error) {
	merged, err := mergeGroupSpecs(ctx, group)
	if err != nil {
		return "", err
	}

	groupKey := buildGroupKey(group)
	merged.ID = uuid.FromStrings(string(gw.UID), groupKey)
	merged.CrossID = uuid.FromStrings(gw.Namespace, gw.Name, groupKey)
	populatePlanIDs(&merged)
	ownerRefs := buildGroupOwnerRefs(group, gw.Namespace)

	gwDef := merged.Api.ToGatewayDefinition()
	jsonSpec, err := json.Marshal(gwDef)
	if err != nil {
		return "", err
	}

	cmName := buildGroupConfigMapName(gw, groupKey)
	cmLabels := map[string]string{
		managedByKey:    core.CRDGroup,
		gioTypeKey:      core.CRDApiDefinitionResource + "." + core.CRDGroup,
		gatewayOwnerKey: gatewayHash(gw),
	}
	for k, v := range groupRouteLabels(group) {
		cmLabels[k] = v
	}

	cm := &v1.ConfigMap{
		ObjectMeta: metaV1.ObjectMeta{
			Name:            cmName,
			Namespace:       gw.Namespace,
			OwnerReferences: ownerRefs,
			Labels:          cmLabels,
		},
		Data: map[string]string{
			definitionKey:    string(jsonSpec),
			apiDefVersionKey: apiDefVersionValue,
			orgKey:           defaultOrgID,
			envKey:           defaultEnvID,
		},
	}

	err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		existing := &v1.ConfigMap{}
		getErr := k8s.GetClient().Get(ctx, client.ObjectKeyFromObject(cm), existing)
		if errors.IsNotFound(getErr) {
			return k8s.GetClient().Create(ctx, cm)
		}
		if getErr != nil {
			return getErr
		}
		existing.Labels = cm.Labels
		existing.Data = cm.Data
		existing.SetOwnerReferences(cm.GetOwnerReferences())
		return k8s.GetClient().Update(ctx, existing)
	})
	if err != nil {
		return "", err
	}

	return cmName, nil
}

// ---------------------------------------------------------------------------
// ProgramMergedAPIs merges HTTPRoutes with overlapping context paths into a
// single ApiV4Definition CR per group (matchAcrossRoutes without
// skipAPIDefinition).
// ---------------------------------------------------------------------------

func ProgramMergedAPIs(ctx context.Context, route *gwAPIv1.HTTPRoute) error {
	for _, gw := range resolveAcceptedParentGateways(ctx, route) {
		if err := programGatewayMergedAPIs(ctx, gw); err != nil {
			return err
		}
	}

	return nil
}

func programGatewayMergedAPIs(ctx context.Context, gw *gwAPIv1.Gateway) error {
	allRoutes, err := listRoutesForGateway(ctx, gw)
	if err != nil {
		return err
	}

	if len(allRoutes) == 0 {
		return deleteAllMergedAPIs(ctx, gw)
	}

	mapped, err := mapAllRoutes(ctx, allRoutes)
	if err != nil {
		return err
	}
	groups := groupByOverlappingPaths(mapped)

	activeNames := make(map[string]bool)
	for _, group := range groups {
		apiName, err := programGroupAPI(ctx, gw, group)
		if err != nil {
			return err
		}
		activeNames[apiName] = true
	}

	return cleanupStaleMergedAPIs(ctx, gw, activeNames)
}

func programGroupAPI(ctx context.Context, gw *gwAPIv1.Gateway, group []mappedRoute) (string, error) {
	merged, err := mergeGroupSpecs(ctx, group)
	if err != nil {
		return "", err
	}

	groupKey := buildGroupKey(group)
	merged.ID = uuid.FromStrings(string(gw.UID), groupKey)
	merged.CrossID = uuid.FromStrings(gw.Namespace, gw.Name, groupKey)
	populatePlanIDs(&merged)

	apiName := buildGroupAPIName(gw, groupKey)
	ownerRefs := buildGroupOwnerRefs(group, gw.Namespace)
	newHash := merged.Hash()

	api := &v1alpha1.ApiV4Definition{}
	api.Name = apiName
	api.Namespace = gw.Namespace

	err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		if err := k8s.GetClient().Get(ctx, client.ObjectKeyFromObject(api), api); client.IgnoreNotFound(err) != nil {
			return err
		}

		return k8s.CreateOrUpdate(ctx, api, func() error {
			api.SetOwnerReferences(ownerRefs)
			if api.Labels == nil {
				api.Labels = make(map[string]string)
			}
			api.Labels[managedByKey] = core.CRDGroup
			api.Labels[gatewayOwnerKey] = gatewayHash(gw)
			for k, v := range groupRouteLabels(group) {
				api.Labels[k] = v
			}

			if api.Annotations != nil && api.Annotations[core.LastSpecHashAnnotation] == newHash {
				return nil
			}
			api.Spec = merged
			return nil
		})
	})
	if err != nil {
		return "", err
	}

	return apiName, nil
}

func mergeGroupSpecs(ctx context.Context, group []mappedRoute) (v1alpha1.ApiV4DefinitionSpec, error) {
	specs := make([]v1alpha1.ApiV4DefinitionSpec, len(group))
	for i, mr := range group {
		if len(group) == 1 {
			spec, err := mapper.MapSpec(ctx, mr.route)
			if err != nil {
				return v1alpha1.ApiV4DefinitionSpec{}, err
			}
			specs[i] = spec
		} else {
			specs[i] = mr.spec
		}
	}
	return mapper.MergeSpecs(specs), nil
}

func populatePlanIDs(spec *v1alpha1.ApiV4DefinitionSpec) {
	if spec.Plans != nil {
		for key, plan := range *spec.Plans {
			plan.ID = uuid.FromStrings(spec.ID, key)
		}
	}
}

func buildGroupOwnerRefs(group []mappedRoute, gwNamespace string) []metaV1.OwnerReference {
	var refs []metaV1.OwnerReference
	for _, mr := range group {
		if mr.route.Namespace == gwNamespace {
			refs = append(refs, metaV1.OwnerReference{
				Kind:       "HTTPRoute",
				APIVersion: gwAPIv1.GroupVersion.String(),
				Name:       mr.route.GetName(),
				UID:        mr.route.GetUID(),
			})
		}
	}
	return refs
}

func buildGroupAPIName(gw *gwAPIv1.Gateway, groupKey string) string {
	gwH := gatewayHash(gw)
	grpH := sha256.Sum256([]byte(groupKey))
	grpSuffix := hex.EncodeToString(grpH[:])[:hashSuffixLength]
	return fmt.Sprintf("%s%s-%s", gatewayAPIPfx, gwH, grpSuffix)
}

func deleteAllMergedAPIs(ctx context.Context, gw *gwAPIv1.Gateway) error {
	return cleanupStaleMergedAPIs(ctx, gw, nil)
}

func cleanupStaleMergedAPIs(ctx context.Context, gw *gwAPIv1.Gateway, activeNames map[string]bool) error {
	apiList := &v1alpha1.ApiV4DefinitionList{}
	labelSelector := labels.SelectorFromSet(map[string]string{
		managedByKey:    core.CRDGroup,
		gatewayOwnerKey: gatewayHash(gw),
	})
	if err := k8s.GetClient().List(ctx, apiList,
		&client.ListOptions{
			Namespace:     gw.Namespace,
			LabelSelector: labelSelector,
		},
	); err != nil {
		return err
	}

	for i := range apiList.Items {
		api := &apiList.Items[i]
		if activeNames != nil && activeNames[api.Name] {
			continue
		}
		if err := k8s.GetClient().Delete(ctx, api); client.IgnoreNotFound(err) != nil {
			return err
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Shared helpers
// ---------------------------------------------------------------------------

type mappedRoute struct {
	route *gwAPIv1.HTTPRoute
	spec  v1alpha1.ApiV4DefinitionSpec
}

func resolveAcceptedParentGateways(ctx context.Context, route *gwAPIv1.HTTPRoute) []*gwAPIv1.Gateway {
	seen := make(map[string]bool)
	var gateways []*gwAPIv1.Gateway

	for _, ref := range route.Spec.ParentRefs {
		if !k8s.IsGatewayKind(ref) {
			continue
		}
		if ps := findParentStatus(route, ref); ps != nil {
			if !k8s.IsAccepted(gateway.WrapRouteParentStatus(ps)) {
				continue
			}
		}
		gw, err := k8s.ResolveGateway(ctx, route.ObjectMeta, ref)
		if err != nil {
			continue
		}
		key := gw.Namespace + "/" + gw.Name
		if !seen[key] {
			seen[key] = true
			gateways = append(gateways, gw)
		}
	}

	return gateways
}

func findParentStatus(route *gwAPIv1.HTTPRoute, ref gwAPIv1.ParentReference) *gwAPIv1.RouteParentStatus {
	refNS := route.Namespace
	if ref.Namespace != nil {
		refNS = string(*ref.Namespace)
	}
	for i := range route.Status.Parents {
		p := &route.Status.Parents[i]
		pNS := route.Namespace
		if p.ParentRef.Namespace != nil {
			pNS = string(*p.ParentRef.Namespace)
		}
		if p.ParentRef.Name == ref.Name && pNS == refNS && sectionNameEquals(p.ParentRef.SectionName, ref.SectionName) {
			return p
		}
	}
	return nil
}

func sectionNameEquals(a, b *gwAPIv1.SectionName) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func mapAllRoutes(ctx context.Context, allRoutes []gwAPIv1.HTTPRoute) ([]mappedRoute, error) {
	sort.Slice(allRoutes, func(i, j int) bool {
		return allRoutes[i].Name < allRoutes[j].Name
	})

	mapped := make([]mappedRoute, 0, len(allRoutes))
	for i := range allRoutes {
		r := &allRoutes[i]
		prefix := r.Name + "-"
		spec, err := mapper.MapSpecWithPrefix(ctx, r, prefix)
		if err != nil {
			return nil, err
		}
		mapped = append(mapped, mappedRoute{route: r, spec: spec})
	}
	return mapped, nil
}

func groupByOverlappingPaths(mapped []mappedRoute) [][]mappedRoute {
	n := len(mapped)
	if n == 0 {
		return nil
	}

	parent := make([]int, n)
	for i := range parent {
		parent[i] = i
	}

	var find = func(i int) int {
		for parent[i] != i {
			parent[i] = parent[parent[i]]
			i = parent[i]
		}
		return i
	}

	union := func(i, j int) {
		ri, rj := find(i), find(j)
		if ri != rj {
			parent[ri] = rj
		}
	}

	pathOwners := make(map[string][]int)
	for i, mr := range mapped {
		for _, p := range mr.spec.GetContextPaths() {
			pathOwners[p] = append(pathOwners[p], i)
		}
	}

	for _, indices := range pathOwners {
		for j := 1; j < len(indices); j++ {
			union(indices[0], indices[j])
		}
	}

	groupMap := make(map[int][]mappedRoute)
	for i, mr := range mapped {
		root := find(i)
		groupMap[root] = append(groupMap[root], mr)
	}

	groups := make([][]mappedRoute, 0, len(groupMap))
	for _, g := range groupMap {
		groups = append(groups, g)
	}

	return groups
}

func buildGroupKey(group []mappedRoute) string {
	names := make([]string, len(group))
	for i, mr := range group {
		names[i] = mr.route.Name
	}
	sort.Strings(names)
	return strings.Join(names, "+")
}

func gatewayHash(gw *gwAPIv1.Gateway) string {
	h := sha256.Sum256([]byte(gw.Namespace + "/" + gw.Name))
	return hex.EncodeToString(h[:])[:hashSuffixLength]
}

func buildGroupConfigMapName(gw *gwAPIv1.Gateway, groupKey string) string {
	gwH := gatewayHash(gw)
	grpH := sha256.Sum256([]byte(groupKey))
	grpSuffix := hex.EncodeToString(grpH[:])[:hashSuffixLength]
	return fmt.Sprintf("%s%s-%s", gatewayConfigMapPfx, gwH, grpSuffix)
}

func listRoutesForGateway(ctx context.Context, gw *gwAPIv1.Gateway) ([]gwAPIv1.HTTPRoute, error) {
	list := &gwAPIv1.HTTPRouteList{}
	if err := k8s.GetClient().List(ctx, list); err != nil {
		return nil, err
	}

	var attached []gwAPIv1.HTTPRoute
	for i := range list.Items {
		if referencesGatewayByNameNs(&list.Items[i], gw) {
			attached = append(attached, list.Items[i])
		}
	}
	return attached, nil
}

func referencesGatewayByNameNs(route *gwAPIv1.HTTPRoute, gw *gwAPIv1.Gateway) bool {
	for _, ref := range route.Spec.ParentRefs {
		if !k8s.IsGatewayKind(ref) {
			continue
		}
		if string(ref.Name) != gw.Name {
			continue
		}
		refNS := ref.Namespace
		if refNS == nil {
			if route.Namespace == gw.Namespace {
				return true
			}
		} else if string(*refNS) == gw.Namespace {
			return true
		}
	}
	return false
}

func deleteAllMergedConfigMaps(ctx context.Context, gw *gwAPIv1.Gateway) error {
	return cleanupStaleMergedConfigMaps(ctx, gw, nil)
}

func cleanupStaleMergedConfigMaps(ctx context.Context, gw *gwAPIv1.Gateway, activeNames map[string]bool) error {
	cmList := &v1.ConfigMapList{}
	labelSelector := labels.SelectorFromSet(map[string]string{
		managedByKey:    core.CRDGroup,
		gioTypeKey:      core.CRDApiDefinitionResource + "." + core.CRDGroup,
		gatewayOwnerKey: gatewayHash(gw),
	})
	if err := k8s.GetClient().List(ctx, cmList,
		&client.ListOptions{
			Namespace:     gw.Namespace,
			LabelSelector: labelSelector,
		},
	); err != nil {
		return err
	}

	for i := range cmList.Items {
		cm := &cmList.Items[i]
		if activeNames != nil && activeNames[cm.Name] {
			continue
		}
		if err := k8s.GetClient().Delete(ctx, cm); client.IgnoreNotFound(err) != nil {
			return err
		}
	}
	return nil
}

func buildConfigMapName(route *gwAPIv1.HTTPRoute) string {
	name := configMapPrefix + route.Name
	if len(name) <= maxK8sNameLength {
		return name
	}
	h := sha256.Sum256([]byte(route.Name))
	suffix := hex.EncodeToString(h[:])[:hashSuffixLength]
	truncateAt := maxK8sNameLength - len(configMapPrefix) - 1 - hashSuffixLength
	return fmt.Sprintf("%s%s-%s", configMapPrefix, route.Name[:truncateAt], suffix)
}

func routeLabels(routes ...*gwAPIv1.HTTPRoute) map[string]string {
	rl := make(map[string]string, len(routes))
	for _, r := range routes {
		rl[routeLabelPrefix+r.Name] = r.Namespace
	}
	return rl
}

func groupRouteLabels(group []mappedRoute) map[string]string {
	rl := make(map[string]string, len(group))
	for _, mr := range group {
		rl[routeLabelPrefix+mr.route.Name] = mr.route.Namespace
	}
	return rl
}

func getOwnerReferences(httpRoute *gwAPIv1.HTTPRoute) []metaV1.OwnerReference {
	return []metaV1.OwnerReference{
		{
			Kind:       "HTTPRoute",
			APIVersion: gwAPIv1.GroupVersion.String(),
			Name:       httpRoute.GetName(),
			UID:        httpRoute.GetUID(),
		},
	}
}
