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

package apim

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-automation-sdk/common/pkg/apicontext"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	ghttp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// probeHRID is the group hrid Probe gets. A Group named probe-not-found in namespace gko has it too: the
// probe then reads that group, which proves the same thing.
const probeHRID = "gko-probe-not-found"

// SDKClient is the APIM Automation API SDK client for one management context.
type SDKClient struct {
	*sdk.APIMClient
	Context core.ContextModel
}

func (c *SDKClient) GetOrgID() string {
	return c.Context.GetOrgID()
}

func (c *SDKClient) GetEnvID() string {
	return c.Context.GetEnvID()
}

// Probe checks that the Automation API is reachable and accepts the context's credentials at its org/env.
// The SDK has no list endpoint, so it gets one group: a 200 or a 404 both mean reachable and authorised.
// APIM answers an unknown org or env with a 400, and bad credentials with a 401, so neither passes.
func (c *SDKClient) Probe(ctx context.Context) error {
	resp, err := c.GetGroupWithResponse(ctx, probeHRID, &sdk.GetGroupParams{HridContainsUUID: new(false)})
	err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	})
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

// NewSDKClient builds the APIM Automation API SDK client for a management context: its secret-held
// credentials injected, the operator's HTTP client (proxy, TLS, trust store), and the same base URL,
// organization, environment and cloud rules as FromContext.
func NewSDKClient(ctx context.Context, mctx core.ContextObject) (*SDKClient, error) {
	if _, err := dynamic.InjectSecretIfAny(ctx, mctx); err != nil {
		return nil, err
	}
	client, err := newSDKClient(mctx)
	if err != nil {
		return nil, err
	}
	return &SDKClient{APIMClient: client, Context: mctx}, nil
}

// SDKFromContextRef resolves a management context reference and builds its SDK client.
func SDKFromContextRef(ctx context.Context, ref core.ObjectRef, parentNs string) (*SDKClient, error) {
	mctx, err := dynamic.ResolveContext(ctx, ref, parentNs)
	if err != nil {
		return nil, err
	}
	return NewSDKClient(ctx, mctx)
}

// errNoContextRef is a NotFound, as a deleted context is: nothing reached APIM without a context either.
var errNoContextRef = &apierrors.StatusError{ErrStatus: metav1.Status{
	Status:  metav1.StatusFailure,
	Code:    http.StatusNotFound,
	Reason:  metav1.StatusReasonNotFound,
	Message: "a management context reference (spec.contextRef) is required",
}}

// ResolveContextRef is the lifecycle ResolveRefs of a resource that requires a management context.
// An unset or deleted context is NotFound, so deleting the resource releases its finalizer instead of
// retrying a call APIM can no longer receive.
func ResolveContextRef[T core.ContextAwareObject](ctx context.Context, obj T, ns string) error {
	if !obj.HasContext() {
		return errNoContextRef
	}
	_, err := dynamic.ResolveContext(ctx, obj.ContextRef(), ns)
	return err
}

// newSDKClient builds the SDK client of a context whose secrets are already injected.
// The SDK checks the authentication when the client is built, not on the first call.
func newSDKClient(mctx core.ContextModel) (*sdk.APIMClient, error) {
	baseURL, err := url.JoinPath(mctx.GetURL(), getAutomationPath(mctx))
	if err != nil {
		return nil, fmt.Errorf("invalid APIM URL [%s] with path [%s]: %w", mctx.GetURL(), getAutomationPath(mctx), err)
	}

	apiContext := sdk.APIContext{
		APIContext: apicontext.APIContext{
			BaseURL: baseURL,
			OrgID:   mctx.GetOrgID(),
			EnvID:   mctx.GetEnvID(),
			Auth:    toSDKAuth(mctx),
		},
		CloudToken: hasEnabledCloud(mctx),
	}

	client, err := sdk.NewAPIMClient(apiContext)
	if err != nil {
		return nil, fmt.Errorf("failed to create APIM SDK client: %w", err)
	}

	httpClient, err := ghttp.NewStdClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	client, err = client.WithHTTPClient(httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create APIM SDK client: %w", err)
	}
	return client, nil
}

// toSDKAuth sets exactly one of bearer or basic: the SDK rejects both. Secret injection always sets
// credentials, possibly empty, so a bearer token wins over them.
func toSDKAuth(mctx core.ContextModel) apicontext.Auth {
	if !mctx.HasAuthentication() {
		return apicontext.Auth{}
	}
	if token := toBearer(mctx.GetAuth()); token != "" {
		value := string(token)
		return apicontext.Auth{BearerToken: &value}
	}
	if basic := toBasicAuth(mctx.GetAuth()); basic != nil {
		return apicontext.Auth{BasicAuth: &apicontext.BasicAuth{Username: basic.Username, Password: basic.Password}}
	}
	return apicontext.Auth{}
}

// hasEnabledCloud is true when the context holds an enabled Gravitee Cloud configuration. A cloud block
// that is present but disabled leaves the context a plain one, as for the cloud defaulting.
func hasEnabledCloud(mctx core.ContextModel) bool {
	cloud, ok := mctx.(core.CloudAwareContext)
	return ok && cloud.HasCloud() && cloud.GetCloud().IsEnabled()
}
