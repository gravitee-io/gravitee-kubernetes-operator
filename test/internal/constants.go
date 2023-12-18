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

	GatewayUrl    = "http://localhost:9001"
	ManagementUrl = "http://localhost:9000/management"
	SamplesPath   = "../samples"

	ClusterContextFile = SamplesPath + "/managementcontext/cluster/management-context-with-secret-ref.yml"
	ClientContextFile  = SamplesPath + "/managementcontext/dev/management-context-with-credentials.yml"

	BasicApiFile                        = SamplesPath + "/apidefinition/v1alpha1/api.yml"
	BasicApiWithRateLimit               = SamplesPath + "/apidefinition/v1alpha1/api-with-rate-limit.yml"
	BasicApiWithDisabledPolicy          = SamplesPath + "/apidefinition/v1alpha1/api-with-disabled-policy.yml"
	BasicApiFileWithTemplate            = SamplesPath + "/apidefinition/v1alpha1/basic-api-with-template.yml"
	BasicApiFileTemplating              = SamplesPath + "/apidefinition/v1alpha1/api-with-templating.yml"
	ExportedApi                         = SamplesPath + "/apidefinition/v1alpha1/api-with-ids.yml"
	ApiWithContextFile                  = SamplesPath + "/apidefinition/v1alpha1/api-with-context.yml"
	ApiWithContextNoPlanFile            = SamplesPath + "/apidefinition/v1alpha1/api-with-no-plan.yml"
	ApiWithDisabledHCFile               = SamplesPath + "/apidefinition/v1alpha1/api-with-health-check-disabled.yml"
	ApiWithHCFile                       = SamplesPath + "/apidefinition/v1alpha1/api-with-health-check.yml"
	ApiWithServiceDiscoveryFile         = SamplesPath + "/apidefinition/v1alpha1/api-with-service-discovery.yml"
	ApiWithMetadataFile                 = SamplesPath + "/apidefinition/v1alpha1/api-with-metadata.yml"
	ApiWithEndpointGroupsFile           = SamplesPath + "/apidefinition/v1alpha1/api-with-endpoint-groups.yml"
	ApiWithLoggingFile                  = SamplesPath + "/apidefinition/v1alpha1/api-with-logging.yml"
	ApiWithApiKeyPlanFile               = SamplesPath + "/apidefinition/v1alpha1/api-with-api-key-plan.yml"
	ApiWithCacheResourceFile            = SamplesPath + "/apidefinition/v1alpha1/api-with-cache-resource.yml"
	ApiWithCacheResourceRefFile         = SamplesPath + "/apidefinition/v1alpha1/api-with-cache-resource-ref.yml"
	ApiWithCacheRedisResourceFile       = SamplesPath + "/apidefinition/v1alpha1/api-with-cache-redis-resource.yml"
	ApiWithCacheRedisResourceRefFile    = SamplesPath + "/apidefinition/v1alpha1/api-with-cache-redis-resource-ref.yml"
	ApiWithOAuth2GenericResourceFile    = SamplesPath + "/apidefinition/v1alpha1/api-with-oauth2-generic-resource.yml"
	ApiWithOAuth2GenericResourceRefFile = SamplesPath + "/apidefinition/v1alpha1/api-with-oauth2-generic-resource-ref.yml"
	ApiWithOauth2AmResourceFile         = SamplesPath + "/apidefinition/v1alpha1/api-with-oauth2-am-resource.yml"
	ApiWithOauth2AmResourceRefFile      = SamplesPath + "/apidefinition/v1alpha1/api-with-oauth2-am-resource-ref.yml"
	ApiWithKeycloakAdapterFile          = SamplesPath + "/apidefinition/v1alpha1/api-with-keycloak-adapter.yml"
	ApiWithKeycloakAdapterRefFile       = SamplesPath + "/apidefinition/v1alpha1/api-with-keycloak-adapter-ref.yml"
	ApiWithLDAPAuthProviderFile         = SamplesPath + "/apidefinition/v1alpha1/api-with-ldap-auth-provider.yml"
	ApiWithLDAPAuthProviderRefFile      = SamplesPath + "/apidefinition/v1alpha1/api-with-ldap-auth-provider-ref.yml"
	ApiWithInlineAuthProviderFile       = SamplesPath + "/apidefinition/v1alpha1/api-with-inline-auth-provider.yml"
	ApiWithInlineAuthProviderRefFile    = SamplesPath + "/apidefinition/v1alpha1/api-with-inline-auth-provider-ref.yml"
	ApiWithHTTPAuthProviderFile         = SamplesPath + "/apidefinition/v1alpha1/api-with-http-auth-provider.yml"
	ApiWithHTTPAuthProviderRefFile      = SamplesPath + "/apidefinition/v1alpha1/api-with-http-auth-provider-ref.yml"

	ApiResourceCacheFile              = SamplesPath + "/apiresource/api-resource-cache.yml"
	ApiResourceCacheRedisFile         = SamplesPath + "/apiresource/api-resource-cache-redis.yml"
	ApiResourceHTTPAuthProviderFile   = SamplesPath + "/apiresource/api-resource-http-auth-provider.yml"
	ApiResourceInlineAuthProviderFile = SamplesPath + "/apiresource/api-resource-inline-auth-provider.yml"
	ApiResourceLDAPAuthProviderFile   = SamplesPath + "/apiresource/api-resource-ldap-auth-provider.yml"
	ApiResourceKeycloakAdapterFile    = SamplesPath + "/apiresource/api-resource-keycloak-adapter.yml"
	ApiResourceOauth2AMFile           = SamplesPath + "/apiresource/api-resource-oauth2-am.yml"
	ApiResourceOauth2GenericFile      = SamplesPath + "/apiresource/api-resource-oauth2-generic.yml"

	ApiTemplateWithApiKeyPlanFile = SamplesPath + "/apitemplate/api-template-with-api-key-plan.yml"

	IngressWithoutTemplateFile = SamplesPath + "/ingress/ingress-without-api-template.yml"
	IngressWithTemplateFile    = SamplesPath + "/ingress/ingress-with-api-template.yml"
	IngressWithMultipleHosts   = SamplesPath + "/ingress/ingress-with-multiple-hosts.yml"
	IngressWithTLS             = SamplesPath + "/ingress/ingress-with-tls.yml"

	BasicApplication = SamplesPath + "/application/basic-application.yml"
)
