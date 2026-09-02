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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/am"
	gcontext "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/context"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.ContextObject = &AMContext{}
var _ core.ContextModel = &AMContext{}
var _ core.Spec = &AMContextSpec{}
var _ core.Status = &AMContextStatus{}
var _ core.ConditionAware = &AMContext{}

// AMContextSpec is a cloud-free, bearer-only copy of ManagementContext.
// +kubebuilder:object:generate=true
type AMContextSpec struct {
	*am.Context `json:",inline"`
}

func (spec *AMContextSpec) Hash() string {
	return hash.Calculate(spec)
}

// AMContextStatus defines the observed state of an AM Context.
type AMContextStatus struct {
	gcontext.Status `json:",inline"`
}

func (s *AMContextStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *AMContext:
		t.Status.DeepCopyInto(s)
		return nil
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

func (s *AMContextStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *AMContext:
		s.DeepCopyInto(&t.Status)
		return nil
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

func (s *AMContextStatus) SetProcessingStatus(_ core.ProcessingStatus) {
	// Not implemented
}

func (s *AMContextStatus) IsFailed() bool {
	return false
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="BaseUrl",type=string,JSONPath=`.spec.baseUrl`
// +kubebuilder:resource:shortName=amcontexts
type AMContext struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AMContextSpec   `json:"spec,omitempty"`
	Status AMContextStatus `json:"status,omitempty"`
}

func (ctx *AMContext) DeepCopyResource() core.Object {
	return ctx.DeepCopy()
}

func (ctx *AMContext) GetSpec() core.Spec {
	return &ctx.Spec
}

func (ctx *AMContext) GetStatus() core.Status {
	return &ctx.Status
}

func (ctx *AMContext) GetAuth() core.Auth {
	return ctx.Spec.Context.Auth
}

func (ctx *AMContext) GetEnvID() string {
	return ctx.Spec.EnvID
}

func (ctx *AMContext) GetOrgID() string {
	return ctx.Spec.OrgID
}

func (ctx *AMContext) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      ctx.Name,
		Namespace: ctx.Namespace,
	}
}

func (ctx *AMContext) GetSecretRef() core.ObjectRef {
	return ctx.Spec.SecretRef()
}

func (ctx *AMContext) GetURL() string {
	return ctx.Spec.BaseUrl
}

func (ctx *AMContext) GetPath() *string {
	return ctx.Spec.Path
}

func (ctx *AMContext) HasAuthentication() bool {
	return ctx.Spec.Auth != nil
}

func (ctx *AMContext) HasSecretRef() bool {
	return ctx.HasAuthentication() && ctx.Spec.Auth.SecretRef != nil
}

func (ctx *AMContext) GetNamespacedName() *refs.NamespacedName {
	return &refs.NamespacedName{Namespace: ctx.Namespace, Name: ctx.Name}
}

func (ctx *AMContext) GetContext() core.ContextModel {
	return ctx.Spec.Context
}

// +kubebuilder:object:root=true
// AMContextList contains a list of AM Contexts.
type AMContextList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AMContext `json:"items"`
}

func (ctx *AMContext) IsBeingDeleted() bool {
	return !ctx.ObjectMeta.DeletionTimestamp.IsZero()
}

func (ctx *AMContext) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(ctx.Status.Conditions)
}

func (ctx *AMContext) SetConditions(conditions []metav1.Condition) {
	ctx.Status.Conditions = conditions
}
