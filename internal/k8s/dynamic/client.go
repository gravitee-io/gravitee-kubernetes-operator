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

package dynamic

import (
	"sync"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"

	ctrl "sigs.k8s.io/controller-runtime"
)

var dynamicClient *dynamic.DynamicClient
var once sync.Once

// UseConfig builds the client from cfg instead of the ambient kubeconfig.
// It must run before the first GetClient call to take effect.
// FIXME: only envtest suites need this; AM should stop resolving secrets
// through this package, then this seam can go.
func UseConfig(cfg *rest.Config) {
	once.Do(func() {
		dynamicClient = dynamic.NewForConfigOrDie(cfg)
	})
}

func GetClient() *dynamic.DynamicClient {
	once.Do(func() {
		dynamicClient = dynamic.NewForConfigOrDie(ctrl.GetConfigOrDie())
	})
	return dynamicClient
}
