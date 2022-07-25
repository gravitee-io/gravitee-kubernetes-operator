package internal

import (
	"context"
	"net/http"
	"strings"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

func GetApimContext(
	ctx context.Context,
	client client.Client,
	api *gio.ApiDefinition,
) (*gio.ManagementContext, error) {
	contextRef := api.Spec.Context

	apimContext := new(gio.ManagementContext)
	ns := types.NamespacedName{Name: contextRef.Name, Namespace: contextRef.Namespace}

	log.FromContext(ctx).Info("Looking for context from", "namespace", contextRef.Namespace, "name", contextRef.Name)

	err := client.Get(ctx, ns, apimContext)

	if err != nil {
		return nil, err
	}

	return apimContext, nil
}

func BuildApimUrl(apimCtx *gio.ManagementContext, path string) string {
	orgId, envId := apimCtx.Spec.OrgId, apimCtx.Spec.EnvId
	baseUrl := strings.TrimSuffix(apimCtx.Spec.BaseUrl, "/")
	url := baseUrl + "/management/organizations/" + orgId
	if envId != "" {
		url = url + "/environments/" + envId
	}
	return url + path
}

func SetApimAuth(apimCtx *gio.ManagementContext, request *http.Request) {
	if apimCtx.Spec.Auth != nil {
		bearerToken := apimCtx.Spec.Auth.BearerToken
		if bearerToken != "" {
			request.Header.Add("Authorization", "Bearer "+bearerToken)
		} else if apimCtx.Spec.Auth.Credentials != nil {
			username := apimCtx.Spec.Auth.Credentials.Username
			password := apimCtx.Spec.Auth.Credentials.Password
			if username != "" {
				request.SetBasicAuth(username, password)
			}
		}
	}
}
