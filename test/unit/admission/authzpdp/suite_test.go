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

package authzpdp_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestAuthzPdpAdmission(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AuthzPdp admission hook tests suite")
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
