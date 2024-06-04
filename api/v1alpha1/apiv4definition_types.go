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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/management"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/uuid"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/types/k8s/custom"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/types/list"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ApiV4DefinitionSpec defines the desired state of ApiDefinition.
// +kubebuilder:object:generate=true
type ApiV4DefinitionSpec struct {
	v4.Api  `json:",inline"`
	Context *refs.NamespacedName `json:"contextRef,omitempty"`
}

// ApiV4DefinitionStatus defines the observed state of API Definition.
type ApiV4DefinitionStatus struct {
	// The organisation ID, if a management context has been defined to sync the API with an APIM instance
	OrgID string `json:"organizationId,omitempty"`
	// The environment ID, if a management context has been defined to sync the API with an APIM instance
	EnvID string `json:"environmentId,omitempty"`
	// The ID of the API definition if a management context has been defined to sync the API with an APIM instance
	ID string `json:"id,omitempty"`
	// The Cross ID of the API definition if a management context has been defined to sync the API with an APIM instance
	CrossID string `json:"crossId,omitempty"`

	// The processing status of the API definition.
	Status custom.ProcessingStatus `json:"processingStatus,omitempty"`

	// The state of the API. Can be either STARTED or STOPPED.
	State string `json:"state,omitempty"`

	// This field is used to store the list of plans that have been created
	// for the API definition if a management context has been defined
	// to sync the API with an APIM instance
	Plans map[string]string `json:"plans,omitempty"`
	// Last generation of the CRD resource
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

var _ custom.ApiDefinition = &ApiV4Definition{}
var _ custom.Status = &ApiDefinitionStatus{}
var _ custom.Spec = &ApiDefinitionV2Spec{}

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

// PickID returns the ID of the API definition, when a context has been defined at the spec level.
// The ID might be returned from the API status, meaning that the API is already known.
// If the API is unknown, the ID is either given from the spec if given,
// or generated from the API UID and the context key to ensure uniqueness
// in case the API is replicated on a same APIM instance.
func (api *ApiV4Definition) PickID(mCtx *management.Context) string {
	if api.Status.ID != "" {
		return api.Status.ID
	}

	if api.Spec.ID != "" {
		return api.Spec.ID
	}

	if mCtx != nil {
		return uuid.FromStrings(api.PickCrossID(), mCtx.OrgId, mCtx.EnvId)
	}

	return string(api.UID)
}

func (api *ApiV4Definition) PickCrossID() string {
	if api.Status.CrossID != "" {
		return api.Status.CrossID
	}

	if api.Spec.CrossID != "" {
		return api.Spec.CrossID
	}

	namespacedName := api.GetNamespacedName()
	return uuid.FromStrings(namespacedName.String())
}

func (api *ApiV4Definition) PickPlanIDs() map[string]*v4.Plan {
	plans := make(map[string]*v4.Plan, len(api.Spec.Plans))
	for key, plan := range api.Spec.Plans {
		p := plan.DeepCopy()
		if id, ok := api.Status.Plans[key]; ok {
			p.Id = id
		} else if plan.Id == "" {
			namespacedName := api.GetNamespacedName()
			p.Id = uuid.FromStrings(namespacedName.String(), key)
		}
		plans[key] = p
	}
	return plans
}

// GetOrGenerateEmptyPlanCrossID For each plan, generate a CrossId from Api Id & Plan Name if not defined.
func (api *ApiV4Definition) GetOrGenerateEmptyPlanCrossID() {
	for name, plan := range api.Spec.Plans {
		if plan.CrossId == "" {
			plan.CrossId = uuid.FromStrings(api.PickCrossID(), "/", name)
		}
	}
}

// EnvID implements custom.ApiDefinition.
func (api *ApiV4Definition) EnvID() string {
	return api.Status.EnvID
}

// ID implements custom.ApiDefinition.
func (api *ApiV4Definition) ID() string {
	return api.Status.ID
}

// OrgID implements custom.ApiDefinition.
func (api *ApiV4Definition) OrgID() string {
	return api.Status.OrgID
}

func (api *ApiV4Definition) GetSpec() custom.Spec {
	return &api.Spec
}

func (api *ApiV4Definition) GetStatus() custom.Status {
	return &api.Status
}

func (api *ApiV4Definition) DeepCopyResource() custom.Resource {
	return api.DeepCopy()
}

func (api *ApiV4Definition) ContextRef() custom.ResourceRef {
	return api.Spec.Context
}

func (api *ApiV4Definition) HasContext() bool {
	return api.Spec.Context != nil
}

func (api *ApiV4Definition) Version() custom.ApiDefinitionVersion {
	return custom.ApiV4
}

func (api *ApiV4Definition) GetNamespacedName() refs.NamespacedName {
	return refs.NamespacedName{Namespace: api.Namespace, Name: api.Name}
}

func (api *ApiV4Definition) GetObjectMeta() *metav1.ObjectMeta {
	return &api.ObjectMeta
}

func (spec *ApiV4DefinitionSpec) Hash() string {
	return hash.Calculate(spec)
}

func (spec *ApiV4DefinitionSpec) GetManagementContext() *refs.NamespacedName {
	return spec.Context
}

func (s *ApiV4DefinitionStatus) SetProcessingStatus(status custom.ProcessingStatus) {
	s.Status = status
}

func (s *ApiV4DefinitionStatus) SetObservedGeneration(g int64) {
	s.ObservedGeneration = g
}

func (s *ApiV4DefinitionStatus) DeepCopyFrom(api client.Object) error {
	switch t := api.(type) {
	case *ApiV4Definition:
		t.Status.DeepCopyInto(s)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *ApiV4DefinitionStatus) DeepCopyTo(api client.Object) error {
	switch t := api.(type) {
	case *ApiV4Definition:
		s.DeepCopyInto(&t.Status)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

// ApiV4DefinitionList contains a list of ApiV4Definition.
// +kubebuilder:object:root=true
type ApiV4DefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApiV4Definition `json:"items"`
}

func (l *ApiV4DefinitionList) GetItems() []list.Item {
	items := make([]list.Item, len(l.Items))
	for i := range l.Items {
		items[i] = &l.Items[i]
	}
	return items
}

func init() {
	SchemeBuilder.Register(&ApiV4Definition{}, &ApiV4DefinitionList{})
}
