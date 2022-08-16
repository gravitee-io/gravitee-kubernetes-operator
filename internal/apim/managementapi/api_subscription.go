package managementapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi/model"
)

func (client *Client) SubscribeToPlan(
	apiId string,
	applicationId string,
	planId string,
) (*model.Subscription, error) {
	queryParams := "?application=" + applicationId + "&plan=" + planId

	url := client.buildUrl("/apis/" + apiId + "/subscriptions" + queryParams)

	req, err := http.NewRequestWithContext(client.ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to subscribe to the plan into the Management API, ApiId=%s, ApplicationId=%s, PlanId=%s",
			apiId, applicationId, planId)
	}

	resp, err := client.http.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf(
			"unable to subscribe to the plan into the Management API, ApiId=%s, ApplicationId=%s, PlanId=%s, HTTP Status: %d",
			apiId, applicationId, planId, resp.StatusCode)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	var subscription model.Subscription
	err = json.Unmarshal(body, &subscription)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}
