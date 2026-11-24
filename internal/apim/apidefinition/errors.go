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

package apidefinition

import (
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// refNotFoundError words a reference the cluster does not hold as the admission webhook always has, and
// unwraps to the Kubernetes NotFound, so apierrors.IsNotFound still holds (drift on an old version, delete).
type refNotFoundError struct {
	msg string
	err error
}

func (e refNotFoundError) Error() string { return e.msg }

func (e refNotFoundError) Unwrap() error { return e.err }

// notFoundRef words err with format when it is a NotFound, and returns any other error as is.
func notFoundRef(err error, format string, args ...any) error {
	if !apierrors.IsNotFound(err) {
		return err
	}
	return refNotFoundError{msg: fmt.Sprintf(format, args...), err: err}
}
