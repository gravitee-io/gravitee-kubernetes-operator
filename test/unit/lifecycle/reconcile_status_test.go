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
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

var _ = Describe("ResourceLifecycle.Reconcile Status hook", func() {
	var r *remote

	BeforeEach(func() {
		r = &remote{}
	})

	// The hook writes the deprecated processingStatus, as ApiV4Definition does.
	withProcessingStatus := func(l groupLifecycle, seen *[]error) groupLifecycle {
		l.Status = func(g *v1alpha1.Group, err error) {
			*seen = append(*seen, err)
			if err != nil {
				g.Status.ProcessingStatus = core.ProcessingStatusFailed
			} else {
				g.Status.ProcessingStatus = core.ProcessingStatusCompleted
			}
		}
		return l
	}

	It("runs on success, and what it sets is saved with the status", func() {
		cli := newCluster(newGroup("g1", testNs))
		var seen []error

		_, err := reconcile(withProcessingStatus(newGroupLifecycle(r), &seen), cli, "g1")

		Expect(err).ToNot(HaveOccurred())
		Expect(seen).To(Equal([]error{nil}))
		Expect(fetch(cli, "g1").Status.ProcessingStatus).To(Equal(core.ProcessingStatusCompleted))
	})

	It("runs on failure with the reconcile error, and what it sets is saved with the status", func() {
		cli := newCluster(newGroup("g1", testNs))
		r.upsertErr = errors.New("boom")
		var seen []error

		_, _ = reconcile(withProcessingStatus(newGroupLifecycle(r), &seen), cli, "g1")

		Expect(seen).To(HaveLen(1))
		Expect(seen[0]).To(MatchError(ContainSubstring("boom")))
		Expect(fetch(cli, "g1").Status.ProcessingStatus).To(Equal(core.ProcessingStatusFailed))
	})

	It("runs once, with the error, when a reference fails to resolve before the upsert", func() {
		cli := newCluster(newGroup("g1", testNs))
		l := newGroupLifecycle(r)
		l.ResolveRefs = func(context.Context, *v1alpha1.Group, string) error { return errors.New("missing ref") }
		var seen []error

		_, _ = reconcile(withProcessingStatus(l, &seen), cli, "g1")

		Expect(r.upserted).To(BeEmpty())
		Expect(seen).To(HaveLen(1))
		Expect(fetch(cli, "g1").Status.ProcessingStatus).To(Equal(core.ProcessingStatusFailed))
	})
})
