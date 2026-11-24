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
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/navigation"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

var _ = Describe("Portal navigation dry-run findings", func() {
	ctx := context.Background()

	It("reports a failing dry run as one severe finding with the error's own message", func() {
		dryRun := navigation.FindingsOnError(func(context.Context, any, any) (*gerrors.AdmissionErrors, error) {
			return nil, errors.New("APIM answered 502")
		})
		findings, err := dryRun(ctx, nil, nil)
		Expect(err).ToNot(HaveOccurred())
		Expect(findings.Severe).To(HaveLen(1))
		Expect(findings.Severe[0].Error()).To(Equal("APIM answered 502"))
		Expect(findings.Warning).To(BeEmpty())
	})

	It("passes the findings of a dry run that succeeds through", func() {
		expected := gerrors.NewAdmissionErrors()
		expected.AddWarning("warning one")
		dryRun := navigation.FindingsOnError(func(context.Context, any, any) (*gerrors.AdmissionErrors, error) {
			return expected, nil
		})
		findings, err := dryRun(ctx, nil, nil)
		Expect(err).ToNot(HaveOccurred())
		Expect(findings).To(BeIdenticalTo(expected))
	})
})
