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
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// FromSDKResponse returns err, or a server error for a 4xx/5xx response of a generated SDK call.
// response is only called when err is nil: on a transport error the SDK returns no response.
// The SDK has already read the response body, so it is passed in and put back for the error to
// report the server's message.
func FromSDKResponse(err error, response func() (*http.Response, []byte)) error {
	if err != nil {
		return err
	}
	resp, body := response()
	resp.Body = io.NopCloser(bytes.NewReader(body))
	return FromResponse(resp)
}

// UnexpectedSDKResponse reports a successful response the SDK could not decode (no JSON body).
func UnexpectedSDKResponse(resp *http.Response) error {
	return fmt.Errorf("unexpected response: status %d, content type %q",
		resp.StatusCode, resp.Header.Get("Content-Type"))
}
