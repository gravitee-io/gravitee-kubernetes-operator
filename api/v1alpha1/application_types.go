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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Application is the main resource handled by the Kubernetes Operator
// +kubebuilder:object:generate=true
type ApplicationSpec struct {
	application.Application `json:",inline"`
	// +kubebuilder:validation:Required
	Context *refs.NamespacedName `json:"contextRef"`
}

// ApplicationStatus defines the observed state of Application.
type ApplicationStatus struct {
	application.Status `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.applicationType`
// +kubebuilder:resource:shortName=graviteeapplications
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

func (app *Application) IsBeingDeleted() bool {
	return !app.ObjectMeta.DeletionTimestamp.IsZero()
}

func init() {
	SchemeBuilder.Register(&Application{}, &ApplicationList{})
}
