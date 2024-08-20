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

package core

import (
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
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
type Object interface {
	client.Object
	GetSpec() Spec
	GetStatus() Status
	GetRef() ObjectRef
}

// +k8s:deepcopy-gen=false
type ObjectOrRef[T any] interface {
	IsRef() bool
	GetRef() ObjectRef
	GetObject() T
	SetObject(obj T)
}

// +k8s:deepcopy-gen=false
type DefinitionContext interface {
	GetOrigin() string
	SetOrigin(string)
}

// +k8s:deepcopy-gen=false
type ApiDefinitionModel interface {
	GetDefinitionVersion() ApiDefinitionVersion
	GetContextPaths() []string
	SetDefinitionContext(DefinitionContext)
	GetDefinitionContext() DefinitionContext
	GetResources() []ObjectOrRef[ResourceModel]
}

// +k8s:deepcopy-gen=false
type ApiDefinitionObject interface {
	ContextAwareObject
	ApiDefinitionModel
	GetDefinition() ApiDefinitionModel
	PopulateIDs(context ContextModel)
	SetDefinitionContext(DefinitionContext)
	GetDefinitionContext() DefinitionContext
}

// +k8s:deepcopy-gen=false
type ApplicationSettings interface {
	IsOAuth() bool
}

// +k8s:deepcopy-gen=false
type ApplicationModel interface {
	GetSettings() ApplicationSettings
}

// +k8s:deepcopy-gen=false
type ApplicationObject interface {
	ContextAwareObject
	GetModel() ApplicationModel
}

// +k8s:deepcopy-gen=false
type Spec interface {
	Hash() string
}

// +k8s:deepcopy-gen=false
type Status interface {
	SetProcessingStatus(status ProcessingStatus)
	DeepCopyFrom(obj client.Object) error
	DeepCopyTo(obj client.Object) error
}

// +k8s:deepcopy-gen=false
type ContextAwareObject interface {
	Object
	ContextRef() ObjectRef
	HasContext() bool
	GetID() string
	GetOrgID() string
	GetEnvID() string
}

// +k8s:deepcopy-gen=false
type SecretAware interface {
	GetSecretRef() ObjectRef
	HasSecretRef() bool
}

// +k8s:deepcopy-gen=false
type ContextModel interface {
	SecretAware
	GetURL() string
	GetEnvID() string
	GetOrgID() string
	HasAuthentication() bool
	GetAuth() Auth
}

// +k8s:deepcopy-gen=false
type ContextObject interface {
	ContextModel
	Object
}

// +k8s:deepcopy-gen=false
type Auth interface {
	GetBearerToken() string
	HasCredentials() bool
	GetCredentials() BasicAuth
	GetSecretRef() ObjectRef
	SetCredentials(username, password string)
	SetToken(token string)
}

// +k8s:deepcopy-gen=false
type BasicAuth interface {
	GetUsername() string
	GetPassword() string
}

// +k8s:deepcopy-gen=false
type ObjectRef interface {
	fmt.Stringer
	NamespacedName() types.NamespacedName
	GetName() string
	GetNamespace() string
	HasNameSpace() bool
	IsMissingNamespace() bool
	SetNamespace(ns string)
}

// +k8s:deepcopy-gen=false
type ResourceObject interface {
	ResourceModel
	Object
}

// +k8s:deepcopy-gen=false
type ResourceModel interface {
	GetType() string
	GetResourceName() string
	GetConfig() *utils.GenericStringMap
}
