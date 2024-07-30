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
	SamplesPath   = "../../../examples"

	ContextSecretFile             = SamplesPath + "/management_context/dev/management-context-secret.yml"
	ContextWithSecretFile         = SamplesPath + "/management_context/dev/management-context-with-secret-ref.yml"
	ContextWithCredentialsFile    = SamplesPath + "/management_context/dev/management-context-with-credentials.yml"
	ContextWithBadCredentialsFile = SamplesPath + "/management_context/dev/management-context-with-bearer-token.yml"
	ContextWithBadURLFile         = SamplesPath + "/management_context/debug/management-context-with-credentials.yml"

	// V2 APIs.
	Api                                 = SamplesPath + "/apim/api_definition/v2/api.yml"
	ApiWithRateLimit                    = SamplesPath + "/apim/api_definition/v2/api-with-rate-limit.yml"
	ApiWithStateStopped                 = SamplesPath + "/apim/api_definition/v2/api-with-state-stopped.yml"
	ApiWithSyncFromAPIM                 = SamplesPath + "/apim/api_definition/v2/api-with-sync-from-apim.yml"
	ApiWithIDs                          = SamplesPath + "/apim/api_definition/v2/api-with-ids.yml"
	ApiWithDisabledPolicy               = SamplesPath + "/apim/api_definition/v2/api-with-disabled-policy.yml"
	ApiWithTemplatingFile               = SamplesPath + "/apim/api_definition/v2/api-with-templating.yml"
	ApiWithTemplatingSecretFile         = SamplesPath + "/apim/api_definition/v2/api-with-templating-secret.yml"
	ApiWithTemplatingConfigMapFile      = SamplesPath + "/apim/api_definition/v2/api-with-templating-config-map.yml"
	ApiWithContextFile                  = SamplesPath + "/apim/api_definition/v2/api-with-context.yml"
	ApiWithContextNoPlanFile            = SamplesPath + "/apim/api_definition/v2/api-with-no-plan.yml"
	ApiWithDisabledHCFile               = SamplesPath + "/apim/api_definition/v2/api-with-health-check-disabled.yml"
	ApiWithHCFile                       = SamplesPath + "/apim/api_definition/v2/api-with-health-check.yml"
	ApiWithServiceDiscoveryFile         = SamplesPath + "/apim/api_definition/v2/api-with-service-discovery.yml"
	ApiWithMetadataFile                 = SamplesPath + "/apim/api_definition/v2/api-with-metadata.yml"
	ApiWithEndpointGroupsFile           = SamplesPath + "/apim/api_definition/v2/api-with-endpoint-groups.yml"
	ApiWithLoggingFile                  = SamplesPath + "/apim/api_definition/v2/api-with-logging.yml"
	ApiWithApiKeyPlanFile               = SamplesPath + "/apim/api_definition/v2/api-with-api-key-plan.yml"
	ApiWithCacheResourceFile            = SamplesPath + "/apim/api_definition/v2/api-with-cache-resource.yml"
	ApiWithCacheResourceRefFile         = SamplesPath + "/apim/api_definition/v2/api-with-cache-resource-ref.yml"
	ApiWithCacheRedisResourceFile       = SamplesPath + "/apim/api_definition/v2/api-with-cache-redis-resource.yml"
	ApiWithCacheRedisResourceRefFile    = SamplesPath + "/apim/api_definition/v2/api-with-cache-redis-resource-ref.yml"
	ApiWithOAuth2GenericResourceFile    = SamplesPath + "/apim/api_definition/v2/api-with-oauth2-generic-resource.yml"
	ApiWithOAuth2GenericResourceRefFile = SamplesPath + "/apim/api_definition/v2/api-with-oauth2-generic-resource-ref.yml"
	ApiWithOauth2AmResourceFile         = SamplesPath + "/apim/api_definition/v2/api-with-oauth2-am-resource.yml"
	ApiWithOauth2AmResourceRefFile      = SamplesPath + "/apim/api_definition/v2/api-with-oauth2-am-resource-ref.yml"
	ApiWithKeycloakAdapterFile          = SamplesPath + "/apim/api_definition/v2/api-with-keycloak-adapter.yml"
	ApiWithKeycloakAdapterRefFile       = SamplesPath + "/apim/api_definition/v2/api-with-keycloak-adapter-ref.yml"
	ApiWithLDAPAuthProviderFile         = SamplesPath + "/apim/api_definition/v2/api-with-ldap-auth-provider.yml"
	ApiWithLDAPAuthProviderRefFile      = SamplesPath + "/apim/api_definition/v2/api-with-ldap-auth-provider-ref.yml"
	ApiWithInlineAuthProviderFile       = SamplesPath + "/apim/api_definition/v2/api-with-inline-auth-provider.yml"
	ApiWithInlineAuthProviderRefFile    = SamplesPath + "/apim/api_definition/v2/api-with-inline-auth-provider-ref.yml"
	ApiWithHTTPAuthProviderFile         = SamplesPath + "/apim/api_definition/v2/api-with-http-auth-provider.yml"
	ApiWithHTTPAuthProviderRefFile      = SamplesPath + "/apim/api_definition/v2/api-with-http-auth-provider-ref.yml"
	ApiWithMarkdownPage                 = SamplesPath + "/apim/api_definition/v2/api-with-page-markdown.yml"
	ApiWithSwaggerHTTPFetcher           = SamplesPath + "/apim/api_definition/v2/api-with-page-swagger-http-fetcher.yml"
	ApiWithMembersAndGroups             = SamplesPath + "/apim/api_definition/v2/api-with-groups-members.yml"

	ApiResourceCacheFile              = SamplesPath + "/apim/api_resource/api-resource-cache.yml"
	ApiResourceCacheRedisFile         = SamplesPath + "/apim/api_resource/api-resource-cache-redis.yml"
	ApiResourceHTTPAuthProviderFile   = SamplesPath + "/apim/api_resource/api-resource-http-auth-provider.yml"
	ApiResourceInlineAuthProviderFile = SamplesPath + "/apim/api_resource/api-resource-inline-auth-provider.yml"
	ApiResourceLDAPAuthProviderFile   = SamplesPath + "/apim/api_resource/api-resource-ldap-auth-provider.yml"
	ApiResourceKeycloakAdapterFile    = SamplesPath + "/apim/api_resource/api-resource-keycloak-adapter.yml"
	ApiResourceOauth2AMFile           = SamplesPath + "/apim/api_resource/api-resource-oauth2-am.yml"
	ApiResourceOauth2GenericFile      = SamplesPath + "/apim/api_resource/api-resource-oauth2-generic.yml"

	ApiWithTemplateAnnotation = SamplesPath + "/apim/api_definition/v2/api-with-template-annotation.yml"

	// V4 APIS.
	ApiV4                              = SamplesPath + "/apim/api_definition/v4/api-v4.yml"
	ApiV4WithSyncFromAPIM              = SamplesPath + "/apim/api_definition/v4/api-v4-with-sync-from-apim.yml"
	ApiV4WithTemplatingFile            = SamplesPath + "/apim/api_definition/v4/api-v4-with-templating.yml"
	ApiV4WithRateLimit                 = SamplesPath + "/apim/api_definition/v4/api-v4-with-rate-limit.yml"
	ApiV4WithDisabledPolicy            = SamplesPath + "/apim/api_definition/v4/api-v4-with-disabled-policy.yml"
	ApiV4WithContextFile               = SamplesPath + "/apim/api_definition/v4/api-v4-with-context.yml"
	ApiV4WithHCFile                    = SamplesPath + "/apim/api_definition/v4/api-v4-with-health-check.yml"
	ApiV4WithDisabledHCFile            = SamplesPath + "/apim/api_definition/v4/api-v4-with-health-check-disabled.yml"
	ApiV4WithLoggingFile               = SamplesPath + "/apim/api_definition/v4/api-v4-with-logging.yml"
	ApiV4WithMetadataFile              = SamplesPath + "/apim/api_definition/v4/api-v4-with-metadata.yml"
	ApiV4WithCacheRedisResourceFile    = SamplesPath + "/apim/api_definition/v4/api-v4-with-cache-redis-resource.yml"
	ApiV4WithOAuth2GenericResourceFile = SamplesPath + "/apim/api_definition/v4/api-v4-with-oauth2-generic-resource.yml"
	ApiV4WithOauth2AmResourceFile      = SamplesPath + "/apim/api_definition/v4/api-v4-with-oauth2-am-resource.yml"
	ApiV4WithKeycloakAdapterFile       = SamplesPath + "/apim/api_definition/v4/api-v4-with-keycloak-adapter.yml"
	ApiV4WithLDAPAuthProviderFile      = SamplesPath + "/apim/api_definition/v4/api-v4-with-ldap-auth-provider.yml"
	ApiV4WithInlineAuthProviderFile    = SamplesPath + "/apim/api_definition/v4/api-v4-with-inline-auth-provider.yml"
	ApiV4WithHTTPAuthProviderFile      = SamplesPath + "/apim/api_definition/v4/api-v4-with-http-auth-provider.yml"
	ApiV4WithStateStopped              = SamplesPath + "/apim/api_definition/v4/api-v4-with-state-stopped.yml"
	ApiV4WithApiKeyPlanFile            = SamplesPath + "/apim/api_definition/v4/api-v4-with-api-key-plan.yml"
	ApiV4WithCacheRedisResourceRef     = SamplesPath + "/apim/api_definition/v4/api-v4-with-cache-redis-resource-ref.yml"
	ApiV4WithOAuth2GenericResRef       = SamplesPath + "/apim/api_definition/v4/api-v4-with-oauth2-generic-res-ref.yml"
	ApiV4WithOauth2AmResourceRefFile   = SamplesPath + "/apim/api_definition/v4/api-v4-with-oauth2-am-resource-ref.yml"
	ApiV4WithKeycloakAdapterRefFile    = SamplesPath + "/apim/api_definition/v4/api-v4-with-keycloak-adapter-ref.yml"
	ApiV4WithLDAPAuthProviderRefFile   = SamplesPath + "/apim/api_definition/v4/api-v4-with-ldap-auth-provider-ref.yml"
	ApiV4WithInlineAuthProviderRef     = SamplesPath + "/apim/api_definition/v4/api-v4-with-inline-auth-provider-ref.yml"
	ApiV4WithHTTPAuthProviderRefFile   = SamplesPath + "/apim/api_definition/v4/api-v4-with-http-auth-provider-ref.yml"
	ApiV4WithMarkdownPage              = SamplesPath + "/apim/api_definition/v4/api-v4-with-page-markdown.yml"
	ApiV4WithSwaggerHTTPFetcher        = SamplesPath + "/apim/api_definition/v4/api-v4-with-page-swagger-http-fetcher.yml"

	IngressPEMRegistry         = SamplesPath + "/ingress/ingress-pem-registry.yml"
	Ingress404ResponseTemplate = SamplesPath + "/ingress/ingress-response-404-config-map.yml"
	IngressWithoutTemplateFile = SamplesPath + "/ingress/ingress-without-api-template.yml"
	IngressWithTemplateFile    = SamplesPath + "/ingress/ingress-with-api-template.yml"
	IngressWithMultipleHosts   = SamplesPath + "/ingress/ingress-with-multiple-hosts.yml"
	IngressWithTLS             = SamplesPath + "/ingress/ingress-with-tls.yml"
	IngressWithTLSSecretFile   = SamplesPath + "/ingress/ingress-with-tls-secret.yml"
	IngressResponse404CMFile   = SamplesPath + "/ingress/ingress-response-404-config-map.yaml"

	Application = SamplesPath + "/apim/application/application.yml"
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
