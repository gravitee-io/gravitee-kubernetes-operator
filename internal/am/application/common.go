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

package application

import (
	"strings"

	amsdk "github.com/gravitee-io/gravitee-automation-sdk/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/mapper"
)

// uuidPrefix marks a reference AM gives by id (id:<uuid>), not by key.
const uuidPrefix = "id:"

// Application is the AM payload of an application, with the key of the domain it lives under.
type Application struct {
	amsdk.Application
	DomainKey string
}

// Response is the upsert response of an application: the sub-resource keys and the clientId AM stored
// (spec.clientId, or the one it generated).
type Response struct {
	am.DomainSubResourceResponse
	ClientID string
}

// ToApplicationDTO converts an AMApplication object into an Application DTO, applying defaults and keys.
func ToApplicationDTO(obj *v1alpha1.AMApplication) (Application, error) {
	app := obj.Spec.Application
	dto, err := mapper.MapViaJSON[amsdk.Application](app)
	if err != nil {
		return Application{}, err
	}
	dto = dto.WithDefaults()
	dto.Key = refs.NewNamespacedNameFromObject(obj).HRID()

	// set keys from ref if necessary
	if app.Certificate == nil {
		dto.Certificate = refKey(obj, app.CertificateRef)
	}

	for i, idp := range obj.Spec.IdentityProviders {
		if idp.Identity == nil {
			dto.IdentityProviders[i].Identity = refKey(obj, idp.IdentityRef)
		}
	}

	return Application{
		Application: dto,
		DomainKey:   am.DomainKey(obj),
	}, nil
}

// CertificateKey returns the AM key of the application's certificate, from certificate or certificateRef.
// A reference by id (id:<uuid>) has no key: it gives "".
func CertificateKey(obj *v1alpha1.AMApplication) string {
	if ref := obj.Spec.CertificateRef; ref != nil {
		return refKey(obj, ref)
	}
	return amKey(obj.Spec.Certificate)
}

// IdentityKeys returns the AM keys of the application's identity providers, from identity or identityRef.
// References by id (id:<uuid>) are left out.
func IdentityKeys(obj *v1alpha1.AMApplication) []string {
	keys := make([]string, 0, len(obj.Spec.IdentityProviders))
	for _, idp := range obj.Spec.IdentityProviders {
		key := amKey(idp.Identity)
		if idp.IdentityRef != nil {
			key = refKey(obj, idp.IdentityRef)
		}
		if key != "" {
			keys = append(keys, key)
		}
	}
	return keys
}

// refKey is the AM key of a sibling CR: its HRID, always in the application's namespace
// (a cross-namespace ref is rejected by admission).
func refKey(obj *v1alpha1.AMApplication, ref *refs.NamespacedName) string {
	key := refs.NewNamespacedName(obj.GetNamespace(), ref.Name)
	return key.HRID()
}

func amKey(key *string) string {
	if key == nil || strings.HasPrefix(*key, uuidPrefix) {
		return ""
	}
	return *key
}
