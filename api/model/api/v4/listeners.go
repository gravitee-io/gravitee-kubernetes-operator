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
	HTTPListenerType         = ListenerType("http")
	SubscriptionListenerType = ListenerType("subscription")
	TCPListenerType          = ListenerType("tcp")
)

type QOS string

const (
	AutoQOS        = QOS("auto")
	NoQOS          = QOS("none")
	AtMostOnceQOS  = QOS("at-most-once")
	AtLeastOnceQOS = QOS("at-least-once")
)

type DLQ struct {
	Endpoint string `json:"endpoint,omitempty"`
}

type EntryPointType string

const (
	EntryPointTypeHTTP = EntryPointType("http-proxy")
)

type Entrypoint struct {
	Type   string                  `json:"type,omitempty"`
	QOS    QOS                     `json:"qos,omitempty"`
	DLQ    *DLQ                    `json:"dlq,omitempty"`
	Config *utils.GenericStringMap `json:"configuration,omitempty"`
}

func NewHttpEntryPoint() Entrypoint {
	return Entrypoint{
		Type:   string(EntryPointTypeHTTP),
		QOS:    AutoQOS,
		Config: utils.NewGenericStringMap(),
	}
}

type Path struct {
	Host string `json:"host,omitempty"`
	// +kubebuilder:default:="/"
	Path           string `json:"path,omitempty"`
	OverrideAccess bool   `json:"overrideAccess,omitempty"`
}

func NewPath(host, path string) *Path {
	return &Path{
		Host: host,
		Path: path,
	}
}

func NewHttpListenerBase() *Listener {
	impl := utils.NewGenericStringMap()
	impl.Put("type", string(HTTPListenerType))
	impl.Put("entrypoints", []Entrypoint{NewHttpEntryPoint()})
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
