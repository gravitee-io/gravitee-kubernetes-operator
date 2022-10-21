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

import (
	"errors"
	"strconv"
)

type IllegalStateError struct {
	message string
}

func (e IllegalStateError) Error() string {
	if e.message == "" {
		return "ILLEGAL STATE"
	}
	return e.message
}

func NewAmbiguousCrossIdError(crossId string, count int) IllegalStateError {
	return IllegalStateError{message: "Expected one API with CrossId " + crossId + "but found " + strconv.Itoa(count)}
}

func IsIllegalState(err error) bool {
	return errors.As(err, &IllegalStateError{})
}
