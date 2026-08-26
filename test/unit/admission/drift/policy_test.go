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
	"context"
	stderrors "errors"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/dictionary"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	admdrift "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Drift policies fetch policies", func() {
	var original env.DriftDetection

	BeforeEach(func() {
		original = env.Config.DriftDetection
		env.Config.DriftDetection.Enabled = true
		env.Config.DriftDetection.OnRemoteMissing = env.DriftPolicyDeny
		env.Config.DriftDetection.OnFetchFailure = env.DriftPolicyDeny
	})

	AfterEach(func() {
		env.Config.DriftDetection = original
	})

	DescribeTable("not found (HTTP 404) uses OnRemoteMissing",
		func(policy env.DriftPolicy, expectSevere, expectWarning bool) {
			env.Config.DriftDetection.OnRemoteMissing = policy
			errs := validateWithRemoteError(errors.NewNotFoundError())

			if expectSevere {
				Expect(errs.Severe).To(HaveLen(1))
				Expect(errs.Severe[0].Message).To(ContainSubstring("not found during drift detection"))
			} else {
				Expect(errs.Severe).To(BeEmpty())
			}
			if expectWarning {
				Expect(errs.Warning).To(HaveLen(1))
				Expect(errs.Warning[0].Message).To(ContainSubstring("not found during drift detection"))
			} else {
				Expect(errs.Warning).To(BeEmpty())
			}
		},
		Entry("deny blocks the update", env.DriftPolicyDeny, true, false),
		Entry("warn allows with a kubectl warning", env.DriftPolicyWarn, false, true),
		Entry("allow reports nothing", env.DriftPolicyAllow, false, false),
	)

	DescribeTable("non-404 fetch failure uses OnFetchFailure",
		func(policy env.DriftPolicy, expectSevere, expectWarning bool) {
			env.Config.DriftDetection.OnFetchFailure = policy
			errs := validateWithRemoteError(stderrors.New("connection refused"))

			if expectSevere {
				Expect(errs.Severe).To(HaveLen(1))
				Expect(errs.Severe[0].Message).To(ContainSubstring("failed to fetch remote [Dictionary]"))
				Expect(errs.Severe[0].Message).To(ContainSubstring("connection refused"))
			} else {
				Expect(errs.Severe).To(BeEmpty())
			}
			if expectWarning {
				Expect(errs.Warning).To(HaveLen(1))
				Expect(errs.Warning[0].Message).To(ContainSubstring("failed to fetch remote [Dictionary]"))
			} else {
				Expect(errs.Warning).To(BeEmpty())
			}
		},
		Entry("deny blocks the update", env.DriftPolicyDeny, true, false),
		Entry("warn allows with a kubectl warning", env.DriftPolicyWarn, false, true),
		Entry("allow reports nothing", env.DriftPolicyAllow, false, false),
	)

	DescribeTable("When remote is fetch drift policy is used", func(policy env.DriftPolicy, expectSevere, expectWarning bool) {
		env.Config.DriftDetection.Policy = policy
		errs := validateWithDrift()

		if expectSevere {
			Expect(errs.Severe).To(HaveLen(1))
			Expect(errs.Severe[0].Message).To(ContainSubstring("drift detected"))
		} else {
			Expect(errs.Severe).To(BeEmpty())
		}
		if expectWarning {
			Expect(errs.Warning).To(HaveLen(1))
			Expect(errs.Warning[0].Message).To(ContainSubstring("drift detected"))
		} else {
			Expect(errs.Warning).To(BeEmpty())
		}
	},
		Entry("deny reports the drift", env.DriftPolicyDeny, true, false),
		Entry("warn reports the drift with a kubectl warning", env.DriftPolicyWarn, false, true),
		Entry("allow reports nothing", env.DriftPolicyAllow, false, false),
	)
})

func validateWithRemoteError(remoteErr error) *errors.AdmissionErrors {
	dict := &v1alpha1.Dictionary{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "test"},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Dictionary",
			APIVersion: "gravitee.io/v1alpha1",
		},
	}
	return admdrift.ValidateDriftWithContext(
		context.Background(),
		dict,
		dict.DeepCopy(),
		func(context.Context) (*apim.APIM, error) {
			return &apim.APIM{}, nil
		},
		func(context.Context, *v1alpha1.Dictionary) error {
			return nil
		},
		func(*apim.APIM, *v1alpha1.Dictionary) (any, error) {
			return nil, remoteErr
		},
		func(*v1alpha1.Dictionary) any {
			return struct{}{}
		},
	)
}

func validateWithDrift() *errors.AdmissionErrors {
	stored := &v1alpha1.Dictionary{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "test"},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Dictionary",
			APIVersion: "gravitee.io/v1alpha1",
		},
		Spec: v1alpha1.DictionarySpec{
			Type: dictionary.Type{
				Name: "original",
			},
		},
	}
	updated := stored.DeepCopy()
	updated.Spec.Name = "updated"
	remote := stored.DeepCopy()
	remote.Spec.Name = "remote"

	return admdrift.ValidateDriftWithContext(
		context.Background(),
		stored,
		updated,
		func(context.Context) (*apim.APIM, error) {
			return &apim.APIM{}, nil
		},
		func(context.Context, *v1alpha1.Dictionary) error {
			return nil
		},
		func(*apim.APIM, *v1alpha1.Dictionary) (any, error) {
			return model.ToDictionaryDTO(remote.Spec.Type, "default-test"), nil
		},
		func(crd *v1alpha1.Dictionary) any {
			return model.ToDictionaryDTO(crd.Spec.Type, "default-test")
		},
	)
}
