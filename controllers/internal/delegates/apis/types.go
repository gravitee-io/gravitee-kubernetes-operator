package apis

import (
	"context"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Delegate struct {
	cli  client.Client
	ctx  context.Context
	log  logr.Logger
	http http.Client
}

func NewDelegate(ctx context.Context, client client.Client) *Delegate {
	log := log.FromContext(ctx)
	http := http.Client{Timeout: requestTimeoutSeconds * time.Second}

	return &Delegate{
		client, ctx, log, http,
	}
}
