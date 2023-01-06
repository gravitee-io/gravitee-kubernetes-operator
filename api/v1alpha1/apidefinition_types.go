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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/uuid"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kUtil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// The API definition is the main resource handled by the Kubernetes Operator
// Most of the configuration properties defined here are already documented
// in the APIM Console API Reference.
// See https://docs.gravitee.io/apim/3.x/apim_installguide_rest_apis_documentation.html
// +kubebuilder:object:generate=true
type ApiDefinitionSpec struct {
	model.Api `json:",inline"`

	// The contextRef refers to the namespace and the name of a ManagementContext used for
	// synchronizing API definitions with a Gravitee API Management instance.
	Contexts []model.NamespacedName `json:"contexts,omitempty"`
}

type StatusContext struct {
	OrgID string `json:"organizationId"`
	EnvID string `json:"environmentId"`
	// The ID of the API definition in the Gravitee API Management instance (if a management context has been configured).
	ID      string `json:"id"`
	CrossID string `json:"crossId"`
	// The processing status of the API definition.
	Status ProcessingStatus `json:"status,omitempty"`
	// The state of the API. Can be either STARTED or STOPPED.
	State string `json:"state,omitempty"`
}

// ApiDefinitionStatus defines the observed state of API Definition.
type ApiDefinitionStatus struct {
	Contexts           map[string]StatusContext `json:"contexts,omitempty"`
	ObservedGeneration int64                    `json:"observedGeneration,omitempty"`
}

func (s *ApiDefinitionStatus) Initialize() {
	if s.Contexts == nil {
		s.Contexts = make(map[string]StatusContext)
	}
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Entrypoint",type=string,JSONPath=`.spec.proxy.virtual_hosts[*].path`,description="API entrypoint."
// +kubebuilder:printcolumn:name="Endpoint",type=string,JSONPath=`.spec.proxy.groups[*].endpoints[*].target`,description="API endpoint."
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`,description="API version."
// +kubebuilder:resource:shortName=graviteeapis
// ApiDefinition is the Schema for the apidefinitions API.
type ApiDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApiDefinitionSpec   `json:"spec,omitempty"`
	Status ApiDefinitionStatus `json:"status,omitempty"`
}

// +kubebuilder:validation:Enum=Completed;Failed;
type ProcessingStatus string

const (
	ProcessingStatusCompleted ProcessingStatus = "Completed"
	ProcessingStatusFailed    ProcessingStatus = "Failed"
)

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
func (api *ApiDefinition) PickID(statusKey string) string {
	status, ok := api.Status.Contexts[statusKey]

	if ok && status.ID != "" {
		return status.ID
	}

	if api.Spec.ID != "" {
		return api.Spec.ID
	}

	return uuid.FromStrings(string(api.UID), statusKey)
}

// GetID returns the ID of the API definition.
// This method should not be used when a context has been defined at the spec level
// and is only used when generating the config map definition, because the gateway
// required an ID to be defined on deserialization.
func (api *ApiDefinition) GetID() string {
	if api.Spec.ID != "" {
		return api.Spec.ID
	}

	return string(api.UID)
}

func (api *ApiDefinition) PickCrossID(statusKey string) string {
	status, ok := api.Status.Contexts[statusKey]

	if ok && status.CrossID != "" {
		return status.CrossID
	}

	return api.GetOrGenerateCrossID()
}

func (api *ApiDefinition) GetOrGenerateCrossID() string {
	if api.Spec.CrossID != "" {
		return api.Spec.CrossID
	}

	return uuid.FromStrings(api.GetNamespacedName().String())
}

func (api *ApiDefinition) GetNamespacedName() model.NamespacedName {
	return model.NamespacedName{Namespace: api.Namespace, Name: api.Name}
}

func (spec *ApiDefinitionSpec) SetDefinitionContext() {
	spec.DefinitionContext = &model.DefinitionContext{
		Mode:   model.ModeFullyManaged,
		Origin: model.OriginKubernetes,
	}
}

// +kubebuilder:object:root=true
// ApiDefinitionList contains a list of ApiDefinition.
type ApiDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApiDefinition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApiDefinition{}, &ApiDefinitionList{})
}
