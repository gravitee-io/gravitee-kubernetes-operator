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
<<<<<<< HEAD
	"strconv"
)

const (
	CMTemplate404Name        = "TEMPLATE_404_CONFIG_MAP_NAME"
	CMTemplate404NS          = "TEMPLATE_404_CONFIG_MAP_NAMESPACE"
	DisableJSONLogs          = "DISABLE_JSON_LOGS"
	LogFormat                = "LOG_FORMAT"
	LogLevel                 = "LOG_LEVEL"
	LogTimestampField        = "LOG_TIMESTAMP_FIELD"
	LogTimestampFormat       = "LOG_TIMESTAMP_FORMAT"
	WatchNS                  = "WATCH_NAMESPACE"
	ApplyCRDs                = "APPLY_CRDS"
	EnableMetrics            = "ENABLE_METRICS"
	SecureMetrics            = "SECURE_METRICS"
	MetricsCertDir           = "METRICS_CERT_DIR"
	MetricsPort              = "METRICS_PORT"
	ProbePort                = "PROBE_PORT"
	HTTPClientSkipCertVerify = "HTTP_CLIENT_SKIP_CERT_VERIFY"
	HTTPClientTimeoutSeconds = "HTTP_CLIENT_TIMEOUT_SECONDS"
	EnableWebhook            = "ENABLE_WEBHOOK"
	WebhookNS                = "WEBHOOK_NAMESPACE"
	WebhookServiceName       = "WEBHOOK_SERVICE_NAME"
	WebhookPort              = "WEBHOOK_SERVICE_PORT"
	WebhookCertSecret        = "WEBHOOK_CERT_SECRET_NAME" //nolint:gosec // This is not a hardcoded secret
	EnableLeaderElection     = "ENABLE_LEADER_ELECTION"

	trueString = "true"

	// This default are applied when running the app locally.
	defaultWebhookPort              = 9443
	defaultMetricsPort              = 8080
	defaultProbesPort               = 8080
	defaultLogFormat                = "console"
	defaultLogLevel                 = "debug"
	defaultLogTimestampField        = "timestamp"
	defaultLogTimestampFormat       = "iso-8601"
	defaultLogTraceIdField          = "reconcile-id"
	defaultHTTPClientTimeoutSeconds = 10
)

var Config = struct {
	WatchNS                  string
	ReleaseNS                string
	ApplyCRDs                bool
	EnableMetrics            bool
	SecureMetrics            bool
	MetricsCertDir           string
	MetricsPort              int
	ProbePort                int
	DisableJSONLogs          bool
	LogFormat                string
	LogLevel                 string
	LogTimestampField        string
	LogTimestampFormat       string
	CMTemplate404Name        string
	CMTemplate404NS          string
	HTTPClientSkipCertVerify bool
	HTTPClientTimeoutSeconds int
	EnableWebhook            bool
	WebhookNS                string
	WebhookService           string
	WebhookPort              int
	WebhookCertSecret        string
	EnableLeaderElection     bool
=======
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
	InsecureSkipCertVerify = "INSECURE_SKIP_CERT_VERIFY"
	TrueString             = "true"
	IngressClasses         = "INGRESS_CLASSES"
)

var Config = struct {
	NS                 string
	ApplyCRDs          bool
	EnableMetrics      bool
	Development        bool
	CMTemplate404Name  string
	CMTemplate404NS    string
	InsecureSkipVerify bool
	IngressClasses     []string
>>>>>>> c3c1a19 (feat: introduce pem registry)
}{}

func GetMetricsAddr() string {
	if !Config.EnableMetrics {
		return "0" // disables metrics
	}
	return fmt.Sprintf(":%d", Config.MetricsPort)
}

func init() {
<<<<<<< HEAD
	Config.WatchNS = os.Getenv(WatchNS)
	Config.ReleaseNS = os.Getenv(WebhookNS)
	Config.ApplyCRDs = os.Getenv(ApplyCRDs) == trueString
	Config.DisableJSONLogs = os.Getenv(DisableJSONLogs) == trueString
	Config.LogFormat = getStringOrDefault(LogFormat, defaultLogFormat)
	Config.LogLevel = getStringOrDefault(LogLevel, defaultLogLevel)
	Config.LogTimestampField = getStringOrDefault(Config.LogTimestampField, defaultLogTimestampField)
	Config.LogTimestampFormat = getStringOrDefault(LogTimestampFormat, defaultLogTimestampFormat)
	Config.CMTemplate404Name = os.Getenv(CMTemplate404Name)
	Config.CMTemplate404NS = os.Getenv(CMTemplate404NS)
	Config.HTTPClientSkipCertVerify = os.Getenv(HTTPClientSkipCertVerify) == trueString
	Config.HTTPClientTimeoutSeconds = parseInt(HTTPClientTimeoutSeconds, defaultHTTPClientTimeoutSeconds)
	Config.EnableMetrics = os.Getenv(EnableMetrics) == trueString
	Config.SecureMetrics = os.Getenv(SecureMetrics) == trueString
	Config.MetricsCertDir = os.Getenv(MetricsCertDir)
	Config.MetricsPort = parseInt(MetricsPort, defaultMetricsPort)
	Config.ProbePort = parseInt(ProbePort, defaultProbesPort)
	Config.EnableWebhook = os.Getenv(EnableWebhook) == trueString
	Config.WebhookNS = os.Getenv(WebhookNS)
	Config.WebhookService = os.Getenv(WebhookServiceName)
	Config.WebhookCertSecret = os.Getenv(WebhookCertSecret)
	Config.WebhookPort = parseInt(WebhookPort, defaultWebhookPort)
	Config.EnableLeaderElection = getStringOrDefault(EnableLeaderElection, trueString) == trueString
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
=======
	Config.NS = os.Getenv(NS)
	Config.ApplyCRDs = os.Getenv(ApplyCRDs) == TrueString
	Config.Development = os.Getenv(Development) == TrueString
	Config.CMTemplate404Name = os.Getenv(CMTemplate404Name)
	Config.CMTemplate404NS = os.Getenv(CMTemplate404NS)
	Config.InsecureSkipVerify = os.Getenv(InsecureSkipCertVerify) == TrueString
	Config.EnableMetrics = os.Getenv(EnableMetrics) == TrueString
	var ingressClass string
	if ingressClass = keys.IngressClassAnnotationValue; os.Getenv(IngressClasses) != "" {
		ingressClass = os.Getenv(IngressClasses)
	}
	Config.IngressClasses = strings.Split(ingressClass, ",")
>>>>>>> c3c1a19 (feat: introduce pem registry)
}
