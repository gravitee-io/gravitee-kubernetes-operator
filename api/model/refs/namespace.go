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

// +kubebuilder:object:generate=true
package refs

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/types/k8s/custom"
	"k8s.io/apimachinery/pkg/types"
)

var _ custom.ResourceRef = NamespacedName{}

type NamespacedName struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

func NewNamespacedName(namespace, name string) NamespacedName {
	return NamespacedName{Namespace: namespace, Name: name}
}

func (n NamespacedName) NamespacedName() types.NamespacedName {
	return types.NamespacedName{Namespace: n.Namespace, Name: n.Name}
}

func (n NamespacedName) String() string {
	return n.Namespace + "/" + n.Name
}
