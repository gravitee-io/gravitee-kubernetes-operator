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

package fixture

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (o *Objects) Apply() *Objects {
	cli := manager.Client()
	ctx := context.Background()

	for _, sec := range o.Secrets {
		err := cli.Create(ctx, sec)
		Expect(client.IgnoreAlreadyExists(err)).ToNot(HaveOccurred(), sec.Name)
	}

	for _, cm := range o.ConfigMaps {
		err := cli.Create(ctx, cm)
		Expect(client.IgnoreAlreadyExists(err)).ToNot(HaveOccurred(), cm.Name)
	}

	if o.Context != nil {
		o.applyContext(cli, ctx)
	}

	if o.Resource != nil {
		o.applyResource(cli, ctx)
	}

	if o.Application != nil {
		o.applyApplication(cli, ctx)
	}

	if o.API != nil {
		o.applyAPI(cli, ctx)
	}

	if o.APIv4 != nil {
		o.applyAPIv4(cli, ctx)
	}

	if o.Subscription != nil {
		o.applySubscription(cli, ctx)
	}

	if o.Ingress != nil {
		o.applyIngress(cli, ctx)
	}

	return o
}

func (o *Objects) applyIngress(cli client.Client, ctx context.Context) {
	Expect(cli.Create(ctx, o.Ingress)).ToNot(HaveOccurred())
	Eventually(ctx, func() error {
		err := manager.GetLatest(ctx, o.Ingress)
		if err != nil {
			return err
		}
		return assert.HasFinalizer(o.Ingress, core.IngressFinalizer)
	}, constants.EventualTimeout, constants.Interval).Should(Succeed(), o.Ingress.Name)
}

func (o *Objects) applyAPIv4(cli client.Client, ctx context.Context) {
	Expect(cli.Create(ctx, o.APIv4)).ToNot(HaveOccurred())
	Eventually(ctx, func() error {
		err := manager.GetLatest(ctx, o.APIv4)
		if err != nil {
			return err
		}
		if isTemplate(o.APIv4) {
			return assert.HasFinalizer(o.APIv4, core.ApiDefinitionTemplateFinalizer)
		}
		if err = assert.ApiV4Completed(o.APIv4); err != nil {
			return assert.ApiV4Failed(o.APIv4)
		}
		return nil
	}, constants.EventualTimeout, constants.Interval).Should(Succeed(), o.APIv4.Name)
}

func (o *Objects) applyAPI(cli client.Client, ctx context.Context) {
	Expect(cli.Create(ctx, o.API)).ToNot(HaveOccurred())
	Eventually(ctx, func() error {
		err := manager.GetLatest(ctx, o.API)
		if err != nil {
			return err
		}
		if isTemplate(o.API) {
			return assert.HasFinalizer(o.API, core.ApiDefinitionTemplateFinalizer)
		}
		if err = assert.ApiCompleted(o.API); err != nil {
			return assert.ApiFailed(o.API)
		}
		return nil
	}, constants.EventualTimeout, constants.Interval).Should(Succeed(), o.API.Name)
}

func (o *Objects) applyApplication(cli client.Client, ctx context.Context) {
	Expect(cli.Create(ctx, o.Application)).ToNot(HaveOccurred())
	Eventually(ctx, func() error {
		err := manager.GetLatest(ctx, o.Application)
		if err != nil {
			return err
		}
		if err = assert.ApplicationCompleted(o.Application); err != nil {
			return assert.ApplicationFailed(o.Application)
		}
		return nil
	}, constants.EventualTimeout, constants.Interval).Should(Succeed(), o.Application.Name)
}

func (o *Objects) applyResource(cli client.Client, ctx context.Context) {
	Expect(cli.Create(ctx, o.Resource)).ToNot(HaveOccurred())
	Eventually(ctx, func() error {
		err := manager.GetLatest(ctx, o.Resource)
		if err != nil {
			return err
		}
		return assert.HasFinalizer(o.Resource, core.ApiResourceFinalizer)
	}, constants.EventualTimeout, constants.Interval).Should(Succeed(), o.Resource.Name)
}

func (o *Objects) applyContext(cli client.Client, ctx context.Context) {
	Expect(cli.Create(ctx, o.Context)).ToNot(HaveOccurred())
	Eventually(ctx, func() error {
		err := manager.GetLatest(ctx, o.Context)
		if err != nil {
			return err
		}
		return assert.HasFinalizer(o.Context, core.ManagementContextFinalizer)
	}, constants.EventualTimeout, constants.Interval).Should(Succeed(), o.Context.Name)
}

func (o *Objects) applySubscription(cli client.Client, ctx context.Context) {
	Expect(cli.Create(ctx, o.Subscription)).ToNot(HaveOccurred())
	Eventually(ctx, func() error {
		err := manager.GetLatest(ctx, o.Subscription)
		if err != nil {
			return err
		}
		if err = assert.SubscriptionCompleted(o.Subscription); err != nil {
			return assert.SubscriptionFailed(o.Subscription)
		}
		return nil
	}, constants.EventualTimeout, constants.Interval).Should(Succeed(), o.Subscription.Name)
}
