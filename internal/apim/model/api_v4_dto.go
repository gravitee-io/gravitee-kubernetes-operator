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

package model

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	nav "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/navigation"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
)

type APIV4DTO struct {
	ID                               string                                          `json:"id,omitempty" drift:"ignore"`
	HRID                             string                                          `json:"hrid,omitempty"`
	CrossID                          string                                          `json:"crossId,omitempty" drift:"ignore"`
	Name                             string                                          `json:"name"`
	Version                          string                                          `json:"version"`
	State                            base.ApiState                                   `json:"state,omitempty"`
	Tags                             []string                                        `json:"tags" drift:"empty-is-nil"`
	Labels                           []string                                        `json:"labels" drift:"empty-is-nil"`
	Visibility                       base.ApiVisibility                              `json:"visibility,omitempty"`
	Properties                       []*APIV4PropertyDTO                             `json:"properties" drift:"empty-is-nil"`
	Metadata                         []*APIV4MetadataEntryDTO                        `json:"metadata" drift:"ignore-remote-only-metadata"`
	Resources                        []*APIV4ResourceDTO                             `json:"resources" drift:"empty-is-nil"`
	Groups                           []string                                        `json:"groups" drift:"ignore-unknown-crd-groups"`
	Categories                       []string                                        `json:"categories" drift:"empty-is-nil"`
	NotifyMembers                    bool                                            `json:"notifyMembers" drift:"empty-is-true"`
	Description                      *string                                         `json:"description,omitempty"`
	DefinitionVersion                base.DefinitionVersion                          `json:"definitionVersion,omitempty" drift:"ignore"`
	DefinitionContext                *APIV4DefinitionContextDTO                      `json:"definitionContext,omitempty" drift:"ignore"`
	LifecycleState                   v4.ApiV4LifecycleState                          `json:"lifecycleState,omitempty"`
	Type                             v4.ApiType                                      `json:"type"`
	Listeners                        []*APIV4GenericListenerDTO                      `json:"listeners"`
	EndpointGroups                   []*APIV4EndpointGroupDTO                        `json:"endpointGroups"`
	FlowExecution                    *APIV4FlowExecutionDTO                          `json:"flowExecution,omitempty"`
	Flows                            []*APIV4FlowDTO                                 `json:"flows" drift:"empty-is-nil"`
	Analytics                        *APIV4AnalyticsDTO                              `json:"analytics,omitempty"`
	Services                         *APIV4ApiServicesDTO                            `json:"services,omitempty"`
	ResponseTemplates                map[string]map[string]*APIV4ResponseTemplateDTO `json:"responseTemplates,omitempty"`
	AllowedInApiProducts             *bool                                           `json:"allowedInApiProducts,omitempty"`
	AllowMultiJwtOauth2Subscriptions *bool                                           `json:"allowMultiJwtOauth2Subscriptions,omitempty"`
	Members                          []*APIV4MemberDTO                               `json:"members,omitempty" drift:"empty-is-nil"`
	Failover                         *APIV4FailoverDTO                               `json:"failover,omitempty"`
	PortalNavigation                 []*APIV4NavigationPathDTO                       `json:"portalNavigation,omitempty" drift:"empty-is-nil"`
	ConsoleNotification              *APIV4ConsoleNotificationDTO                    `json:"consoleNotification,omitempty"`
	Pages                            []*APIV4PageDTO                                 `json:"pages" drift:"empty-is-nil"`
	Plans                            []*APIV4PlanDTO                                 `json:"plans" drift:"empty-is-nil"`
}

type APIV4PropertyDTO struct {
	Key         *string `json:"key,omitempty"`
	Value       *string `json:"value,omitempty"`
	Encrypted   *bool   `json:"encrypted,omitempty"`
	Dynamic     *bool   `json:"dynamic,omitempty" drift:"empty-is-nil"`
	Encryptable *bool   `json:"encryptable,omitempty" drift:"empty-is-nil"`
}

type APIV4MetadataEntryDTO struct {
	BaseMetadata `json:",inline"`
	Key          string              `json:"key"`
	Format       base.MetadataFormat `json:"format"`
}

type APIV4ResourceDTO struct {
	Enabled       bool                    `json:"enabled"`
	Name          *string                 `json:"name,omitempty"`
	Type          *string                 `json:"type,omitempty"`
	Configuration *utils.GenericStringMap `json:"configuration,omitempty" drift:"unstructured"`
}

type APIV4ResponseTemplateDTO struct {
	StatusCode              *int              `json:"status,omitempty"`
	Headers                 map[string]string `json:"headers,omitempty"`
	Body                    *string           `json:"body,omitempty"`
	PropagateErrorKeyToLogs *bool             `json:"propagateErrorKeyToLogs,omitempty"`
}

