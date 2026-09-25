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

package am

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/v2/pkg"
	amsdk "github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io-labs/gravitee-automation-tools/common/pkg/apicontext"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	ghttp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

// Client wraps the Automation API services needed to talk to one Client environment.
type Client struct {
	*pkg.AMClient
	Context core.ContextModel
}

func NewSDKClient(ctx context.Context, obj *v1alpha1.AMContext) (*Client, error) {

	if _, err := dynamic.InjectSecretIfAny(ctx, obj); err != nil {
		return nil, err
	}

	client, err := ghttp.NewStdClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}
	baseClient, err := pkg.NewClient(toSDKContext(obj.Spec))
	if err != nil {
		return nil, fmt.Errorf("failed to create AMClient: %w", err)
	}
	amClient, err := baseClient.WithHTTPClient(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create AMClient: %w", err)
	}

	return &Client{
		AMClient: amClient,
		Context:  obj,
	}, nil
}

func (c *Client) GetOrgID() string {
	return c.Context.GetOrgID()
}

func (c *Client) GetEnvID() string {
	return c.Context.GetEnvID()
}

// Probe checks that AM's Automation API is reachable at this org/env.
// A 200 is the whole contract: the body is discarded.
func (c *Client) Probe(ctx context.Context) error {
	resp, err := c.ListDomainsWithResponse(ctx, func(ctx context.Context, req *http.Request) error {
		// Add page and length
		req.URL.RawQuery = url.Values{"size": {"1"}}.Encode()
		return nil
	})

	if err := HasErrors(err, func() *http.Response {
		return resp.HTTPResponse
	}); err != nil {
		return err
	}

	return err
}

func toSDKContext(spec v1alpha1.AMContextSpec) apicontext.APIContext {
	baseUrl := spec.BaseUrl
	path := "/automation"
	if spec.Path != nil && strings.TrimSpace(*spec.Path) != "" {
		path = *spec.Path
	}
	baseUrl = fmt.Sprintf("%s%s", baseUrl, path)
	return apicontext.APIContext{
		BaseURL: baseUrl,
		OrgID:   spec.OrgID,
		EnvID:   spec.EnvID,
		Auth: apicontext.Auth{
			BearerToken: new(spec.GetAuth().GetBearerToken()),
		},
	}
}

func ToAdmissionErrors(errors []amsdk.DryRunError) *gerrors.AdmissionErrors {
	errs := gerrors.NewAdmissionErrors()
	for _, err := range errors {
		if err.Severity != nil && err.Message != nil {
			switch *err.Severity {
			case amsdk.SeverityError:
				errs.AddSevere(*err.Message)
			case amsdk.SeverityWarning:
				errs.AddWarning(*err.Message)
			}
		}
	}
	return errs
}

func HasErrors(err error, response func() *http.Response) error {
	if err != nil {
		return err
	}
	return gerrors.FromResponse(response())
}
