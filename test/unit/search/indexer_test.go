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
package search_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
)

// recordingCache keeps the index functions InitCache registers. Only IndexField is used.
type recordingCache struct {
	cache.Cache
	indexers map[string]client.IndexerFunc
}

func (c *recordingCache) IndexField(_ context.Context, _ client.Object, field string, f client.IndexerFunc) error {
	c.indexers[field] = f
	return nil
}

var _ = Describe("InitCache", func() {
	It("indexes an AMSecurityDomain by its AMContext", func() {
		c := &recordingCache{indexers: map[string]client.IndexerFunc{}}
		Expect(search.InitCache(context.Background(), c)).To(Succeed())

		domain := &v1alpha1.AMSecurityDomain{
			ObjectMeta: metav1.ObjectMeta{Name: "domain", Namespace: "ns"},
			Spec:       v1alpha1.AMSecurityDomainSpec{Context: &refs.NamespacedName{Name: "am-ctx"}},
		}

		index := c.indexers[search.AMSecurityContextField.String()]
		Expect(index).ToNot(BeNil())
		Expect(index(domain)).To(ConsistOf("ns/am-ctx"))
	})
})
