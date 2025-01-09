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

package mctx

import (
	"context"
	"fmt"
	"slices"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	gioerr "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"k8s.io/apimachinery/pkg/runtime"
)

const cloudGateUrlTemplate = "https://%s.cloudgate.gravitee.io"

// SetDefaults when cloud mode is enabled.
func SetDefaults(ctx context.Context, obj runtime.Object) error {
	if contextObject, ok := obj.(core.ContextObject); ok {
		if contextObject.HasCloud() && contextObject.GetCloud().IsEnabled() {
			return defaultContextUsingCloudToken(ctx, contextObject)
		}
	}

	return nil
}

func defaultContextUsingCloudToken(ctx context.Context, contextObject core.ContextObject) error {
	var cloudToken string

	if contextObject.GetCloud().HasSecretRef() {
		secret, err := dynamic.ResolveSecret(ctx, contextObject.GetCloud().GetSecretRef(), contextObject.GetNamespace())
		if err != nil {
			return gioerr.NewSeveref("secret [%v] doesn't exist in the cluster", contextObject.GetCloud().GetSecretRef())
		}
		cloudToken = string(secret.Data[core.CloudTokenSecretKey])
	} else {
		cloudToken = contextObject.GetCloud().GetToken()
	}

	jwtData, err := extractCloudTokenData(cloudToken)
	if err != nil {
		return gioerr.NewSeveref("cannot parse cloud token: %s", err)
	}

	var url string
	var orgID string
	var envID string

	if contextObject.GetURL() == "" {
		url = jwtData.baseUrl()
	} else {
		url = contextObject.GetURL()
	}

	orgID = jwtData.Org

	switch {
	case contextObject.GetEnvID() == "":
		if len(jwtData.Envs) > 1 {
			return gioerr.NewSeveref(
				"cloud token contains more than one environment (%d), environmentId is then required",
				len(jwtData.Envs))
		}
		envID = jwtData.Envs[0]
	case !slices.Contains(jwtData.Envs, contextObject.GetEnvID()):
		return gioerr.NewSeveref("cloud token does not contain environment [%s], it must be one of: %s",
			contextObject.GetEnvID(),
			jwtData.Envs)
	default:
		envID = contextObject.GetEnvID()
	}

	contextObject.ConfigureCloud(url, orgID, envID)

	return nil
}

func extractCloudTokenData(jwtToken string) (CloudTokenClaimsData, error) {
	claims := &CloudTokenClaims{}

	_, _, err := jwt.NewParser().ParseUnverified(jwtToken, claims)

	if err != nil {
		return CloudTokenClaimsData{}, err
	}

	if !claims.isValid() {
		return CloudTokenClaimsData{}, gioerr.NewSeveref("cloud token does not contains all required claims, " +
			"are you sure this is a Gravitee cloud token ?")
	}

	return claims.CloudTokenClaimsData, nil
}

type CloudTokenClaims struct {
	jwt.RegisteredClaims
	CloudTokenClaimsData
}

type CloudTokenClaimsData struct {
	Org       string   `json:"org"`
	Envs      []string `json:"envs"`
	Geography string   `json:"cpg"`
}

func (d CloudTokenClaimsData) baseUrl() string {
	return fmt.Sprintf(cloudGateUrlTemplate, d.Geography)
}

func (d CloudTokenClaimsData) isValid() bool {
	return d.Org != "" && d.Envs != nil && len(d.Envs) > 0 && d.Geography != ""
}
