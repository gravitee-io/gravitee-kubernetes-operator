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
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const separator = "/"

// ApiV4DefinitionSpec defines the desired state of ApiDefinition.
// +kubebuilder:object:generate=true
type ApiV4DefinitionSpec struct {
	v4.Api  `json:",inline"`
	Context *refs.NamespacedName `json:"contextRef,omitempty"`
}

// ApiV4DefinitionStatus defines the observed state of API Definition.
type ApiV4DefinitionStatus struct {
	base.Status `json:",inline"`
}

var _ core.ApiDefinitionObject = &ApiV4Definition{}
var _ core.SubscribableStatus = &ApiV4DefinitionStatus{}
var _ core.Spec = &ApiDefinitionV2Spec{}

// ApiV4Definition is the Schema for the v4 apidefinitions API.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="State",type=string,JSONPath=`.spec.state`,description="State"
// +kubebuilder:printcolumn:name="Lifecycle State",type=string,JSONPath=`.spec.lifecycleState`,description="Lifecycle State"
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`,description="API version."
// +kubebuilder:resource:shortName=graviteev4apis
// +kubebuilder:storageversion
type ApiV4Definition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApiV4DefinitionSpec   `json:"spec,omitempty"`
	Status ApiV4DefinitionStatus `json:"status,omitempty"`
}

func (api *ApiV4Definition) ToGatewayDefinition() v4.GatewayDefinitionApi {
	cp := api.DeepCopy()
	return cp.Spec.Api.ToGatewayDefinition()
}

func (api *ApiV4Definition) IsBeingDeleted() bool {
	return !api.ObjectMeta.DeletionTimestamp.IsZero()
}

func (api *ApiV4Definition) GetResources() []core.ObjectOrRef[core.ResourceModel] {
	return api.Spec.GetResources()
}

func (api *ApiV4Definition) GetState() string {
	return api.Spec.GetState()
}

func (api *ApiV4Definition) HasPlans() bool {
	return api.Spec.HasPlans()
}

func (api *ApiV4Definition) GetPlan(name string) core.PlanModel {
	return api.Spec.GetPlan(name)
}

func (api *ApiV4Definition) IsStopped() bool {
	return api.Spec.IsStopped()
}

func (api *ApiV4Definition) IsSyncFromManagement() bool {
	defCtx := api.Spec.DefinitionContext
	return defCtx == nil || defCtx.SyncFrom == v4.OriginManagement
}

func (api *ApiV4Definition) PopulateIDs(context core.ContextModel) {
	api.Spec.ID = api.pickID(context)
	api.Spec.CrossID = api.pickCrossID()
	api.Spec.Pages = api.pickPageIDs()
	api.Spec.Plans = api.pickPlanIDs()
}

// pickID returns the ID of the API definition, when a context has been defined at the spec level.
// The ID might be returned from the API status, meaning that the API is already known.
// If the API is unknown, the ID is either given from the spec if given,
// or generated from the API UID and the context key to ensure uniqueness
// in case the API is replicated on a same APIM instance.
func (api *ApiV4Definition) pickID(mCtx core.ContextModel) string {
	if api.Status.ID != "" {
		return api.Status.ID
	}

	if api.Spec.ID != "" {
		return api.Spec.ID
	}

	if mCtx != nil {
		return uuid.FromStrings(api.pickCrossID(), mCtx.GetOrgID(), mCtx.GetEnvID())
	}

	return string(api.UID)
}

func (api *ApiV4Definition) pickCrossID() string {
	if api.Status.CrossID != "" {
		return api.Status.CrossID
	}

	if api.Spec.CrossID != "" {
		return api.Spec.CrossID
	}

	namespacedName := api.GetNamespacedName()
	return uuid.FromStrings(namespacedName.String())
}

func (api *ApiV4Definition) pickPlanIDs() *map[string]*v4.Plan {
	if !api.HasPlans() {
		return nil
	}

	plans := make(map[string]*v4.Plan, len(*api.Spec.Plans))
	for key, plan := range *api.Spec.Plans {
		p := plan.DeepCopy()
		if id, ok := api.Status.Plans[key]; ok {
			p.ID = id
		} else if plan.ID == "" {
			namespacedName := api.GetNamespacedName()
			p.ID = uuid.FromStrings(namespacedName.String(), key)
		}
		plans[key] = p
	}
	return &plans
}

func (api *ApiV4Definition) pickPageIDs() *map[string]*v4.Page {
	if api.Spec.Pages == nil {
		return nil
	}

	pages := make(map[string]*v4.Page, len(*api.Spec.Pages))
	for name, page := range *api.Spec.Pages {
		p := page.DeepCopy()

		p.API = &api.Spec.ID
		apiName := api.GetNamespacedName().String()
		if page.ID == "" {
			p.ID = uuid.FromStrings(api.Spec.ID, separator, name)
		}
		if page.CrossID == "" {
			p.CrossID = uuid.FromStrings(apiName, separator, name)
		}
		if page.Parent != nil {
			pID := uuid.FromStrings(api.Spec.ID, separator, *page.Parent)
			p.ParentID = &pID
		}

		pages[name] = p
	}
	return &pages
}

// GetEnvID implements custom.ApiDefinition.
func (api *ApiV4Definition) GetEnvID() string {
	return api.Status.EnvID
}

// GetID implements custom.ApiDefinition.
func (api *ApiV4Definition) GetID() string {
	return api.Status.ID
}

// GetOrgID implements custom.ApiDefinition.
func (api *ApiV4Definition) GetOrgID() string {
	return api.Status.OrgID
}

func (api *ApiV4Definition) GetSpec() core.Spec {
	return &api.Spec
}

func (api *ApiV4Definition) GetStatus() core.Status {
	return &api.Status
}

func (api *ApiV4Definition) DeepCopyResource() core.Object {
	return api.DeepCopy()
}

func (api *ApiV4Definition) ContextRef() core.ObjectRef {
	return api.Spec.Context
}

func (api *ApiV4Definition) HasContext() bool {
	return api.Spec.Context != nil
}

func (api *ApiV4Definition) Version() core.ApiDefinitionVersion {
	return core.ApiV4
}

func (api *ApiV4Definition) GetNamespacedName() *refs.NamespacedName {
	return &refs.NamespacedName{Namespace: api.Namespace, Name: api.Name}
}

func (api *ApiV4Definition) GetObjectMeta() *metav1.ObjectMeta {
	return &api.ObjectMeta
}

func (api *ApiV4Definition) GetContextPaths() []string {
	return api.Spec.GetContextPaths()
}

func (api *ApiV4Definition) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      api.Name,
		Namespace: api.Namespace,
	}
}

func (api *ApiV4Definition) GetDefinitionVersion() core.ApiDefinitionVersion {
	return core.ApiV4
}

func (api *ApiV4Definition) GetDefinitionContext() core.DefinitionContext {
	return api.Spec.GetDefinitionContext()
}

func (api *ApiV4Definition) SetDefinitionContext(ctx core.DefinitionContext) {
	api.Spec.SetDefinitionContext(ctx)
}

func (api *ApiV4Definition) GetDefinition() core.ApiDefinitionModel {
	return &api.Spec.Api
}

func (spec *ApiV4DefinitionSpec) Hash() string {
	return hash.Calculate(spec)
}

func (spec *ApiV4DefinitionSpec) GetManagementContext() *refs.NamespacedName {
	return spec.Context
}

func (s *ApiV4DefinitionStatus) SetProcessingStatus(status core.ProcessingStatus) {
	s.ProcessingStatus = status
}

func (s *ApiV4DefinitionStatus) IsFailed() bool {
	return s.ProcessingStatus == core.ProcessingStatusFailed
}

func (s *ApiV4DefinitionStatus) DeepCopyFrom(api client.Object) error {
	switch t := api.(type) {
	case *ApiV4Definition:
		subscriptionCount := s.Status.SubscriptionCount
		t.Status.DeepCopyInto(s)
		s.Status.SubscriptionCount = subscriptionCount
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *ApiV4DefinitionStatus) DeepCopyTo(api client.Object) error {
	switch t := api.(type) {
	case *ApiV4Definition:
		subscriptionCount := t.Status.SubscriptionCount
		s.DeepCopyInto(&t.Status)
		t.Status.SubscriptionCount = subscriptionCount
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *ApiV4DefinitionStatus) AddSubscription() {
	s.SubscriptionCount += 1
}

// GetSubscriptionCount implements core.SubscribableStatus.
func (s *ApiV4DefinitionStatus) GetSubscriptionCount() uint {
	return s.SubscriptionCount
}

// RemoveSubscription implements core.SubscribableStatus.
func (s *ApiV4DefinitionStatus) RemoveSubscription() {
	if s.SubscriptionCount > 0 {
		s.SubscriptionCount -= 1
	}
}

// ApiV4DefinitionList contains a list of ApiV4Definition.
// +kubebuilder:object:root=true
type ApiV4DefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApiV4Definition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApiV4Definition{}, &ApiV4DefinitionList{})
}
