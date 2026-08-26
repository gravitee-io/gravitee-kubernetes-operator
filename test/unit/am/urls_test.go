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

package am_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	amclient "github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/client"
)

var _ = Describe("AM automation URLs", func() {
	It("builds the default automation root", func() {
		urls, err := amclient.NewURLs("http://am.example", "/automation", "DEFAULT", "DEFAULT")
		Expect(err).ToNot(HaveOccurred())
		Expect(urls.Automation.String()).To(Equal(
			"http://am.example/automation/organizations/DEFAULT/environments/DEFAULT",
		))
	})

	It("honours a spec.path override of the /automation prefix", func() {
		urls, err := amclient.NewURLs("http://am.example", "/custom", "my-org", "my-env")
		Expect(err).ToNot(HaveOccurred())
		Expect(urls.Automation.String()).To(Equal(
			"http://am.example/custom/organizations/my-org/environments/my-env",
		))
	})

	It("appends the domains probe path and size query", func() {
		urls, err := amclient.NewURLs("http://am.example", "/automation", "DEFAULT", "DEFAULT")
		Expect(err).ToNot(HaveOccurred())
		c := &amclient.Client{URLs: urls}
		Expect(c.AutomationTarget("domains").WithQueryParam("size", "1").String()).To(Equal(
			"http://am.example/automation/organizations/DEFAULT/environments/DEFAULT/domains?size=1",
		))
	})
})
