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

package extensions

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	apiExtensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/types"
)

var versionsConversions = map[string][]string{
	"apidefinitions.gravitee.io": {"v1alpha1", "v1beta1"},
}

var webhookContextPath = "/convert"

func ExtendWithWebhook(crd *apiExtensions.CustomResourceDefinition) {
	versions := getVersions(crd.Name)
	if versions != nil {
		setConversion(crd, versions)
	}
}

func InjectCA(crd *apiExtensions.CustomResourceDefinition) {
	infectFrom := types.NamespacedName{
		Namespace: env.Config.WebhookNS,
		Name:      env.Config.WebhookCertSecret,
	}
	crd.Annotations[keys.InjectCAAnnotation] = infectFrom.String()
}

func getVersions(name string) []string {
	versions, ok := versionsConversions[name]
	if !ok {
		return nil
	}
	return versions
}

func setConversion(crd *apiExtensions.CustomResourceDefinition, versions []string) {
	crd.Spec.Conversion = &apiExtensions.CustomResourceConversion{
		Strategy: apiExtensions.WebhookConverter,
		Webhook: &apiExtensions.WebhookConversion{
			ClientConfig: &apiExtensions.WebhookClientConfig{
				Service: &apiExtensions.ServiceReference{
					Name:      env.Config.WebhookService,
					Namespace: env.Config.WebhookNS,
					Path:      &webhookContextPath,
				},
			},
			ConversionReviewVersions: versions,
		},
	}
}
