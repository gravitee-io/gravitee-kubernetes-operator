package apis

import (
	"context"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/apim"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Delegate struct {
	ctx        context.Context
	apimCtx    *gio.ManagementContext
	apimClient *apim.Client
	cli        client.Client
	log        logr.Logger
}

func NewDelegate(ctx context.Context, apimCtx *gio.ManagementContext, client client.Client) *Delegate {
	log := log.FromContext(ctx)

	var apimClient *apim.Client

	if apimCtx != nil {
		httpClient := http.Client{Timeout: requestTimeoutSeconds * time.Second}
		apimClient = apim.NewClient(ctx, apimCtx, httpClient)
	}

	return &Delegate{
		ctx, apimCtx, apimClient, client, log,
	}
}
