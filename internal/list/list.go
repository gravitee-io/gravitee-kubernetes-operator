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

package list

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/types/list"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func OfKind(gvk schema.GroupVersionKind) (list.Interface, error) {
	switch {
	case cmp.Equal(gvk, v1alpha1.ApiDefinitionKind):
		return &v1alpha1.ApiDefinitionList{}, nil
	case cmp.Equal(gvk, v1alpha1.ApiDefinitionListKind):
		return &v1alpha1.ApiDefinitionList{}, nil
	default:
		return nil, fmt.Errorf("unknown kind %s", gvk)
	}
}
