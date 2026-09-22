package securitydomain

import (
	"context"
	"fmt"
	"net/http"

	amsdk "github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/pkg/sdk/domain"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

func Delete(ctx context.Context, amClient *am.Client, dto amsdk.Domain) error {
	resp, err := amClient.Domains.DeleteDomainWithResponse(ctx, dto.Identity())
	return am.HasErrors(err, func() *http.Response {
		return resp.HTTPResponse
	})
}

func Upsert(ctx context.Context, client *am.Client, dto amsdk.Domain) (DomainResponse, error) {
	resp, err := client.Domains.UpsertDomainWithResponse(ctx, nil, dto)

	if err = am.HasErrors(err, func() *http.Response {
		return resp.HTTPResponse
	}); err != nil {
		return DomainResponse{}, err
	}

	return DomainResponse{
		Key: resp.JSON200.Key,
		OrgEnv: am.OrgEnv{
			OrgID: client.GetOrgID(),
			EnvID: client.GetEnvID(),
		},
	}, nil
}

func UpdateStatus(_ context.Context, obj *v1alpha1.AMSecurityDomain, resp DomainResponse) error {
	obj.Status.Status.ID = resp.Key
	obj.Status.Status.OrgID = resp.GetOrgID()
	obj.Status.Status.EnvID = resp.GetEnvID()
	return nil
}

func CreateAMClient(ctx context.Context, obj *v1alpha1.AMSecurityDomain) (*am.Client, error) {
	if obj.HasContext() {
		return nil, fmt.Errorf("contextRef empty on %s [%s/%s]", obj.Kind, obj.GetName(), obj.GetNamespace())
	}

	amContext := &v1alpha1.AMContext{}
	err := k8s.GetClient().Get(ctx, obj.ContextRef().NamespacedName(), amContext)
	if err != nil {
		return nil, fmt.Errorf("AMContext [%s] not found", obj.ContextRef().NamespacedName().String())
	}

	return am.NewSDKClient(ctx, amContext)
}

func DryRun(ctx context.Context, client *am.Client, dto amsdk.Domain) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	resp, err := client.Domains.UpsertDomainWithResponse(ctx, new(amsdk.UpsertDomainParams{
		DryRun: new(true),
	}), dto)

	if err = am.HasErrors(err, func() *http.Response {
		return resp.HTTPResponse
	}); err != nil {
		errs.AddSevere(err.Error())
		return errs
	}
	if resp.JSON200.DryRunErrors == nil {
		return errs
	}
	return am.ToAdmissionErrors(*resp.JSON200.DryRunErrors)
}

func GetRemote(ctx context.Context, client *am.Client, dto amsdk.Domain) (amsdk.Domain, error) {
	resp, err := client.Domains.GetDomainWithResponse(ctx, dto.Key)
	if err = am.HasErrors(err, func() *http.Response {
		return resp.HTTPResponse
	}); err != nil {
		return amsdk.Domain{}, err
	}
	return *resp.JSON200, nil
}
