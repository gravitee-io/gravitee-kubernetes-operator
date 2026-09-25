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
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
)

const (
	testNs        = "default"
	testFinalizer = "finalizers.gravitee.io/lifecycle-test"
)

type upsertResponse struct{ ID string }

func (upsertResponse) GetOrgID() string { return "" }
func (upsertResponse) GetEnvID() string { return "" }

type groupLifecycle = lifecycle.ResourceLifecycle[*v1alpha1.Group, testDTO, *testClient, upsertResponse]

// remote records what the lifecycle sent to the fake remote API.
type remote struct {
	upserted  []testDTO
	deleted   []testDTO
	upsertErr error
}

func newGroupLifecycle(r *remote) groupLifecycle {
	return groupLifecycle{
		Finalizer:     testFinalizer,
		ResolveRefs:   noopRefs,
		ClientFactory: resolveCtx(&testClient{}),
		ToDTO:         toDTO,
		Upsert: func(_ context.Context, _ *testClient, dto testDTO) (upsertResponse, error) {
			r.upserted = append(r.upserted, dto)
			return upsertResponse{ID: "remote-id"}, r.upsertErr
		},
		Delete: func(_ context.Context, _ *testClient, dto testDTO) error {
			r.deleted = append(r.deleted, dto)
			return nil
		},
	}
}

func newCluster(objects ...client.Object) client.Client {
	scheme := runtime.NewScheme()
	Expect(clientgoscheme.AddToScheme(scheme)).To(Succeed())
	Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
	cli := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(objects...).
		WithStatusSubresource(&v1alpha1.Group{}).
		Build()
	k8s.RegisterClient(cli)
	return cli
}

func deletingGroup(name string, finalizers ...string) *v1alpha1.Group {
	g := newGroup(name, testNs)
	g.Finalizers = finalizers
	g.DeletionTimestamp = &metav1.Time{Time: time.Now()}
	return g
}

func reconcile(l groupLifecycle, cli client.Client, name string) (ctrl.Result, error) {
	req := ctrl.Request{NamespacedName: types.NamespacedName{Namespace: testNs, Name: name}}
	return l.Reconcile(context.Background(), cli, record.NewFakeRecorder(20), req, &v1alpha1.Group{})
}

func fetch(cli client.Client, name string) *v1alpha1.Group {
	g := &v1alpha1.Group{}
	Expect(cli.Get(context.Background(), types.NamespacedName{Namespace: testNs, Name: name}, g)).To(Succeed())
	return g
}

func condition(g *v1alpha1.Group, conditionType string) metav1.Condition {
	c, ok := g.GetConditions()[conditionType]
	Expect(ok).To(BeTrue(), "missing condition %s", conditionType)
	return c
}

func reconcileErrorType(err error) gerrors.ErrorType {
	var re gerrors.ReconcileError
	Expect(errors.As(err, &re)).To(BeTrue(), "not a ReconcileError: %v", err)
	return re.Type
}

