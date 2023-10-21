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
}{}

func init() {
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
}
