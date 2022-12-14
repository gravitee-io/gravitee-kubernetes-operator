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

package clienterror

import "errors"

type NotFoundError struct {
	message string
}

func (e NotFoundError) Error() string {
	if e.message == "" {
		return "NOT FOUND"
	}
	return e.message
}

func NewCrossIdNotFoundError(crossId string) NotFoundError {
	return NotFoundError{message: "No API found for CrossId " + crossId}
}

func NewApiNotFoundError(apiId string) NotFoundError {
	return NotFoundError{message: "No API found for ID " + apiId}
}

func IsNotFound(err error) bool {
	return errors.As(err, &NotFoundError{})
}
