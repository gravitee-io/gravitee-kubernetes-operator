package managementapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi/clienterror"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi/model"
)

func (client *Client) GetByCrossId(
	crossId string,
) (*model.ApiListItem, error) {
	req, err := http.NewRequestWithContext(
		client.ctx,
		http.MethodGet,
		client.buildUrl("/apis?crossId="+crossId),
		nil,
	)

	if err != nil && crossId == "" {
		return nil, fmt.Errorf("unable to look for apis matching cross id %s (%w)", crossId, err)
	}

	resp, err := client.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("an error as occurred while performing GetByCrossId request")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// TODO parse response body as a map and log
		return nil, fmt.Errorf("an error as occurred trying to find API %s, HTTP Status: %d ", crossId, resp.StatusCode)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	var apis []model.ApiListItem

	err = json.Unmarshal(body, &apis)
	if err != nil {
		return nil, err
	}

	if len(apis) == 0 {
		return nil, &clienterror.CrossIdNotFoundError{CrossId: crossId}
	}

	if len(apis) > 1 {
		return nil, &clienterror.CrossIdMultipleFoundError{CrossId: crossId, Apis: apis}
	}

	return &apis[0], nil
}

func (client *Client) GetApiById(
	apiId string,
) (*model.ApiEntity, error) {
	req, err := http.NewRequestWithContext(
		client.ctx,
		http.MethodGet,
		client.buildUrl("/apis/"+apiId),
		nil,
	)

	if err != nil && apiId == "" {
		return nil, fmt.Errorf("unable to look for apis matching id %s (%w)", apiId, err)
	}
	resp, err := client.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("an error as occurred while performing GetApiById request")
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, &clienterror.ApiNotFoundError{ApiId: apiId}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"an error as occurred trying to get API matching id %s, HTTP Status: %d ",
			apiId, resp.StatusCode)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	var api model.ApiEntity

	err = json.Unmarshal(body, &api)
	if err != nil {
		return nil, err
	}

	return &api, nil
}

func (client *Client) CreateApi(
	apiJson []byte,
) (*model.ApiEntity, error) {
	return importApi(client, http.MethodPost, apiJson)
}

func (client *Client) UpdateApi(
	apiJson []byte,
) (*model.ApiEntity, error) {
	return importApi(client, http.MethodPut, apiJson)
}

func importApi(
	client *Client,
	importHttpMethod string,
	apiJson []byte,
) (*model.ApiEntity, error) {
	url := client.buildUrl("/apis/import?definitionVersion=2.0.0")
	req, err := http.NewRequestWithContext(client.ctx, importHttpMethod, url, bytes.NewBuffer(apiJson))

	if err != nil {
		return nil, fmt.Errorf("unable to import the api into the Management API")
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := client.http.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("management has returned a %d code", resp.StatusCode)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	var api model.ApiEntity

	err = json.Unmarshal(body, &api)
	if err != nil {
		return nil, err
	}

	return &api, nil
}

func (client *Client) UpdateApiState(
	apiId string,
	action model.Action,
) error {
	url := client.buildUrl("/apis/" + apiId + "?action=" + string(action))
	req, err := http.NewRequestWithContext(client.ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("unable to update the api state into the Management API. Action: %s", action)
	}

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

func (client *Client) DeleteApi(
	apiId string,
) error {
	url := client.buildUrl("/apis/" + apiId + "?closePlans=true")
	req, err := http.NewRequestWithContext(client.ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("unable to delete the api into the Management API")
	}

	resp, err := client.http.Do(req)
	if err != nil {
		return err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode == http.StatusNotFound {
		return &clienterror.ApiNotFoundError{ApiId: apiId}
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("management has returned a %d code", resp.StatusCode)
	}

	return err
}
