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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.ApplicationObject = &Application{}
var _ core.Spec = &ApplicationSpec{}

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

// GetSpec implements custom.Resource.
func (app *Application) GetSpec() core.Spec {
	return &app.Spec
}

// GetStatus implements custom.Resource.
func (app *Application) GetStatus() core.Status {
	return &app.Status
}

func (app *Application) ContextRef() core.ObjectRef {
	return app.Spec.Context
}

func (app *Application) HasContext() bool {
	return app.Spec.Context != nil
}

func (app *Application) PopulateIDs(mCtx core.ContextModel) {
	if app.Status.ID != "" {
		app.Spec.ID = app.Status.ID
	}

	if app.Spec.ID == "" {
		app.Spec.ID = string(app.UID)
	}
}

func (app *Application) GetID() string {
	return app.Status.ID
}

func (app *Application) GetOrgID() string {
	return app.Status.OrgID
}

func (app *Application) GetModel() core.ApplicationModel {
	return &app.Spec.Application
}

func (app *Application) GetEnvID() string {
	return app.Status.EnvID
}

func (app *Application) DeepCopyResource() core.Object {
	return app.DeepCopy()
}

func (app *Application) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      app.Name,
		Namespace: app.Namespace,
	}
}

func (spec *ApplicationSpec) Hash() string {
	return hash.Calculate(spec)
}

func (s *ApplicationStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *Application:
		t.Status.DeepCopyInto(s)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *ApplicationStatus) DeepCopyTo(api client.Object) error {
	switch t := api.(type) {
	case *Application:
		s.DeepCopyInto(&t.Status)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *ApplicationStatus) SetProcessingStatus(status core.ProcessingStatus) {
	s.Status.ProcessingStatus = status
}

func (s *ApplicationStatus) IsFailed() bool {
	return s.ProcessingStatus == core.ProcessingStatusFailed
}
