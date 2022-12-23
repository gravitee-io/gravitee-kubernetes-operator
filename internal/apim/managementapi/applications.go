// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package managementapi

import (
	"encoding/json"
	"fmt"
	"io"
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
		client.envUrl()+"/applications?status="+status+"&query="+query,
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

	body, readErr := io.ReadAll(resp.Body)
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
