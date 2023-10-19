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
	"context"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	admissionRegistration "k8s.io/api/admissionregistration/v1"
	apiExtensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var log = ctrl.Log.WithName("extensions")

var (
	conversionPath = "/convert"
)

var (
	validationResourceName = "gko-validating-webhook-configurations"
	validationVersion      = "v1"
	validationSideEffects  = admissionRegistration.SideEffectClassNone
	validationPolicy       = admissionRegistration.Fail
	validatedOperations    = []admissionRegistration.OperationType{"CREATE", "UPDATE", "DELETE"}
)

var conversions = map[string][]string{
	"apidefinitions.gravitee.io": {"v1alpha1", "v1beta1"},
}

var admissions = map[string][]string{
	"apidefinitions.gravitee.io": {"v1alpha1", "v1beta1"},
}

func AddConversionWebhook(crd *apiExtensions.CustomResourceDefinition) {
	conversions := getConversions(crd.Name)
	if conversions != nil {
		setConversion(crd, conversions)
	}
}

func InjectCA(crd *apiExtensions.CustomResourceDefinition) {
	infectFrom := types.NamespacedName{
		Namespace: env.Config.WebhookNS,
		Name:      env.Config.WebhookCertSecret,
	}
	crd.Annotations[keys.InjectCAAnnotation] = infectFrom.String()
}

func CreateValidatingWebhooks(
	ctx context.Context,
	cli client.Client,
	crd *apiExtensions.CustomResourceDefinition,
) error {
	config := GetValidationWebhookConfig()
	if err := cli.Get(ctx, types.NamespacedName{Name: config.Name}, config); err != nil {
		return err
	}
	return cli.Update(ctx, AddValidatingWebhook(config, crd))
}

func InitValidatingWebhookConfig(
	ctx context.Context,
	cli client.Client,
) error {
	config := GetValidationWebhookConfig()
	err := cli.Get(ctx, types.NamespacedName{Name: config.Name}, config)

	if errors.IsNotFound(err) {
		log.Info("creating validating webhook", "name", config.Name)
		if err = cli.Create(ctx, config); err != nil {
			log.Error(err, "unable to create validating webhook", "name", config.Name)
			return err
		}
	} else if err != nil {
		log.Error(err, "unable to get validating webhook", "name", config.Name)
		return err
	}
	config.Webhooks = []admissionRegistration.ValidatingWebhook{}
	return cli.Update(ctx, config)
}

func getConversions(name string) []string {
	versions, ok := conversions[name]
	if !ok {
		return nil
	}
	return versions
}

func GetValidationWebhookConfig() *admissionRegistration.ValidatingWebhookConfiguration {
	log.Info("Creating validation webhook configuration", "name", validationResourceName)
	return &admissionRegistration.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: validationResourceName,
			Annotations: map[string]string{
				keys.InjectCAAnnotation: types.NamespacedName{
					Namespace: env.Config.WebhookNS,
					Name:      env.Config.WebhookCertSecret,
				}.String(),
			},
		},
	}
}

func AddValidatingWebhook(
	config *admissionRegistration.ValidatingWebhookConfiguration,
	crd *apiExtensions.CustomResourceDefinition,
) *admissionRegistration.ValidatingWebhookConfiguration {
	versions, ok := admissions[crd.Name]
	if !ok {
		log.Info("No validation webhook to add for resource", "resource", crd.Spec.Names.Singular)
		return config
	}

	for _, version := range versions {
		log.Info("Adding validation webhook configuration", "name", validationResourceName)
		config.Webhooks = append(config.Webhooks, buildValidatingWebhook(crd, version))
	}

	return config
}

func buildValidatingWebhook(
	crd *apiExtensions.CustomResourceDefinition, version string,
) admissionRegistration.ValidatingWebhook {
	return admissionRegistration.ValidatingWebhook{
		AdmissionReviewVersions: []string{validationVersion},
		SideEffects:             &validationSideEffects,
		Name:                    buildValidatingWebhookName(crd, version),
		ClientConfig: admissionRegistration.WebhookClientConfig{
			Service: &admissionRegistration.ServiceReference{
				Name:      env.Config.WebhookService,
				Namespace: env.Config.WebhookNS,
				Path:      buildValidatingWebhookPath(crd, version),
			},
		},
		FailurePolicy: &validationPolicy,
		Rules: []admissionRegistration.RuleWithOperations{
			{
				Operations: validatedOperations,
				Rule: admissionRegistration.Rule{
					APIGroups:   []string{keys.CrdGroup},
					APIVersions: []string{version},
					Resources:   []string{crd.Spec.Names.Plural}},
			},
		},
	}
}

func buildValidatingWebhookName(
	crd *apiExtensions.CustomResourceDefinition,
	version string,
) string {
	return version + "." + crd.Spec.Group + "." + crd.Spec.Names.Singular
}

func buildValidatingWebhookPath(
	crd *apiExtensions.CustomResourceDefinition,
	version string,
) *string {
	group := strings.ReplaceAll(crd.Spec.Group, ".", "-")
	name := crd.Spec.Names.Singular
	path := "/validate-" + group + "-" + version + "-" + name
	return &path
}

func setConversion(crd *apiExtensions.CustomResourceDefinition, versions []string) {
	crd.Spec.Conversion = &apiExtensions.CustomResourceConversion{
		Strategy: apiExtensions.WebhookConverter,
		Webhook: &apiExtensions.WebhookConversion{
			ClientConfig: &apiExtensions.WebhookClientConfig{
				Service: &apiExtensions.ServiceReference{
					Name:      env.Config.WebhookService,
					Namespace: env.Config.WebhookNS,
					Path:      &conversionPath,
				},
			},
			ConversionReviewVersions: versions,
		},
	}
}
