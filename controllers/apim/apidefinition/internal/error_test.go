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
	"fmt"
	"testing"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi/clienterror"
)

func Test_wrapError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			"Should be recoverable with raw error",
			fmt.Errorf("raw error"),
			true,
		},
		{
			"Should be recoverable with not found error",
			clienterror.NewCrossIdNotFoundError("cross-id"),
			true,
		},
		{
			"Should be recoverable with illegal state error",
			clienterror.NewAmbiguousCrossIdError("cross-id", 2),
			true,
		},
		{
			"Should not be recoverable with unauthorized api request error",
			clienterror.NewUnauthorizedApiRequestError("api-id"),
			false,
		},
		{
			"Should not be recoverable with unauthorized cross ID request error",
			clienterror.NewUnauthorizedCrossIdRequestError("cross-id"),
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			given := IsRecoverableError(wrapError(test.err))
			if given != test.expected {
				t.Fail()
				t.Logf("Expected %t to be %t", given, test.expected)
			}
		})
	}
}
