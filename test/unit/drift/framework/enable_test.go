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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Drift enablement", func() {
	var originalConfig bool

	BeforeEach(func() {
		originalConfig = env.Config.DriftDetection.Enabled
	})

	AfterEach(func() {
		env.Config.DriftDetection.Enabled = originalConfig
	})

	DescribeTable(
		"unsupported predicate always excludes (legacy group)",
		func(cfg bool, annotatedValue *string) {
			env.Config.DriftDetection.Enabled = cfg

			g := &v1alpha1.Group{}
			if annotatedValue != nil {
				g.Annotations = map[string]string{core.DriftDetectionAnnotation: *annotatedValue}
			}

			Expect(drift.IsDriftEnabled(g)).To(BeFalse())
		},
		Entry("config=true, annotation=true", true, new(env.TrueString)),
		Entry("config=true, annotation=false", true, new(env.FalseString)),
		Entry("config=true, annotation missing", true, nil),
		Entry("config=false, annotation=true", false, new(env.TrueString)),
	)

	DescribeTable(
		"unsupported predicate always excludes (subscription on definition-v2 API)",
		func(cfg bool, annotatedValue *string) {
			env.Config.DriftDetection.Enabled = cfg

			sub := &v1alpha1.Subscription{}
			sub.Spec.API.Kind = "ApiDefinition"
			if annotatedValue != nil {
				sub.Annotations = map[string]string{core.DriftDetectionAnnotation: *annotatedValue}
			}

			Expect(drift.IsDriftEnabled(sub)).To(BeFalse())
		},
		Entry("config=true, annotation=true", true, new(env.TrueString)),
		Entry("config=true, annotation=false", true, new(env.FalseString)),
		Entry("config=true, annotation missing", true, nil),
		Entry("config=false, annotation=true", false, new(env.TrueString)),
	)

	It("keeps drift enabled for a subscription on a v4 API", func() {
		env.Config.DriftDetection.Enabled = true
		sub := &v1alpha1.Subscription{}
		sub.Spec.API.Kind = "ApiV4Definition"
		Expect(drift.IsDriftEnabled(sub)).To(BeTrue())
	})

	DescribeTable(
		"supported CRD uses annotation if present, otherwise falls back to config",
		func(cfg bool, annotatedValue string, expected bool) {
			env.Config.DriftDetection.Enabled = cfg

			d := &v1alpha1.Dictionary{}
			if annotatedValue != "" {
				d.Annotations = map[string]string{core.DriftDetectionAnnotation: annotatedValue}
			}

			Expect(drift.IsDriftEnabled(d)).To(Equal(expected))
		},
		Entry("config=true, annotation missing -> enabled", true, "", true),
		Entry("config=false, annotation missing -> disabled", false, "", false),
		Entry("config=false, annotation=true -> enabled (annotation overrides config)", false, env.TrueString, true),
		Entry("config=true, annotation=false -> disabled (annotation overrides config)", true, env.FalseString, false),
		Entry("config=true, annotation invalid -> enabled (falls back to config)", true, "invalid", true),
		Entry("config=false, annotation invalid -> disabled (falls back to config)", false, "invalid", false),
	)
})
