package managementapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi/model"
)

func (client *Client) SearchApplications(
	query string,
	status string,
) ([]model.Application, error) {
	req, err := http.NewRequestWithContext(
		client.ctx,
		http.MethodGet,
		client.buildUrl("/applications?status="+status+"&query="+query),
		nil,
	)

	if err != nil && status == "" {
		return nil, fmt.Errorf("unable to look for applications (%w)", err)
	}

	resp, err := client.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("an error as occurred while performing FindApplications request")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("an error as occurred trying to find applications, HTTP Status: %d ", resp.StatusCode)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	var applications []model.Application

	err = json.Unmarshal(body, &applications)
	if err != nil {
		return nil, err
	}

	return applications, nil
}
