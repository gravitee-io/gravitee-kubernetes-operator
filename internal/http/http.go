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

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

const (
	requestTimeoutSeconds = 5
)

type Client struct {
	ctx  context.Context
	http http.Client
}

func (client *Client) Context() context.Context {
	return client.ctx
}

// RequestTransformer is a function that can be used to mutate a request
// before it is sent (e.g. setting headers).
type RequestTransformer = func(*http.Request)

func WithHost(host string) RequestTransformer {
	return func(req *http.Request) {
		req.Host = host
	}
}

// Get returns the result of a GET request to the specified URL, marshaled into the target.
// If the target is nil, the response is discarded.
func (client *Client) Get(url URL, target any, transformers ...RequestTransformer) error {
	req, err := client.prepareGet(url.String())
	if err != nil {
		return errors.FromNewRequestError(http.MethodGet, url.String(), err)
	}

	for _, transform := range transformers {
		transform(req)
	}

	return client.do(req, target)
}

// Special GET method for YAML content type (used in tests).
func (client *Client) GetYAML(url URL, target any, transformers ...RequestTransformer) error {
	req, err := client.prepareGet(url.String())
	if err != nil {
		return errors.FromNewRequestError(http.MethodGet, url.String(), err)
	}

	for _, transform := range transformers {
		transform(req)
	}

	return client.doYAML(req, target)
}

// Post returns the result of a POST request to the specified URL,
// using entity as the body of the request, marshaling the result into target.
// If the target is nil, the response is discarded.
func (client *Client) Post(url URL, entity, target any, transformers ...RequestTransformer) error {
	req, err := client.preparePost(url.String(), entity)
	if err != nil {
		return errors.FromNewRequestError(http.MethodPost, url.String(), err)
	}

	for _, transform := range transformers {
		transform(req)
	}

	return client.do(req, target)
}

// Put returns the result of a PUT request to the specified URL,
// using entity as the body of the request, marshaling the result into target.
// If the target is nil, the response is discarded.
func (client *Client) Put(url URL, entity, target any, transformers ...RequestTransformer) error {
	req, err := client.preparePut(url.String(), entity)
	if err != nil {
		return errors.FromNewRequestError(http.MethodPut, url.String(), err)
	}

	for _, transform := range transformers {
		transform(req)
	}

	return client.do(req, target)
}

// Delete returns the result of a DELETE request to the specified URL, marshaling the result into target.
// If the target is nil, the response is discarded.
func (client *Client) Delete(url URL, target any, transformers ...RequestTransformer) error {
	req, err := client.prepareDelete(url.String())
	if err != nil {
		return errors.FromNewRequestError(http.MethodDelete, url.String(), err)
	}

	for _, transform := range transformers {
		transform(req)
	}

	return client.do(req, target)
}

func (client *Client) do(req *http.Request, target any) error {
	resp, err := client.http.Do(req)
	if err != nil {
		return errors.FromDoRequestError(req, err)
	}

	defer resp.Body.Close()

	if err = errors.FromResponse(resp); err != nil {
		return err
	}

	return WriteJSON(resp, target)
}

func (client *Client) doYAML(req *http.Request, target any) error {
	resp, err := client.http.Do(req)
	if err != nil {
		return errors.FromDoRequestError(req, err)
	}

	defer resp.Body.Close()

	if err = errors.FromResponse(resp); err != nil {
		return err
	}

	return WriteYAML(resp, target)
}

func (client *Client) preparePost(url string, entity any) (*http.Request, error) {
	if entity == nil {
		return http.NewRequestWithContext(client.ctx, http.MethodPost, url, nil)
	}
	return client.newJSONRequest(http.MethodPost, url, entity)
}

func (client *Client) prepareGet(url string) (*http.Request, error) {
	return http.NewRequestWithContext(client.ctx, http.MethodGet, url, nil)
}

func (client *Client) preparePut(url string, entity any) (*http.Request, error) {
	if entity == nil {
		return http.NewRequestWithContext(client.ctx, http.MethodPut, url, nil)
	}
	return client.newJSONRequest(http.MethodPut, url, entity)
}

func (client *Client) prepareDelete(url string) (*http.Request, error) {
	return http.NewRequestWithContext(client.ctx, http.MethodDelete, url, nil)
}

func (client *Client) newJSONRequest(method, url string, entity any) (*http.Request, error) {
	reader, err := ReadJSON(entity)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(client.ctx, method, url, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Add(ContentTypeHeader, ContentTypeJSON)
	return req, nil
}

func NewNoAuthClient(ctx context.Context) *Client {
	return NewClient(ctx, nil)
}

func NewClient(ctx context.Context, auth *Auth) *Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{
		// #nosec G402
		InsecureSkipVerify: env.Config.InsecureSkipVerify,
	}

	httpClient := http.Client{Timeout: requestTimeoutSeconds * time.Second, Transport: transport}

	if auth != nil {
		authRoundTripper := NewAuthenticatedRoundTripper(auth, transport)
		httpClient.Transport = authRoundTripper
	}

	return &Client{ctx, httpClient}
}
