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

package apim

import (
	driftapim "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
)

var _ = Describe("IgnoreRemoteOnlyMetadata", func() {
	DescribeTable("should report equivalence and skip for empty slices",
		func(crd, remote any) {
			Expect(driftapim.IgnoreRemoteOnlyMetadata(crd, remote, drift.DriftContext{})).To(Equal(
				drift.Equivalence{Equivalent: drift.Equivalent, Skip: true},
			))
		},
		Entry("nil vs nil", nil, nil),
		Entry("nil vs empty slice", nil, []model.BaseMetadata{}),
		Entry("empty slice vs nil", []model.BaseMetadata{}, nil),
		Entry("empty slice vs empty slice", []model.BaseMetadata{}, []model.BaseMetadata{}),
	)

	DescribeTable("should report cannot compare without filter",
		func(crd, remote any) {
			e := driftapim.IgnoreRemoteOnlyMetadata(crd, remote, drift.DriftContext{})
			Expect(e.Equivalent).To(Equal(drift.CannotCompare))
			Expect(e.RemoteItemsFilterFunc).To(BeNil())
		},
		Entry("same metadata names", []model.BaseMetadata{{Name: "owner"}}, []model.BaseMetadata{{Name: "owner"}}),
		Entry("non-metadata slice items", []string{"foo"}, []string{"foo"}),
		Entry("crd metadata missing from remote", []model.BaseMetadata{{Name: "owner"}}, []model.BaseMetadata{}),
	)

	It("provides a filter that removes remote-only metadata by name", func() {
		crd := []model.BaseMetadata{{Name: "owner"}}
		remote := []model.BaseMetadata{
			{Name: "owner"},
			{Name: "sync-id"},
		}

		e := driftapim.IgnoreRemoteOnlyMetadata(crd, remote, drift.DriftContext{})
		Expect(e.Equivalent).To(Equal(drift.CannotCompare))
		Expect(e.RemoteItemsFilterFunc).NotTo(BeNil())

		filtered := e.RemoteItemsFilterFunc(remote)
		Expect(filtered).To(ConsistOf(model.BaseMetadata{Name: "owner"}))
	})

	It("filters remote-only metadata when crd metadata is empty", func() {
		remote := []model.BaseMetadata{{Name: "sync-id"}}

		e := driftapim.IgnoreRemoteOnlyMetadata(nil, remote, drift.DriftContext{})
		Expect(e.Equivalent).To(Equal(drift.CannotCompare))
		Expect(e.RemoteItemsFilterFunc).NotTo(BeNil())

		filtered := e.RemoteItemsFilterFunc(remote)
		Expect(filtered).To(BeEmpty())
	})

	It("keeps non-metadata items when filtering remote-only metadata", func() {
		crd := []model.BaseMetadata{{Name: "owner"}}
		remote := []model.BaseMetadata{
			{Name: "owner"},
			{Name: "sync-id"},
		}

		e := driftapim.IgnoreRemoteOnlyMetadata(crd, remote, drift.DriftContext{})
		Expect(e.RemoteItemsFilterFunc).NotTo(BeNil())

		mixed := []any{
			model.BaseMetadata{Name: "owner"},
			model.BaseMetadata{Name: "sync-id"},
			"keep-me",
		}
		filtered := e.RemoteItemsFilterFunc(mixed)
		Expect(filtered).To(ConsistOf(
			model.BaseMetadata{Name: "owner"},
			"keep-me",
		))
	})
})

var _ = Describe("IgnoreUnknownCRDGroups", func() {
	ctx := drift.DriftContext{Namespace: "test-namespace"}

	DescribeTable("should report equivalence and skip for empty slices",
		func(crd, remote any) {
			Expect(driftapim.IgnoreUnknownCRDGroups(crd, remote, ctx)).To(Equal(
				drift.Equivalence{Equivalent: drift.Equivalent, Skip: true},
			))
		},
		Entry("nil vs nil", nil, nil),
		Entry("nil vs empty slice", nil, []string{}),
		Entry("empty slice vs nil", []string{}, nil),
		Entry("empty slice vs empty slice", []string{}, []string{}),
	)

	It("should be equivalent when crd items match remote after filtering", func() {
		crd := []string{"group1", "group2"}
		remote := []string{"group1", "group2"}

		e := driftapim.IgnoreUnknownCRDGroups(crd, remote, ctx)
		Expect(e.Equivalent).To(Equal(drift.Equivalent))
		Expect(e.Skip).To(BeTrue())
	})

	It("should filter out CRD-only items not in remote", func() {
		crd := []string{"group1", "group2", "crd-only"}
		remote := []string{"group1", "group2"}

		e := driftapim.IgnoreUnknownCRDGroups(crd, remote, ctx)
		Expect(e.Equivalent).To(Equal(drift.Equivalent))
		Expect(e.Skip).To(BeTrue())
	})

	It("should be equivalent when comparing with namespace prefix ignored", func() {
		crd := []string{"group1"}
		remote := []string{"test-namespace-group1"}

		e := driftapim.IgnoreUnknownCRDGroups(crd, remote, ctx)
		Expect(e.Equivalent).To(Equal(drift.Equivalent))
		Expect(e.Skip).To(BeTrue())
	})

	It("should be equivalent when crd has namespace-prefixed items matching remote", func() {
		crd := []string{"test-namespace-group1"}
		remote := []string{"group1"}

		e := driftapim.IgnoreUnknownCRDGroups(crd, remote, ctx)
		Expect(e.Equivalent).To(Equal(drift.Equivalent))
		Expect(e.Skip).To(BeTrue())
	})

	It("should detect drift when remote membership differs after filtering CRD-only groups", func() {
		crd := []string{"group1", "group2"}
		remote := []string{"group1", "group3"}

		e := driftapim.IgnoreUnknownCRDGroups(crd, remote, ctx)
		Expect(e.Equivalent).To(Equal(drift.Inequivalent))
		Expect(e.Skip).To(BeTrue())
	})

	It("should detect drift when remote has extra groups", func() {
		crd := []string{"group1"}
		remote := []string{"group1", "group3"}

		e := driftapim.IgnoreUnknownCRDGroups(crd, remote, ctx)
		Expect(e.Equivalent).To(Equal(drift.Inequivalent))
		Expect(e.Skip).To(BeTrue())
	})

	It("should handle non-string slice items", func() {
		crd := []int{1, 2}
		remote := []int{1, 2}

		e := driftapim.IgnoreUnknownCRDGroups(crd, remote, ctx)
		Expect(e.Equivalent).To(Equal(drift.CannotCompare))
	})

	It("Detect reports drift when remote group membership differs", func() {
		crd := model.ApplicationDTO{Groups: []string{"group1", "group2"}}
		remote := model.ApplicationDTO{Groups: []string{"group1", "group3"}}
		result := drift.DetectWithNamespace(crd, remote, "test-namespace")
		Expect(result.DriftDetected()).To(BeTrue())
	})
})
