package apim

import (
	"context"
	"net/http"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

type Client struct {
	ctx      context.Context
	buildUrl func(string) string
	http     http.Client
}

type AuthenticatedRoundTripper struct {
	apimCtx   *gio.ManagementContext
	transport http.RoundTripper
}

func newAuthenticatedRoundTripper(
	apimCtx *gio.ManagementContext,
	transport http.RoundTripper,
) *AuthenticatedRoundTripper {
	return &AuthenticatedRoundTripper{
		apimCtx, transport,
	}
}

func (t *AuthenticatedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	t.apimCtx.Spec.Authenticate(req)
	return t.transport.RoundTrip(req)
}

func NewClient(ctx context.Context, apimCtx *gio.ManagementContext, httpCli http.Client) *Client {
	buildUrl := apimCtx.Spec.BuildUrl
	authRoundTripper := newAuthenticatedRoundTripper(apimCtx, http.DefaultTransport)
	httpCli.Transport = authRoundTripper

	return &Client{
		ctx, buildUrl, httpCli,
	}
}
