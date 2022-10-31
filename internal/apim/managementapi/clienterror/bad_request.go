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
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type BadRequestError struct {
	Message string `json:"message"`
}

func (e BadRequestError) Error() string {
	if e.Message == "" {
		return "BAD REQUEST"
	}
	return e.Message
}

func NewBadRequestError(resp *http.Response) BadRequestError {
	if resp.Body == nil {
		return BadRequestError{}
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return BadRequestError{}
	}

	var badRequestErr BadRequestError

	if err = json.Unmarshal(body, &badRequestErr); err != nil {
		return BadRequestError{}
	}

	return badRequestErr
}

func IsBadRequest(err error) bool {
	return errors.As(err, &BadRequestError{})
}
