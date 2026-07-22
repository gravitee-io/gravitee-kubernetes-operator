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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
)

type APIV4DTO struct {
	ID                               string                                       `json:"id,omitempty" drift:"ignore"`
	HRID                             string                                       `json:"hrid,omitempty"`
	CrossID                          string                                       `json:"crossId,omitempty" drift:"ignore"`
	Name                             string                                       `json:"name"`
	Version                          string                                       `json:"version"`
	State                            base.ApiState                                `json:"state,omitempty"`
	Tags                             []string                                     `json:"tags" drift:"empty-is-nil"`
	Labels                           []string                                     `json:"labels" drift:"empty-is-nil"`
	Visibility                       base.ApiVisibility                           `json:"visibility,omitempty"`
	Properties                       []*APIV4Property                             `json:"properties" drift:"empty-is-nil"`
	Metadata                         []*APIV4MetadataEntry                        `json:"metadata" drift:"ignore-remote-only-metadata"`
	Resources                        []*APIV4Resource                             `json:"resources" drift:"empty-is-nil"`
	Groups                           []string                                     `json:"groups" drift:"ignore-crd-only-and-namespace-prefix"`
	Categories                       []string                                     `json:"categories" drift:"empty-is-nil"`
	NotifyMembers                    bool                                         `json:"notifyMembers" drift:"empty-is-true"`
	Description                      *string                                      `json:"description,omitempty"`
	DefinitionVersion                base.DefinitionVersion                       `json:"definitionVersion,omitempty" drift:"ignore"`
	DefinitionContext                *APIV4DefinitionContext                      `json:"definitionContext,omitempty" drift:"ignore"`
	LifecycleState                   v4.ApiV4LifecycleState                       `json:"lifecycleState,omitempty"`
	Type                             v4.ApiType                                   `json:"type"`
	Listeners                        []*APIV4GenericListener                      `json:"listeners"`
	EndpointGroups                   []*APIV4EndpointGroup                        `json:"endpointGroups"`
	FlowExecution                    *APIV4FlowExecution                          `json:"flowExecution,omitempty"`
	Flows                            []*APIV4Flow                                 `json:"flows" drift:"empty-is-nil"`
	Analytics                        *APIV4Analytics                              `json:"analytics,omitempty"`
	Services                         *APIV4ApiServices                            `json:"services,omitempty"`
	ResponseTemplates                map[string]map[string]*APIV4ResponseTemplate `json:"responseTemplates,omitempty"`
	AllowedInApiProducts             *bool                                        `json:"allowedInApiProducts,omitempty"`
	AllowMultiJwtOauth2Subscriptions *bool                                        `json:"allowMultiJwtOauth2Subscriptions,omitempty"`
	Members                          []*APIV4Member                               `json:"members,omitempty" drift:"empty-is-nil"`
	Failover                         *APIV4Failover                               `json:"failover,omitempty"`
	PortalNavigation                 []*APIV4NavigationPath                       `json:"portalNavigation,omitempty" drift:"empty-is-nil"`
	ConsoleNotification              *APIV4ConsoleNotification                    `json:"consoleNotification,omitempty"`
	Pages                            []*APIV4Page                                 `json:"pages" drift:"empty-is-nil"`
	Plans                            []*APIV4Plan                                 `json:"plans" drift:"empty-is-nil"`
}

type APIV4Property struct {
	Key         *string `json:"key,omitempty"`
	Value       *string `json:"value,omitempty"`
	Encrypted   *bool   `json:"encrypted,omitempty"`
	Dynamic     *bool   `json:"dynamic,omitempty" drift:"empty-is-nil"`
	Encryptable *bool   `json:"encryptable,omitempty" drift:"empty-is-nil"`
}

type APIV4MetadataEntry struct {
	BaseMetadata `json:",inline"`
	Key          string              `json:"key"`
	Format       base.MetadataFormat `json:"format"`
}

type APIV4Resource struct {
	Enabled       bool                    `json:"enabled"`
	Name          *string                 `json:"name,omitempty"`
	Type          *string                 `json:"type,omitempty"`
	Configuration *utils.GenericStringMap `json:"configuration,omitempty" drift:"unstructured"`
}

