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

package apidefinition_test

import (
	"encoding/json"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/yaml"
)

// The API server does not default anything below spec.listeners (preserve-unknown-fields),
// so a stored KAFKA entrypoint without qos reaches ToGatewayDefinition with a nil Qos.
const nativeKafkaWithoutQos = `
apiVersion: gravitee.io/v1alpha1
kind: ApiV4Definition
metadata:
  name: e2e-kafka-through-gateway
  namespace: default
spec:
  name: e2e-kafka-through-gateway
  version: "1.0.0"
  type: NATIVE
  state: STARTED
  definitionContext:
    origin: KUBERNETES
    syncFrom: KUBERNETES
  listeners:
    - type: KAFKA
      host: orders
      port: 9092
      entrypoints:
        - type: native-kafka
  endpointGroups:
    - name: local-group
      type: native-kafka
      endpoints:
        - name: local-broker
          type: native-kafka
          inheritConfiguration: false
          configuration:
            bootstrapServers: "gck-kafka.gravitee.svc.cluster.local:9092"
  plans:
    KeyLess:
      name: Free plan
      status: PUBLISHED
      security:
        type: KEY_LESS
`

var _ = Describe("ToGatewayDefinition with a KAFKA entrypoint without qos", func() {
	It("defaults the entrypoint qos to AUTO instead of panicking", func() {
		api := new(v1alpha1.ApiV4Definition)
		Expect(yaml.Unmarshal([]byte(nativeKafkaWithoutQos), api)).To(Succeed())

		var def v4.GatewayDefinitionApi
		Expect(func() { def = api.ToGatewayDefinition() }).NotTo(Panic())

		Expect(def.Listeners).To(HaveLen(1))
		body, err := json.Marshal(def.Listeners[0])
		Expect(err).NotTo(HaveOccurred())
		listener := new(v4.KafkaListener)
		Expect(json.Unmarshal(body, listener)).To(Succeed())
		Expect(listener.Entrypoints).To(HaveLen(1))
		Expect(listener.Entrypoints[0].Qos).NotTo(BeNil())
		Expect(*listener.Entrypoints[0].Qos).To(Equal(v4.QosType(v4.Enum(v4.AutoQOS).ToGatewayDefinition())))
	})
})
