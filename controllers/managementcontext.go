package controllers

import (
	"context"

	"github.com/go-logr/logr"
	graviteeiov1alpha1 "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const defaultNamespace = "default"

func getManagementContext(
	ctx context.Context,
	client client.Client,
	log logr.Logger,
	apiDefinition *graviteeiov1alpha1.ApiDefinition,
) (graviteeiov1alpha1.ManagementContext, error) {
	contextRef := apiDefinition.Spec.Context

	// If namespace is not specified in contextRef, use default namespace
	if contextRef.Namespace == "" {
		log.Info("Context namespace is not specified, using default")

		contextRef.Namespace = defaultNamespace
	}

	var mgmtContext graviteeiov1alpha1.ManagementContext
	var ns = types.NamespacedName{Name: contextRef.Name, Namespace: contextRef.Namespace}

	log.Info("Lookup for context from", "namespace", contextRef.Namespace, "name", contextRef.Name)

	err := client.Get(ctx, ns, &mgmtContext)

	if err != nil {
		return mgmtContext, err
	}

	return mgmtContext, nil
}
