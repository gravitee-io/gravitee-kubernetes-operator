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

package lifecycle_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
)

type groupAdmission = lifecycle.AdmissionLifecycle[*v1alpha1.Group, testDTO, *testClient]

func dryRunOK(context.Context, *testClient, testDTO) *gerrors.AdmissionErrors { return nil }

var _ = Describe("NewResourceLifecycle", func() {
	It("returns a lifecycle with all required holes", func() {
		l := newGroupLifecycle(&remote{})
		Expect(func() { lifecycle.NewResourceLifecycle(l) }).ToNot(Panic())
	})

	It("panics on the first missing required hole", func() {
		l := newGroupLifecycle(&remote{})
		l.Upsert = nil
		Expect(func() { lifecycle.NewResourceLifecycle(l) }).
			To(PanicWith("lifecycle: ResourceLifecycle.Upsert is required"))
	})
})

func getRemoteOK(context.Context, *testClient, testDTO) (testDTO, error) { return testDTO{}, nil }

func fullAdmission() groupAdmission {
	a := minimalAdmission(&testClient{})
	a.DryRun = dryRunOK
	a.GetRemote = getRemoteOK
	return a
}

var _ = Describe("NewAdmissionLifecycle", func() {
	It("returns a lifecycle with all required holes", func() {
		Expect(func() { lifecycle.NewAdmissionLifecycle(fullAdmission()) }).ToNot(Panic())
	})

	DescribeTable("panics on a missing required hole",
		func(unset func(*groupAdmission), message string) {
			a := fullAdmission()
			unset(&a)
			Expect(func() { lifecycle.NewAdmissionLifecycle(a) }).To(PanicWith(message))
		},
		Entry("ClientFactory", func(a *groupAdmission) { a.ClientFactory = nil },
			"lifecycle: AdmissionLifecycle.ClientFactory is required"),
		Entry("ToDTO", func(a *groupAdmission) { a.ToDTO = nil },
			"lifecycle: AdmissionLifecycle.ToDTO is required"),
		Entry("DryRun", func(a *groupAdmission) { a.DryRun = nil },
			"lifecycle: AdmissionLifecycle.DryRun is required"),
		Entry("GetRemote", func(a *groupAdmission) { a.GetRemote = nil },
			"lifecycle: AdmissionLifecycle.GetRemote is required"),
	)
})
