package apis

import (
	"fmt"
	"net/http"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

func (d *Delegate) findByCrossId(
	apimCtx *gio.ManagementContext,
	apiId string,
) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		d.ctx,
		http.MethodGet,
		apimCtx.Spec.BuildUrl("/apis?crossId="+apiId),
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("an error as occured while trying to create new findApisByCrossId request")
	}

	apimCtx.Spec.Authenticate(req)
	resp, err := d.http.Do(req)

	return resp, err
}
