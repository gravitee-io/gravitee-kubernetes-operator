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

package env

import (
	"os"
	"strconv"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
)

const (
	CMTemplate404Name      = "TEMPLATE_404_CONFIG_MAP_NAME"
	CMTemplate404NS        = "TEMPLATE_404_CONFIG_MAP_NAMESPACE"
	Development            = "DEV_MODE"
	NS                     = "NAMESPACE"
	ApplyCRDs              = "APPLY_CRDS"
	EnableMetrics          = "ENABLE_METRICS"
	EnableWebhook          = "ENABLE_WEBHOOK"
	WebhookNS              = "WEBHOOK_NAMESPACE"
	WebhookServiceName     = "WEBHOOK_SERVICE_NAME"
	WebhookPort            = "WEBHOOK_SERVICE_PORT"
	WebhookCertSecret      = "WEBHOOK_CERT_SECRET_NAME" //nolint:gosec // This is not a hardcoded secret
	InsecureSkipCertVerify = "INSECURE_SKIP_CERT_VERIFY"
	TrueString             = "true"
	IngressClasses         = "INGRESS_CLASSES"

	// This default are applied when running the app locally.
	defaultWebhookPort = 9443
)

var Config = struct {
	NS                 string
	ApplyCRDs          bool
	EnableMetrics      bool
	EnableWebhook      bool
	WebhookNS          string
	WebhookService     string
	WebhookPort        int
	WebhookCertSecret  string
	Development        bool
	CMTemplate404Name  string
	CMTemplate404NS    string
	InsecureSkipVerify bool
	IngressClasses     []string
}{}

func init() {
	Config.NS = os.Getenv(NS)
	Config.ApplyCRDs = os.Getenv(ApplyCRDs) == TrueString
	Config.Development = os.Getenv(Development) == TrueString
	Config.CMTemplate404Name = os.Getenv(CMTemplate404Name)
	Config.CMTemplate404NS = os.Getenv(CMTemplate404NS)
	Config.InsecureSkipVerify = os.Getenv(InsecureSkipCertVerify) == TrueString
	Config.EnableMetrics = os.Getenv(EnableMetrics) == TrueString
	Config.EnableWebhook = os.Getenv(EnableWebhook) == TrueString
	Config.WebhookNS = os.Getenv(WebhookNS)
	Config.WebhookService = os.Getenv(WebhookServiceName)
	Config.WebhookCertSecret = os.Getenv(WebhookCertSecret)
	Config.WebhookPort = parseInt(WebhookPort, defaultWebhookPort)
	var ingressClass string
	if ingressClass = keys.IngressClassAnnotationValue; os.Getenv(IngressClasses) != "" {
		ingressClass = os.Getenv(IngressClasses)
	}
	Config.IngressClasses = strings.Split(ingressClass, ",")
}

func parseInt(key string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return defaultValue
	}
	return value
}
