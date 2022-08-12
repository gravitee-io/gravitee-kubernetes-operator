package apis

import (
	"context"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	apim "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Delegate struct {
	ctx               context.Context
	managementContext *gio.ManagementContext
	apimClient        *apim.Client
	k8sClient         client.Client
	log               logr.Logger
}

func NewDelegate(ctx context.Context, client client.Client) *Delegate {
	log := log.FromContext(ctx)

	return &Delegate{
		ctx, nil, nil, client, log,
	}
}

func (d *Delegate) SetManagementContext(managementContext *gio.ManagementContext) {
	if managementContext == nil {
		return
	}

	d.managementContext = managementContext

	httpClient := http.Client{Timeout: requestTimeoutSeconds * time.Second}
	d.apimClient = apim.NewClient(d.ctx, d.managementContext, httpClient)
}

func (d *Delegate) IsConnectedToManagementApi() bool {
	return d.apimClient != nil
}
