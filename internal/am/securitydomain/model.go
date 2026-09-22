package securitydomain

import (
	"github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/pkg/sdk/domain"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/mapper"
)

func ToDomainDTO(obj *v1alpha1.AMSecurityDomain) domain.Domain {
	return mapper.MapViaJSON[domain.Domain](obj)
}

type DomainResponse struct {
	am.OrgEnv
	Key string `json:"key"`
}
