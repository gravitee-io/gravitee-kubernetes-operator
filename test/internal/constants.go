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

import "time"

const (
	GatewayUrl                   = "http://localhost:9000/gateway"
	SamplesPath                  = "../config/samples"
	ContextWithSecretFile        = SamplesPath + "/context/dev/management-context-with-secret-ref.yaml"
	BasicApiFile                 = SamplesPath + "/apim/basic-api.yml"
	ApiWithDisabledHCFile        = SamplesPath + "/apim/api-with-health-check-disabled.yml"
	ApiWithHCFile                = SamplesPath + "/apim/api-with-health-check.yml"
	ApiWithServiceDiscoveryFile  = SamplesPath + "/apim/api-with-service-discovery.yml"
	ApiWithMetadataFile          = SamplesPath + "/apim/api-with-metadata.yml"
	ApiWithEndpointGroupsFile    = SamplesPath + "/apim/api-with-endpoint-groups.yml"
	ApiWithLoggingFile           = SamplesPath + "/apim/api-with-logging.yml"
	ApiWithApiKeyPlanFile        = SamplesPath + "/apim/api-with-api-key-plan.yml"
	ApiWithCacheResource         = SamplesPath + "/apim/api-with-cache-resource.yml"
	ApiWithOAuth2GenericResource = SamplesPath + "/apim/api-with-oauth2-generic-resource.yml"

	contextWithCredentialsFile = SamplesPath + "/context/dev/management-context-with-credentials.yaml"
	apimClientTimeout          = 5 * time.Second
)
