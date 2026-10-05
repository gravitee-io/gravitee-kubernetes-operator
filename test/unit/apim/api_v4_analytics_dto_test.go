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


package apim_test

import (
	"encoding/json"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("API v4 analytics payload", func() {
	sentAnalytics := func(analytics *v4.Analytics) map[string]json.RawMessage {
		api := &v4.Api{V4BaseApi: &v4.V4BaseApi{ApiBase: &base.ApiBase{}, Analytics: analytics}}
		payload, err := json.Marshal(model.ToAPIV4DTO(api))
		Expect(err).ToNot(HaveOccurred())
		var sent struct {
			Analytics map[string]json.RawMessage `json:"analytics"`
		}
		Expect(json.Unmarshal(payload, &sent)).To(Succeed())
		return sent.Analytics
	}

	It("sends the selected connection events to APIM", func() {
		events := []v4.ConnectionEvent{v4.ConnectionEventConnected, v4.ConnectionEventDisconnected}

		analytics := sentAnalytics(&v4.Analytics{Enabled: true, ConnectionEvents: &events})

		Expect(analytics).To(HaveKey("connectionEvents"))
		Expect(analytics["connectionEvents"]).To(MatchJSON(`["CONNECTED", "DISCONNECTED"]`))
	})

	It("sends an empty selection, which reports no connection event", func() {
		analytics := sentAnalytics(&v4.Analytics{Enabled: true, ConnectionEvents: &[]v4.ConnectionEvent{}})

		Expect(analytics).To(HaveKey("connectionEvents"))
		Expect(analytics["connectionEvents"]).To(MatchJSON(`[]`))
	})

	It("sends an unset selection as null, which keeps the gateway default", func() {
		analytics := sentAnalytics(&v4.Analytics{Enabled: true})

		Expect(analytics).To(HaveKeyWithValue("connectionEvents", json.RawMessage("null")))
	})
})
