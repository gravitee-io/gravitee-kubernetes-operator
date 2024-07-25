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

package custom

import (
	"fmt"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ProcessingStatus string

const (
	ProcessingStatusCompleted ProcessingStatus = "Completed"
	ProcessingStatusFailed    ProcessingStatus = "Failed"
)

type ApiDefinitionVersion string

const (
	ApiV2 = ApiDefinitionVersion("V2")
	ApiV4 = ApiDefinitionVersion("V4")
)

// +k8s:deepcopy-gen=false
type Resource interface {
	client.Object
	GetSpec() Spec
	GetStatus() Status
	DeepCopyResource() Resource
}

// +k8s:deepcopy-gen=false
type ApiDefinition interface {
	ContextAwareResource
	Version() ApiDefinitionVersion
	OrgID() string
	EnvID() string
}

// +k8s:deepcopy-gen=false
type Spec interface {
	Hash() string
}

// +k8s:deepcopy-gen=false
type Status interface {
	SetProcessingStatus(status ProcessingStatus)
	SetObservedGeneration(g int64)
	DeepCopyFrom(obj client.Object) error
	DeepCopyTo(obj client.Object) error
}

// +k8s:deepcopy-gen=false
type ContextAwareResource interface {
	Resource
	ContextRef() ResourceRef
	HasContext() bool
	ID() string
}

// +k8s:deepcopy-gen=false
type ResourceRef interface {
	fmt.Stringer
	NamespacedName() types.NamespacedName
	GetName() string
	GetNamespace() string
	HasNameSpace() bool
	IsMissingNamespace() bool
	SetNamespace(ns string)
}
