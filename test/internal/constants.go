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
	GatewayUrl = "http://localhost:9000/gateway"

	SamplesPath = "../config/samples"

	ContextWithSecretFile      = SamplesPath + "/context/dev/management-context-with-secret-ref.yml"
	ContextWithCredentialsFile = SamplesPath + "/context/dev/management-context-with-credentials.yml"

	BasicApiFile                        = SamplesPath + "/apim/basic-api.yml"
	ApiWithContextFile                  = SamplesPath + "/apim/api-with-context.yml"
	ApiWithContextNoPlanFile            = SamplesPath + "/apim/api-with-no-plan.yml"
	ApiWithDisabledHCFile               = SamplesPath + "/apim/api-with-health-check-disabled.yml"
	ApiWithHCFile                       = SamplesPath + "/apim/api-with-health-check.yml"
	ApiWithServiceDiscoveryFile         = SamplesPath + "/apim/api-with-service-discovery.yml"
	ApiWithMetadataFile                 = SamplesPath + "/apim/api-with-metadata.yml"
	ApiWithEndpointGroupsFile           = SamplesPath + "/apim/api-with-endpoint-groups.yml"
	ApiWithLoggingFile                  = SamplesPath + "/apim/api-with-logging.yml"
	ApiWithApiKeyPlanFile               = SamplesPath + "/apim/api-with-api-key-plan.yml"
	ApiWithCacheResourceFile            = SamplesPath + "/apim/api-with-cache-resource.yml"
	ApiWithCacheResourceRefFile         = SamplesPath + "/apim/api-with-cache-resource-ref.yml"
	ApiWithCacheRedisResourceFile       = SamplesPath + "/apim/api-with-cache-redis-resource.yml"
	ApiWithCacheRedisResourceRefFile    = SamplesPath + "/apim/api-with-cache-redis-resource-ref.yml"
	ApiWithOAuth2GenericResourceFile    = SamplesPath + "/apim/api-with-oauth2-generic-resource.yml"
	ApiWithOAuth2GenericResourceRefFile = SamplesPath + "/apim/api-with-oauth2-generic-resource-ref.yml"
	ApiWithOauth2AmResourceFile         = SamplesPath + "/apim/api-with-oauth2-am-resource.yml"
	ApiWithOauth2AmResourceRefFile      = SamplesPath + "/apim/api-with-oauth2-am-resource-ref.yml"
	ApiWithKeycloakAdapterFile          = SamplesPath + "/apim/api-with-keycloak-adapter.yml"
	ApiWithKeycloakAdapterRefFile       = SamplesPath + "/apim/api-with-keycloak-adapter-ref.yml"
	ApiWithLDAPAuthProviderFile         = SamplesPath + "/apim/api-with-ldap-auth-provider.yml"
	ApiWithLDAPAuthProviderRefFile      = SamplesPath + "/apim/api-with-ldap-auth-provider-ref.yml"
	ApiWithInlineAuthProviderFile       = SamplesPath + "/apim/api-with-inline-auth-provider.yml"
	ApiWithInlineAuthProviderRefFile    = SamplesPath + "/apim/api-with-inline-auth-provider-ref.yml"
	ApiWithHTTPAuthProviderFile         = SamplesPath + "/apim/api-with-http-auth-provider.yml"
	ApiWithHTTPAuthProviderRefFile      = SamplesPath + "/apim/api-with-http-auth-provider-ref.yml"
	ApiResourceCacheFile                = SamplesPath + "/apim/api-resource-cache.yml"
	ApiResourceCacheRedisFile           = SamplesPath + "/apim/api-resource-cache-redis.yml"
	ApiResourceHTTPAuthProviderFile     = SamplesPath + "/apim/api-resource-http-auth-provider.yml"
	ApiResourceInlineAuthProviderFile   = SamplesPath + "/apim/api-resource-inline-auth-provider.yml"
	ApiResourceLDAPAuthProviderFile     = SamplesPath + "/apim/api-resource-ldap-auth-provider.yml"
	ApiResourceKeycloakAdapterFile      = SamplesPath + "/apim/api-resource-keycloak-adapter.yml"
	ApiResourceOauth2AMFile             = SamplesPath + "/apim/api-resource-oauth2-am.yml"
	ApiResourceOauth2GenericFile        = SamplesPath + "/apim/api-resource-oauth2-generic.yml"
	IngressWithoutTemplateFile          = SamplesPath + "/ingress/ingress-without-api-template.yml"
	IngressWithMultipleHosts            = SamplesPath + "/ingress/ingress-with-multiple-hosts.yml"
)
