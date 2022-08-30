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