var _ = Describe("ResourceLifecycle.Reconcile", func() {
	var r *remote

	BeforeEach(func() {
		r = &remote{}
	})

	Describe("create or update", func() {
		It("adds the finalizer and spec hash, upserts the DTO and reports success", func() {
			cli := newCluster(newGroup("g1", testNs))

			result, err := reconcile(newGroupLifecycle(r), cli, "g1")

			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(ctrl.Result{}))
			Expect(r.upserted).To(ConsistOf(testDTO{Key: "g1", Name: "g1"}))

			g := fetch(cli, "g1")
			Expect(g.Finalizers).To(ContainElement(testFinalizer))
			Expect(g.Annotations).To(HaveKeyWithValue(core.LastSpecHashAnnotation, hash.Calculate(g.GetSpec())))
			Expect(condition(g, k8s.ConditionAccepted).Status).To(Equal(metav1.ConditionTrue))
			Expect(condition(g, k8s.ConditionResolvedRefs).Status).To(Equal(metav1.ConditionTrue))
		})

		It("persists status set by PostUpsert", func() {
			cli := newCluster(newGroup("g1", testNs))
			l := newGroupLifecycle(r)
			l.PostUpsert = func(_ context.Context, g *v1alpha1.Group, resp upsertResponse) error {
				g.Status.ID = resp.ID
				return nil
			}

			_, err := reconcile(l, cli, "g1")

			Expect(err).ToNot(HaveOccurred())
			Expect(fetch(cli, "g1").Status.ID).To(Equal("remote-id"))
		})

		It("does not wrap a PostUpsert error as a control plane error", func() {
			cli := newCluster(newGroup("g1", testNs))
			l := newGroupLifecycle(r)
			l.PostUpsert = func(context.Context, *v1alpha1.Group, upsertResponse) error {
				return errors.New("post upsert failed")
			}

			_, err := reconcile(l, cli, "g1")

			Expect(err).To(MatchError("post upsert failed"))
			accepted := condition(fetch(cli, "g1"), k8s.ConditionAccepted)
			Expect(accepted.Status).To(Equal(metav1.ConditionFalse))
			Expect(accepted.Reason).To(Equal("ReconcileFailed"))
		})
	})

	DescribeTable("classifies Upsert errors",
		func(upsertErr error, wantType gerrors.ErrorType, wantRequeue bool) {
			cli := newCluster(newGroup("g1", testNs))
			r.upsertErr = upsertErr

			_, err := reconcile(newGroupLifecycle(r), cli, "g1")

			if wantRequeue {
				Expect(reconcileErrorType(err)).To(Equal(wantType))
			} else {
				Expect(err).ToNot(HaveOccurred())
			}
			accepted := condition(fetch(cli, "g1"), k8s.ConditionAccepted)
			Expect(accepted.Status).To(Equal(metav1.ConditionFalse))
			Expect(accepted.Reason).To(Equal(string(wantType)))
		},
		Entry("plain error is wrapped as a recoverable control plane error",
			errors.New("boom"), gerrors.ControlPlaneError, true),
		Entry("recoverable server error requeues",
			gerrors.ServerError{StatusCode: 503}, gerrors.ControlPlaneError, true),
		Entry("unrecoverable server error aborts",
			gerrors.ServerError{StatusCode: 400}, gerrors.ControlPlaneError, false),
		Entry("a ReconcileError passes through unchanged",
			gerrors.NewIllegalStateError(errors.New("bad state")), gerrors.IllegalStateError, false),
	)

	Describe("delete", func() {
		It("deletes the remote DTO and removes the finalizer", func() {
			cli := newCluster(deletingGroup("g1", testFinalizer))

			_, err := reconcile(newGroupLifecycle(r), cli, "g1")

			Expect(err).ToNot(HaveOccurred())
			Expect(r.deleted).To(ConsistOf(testDTO{Key: "g1", Name: "g1"}))
			err = cli.Get(context.Background(), types.NamespacedName{Namespace: testNs, Name: "g1"}, &v1alpha1.Group{})
			Expect(apierrors.IsNotFound(err)).To(BeTrue())
		})

		It("keeps the finalizer and skips Delete when DeleteGuard refuses", func() {
			cli := newCluster(deletingGroup("g1", testFinalizer))
			l := newGroupLifecycle(r)
			l.DeleteGuard = func(context.Context, *v1alpha1.Group) error {
				return errors.New("still referenced")
			}

			_, err := reconcile(l, cli, "g1")

			Expect(err).To(MatchError("still referenced"))
			Expect(r.deleted).To(BeEmpty())
			g := fetch(cli, "g1")
			Expect(g.Finalizers).To(ContainElement(testFinalizer))
			Expect(condition(g, k8s.ConditionAccepted).Message).To(ContainSubstring("still referenced"))
		})
	})

	It("does not delete when its finalizer is already gone", func() {
		cli := newCluster(deletingGroup("g1", "other/finalizer"))

		_, err := reconcile(newGroupLifecycle(r), cli, "g1")

		Expect(err).ToNot(HaveOccurred())
		Expect(r.deleted).To(BeEmpty())
	})

	Describe("delete with a reference already gone", func() {
		It("releases the finalizer without calling Delete", func() {
			cli := newCluster(deletingGroup("g1", testFinalizer))
			l := newGroupLifecycle(r)
			l.ResolveRefs = func(context.Context, *v1alpha1.Group, string) error {
				notFound := apierrors.NewNotFound(corev1.Resource("secrets"), "token")
				return fmt.Errorf("lifecycle/ref: resolve: %w", notFound)
			}

			_, err := reconcile(l, cli, "g1")

			Expect(err).ToNot(HaveOccurred())
			Expect(r.deleted).To(BeEmpty())
			err = cli.Get(context.Background(), types.NamespacedName{Namespace: testNs, Name: "g1"}, &v1alpha1.Group{})
			Expect(apierrors.IsNotFound(err)).To(BeTrue())
		})

		It("still requeues on any other resolution error", func() {
			cli := newCluster(deletingGroup("g1", testFinalizer))
			l := newGroupLifecycle(r)
			l.ResolveRefs = func(context.Context, *v1alpha1.Group, string) error {
				return errors.New("api server unavailable")
			}

			_, err := reconcile(l, cli, "g1")

			Expect(reconcileErrorType(err)).To(Equal(gerrors.ResolveRefError))
			Expect(r.deleted).To(BeEmpty())
			Expect(fetch(cli, "g1").Finalizers).To(ContainElement(testFinalizer))
		})
	})

	Describe("client factory failure", func() {
		It("requeues a create with a ResolveRef error and does not upsert", func() {
			cli := newCluster(newGroup("g1", testNs))
			l := newGroupLifecycle(r)
			l.ClientFactory = failResolveCtx("context not found")

			_, err := reconcile(l, cli, "g1")

			Expect(reconcileErrorType(err)).To(Equal(gerrors.ResolveRefError))
			Expect(r.upserted).To(BeEmpty())
			refs := condition(fetch(cli, "g1"), k8s.ConditionResolvedRefs)
			Expect(refs.Status).To(Equal(metav1.ConditionFalse))
			Expect(refs.Message).To(ContainSubstring("context not found"))
		})

		It("requeues a delete, keeping the finalizer and skipping Delete", func() {
			cli := newCluster(deletingGroup("g1", testFinalizer))
			l := newGroupLifecycle(r)
			l.ClientFactory = failResolveCtx("context not found")

			_, err := reconcile(l, cli, "g1")

			Expect(reconcileErrorType(err)).To(Equal(gerrors.ResolveRefError))
			Expect(r.deleted).To(BeEmpty())
			Expect(fetch(cli, "g1").Finalizers).To(ContainElement(testFinalizer))
		})
	})

	Describe("ResolveRefs", func() {
		It("sets no ResolvedRefs condition on success when ResolveRefs is nil", func() {
			cli := newCluster(newGroup("g1", testNs))
			l := newGroupLifecycle(r)
			l.ResolveRefs = nil

			_, err := reconcile(l, cli, "g1")

			Expect(err).ToNot(HaveOccurred())
			Expect(r.upserted).To(HaveLen(1))
			g := fetch(cli, "g1")
			Expect(condition(g, k8s.ConditionAccepted).Status).To(Equal(metav1.ConditionTrue))
			Expect(g.GetConditions()).ToNot(HaveKey(k8s.ConditionResolvedRefs))
		})

		It("sets no ResolvedRefs condition on failure when ResolveRefs is nil", func() {
			cli := newCluster(newGroup("g1", testNs))
			l := newGroupLifecycle(r)
			l.ResolveRefs = nil
			r.upsertErr = errors.New("boom")

			_, err := reconcile(l, cli, "g1")

			Expect(err).To(HaveOccurred())
			g := fetch(cli, "g1")
			Expect(condition(g, k8s.ConditionAccepted).Status).To(Equal(metav1.ConditionFalse))
			Expect(g.GetConditions()).ToNot(HaveKey(k8s.ConditionResolvedRefs))
		})

		It("fails with a ResolveRef error and does not upsert", func() {
			cli := newCluster(newGroup("g1", testNs))
			l := newGroupLifecycle(r)
			l.ResolveRefs = func(context.Context, *v1alpha1.Group, string) error {
				return errors.New("secret missing")
			}

			_, err := reconcile(l, cli, "g1")

			Expect(reconcileErrorType(err)).To(Equal(gerrors.ResolveRefError))
			Expect(r.upserted).To(BeEmpty())
			refs := condition(fetch(cli, "g1"), k8s.ConditionResolvedRefs)
			Expect(refs.Status).To(Equal(metav1.ConditionFalse))
			Expect(refs.Message).To(ContainSubstring("secret missing"))
		})
	})

	It("requeues on a template compile error without upserting", func() {
		g := newGroup("g1", testNs)
		g.Spec.Type = &group.Type{Name: "[[ secret `missing/token` ]]"}
		cli := newCluster(g)

		_, err := reconcile(newGroupLifecycle(r), cli, "g1")

		// the Secret may be created later: retry like the other controllers do
		Expect(err).To(HaveOccurred())
		Expect(r.upserted).To(BeEmpty())
		Expect(condition(fetch(cli, "g1"), k8s.ConditionAccepted).Status).To(Equal(metav1.ConditionFalse))
	})
})
