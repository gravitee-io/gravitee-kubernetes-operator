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

package http

import "net/url"

// URL is a wrapper around url.URL that provides a fluent API for building URLs.
type URL struct {
	base *url.URL
}

// WithPath appends the given segments to the URL path.
func (u *URL) WithPath(segments ...string) *URL {
	return &URL{u.base.JoinPath(segments...)}
}

// WithQueryParam adds the given key-value pair to the URL query.
func (u *URL) WithQueryParam(k, v string) *URL {
	base := u.base
	query := base.Query()
	query.Add(k, v)
	base.RawQuery = query.Encode()
	return &URL{base}
}

// WithQueryParams adds the given key-value pairs to the URL query.
func (u *URL) WithQueryParams(params map[string]string) *URL {
	base := u.base
	query := base.Query()
	for k, v := range params {
		query.Add(k, v)
	}
	base.RawQuery = query.Encode()
	return &URL{base}
}

// String returns the URL as a string.
func (u *URL) String() string {
	return u.base.String()
}

func NewURL(baseUrl string) (*URL, error) {
	base, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	return &URL{base}, nil
}
