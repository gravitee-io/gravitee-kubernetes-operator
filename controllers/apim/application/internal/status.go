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

package internal

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/kube/custom"
	"k8s.io/apimachinery/pkg/types"
)

func UpdateStatusSuccess(ctx context.Context, application *v1alpha1.Application) error {
	if application.IsBeingDeleted() {
		return nil
	}

	app := &v1alpha1.Application{}
	cli := k8s.GetClient()

	if err := cli.Get(
		ctx, types.NamespacedName{Namespace: application.Namespace, Name: application.Name}, app,
	); err != nil {
		return err
	}

	application.Status.ObservedGeneration = application.ObjectMeta.Generation
	application.Status.DeepCopyInto(&app.Status)
	return cli.Status().Update(ctx, app)
}

func UpdateStatusFailure(ctx context.Context, application *v1alpha1.Application) error {
	app := &v1alpha1.Application{}
	cli := k8s.GetClient()

	if err := cli.Get(
		ctx, types.NamespacedName{Namespace: application.Namespace, Name: application.Name}, app,
	); err != nil {
		return err
	}

	application.Status.Status = custom.ProcessingStatusFailed
	application.Status.DeepCopyInto(&app.Status)
	return cli.Status().Update(ctx, app)
}
