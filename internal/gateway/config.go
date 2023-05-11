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

package gateway

import (
	"fmt"
	"strings"

	coreV1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	graviteeKubeScheme = "kubernetes://"
	jksKeystoreType    = "jks"

	expectedKubeFormat             = "$NS/(secrets|configmaps)/$NAME/$KEY"
	expectedKubePathComponentCount = 4

	SecretKubePropertyType = "secrets"
	ConfigKubePropertyType = "configmaps"
)

// Config is the configuration of the Gravitee Gateway.
// Currently, we only support the HTTP bit of this config,
// for keystore discovery when TLS is enabled on an ingress.
type Config struct {
	HTTP HTTPServerConfig `yaml:"http"`
}

// HTTPServerConfig if the HTTP server configuration of the Gravitee Gateway.
// see https://docs.gravitee.io/apim/3.x/apim_installguide_gateway_configuration.html#api-gateway-http-server
type HTTPServerConfig struct {
	TLS TLSConfig `yaml:"ssl"`
}

type TLSConfig struct {
	Keystore KeystoreConfig `yaml:"keystore,omitempty"`
}

type KeystoreConfig struct {
	Type     string               `yaml:"type"`
	Password string               `yaml:"password"`
	Location GraviteeKubeProperty `yaml:"kubernetes"`
}

func (gkc KeystoreConfig) Validate() error {
	if strings.ToLower(gkc.Type) != jksKeystoreType {
		return fmt.Errorf("expected keystore type jks, got %s", gkc.Type)
	}

	if !gkc.Location.IsValid() {
		return fmt.Errorf(
			"expected kubernetes location format /%s, got %s",
			expectedKubeFormat, gkc.Location,
		)
	}

	if gkc.Password == "" {
		return fmt.Errorf("password is required")
	}

	if strings.HasPrefix(gkc.Password, graviteeKubeScheme) {
		if !GraviteeKubeProperty(gkc.Password).IsValid() {
			return fmt.Errorf(
				"expected password location format %s%s, got %s",
				graviteeKubeScheme, expectedKubeFormat, gkc.Password,
			)
		}
	}

	return nil
}

type GraviteeKubeProperty string

func (gkp GraviteeKubeProperty) IsValid() bool {
	if len(strings.Split(gkp.TrimPrefix(), "/")) != expectedKubePathComponentCount {
		return false
	}
	return gkp.Type() == SecretKubePropertyType || gkp.Type() == ConfigKubePropertyType
}

func (gkp GraviteeKubeProperty) NewReceiver() client.Object {
	if gkp.Type() == SecretKubePropertyType {
		return &coreV1.Secret{}
	} else if gkp.Type() == ConfigKubePropertyType {
		return &coreV1.ConfigMap{}
	}
	return nil
}

func (gkp GraviteeKubeProperty) Get(receiver client.Object) []byte {
	switch r := receiver.(type) {
	case *coreV1.Secret:
		return r.Data[gkp.Key()]
	case *coreV1.ConfigMap:
		return []byte(r.Data[gkp.Key()])
	default:
		return nil
	}
}

func (gkp GraviteeKubeProperty) TrimPrefix() string {
	withoutScheme := strings.TrimPrefix(string(gkp), graviteeKubeScheme)
	return strings.TrimPrefix(withoutScheme, "/")
}

func (gkp GraviteeKubeProperty) Namespace() string {
	return strings.Split(gkp.TrimPrefix(), "/")[0]
}

func (gkp GraviteeKubeProperty) Type() string {
	return strings.Split(gkp.TrimPrefix(), "/")[1]
}

func (gkp GraviteeKubeProperty) Name() string {
	return strings.Split(gkp.TrimPrefix(), "/")[2]
}

func (gkp GraviteeKubeProperty) Key() string {
	return strings.Split(gkp.TrimPrefix(), "/")[3]
}

func (gkp GraviteeKubeProperty) String() string {
	return string(gkp)
}
