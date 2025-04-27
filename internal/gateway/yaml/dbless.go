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

package yaml

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"k8s.io/apimachinery/pkg/util/runtime"
)

var DBLess = utils.NewGenericStringMap()

func init() {
	runtime.Must(unmarshalDBLess())
}

func unmarshalDBLess() error {
	return DBLess.UnmarshalYAML(
		[]byte(`
servers: []

http: {}

management:
  type: none

ratelimit:
  type: none

reporters:
  elasticsearch:
    enabled: false

services:
  core:
    http:
      enabled: true
      port: 18082
      host: 0.0.0.0
      authentication:
        type: basic
        users:
          admin: admin
      secured: false

  sync: 
    enabled: true
    kubernetes:
      enabled: true

  monitoring:
    delay: 5000
    unit: MILLISECONDS

  heartbeat: 
    delay: 5000
    enabled: true
    unit: MILLISECONDS

  metrics:
    enabled: false
    prometheus:
      enabled: true

api:
  validateSubscription: false
  allowOverlappingContext: true

gracefulShutdown:
  delay: 0
  unit: MILLISECONDS

secrets:
  kubernetes:
    enabled: true
`))
}
