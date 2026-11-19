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

package apim_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	spgmodel "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/sharedpolicygroups"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/apidefinition"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/authzentity"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/authzpolicy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/authzschema"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/dictionary"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/sharedpolicygroups"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
)

// The refusing APIM answers a DELETE with a 401, which a create or an update does not retry. A
// delete is retried whatever the status code, so the CR never stays Terminating with nothing left to
// retry it, and Accepted names the control plane's refusal.
var _ = DescribeTable("A delete the control plane refused",
	func(obj client.Object, finalizer, path string, newReconciler func(client.Client) reconcile.Reconciler) {
		obj.SetFinalizers([]string{finalizer})
		obj.SetDeletionTimestamp(&metav1.Time{Time: time.Now()})
		cli := deletingCluster(obj)
		key := client.ObjectKeyFromObject(obj)

		_, err := newReconciler(cli).Reconcile(context.Background(), ctrl.Request{NamespacedName: key})

		Expect(err).To(HaveOccurred())
		Expect(cli.Get(context.Background(), key, obj)).To(Succeed())
		Expect(obj.GetFinalizers()).To(ContainElement(finalizer))
		conditionAware, ok := obj.(core.ConditionAware)
		Expect(ok).To(BeTrue())
		accepted := k8s.GetCondition(conditionAware, k8s.ConditionAccepted)
		Expect(accepted).ToNot(BeNil())
		Expect(accepted.Status).To(Equal(metav1.ConditionFalse))
		Expect(accepted.Reason).To(Equal("ControlPlaneError"))
		Expect(accepted.Message).To(ContainSubstring(refusedDeleteMessage))
		Expect(accepted.Message).To(ContainSubstring("[DELETE] " + platformURL + path))
	},
	Entry("ApiV4Definition", apiV4On(refusingCtx), core.ApiDefinitionFinalizer, automation+"/apis/gravitee-orders?",
		func(cli client.Client) reconcile.Reconciler {
			return &apidefinition.V4Reconciler{Client: cli, Recorder: record.NewFakeRecorder(20)}
		}),
	Entry("ApiDefinition", apiV2On(refusingCtx), core.ApiDefinitionFinalizer,
		"/management/organizations/DEFAULT/environments/DEFAULT/apis/8f0c",
		func(cli client.Client) reconcile.Reconciler {
			return &apidefinition.Reconciler{Client: cli, Recorder: record.NewFakeRecorder(20)}
		}),
	Entry("Application", applicationOn(refusingCtx), core.ApplicationFinalizer,
		automation+"/applications/gravitee-orders-app?",
		func(cli client.Client) reconcile.Reconciler {
			return &application.Reconciler{Client: cli, Recorder: record.NewFakeRecorder(20)}
		}),
	Entry("Dictionary", dictionaryOn(refusingCtx), core.DictionaryFinalizer, automation+"/dictionaries/gravitee-orders-dict",
		func(cli client.Client) reconcile.Reconciler {
			return &dictionary.Reconciler{Client: cli, Recorder: record.NewFakeRecorder(20)}
		}),
	Entry("SharedPolicyGroup", sharedPolicyGroupOn(refusingCtx), core.SharedPolicyGroupFinalizer,
		automation+"/shared-policy-groups/gravitee-rate-limits?",
		func(cli client.Client) reconcile.Reconciler {
			return &sharedpolicygroups.Reconciler{Client: cli, Recorder: record.NewFakeRecorder(20)}
		}),
	Entry("AuthzEntity", authzEntityOn(refusingCtx), core.AuthzEntityFinalizer,
		automation+"/authz/entities/gravitee-agents-group",
		func(cli client.Client) reconcile.Reconciler {
			return &authzentity.Reconciler{Client: cli, Recorder: record.NewFakeRecorder(20)}
		}),
	Entry("AuthzPdp", authzPdpOn(refusingCtx), core.AuthzPdpFinalizer, automation+"/authz/pdps/gravitee-orders-eu",
		func(cli client.Client) reconcile.Reconciler {
			return &authzpdp.Reconciler{Client: cli, Recorder: record.NewFakeRecorder(20)}
		}),
	Entry("AuthzPolicy", authzPolicyOn(refusingCtx), core.AuthzPolicyFinalizer,
		automation+"/authz/policies/gravitee-hotel-booking-read",
		func(cli client.Client) reconcile.Reconciler {
			return &authzpolicy.Reconciler{Client: cli, Recorder: record.NewFakeRecorder(20)}
		}),
	Entry("AuthzSchema", authzSchemaOn(refusingCtx), core.AuthzSchemaFinalizer,
		automation+"/authz/schemas/gravitee-travel-schema",
		func(cli client.Client) reconcile.Reconciler {
			return &authzschema.Reconciler{Client: cli, Recorder: record.NewFakeRecorder(20)}
		}),
)