type APIV4DefinitionContextDTO struct {
	Origin   v4.DefinitionContextOrigin `json:"origin,omitempty"`
	SyncFrom v4.DefinitionContextOrigin `json:"syncFrom,omitempty"`
}

type APIV4FailoverDTO struct {
	Enabled           *bool  `json:"enabled,omitempty"`
	MaxRetries        *int   `json:"maxRetries,omitempty"`
	SlowCallDuration  *int64 `json:"slowCallDuration,omitempty"`
	OpenStateDuration *int64 `json:"openStateDuration,omitempty"`
	MaxFailures       *int   `json:"maxFailures,omitempty"`
	PerSubscription   *bool  `json:"perSubscription,omitempty"`
}

type APIV4ConsoleNotificationDTO struct {
	Events []string `json:"events" drift:"empty-is-nil"`
	Groups []string `json:"groups" drift:"ignore-unknown-crd-groups"`
}

type APIV4NavigationPathDTO struct {
	Path        string         `json:"path"`
	DisplayName *string        `json:"displayName,omitempty"`
	Order       *int32         `json:"order,omitempty"`
	Visibility  nav.Visibility `json:"visibility,omitempty" drift:"ignore-unset"`
}

type APIV4MemberDTO struct {
	Source   string `json:"source"`
	SourceID string `json:"sourceId"`
	Role     string `json:"role,omitempty"`
}

type APIV4GenericListenerDTO struct {
	*utils.GenericStringMap `json:",inline" drift:"unstructured"`
}

func (l *APIV4GenericListenerDTO) UnmarshalJSON(data []byte) error {
	if l.GenericStringMap == nil {
		l.GenericStringMap = utils.NewGenericStringMap()
	}
	return l.GenericStringMap.UnmarshalJSON(data)
}

func (l *APIV4GenericListenerDTO) MarshalJSON() ([]byte, error) {
	if l.GenericStringMap == nil {
		return []byte("{}"), nil
	}
	return l.GenericStringMap.MarshalJSON()
}

type APIV4EndpointDTO struct {
	Name           *string                   `json:"name,omitempty"`
	Type           string                    `json:"type,omitempty"`
	Weight         *int32                    `json:"weight,omitempty" drift:"empty-is-nil"`
	Inherit        bool                      `json:"inheritConfiguration"`
	Config         *utils.GenericStringMap   `json:"configuration,omitempty" drift:"unstructured"`
	ConfigOverride *utils.GenericStringMap   `json:"sharedConfigurationOverride,omitempty" drift:"unstructured"`
	Services       *APIV4EndpointServicesDTO `json:"services,omitempty"`
	Secondary      *bool                     `json:"secondary,omitempty"`
	Tenants        []string                  `json:"tenants" drift:"empty-is-nil"`
}

type APIV4LoadBalancerDTO struct {
	Type v4.LoadBalancerType `json:"type"`
}

type APIV4EndpointServicesDTO struct {
	HealthCheck *APIV4ServiceDTO `json:"healthCheck,omitempty"`
}

type APIV4EndpointGroupServicesDTO struct {
	Discovery   *APIV4ServiceDTO `json:"discovery,omitempty"`
	HealthCheck *APIV4ServiceDTO `json:"healthCheck,omitempty"`
}

type APIV4ServiceDTO struct {
	Enabled        bool                    `json:"enabled"`
	Type           *string                 `json:"type,omitempty"`
	OverrideConfig bool                    `json:"overrideConfiguration"`
	Config         *utils.GenericStringMap `json:"configuration,omitempty" drift:"unstructured"`
}

type APIV4HttpClientOptionsDTO struct {
	IdleTimeout                   *uint64              `json:"idleTimeout,omitempty"`
	ConnectTimeout                *uint64              `json:"connectTimeout,omitempty"`
	KeepAlive                     bool                 `json:"keepAlive"`
	KeepAliveTimeout              uint64               `json:"keepAliveTimeout"`
	ReadTimeout                   *uint64              `json:"readTimeout,omitempty"`
	Pipelining                    bool                 `json:"pipelining"`
	MaxConcurrentConnections      *int                 `json:"maxConcurrentConnections,omitempty"`
	UseCompression                bool                 `json:"useCompression"`
	PropagateClientAcceptEncoding bool                 `json:"propagateClientAcceptEncoding"`
	FollowRedirects               bool                 `json:"followRedirects"`
	ClearTextUpgrade              bool                 `json:"clearTextUpgrade"`
	ProtocolVersion               base.ProtocolVersion `json:"version,omitempty"`
	MaxHeaderSize                 *int                 `json:"maxHeaderSize,omitempty"`
	MaxChunkSize                  *int                 `json:"maxChunkSize,omitempty"`
}

