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

package utils

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

type Event struct {
	k8sEventRecorder record.EventRecorder
}

func NewEvent(recorder record.EventRecorder) *Event {
	return &Event{
		recorder,
	}
}

// K8s Normal event
// 'reason' is the reason this event is generated. 'reason' should be short and unique; it
// should be in UpperCamelCase format (starting with a capital letter). "reason" will be used
// to automate handling of events, so imagine people writing switch statements to handle them.
func (e *Event) NormalEvent(obj runtime.Object, reason string, message string) {
	e.k8sEventRecorder.Event(obj, "Normal", reason, message)
}

// K8s Warning event
// 'reason' is the reason this event is generated. 'reason' should be short and unique; it
// should be in UpperCamelCase format (starting with a capital letter). "reason" will be used
// to automate handling of events, so imagine people writing switch statements to handle them.
func (e *Event) WarningEvent(obj runtime.Object, reason string, message string) {
	e.k8sEventRecorder.Event(obj, "Warning", reason, message)
}
