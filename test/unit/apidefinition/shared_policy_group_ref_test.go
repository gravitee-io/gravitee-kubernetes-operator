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

package apidefinition_test

import (
	"context"
	"encoding/json"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/sharedpolicygroups"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apidefinition"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"github.com/onsi/gomega/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var _ = Describe("PrepareV4SpecForAutomation with a sharedPolicyGroupRef", func() {
	var (
		ctx             context.Context
		scheme          *runtime.Scheme
		prevEnableTempl bool
	)

	BeforeEach(func() {
		ctx = context.Background()
		prevEnableTempl = env.Config.EnableTemplating
		env.Config.EnableTemplating = true

		scheme = runtime.NewScheme()
		Expect(clientgoscheme.AddToScheme(scheme)).To(Succeed())
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
	})

	AfterEach(func() {
		env.Config.EnableTemplating = prevEnableTempl
	})

	registerClient := func(objects ...client.Object) {
		k8s.RegisterClient(fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).Build())
	}

	newSPG := func() *v1alpha1.SharedPolicyGroup {
		phase := sharedpolicygroups.FlowPhase("REQUEST")
		return &v1alpha1.SharedPolicyGroup{
			ObjectMeta: metav1.ObjectMeta{Name: "spg", Namespace: "default"},
			Spec: v1alpha1.SharedPolicyGroupSpec{
				SharedPolicyGroup: &sharedpolicygroups.SharedPolicyGroup{Name: "spg", Phase: &phase},
			},
		}
	}

	spgRefStep := func() *v4.FlowStep {
		return &v4.FlowStep{
			FlowStep:          base.FlowStep{Enabled: true},
			SharedPolicyGroup: &refs.NamespacedName{Name: "spg"},
		}
	}

	newAPI := func(flows []*v4.Flow, plans *map[string]*v4.Plan) *v1alpha1.ApiV4Definition {
		return &v1alpha1.ApiV4Definition{
			ObjectMeta: metav1.ObjectMeta{Name: "api", Namespace: "default"},
			Spec: v1alpha1.ApiV4DefinitionSpec{
				Api: v4.Api{
					V4BaseApi: &v4.V4BaseApi{
						ApiBase: &base.ApiBase{Name: "api", Version: "1.0"},
						Flows:   flows,
					},
					Plans: plans,
				},
			},
		}
	}

	DescribeTable("sends the SPG step with its policy and identifier",
		func(setState func(*v1alpha1.SharedPolicyGroup), expectedStep string) {
			spg := newSPG()
			setState(spg)
			registerClient(spg)
			api := newAPI([]*v4.Flow{{Enabled: true, Request: []*v4.FlowStep{spgRefStep()}}}, nil)

			Expect(apidefinition.PrepareV4SpecForAutomation(ctx, api, false)).To(Succeed())

			dto := model.ToAPIV4DTO(&api.Spec.Api)
			Expect(json.Marshal(dto.Flows[0].Request[0])).To(MatchJSON(expectedStep))
		},
		Entry("managed by the Automation API",
			func(spg *v1alpha1.SharedPolicyGroup) { k8s.AddAutomationAPIManagedCondition(spg) },
			`{"enabled":true,"policy":"shared-policy-group-policy","name":"spg","configuration":{"hrid":"default-spg"}}`,
		),
		Entry("legacy, addressed by its cross ID",
			func(spg *v1alpha1.SharedPolicyGroup) { spg.Status.CrossID = "abc" },
			`{"enabled":true,"policy":"shared-policy-group-policy","name":"spg",`+
				`"configuration":{"sharedPolicyGroupId":"abc"}}`,
		),
		Entry("not reconciled yet",
			func(*v1alpha1.SharedPolicyGroup) {},
			`{"enabled":true,"policy":"shared-policy-group-policy","name":"spg","configuration":{"hrid":"default-spg"}}`,
		),
	)

	It("resolves an SPG ref in a plan flow", func() {
		spg := newSPG()
		k8s.AddAutomationAPIManagedCondition(spg)
		registerClient(spg)
		plans := map[string]*v4.Plan{
			"default": {
				Plan:  &base.Plan{},
				Name:  "default",
				Flows: []*v4.Flow{{Enabled: true, Response: []*v4.FlowStep{spgRefStep()}}},
			},
		}
		api := newAPI(nil, &plans)

		Expect(apidefinition.PrepareV4SpecForAutomation(ctx, api, false)).To(Succeed())

		dto := model.ToAPIV4DTO(&api.Spec.Api)
		Expect(json.Marshal(dto.Plans[0].Flows[0].Response[0])).To(MatchJSON(
			`{"enabled":true,"policy":"shared-policy-group-policy","name":"spg","configuration":{"hrid":"default-spg"}}`,
		))
	})

	DescribeTable("touches Secrets templated in the SPG only when updateMetadata is true",
		func(updateMetadata bool, annotations types.GomegaMatcher) {
			spg := newSPG()
			spg.Spec.Steps = []*sharedpolicygroups.Step{{
				Enabled: true,
				Configuration: utils.ToGenericStringMap(map[string]interface{}{
					"value": "[[ secret `spg-secret/key` ]]",
				}),
			}}
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "spg-secret", Namespace: "default"},
				Data:       map[string][]byte{"key": []byte("s3cr3t")},
			}
			registerClient(spg, secret)
			api := newAPI([]*v4.Flow{{Enabled: true, Request: []*v4.FlowStep{spgRefStep()}}}, nil)

			Expect(apidefinition.PrepareV4SpecForAutomation(ctx, api, updateMetadata)).To(Succeed())

			Expect(k8s.GetClient().Get(ctx, client.ObjectKeyFromObject(secret), secret)).To(Succeed())
			Expect(util.ContainsFinalizer(secret, core.TemplatingFinalizer)).To(Equal(updateMetadata))
			Expect(secret.Annotations).To(annotations)
		},
		Entry("admission (false) leaves the Secret untouched", false, BeEmpty()),
		Entry("reconcile (true) tracks the reference", true, HaveKeyWithValue("gravitee.io/references", "1")),
	)
})
