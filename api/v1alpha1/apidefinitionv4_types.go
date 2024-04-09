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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/uuid"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/types/list"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ApiDefinitionV4Spec defines the desired state of ApiDefinition.
// +kubebuilder:object:generate=true
type ApiDefinitionV4Spec struct {
	v4.Api  `json:",inline"`
	Context *refs.NamespacedName `json:"contextRef,omitempty"`
}

// ApiDefinitionV4Status defines the observed state of API Definition.
type ApiDefinitionV4Status struct {
	// The organisation ID, if a management context has been defined to sync the API with an APIM instance
	OrgID string `json:"organizationId,omitempty"`
	// The environment ID, if a management context has been defined to sync the API with an APIM instance
	EnvID string `json:"environmentId,omitempty"`
	// The ID of the API definition if a management context has been defined to sync the API with an APIM instance
	ID string `json:"id,omitempty"`
	// The Cross ID of the API definition if a management context has been defined to sync the API with an APIM instance
	CrossID string `json:"crossId,omitempty"`

	// The processing status of the API definition.
	Status ProcessingStatus `json:"processingStatus,omitempty"`

	// The state of the API. Can be either STARTED or STOPPED.
	State string `json:"state,omitempty"`

	// This field is used to store the list of plans that have been created
	// for the API definition if a management context has been defined
	// to sync the API with an APIM instance
	Plans map[string]string `json:"plans,omitempty"`
	// Last generation of the CRD resource
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// ApiDefinitionV4 is the Schema for the v4 apidefinitions API.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Entrypoint",type=string,JSONPath=`.spec.proxy.virtual_hosts[*].path`,description="API entrypoint."
// +kubebuilder:printcolumn:name="Endpoint",type=string,JSONPath=`.spec.proxy.groups[*].endpoints[*].target`,description="API endpoint."
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`,description="API version."
// +kubebuilder:resource:shortName=graviteev4apis
// +kubebuilder:storageversion
type ApiDefinitionV4 struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApiDefinitionV4Spec   `json:"spec,omitempty"`
	Status ApiDefinitionV4Status `json:"status,omitempty"`
}

func (api *ApiDefinitionV4) ToGatewayDefinition() v4.GatewayDefinitionApi {
	cp := api.DeepCopy()
	cp.Spec.ID = api.PickID()
	return cp.Spec.Api.ToGatewayDefinition()
}

func (api *ApiDefinitionV4) IsBeingDeleted() bool {
	return !api.ObjectMeta.DeletionTimestamp.IsZero()
}

// PickID returns the ID of the API definition, when a context has been defined at the spec level.
// The ID might be returned from the API status, meaning that the API is already known.
// If the API is unknown, the ID is either given from the spec if given,
// or generated from the API UID and the context key to ensure uniqueness
// in case the API is replicated on a same APIM instance.
func (api *ApiDefinitionV4) PickID() string {
	if api.Status.ID != "" {
		return api.Status.ID
	}

	if api.Spec.ID != "" {
		return api.Spec.ID
	}

	return string(api.UID)
}

func (api *ApiDefinitionV4) PickCrossID() string {
	if api.Status.CrossID != "" {
		return api.Status.CrossID
	}

	if api.Spec.CrossID != "" {
		return api.Spec.CrossID
	}

	namespacedName := api.GetNamespacedName()
	return uuid.FromStrings(namespacedName.String())
}

func (api *ApiDefinitionV4) PickPlanIDs() map[string]*v4.Plan {
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
func (api *ApiDefinitionV4) GetOrGenerateEmptyPlanCrossID() {
	for _, plan := range api.Spec.Plans {
		if plan.CrossId == "" {
			plan.CrossId = uuid.FromStrings(api.Spec.ID, "/", plan.Name)
		}
	}
}

func (api *ApiDefinitionV4) GetNamespacedName() refs.NamespacedName {
	return refs.NamespacedName{Namespace: api.Namespace, Name: api.Name}
}

func (spec *ApiDefinitionV4Spec) SetDefinitionContext() {
	spec.DefinitionContext = &v4.DefinitionContext{
		Origin:   v4.OriginKubernetes,
		SyncFrom: v4.OriginKubernetes,
	}
}

func (api *ApiDefinitionV4) DeepCopyCrd() CRD {
	return api.DeepCopy()
}

func (api *ApiDefinitionV4) GetSpec() Spec {
	return &api.Spec
}

func (api *ApiDefinitionV4) GetApiDefinitionSpec() ContextAwareSpec {
	return &api.Spec
}

func (spec *ApiDefinitionV4Spec) Hash() string {
	return hash.Calculate(spec)
}

func (spec *ApiDefinitionV4Spec) GetManagementContext() *refs.NamespacedName {
	return spec.Context
}

func (api *ApiDefinitionV4) GetStatus() Status {
	return &api.Status
}

func (s *ApiDefinitionV4Status) SetProcessingStatus(status ProcessingStatus) {
	s.Status = status
}

func (s *ApiDefinitionV4Status) SetObservedGeneration(g int64) {
	s.ObservedGeneration = g
}

func (s *ApiDefinitionV4Status) DeepCopyFrom(api client.Object) error {
	switch t := api.(type) {
	case *ApiDefinitionV4:
		t.Status.DeepCopyInto(s)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *ApiDefinitionV4Status) DeepCopyTo(api client.Object) error {
	switch t := api.(type) {
	case *ApiDefinitionV4:
		s.DeepCopyInto(&t.Status)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

// ApiDefinitionV4List contains a list of ApiDefinitionV4.
// +kubebuilder:object:root=true
type ApiDefinitionV4List struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApiDefinitionV4 `json:"items"`
}

func (l *ApiDefinitionV4List) GetItems() []list.Item {
	items := make([]list.Item, len(l.Items))
	for i := range l.Items {
		items[i] = &l.Items[i]
	}
	return items
}

func init() {
	SchemeBuilder.Register(&ApiDefinitionV4{}, &ApiDefinitionV4List{})
}
