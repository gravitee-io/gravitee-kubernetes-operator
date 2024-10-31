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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	kErrors "k8s.io/apimachinery/pkg/util/errors"
)

type ServerError struct {
	StatusCode int    `json:"status"`
	URL        string `json:"url"`
	Method     string `json:"method"`
	Message    string `json:"message"`
	Body       string `json:"body"`
}

type ContextError struct {
	error
}

type UnrecoverableError struct {
	message string
}

func (err UnrecoverableError) Error() string {
	return err.message
}

func NewUnrecoverableError(message string) error {
	return UnrecoverableError{message}
}

var unRecoverableStatusCodes = []int{http.StatusBadRequest, http.StatusUnauthorized}

func (err ServerError) Error() string {
	message := err.Message
	if message == "" {
		message = http.StatusText(err.StatusCode)
	}
	return fmt.Sprintf("request [%s] %s failed with status %d (%s)", err.Method, err.URL, err.StatusCode, message)
}

func NewServerError(resp *http.Response) ServerError {
	statusCode := resp.StatusCode
	url := resp.Request.URL
	method := resp.Request.Method

	serverError := ServerError{StatusCode: statusCode, URL: url.String(), Method: method}

	if resp.Body == nil {
		return serverError
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return serverError
	}

	serverError.Body = string(body)

	if err = json.Unmarshal(body, &serverError); err != nil {
		return serverError
	}

	return serverError
}

func IsServerError(err error) bool {
	return errors.As(err, &ServerError{})
}

func (err ServerError) IsRecoverable() bool {
	for _, code := range unRecoverableStatusCodes {
		if err.StatusCode == code {
			return false
		}
	}
	return true
}

func FromResponse(resp *http.Response) error {
	if resp.StatusCode < http.StatusBadRequest {
		return nil
	}
	return NewServerError(resp)
}

func FromDoRequestError(req *http.Request, err error) error {
	return fmt.Errorf(
		"unable to perform request [%s] %s: (%w)",
		req.Method, req.URL, err,
	)
}

func FromNewRequestError(method, url string, err error) error {
	return fmt.Errorf(
		"unable to create request [%s] %s: (%w)",
		method, url, err,
	)
}

func NewNotFoundError() error {
	return ServerError{StatusCode: http.StatusNotFound}
}

func IsNotFound(err error) bool {
	serverError := &ServerError{}
	if errors.As(err, serverError) {
		return serverError.StatusCode == http.StatusNotFound
	}
	return false
}

func IsBadRequest(err error) bool {
	serverError := &ServerError{}
	if errors.As(err, serverError) {
		return serverError.StatusCode == http.StatusBadRequest
	}
	return false
}

func IsNetworkError(err error) bool {
	opErr := new(net.OpError)
	return errors.As(err, &opErr)
}

func IsUnauthorized(err error) bool {
	serverError := &ServerError{}
	if errors.As(err, serverError) {
		return serverError.StatusCode == http.StatusUnauthorized
	}
	return false
}

func IgnoreNotFound(err error) error {
	if IsNotFound(err) {
		return nil
	}

	return err
}

// Redirects the behavior of Is to As
// because Is is not implemented for k8s.io errors aggregate.
func (e ContextError) Is(err error) bool {
	return errors.As(err, new(ContextError))
}

func NewContextError(err error) error {
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
	if errors.As(err, &UnrecoverableError{}) {
		return false
	}
	contextError := &ContextError{}
	if errors.As(err, contextError) {
		cause := contextError.error
		serverError := &ServerError{}
		if errors.As(cause, serverError) {
			return serverError.IsRecoverable()
		}
	}

	return true
}
