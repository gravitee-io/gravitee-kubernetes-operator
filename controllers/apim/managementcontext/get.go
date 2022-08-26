package managementcontext

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Get(
	ctx context.Context,
	k8sClient client.Client,
	log logr.Logger,
	contextRef *model.ContextRef,
) (*gio.ManagementContext, error) {
	apimContext := new(gio.ManagementContext)
	ns := types.NamespacedName{Name: contextRef.Name, Namespace: contextRef.Namespace}

	log.Info("Looking for context from", "namespace", contextRef.Namespace, "name", contextRef.Name)

	err := k8sClient.Get(ctx, ns, apimContext)

	if err != nil {
		return nil, err
	}

	return apimContext, nil
}
