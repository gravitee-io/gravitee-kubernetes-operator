// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package authzentity_test

import (
	"context"
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzentity"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/authzentity"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
)

const ownershipFinding = `entityId [agents] is held by Group::"agents", which this manifest did not create ` +
	`(source: local). Declare another entityId, or delete that entity first`

const ns = "travel"

var devCtx = refs.NamespacedName{Name: "dev-ctx"}

func newEntity(name, entityID string, targets ...authzpdp.Target) *v1alpha1.AuthzEntity {
	return &v1alpha1.AuthzEntity{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns},
		Spec: v1alpha1.AuthzEntitySpec{
			Context: &devCtx,
			Type: authzentity.Type{
				EntityID: entityID, Kind: "PRINCIPAL", EntityType: "Group", Targets: targets,
			},
		},
	}
}

func byRef(name string) authzpdp.Target {
	return authzpdp.Target{PdpRef: &refs.NamespacedName{Name: name}}
}

func newPdp(name string, contextRef refs.NamespacedName, scope string) *v1alpha1.AuthzPdp {
	return &v1alpha1.AuthzPdp{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns},
		Spec: v1alpha1.AuthzPdpSpec{
			Context: &contextRef,
			Type:    authzpdp.Type{Name: name, TargetPdpID: name},
		},
		Status: v1alpha1.AuthzPdpStatus{Status: authzpdp.Status{RoutingScope: scope}},
	}
}

// cluster registers a fake client holding objects, AuthzEntity indexed as the manager's cache
// indexes it.
func cluster(objects ...client.Object) {
	c := &recordingCache{indexers: map[string]client.IndexerFunc{}}
	Expect(search.InitCache(context.Background(), c)).To(Succeed())
	scheme := runtime.NewScheme()
	Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
	builder := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...)
	for _, field := range []search.IndexField{
		search.AuthzEntityContextField, search.AuthzEntityAuthzPdpField, search.AuthzEntityParentField,
	} {
		builder = builder.WithIndex(&v1alpha1.AuthzEntity{}, field.String(), c.indexers[field.String()])
	}
	k8s.RegisterClient(builder.Build())
}

var _ = Describe("Validating an AuthzEntity", func() {
	ctx := context.Background()
	ctrl := admission.AdmissionCtrl{}

	BeforeEach(func() {
		platform.reset()
	})

	It("warns once, naming every pending AuthzPdp, and skips the dry run", func() {
		cluster(newPdp("orders-us", devCtx, ""))

		warnings, err := ctrl.ValidateCreate(ctx,
			newEntity("agents-group", "agents", byRef("orders-eu"), authzpdp.Target{Scope: "default"}, byRef("orders-us")))

		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(ConsistOf("authz pdp not found: [travel/orders-eu]; " +
			"authz pdp not synced: [travel/orders-us] has no routing scope yet: " +
			"the entity is checked against the platform once every AuthzPdp it targets has synced"))
		Expect(platform.requests()).To(BeEmpty())
	})

	// The controller's own finalizer and annotation writes go through the webhook: refusing them
	// while an AuthzPdp is pending would stop the controller from saving its finalizer.
	It("admits a metadata-only update while an AuthzPdp is pending", func() {
		cluster()
		oldObj := newEntity("agents-group", "agents", byRef("orders-eu"))
		oldObj.Status.ID = "3b9d0c2b"
		newObj := oldObj.DeepCopy()
		newObj.Finalizers = []string{"finalizers.gravitee.io/authzentities"}

		_, err := ctrl.ValidateUpdate(ctx, oldObj, newObj)

		Expect(err).ToNot(HaveOccurred())
		Expect(platform.requests()).To(BeEmpty())
	})

	It("refuses a pdpRef naming another kind, naming its index", func() {
		cluster()
		mcp := authzpdp.Target{PdpRef: &refs.NamespacedName{Name: "orders-eu", Kind: "McpProxy"}}

		_, err := ctrl.ValidateCreate(ctx, newEntity("agents-group", "agents", authzpdp.Target{Scope: "default"}, mcp))

		Expect(err).To(MatchError(ContainSubstring(
			"targets[1].pdpRef [travel/orders-eu]: kind must be AuthzPdp, got [McpProxy]")))
		Expect(platform.requests()).To(BeEmpty())
	})

	It("refuses an AuthzPdp of another management context", func() {
		cluster(newPdp("orders-eu", refs.NamespacedName{Name: "prod-ctx"}, "orders@eu"))

		_, err := ctrl.ValidateCreate(ctx, newEntity("agents-group", "agents", byRef("orders-eu")))

		Expect(err).To(MatchError(ContainSubstring(
			"authz pdp [travel/orders-eu] belongs to management context [travel/prod-ctx], not [travel/dev-ctx]")))
		Expect(platform.requests()).To(BeEmpty())
	})

	It("dry-runs the resolved targets in declaration order", func() {
		cluster(newPdp("orders-eu", devCtx, "orders@eu"), newPdp("orders-us", devCtx, "orders@us"))
		platform.answer(http.StatusOK, `{"hrid":"travel-agents-group"}`)

		_, err := ctrl.ValidateCreate(ctx,
			newEntity("agents-group", "agents", byRef("orders-us"), authzpdp.Target{Scope: "default"}, byRef("orders-eu")))

		Expect(err).ToNot(HaveOccurred())
		Expect(platform.requests()).To(HaveLen(1))
		sent := struct {
			Targets []string `json:"targets"`
		}{}
		Expect(json.Unmarshal([]byte(platform.requests()[0]), &sent)).To(Succeed())
		Expect(sent.Targets).To(Equal([]string{"orders@us", "default", "orders@eu"}))
	})

	It("refuses an entityId another source holds, naming that source", func() {
		cluster(newPdp("orders-eu", devCtx, "orders@eu"))
		finding, err := json.Marshal(ownershipFinding)
		Expect(err).ToNot(HaveOccurred())
		platform.answer(http.StatusOK, `{"hrid":"travel-agents-group","errors":{"severe":[`+string(finding)+`]}}`)

		_, err = ctrl.ValidateCreate(ctx, newEntity("agents-group", "agents", byRef("orders-eu")))

		Expect(err).To(MatchError(ContainSubstring(ownershipFinding)))
		Expect(err).To(MatchError(ContainSubstring("(source: local)")))
	})
})

