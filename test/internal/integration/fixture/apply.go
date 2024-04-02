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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	. "github.com/onsi/gomega"
	netV1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
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
			cp := new(v1alpha1.ManagementContext)
			key := types.NamespacedName{
				Namespace: o.Context.GetNamespace(),
				Name:      o.Context.GetName(),
			}
			if err := cli.Get(ctx, key, cp); err != nil {
				return err
			}
			o.Context.Status = cp.Status
			o.Context.ObjectMeta = cp.ObjectMeta
			return nil
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), o.Context.Name)
	}

	if o.Resource != nil {
		Expect(cli.Create(ctx, o.Resource)).ToNot(HaveOccurred())
		Eventually(func() error {
			cp := new(v1alpha1.ApiResource)
			key := types.NamespacedName{
				Namespace: o.Resource.GetNamespace(),
				Name:      o.Resource.GetName(),
			}
			if err := cli.Get(ctx, key, cp); err != nil {
				return err
			}
			o.Resource.Status = cp.Status
			o.Resource.ObjectMeta = cp.ObjectMeta
			return nil
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), o.Resource.Name)
	}

	if o.Application != nil {
		Expect(cli.Create(ctx, o.Application)).ToNot(HaveOccurred())
		Eventually(func() error {
			cp := new(v1alpha1.Application)
			key := types.NamespacedName{
				Namespace: o.Application.GetNamespace(),
				Name:      o.Application.GetName(),
			}
			if err := cli.Get(ctx, key, cp); err != nil {
				return err
			}
			if err := assert.NotEmptyString("status", string(cp.Status.Status)); err != nil {
				return err
			}
			o.Application.Status = cp.Status
			o.Application.ObjectMeta = cp.ObjectMeta
			return nil
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), o.Application.Name)
	}

	if o.API != nil {
		Expect(cli.Create(ctx, o.API)).ToNot(HaveOccurred())
		Eventually(func() error {
			cp := new(v1alpha1.ApiDefinition)
			key := types.NamespacedName{
				Namespace: o.API.Namespace,
				Name:      o.API.Name,
			}
			if err := cli.Get(ctx, key, cp); err != nil {
				return err
			}
			if err := assert.NotEmptyString("status", string(cp.Status.Status)); err != nil {
				return err
			}
			o.API.Status = cp.Status
			o.API.ObjectMeta = cp.ObjectMeta
			return nil
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), o.API.Name)
	}

	if o.Ingress != nil {
		Expect(cli.Create(ctx, o.Ingress)).ToNot(HaveOccurred())
		Eventually(func() error {
			cp := new(netV1.Ingress)
			key := types.NamespacedName{
				Namespace: o.Ingress.GetNamespace(),
				Name:      o.Ingress.GetName(),
			}
			if err := cli.Get(ctx, key, cp); err != nil {
				return err
			}
			o.Ingress.Status = cp.Status
			o.Ingress.ObjectMeta = cp.ObjectMeta
			return nil
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), o.Ingress.Name)
	}

	return o
}
