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

package apim_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portaltheme"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

func newPortalTheme(definition portaltheme.Definition) *v1alpha1.PortalTheme {
	return &v1alpha1.PortalTheme{
		ObjectMeta: metav1.ObjectMeta{Name: "default-theme", Namespace: "default"},
		Spec: v1alpha1.PortalThemeSpec{
			Type: portaltheme.Type{
				Name:       "Default Theme",
				Definition: definition,
			},
		},
	}
}

var _ = Describe("Portal theme spec", func() {
	It("should be reflected in the spec hash", func() {
		primary := newPortalTheme(portaltheme.Definition{
			Color: &portaltheme.Color{Primary: new("#123456")},
		})
		secondary := newPortalTheme(portaltheme.Definition{
			Color: &portaltheme.Color{Primary: new("#654321")},
		})
		empty := newPortalTheme(portaltheme.Definition{})

		Expect(primary.Spec.Hash()).ToNot(Equal(secondary.Spec.Hash()))
		Expect(empty.Spec.Hash()).ToNot(Equal(primary.Spec.Hash()))
	})
})
