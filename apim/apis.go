package apim

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/apim/model"
)

func (client *Client) FindByCrossId(
	crossId string,
) ([]model.Api, error) {
	req, err := http.NewRequestWithContext(
		client.ctx,
		http.MethodGet,
		client.buildUrl("/apis?crossId="+crossId),
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("an error as occurred while trying to create new findApisByCrossId request")
	}

	resp, err := client.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("an error as occurred while performing findApisByCrossId request")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// TODO parse response body as a map and log
		return nil, fmt.Errorf("an error as occurred trying to find API %s, HTTP Status: %d ", crossId, resp.StatusCode)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	var apis []model.Api

	err = json.Unmarshal(body, &apis)
	if err != nil {
		return nil, err
	}

	return apis, err
}
