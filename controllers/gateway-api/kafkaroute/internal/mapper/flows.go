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
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/kafka"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

var (
	kafkaFlowName = "filters"
)

func buildFlows(
	route *v1alpha1.KafkaRoute,
) []*v4.Flow {
	flow := buildFlow(route)
	if len(flow.Interact) > 0 {
		return []*v4.Flow{flow}
	}
	return []*v4.Flow{}
}

func buildFlow(
	route *v1alpha1.KafkaRoute,
) *v4.Flow {
	return &v4.Flow{
		Name:     &kafkaFlowName,
		Interact: buildInteract(route.Spec.Filters),
		Enabled:  true,
	}
}

func buildInteract(filters []kafka.KafkaRouteFilter) []*v4.FlowStep {
	steps := []*v4.FlowStep{}
	for _, filter := range filters {
		if filter.Type == kafka.KafkaRouteFilterTypeAccessControlList {
			steps = append(steps, buildACL(filter.ACL))
		}
	}
	return steps
}
