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

package drift

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ = Describe("Drift enablement", func() {
	var originalConfig bool

	BeforeEach(func() {
		originalConfig = env.Config.DriftDetection
	})

	AfterEach(func() {
		env.Config.DriftDetection = originalConfig
	})

	DescribeTable(
		"unsupported predicate always excludes (legacy group)",
		func(cfg bool, annotatedValue *string) {
			env.Config.DriftDetection = cfg

			g := &v1alpha1.Group{}
			if annotatedValue != nil {
				g.Annotations = map[string]string{core.DriftDetectionAnnotation: *annotatedValue}
			}

			Expect(drift.IsDriftEnabled(g)).To(BeFalse())
		},
		Entry("config=true, annotation=true", true, ptr(env.TrueString)),
		Entry("config=true, annotation=false", true, ptr(env.FalseString)),
		Entry("config=true, annotation missing", true, nil),
		Entry("config=false, annotation=true", false, ptr(env.TrueString)),
	)

	DescribeTable(
		"disabled predicate excludes only when annotation is missing/invalid",
		func(makeCRD func() runtime.Object, cfg bool, annotatedValue string, expected bool) {
			env.Config.DriftDetection = cfg

			p := makeCRD()
			if annotatedValue != "" {
				switch t := p.(type) {
				case *v1alpha1.Portal:
					t.Annotations = map[string]string{core.DriftDetectionAnnotation: annotatedValue}
				case *v1alpha1.Documentation:
					t.Annotations = map[string]string{core.DriftDetectionAnnotation: annotatedValue}
				case *v1alpha1.PortalListing:
					t.Annotations = map[string]string{core.DriftDetectionAnnotation: annotatedValue}
				default:
					Fail("unexpected CRD type in test table")
				}
			}

			Expect(drift.IsDriftEnabled(p)).To(Equal(expected))
		},
		Entry("Portal: config=true, annotation missing -> excluded by predicate",
			func() runtime.Object { return &v1alpha1.Portal{} }, true, "", false),
		Entry("Documentation: config=true, annotation missing -> excluded by predicate",
			func() runtime.Object { return &v1alpha1.Documentation{} }, true, "", false),
		Entry("PortalListing: config=true, annotation missing -> excluded by predicate",
			func() runtime.Object { return &v1alpha1.PortalListing{} }, true, "", false),
		Entry("Portal: config=true, annotation invalid -> excluded by predicate",
			func() runtime.Object { return &v1alpha1.Portal{} }, true, "invalid", false),
		Entry("Portal: config=true, annotation=true -> enabled (annotation overrides predicate)",
			func() runtime.Object { return &v1alpha1.Portal{} }, true, env.TrueString, true),
		Entry("Portal: config=true, annotation=false -> disabled (annotation overrides predicate)",
			func() runtime.Object { return &v1alpha1.Portal{} }, true, env.FalseString, false),
		Entry("Portal: config=false, annotation=true -> enabled (annotation overrides predicate)",
			func() runtime.Object { return &v1alpha1.Portal{} }, false, env.TrueString, true),
	)

	DescribeTable(
		"supported CRD uses annotation if present, otherwise falls back to config",
		func(cfg bool, annotatedValue string, expected bool) {
			env.Config.DriftDetection = cfg

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
