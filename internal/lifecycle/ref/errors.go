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

package ref

import (
	"errors"
	"fmt"
)

var (
	ErrRefTagInvalid      = errors.New("ref tag must be kind[,refField[,key]]")
	ErrUnknownKind        = errors.New("unknown kind")
	ErrStructKind         = errors.New("not a struct")
	ErrRefTypeUnsupported = errors.New("ref field is not refs.NamespacedName")
	ErrNotASecret         = errors.New("object is not a secret")
	ErrCannotSetField     = errors.New("cannot set field")
	ErrTypeMismatch       = errors.New("type mismatch")
	ErrMissingRef         = errors.New("sibling ref is missing")
	ErrSecretKeyMissing   = errors.New("secret data key is missing")
)

// Error is a ref-resolution failure. Unwrap is the cause (sentinel or API error).
// errors.Is(err, ErrUnknownKind) follows Unwrap.
// errors.Is(err, &Error{Op: "resolve"}) matches on any set field of the target.
type Error struct {
	Op    string
	Kind  string
	Field string
	Err   error
}

func (e *Error) Error() string {
	if e == nil {
		return "lifecycle/ref: <nil>"
	}
	msg := "lifecycle/ref"
	if e.Op != "" {
		msg += ": " + e.Op
	}
	if e.Kind != "" {
		msg += fmt.Sprintf(": kind %s", e.Kind)
	}
	if e.Field != "" {
		msg += fmt.Sprintf(": field %s", e.Field)
	}
	if e.Err != nil {
		msg += ": " + e.Err.Error()
	}
	return msg
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func (e *Error) Is(target error) bool {
	var t *Error
	ok := errors.As(target, &t)
	if !ok || e == nil || t == nil {
		return false
	}
	if t.Err != nil && e.Err != t.Err {
		return false
	}
	if t.Op != "" && e.Op != t.Op {
		return false
	}
	if t.Kind != "" && e.Kind != t.Kind {
		return false
	}
	if t.Field != "" && e.Field != t.Field {
		return false
	}
	return true
}

func NewWrappedError(op, kind, field string, err error) error {
	if err == nil {
		return nil
	}
	var existing *Error
	if errors.As(err, &existing) {
		return err
	}
	return &Error{Op: op, Kind: kind, Field: field, Err: err}
}
