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
	"sync"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1beta1"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var cli client.Client
var once sync.Once

func ClusterClient() client.Client {
	once.Do(func() {
		c, err := client.New(ctrl.GetConfigOrDie(), client.Options{Scheme: scheme.Scheme})
		Expect(err).ToNot(HaveOccurred())
		cli = c
	})
	return cli
}

func init() {
	must(scheme.AddToScheme(scheme.Scheme))
	must(v1alpha1.AddToScheme(scheme.Scheme))
	must(v1beta1.AddToScheme(scheme.Scheme))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