// The accepting APIM answers a DELETE with a 404, as for a dictionary deleted in the console.
var _ = It("releases a Dictionary APIM no longer holds", func() {
	dict := dictionaryOn(&refs.NamespacedName{Name: "accepting-ctx"})
	dict.SetFinalizers([]string{core.DictionaryFinalizer})
	dict.SetDeletionTimestamp(&metav1.Time{Time: time.Now()})
	cli := deletingCluster(dict)
	key := client.ObjectKeyFromObject(dict)

	_, err := (&dictionary.Reconciler{Client: cli, Recorder: record.NewFakeRecorder(20)}).
		Reconcile(context.Background(), ctrl.Request{NamespacedName: key})

	Expect(err).ToNot(HaveOccurred())
	Expect(apierrors.IsNotFound(cli.Get(context.Background(), key, &v1alpha1.Dictionary{}))).To(BeTrue())
})

var refusingCtx = &refs.NamespacedName{Name: "refusing-ctx"}

const automation = "/automation/organizations/DEFAULT/environments/DEFAULT"

// deletingCluster registers a fake client holding obj. An API or shared policy group delete first
// checks that nothing references it: those indexes hold nothing.
func deletingCluster(obj client.Object) client.Client {
	scheme := runtime.NewScheme()
	Expect(clientgoscheme.AddToScheme(scheme)).To(Succeed())
	Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
	nothing := func(client.Object) []string { return nil }
	cli := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(obj).
		WithStatusSubresource(obj).
		WithIndex(&v1alpha1.Documentation{}, search.DocumentationApiField.String(), nothing).
		WithIndex(&v1alpha1.PortalListing{}, search.PortalListingApiField.String(), nothing).
		WithIndex(&v1alpha1.PortalLink{}, search.PortalLinkApiField.String(), nothing).
		WithIndex(&v1alpha1.ApiV4Definition{}, search.ApiV4SharedPolicyGroupsField.String(), nothing).
		Build()
	k8s.RegisterClient(cli)
	return cli
}

func apiV4On(ctxRef *refs.NamespacedName) *v1alpha1.ApiV4Definition {
	api := &v1alpha1.ApiV4Definition{ObjectMeta: metav1.ObjectMeta{Name: "orders", Namespace: "gravitee"}}
	api.Spec.V4BaseApi = &v4.V4BaseApi{ApiBase: &base.ApiBase{}}
	api.Spec.Context = ctxRef
	return api
}

func apiV2On(ctxRef *refs.NamespacedName) *v1alpha1.ApiDefinition {
	api := &v1alpha1.ApiDefinition{ObjectMeta: metav1.ObjectMeta{Name: "orders-v2", Namespace: "gravitee"}}
	api.Spec.Api = v2.Api{ApiBase: &base.ApiBase{}}
	api.Spec.Context = ctxRef
	api.Status.ID = "8f0c"
	return api
}

func applicationOn(ctxRef *refs.NamespacedName) *v1alpha1.Application {
	app := &v1alpha1.Application{ObjectMeta: metav1.ObjectMeta{Name: "orders-app", Namespace: "gravitee"}}
	app.Spec.Context = ctxRef
	return app
}

func dictionaryOn(ctxRef *refs.NamespacedName) *v1alpha1.Dictionary {
	dict := &v1alpha1.Dictionary{ObjectMeta: metav1.ObjectMeta{Name: "orders-dict", Namespace: "gravitee"}}
	dict.Spec.Context = ctxRef
	return dict
}

func sharedPolicyGroupOn(ctxRef *refs.NamespacedName) *v1alpha1.SharedPolicyGroup {
	spg := &v1alpha1.SharedPolicyGroup{ObjectMeta: metav1.ObjectMeta{Name: "rate-limits", Namespace: "gravitee"}}
	spg.Spec.SharedPolicyGroup = &spgmodel.SharedPolicyGroup{}
	spg.Spec.Context = ctxRef
	return spg
}

func authzEntityOn(ctxRef *refs.NamespacedName) *v1alpha1.AuthzEntity {
	entity := newAuthzEntity()
	entity.Spec.Context = ctxRef
	return entity
}

func authzPdpOn(ctxRef *refs.NamespacedName) *v1alpha1.AuthzPdp {
	pdp := newAuthzPdp()
	pdp.Spec.Context = ctxRef
	return pdp
}

func authzPolicyOn(ctxRef *refs.NamespacedName) *v1alpha1.AuthzPolicy {
	policy := newAuthzPolicy()
	policy.Spec.Context = ctxRef
	return policy
}

func authzSchemaOn(ctxRef *refs.NamespacedName) *v1alpha1.AuthzSchema {
	schema := newAuthzSchema()
	schema.Spec.Context = ctxRef
	return schema
}
