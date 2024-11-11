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
	"encoding/json"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
)

// +kubebuilder:validation:Enum=HTTP;SUBSCRIPTION;TCP;KAFKA;
type ListenerType string

// +kubebuilder:validation:Enum=NONE;AUTO;AT_MOST_ONCE;AT_LEAST_ONCE;
type QosType string

const (
	HTTPListenerType         ListenerType = "HTTP"
	SubscriptionListenerType ListenerType = "SUBSCRIPTION"
	TCPListenerType          ListenerType = "TCP"
	KafkaType                ListenerType = "KAFKA"
)

const (
	AutoQOS        QosType = "AUTO"
	NoQOS          QosType = "NONE"
	AtMostOnceQOS  QosType = "AT_MOST_ONCE"
	AtLeastOnceQOS QosType = "AT_LEAST_ONCE"
)

type DLQ struct {
	// The endpoint to use when a message should be sent to the dead letter queue.
	// +kubebuilder:validation:Optional
	Endpoint *string `json:"endpoint,omitempty"`
}

type EntryPointType string

type GenericListener struct {
	*utils.GenericStringMap `json:",inline"`
}

func (l *GenericListener) UnmarshalJSON(data []byte) error {
	if l.GenericStringMap == nil {
		l.GenericStringMap = utils.NewGenericStringMap()
	}
	return l.GenericStringMap.UnmarshalJSON(data)
}

func (l *GenericListener) ListenerType() ListenerType {
	return ListenerType(l.GetString("type"))
}

func (l *GenericListener) ToListener() Listener {
	body, _ := json.Marshal(l)
	var listener Listener
	switch l.ListenerType() {
	case HTTPListenerType:
		listener = new(HttpListener)
	case SubscriptionListenerType:
		listener = new(SubscriptionListener)
	case TCPListenerType:
		listener = new(TCPListener)
	case KafkaType:
		listener = new(KafkaListener)
	}

	_ = json.Unmarshal(body, listener)
	return listener
}

func ToGenericListener(l Listener) *GenericListener {
	body, _ := json.Marshal(l)
	obj := new(GenericListener)
	_ = json.Unmarshal(body, obj)
	return obj
}

// +k8s:deepcopy-gen=false
type Listener interface {
	ListenerType() ListenerType
}

type AbstractListener struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default:=`HTTP`
	Type ListenerType `json:"type"`
	// +kubebuilder:validation:Required
	Entrypoints []*Entrypoint `json:"entrypoints"`
	// +kubebuilder:validation:Optional
	Servers []string `json:"servers"`
}

type HttpListener struct {
	*AbstractListener `json:",inline"`
	// +kubebuilder:validation:Required
	Paths []*Path `json:"paths"`
	// +kubebuilder:validation:Optional
	PathMappings []string `json:"pathMappings"`
}

func (l *AbstractListener) ToGatewayDefinition() *AbstractListener {
	listener := l.DeepCopy()
	listener.Type = ListenerType(Enum(l.Type).ToGatewayDefinition())
	ep := make([]*Entrypoint, len(l.Entrypoints))
	for i, l := range l.Entrypoints {
		ep[i] = l.ToGatewayDefinition()
	}
	listener.Entrypoints = ep

	return listener
}

func (l *HttpListener) ListenerType() ListenerType {
	return l.Type
}

func (l *HttpListener) ToGatewayDefinition() *HttpListener {
	listener := l.DeepCopy()
	listener.AbstractListener = l.AbstractListener.ToGatewayDefinition()

	return listener
}

type SubscriptionListener struct {
	*AbstractListener `json:",inline"`
}

func (l *SubscriptionListener) ToGatewayDefinition() *SubscriptionListener {
	listener := l.DeepCopy()
	listener.AbstractListener = l.AbstractListener.ToGatewayDefinition()

	return listener
}

func (l *SubscriptionListener) ListenerType() ListenerType {
	return l.Type
}

type TCPListener struct {
	*AbstractListener `json:",inline"`
	// +kubebuilder:validation:Required
	Hosts []string `json:"hosts"`
}

func (l *TCPListener) ToGatewayDefinition() *TCPListener {
	listener := l.DeepCopy()
	listener.AbstractListener = l.AbstractListener.ToGatewayDefinition()

	return listener
}

func (l *TCPListener) ListenerType() ListenerType {
	return l.Type
}

type KafkaListener struct {
	*AbstractListener `json:",inline"`
	// Kafka server hostname
	// +kubebuilder:validation:Required
	Host string `json:"host"`

	// Kafka server port number
	// +kubebuilder:validation:Required
	Port int `json:"port"`
}

func (l *KafkaListener) ToGatewayDefinition() *KafkaListener {
	listener := l.DeepCopy()
	listener.AbstractListener = l.AbstractListener.ToGatewayDefinition()

	return listener
}

func (l *KafkaListener) ListenerType() ListenerType {
	return l.Type
}

type Path struct {
	// +kubebuilder:validation:Optional
	Host string `json:"host,omitempty"`
	// +kubebuilder:validation:Required
	Path string `json:"path"`
}

type Entrypoint struct {
	// +kubebuilder:validation:Required
	Type string `json:"type"`
	// +kubebuilder:validation:Optional
	Qos *QosType `json:"qos,omitempty"`
	Dlq *DLQ     `json:"dlq,omitempty"`
	// +kubebuilder:validation:Optional
	Configuration *utils.GenericStringMap `json:"configuration,omitempty"`
}

func (ep *Entrypoint) ToGatewayDefinition() *Entrypoint {
	entryPoint := ep.DeepCopy()
	qosType := QosType(Enum(*ep.Qos).ToGatewayDefinition())
	entryPoint.Qos = &qosType
	if ep.Dlq != nil && ep.Dlq.Endpoint == nil {
		entryPoint.Dlq = nil
	}
	return entryPoint
}

func ToListenerGatewayDefinition(l Listener) *GenericListener {
	switch t := l.(type) {
	case *GenericListener:
		return ToListenerGatewayDefinition(t.ToListener())
	case *HttpListener:
		return ToGenericListener(t.ToGatewayDefinition())
	case *SubscriptionListener:
		return ToGenericListener(t.ToGatewayDefinition())
	case *TCPListener:
		return ToGenericListener(t.ToGatewayDefinition())
	case *KafkaListener:
		return ToGenericListener(t.ToGatewayDefinition())
	}

	return nil
}
