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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func buildListeners(
	route *v1alpha1.KafkaRoute,
	gw *gwAPIv1.Gateway,
) []*v4.GenericListener {
	gwListener := k8s.GetKafkaListener(gw)
	listener := v4.NewKafkaListener(string(*route.Spec.Hostname), int(gwListener.Port))
	return []*v4.GenericListener{
		v4.ToGenericListener(listener),
	}
}
