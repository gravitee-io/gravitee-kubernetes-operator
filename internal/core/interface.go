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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	IsBeingDeleted() bool
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
type PlanModel interface {
	GetSecurityType() string
	GetID() string
}

// +k8s:deepcopy-gen=false
type ApiDefinitionModel interface {
	GetName() string
	GetDefinitionVersion() ApiDefinitionVersion
	GetContextPaths() []string
	SetDefinitionContext(DefinitionContext)
	GetDefinitionContext() DefinitionContext
	GetResources() []ObjectOrRef[ResourceModel]
	GetState() string
	HasPlans() bool
	GetPlan(string) PlanModel
	IsStopped() bool
	GetType() string
	GetGroups() []string
	SetGroups([]string)
	GetGroupRefs() []ObjectRef
	GetNotificationRefs() []ObjectRef
	SetConsoleNotification(ConsoleNotificationSettingsObject)
	GetTags() []string
}

// +k8s:deepcopy-gen=false
type ConsoleNotificationSettingsObject interface {
	IsConsoleNotification() bool
}

// +k8s:deepcopy-gen=false
type ApiDefinitionObject interface {
	ContextAwareObject
	ApiDefinitionModel
	GetDefinition() ApiDefinitionModel
	SetDefinitionContext(DefinitionContext)
	GetDefinitionContext() DefinitionContext
	IsSyncFromManagement() bool
}

// +k8s:deepcopy-gen=false
type ApplicationSettings interface {
	IsOAuth() bool
	GetOAuthType() string
	IsSimple() bool
	HasTLS() bool
	GetClientID() string
	GetClientCertificate() string
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
	IsFailed() bool
	DeepCopyFrom(obj client.Object) error
	DeepCopyTo(obj client.Object) error
}

type SubscribableStatus interface {
	Status
	GetSubscriptionCount() uint
	AddSubscription()
	RemoveSubscription()
}

// +k8s:deepcopy-gen=false
type ContextAwareObject interface {
	Object
	ContextRef() ObjectRef
	HasContext() bool
	GetID() string
	PopulateIDs(context ContextModel)
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
	GetPath() *string
	GetEnvID() string
	GetOrgID() string
	HasAuthentication() bool
	GetAuth() Auth
	HasCloud() bool
	GetCloud() Cloud
	ConfigureCloud(url string, orgID string, envID string)
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
	SetSecretRef(ref ObjectRef)
}

// +k8s:deepcopy-gen=false
type BasicAuth interface {
	GetUsername() string
	GetPassword() string
}

// +k8s:deepcopy-gen=false
type Cloud interface {
	SecretAware
	GetToken() string
	IsEnabled() bool
}

// +k8s:deepcopy-gen=false
type ObjectRef interface {
	fmt.Stringer
	NamespacedName() types.NamespacedName
	GetName() string
	GetNamespace() string
	GetKind() string
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

type SubscriptionObject interface {
	Object
	SubscriptionModel
}

type SubscriptionModel interface {
	GetAppRef() ObjectRef
	GetApiRef() ObjectRef
	SetApiKind(string)
	GetPlan() string
	GetEndingAt() *string
}

type ConditionAware interface {
	GetConditions() map[string]metav1.Condition
	SetConditions([]metav1.Condition)
}
