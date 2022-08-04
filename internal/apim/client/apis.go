package apim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client/model"
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
		return nil, fmt.Errorf("unable to look for apis matching cross id %s (%w)", crossId, err)
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

func (client *Client) Import(
	importHttpMethod string,
	apiJson []byte,
) error {
	url := client.buildUrl("/apis/import?definitionVersion=2.0.0")
	req, err := http.NewRequestWithContext(client.ctx, importHttpMethod, url, bytes.NewBuffer(apiJson))

	if err != nil {
		return fmt.Errorf("unable to import the api into the Management API")
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := client.http.Do(req)

	if err != nil {
		return err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("management has returned a %d code", resp.StatusCode)
	}

	return err
}
