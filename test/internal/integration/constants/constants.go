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

package constants

import (
	"time"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

const (
	Namespace = "default"

	ConsistentTimeout = time.Second * 2
	EventualTimeout   = time.Second * 30
	Interval          = time.Millisecond * 250

	GatewayHost   = "localhost"
	GatewayPort   = "30082"
	GatewayUrl    = "http://" + GatewayHost + ":" + GatewayPort
	ManagementUrl = "http://localhost:30083/management"

	ContextSecretFile             = "management_context/dev/management-context-secret.yml"
	ContextWithSecretFile         = "management_context/dev/management-context-with-secret-ref.yml"
	ContextWithCredentialsFile    = "management_context/dev/management-context-with-credentials.yml"
	ContextWithBadCredentialsFile = "management_context/dev/management-context-with-bearer-token.yml"
	ContextWithBadURLFile         = "management_context/debug/management-context-with-credentials.yml"

	// V2 APIs.
	Api                                 = "apim/api_definition/v2/api.yml"
	ApiWithRateLimit                    = "apim/api_definition/v2/api-with-rate-limit.yml"
	ApiWithStateStopped                 = "apim/api_definition/v2/api-with-state-stopped.yml"
	ApiWithSyncFromAPIM                 = "apim/api_definition/v2/api-with-sync-from-apim.yml"
	ApiWithIDs                          = "apim/api_definition/v2/api-with-ids.yml"
	ApiWithDisabledPolicy               = "apim/api_definition/v2/api-with-disabled-policy.yml"
	ApiWithTemplatingFile               = "apim/api_definition/v2/api-with-templating.yml"
	ApiWithTemplatingSecretFile         = "apim/api_definition/v2/api-with-templating-secret.yml"
	ApiWithTemplatingConfigMapFile      = "apim/api_definition/v2/api-with-templating-config-map.yml"
	ApiWithContextFile                  = "apim/api_definition/v2/api-with-context.yml"
	ApiWithContextNoPlanFile            = "apim/api_definition/v2/api-with-no-plan.yml"
	ApiWithDisabledHCFile               = "apim/api_definition/v2/api-with-health-check-disabled.yml"
	ApiWithHCFile                       = "apim/api_definition/v2/api-with-health-check.yml"
	ApiWithServiceDiscoveryFile         = "apim/api_definition/v2/api-with-service-discovery.yml"
	ApiWithMetadataFile                 = "apim/api_definition/v2/api-with-metadata.yml"
	ApiWithEndpointGroupsFile           = "apim/api_definition/v2/api-with-endpoint-groups.yml"
	ApiWithLoggingFile                  = "apim/api_definition/v2/api-with-logging.yml"
	ApiWithApiKeyPlanFile               = "apim/api_definition/v2/api-with-api-key-plan.yml"
	ApiWithCacheResourceFile            = "apim/api_definition/v2/api-with-cache-resource.yml"
	ApiWithCacheResourceRefFile         = "apim/api_definition/v2/api-with-cache-resource-ref.yml"
	ApiWithCacheRedisResourceFile       = "apim/api_definition/v2/api-with-cache-redis-resource.yml"
	ApiWithCacheRedisResourceRefFile    = "apim/api_definition/v2/api-with-cache-redis-resource-ref.yml"
	ApiWithOAuth2GenericResourceFile    = "apim/api_definition/v2/api-with-oauth2-generic-resource.yml"
	ApiWithOAuth2GenericResourceRefFile = "apim/api_definition/v2/api-with-oauth2-generic-resource-ref.yml"
	ApiWithOauth2AmResourceFile         = "apim/api_definition/v2/api-with-oauth2-am-resource.yml"
	ApiWithOauth2AmResourceRefFile      = "apim/api_definition/v2/api-with-oauth2-am-resource-ref.yml"
	ApiWithKeycloakAdapterFile          = "apim/api_definition/v2/api-with-keycloak-adapter.yml"
	ApiWithKeycloakAdapterRefFile       = "apim/api_definition/v2/api-with-keycloak-adapter-ref.yml"
	ApiWithLDAPAuthProviderFile         = "apim/api_definition/v2/api-with-ldap-auth-provider.yml"
	ApiWithLDAPAuthProviderRefFile      = "apim/api_definition/v2/api-with-ldap-auth-provider-ref.yml"
	ApiWithInlineAuthProviderFile       = "apim/api_definition/v2/api-with-inline-auth-provider.yml"
	ApiWithInlineAuthProviderRefFile    = "apim/api_definition/v2/api-with-inline-auth-provider-ref.yml"
	ApiWithHTTPAuthProviderFile         = "apim/api_definition/v2/api-with-http-auth-provider.yml"
	ApiWithHTTPAuthProviderRefFile      = "apim/api_definition/v2/api-with-http-auth-provider-ref.yml"
	ApiWithMarkdownPage                 = "apim/api_definition/v2/api-with-page-markdown.yml"
	ApiWithSwaggerHTTPFetcher           = "apim/api_definition/v2/api-with-page-swagger-http-fetcher.yml"
	ApiWithMembersAndGroups             = "apim/api_definition/v2/api-with-groups-members.yml"

	ApiResourceCacheFile              = "apim/api_resource/api-resource-cache.yml"
	ApiResourceCacheRedisFile         = "apim/api_resource/api-resource-cache-redis.yml"
	ApiResourceHTTPAuthProviderFile   = "apim/api_resource/api-resource-http-auth-provider.yml"
	ApiResourceInlineAuthProviderFile = "apim/api_resource/api-resource-inline-auth-provider.yml"
	ApiResourceLDAPAuthProviderFile   = "apim/api_resource/api-resource-ldap-auth-provider.yml"
	ApiResourceKeycloakAdapterFile    = "apim/api_resource/api-resource-keycloak-adapter.yml"
	ApiResourceOauth2AMFile           = "apim/api_resource/api-resource-oauth2-am.yml"
	ApiResourceOauth2GenericFile      = "apim/api_resource/api-resource-oauth2-generic.yml"

	ApiWithTemplateAnnotation = "apim/api_definition/v2/api-with-template-annotation.yml"

	// V4 APIS.
	ApiV4                              = "apim/api_definition/v4/api-v4.yml"
	ApiV4WithSyncFromAPIM              = "apim/api_definition/v4/api-v4-with-sync-from-apim.yml"
	ApiV4WithTemplatingFile            = "apim/api_definition/v4/api-v4-with-templating.yml"
	ApiV4WithRateLimit                 = "apim/api_definition/v4/api-v4-with-rate-limit.yml"
	ApiV4WithDisabledPolicy            = "apim/api_definition/v4/api-v4-with-disabled-policy.yml"
	ApiV4WithContextFile               = "apim/api_definition/v4/api-v4-with-context.yml"
	ApiV4WithHCFile                    = "apim/api_definition/v4/api-v4-with-health-check.yml"
	ApiV4WithDisabledHCFile            = "apim/api_definition/v4/api-v4-with-health-check-disabled.yml"
	ApiV4WithLoggingFile               = "apim/api_definition/v4/api-v4-with-logging.yml"
	ApiV4WithMetadataFile              = "apim/api_definition/v4/api-v4-with-metadata.yml"
	ApiV4WithCacheRedisResourceFile    = "apim/api_definition/v4/api-v4-with-cache-redis-resource.yml"
	ApiV4WithOAuth2GenericResourceFile = "apim/api_definition/v4/api-v4-with-oauth2-generic-resource.yml"
	ApiV4WithOauth2AmResourceFile      = "apim/api_definition/v4/api-v4-with-oauth2-am-resource.yml"
	ApiV4WithKeycloakAdapterFile       = "apim/api_definition/v4/api-v4-with-keycloak-adapter.yml"
	ApiV4WithLDAPAuthProviderFile      = "apim/api_definition/v4/api-v4-with-ldap-auth-provider.yml"
	ApiV4WithInlineAuthProviderFile    = "apim/api_definition/v4/api-v4-with-inline-auth-provider.yml"
	ApiV4WithHTTPAuthProviderFile      = "apim/api_definition/v4/api-v4-with-http-auth-provider.yml"
	ApiV4WithStateStopped              = "apim/api_definition/v4/api-v4-with-state-stopped.yml"
	ApiV4WithApiKeyPlanFile            = "apim/api_definition/v4/api-v4-with-api-key-plan.yml"
	ApiV4WithCacheRedisResourceRef     = "apim/api_definition/v4/api-v4-with-cache-redis-resource-ref.yml"
	ApiV4WithOAuth2GenericResRef       = "apim/api_definition/v4/api-v4-with-oauth2-generic-res-ref.yml"
	ApiV4WithOauth2AmResourceRefFile   = "apim/api_definition/v4/api-v4-with-oauth2-am-resource-ref.yml"
	ApiV4WithKeycloakAdapterRefFile    = "apim/api_definition/v4/api-v4-with-keycloak-adapter-ref.yml"
	ApiV4WithLDAPAuthProviderRefFile   = "apim/api_definition/v4/api-v4-with-ldap-auth-provider-ref.yml"
	ApiV4WithInlineAuthProviderRef     = "apim/api_definition/v4/api-v4-with-inline-auth-provider-ref.yml"
	ApiV4WithHTTPAuthProviderRefFile   = "apim/api_definition/v4/api-v4-with-http-auth-provider-ref.yml"
	ApiV4WithMarkdownPage              = "apim/api_definition/v4/api-v4-with-page-markdown.yml"
	ApiV4WithSwaggerHTTPFetcher        = "apim/api_definition/v4/api-v4-with-page-swagger-http-fetcher.yml"

	IngressPEMRegistry         = "ingress/ingress-pem-registry.yml"
	Ingress404ResponseTemplate = "ingress/ingress-response-404-config-map.yml"
	IngressWithoutTemplateFile = "ingress/ingress-without-api-template.yml"
	IngressWithTemplateFile    = "ingress/ingress-with-api-template.yml"
	IngressWithMultipleHosts   = "ingress/ingress-with-multiple-hosts.yml"
	IngressWithTLS             = "ingress/ingress-with-tls.yml"
	IngressWithTLSSecretFile   = "ingress/ingress-with-tls-secret.yml"
	IngressResponse404CMFile   = "ingress/ingress-response-404-config-map.yaml"

	Application = "apim/application/application.yml"
)

func BuildAPIEndpoint(api *v1alpha1.ApiDefinition) string {
	return GatewayUrl + api.Spec.Proxy.VirtualHosts[0].Path
}

func BuildAPIV4Endpoint(l v4.Listener) string {
	switch t := l.(type) {
	case *v4.GenericListener:
		return BuildAPIV4Endpoint(t.ToListener())
	case *v4.HttpListener:
		return GatewayUrl + t.Paths[0].Path
	case *v4.TCPListener:
		return t.Hosts[0]
	}

	return ""
}
