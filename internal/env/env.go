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
)

const (
	CMTemplate404Name      = "TEMPLATE_404_CONFIG_MAP_NAME"
	CMTemplate404NS        = "TEMPLATE_404_CONFIG_MAP_NAMESPACE"
	DisableJSONLogs        = "DISABLE_JSON_LOGS"
	LogFormat              = "LOG_FORMAT"
	LogLevel               = "LOG_LEVEL"
	LogTimeFormat          = "LOG_TIME_FORMAT"
	LogReconcileIdField    = "LOG_RECONCILE_ID_FIELD"
	WatchNS                = "WATCH_NAMESPACE"
	ApplyCRDs              = "APPLY_CRDS"
	EnableMetrics          = "ENABLE_METRICS"
	InsecureSkipCertVerify = "INSECURE_SKIP_CERT_VERIFY"
	EnableWebhook          = "ENABLE_WEBHOOK"
	WebhookNS              = "WEBHOOK_NAMESPACE"
	WebhookServiceName     = "WEBHOOK_SERVICE_NAME"
	WebhookPort            = "WEBHOOK_SERVICE_PORT"
	WebhookCertSecret      = "WEBHOOK_CERT_SECRET_NAME" //nolint:gosec // This is not an hardcoded secret

	trueString             = "true"
	defaultWebhookPort     = 9443
	defaultMetricsPort     = 8080
	defaultProbesPort      = 8080
	defaultLogFormat       = "json"
	defaultLogLevel        = "info"
	defaultLogTimeFormat   = "EpochMillis"
	defaultLogTraceIdField = "reconcile-id"
)

var Config = struct {
	WatchNS            string
	ReleaseNS          string
	ApplyCRDs          bool
	EnableMetrics      bool
	DisableJSONLogs    bool
	LogFormat          string
	LogLevel           string
	LogTimeFormat      string
	LogTraceIdField    string
	CMTemplate404Name  string
	CMTemplate404NS    string
	InsecureSkipVerify bool
	EnableWebhook      bool
	WebhookNS          string
	WebhookService     string
	WebhookPort        int
	WebhookCertSecret  string
}{}

func init() {
	Config.WatchNS = os.Getenv(WatchNS)
	Config.ReleaseNS = os.Getenv(WebhookNS)
	Config.ApplyCRDs = os.Getenv(ApplyCRDs) == trueString
	Config.DisableJSONLogs = os.Getenv(DisableJSONLogs) == trueString
	Config.LogFormat = getStringOrDefault(LogFormat, defaultLogFormat)
	Config.LogLevel = getStringOrDefault(LogLevel, defaultLogLevel)
	Config.LogTimeFormat = getStringOrDefault(LogTimeFormat, defaultLogTimeFormat)
	Config.LogTraceIdField = getStringOrDefault(LogReconcileIdField, defaultLogTraceIdField)
	Config.CMTemplate404Name = os.Getenv(CMTemplate404Name)
	Config.CMTemplate404NS = os.Getenv(CMTemplate404NS)
	Config.InsecureSkipVerify = os.Getenv(InsecureSkipCertVerify) == trueString
	Config.EnableMetrics = os.Getenv(EnableMetrics) == trueString
	Config.EnableWebhook = os.Getenv(EnableWebhook) == trueString
	Config.WebhookNS = os.Getenv(WebhookNS)
	Config.WebhookService = os.Getenv(WebhookServiceName)
	Config.WebhookCertSecret = os.Getenv(WebhookCertSecret)
	Config.WebhookPort = parseInt(WebhookPort, defaultWebhookPort)
}

func getStringOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func parseInt(key string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return defaultValue
	}
	return value
}
