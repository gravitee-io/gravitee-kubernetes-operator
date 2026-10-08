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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/management"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
)

var _ = Describe("EnvTarget", func() {
	path := "/proxied"

	DescribeTable("returns the Automation API environment URL",
		func(spec management.Context, expected string) {
			mctx := &v1alpha1.ManagementContext{Spec: v1alpha1.ManagementContextSpec{Context: &spec}}
			target, err := apim.EnvTarget(mctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(target).To(Equal(expected))
		},
		Entry("for a self-hosted context",
			management.Context{BaseUrl: "http://localhost:30083", OrgID: "DEFAULT", EnvID: "DEFAULT"},
			"http://localhost:30083/automation/organizations/DEFAULT/environments/DEFAULT"),
		Entry("for a context with a path override",
			management.Context{BaseUrl: "http://localhost:30083", OrgID: "DEFAULT", EnvID: "DEFAULT", Path: &path},
			"http://localhost:30083/proxied/organizations/DEFAULT/environments/DEFAULT"),
		Entry("for a Gravitee Cloud context",
			management.Context{
				BaseUrl: "https://eu.cloudgate.gravitee.io", OrgID: "org", EnvID: "env",
				Cloud: &management.Cloud{Token: "token"},
			},
			"https://eu.cloudgate.gravitee.io/apim/automation/organizations/org/environments/env"),
	)
})
