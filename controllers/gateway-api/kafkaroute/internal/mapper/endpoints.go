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

package mapper

import (
	"fmt"
	"strings"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/kafka"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

const (
	serviceBootstrapServerPattern = "%s.%s.svc.cluster.local:%d"
	endpointName                  = "native-kafka"
)

var plainTextSecurity = map[string]any{
	"protocol": "PLAINTEXT",
}

func buildEndpointGroups(
	route *v1alpha1.KafkaRoute,
) []*v4.EndpointGroup {
	return []*v4.EndpointGroup{buildEndpointGroup(route)}
}

func buildEndpointGroup(route *v1alpha1.KafkaRoute) *v4.EndpointGroup {
	group := v4.NewKafkaEndpointGroup(endpointName)
	group.Endpoints = []*v4.Endpoint{buildEndpoint(route)}
	group.SharedConfig.Put("security", plainTextSecurity)
	return group
}

func buildEndpoint(route *v1alpha1.KafkaRoute) *v4.Endpoint {
	endpoint := v4.NewKafkaEndpoint(endpointName)
	endpoint.Config.Put("bootstrapServers", buildBootrapServers(route.Spec.BackendRefs))
	endpoint.Config.Put("weight", 1)
	endpoint.Config.Put("secondary", false)
	endpoint.Inherit = true
	return endpoint
}

func buildBootrapServers(backendRefs []kafka.KafkaBackendRef) string {
	return strings.Join(getBootsrapServers(backendRefs), ",")
}

func getBootsrapServers(backendRefs []kafka.KafkaBackendRef) []string {
	servers := make([]string, len(backendRefs))
	for i := range backendRefs {
		servers[i] = buildBootsrapServer(backendRefs[i])
	}
	return servers
}

func buildBootsrapServer(backendRef kafka.KafkaBackendRef) string {
	return fmt.Sprintf(
		serviceBootstrapServerPattern,
		backendRef.Name,
		string(*backendRef.Namespace),
		*backendRef.Port,
	)
}
