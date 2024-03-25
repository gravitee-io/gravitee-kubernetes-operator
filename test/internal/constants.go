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

package internal

const (
	Namespace = "default"

	GatewayUrl    = "http://localhost:30082"
	ManagementUrl = "http://localhost:30083/management"
	SamplesPath   = "../examples"

	ContextSecretFile          = SamplesPath + "/management_context/dev/management-context-secret.yml"
	ContextWithSecretFile      = SamplesPath + "/management_context/dev/management-context-with-secret-ref.yml"
	ContextWithCredentialsFile = SamplesPath + "/management_context/dev/management-context-with-credentials.yml"

	BasicApiFile                        = SamplesPath + "/apim/api_definition/basic-api.yml"
	BasicApiWithRateLimit               = SamplesPath + "/apim/api_definition/basic-api-with-rate-limit.yml"
	BasicApiWithDisabledPolicy          = SamplesPath + "/apim/api_definition/basic-api-with-disabled-policy.yml"
	ApiWithTemplatingFile               = SamplesPath + "/apim/api_definition/api-with-templating.yml"
	ApiWithTemplatingSecretFile         = SamplesPath + "/apim/api_definition/api-with-templating-secret.yml"
	ApiWithTemplatingConfigMapFile      = SamplesPath + "/apim/api_definition/api-with-templating-config-map.yml"
	ExportedApi                         = SamplesPath + "/apim/api_definition/exported-api.yml"
	ApiWithContextFile                  = SamplesPath + "/apim/api_definition/api-with-context.yml"
	ApiWithContextNoPlanFile            = SamplesPath + "/apim/api_definition/api-with-no-plan.yml"
	ApiWithDisabledHCFile               = SamplesPath + "/apim/api_definition/api-with-health-check-disabled.yml"
	ApiWithHCFile                       = SamplesPath + "/apim/api_definition/api-with-health-check.yml"
	ApiWithServiceDiscoveryFile         = SamplesPath + "/apim/api_definition/api-with-service-discovery.yml"
	ApiWithMetadataFile                 = SamplesPath + "/apim/api_definition/api-with-metadata.yml"
	ApiWithEndpointGroupsFile           = SamplesPath + "/apim/api_definition/api-with-endpoint-groups.yml"
	ApiWithLoggingFile                  = SamplesPath + "/apim/api_definition/api-with-logging.yml"
	ApiWithApiKeyPlanFile               = SamplesPath + "/apim/api_definition/api-with-api-key-plan.yml"
	ApiWithCacheResourceFile            = SamplesPath + "/apim/api_definition/api-with-cache-resource.yml"
	ApiWithCacheResourceRefFile         = SamplesPath + "/apim/api_definition/api-with-cache-resource-ref.yml"
	ApiWithCacheRedisResourceFile       = SamplesPath + "/apim/api_definition/api-with-cache-redis-resource.yml"
	ApiWithCacheRedisResourceRefFile    = SamplesPath + "/apim/api_definition/api-with-cache-redis-resource-ref.yml"
	ApiWithOAuth2GenericResourceFile    = SamplesPath + "/apim/api_definition/api-with-oauth2-generic-resource.yml"
	ApiWithOAuth2GenericResourceRefFile = SamplesPath + "/apim/api_definition/api-with-oauth2-generic-resource-ref.yml"
	ApiWithOauth2AmResourceFile         = SamplesPath + "/apim/api_definition/api-with-oauth2-am-resource.yml"
	ApiWithOauth2AmResourceRefFile      = SamplesPath + "/apim/api_definition/api-with-oauth2-am-resource-ref.yml"
	ApiWithKeycloakAdapterFile          = SamplesPath + "/apim/api_definition/api-with-keycloak-adapter.yml"
	ApiWithKeycloakAdapterRefFile       = SamplesPath + "/apim/api_definition/api-with-keycloak-adapter-ref.yml"
	ApiWithLDAPAuthProviderFile         = SamplesPath + "/apim/api_definition/api-with-ldap-auth-provider.yml"
	ApiWithLDAPAuthProviderRefFile      = SamplesPath + "/apim/api_definition/api-with-ldap-auth-provider-ref.yml"
	ApiWithInlineAuthProviderFile       = SamplesPath + "/apim/api_definition/api-with-inline-auth-provider.yml"
	ApiWithInlineAuthProviderRefFile    = SamplesPath + "/apim/api_definition/api-with-inline-auth-provider-ref.yml"
	ApiWithHTTPAuthProviderFile         = SamplesPath + "/apim/api_definition/api-with-http-auth-provider.yml"
	ApiWithHTTPAuthProviderRefFile      = SamplesPath + "/apim/api_definition/api-with-http-auth-provider-ref.yml"

	ApiResourceCacheFile              = SamplesPath + "/apim/api_resource/api-resource-cache.yml"
	ApiResourceCacheRedisFile         = SamplesPath + "/apim/api_resource/api-resource-cache-redis.yml"
	ApiResourceHTTPAuthProviderFile   = SamplesPath + "/apim/api_resource/api-resource-http-auth-provider.yml"
	ApiResourceInlineAuthProviderFile = SamplesPath + "/apim/api_resource/api-resource-inline-auth-provider.yml"
	ApiResourceLDAPAuthProviderFile   = SamplesPath + "/apim/api_resource/api-resource-ldap-auth-provider.yml"
	ApiResourceKeycloakAdapterFile    = SamplesPath + "/apim/api_resource/api-resource-keycloak-adapter.yml"
	ApiResourceOauth2AMFile           = SamplesPath + "/apim/api_resource/api-resource-oauth2-am.yml"
	ApiResourceOauth2GenericFile      = SamplesPath + "/apim/api_resource/api-resource-oauth2-generic.yml"

	IngressWithoutTemplateFile    = SamplesPath + "/ingress/ingress-without-api-template.yml"
	ApiTemplateWithApiKeyPlanFile = SamplesPath + "/ingress/api-template-with-api-key-plan.yml"
	IngressWithTemplateFile       = SamplesPath + "/ingress/ingress-with-api-template.yml"
	IngressWithMultipleHosts      = SamplesPath + "/ingress/ingress-with-multiple-hosts.yml"
	IngressWithTLS                = SamplesPath + "/ingress/ingress-with-tls.yml"
	IngressWithTLSSecretFile      = SamplesPath + "/ingress/ingress-with-tls-secret.yml"
	IngressResponse404CMFile      = SamplesPath + "/ingress/ingress-response-404-config-map.yaml"

	BasicApplication = SamplesPath + "/apim/application/basic-application.yml"
)
