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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Gravitee.io CRDs.

var ApplicationGVR = schema.GroupVersionResource{
	Group:    core.CRDGroup,
	Version:  core.CRDVersion,
	Resource: core.CRDApplicationResource,
}

var ApiGVR = schema.GroupVersionResource{
	Group:    core.CRDGroup,
	Version:  core.CRDVersion,
	Resource: core.CRDApiDefinitionResource,
}

var ApiV4GVR = schema.GroupVersionResource{
	Group:    core.CRDGroup,
	Version:  core.CRDVersion,
	Resource: core.CRDApiV4DefinitionResource,
}

var ManagementContextGVR = schema.GroupVersionResource{
	Group:    core.CRDGroup,
	Version:  core.CRDVersion,
	Resource: core.CRDManagementContextResource,
}

var ResourceGVR = schema.GroupVersionResource{
	Group:    core.CRDGroup,
	Version:  core.CRDVersion,
	Resource: core.CRDResourceResource,
}

var NotificationGVR = schema.GroupVersionResource{
	Group:    core.CRDGroup,
	Version:  core.CRDVersion,
	Resource: core.CRDResourceNotification,
}

var GroupGVR = schema.GroupVersionResource{
	Group:    core.CRDGroup,
	Version:  core.CRDVersion,
	Resource: core.CRDResourceGroup,
}

var SecretGVR = schema.GroupVersionResource{
	Group:    "",
	Version:  "v1",
	Resource: "secrets",
}
