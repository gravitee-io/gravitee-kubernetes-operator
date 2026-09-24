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

package envtest

import (
	"os"
	"path/filepath"
	"runtime"

	runtimeUtil "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

// Start boots kube-apiserver and etcd with the gravitee.io CRDs installed.
func Start() (*envtest.Environment, *rest.Config) {
	if os.Getenv("KUBEBUILDER_ASSETS") == "" {
		panic("KUBEBUILDER_ASSETS is not set: run make envtest")
	}

	env := &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join(repoRoot(), "crds", "gravitee.io")},
		ErrorIfCRDPathMissing: true,
	}

	cfg, err := env.Start()
	runtimeUtil.Must(err)

	return env, cfg
}

func Stop(env *envtest.Environment) {
	runtimeUtil.Must(env.Stop())
}

func repoRoot() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..")
}
