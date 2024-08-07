// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package k8s

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	netV1 "k8s.io/api/networking/v1"
)

func IsGraviteeIngress(ingress *netV1.Ingress) bool {
	var ingressClassName string
	if ingressClassName = ingress.GetAnnotations()[core.IngressClassAnnotation]; ingress.Spec.IngressClassName != nil {
		ingressClassName = *(ingress.Spec.IngressClassName)
	}

	for _, ingressClass := range env.Config.IngressClasses {
		if ingressClassName == ingressClass {
			return true
		}
	}
	return false
}
