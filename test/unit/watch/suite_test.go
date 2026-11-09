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

package watch_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
)

func TestWatch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Watch unit tests suite")
}

// recordingCache keeps the index functions InitCache registers, for the fake client to use.
type recordingCache struct {
	cache.Cache
	indexers map[string]client.IndexerFunc
}

func (c *recordingCache) IndexField(_ context.Context, _ client.Object, field string, f client.IndexerFunc) error {
	c.indexers[field] = f
	return nil
}

// cluster registers a fake client holding objects, AuthzEntity indexed as the manager's cache
// indexes it, and returns it.
func cluster(objects ...client.Object) client.Client {
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
	cli := builder.Build()
	k8s.RegisterClient(cli)
	return cli
}

func newQueue() workqueue.TypedRateLimitingInterface[reconcile.Request] {
	q := workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[reconcile.Request]())
	DeferCleanup(q.ShutDown)
	return q
}

// queued drains q and returns what it held.
func queued(q workqueue.TypedRateLimitingInterface[reconcile.Request]) []reconcile.Request {
	var requests []reconcile.Request
	for q.Len() > 0 {
		request, _ := q.Get()
		q.Done(request)
		requests = append(requests, request)
	}
	return requests
}
