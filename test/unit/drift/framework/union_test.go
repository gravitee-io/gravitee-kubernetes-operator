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

package framework

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// The shape of a generated SDK union (OneOf*): one pointer per variant, each tagged with its type.
type unionNamed struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type unionEndpoint struct {
	Type     string  `json:"type"`
	Endpoint *string `json:"endpoint,omitempty"`
}

type union struct {
	Named    *unionNamed    `oneof:"NAMED" json:"named,omitempty"`
	Endpoint *unionEndpoint `oneof:"ENDPOINT" json:"endpoint,omitempty"`
}

type withUnion struct {
	Auth  *union  `json:"auth,omitempty"`
	Items []union `json:"items,omitempty"`
}

func named(name string) union {
	return union{Named: &unionNamed{Type: "NAMED", Name: name}}
}

func endpoint(value *string) union {
	return union{Endpoint: &unionEndpoint{Type: "ENDPOINT", Endpoint: value}}
}

var _ = Describe("SDK unions", func() {

	It("reports a variant field at the path the wire carries", func() {
		crd := withUnion{Auth: new(named("foo"))}
		remote := withUnion{Auth: new(named("bar"))}

		expectDrift(drift.DetectWithNamespace(crd, remote, ""), `auth:
  name: "foo" != "bar"`)
	})

	It("reports a variant change", func() {
		crd := withUnion{Auth: new(named("foo"))}
		remote := withUnion{Auth: new(endpoint(nil))}

		result := drift.DetectWithNamespace(crd, remote, "")
		Expect(result.DriftDetected()).To(BeTrue())
	})

	Describe("in a list, at admission", func() {
		admission := func(crd, remote any) drift.Result {
			return drift.Merge(drift.DetectWithNamespace(crd, remote, ""), drift.DetectWithNamespace(crd, remote, ""))
		}

		It("reports a variant field the resource leaves out and the remote sets", func() {
			crd := withUnion{Items: []union{named("foo"), endpoint(nil)}}
			remote := withUnion{Items: []union{named("foo"), endpoint(new("/me"))}}

			expectDrift(admission(crd, remote), `items[1]:
  endpoint: <nil> != "/me"`)
		})

		It("reports an item only the remote has", func() {
			crd := withUnion{Items: []union{named("foo")}}
			remote := withUnion{Items: []union{named("foo"), endpoint(nil)}}

			result := admission(crd, remote)
			Expect(result.DriftDetected()).To(BeTrue())
		})

		It("admits a variant field the resource leaves out outside a list", func() {
			crd := withUnion{Auth: new(endpoint(nil))}
			remote := withUnion{Auth: new(endpoint(new("/me")))}

			expectNoDrift(admission(crd, remote))
		})
	})
})
