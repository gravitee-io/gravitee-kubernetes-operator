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

	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
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
		Expect(cli.Create(ctx, o.Context)).ToNot(HaveOccurred())
		Eventually(func() error {
			err := manager.GetLatest(o.Context)
			if err != nil {
				return err
			}
			return assert.HasFinalizer(o.Context, keys.ManagementContextFinalizer)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), o.Context.Name)
	}

	if o.Resource != nil {
		Expect(cli.Create(ctx, o.Resource)).ToNot(HaveOccurred())
		Eventually(func() error {
			err := manager.GetLatest(o.Resource)
			if err != nil {
				return err
			}
			return assert.HasFinalizer(o.Resource, keys.ApiResourceFinalizer)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), o.Resource.Name)
	}

	if o.Application != nil {
		Expect(cli.Create(ctx, o.Application)).ToNot(HaveOccurred())
		Eventually(func() error {
			err := manager.GetLatest(o.Application)
			if err != nil {
				return err
			}
			if err = assert.ApplicationCompleted(o.Application); err != nil {
				return assert.ApplicationFailed(o.Application)
			}
			return nil
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), o.Application.Name)
	}

	if o.API != nil {
		Expect(cli.Create(ctx, o.API)).ToNot(HaveOccurred())
		Eventually(func() error {
			err := manager.GetLatest(o.API)
			if err != nil {
				return err
			}
			if isTemplate(o.API) {
				return assert.HasFinalizer(o.API, keys.ApiDefinitionTemplateFinalizer)
			}
			if err = assert.ApiCompleted(o.API); err != nil {
				return assert.ApiFailed(o.API)
			}
			return nil
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), o.API.Name)
	}

	if o.APIv4 != nil {
		Expect(cli.Create(ctx, o.APIv4)).ToNot(HaveOccurred())
		Eventually(func() error {
			err := manager.GetLatest(o.APIv4)
			if err != nil {
				return err
			}
			if isTemplate(o.APIv4) {
				return assert.HasFinalizer(o.APIv4, keys.ApiDefinitionTemplateFinalizer)
			}
			if err = assert.ApiV4Completed(o.APIv4); err != nil {
				return assert.ApiV4Failed(o.APIv4)
			}
			return nil
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), o.APIv4.Name)
	}

	if o.Ingress != nil {
		Expect(cli.Create(ctx, o.Ingress)).ToNot(HaveOccurred())
		Eventually(func() error {
			err := manager.GetLatest(o.Ingress)
			if err != nil {
				return err
			}
			return assert.HasFinalizer(o.Ingress, keys.IngressFinalizer)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), o.Ingress.Name)
	}

	return o
}
