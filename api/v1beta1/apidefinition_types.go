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

package v1beta1

import (
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/uuid"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/types/list"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kUtil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (*ApiDefinition) Hub() {}

// ApiDefinitionSpec defines the desired state of ApiDefinition.
type ApiDefinitionSpec struct {
	v4.Api  `json:",inline"`
	Context *refs.NamespacedName `json:"contextRef,omitempty"`
}

// +kubebuilder:validation:Enum=Completed;Failed;
type ProcessingStatus string

const (
	ProcessingStatusCompleted ProcessingStatus = "Completed"
	ProcessingStatusFailed    ProcessingStatus = "Failed"
)

// ApiDefinitionStatus defines the observed state of API Definition.
type ApiDefinitionStatus struct {
	OrgID string `json:"organizationId,omitempty"`

	EnvID string `json:"environmentId,omitempty"`

	// The ID of the API definition in the Gravitee API Management instance (if an API context has been configured).
	ID string `json:"id,omitempty"`

	CrossID string `json:"crossId,omitempty"`

	// The processing status of the API definition.
	Status ProcessingStatus `json:"processingStatus,omitempty"`

	// The state of the API. Can be either STARTED or STOPPED.
	State string `json:"state,omitempty"`

	// This field is used to store the list of plans that have been created for the API definition.
	// Especially when the API is synced with an APIM instance
	Plans map[string]string `json:"plans,omitempty"`

	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// ApiDefinition is the Schema for the apidefinitions API.
// The v1beta1 API version is compatible with APIM 4.x features.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
type ApiDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApiDefinitionSpec   `json:"spec,omitempty"`
	Status ApiDefinitionStatus `json:"status,omitempty"`
}

func (api *ApiDefinition) ToGatewayDefinition() v4.GatewayDefinitionApi {
	cp := api.DeepCopy()
	cp.Spec.ID = api.PickID()
	return cp.Spec.Api.ToGatewayDefinition()
}

//+kubebuilder:object:root=true

// ApiDefinitionList contains a list of ApiDefinition.
type ApiDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApiDefinition `json:"items"`
}

func (api *ApiDefinition) IsMissingDeletionFinalizer() bool {
	return !kUtil.ContainsFinalizer(api, keys.ApiDefinitionDeletionFinalizer)
}

func (api *ApiDefinition) IsBeingDeleted() bool {
	return !api.ObjectMeta.DeletionTimestamp.IsZero()
}

// PickID returns the ID of the API definition, when a context has been defined at the spec level.
// The ID might be returned from the API status, meaning that the API is already known.
// If the API is unknown, the ID is either given from the spec if given,
// or generated from the API UID and the context key to ensure uniqueness
// in case the API is replicated on a same APIM instance.
func (api *ApiDefinition) PickID() string {
	if api.Status.ID != "" {
		return api.Status.ID
	}

	if api.Spec.ID != "" {
		return api.Spec.ID
	}

	return string(api.UID)
}

func (api *ApiDefinition) PickCrossID() string {
	if api.Status.CrossID != "" {
		return api.Status.CrossID
	}

	if api.Spec.CrossID != "" {
		return api.Spec.CrossID
	}

	return uuid.FromStrings(api.GetNamespacedName().String())
}

func (api *ApiDefinition) PickPlanIDs() map[string]*v4.Plan {
	plans := make(map[string]*v4.Plan, len(api.Spec.Plans))
	for key, plan := range api.Spec.Plans {
		p := plan.DeepCopy()
		if id, ok := api.Status.Plans[key]; ok {
			p.ID = id
		} else if plan.ID == "" {
			p.ID = uuid.FromStrings(api.GetNamespacedName().String(), key)
		}
		plans[key] = p
	}
	return plans
}

func (api *ApiDefinition) GetNamespacedName() refs.NamespacedName {
	return refs.NamespacedName{Namespace: api.Namespace, Name: api.Name}
}

var _ list.Interface = &ApiDefinitionList{}

func (l *ApiDefinitionList) GetItems() []list.Item {
	items := make([]list.Item, len(l.Items))
	for i := range l.Items {
		items[i] = &l.Items[i]
	}
	return items
}

func init() {
	SchemeBuilder.Register(&ApiDefinition{}, &ApiDefinitionList{})
}
