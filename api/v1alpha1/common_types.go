// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CRD is a common interface that can be used for generic use cases
// +k8s:deepcopy-gen=false
type CRD interface {
	client.Object
	schema.ObjectKind
	GetObjectMeta() *metav1.ObjectMeta
	GetSpec() Spec
	GetStatus() Status
	DeepCopyCrd() CRD
}

// ApiDefinitionCRD
// +k8s:deepcopy-gen=false
type ApiDefinitionCRD interface {
	CRD
	GetApiDefinitionSpec() ContextAwareSpec
}

// Spec
// +k8s:deepcopy-gen=false
type Spec interface {
	Hash() string
}

// ContextAwareSpec
// +k8s:deepcopy-gen=false
type ContextAwareSpec interface {
	Spec
	GetManagementContext() *refs.NamespacedName
}

// Status is a common interface that can be used for generic use cases
// +k8s:deepcopy-gen=false
type Status interface {
	SetProcessingStatus(status ProcessingStatus)
	SetObservedGeneration(g int64)
	DeepCopyFrom(obj client.Object) error
	DeepCopyTo(obj client.Object) error
}