type APIV4ResponseTemplate struct {
	StatusCode              *int              `json:"status,omitempty"`
	Headers                 map[string]string `json:"headers,omitempty"`
	Body                    *string           `json:"body,omitempty"`
	PropagateErrorKeyToLogs *bool             `json:"propagateErrorKeyToLogs,omitempty"`
}

type APIV4DefinitionContext struct {
	Origin   v4.DefinitionContextOrigin `json:"origin,omitempty"`
	SyncFrom v4.DefinitionContextOrigin `json:"syncFrom,omitempty"`
}

type APIV4Failover struct {
	Enabled           *bool  `json:"enabled,omitempty"`
	MaxRetries        *int   `json:"maxRetries,omitempty"`
	SlowCallDuration  *int64 `json:"slowCallDuration,omitempty"`
	OpenStateDuration *int64 `json:"openStateDuration,omitempty"`
	MaxFailures       *int   `json:"maxFailures,omitempty"`
	PerSubscription   *bool  `json:"perSubscription,omitempty"`
}

type APIV4ConsoleNotification struct {
	Events []string `json:"events" drift:"empty-is-nil"`
	Groups []string `json:"groups" drift:"ignore-crd-only-and-namespace-prefix"`
}

type APIV4NavigationPath struct {
	Path        string  `json:"path"`
	DisplayName *string `json:"displayName,omitempty"`
	Order       *int32  `json:"order,omitempty"`
}

type APIV4Member struct {
	Source   string `json:"source"`
	SourceID string `json:"sourceId"`
	Role     string `json:"role,omitempty"`
}

type APIV4GenericListener struct {
	*utils.GenericStringMap `json:",inline" drift:"unstructured"`
}

func (l *APIV4GenericListener) UnmarshalJSON(data []byte) error {
	if l.GenericStringMap == nil {
		l.GenericStringMap = utils.NewGenericStringMap()
	}
	return l.GenericStringMap.UnmarshalJSON(data)
}

func (l *APIV4GenericListener) MarshalJSON() ([]byte, error) {
	if l.GenericStringMap == nil {
		return []byte("{}"), nil
	}
	return l.GenericStringMap.MarshalJSON()
}

type APIV4Endpoint struct {
	Name           *string                 `json:"name,omitempty"`
	Type           string                  `json:"type,omitempty"`
	Weight         *int32                  `json:"weight,omitempty" drift:"empty-is-nil"`
	Inherit        bool                    `json:"inheritConfiguration"`
	Config         *utils.GenericStringMap `json:"configuration,omitempty" drift:"unstructured"`
	ConfigOverride *utils.GenericStringMap `json:"sharedConfigurationOverride,omitempty" drift:"unstructured"`
	Services       *APIV4EndpointServices  `json:"services,omitempty"`
	Secondary      *bool                   `json:"secondary,omitempty"`
	Tenants        []string                `json:"tenants" drift:"empty-is-nil"`
}

type APIV4LoadBalancer struct {
	Type v4.LoadBalancerType `json:"type"`
}

type APIV4EndpointServices struct {
	HealthCheck *APIV4Service `json:"healthCheck,omitempty"`
}

type APIV4EndpointGroupServices struct {
	Discovery   *APIV4Service `json:"discovery,omitempty"`
	HealthCheck *APIV4Service `json:"healthCheck,omitempty"`
}

type APIV4Service struct {
	Enabled        bool                    `json:"enabled"`
	Type           *string                 `json:"type,omitempty"`
	OverrideConfig bool                    `json:"overrideConfiguration"`
	Config         *utils.GenericStringMap `json:"configuration,omitempty" drift:"unstructured"`
}

