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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/kafka"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
)

var (
	kafkaACLPolicyName = "kafka-acl"

	resourceMap = map[kafka.KafkaAcccessControlResourceType]string{
		kafka.KafkaAcccessControlTypeCluster:             "CLUSTER",
		kafka.KafkaAcccessControlTypeTopic:               "TOPIC",
		kafka.KafkaAcccessControlTypeGroup:               "GROUP",
		kafka.KafkaAcccessControlTransactionalIdentifier: "TRANSACTIONAL_ID",
	}

	opsMap = map[kafka.KafkaAccessControlOperation]string{
		kafka.KafkaAccessControlOperationCreate:          "CREATE",
		kafka.KafkaAccessControlOperationRead:            "READ",
		kafka.KafkaAccessControlOperationWrite:           "WRITE",
		kafka.KafkaAccessControlOperationDelete:          "DELETE",
		kafka.KafkaAccessControlOperationAlter:           "ALTER",
		kafka.KafkaAccessControlOperationAlterConfigs:    "ALTER_CONFIGS",
		kafka.KafkaAccessControlOperationDescribe:        "DESCRIBE",
		kafka.KafkaAccessControlOperationDescribeConfigs: "DESCRIBE_CONFIGS",
		kafka.KafkaAccessControlOperationClusterAction:   "CLUSTER_ACTION",
	}

	matchMap = map[kafka.KafkaAccessControlResourceMatchType]string{
		kafka.KafkaResourceMatchTypeExact:             "LITERAL",
		kafka.KafkaResourceMatchTypePrefix:            "PREFIXED",
		kafka.KafkaResourceMatchTypeRegularExpression: "MATCH",
	}
)

func buildACL(acl kafka.KafkaACLFilter) *v4.FlowStep {
	return v4.NewFlowStep(base.FlowStep{
		Policy:  &kafkaACLPolicyName,
		Enabled: true,
		Configuration: utils.NewGenericStringMap().
			Put("authorizations", mapKafkaAccessControlRules(acl)),
	})
}

func mapKafkaAccessControlRules(acl kafka.KafkaACLFilter) []any {
	authz := make([]any, len(acl.Rules))
	for i := range acl.Rules {
		authz[i] = mapAccessControlRule(acl.Rules[i].Resources)
	}
	return authz
}

func mapAccessControlRule(acl []kafka.KafkaAccessControl) any {
	resources := make([]any, len(acl))
	for i := range acl {
		resources[i] = mapAccessControl(acl[i])
	}
	return utils.NewGenericStringMap().Put("resources", resources)
}

func mapAccessControl(accessControl kafka.KafkaAccessControl) map[string]any {
	ac := map[string]any{}
	ac["type"] = resourceMap[accessControl.Type]
	ac["operations"] = mapOperations(accessControl)
	if accessControl.Type == kafka.KafkaAcccessControlTypeCluster {
		return ac
	}
	if accessControl.Match == nil {
		ac["resourcePatternType"] = "ANY"
		return ac
	}
	ac["resourcePatternType"] = matchMap[accessControl.Match.Type]
	ac["resourcePattern"] = accessControl.Match.Value
	return ac
}

func mapOperations(
	accessControl kafka.KafkaAccessControl,
) []any {
	resourceType := accessControl.Type
	ops := make([]any, len(accessControl.Operations))
	for i := range accessControl.Operations {
		op := accessControl.Operations[i]
		ops[i] = resourceMap[resourceType] + "_" + opsMap[op]
	}
	return ops
}
