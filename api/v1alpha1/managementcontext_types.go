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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/management"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/types/k8s/custom"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ custom.ContextResource = &ManagementContext{}
var _ custom.Spec = &ManagementContextSpec{}
var _ custom.Status = &ManagementContextStatus{}

// ManagementContext represents the configuration for a specific environment
// +kubebuilder:object:generate=true
type ManagementContextSpec struct {
	*management.Context `json:",inline"`
}

// Hash implements custom.Spec.
func (spec *ManagementContextSpec) Hash() string {
	return hash.Calculate(spec)
}

// ManagementContextStatus defines the observed state of an API Context.
type ManagementContextStatus struct {
}

// DeepCopyFrom implements custom.Status.
func (st *ManagementContextStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *ManagementContext:
		t.Status.DeepCopyInto(st)
		return nil
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

// DeepCopyTo implements custom.Status.
func (st *ManagementContextStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *ManagementContext:
		st.DeepCopyInto(&t.Status)
		return nil
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

// SetObservedGeneration implements custom.Status.
func (st *ManagementContextStatus) SetObservedGeneration(g int64) {
	// Not implemented
}

// SetProcessingStatus implements custom.Status.
func (st *ManagementContextStatus) SetProcessingStatus(status custom.ProcessingStatus) {
	// Not implemented
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="BaseUrl",type=string,JSONPath=`.spec.baseUrl`
// +kubebuilder:resource:shortName=graviteecontexts
type ManagementContext struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManagementContextSpec   `json:"spec,omitempty"`
	Status ManagementContextStatus `json:"status,omitempty"`
}

// DeepCopyResource implements custom.Context.
func (ctx *ManagementContext) DeepCopyResource() custom.Resource {
	return ctx.DeepCopy()
}

// GetSpec implements custom.Context.
func (ctx *ManagementContext) GetSpec() custom.Spec {
	return &ctx.Spec
}

// GetStatus implements custom.Context.
func (ctx *ManagementContext) GetStatus() custom.Status {
	return &ctx.Status
}

// GetAuth implements custom.Context.
func (ctx *ManagementContext) GetAuth() custom.Auth {
	return ctx.Spec.Context.Auth
}

// GetEnv implements custom.Context.
func (ctx *ManagementContext) GetEnv() string {
	return ctx.Spec.EnvId
}

// GetOrg implements custom.Context.
func (ctx *ManagementContext) GetOrg() string {
	return ctx.Spec.OrgId
}

// GetSecretRef implements custom.Context.
func (ctx *ManagementContext) GetSecretRef() custom.ResourceRef {
	return ctx.Spec.SecretRef()
}

// GetURL implements custom.Context.
func (ctx *ManagementContext) GetURL() string {
	return ctx.Spec.BaseUrl
}

// HasAuthentication implements custom.Context.
func (ctx *ManagementContext) HasAuthentication() bool {
	return ctx.Spec.Auth != nil
}

// HasSecretRef implements custom.Context.
func (ctx *ManagementContext) HasSecretRef() bool {
	return ctx.HasAuthentication() && ctx.Spec.Auth.SecretRef != nil
}

func (ctx *ManagementContext) GetNamespacedName() *refs.NamespacedName {
	return &refs.NamespacedName{Namespace: ctx.Namespace, Name: ctx.Name}
}

// +kubebuilder:object:root=true
// ManagementContextList contains a list of API Contexts.
type ManagementContextList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ManagementContext `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ManagementContext{}, &ManagementContextList{})
}

func (ctx *ManagementContext) IsBeingDeleted() bool {
	return !ctx.ObjectMeta.DeletionTimestamp.IsZero()
}
