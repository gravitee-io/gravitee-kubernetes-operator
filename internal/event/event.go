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

package event

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

type Detail struct {
	BeforeReason        string
	BeforeMessage       string
	AfterSuccessReason  string
	AfterSuccessMessage string
	AfterFailureReason  string
}

type Action string

const (
	Delete Action = "DELETE"
	Update Action = "UPDATE"
)

type Kind string

const (
	Normal  Kind = "Normal"
	Warning Kind = "Warning"
)

var EventDetails = map[Action]Detail{
	Delete: {
		BeforeReason:        "DeleteStarted",
		BeforeMessage:       "Delete started",
		AfterSuccessReason:  "DeleteSucceeded",
		AfterSuccessMessage: "Delete succeeded",
		AfterFailureReason:  "DeleteFailed",
	},
	Update: {
		BeforeReason:        "UpdateStarted",
		BeforeMessage:       "Update started",
		AfterSuccessReason:  "UpdateSucceeded",
		AfterSuccessMessage: "Update succeeded",
		AfterFailureReason:  "UpdateFailed",
	},
}

type Recorder struct {
	k8sEventRecorder record.EventRecorder
}

func NewRecorder(recorder record.EventRecorder) *Recorder {
	return &Recorder{
		recorder,
	}
}

// Wrap a function call with event records
// if the function call returns an error, the event will be a Warning event
// otherwise, the event will be a Normal event.
func (e *Recorder) Record(action Action, obj runtime.Object, do func() error) error {
	details := EventDetails[action]

	e.info(obj, details.BeforeReason, details.BeforeMessage)

	err := do()

	if err != nil {
		e.warn(obj, details.AfterFailureReason, err.Error())
	} else {
		e.info(obj, details.AfterSuccessReason, details.AfterSuccessMessage)
	}

	return err
}

func (e *Recorder) info(obj runtime.Object, reason string, message string) {
	e.k8sEventRecorder.Event(obj, string(Normal), reason, message)
}

func (e *Recorder) warn(obj runtime.Object, reason string, message string) {
	e.k8sEventRecorder.Event(obj, string(Warning), reason, message)
}
