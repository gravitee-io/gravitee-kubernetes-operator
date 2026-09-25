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
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
)

// testDTO is the minimal wire payload used across admission tests.
type testDTO struct {
	Key  string
	Name string
}

func (d testDTO) Identity() string { return d.Key }

// testClient stands in for *am.Client or *apim.APIM.
type testClient struct{}

func (c *testClient) Probe(_ context.Context) error { return nil }

func newGroup(name, ns string) *v1alpha1.Group {
	return &v1alpha1.Group{
		TypeMeta: metav1.TypeMeta{Kind: "Group", APIVersion: "gravitee.io/v1alpha1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: v1alpha1.GroupSpec{
			Context: &refs.NamespacedName{Name: "my-ctx", Namespace: ns},
		},
	}
}

func toDTO(g *v1alpha1.Group) testDTO {
	return testDTO{Key: g.Name, Name: g.Name}
}

func resolveCtx(tc *testClient) lifecycle.ClientFactoryFunc[*v1alpha1.Group, *testClient] {
	return func(_ context.Context, _ *v1alpha1.Group) (*testClient, error) {
		return tc, nil
	}
}

func failResolveCtx(msg string) lifecycle.ClientFactoryFunc[*v1alpha1.Group, *testClient] {
	return func(_ context.Context, _ *v1alpha1.Group) (*testClient, error) {
		return nil, fmt.Errorf("%s", msg)
	}
}

func noopRefs(_ context.Context, _ *v1alpha1.Group, _ string) error { return nil }

func minimalAdmission(tc *testClient) lifecycle.AdmissionLifecycle[*v1alpha1.Group, testDTO, *testClient] {
	return lifecycle.AdmissionLifecycle[*v1alpha1.Group, testDTO, *testClient]{
		ResolveRefs:   noopRefs,
		ClientFactory: resolveCtx(tc),
		ToDTO:         toDTO,
	}
}

// refGroup is a Group carrying one ref-tagged field. With SecretRef unset, the generic
// resolver fails, which makes "was resolution attempted" observable.
type refGroup struct {
	*v1alpha1.Group
	Secret    []byte               `json:"secret,omitempty" ref:"secret,SecretRef,key"`
	SecretRef *refs.NamespacedName `json:"secretRef,omitempty"`
}

var _ = Describe("AdmissionLifecycle without ResolveRefs", func() {
	It("does not resolve refs", func() {
		a := lifecycle.AdmissionLifecycle[*refGroup, testDTO, *testClient]{}

		errs := a.ValidateCreate(context.Background(), &refGroup{Group: newGroup("g1", "ns")})

		Expect(errs.IsSevere()).To(BeFalse(), "unexpected severe errors: %v", errs.Severe)
	})
})

