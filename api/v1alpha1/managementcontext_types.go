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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.ContextObject = &ManagementContext{}
var _ core.ContextModel = &ManagementContext{}
var _ core.Spec = &ManagementContextSpec{}
var _ core.Status = &ManagementContextStatus{}

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
func (s *ManagementContextStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *ManagementContext:
		t.Status.DeepCopyInto(s)
		return nil
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

// DeepCopyTo implements custom.Status.
func (s *ManagementContextStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *ManagementContext:
		s.DeepCopyInto(&t.Status)
		return nil
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

// SetProcessingStatus implements custom.Status.
func (s *ManagementContextStatus) SetProcessingStatus(status core.ProcessingStatus) {
	// Not implemented
}

func (s *ManagementContextStatus) IsFailed() bool {
	return false
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
func (ctx *ManagementContext) DeepCopyResource() core.Object {
	return ctx.DeepCopy()
}

// GetSpec implements custom.Context.
func (ctx *ManagementContext) GetSpec() core.Spec {
	return &ctx.Spec
}

// GetStatus implements custom.Context.
func (ctx *ManagementContext) GetStatus() core.Status {
	return &ctx.Status
}

// GetAuth implements custom.Context.
func (ctx *ManagementContext) GetAuth() core.Auth {
	return ctx.Spec.Context.Auth
}

// GetEnvID implements custom.Context.
func (ctx *ManagementContext) GetEnvID() string {
	return ctx.Spec.EnvID
}

// GetOrgID implements custom.Context.
func (ctx *ManagementContext) GetOrgID() string {
	return ctx.Spec.OrgID
}

func (ctx *ManagementContext) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      ctx.Name,
		Namespace: ctx.Namespace,
	}
}

// GetSecretRef implements custom.Context.
func (ctx *ManagementContext) GetSecretRef() core.ObjectRef {
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

func (ctx *ManagementContext) HasCloud() bool {
	return ctx.Spec.HasCloud()
}

func (ctx *ManagementContext) GetCloud() core.Cloud {
	return ctx.Spec.Cloud
}

func (ctx *ManagementContext) ConfigureCloud(url string, orgID string, envID string) {
	ctx.Spec.ConfigureCloud(url, orgID, envID)
}

func (ctx *ManagementContext) GetContext() core.ContextModel {
	return ctx.Spec.Context
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
