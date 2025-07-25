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

package core

const (
	CRDGroup   = "gravitee.io"
	CRDVersion = "v1alpha1"

	CRDManagementContextResource          = "managementcontexts"
	CRDApplicationResource                = "applications"
	CRDApiDefinitionResource              = "apidefinitions"
	CRDApiV4DefinitionResource            = "apiv4definitions"
	CRDResourceResource                   = "apiresources"
	CRDResourceSharedPolicyGroupsResource = "sharedpolicygroups"
	CRDResourceNotification               = "notifications"
	CRDResourceGroup                      = "groups"

	GraviteeComponentLabel      = "gravitee.io/component"
	IngressLabel                = "gravitee.io/ingress"
	IngressLabelValue           = "graviteeio"
	IngressClassAnnotation      = "kubernetes.io/ingress.class"
	IngressClassAnnotationValue = "graviteeio"
	IngressTemplateAnnotation   = "gravitee.io/template"
	GraviteePemRegistryLabel    = "kubernetes-pem-registry"
	LastSpecHashAnnotation      = "gravitee.io/last-spec-hash"

	Extends = "gravitee.io/extends"

	NotificationFinalizer = "finalizers.gravitee.io/notification"
	//nolint:gosec // This is not an hardcoded secret
	ManagementContextSecretFinalizer = "finalizers.gravitee.io/managementcontextSecret"
	ApiDefinitionFinalizer           = "finalizers.gravitee.io/apidefinitiondeletion"
	ApiDefinitionTemplateFinalizer   = "finalizers.gravitee.io/apidefinitiontemplate"
	ManagementContextFinalizer       = "finalizers.gravitee.io/managementcontextdeletion"
	ApiResourceFinalizer             = "finalizers.gravitee.io/apiresource"
	IngressFinalizer                 = "finalizers.gravitee.io/ingress"
	KeyPairFinalizer                 = "finalizers.gravitee.io/keypair"
	ApplicationFinalizer             = "finalizers.gravitee.io/applicationdeletion"
	SubscriptionFinalizer            = "finalizers.gravitee.io/subscriptions"
	GroupFinalizer                   = "finalizers.gravitee.io/groups"
	TemplatingFinalizer              = "finalizers.gravitee.io/templating"
	SharedPolicyGroupFinalizer       = "finalizers.gravitee.io/sharedpolicygroups"
	GraviteeClassParametersFinalizer = "finalizers.gravitee.io/gatewayclassparameters"
	GatewayClassFinalizer            = "finalizers.gravitee.io/gatewayclass"
	GatewayFinalizer                 = "finalizers.gravitee.io/gateway"
	HTTPRouteFinalizer               = "finalizers.gravitee.io/httproute"
	KafkaRouteFinalizer              = "finalizers.gravitee.io/kafkaroute"

	GraviteeGatewayClassController = "apim.gravitee.io/gateway"

	CloudTokenSecretKey  = "cloudToken"
	BearerTokenSecretKey = "bearerToken"
	UsernameSecretKey    = "username"
	PasswordSecretKey    = "password"
)
