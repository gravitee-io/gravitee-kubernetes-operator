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

package v4

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
)

type ListenerType string

const (
	HTTPListenerType         = "HTTP"
	SubscriptionListenerType = "SUBSCRIPTION"
	TCPListenerType          = "TCP"
)

const (
	AutoQOS        = "AUTO"
	NoQOS          = "NONE"
	AtMostOnceQOS  = "AT_MOST_ONCE"
	AtLeastOnceQOS = "AT_LEAST_ONCE"
)

type DLQ struct {
	Endpoint string `json:"endpoint,omitempty"`
}

type EntryPointType string

const (
	EntryPointTypeHTTP = EntryPointType("http-proxy")
)

func NewHttpEntryPoint() map[string]interface{} {
	return map[string]interface{}{
		"type":          string(EntryPointTypeHTTP),
		"qos":           string(AutoQOS),
		"configuration": map[string]interface{}{},
	}
}

func NewPath(host, path string) map[string]interface{} {
	out := map[string]interface{}{"path": path}
	if host != "" {
		out["host"] = host
	}
	return out
}

func NewHttpListenerBase() *Listener {
	impl := utils.NewGenericStringMap()
	impl.Put("type", string(HTTPListenerType))
	impl.Put("entrypoints", []interface{}{NewHttpEntryPoint()})
	return &Listener{impl}
}

type Listener struct {
	*utils.GenericStringMap `json:",inline"`
}

func (l *Listener) UnmarshalJSON(data []byte) error {
	if l.GenericStringMap == nil {
		l.GenericStringMap = utils.NewGenericStringMap()
	}
	return l.GenericStringMap.UnmarshalJSON(data)
}

func (l *Listener) ToGatewayDefinition() *Listener {
	listener := l.DeepCopy()
	listener.Put("type", Enum(listener.GetString("type")).ToGatewayDefinition())
	listener.Put("entrypoints", listener.GetGatewayDefinitionEntryPoints())
	return listener
}

func (l *Listener) GetGatewayDefinitionEntryPoints() []interface{} {
	var entrypoints []interface{}
	for _, entrypoint := range l.GetSlice("entrypoints") {
		entrypoint := utils.ToGenericStringMap(entrypoint)
		qos := Enum(entrypoint.GetString("qos")).ToGatewayDefinition()
		entrypoint.Put("qos", Enum(qos).ToGatewayDefinition())
		entrypoints = append(entrypoints, entrypoint)
	}
	return entrypoints
}
