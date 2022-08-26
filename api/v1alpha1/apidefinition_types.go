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
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// APIDefinition represents the configuration for a single proxied API and it's versions.
// +kubebuilder:object:generate=true
type ApiDefinitionSpec struct {
	model.Api `json:",inline"`

	// The context is specifying the namespace and the name of a ManagementContext used for
	// managing the APIDefinition from the ManagementAPI
	Context *model.ContextRef `json:"contextRef,omitempty"`
}

// ApiDefinitionStatus defines the observed state of ApiDefinition.
type ApiDefinitionStatus struct {
	ID         string      `json:"id"`
	CrossID    string      `json:"crossId"`
	State      model.State `json:"state"`
	Generation int64       `json:"generation"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`,description="API state (STARTED or STOPPED)."
// +kubebuilder:printcolumn:name="Entrypoint",type=string,JSONPath=`.spec.proxy.virtual_hosts[*].path`,description="API entrypoint."
// +kubebuilder:printcolumn:name="Endpoint",type=string,JSONPath=`.spec.proxy.groups[*].endpoints[*].target`,description="API endpoint."
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`,description="API version."
// +kubebuilder:printcolumn:name="ManagementContext",type=string,JSONPath=`.spec.contextRef.name`,description="Management context name."
// +kubebuilder:resource:shortName=graviteeapis
// ApiDefinition is the Schema for the apidefinitions API.
type ApiDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApiDefinitionSpec   `json:"spec,omitempty"`
	Status ApiDefinitionStatus `json:"status,omitempty"`
}

func (api *ApiDefinition) HasFinalizer() bool {
	return util.ContainsFinalizer(api, keys.ApiDefinitionDeletionFinalizer)
}

func (api *ApiDefinition) IsBeingDeleted() bool {
	return !api.ObjectMeta.DeletionTimestamp.IsZero()
}

func (api *ApiDefinition) IsBeingUpdated() bool {
	return api.Status.Generation != api.ObjectMeta.Generation
}

func (api *ApiDefinition) IsBeingCreated() bool {
	return api.Status.CrossID == ""
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
