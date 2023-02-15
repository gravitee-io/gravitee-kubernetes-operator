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

package internal

import (
	"errors"

	apimError "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	kErrors "k8s.io/apimachinery/pkg/util/errors"
)

type ContextError struct {
	error
}

// Redirects the behavior of Is to As
// because Is is not implemented for k8s.io errors aggregate.
func (e ContextError) Is(err error) bool {
	return errors.As(err, new(ContextError))
}

func newContextError(err error) error {
	return ContextError{err}
}

func IsRecoverable(err error) bool {
	errs := make([]error, 0)

	//nolint:errorlint // type assertion is intended here (Aggregate is an interface)
	if agg, ok := err.(kErrors.Aggregate); ok {
		errs = kErrors.Flatten(agg).Errors()
	} else {
		errs = append(errs, err)
	}

	for _, e := range errs {
		if isRecoverable(e) {
			return true
		}
	}

	return false
}

func isRecoverable(err error) bool {
	contextError := &ContextError{}
	if errors.As(err, contextError) {
		cause := contextError.error
		serverError := &apimError.ServerError{}
		if errors.As(cause, serverError) {
			return serverError.IsRecoverable()
		}
	}

	return true
}