type APIV4HttpClientSslOptionsDTO struct {
	TrustAll         bool                  `json:"trustAll"`
	HostnameVerifier bool                  `json:"hostnameVerifier"`
	TrustStore       *APIV4TrustStoreDTO   `json:"trustStore,omitempty"`
	KeyStore         *APIV4KeyStoreDTO     `json:"keyStore,omitempty"`
	Headers          []*APIV4HttpHeaderDTO `json:"headers,omitempty"`
}

type APIV4TrustStoreDTO struct {
	TrustStoreType base.KeyStoreType `json:"type,omitempty"`
	Path           *string           `json:"path,omitempty"`
	Content        *string           `json:"content,omitempty"`
	Password       *string           `json:"password,omitempty"`
}

type APIV4KeyStoreDTO struct {
	KeyStoreType base.KeyStoreType `json:"type,omitempty"`
	Path         *string           `json:"path,omitempty"`
	Content      *string           `json:"content,omitempty"`
	Password     *string           `json:"password,omitempty"`
	KeyPath      *string           `json:"keyPath,omitempty"`
	KeyContent   *string           `json:"keyContent,omitempty"`
	CertPath     *string           `json:"certPath,omitempty"`
	CertContent  *string           `json:"certContent,omitempty"`
}

type APIV4HttpHeaderDTO struct {
	Name  *string `json:"name,omitempty"`
	Value *string `json:"value,omitempty"`
}

type APIV4EndpointGroupDTO struct {
	Name                 string                         `json:"name"`
	Type                 *string                        `json:"type,omitempty"`
	LoadBalancer         *APIV4LoadBalancerDTO          `json:"loadBalancer,omitempty"`
	SharedConfig         *utils.GenericStringMap        `json:"sharedConfiguration,omitempty" drift:"unstructured"`
	Endpoints            []*APIV4EndpointDTO            `json:"endpoints"`
	Services             *APIV4EndpointGroupServicesDTO `json:"services,omitempty"`
	HttpClientOptions    *APIV4HttpClientOptionsDTO     `json:"http,omitempty"`
	HttpClientSslOptions *APIV4HttpClientSslOptionsDTO  `json:"ssl,omitempty"`
	Headers              map[string]string              `json:"headers,omitempty"`
}

type APIV4FlowExecutionDTO struct {
	Mode          v4.FlowMode `json:"mode,omitempty" drift:"ignore-remote:DEFAULT"`
	MatchRequired bool        `json:"matchRequired"`
}

type APIV4FlowDTO struct {
	ID        string                  `json:"id,omitempty" drift:"ignore"`
	Name      *string                 `json:"name,omitempty"`
	Enabled   bool                    `json:"enabled"`
	Selectors []*APIV4FlowSelectorDTO `json:"selectors,omitempty" drift:"empty-is-nil"`
	Request   []*APIV4FlowStepDTO     `json:"request,omitempty" drift:"empty-is-nil"`
	Response  []*APIV4FlowStepDTO     `json:"response,omitempty" drift:"empty-is-nil"`
	Subscribe []*APIV4FlowStepDTO     `json:"subscribe,omitempty" drift:"empty-is-nil"`
	Publish   []*APIV4FlowStepDTO     `json:"publish,omitempty" drift:"empty-is-nil"`
	Connect   []*APIV4FlowStepDTO     `json:"connect,omitempty" drift:"empty-is-nil"`
	Interact  []*APIV4FlowStepDTO     `json:"interact,omitempty" drift:"empty-is-nil"`
	Tags      []string                `json:"tags,omitempty" drift:"empty-is-nil"`
}

type APIV4FlowStepDTO struct {
	Enabled          bool                    `json:"enabled"`
	Policy           *string                 `json:"policy,omitempty"`
	Name             *string                 `json:"name,omitempty"`
	Description      *string                 `json:"description,omitempty"`
	Configuration    *utils.GenericStringMap `json:"configuration,omitempty" drift:"unstructured"`
	Condition        *string                 `json:"condition,omitempty"`
	MessageCondition *string                 `json:"messageCondition,omitempty"`
}

type APIV4FlowSelectorDTO struct {
	*utils.GenericStringMap `json:",inline" drift:"unstructured"`
}

func (s *APIV4FlowSelectorDTO) UnmarshalJSON(data []byte) error {
	if s.GenericStringMap == nil {
		s.GenericStringMap = utils.NewGenericStringMap()
	}
	return s.GenericStringMap.UnmarshalJSON(data)
}

func (s *APIV4FlowSelectorDTO) MarshalJSON() ([]byte, error) {
	if s.GenericStringMap == nil {
		return []byte("{}"), nil
	}
	return s.GenericStringMap.MarshalJSON()
}

type APIV4LoggingPhaseDTO struct {
	Request  bool `json:"request"`
	Response bool `json:"response"`
}

