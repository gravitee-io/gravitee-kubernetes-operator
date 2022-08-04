package context

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
)

type Delegate struct {
	ctx context.Context
	cli client.Client
	log logr.Logger
}

func NewDelegate(ctx context.Context, cli client.Client) *Delegate {
	log := log.FromContext(ctx)
	return &Delegate{ctx, cli, log}
}