var _ = Describe("AdmissionLifecycle", func() {

	ctx := context.Background()

	Describe("ValidateCreate", func() {
		It("passes with no optional hooks", func() {
			tc := &testClient{}
			a := minimalAdmission(tc)
			errs := a.ValidateCreate(ctx, newGroup("g1", "ns"))
			Expect(errs.IsSevere()).To(BeFalse())
		})

		It("stops on ResolveContext failure", func() {
			a := lifecycle.AdmissionLifecycle[*v1alpha1.Group, testDTO, *testClient]{
				ResolveRefs:   noopRefs,
				ClientFactory: failResolveCtx("unreachable"),
				ToDTO:         toDTO,
			}
			errs := a.ValidateCreate(ctx, newGroup("g1", "ns"))
			Expect(errs.IsSevere()).To(BeTrue())
			Expect(errs.Severe[0].Message).To(ContainSubstring("unreachable"))
		})

		It("runs PreCheck and stops on severe", func() {
			tc := &testClient{}
			a := minimalAdmission(tc)
			a.PreCheck = func(_ context.Context, _ *v1alpha1.Group) *gerrors.AdmissionErrors {
				e := gerrors.NewAdmissionErrors()
				e.AddSevere("pre-check failed")
				return e
			}
			errs := a.ValidateCreate(ctx, newGroup("g1", "ns"))
			Expect(errs.IsSevere()).To(BeTrue())
			Expect(errs.Severe[0].Message).To(Equal("pre-check failed"))
		})

		It("runs DryRun and reports errors", func() {
			tc := &testClient{}
			a := minimalAdmission(tc)
			a.DryRun = func(_ context.Context, _ *testClient, _ testDTO) *gerrors.AdmissionErrors {
				e := gerrors.NewAdmissionErrors()
				e.AddSevere("dry run rejected")
				return e
			}
			errs := a.ValidateCreate(ctx, newGroup("g1", "ns"))
			Expect(errs.IsSevere()).To(BeTrue())
			Expect(errs.Severe[0].Message).To(ContainSubstring("dry run rejected"))
		})

		It("runs PostCheck after DryRun", func() {
			tc := &testClient{}
			var order []string
			a := minimalAdmission(tc)
			a.DryRun = func(_ context.Context, _ *testClient, _ testDTO) *gerrors.AdmissionErrors {
				order = append(order, "dryrun")
				return gerrors.NewAdmissionErrors()
			}
			a.PostCheck = func(_ context.Context, _ *v1alpha1.Group) *gerrors.AdmissionErrors {
				order = append(order, "postcheck")
				return gerrors.NewAdmissionErrors()
			}
			errs := a.ValidateCreate(ctx, newGroup("g1", "ns"))
			Expect(errs.IsSevere()).To(BeFalse())
			Expect(order).To(Equal([]string{"dryrun", "postcheck"}))
		})

		It("collects warnings without rejecting", func() {
			tc := &testClient{}
			a := minimalAdmission(tc)
			a.PreCheck = func(_ context.Context, _ *v1alpha1.Group) *gerrors.AdmissionErrors {
				e := gerrors.NewAdmissionErrors()
				e.AddWarning("heads up")
				return e
			}
			errs := a.ValidateCreate(ctx, newGroup("g1", "ns"))
			Expect(errs.IsSevere()).To(BeFalse())
			Expect(errs.Warning).To(HaveLen(1))
		})
	})

	Describe("ValidateUpdate", func() {
		It("skips validation for an object being deleted", func() {
			a := minimalAdmission(&testClient{})
			a.DryRun = func(context.Context, *testClient, testDTO) *gerrors.AdmissionErrors {
				e := gerrors.NewAdmissionErrors()
				e.AddSevere("remote already deleted")
				return e
			}
			oldObj := newGroup("g1", "ns")
			newObj := oldObj.DeepCopy()
			newObj.DeletionTimestamp = &metav1.Time{Time: time.Now()}

			errs := a.ValidateUpdate(ctx, oldObj, newObj)

			Expect(errs.IsSevere()).To(BeFalse())
		})

		It("runs ImmutableFields between PreCheck and DryRun", func() {
			tc := &testClient{}
			var order []string
			a := minimalAdmission(tc)
			a.PreCheck = func(_ context.Context, _ *v1alpha1.Group) *gerrors.AdmissionErrors {
				order = append(order, "precheck")
				return gerrors.NewAdmissionErrors()
			}
			a.ImmutableFields = func(_, _ *v1alpha1.Group) *gerrors.AdmissionErrors {
				order = append(order, "immutable")
				return gerrors.NewAdmissionErrors()
			}
			a.DryRun = func(_ context.Context, _ *testClient, _ testDTO) *gerrors.AdmissionErrors {
				order = append(order, "dryrun")
				return gerrors.NewAdmissionErrors()
			}
			oldG := newGroup("g1", "ns")
			newG := newGroup("g1", "ns")
			errs := a.ValidateUpdate(ctx, oldG, newG)
			Expect(errs.IsSevere()).To(BeFalse())
			Expect(order).To(Equal([]string{"precheck", "immutable", "dryrun"}))
		})

		It("stops on ImmutableFields severe error", func() {
			tc := &testClient{}
			a := minimalAdmission(tc)
			dryRunCalled := false
			a.ImmutableFields = func(_, _ *v1alpha1.Group) *gerrors.AdmissionErrors {
				e := gerrors.NewAdmissionErrors()
				e.AddSevere("field is immutable")
				return e
			}
			a.DryRun = func(_ context.Context, _ *testClient, _ testDTO) *gerrors.AdmissionErrors {
				dryRunCalled = true
				return gerrors.NewAdmissionErrors()
			}
			errs := a.ValidateUpdate(ctx, newGroup("g1", "ns"), newGroup("g1", "ns"))
			Expect(errs.IsSevere()).To(BeTrue())
			Expect(dryRunCalled).To(BeFalse())
		})
	})

	Describe("ValidateDelete", func() {
		It("passes with no DeleteGuard", func() {
			tc := &testClient{}
			a := minimalAdmission(tc)
			errs := a.ValidateDelete(ctx, newGroup("g1", "ns"))
			Expect(errs.IsSevere()).To(BeFalse())
		})

		It("rejects when DeleteGuard returns an error", func() {
			tc := &testClient{}
			a := minimalAdmission(tc)
			a.DeleteGuard = func(_ context.Context, _ *v1alpha1.Group) error {
				return fmt.Errorf("still referenced")
			}
			errs := a.ValidateDelete(ctx, newGroup("g1", "ns"))
			Expect(errs.IsSevere()).To(BeTrue())
			Expect(errs.Severe[0].Message).To(ContainSubstring("still referenced"))
		})
	})

	Describe("Pipeline ordering", func() {
		It("executes full create pipeline in order", func() {
			tc := &testClient{}
			var order []string
			a := lifecycle.AdmissionLifecycle[*v1alpha1.Group, testDTO, *testClient]{
				ResolveRefs: func(_ context.Context, _ *v1alpha1.Group, _ string) error {
					order = append(order, "refs")
					return nil
				},
				ClientFactory: func(_ context.Context, _ *v1alpha1.Group) (*testClient, error) {
					order = append(order, "context")
					return tc, nil
				},
				ToDTO: toDTO,
				PreCheck: func(_ context.Context, _ *v1alpha1.Group) *gerrors.AdmissionErrors {
					order = append(order, "precheck")
					return gerrors.NewAdmissionErrors()
				},
				DryRun: func(_ context.Context, _ *testClient, _ testDTO) *gerrors.AdmissionErrors {
					order = append(order, "dryrun")
					return gerrors.NewAdmissionErrors()
				},
				PostCheck: func(_ context.Context, _ *v1alpha1.Group) *gerrors.AdmissionErrors {
					order = append(order, "postcheck")
					return gerrors.NewAdmissionErrors()
				},
			}
			errs := a.ValidateCreate(ctx, newGroup("g1", "ns"))
			Expect(errs.IsSevere()).To(BeFalse())
			Expect(order).To(Equal([]string{"refs", "context", "precheck", "dryrun", "postcheck"}))
		})

		It("skips DryRun and drift when DryRun callback is nil", func() {
			tc := &testClient{}
			var order []string
			a := lifecycle.AdmissionLifecycle[*v1alpha1.Group, testDTO, *testClient]{
				ResolveRefs:   noopRefs,
				ClientFactory: resolveCtx(tc),
				ToDTO:         toDTO,
				PreCheck: func(_ context.Context, _ *v1alpha1.Group) *gerrors.AdmissionErrors {
					order = append(order, "precheck")
					return gerrors.NewAdmissionErrors()
				},
				PostCheck: func(_ context.Context, _ *v1alpha1.Group) *gerrors.AdmissionErrors {
					order = append(order, "postcheck")
					return gerrors.NewAdmissionErrors()
				},
			}
			errs := a.ValidateCreate(ctx, newGroup("g1", "ns"))
			Expect(errs.IsSevere()).To(BeFalse())
			Expect(order).To(Equal([]string{"precheck", "postcheck"}))
		})
	})
})

