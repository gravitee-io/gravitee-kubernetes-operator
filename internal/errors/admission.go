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

package errors

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type Severity string

const (
	Severe  = Severity("severe")
	Warning = Severity("warning")
)

type AdmissionError struct {
	Severity Severity
	Message  string
}

type AdmissionErrors struct {
	Warning []*AdmissionError
	Severe  []*AdmissionError
}

func NewAdmissionErrors() *AdmissionErrors {
	return &AdmissionErrors{
		Severe:  make([]*AdmissionError, 0),
		Warning: make([]*AdmissionError, 0),
	}
}

func (errs *AdmissionErrors) Map() (admission.Warnings, error) {
	warnings := admission.Warnings{}
	for _, w := range errs.Warning {
		warnings = append(warnings, w.Error())
	}
	if len(errs.Severe) == 0 {
		return warnings, nil
	}
	return warnings, errs.Severe[0]
}

func (errs *AdmissionErrors) MergeWith(other *AdmissionErrors) {
	if other == nil {
		return
	}
	errs.Warning = append(errs.Warning, other.Warning...)
	errs.Severe = append(errs.Severe, other.Severe...)
}

func (errs *AdmissionErrors) IsSevere() bool {
	return len(errs.Severe) > 0
}

func (errs *AdmissionErrors) AddSevere(format string, args ...any) {
	errs.Severe = append(errs.Severe, NewSevere(format, args...))
}

func (errs *AdmissionErrors) AddWarning(format string, args ...any) {
	errs.Warning = append(errs.Warning, NewWarning(format, args...))
}

func (errs *AdmissionErrors) Add(err *AdmissionError) {
	if err == nil {
		return
	}
	if err.Severity == Severe {
		errs.Severe = append(errs.Severe, err)
	}
	if err.Severity == Warning {
		errs.Warning = append(errs.Warning, err)
	}
}

func (err *AdmissionError) Error() string {
	return err.Message
}

func NewSevere(format string, args ...any) *AdmissionError {
	return &AdmissionError{
		Severity: Severe,
		Message:  fmt.Sprintf(format, args...),
	}
}

func NewWarning(format string, args ...any) *AdmissionError {
	return &AdmissionError{
		Severity: Warning,
		Message:  fmt.Sprintf(format, args...),
	}
}
