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

package kafka

import gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"

// +kubebuilder:validation:Enum:=Topic;Cluster;Group;TransactionalIdentifier;
type KafkaAcccessControlResourceType string

const (
	KafkaAcccessControlTypeTopic               = KafkaAcccessControlResourceType("Topic")
	KafkaAcccessControlTypeCluster             = KafkaAcccessControlResourceType("Cluster")
	KafkaAcccessControlTypeGroup               = KafkaAcccessControlResourceType("Group")
	KafkaAcccessControlTransactionalIdentifier = KafkaAcccessControlResourceType("TransactionalIdentifier")
)

// +kubebuilder:validation:Enum:=Create;Read;Write;Delete;Alter;AlterConfigs;Describe;DescribeConfigs;ClusterAction;
type KafkaAccessControlOperation string

const (
	KafkaAccessControlOperationCreate          = KafkaAccessControlOperation("Create")
	KafkaAccessControlOperationRead            = KafkaAccessControlOperation("Read")
	KafkaAccessControlOperationWrite           = KafkaAccessControlOperation("Write")
	KafkaAccessControlOperationDelete          = KafkaAccessControlOperation("Delete")
	KafkaAccessControlOperationAlter           = KafkaAccessControlOperation("Alter")
	KafkaAccessControlOperationAlterConfigs    = KafkaAccessControlOperation("AlterConfigs")
	KafkaAccessControlOperationDescribe        = KafkaAccessControlOperation("Describe")
	KafkaAccessControlOperationDescribeConfigs = KafkaAccessControlOperation("DescribeConfigs")
	KafkaAccessControlOperationClusterAction   = KafkaAccessControlOperation("ClusterAction")
)

// +kubebuilder:validation:Enum:=Exact;Prefix;RegularExpression;
type KafkaAccessControlResourceMatchType string

const (
	KafkaResourceMatchTypeExact             = KafkaAccessControlResourceMatchType("Exact")
	KafkaResourceMatchTypePrefix            = KafkaAccessControlResourceMatchType("Prefix")
	KafkaResourceMatchTypeRegularExpression = KafkaAccessControlResourceMatchType("RegularExpression")
)

type KafkaAccessControlMatch struct {
	// Type specifies the semantics of how the resource should be matched.

	// Valid PathMatchType values, along with their support levels, are:
	//
	// * "Exact" Resources whose name is an exact match to the specified string receive the ACL.
	// * "Prefix" Resources whose name starts with the specified string receive the ACL.
	// * "RegularExpression" Resources that match the specified expression receive the ACL.
	Type KafkaAccessControlResourceMatchType `json:"type"`
	// Value of the resource to match against.
	Value string `json:"value"`
}

type KafkaAccessControl struct {
	Type KafkaAcccessControlResourceType `json:"type"`
	// Operations specifies the set of operations / verbs to allow for the resource
	// under access control.
	Operations []KafkaAccessControlOperation `json:"operations"`
	//+optional
	// Match describes how to select the resource that will be subject to the access control.
	// If not specified, any resource will be matched.
	Match *KafkaAccessControlMatch `json:"match,omitempty"`
}

type KafkaAccessControlRules struct {
	// A resource group together a type of matched resource and a set of operations
	// to be granted by the access control for that resource.
	// +kubebuilder:validation:MinItems=1
	Resources []KafkaAccessControl `json:"resources"`
	// Options allow to specify implementation specific behaviours
	// for a set of rules.
	// +optional
	// +kubebuilder:validation:MaxProperties=16
	Options map[gwAPIv1.AnnotationKey]gwAPIv1.AnnotationValue `json:"options,omitempty"`
}

type KafkaACLFilter struct {
	// Rules define a set of rules that can be use to group a set of resources together with
	// access control rules to be applied.
	// ACLs are restrictive because once they are applied,
	// proxy clients must be authorized to perform the actions they are taking.
	// If there is no ACL defined for the action taken by the user, the action is prohibited.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	Rules []KafkaAccessControlRules `json:"rules"`
}