type APIV4LoggingModeDTO struct {
	Entrypoint bool `json:"entrypoint"`
	Endpoint   bool `json:"endpoint"`
}

type APIV4LoggingContentDTO struct {
	Headers         bool `json:"headers"`
	MessageHeaders  bool `json:"messageHeaders"`
	Payload         bool `json:"payload"`
	MessagePayload  bool `json:"messagePayload"`
	MessageMetadata bool `json:"messageMetadata"`
}

type APIV4LoggingDTO struct {
	Condition        *string                 `json:"condition,omitempty"`
	MessageCondition *string                 `json:"messageCondition,omitempty"`
	Content          *APIV4LoggingContentDTO `json:"content,omitempty"`
	Mode             *APIV4LoggingModeDTO    `json:"mode,omitempty"`
	Phase            *APIV4LoggingPhaseDTO   `json:"phase,omitempty"`
}

type APIV4OtelLogsDTO struct {
	Enabled *bool `json:"enabled,omitempty"`
}

type APIV4SamplingDTO struct {
	Type  v4.SamplingType `json:"type"`
	Value string          `json:"value"`
}

type APIV4TracingDTO struct {
	Enabled *bool `json:"enabled,omitempty"`
	Verbose *bool `json:"verbose,omitempty"`
}

type APIV4AnalyticsDTO struct {
	Enabled                bool              `json:"enabled" drift:"empty-is-true"`
	ReporterMetricsEnabled *bool             `json:"reporterMetricsEnabled,omitempty" drift:"empty-is-true"`
	OtelLogs               *APIV4OtelLogsDTO `json:"otelLogs,omitempty"`
	Sampling               *APIV4SamplingDTO `json:"sampling,omitempty"`
	Logging                *APIV4LoggingDTO  `json:"logging,omitempty"`
	Tracing                *APIV4TracingDTO  `json:"tracing,omitempty"`
}

type APIV4ApiServicesDTO struct {
	DynamicProperty *APIV4ServiceDTO `json:"dynamicProperty,omitempty"`
}

type APIV4PlanDTO struct {
	ID                    string                `json:"id,omitempty" drift:"ignore"`
	HRID                  string                `json:"hrid,omitempty"`
	CrossID               string                `json:"crossId,omitempty" drift:"ignore"`
	Tags                  []string              `json:"tags" drift:"empty-is-nil"`
	Status                base.PlanStatus       `json:"status,omitempty"`
	Characteristics       []string              `json:"characteristics" drift:"empty-is-nil"`
	Validation            base.PlanValidation   `json:"validation,omitempty"`
	CommentRequired       *bool                 `json:"comment_required,omitempty"`
	Order                 *int                  `json:"order,omitempty"`
	Type                  base.PlanType         `json:"type,omitempty"`
	Name                  string                `json:"name"`
	Description           *string               `json:"description,omitempty"`
	DefinitionVersion     v4.DefinitionVersion  `json:"definitionVersion,omitempty" drift:"ignore"`
	Security              *APIV4PlanSecurityDTO `json:"security,omitempty"`
	Mode                  v4.PlanMode           `json:"mode,omitempty"`
	SelectionRule         *string               `json:"selectionRule,omitempty"`
	Flows                 []*APIV4FlowDTO       `json:"flows"`
	ExcludedGroups        []string              `json:"excludedGroups"`
	GeneralConditionsHRID *string               `json:"generalConditionsHrid,omitempty"`
	BootstrapPort         *int                  `json:"bootstrapPort,omitempty"`
	BrokerRangeStart      *int                  `json:"brokerRangeStart,omitempty"`
	BrokerRangeEnd        *int                  `json:"brokerRangeEnd,omitempty"`
}

type APIV4PlanSecurityDTO struct {
	Type   string                  `json:"type" drift:"case-insensitive"`
	Config *utils.GenericStringMap `json:"configuration,omitempty" drift:"unstructured"`
}

type APIV4PageSourceDTO struct {
	Type          string                  `json:"type"`
	Configuration *utils.GenericStringMap `json:"configuration" drift:"unstructured"`
}

type APIV4PageDTO struct {
	HRID          string              `json:"hrid,omitempty"`
	CrossID       string              `json:"crossId,omitempty" drift:"ignore"`
	Name          string              `json:"name,omitempty"`
	Type          string              `json:"type"`
	Content       *string             `json:"content,omitempty"`
	Order         *uint64             `json:"order,omitempty"`
	Published     bool                `json:"published"`
	Visibility    string              `json:"visibility,omitempty"`
	HomePage      bool                `json:"homepage"`
	ParentHRID    *string             `json:"parentHrid,omitempty"`
	API           *string             `json:"api,omitempty" drift:"ignore"`
	Source        *APIV4PageSourceDTO `json:"source,omitempty"`
	Configuration map[string]string   `json:"configuration,omitempty"`
}
