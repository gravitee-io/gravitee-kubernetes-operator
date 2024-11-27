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
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/examples"
)

func NewClient() *http.Client {
	tr := &http.Transport{
		// #nosec G402
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
	}
}

func NewMTLSClient(caFile, crtFile, keyFile string) *http.Client {
	ca, err := examples.FS.ReadFile(caFile)
	if err != nil {
		panic("unable to load sever CA for client authentication")
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(ca)

	crt, err := examples.FS.ReadFile(crtFile)
	if err != nil {
		panic("unable to load crt for client authentication")
	}
	key, err := examples.FS.ReadFile(keyFile)
	if err != nil {
		panic("unable to load key for client authentication")
	}
	keyPair, err := tls.X509KeyPair(crt, key)
	if err != nil {
		panic("unable to parse key pair for client authentication")
	}
	tr := &http.Transport{
		// #nosec G402
		TLSClientConfig: &tls.Config{
			RootCAs:      caCertPool,
			Certificates: []tls.Certificate{keyPair},
		},
	}
	return &http.Client{
		Transport: tr,
		Timeout:   3000 * time.Second,
	}
}