var _ = Describe("AdmissionLifecycle drift", func() {
	ctx := context.Background()
	var enabled bool

	BeforeEach(func() {
		enabled = env.Config.DriftDetection.Enabled
		env.Config.DriftDetection.Enabled = true
	})

	AfterEach(func() {
		env.Config.DriftDetection.Enabled = enabled
	})

	withRemote := func(remote testDTO, err error) groupAdmission {
		a := minimalAdmission(&testClient{})
		a.GetRemote = func(context.Context, *testClient, testDTO) (testDTO, error) { return remote, err }
		return a
	}

	It("applies the remote-missing policy when the remote is not found", func() {
		a := withRemote(testDTO{}, gerrors.ServerError{StatusCode: 404})
		g := newGroup("g1", "ns")

		errs := a.ValidateUpdate(ctx, g, g.DeepCopy())

		Expect(errs.IsSevere()).To(BeTrue())
		Expect(errs.Severe[0].Error()).To(ContainSubstring("not found during drift detection"))
	})

	It("rejects an unchanged CR when the remote changed", func() {
		a := withRemote(testDTO{Key: "g1", Name: "changed remotely"}, nil)
		g := newGroup("g1", "ns")

		errs := a.ValidateUpdate(ctx, g, g.DeepCopy())

		Expect(errs.IsSevere()).To(BeTrue())
		Expect(errs.Severe[0].Error()).To(ContainSubstring("drift detected"))
	})

	It("accepts when the remote matches", func() {
		a := withRemote(testDTO{Key: "g1", Name: "g1"}, nil)
		g := newGroup("g1", "ns")

		errs := a.ValidateUpdate(ctx, g, g.DeepCopy())

		Expect(errs.IsSevere()).To(BeFalse())
	})
})
