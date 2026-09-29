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

package am

import (
	"fmt"
	"net/http"
)

// BaseResponse is what an AM upsert reports back to the CR status: the AM key, organization and environment.
type BaseResponse struct {
	OrgEnv `json:",inline"`
	Key    string `json:"key"`
}

// UnexpectedResponse reports an AM response the SDK could not decode (no JSON body for the status).
func UnexpectedResponse(resp *http.Response) error {
	return fmt.Errorf("unexpected AM response: status %d, content type %q",
		resp.StatusCode, resp.Header.Get("Content-Type"))
}
