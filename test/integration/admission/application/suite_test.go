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

package application

import (
	"testing"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	//+kubebuilder:scaffold:imports
)

func TestResources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Application Admission Suite")
}

var _ = SynchronizedBeforeSuite(func() {
	// NOSONAR mandatory noop
}, func() {
	fixture.Builder().
		AddSecret(constants.ContextSecretFile).
		AddSecret(constants.ApiWithTemplatingSecretFile).
		AddConfigMap(constants.ApiWithTemplatingConfigMapFile).
		Build().
		Apply()
})

var _ = SynchronizedAfterSuite(func() {
	By("Tearing down the test environment")
	gexec.KillAndWait(5 * time.Second)
}, func() {
	// NOSONAR ignore this noop func
})
