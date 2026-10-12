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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/notification"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apidefinition"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("ResolveConsoleNotificationRefs", func() {
	ctx := context.Background()

	registerNotification := func(events ...notification.ApiEvent) {
		scheme := runtime.NewScheme()
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
		notif := &v1alpha1.Notification{
			ObjectMeta: metav1.ObjectMeta{Name: "notification", Namespace: "default"},
			Spec: v1alpha1.NotificationSpec{Type: &notification.Type{
				Target:    notification.TargetConsole,
				EventType: notification.EventTypeAPI,
				Console:   notification.Console{APIEvents: events},
			}},
		}
		k8s.RegisterClient(fake.NewClientBuilder().WithScheme(scheme).WithObjects(notif).Build())
	}

	apiBase := func() *base.ApiBase {
		return &base.ApiBase{
			Name:              "api",
			Version:           "1.0",
			NotificationsRefs: []refs.NamespacedName{{Name: "notification"}},
		}
	}

	newV4API := func() *v1alpha1.ApiV4Definition {
		return &v1alpha1.ApiV4Definition{
			ObjectMeta: metav1.ObjectMeta{Name: "api", Namespace: "default"},
			Spec: v1alpha1.ApiV4DefinitionSpec{
				Api: v4.Api{V4BaseApi: &v4.V4BaseApi{ApiBase: apiBase()}},
			},
		}
	}

	It("sends no console notification for a v4 API when the notification has no events", func() {
		registerNotification()
		api := newV4API()

		Expect(apidefinition.ResolveConsoleNotificationRefs(ctx, api)).To(Succeed())

		Expect(api.Spec.ConsoleNotification).To(BeNil())
	})

	It("sends no console notification for a v2 API when the notification has no events", func() {
		registerNotification()
		api := &v1alpha1.ApiDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: "api", Namespace: "default"},
			Spec: v1alpha1.ApiDefinitionV2Spec{
				Api: v2.Api{ApiBase: apiBase()},
			},
		}

		Expect(apidefinition.ResolveConsoleNotificationRefs(ctx, api)).To(Succeed())

		Expect(api.Spec.ConsoleNotification).To(BeNil())
	})

	It("sends the notification events for a v4 API", func() {
		registerNotification("API_STARTED")
		api := newV4API()

		Expect(apidefinition.ResolveConsoleNotificationRefs(ctx, api)).To(Succeed())

		Expect(api.Spec.ConsoleNotification).NotTo(BeNil())
		Expect(api.Spec.ConsoleNotification.Events).To(Equal([]string{"API_STARTED"}))
	})
})
