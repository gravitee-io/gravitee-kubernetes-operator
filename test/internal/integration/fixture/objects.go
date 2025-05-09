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

package fixture

import (
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	coreV1 "k8s.io/api/core/v1"
	netV1 "k8s.io/api/networking/v1"
)

type Objects struct {
	Secrets    []*coreV1.Secret
	ConfigMaps []*coreV1.ConfigMap

	Ingress *netV1.Ingress

	Context           *v1alpha1.ManagementContext
	Resource          *v1alpha1.ApiResource
	API               *v1alpha1.ApiDefinition
	APIv4             *v1alpha1.ApiV4Definition
	Application       *v1alpha1.Application
	Subscription      *v1alpha1.Subscription
	Group             *v1alpha1.Group
	SharedPolicyGroup *v1alpha1.SharedPolicyGroup
	Notification      *v1alpha1.Notification

	randomSuffix string
}

func (o *Objects) GetGeneratedSuffix() string {
	return o.randomSuffix
}

func (o *Objects) GetIngressPEMRegistryKey() string {
	if o.Ingress == nil {
		return ""
	}
	return fmt.Sprintf("%s-%s", o.Ingress.Namespace, o.Ingress.Name)
}
