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
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
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
func (client *Client) Get(url string, target any, transformers ...RequestTransformer) error {
	req, err := client.prepareGet(url)
	if err != nil {
		return errors.FromNewRequestError(http.MethodGet, url, err)
	}

	for _, transform := range transformers {
		transform(req)
	}

	return client.do(req, target)
}

// Special GET method for YAML content type (used in tests).
func (client *Client) GetYAML(url string, target any, transformers ...RequestTransformer) error {
	req, err := client.prepareGet(url)
	if err != nil {
		return errors.FromNewRequestError(http.MethodGet, url, err)
	}

	for _, transform := range transformers {
		transform(req)
	}

	return client.doYAML(req, target)
}

// Post returns the result of a POST request to the specified URL,
// using entity as the body of the request, marshaling the result into target.
// If the target is nil, the response is discarded.
func (client *Client) Post(url string, entity, target any, transformers ...RequestTransformer) error {
	req, err := client.preparePost(url, entity)
	if err != nil {
		return errors.FromNewRequestError(http.MethodPost, url, err)
	}

	for _, transform := range transformers {
		transform(req)
	}

	return client.do(req, target)
}

// Put returns the result of a PUT request to the specified URL,
// using entity as the body of the request, marshaling the result into target.
// If the target is nil, the response is discarded.
func (client *Client) Put(url string, entity, target any, transformers ...RequestTransformer) error {
	req, err := client.preparePut(url, entity)
	if err != nil {
		return errors.FromNewRequestError(http.MethodPut, url, err)
	}

	for _, transform := range transformers {
		transform(req)
	}

	return client.do(req, target)
}

// Delete returns the result of a DELETE request to the specified URL, marshaling the result into target.
// If the target is nil, the response is discarded.
func (client *Client) Delete(url string, target any, transformers ...RequestTransformer) error {
	req, err := client.prepareDelete(url)
	if err != nil {
		return errors.FromNewRequestError(http.MethodDelete, url, err)
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

func NewNoAuthClient(ctx context.Context) (*Client, error) {
	return NewClient(ctx, nil)
}

func NewClient(ctx context.Context, auth *Auth) (*Client, error) {
	httpClient, err := NewStdClient()
	if err != nil {
		return nil, err
	}

	if auth != nil {
		httpClient.Transport = NewAuthenticatedRoundTripper(auth, httpClient.Transport)
	}

	return &Client{ctx, *httpClient}, nil
}

// NewStdClient returns a net/http client configured from the operator settings: HTTP proxy,
// TLS verification, trust store and timeout. Use it for any client the operator does not build
// itself, such as the AM SDK. Every client shares one transport, so connections to APIM and AM
// are reused across reconciles instead of being opened again for each client.
func NewStdClient() (*http.Client, error) {
	transport, err := getTransport()
	if err != nil {
		return nil, err
	}
	timeout := time.Duration(env.Config.HTTPClientTimeoutSeconds) * time.Second
	return &http.Client{Timeout: timeout, Transport: transport}, nil
}

// transportSettings are the operator settings a transport is built from, plus the trust store
// file's size and modification time, so that a rotated trust store is loaded again.
type transportSettings struct {
	proxy              env.HttpProxy
	insecureSkipVerify bool
	trustStorePath     string
	trustStoreSize     int64
	trustStoreModTime  int64
}

var sharedTransport struct {
	sync.Mutex
	settings  transportSettings
	transport *http.Transport
}

// getTransport returns the transport built from the current settings. It is built on first use,
// and again when a setting or the trust store file changes; the previous one then closes its idle
// connections, and requests already running on it finish normally.
func getTransport() (*http.Transport, error) {
	settings, err := currentTransportSettings()
	if err != nil {
		return nil, err
	}

	sharedTransport.Lock()
	defer sharedTransport.Unlock()
	if sharedTransport.transport != nil && sharedTransport.settings == settings {
		return sharedTransport.transport, nil
	}

	transport, err := newTransport(settings)
	if err != nil {
		return nil, err
	}
	if sharedTransport.transport != nil {
		sharedTransport.transport.CloseIdleConnections()
	}
	sharedTransport.settings, sharedTransport.transport = settings, transport
	return transport, nil
}

func currentTransportSettings() (transportSettings, error) {
	settings := transportSettings{
		proxy:              env.Config.HttpProxy,
		insecureSkipVerify: env.Config.HTTPClientInsecureSkipVerify,
		trustStorePath:     env.Config.HttpClientTrustStorePath,
	}
	if settings.trustStorePath == "" {
		return settings, nil
	}
	info, err := os.Stat(settings.trustStorePath)
	if err != nil {
		return transportSettings{}, fmt.Errorf("failed to read trust store file %s: %w", settings.trustStorePath, err)
	}
	settings.trustStoreSize, settings.trustStoreModTime = info.Size(), info.ModTime().UnixNano()
	return settings, nil
}

func newTransport(settings transportSettings) (*http.Transport, error) {
	defaultTransport, _ := http.DefaultTransport.(*http.Transport)
	transport := defaultTransport.Clone()
	if settings.proxy.Enabled {
		if settings.proxy.UseSystemProxy {
			transport.Proxy = http.ProxyFromEnvironment
		} else {
			index := strings.Index(settings.proxy.URL, "://")
			rawUrl := settings.proxy.URL
			username := settings.proxy.Username
			password := settings.proxy.Password
			if username != "" {
				rawUrl = rawUrl[:index+3] + username + ":" + password + "@" + rawUrl[index+3:]
			}
			proxyURL, err := url.Parse(rawUrl)
			if err != nil {
				return nil, err
			}
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	transport.TLSClientConfig = &tls.Config{
		// #nosec G402
		InsecureSkipVerify: settings.insecureSkipVerify,
	}

	// Load your trust store (CA certificates)
	if settings.trustStorePath != "" {
		rootCAs, _ := x509.SystemCertPool()
		if rootCAs == nil {
			rootCAs = x509.NewCertPool()
		}

		caCert, err := os.ReadFile(settings.trustStorePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read trust store file %s: %w", settings.trustStorePath, err)
		}

		if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
			return nil, fmt.Errorf("failed to append CA certs from %s to pool", settings.trustStorePath)
		}

		transport.TLSClientConfig.RootCAs = rootCAs
	}

	return transport, nil
}
