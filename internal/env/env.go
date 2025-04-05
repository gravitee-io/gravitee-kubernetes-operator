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
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

const (
	CMTemplate404Name                    = "TEMPLATE_404_CONFIG_MAP_NAME"
	CMTemplate404NS                      = "TEMPLATE_404_CONFIG_MAP_NAMESPACE"
	Development                          = "DEV_MODE"
	NS                                   = "NAMESPACE"
	ApplyCRDs                            = "APPLY_CRDS"
	EnableMetrics                        = "ENABLE_METRICS"
	SecureMetrics                        = "SECURE_METRICS"
	MetricsCertDir                       = "METRICS_CERT_DIR"
	MetricsPort                          = "METRICS_PORT"
	ProbesPort                           = "PROBES_PORT"
	EnableIngress                        = "ENABLE_INGRESS"
	EnableWebhook                        = "ENABLE_WEBHOOK"
	WebhookNS                            = "WEBHOOK_NAMESPACE"
	WebhookServiceName                   = "WEBHOOK_SERVICE_NAME"
	WebhookPort                          = "WEBHOOK_SERVICE_PORT"
	WebhookCertSecret                    = "WEBHOOK_CERT_SECRET_NAME" //nolint:gosec // This is not a hardcoded secret
	HttpCLientInsecureSkipCertVerify     = "HTTP_CLIENT_INSECURE_SKIP_CERT_VERIFY"
	HttpClientTimeoutSeconds             = "HTTP_CLIENT_TIMEOUT_SECONDS"
	TrueString                           = "true"
	IngressClasses                       = "INGRESS_CLASSES"
	CheckApiContextPathConflictInCluster = "CHECK_API_CONTEXT_PATH_CONFLICT_IN_CLUSTER"
	LogsFormat                           = "LOGS_FORMAT"
	LogsLevel                            = "LOGS_LEVEL"
	LogsLevelCase                        = "LOGS_LEVEL_CASE"
	LogsTimestampField                   = "LOGS_TIMESTAMP_FIELD"
	LogsTimestampFormat                  = "LOGS_TIMESTAMP_FORMAT"

	// This default are applied when running the app locally.
	defaultWebhookPort       = 9443
	defaultMetricsPort       = 8080
	defaultProbesPort        = 8081
	defaultHttpClientTimeout = 5

	ReconcileStrategy = "RECONCILE_STRATEGY"
)

var Config = struct {
	NS                                   string
	ApplyCRDs                            bool
	EnableMetrics                        bool
	SecureMetrics                        bool
	MetricsCertDir                       string
	MetricsPort                          int
	ProbesPort                           int
	EnableIngress                        bool
	EnableWebhook                        bool
	WebhookNS                            string
	WebhookService                       string
	WebhookPort                          int
	WebhookCertSecret                    string
	Development                          bool
	CMTemplate404Name                    string
	CMTemplate404NS                      string
	HTTPClientInsecureSkipVerify         bool
	HTTPClientTimeoutSeconds             int
	IngressClasses                       []string
	CheckApiContextPathConflictInCluster bool
	LogsFormat                           string
	LogsLevel                            string
	LogsLevelCase                        string
	LogsTimestampField                   string
	LogsTimestampFormat                  string
	ReconcileStrategy                    string
}{}

func init() {
	Config.NS = os.Getenv(NS)
	Config.ApplyCRDs = os.Getenv(ApplyCRDs) == TrueString
	Config.Development = os.Getenv(Development) == TrueString
	Config.CMTemplate404Name = os.Getenv(CMTemplate404Name)
	Config.CMTemplate404NS = os.Getenv(CMTemplate404NS)
	Config.HTTPClientInsecureSkipVerify = os.Getenv(HttpCLientInsecureSkipCertVerify) == TrueString
	Config.HTTPClientTimeoutSeconds = parseInt(HttpClientTimeoutSeconds, defaultHttpClientTimeout)
	Config.EnableMetrics = os.Getenv(EnableMetrics) == TrueString
	Config.SecureMetrics = os.Getenv(SecureMetrics) == TrueString
	Config.MetricsCertDir = os.Getenv(MetricsCertDir)
	Config.MetricsPort = parseInt(MetricsPort, defaultMetricsPort)
	Config.ProbesPort = parseInt(ProbesPort, defaultProbesPort)
	Config.EnableIngress = os.Getenv(EnableIngress) == TrueString
	Config.EnableWebhook = os.Getenv(EnableWebhook) == TrueString
	Config.WebhookNS = os.Getenv(WebhookNS)
	Config.WebhookService = os.Getenv(WebhookServiceName)
	Config.WebhookCertSecret = os.Getenv(WebhookCertSecret)
	Config.WebhookPort = parseInt(WebhookPort, defaultWebhookPort)
	Config.CheckApiContextPathConflictInCluster = os.Getenv(CheckApiContextPathConflictInCluster) == TrueString
	var ingressClass string
	if ingressClass = core.IngressClassAnnotationValue; os.Getenv(IngressClasses) != "" {
		ingressClass = os.Getenv(IngressClasses)
	}
	Config.IngressClasses = strings.Split(ingressClass, ",")
	Config.LogsFormat = os.Getenv(LogsFormat)
	Config.LogsLevel = os.Getenv(LogsLevel)
	Config.LogsLevelCase = os.Getenv(LogsLevelCase)
	Config.LogsTimestampField = os.Getenv(LogsTimestampField)
	Config.LogsTimestampFormat = os.Getenv(LogsTimestampFormat)
	Config.ReconcileStrategy = os.Getenv(ReconcileStrategy)
}

func GetMetricsAddr() string {
	if !Config.EnableMetrics {
		return "0" // disables metrics
	}
	return fmt.Sprintf(":%d", Config.MetricsPort)
}

func GetProbesAddr() string {
	return fmt.Sprintf(":%d", Config.ProbesPort)
}

func parseInt(key string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return defaultValue
	}
	return value
}
