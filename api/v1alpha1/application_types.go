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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kUtil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// +kubebuilder:object:generate=true
type ApplicationSpec struct {
	application.Application `json:",inline"`
	// The contextRef field allows to reference a previously created ManagementContext.
	// See <a href="#managementcontext">ManagementContext</a> for more information.
	// +kubebuilder:validation:Required
	Context *refs.NamespacedName `json:"contextRef,omitempty"`
}

// ApplicationStatus defines the observed state of Application.
type ApplicationStatus struct {
	// The Organization ID of the Application in the Gravitee API Management instance.
	OrgID string `json:"organizationId,omitempty"`
	// The Environment ID of the Application in the Gravitee API Management instance.
	EnvID string `json:"environmentId,omitempty"`
	// The ID of the Application in the Gravitee API Management instance.
	ID string `json:"id,omitempty"`
	// The processing status of the Application.
	Status ProcessingStatus `json:"processingStatus,omitempty"`

	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.applicationType`
// +kubebuilder:resource:shortName=graviteeapplications
// Applications resource allows to use the operator to create and manage applications in Gravitee API Management.
// This means that a management context must be define in the specification when creating an application.
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec   `json:"spec,omitempty"`
	Status ApplicationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

func (api *Application) IsMissingDeletionFinalizer() bool {
	return !kUtil.ContainsFinalizer(api, keys.ApplicationDeletionFinalizer)
}

func (api *Application) IsBeingDeleted() bool {
	return !api.ObjectMeta.DeletionTimestamp.IsZero()
}

func init() {
	SchemeBuilder.Register(&Application{}, &ApplicationList{})
}
