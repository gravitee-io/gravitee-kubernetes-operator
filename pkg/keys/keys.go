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

package keys

// Kubernetes Ingresses.
const (
	GraviteeComponentLabel      = "gravitee.io/component"
	IngressLabel                = "gravitee.io/ingress"
	IngressLabelValue           = "graviteeio"
	IngressClassAnnotation      = "kubernetes.io/ingress.class"
	IngressClassAnnotationValue = "graviteeio"
	IngressTemplateAnnotation   = "gravitee.io/template"
	GraviteePemRegistryLabel    = "kubernetes-pem-registry"
)

// Gravitee.io CRDs.
const (
	CrdGroup   = "gravitee.io"
	CrdVersion = "v1alpha1"

	CrdManagementContextResource = "managementcontext"
	CrdApiDefinitionResource     = "apidefinitions"
)

const Extends = "gravitee.io/extends"

// Kubernetes Finalizers.
const (
	ApiDefinitionDeletionFinalizer = "finalizers.gravitee.io/apidefinitiondeletion"
	ApiDefinitionTemplateFinalizer = "finalizers.gravitee.io/apidefinitiontemplate"
	ManagementContextFinalizer     = "finalizers.gravitee.io/managementcontextdeletion"
	ApiResourceFinalizer           = "finalizers.gravitee.io/apiresource"
	//nolint:gosec // This is not an hardcoded secret
	ManagementContextSecretFinalizer = "finalizers.gravitee.io/managementcontextSecret"
	IngressFinalizer                 = "finalizers.gravitee.io/ingress"
	KeyPairFinalizer                 = "finalizers.gravitee.io/keypair"
	ApplicationDeletionFinalizer     = "finalizers.gravitee.io/applicationdeletion"
	TemplatingFinalizer              = "finalizers.gravitee.io/templating"
)
