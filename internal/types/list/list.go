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

package list

import (
	"fmt"

	v1 "k8s.io/api/core/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	netv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func OfType(obj interface{}) (client.ObjectList, error) {
	switch obj.(type) {
	case *v1alpha1.ApiDefinitionList:
		return &v1alpha1.ApiDefinitionList{}, nil
	case *v1alpha1.ApiDefinition:
		return &v1alpha1.ApiDefinitionList{}, nil
	case *netv1.IngressList:
		return &netv1.IngressList{}, nil
	case *v1.SecretList:
		return &v1.SecretList{}, nil
	default:
		return nil, fmt.Errorf("unknown type %T", obj)
	}
}
