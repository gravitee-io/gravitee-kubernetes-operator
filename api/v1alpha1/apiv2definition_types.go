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
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"

	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/uuid"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/types/k8s/custom"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/types/list"
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
	// +kubebuilder:default:=true
	IsLocal bool `json:"local"`
}

// ApiDefinitionStatus defines the observed state of API Definition.
type ApiDefinitionStatus struct {
	OrgID string `json:"organizationId,omitempty"`
	EnvID string `json:"environmentId,omitempty"`
	// The ID of the API definition in the Gravitee API Management instance (if an API context has been configured).
	ID      string `json:"id,omitempty"`
	CrossID string `json:"crossId,omitempty"`
	// The processing status of the API definition.
	Status custom.ProcessingStatus `json:"processingStatus,omitempty"`
	// This field is kept for backward compatibility and shall be removed in future versions.
	// Use processingStatus instead.
	DeprecatedStatus custom.ProcessingStatus `json:"status,omitempty"`

	// The state of the API. Can be either STARTED or STOPPED.
	State base.ApiState `json:"state,omitempty"`

	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// This field is kept for backward compatibility and shall be removed in future versions.
	// Use observedGeneration instead.
	DeprecatedObservedGeneration int64 `json:"generation,omitempty"`
}

var _ list.Item = &ApiDefinition{}
var _ custom.ApiDefinition = &ApiDefinition{}
var _ custom.Status = &ApiDefinitionStatus{}
var _ custom.Spec = &ApiDefinitionV2Spec{}

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

	return api.GetOrGenerateCrossID()
}

func (api *ApiDefinition) GetOrGenerateCrossID() string {
	if api.Spec.CrossID != "" {
		return api.Spec.CrossID
	}

	return uuid.FromStrings(api.GetNamespacedName().String())
}

func (api *ApiDefinition) GetNamespacedName() refs.NamespacedName {
	return refs.NamespacedName{Namespace: api.Namespace, Name: api.Name}
}

func (spec *ApiDefinitionV2Spec) SetDefinitionContext() {
	spec.DefinitionContext = &v2.DefinitionContext{
		Mode:     v2.ModeFullyManaged,
		Origin:   v2.OriginKubernetes,
		SyncFrom: v2.OriginKubernetes,
	}
	if !spec.IsLocal || strings.EqualFold(string(v4.OriginManagement), spec.DefinitionContext.SyncFrom) {
		spec.DefinitionContext.SyncFrom = string(v4.OriginManagement)
	}
}

func (api *ApiDefinition) EnvID() string {
	return api.Status.EnvID
}

func (api *ApiDefinition) ID() string {
	return api.Status.ID
}

func (api *ApiDefinition) OrgID() string {
	return api.Status.OrgID
}

func (api *ApiDefinition) Version() custom.ApiDefinitionVersion {
	return custom.ApiV2
}

func (api *ApiDefinition) GetSpec() custom.Spec {
	return &api.Spec
}

func (api *ApiDefinition) GetStatus() custom.Status {
	return &api.Status
}

func (api *ApiDefinition) DeepCopyResource() custom.Resource {
	return api.DeepCopy()
}

func (api *ApiDefinition) ContextRef() custom.ResourceRef {
	return api.Spec.Context
}

func (api *ApiDefinition) HasContext() bool {
	return api.Spec.Context != nil
}

func (spec *ApiDefinitionV2Spec) Hash() string {
	return hash.Calculate(spec)
}

func (s *ApiDefinitionStatus) SetProcessingStatus(status custom.ProcessingStatus) {
	s.Status = status
}

func (s *ApiDefinitionStatus) SetObservedGeneration(g int64) {
	s.ObservedGeneration = g
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

var _ list.Interface = &ApiDefinitionList{}

// ApiDefinitionList contains a list of ApiDefinition.
// +kubebuilder:object:root=true
type ApiDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApiDefinition `json:"items"`
}

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
