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
	"fmt"
	"testing"
)

func Test_IsUnauthorizedError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			"Should not be unauthorized with raw error",
			fmt.Errorf("raw error"),
			false,
		},
		{
			"Should not be unauthorized with nil error",
			nil,
			false,
		},
		{
			"Should be unauthorized with unauthorized API error",
			NewUnauthorizedApiRequestError("api-id"),
			true,
		},
		{
			"Should be unauthorized with unauthorized crossId error",
			NewUnauthorizedCrossIdRequestError("cross-id"),
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			given := IsUnauthorized(test.err)
			if given != test.expected {
				t.Fail()
				t.Logf("Expected %t to be %t", given, test.expected)
			}
		})
	}
}

func Test_IsNotFoundError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			"Should not be not found with raw error",
			fmt.Errorf("raw error"),
			false,
		},
		{
			"Should not be not found with nil error",
			nil,
			false,
		},
		{
			"Should be not found with API error",
			NewApiNotFoundError("api-id"),
			true,
		},
		{
			"Should be not found with crossId error",
			NewCrossIdNotFoundError("cross-id"),
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			given := IsNotFound(test.err)
			if given != test.expected {
				t.Fail()
				t.Logf("Expected %t to be %t", given, test.expected)
			}
		})
	}
}

func Test_IsIllegalStateError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			"Should not be illegal state error with raw error",
			fmt.Errorf("raw error"),
			false,
		},
		{
			"Should not be illegal state error not found with nil error",
			nil,
			false,
		},
		{
			"Should be illegal state error with ambiguous crossId error",
			NewAmbiguousCrossIdError("cross-id", 2),
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			given := IsIllegalState(test.err)
			if given != test.expected {
				t.Fail()
				t.Logf("Expected %t to be %t", given, test.expected)
			}
		})
	}
}
