/*
Copyright 2022 DAVID BRASSELY.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"fmt"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"

	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// The API definition is the main resource handled by the Kubernetes Operator
// Most of the configuration properties defined here are already documented
// in the APIM Console API Reference.
// See https://docs.gravitee.io/apim/3.x/apim_installguide_rest_apis_documentation.html
// +kubebuilder:object:generate=true
type ApiDefinitionV2Spec struct {
	v2.Api  `json:",inline"`
	Context *refs.NamespacedName `json:"contextRef,omitempty"`
	// local defines if the api is local or not.
	//
	// If true, the Operator will create the ConfigMaps for the Gateway and pushes the API to the Management API
	// but without setting the update flag in the datastore.
	//
	// If false, the Operator will not create the ConfigMaps for the Gateway.
	// Instead, it pushes the API to the Management API and forces it to update the event in the datastore.
	// This will cause Gateways to fetch the APIs from the datastore
	//
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=false
	IsLocal bool `json:"local"`
}

// ApiDefinitionStatus defines the observed state of API Definition.
type ApiDefinitionStatus struct {
	base.Status `json:",inline"`
}

var _ core.ApiDefinitionObject = &ApiDefinition{}
var _ core.SubscribableStatus = &ApiDefinitionStatus{}
var _ core.Spec = &ApiDefinitionV2Spec{}

// ApiDefinition is the Schema for the apidefinitions API.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Entrypoint",type=string,JSONPath=`.spec.proxy.virtual_hosts[*].path`,description="API entrypoint."
// +kubebuilder:printcolumn:name="Endpoint",type=string,JSONPath=`.spec.proxy.groups[*].endpoints[*].target`,description="API endpoint."
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`,description="API version."
// +kubebuilder:resource:shortName=graviteeapis
// +kubebuilder:storageversion
type ApiDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApiDefinitionV2Spec `json:"spec,omitempty"`
	Status ApiDefinitionStatus `json:"status,omitempty"`
}

func (api *ApiDefinition) GetNamespacedName() *refs.NamespacedName {
	return &refs.NamespacedName{Namespace: api.Namespace, Name: api.Name}
}

func (spec *ApiDefinitionV2Spec) EnsureDefinitionContext() {
	if spec.DefinitionContext == nil {
		spec.DefinitionContext = &v2.DefinitionContext{
			Mode:   v2.ModeFullyManaged,
			Origin: v2.OriginKubernetes,
		}
	}

	if !spec.IsLocal || strings.EqualFold(v2.OriginManagement, spec.DefinitionContext.SyncFrom) {
		spec.DefinitionContext.SyncFrom = strings.ToUpper(v2.OriginManagement)
	} else {
		spec.DefinitionContext.SyncFrom = strings.ToUpper(v2.OriginKubernetes)
	}
}

func (api *ApiDefinition) GetEnvID() string {
	return api.Status.EnvID
}

func (api *ApiDefinition) GetID() string {
	return api.Status.ID
}

func (api *ApiDefinition) GetOrgID() string {
	return api.Status.OrgID
}

func (api *ApiDefinition) GetDefinitionVersion() core.ApiDefinitionVersion {
	return core.ApiV2
}

func (api *ApiDefinition) GetSpec() core.Spec {
	return &api.Spec
}

func (api *ApiDefinition) GetStatus() core.Status {
	return &api.Status
}

func (api *ApiDefinition) DeepCopyResource() core.Object {
	return api.DeepCopy()
}

func (api *ApiDefinition) ContextRef() core.ObjectRef {
	return api.Spec.Context
}

func (api *ApiDefinition) HasContext() bool {
	return api.Spec.Context != nil
}

func (api *ApiDefinition) GetContextPaths() []string {
	return api.Spec.GetContextPaths()
}

func (api *ApiDefinition) GetDefinition() core.ApiDefinitionModel {
	return &api.Spec.Api
}

func (api *ApiDefinition) GetDefinitionContext() core.DefinitionContext {
	return api.Spec.GetDefinitionContext()
}

func (api *ApiDefinition) SetDefinitionContext(ctx core.DefinitionContext) {
	api.Spec.SetDefinitionContext(ctx)
}

func (api *ApiDefinition) GetState() string {
	return api.Spec.GetState()
}

func (api *ApiDefinition) HasPlans() bool {
	return api.Spec.HasPlans()
}

func (api *ApiDefinition) IsStopped() bool {
	return api.Spec.IsStopped()
}

func (api *ApiDefinition) GetPlan(name string) core.PlanModel {
	return api.Spec.GetPlan(name)
}

func (api *ApiDefinition) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      api.Name,
		Namespace: api.Namespace,
	}
}

func (api *ApiDefinition) PopulateIDs(_ core.ContextModel) {
	api.Spec.ID = api.pickID()
	api.Spec.CrossID = api.pickCrossID()
	api.generateEmptyPlanIDs()
	api.generatePageIDs()
}

func (api *ApiDefinition) GetResources() []core.ObjectOrRef[core.ResourceModel] {
	return api.Spec.GetResources()
}

func (api *ApiDefinition) IsSyncFromManagement() bool {
	isNotLocal := !api.Spec.IsLocal
	defCtx := api.Spec.DefinitionContext
	isSyncFromManagement := defCtx != nil && defCtx.SyncFrom == v2.OriginManagement
	return isNotLocal || isSyncFromManagement
}

// For each plan, generate a Cross id from Api id & Plan Name if not defined.
func (api *ApiDefinition) generateEmptyPlanIDs() {
	plans := api.Spec.Plans

	for _, plan := range plans {
		if plan.CrossID == "" {
			plan.CrossID = uuid.FromStrings(api.Spec.ID, separator, plan.Name)
		}

		if id, ok := api.Status.Plans[plan.CrossID]; ok {
			plan.ID = id
		} else {
			plan.ID = uuid.FromStrings(plan.CrossID, separator, plan.Name)
		}
	}
}

func (api *ApiDefinition) generatePageIDs() {
	if api.Spec.Pages == nil {
		return
	}

	spec := &api.Spec
	pages := spec.Pages
	for name, page := range *pages {
		page.API = &spec.ID
		apiName := api.GetNamespacedName().String()
		if page.CrossID == "" {
			page.CrossID = uuid.FromStrings(apiName, separator, name)
		}
		if page.ID == "" {
			page.ID = uuid.FromStrings(spec.ID, separator, name)
		}
		if page.Parent != nil {
			pID := uuid.FromStrings(spec.ID, separator, *page.Parent)
			page.ParentID = &pID
		}
	}
}

// PickID returns the ID of the API definition, when a context has been defined at the spec level.
// The ID might be returned from the API status, meaning that the API is already known.
// If the API is unknown, the ID is either given from the spec if given,
// or generated from the API UID and the context key to ensure uniqueness
// in case the API is replicated on a same APIM instance.
func (api *ApiDefinition) pickID() string {
	if api.Status.ID != "" {
		return api.Status.ID
	}

	if api.Spec.ID != "" {
		return api.Spec.ID
	}

	return string(api.UID)
}

func (api *ApiDefinition) pickCrossID() string {
	if api.Status.CrossID != "" {
		return api.Status.CrossID
	}

	return api.getOrGenerateCrossID()
}

func (api *ApiDefinition) getOrGenerateCrossID() string {
	if api.Spec.CrossID != "" {
		return api.Spec.CrossID
	}

	return uuid.FromStrings(api.GetNamespacedName().String())
}

func (spec *ApiDefinitionV2Spec) Hash() string {
	return hash.Calculate(spec)
}

func (s *ApiDefinitionStatus) SetProcessingStatus(status core.ProcessingStatus) {
	s.ProcessingStatus = status
}

func (s *ApiDefinitionStatus) IsFailed() bool {
	return s.ProcessingStatus == core.ProcessingStatusFailed
}

func (s *ApiDefinitionStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *ApiDefinition:
		t.Status.DeepCopyInto(s)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *ApiDefinitionStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *ApiDefinition:
		s.DeepCopyInto(&t.Status)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *ApiDefinitionStatus) AddSubscription() {
	s.SubscriptionCount += 1
}

func (s *ApiDefinitionStatus) RemoveSubscription() {
	s.SubscriptionCount -= 1
}

func (s *ApiDefinitionStatus) GetSubscriptionCount() uint {
	return s.SubscriptionCount
}

// ApiDefinitionList contains a list of ApiDefinition.
// +kubebuilder:object:root=true
type ApiDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApiDefinition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApiDefinition{}, &ApiDefinitionList{})
}