type APIV4HttpClientOptions struct {
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

type APIV4HttpClientSslOptions struct {
	TrustAll         bool               `json:"trustAll"`
	HostnameVerifier bool               `json:"hostnameVerifier"`
	TrustStore       *APIV4TrustStore   `json:"trustStore,omitempty"`
	KeyStore         *APIV4KeyStore     `json:"keyStore,omitempty"`
	Headers          []*APIV4HttpHeader `json:"headers,omitempty"`
}

type APIV4TrustStore struct {
	TrustStoreType base.KeyStoreType `json:"type,omitempty"`
	Path           *string           `json:"path,omitempty"`
	Content        *string           `json:"content,omitempty"`
	Password       *string           `json:"password,omitempty"`
}

type APIV4KeyStore struct {
	KeyStoreType base.KeyStoreType `json:"type,omitempty"`
	Path         *string           `json:"path,omitempty"`
	Content      *string           `json:"content,omitempty"`
	Password     *string           `json:"password,omitempty"`
	KeyPath      *string           `json:"keyPath,omitempty"`
	KeyContent   *string           `json:"keyContent,omitempty"`
	CertPath     *string           `json:"certPath,omitempty"`
	CertContent  *string           `json:"certContent,omitempty"`
}

type APIV4HttpHeader struct {
	Name  *string `json:"name,omitempty"`
	Value *string `json:"value,omitempty"`
}

type APIV4EndpointGroup struct {
	Name                 string                      `json:"name"`
	Type                 *string                     `json:"type,omitempty"`
	LoadBalancer         *APIV4LoadBalancer          `json:"loadBalancer,omitempty"`
	SharedConfig         *utils.GenericStringMap     `json:"sharedConfiguration,omitempty" drift:"unstructured"`
	Endpoints            []*APIV4Endpoint            `json:"endpoints"`
	Services             *APIV4EndpointGroupServices `json:"services,omitempty"`
	HttpClientOptions    *APIV4HttpClientOptions     `json:"http,omitempty"`
	HttpClientSslOptions *APIV4HttpClientSslOptions  `json:"ssl,omitempty"`
	Headers              map[string]string           `json:"headers,omitempty"`
}

type APIV4FlowExecution struct {
	Mode          v4.FlowMode `json:"mode,omitempty"`
	MatchRequired bool        `json:"matchRequired"`
}

type APIV4Flow struct {
	ID        string               `json:"id,omitempty"`
	Name      *string              `json:"name,omitempty"`
	Enabled   bool                 `json:"enabled"`
	Selectors []*APIV4FlowSelector `json:"selectors,omitempty"`
	Request   []*APIV4FlowStep     `json:"request,omitempty"`
	Response  []*APIV4FlowStep     `json:"response,omitempty"`
	Subscribe []*APIV4FlowStep     `json:"subscribe,omitempty"`
	Publish   []*APIV4FlowStep     `json:"publish,omitempty"`
	Connect   []*APIV4FlowStep     `json:"connect,omitempty"`
	Interact  []*APIV4FlowStep     `json:"interact,omitempty"`
	Tags      []string             `json:"tags,omitempty" drift:"empty-is-nil"`
}

type APIV4FlowStep struct {
	Enabled          bool                    `json:"enabled"`
	Policy           *string                 `json:"policy,omitempty"`
	Name             *string                 `json:"name,omitempty"`
	Description      *string                 `json:"description,omitempty"`
	Configuration    *utils.GenericStringMap `json:"configuration,omitempty" drift:"unstructured"`
	Condition        *string                 `json:"condition,omitempty"`
	MessageCondition *string                 `json:"messageCondition,omitempty"`
}

type APIV4FlowSelector struct {
	*utils.GenericStringMap `json:",inline" drift:"unstructured"`
}

func (s *APIV4FlowSelector) UnmarshalJSON(data []byte) error {
	if s.GenericStringMap == nil {
		s.GenericStringMap = utils.NewGenericStringMap()
	}
	return s.GenericStringMap.UnmarshalJSON(data)
}

func (s *APIV4FlowSelector) MarshalJSON() ([]byte, error) {
	if s.GenericStringMap == nil {
		return []byte("{}"), nil
	}
	return s.GenericStringMap.MarshalJSON()
}

type APIV4LoggingPhase struct {
	Request  bool `json:"request"`
	Response bool `json:"response"`
}

type APIV4LoggingMode struct {
	Entrypoint bool `json:"entrypoint"`
	Endpoint   bool `json:"endpoint"`
}

type APIV4LoggingContent struct {
	Headers         bool `json:"headers"`
	MessageHeaders  bool `json:"messageHeaders"`
	Payload         bool `json:"payload"`
	MessagePayload  bool `json:"messagePayload"`
	MessageMetadata bool `json:"messageMetadata"`
}

type APIV4Logging struct {
	Condition        *string              `json:"condition,omitempty"`
	MessageCondition *string              `json:"messageCondition,omitempty"`
	Content          *APIV4LoggingContent `json:"content,omitempty"`
	Mode             *APIV4LoggingMode    `json:"mode,omitempty"`
	Phase            *APIV4LoggingPhase   `json:"phase,omitempty"`
}

type APIV4OtelLogs struct {
	Enabled *bool `json:"enabled,omitempty"`
}

type APIV4Sampling struct {
	Type  v4.SamplingType `json:"type"`
	Value string          `json:"value"`
}

type APIV4Tracing struct {
	Enabled *bool `json:"enabled,omitempty"`
	Verbose *bool `json:"verbose,omitempty"`
}

type APIV4Analytics struct {
	Enabled                bool           `json:"enabled" drift:"empty-is-true"`
	ReporterMetricsEnabled *bool          `json:"reporterMetricsEnabled,omitempty" drift:"empty-is-true"`
	OtelLogs               *APIV4OtelLogs `json:"otelLogs,omitempty"`
	Sampling               *APIV4Sampling `json:"sampling,omitempty"`
	Logging                *APIV4Logging  `json:"logging,omitempty"`
	Tracing                *APIV4Tracing  `json:"tracing,omitempty"`
}

type APIV4ApiServices struct {
	DynamicProperty *APIV4Service `json:"dynamicProperty,omitempty"`
}

type APIV4Plan struct {
	ID                    string               `json:"id,omitempty" drift:"ignore"`
	HRID                  string               `json:"hrid,omitempty"`
	CrossID               string               `json:"crossId,omitempty" drift:"ignore"`
	Tags                  []string             `json:"tags" drift:"empty-is-nil"`
	Status                base.PlanStatus      `json:"status,omitempty"`
	Characteristics       []string             `json:"characteristics"`
	Validation            base.PlanValidation  `json:"validation,omitempty"`
	CommentRequired       *bool                `json:"comment_required,omitempty"`
	Order                 *int                 `json:"order,omitempty"`
	Type                  base.PlanType        `json:"type,omitempty"`
	Name                  string               `json:"name"`
	Description           *string              `json:"description,omitempty"`
	DefinitionVersion     v4.DefinitionVersion `json:"definitionVersion,omitempty" drift:"ignore"`
	Security              *APIV4PlanSecurity   `json:"security,omitempty"`
	Mode                  v4.PlanMode          `json:"mode,omitempty"`
	SelectionRule         *string              `json:"selectionRule,omitempty"`
	Flows                 []*APIV4Flow         `json:"flows"`
	ExcludedGroups        []string             `json:"excludedGroups"`
	GeneralConditionsHRID *string              `json:"generalConditionsHrid,omitempty"`
	BootstrapPort         *int                 `json:"bootstrapPort,omitempty"`
	BrokerRangeStart      *int                 `json:"brokerRangeStart,omitempty"`
	BrokerRangeEnd        *int                 `json:"brokerRangeEnd,omitempty"`
}

type APIV4PlanSecurity struct {
	Type   string                  `json:"type"`
	Config *utils.GenericStringMap `json:"configuration,omitempty" drift:"unstructured"`
}

type APIV4PageSource struct {
	Type          string                  `json:"type"`
	Configuration *utils.GenericStringMap `json:"configuration" drift:"unstructured"`
}

type APIV4Page struct {
	HRID          string            `json:"hrid,omitempty"`
	CrossID       string            `json:"crossId,omitempty"`
	Name          string            `json:"name,omitempty"`
	Type          string            `json:"type"`
	Content       *string           `json:"content,omitempty"`
	Order         *uint64           `json:"order,omitempty"`
	Published     bool              `json:"published"`
	Visibility    string            `json:"visibility,omitempty"`
	HomePage      bool              `json:"homepage"`
	ParentHRID    *string           `json:"parentHrid,omitempty"`
	API           *string           `json:"api,omitempty"`
	Source        *APIV4PageSource  `json:"source,omitempty"`
	Configuration map[string]string `json:"configuration,omitempty"`
}
