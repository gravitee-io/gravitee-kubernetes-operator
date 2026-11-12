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

package apidefinition_test

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apidefinition"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const httpListenerTemplate = `{
  "type": "HTTP",
  "paths": [{"path": %q, "overrideAccess": true}, {"host": "api.example.com", "path": "/already/"}],
  "pathMappings": ["/orders/:id"],
  "entrypoints": [{"type": "http-proxy", "qos": "AUTO"}],
  "cors": {
    "enabled": true,
    "allowOrigin": ["https://.*\\.example\\.com"],
    "allowMethods": ["GET", "OPTIONS"],
    "allowHeaders": ["*"],
    "allowCredentials": false,
    "maxAge": 1728000,
    "runPolicies": false
  },
  "requestValidation": {"rejectNullByte": true}
}`

const subscriptionListener = `{
  "type": "SUBSCRIPTION",
  "entrypoints": [{"type": "webhook", "configuration": {"proxy": {"enabled": false}}}]
}`

var _ = Describe("PrepareV4SpecForAutomation listener path normalisation", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
		scheme := runtime.NewScheme()
		Expect(clientgoscheme.AddToScheme(scheme)).To(Succeed())
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
		k8s.RegisterClient(fake.NewClientBuilder().WithScheme(scheme).Build())
	})

	httpListener := func(path string) string {
		return fmt.Sprintf(httpListenerTemplate, path)
	}

	toListener := func(raw string) *v4.GenericListener {
		listener := new(v4.GenericListener)
		Expect(json.Unmarshal([]byte(raw), listener)).To(Succeed())
		return listener
	}

	newAPI := func(listeners ...*v4.GenericListener) *v1alpha1.ApiV4Definition {
		return &v1alpha1.ApiV4Definition{
			ObjectMeta: metav1.ObjectMeta{Name: "api", Namespace: "default"},
			Spec: v1alpha1.ApiV4DefinitionSpec{
				Api: v4.Api{
					V4BaseApi: &v4.V4BaseApi{
						ApiBase:   &base.ApiBase{Name: "api", Version: "1.0"},
						Listeners: listeners,
					},
				},
			},
		}
	}

	DescribeTable("sends the listener with only its path strings changed",
		func(input, expected string) {
			api := newAPI(toListener(input))

			Expect(apidefinition.PrepareV4SpecForAutomation(ctx, api, false)).To(Succeed())

			dto := model.ToAPIV4DTO(&api.Spec.Api)
			Expect(dto.Listeners).To(HaveLen(1))
			Expect(json.Marshal(dto.Listeners[0])).To(MatchJSON(expected))
		},
		Entry("HTTP listener keeps cors, requestValidation and overrideAccess",
			httpListener("/contracts"), httpListener("/contracts/")),
		Entry("HTTP listener whose paths already end with a slash",
			httpListener("/contracts/"), httpListener("/contracts/")),
		Entry("non-HTTP listener is left as is",
			subscriptionListener, subscriptionListener),
		Entry("HTTP listener with an empty path keeps it empty",
			httpListener(""), httpListener("")),
	)

	It("passes a listener with no content through and normalises the others", func() {
		empty := &v4.GenericListener{}
		api := newAPI(empty, toListener(httpListener("/contracts")))

		Expect(apidefinition.PrepareV4SpecForAutomation(ctx, api, false)).To(Succeed())

		Expect(api.Spec.Api.Listeners).To(HaveLen(2))
		Expect(api.Spec.Api.Listeners[0]).To(BeIdenticalTo(empty))
		Expect(json.Marshal(api.Spec.Api.Listeners[1])).To(MatchJSON(httpListener("/contracts/")))
	})

	It("skips a nil listener and normalises the others", func() {
		api := newAPI(nil, toListener(httpListener("/contracts")))

		Expect(apidefinition.PrepareV4SpecForAutomation(ctx, api, false)).To(Succeed())

		Expect(api.Spec.Api.Listeners).To(HaveLen(2))
		Expect(api.Spec.Api.Listeners[0]).To(BeNil())
		Expect(json.Marshal(api.Spec.Api.Listeners[1])).To(MatchJSON(httpListener("/contracts/")))
	})

	It("does not modify the listener object it was given", func() {
		original := toListener(httpListener("/contracts"))
		api := newAPI(original)

		Expect(apidefinition.PrepareV4SpecForAutomation(ctx, api, false)).To(Succeed())

		Expect(json.Marshal(original)).To(MatchJSON(httpListener("/contracts")))
	})
})