var _ = Describe("Deleting an AuthzEntity", func() {
	ctx := context.Background()
	ctrl := admission.AdmissionCtrl{}

	child := func(name string, contextRef refs.NamespacedName, parents ...authzentity.Parent) *v1alpha1.AuthzEntity {
		entity := newEntity(name, name)
		entity.Spec.Context = &contextRef
		entity.Spec.Kind = "RESOURCE"
		entity.Spec.EntityType = "Hotel"
		entity.Spec.Parents = parents
		return entity
	}
	agents := authzentity.Parent{EntityType: "Group", EntityID: "agents"}

	BeforeEach(func() {
		platform.reset()
	})

	It("warns, naming the entities of its context its delete cascades to", func() {
		suite := child("ritz-suite", devCtx, authzentity.Parent{EntityType: "Hotel", EntityID: "ritz"})
		concierge := child("concierge", devCtx, authzentity.Parent{EntityType: "Hotel", EntityID: "ritz"})
		concierge.Spec.Kind, concierge.Spec.EntityType = "PRINCIPAL", "User"
		cluster(
			child("ritz", devCtx, agents),
			suite,
			concierge,
			child("grand-paris", refs.NamespacedName{Name: "dev-ctx", Namespace: ns},
				authzentity.Parent{EntityType: "Group", EntityID: "travel"}, agents),
			child("savoy", refs.NamespacedName{Name: "prod-ctx"}, agents),
			child("negresco", devCtx, authzentity.Parent{EntityType: "Team", EntityID: "agents"}),
		)

		warnings, err := ctrl.ValidateDelete(ctx, newEntity("agents-group", "agents"))

		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(ConsistOf("deleting authz entity [travel/agents-group] cascades on the platform to " +
			"AuthzEntity [travel/concierge, travel/grand-paris, travel/ritz, travel/ritz-suite]: each RESOURCE " +
			"among them is deleted and each PRINCIPAL detached. Each is then applied again: a PRINCIPAL as " +
			"declared, a RESOURCE refused until its parent exists again or is dropped from its parents"))
		Expect(platform.requests()).To(BeEmpty())
	})

	It("admits the delete of an entity no other declares as a parent, with no warning", func() {
		cluster(child("ritz", devCtx, authzentity.Parent{EntityType: "Group", EntityID: "travel"}))

		warnings, err := ctrl.ValidateDelete(ctx, newEntity("agents-group", "agents"))

		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(BeEmpty())
	})
})

var _ = Describe("Detecting drift on an AuthzEntity the platform no longer holds", func() {
	ctx := context.Background()
	ctrl := admission.AdmissionCtrl{}
	notFound := `{"message":"entity not found"}`

	BeforeEach(func() {
		platform.reset()
		platform.answer(http.StatusOK, `{"hrid":"travel-ritz"}`)
		platform.answerGet(http.StatusNotFound, notFound)
		saved := env.Config.DriftDetection
		env.Config.DriftDetection.Enabled = true
		env.Config.DriftDetection.Policy = env.DriftPolicyDeny
		env.Config.DriftDetection.OnRemoteMissing = env.DriftPolicyDeny
		DeferCleanup(func() { env.Config.DriftDetection = saved })
	})

	// ritz was applied under Group::"agents", then deleted with it on the platform.
	ritz := func() *v1alpha1.AuthzEntity {
		entity := newEntity("ritz", "ritz", authzpdp.Target{Scope: "default"})
		entity.Spec.Kind, entity.Spec.EntityType = "RESOURCE", "Hotel"
		entity.Spec.Parents = []authzentity.Parent{{EntityType: "Group", EntityID: "agents"}}
		entity.Status.ID = "3b9d0c2b"
		return entity
	}
	withoutParents := func(entity *v1alpha1.AuthzEntity) *v1alpha1.AuthzEntity {
		edited := entity.DeepCopy()
		edited.Spec.Parents = nil
		return edited
	}

	It("admits the edit that drops the deleted parent of a RESOURCE the platform refused", func() {
		refused := ritz()
		refused.Status.Errors.Severe = []string{`parent Group::"agents" does not exist: a RESOURCE is deleted ` +
			`with its parent, so declare the parent first or drop it from parents`}
		cluster(refused)

		_, err := ctrl.ValidateUpdate(ctx, refused, withoutParents(refused))

		Expect(err).ToNot(HaveOccurred())
		Expect(platform.requests()).To(HaveLen(1), "the new spec is dry-run")
	})

	It("refuses an edit to an entity the platform lost with no refusal recorded", func() {
		applied := ritz()
		cluster(applied)

		_, err := ctrl.ValidateUpdate(ctx, applied, withoutParents(applied))

		Expect(err).To(MatchError(ContainSubstring("[travel/ritz] not found during drift detection")))
	})
})
