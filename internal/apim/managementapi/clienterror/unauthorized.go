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

type UnauthorizedError struct {
	message string
}

func (e UnauthorizedError) Error() string {
	if e.message == "" {
		return "UNAUTHORIZED"
	}
	return e.message
}

func NewUnauthorizedCrossIdRequestError(crossId string) UnauthorizedError {
	return UnauthorizedError{message: "Unauthorized error for CrossId " + crossId}
}

func NewUnauthorizedApiRequestError(apiId string) UnauthorizedError {
	return UnauthorizedError{message: "Unauthorized error for API " + apiId}
}

func IsUnauthorized(err error) bool {
	return errors.As(err, &UnauthorizedError{})
}
