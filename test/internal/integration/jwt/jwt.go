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

package jwt

import (
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gravitee-io/gravitee-kubernetes-operator/examples"
)

type Claims struct {
	ClientID interface{} `json:"client_id,omitempty"`
	jwt.RegisteredClaims
}

func CreateClaims(clientID string) *Claims {
	return &Claims{
		ClientID: clientID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "gravitee.io/kubernetes-operator",
			Subject:   clientID,
			Audience:  []string{clientID},
		},
	}
}

func GetToken(clientID string, pkPath string) (string, error) {
	pemContent, err := examples.FS.ReadFile(pkPath)
	if err != nil {
		return "", err
	}
	alg := jwt.SigningMethodRS256
	pem, _ := pem.Decode(pemContent)
	key, err := x509.ParsePKCS1PrivateKey(pem.Bytes)
	if err != nil {
		return "", err
	}
	return jwt.NewWithClaims(alg, CreateClaims(clientID)).SignedString(key)
}
