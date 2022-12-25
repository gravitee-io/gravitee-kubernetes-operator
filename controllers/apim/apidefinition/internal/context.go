// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"text/template"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	apim "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

const (
	bearerTokenSecretKey = "bearerToken"
	usernameSecretKey    = "username"
	passwordSecretKey    = "password"
)

var templateEngine = template.New("api-definitions")

type DelegateContext struct {
	*apim.Client
	// Location is the namespace and name of the API in the form namespace/name
	Location string
	// Values is a map of values that can be passed as a golang template context to the API spec
	Values map[string]string
}

func (c *DelegateContext) compile(api *gio.ApiDefinition) (*gio.ApiDefinition, error) {
	if !c.hasValues() {
		return api, nil
	}

	cp := api.DeepCopy()
	spec := cp.Spec

	jsonSpec, mErr := json.Marshal(spec)
	if mErr != nil {
		return api, mErr
	}

	template, pErr := templateEngine.Parse(string(jsonSpec))
	if pErr != nil {
		return api, pErr
	}

	var buff bytes.Buffer
	if tErr := template.Execute(&buff, c); tErr != nil {
		return api, tErr
	}

	if uErr := json.Unmarshal(buff.Bytes(), &cp.Spec); uErr != nil {
		return api, uErr
	}

	return cp, nil
}

func (c *DelegateContext) hasManagement() bool {
	return c.Client != nil
}

func (c *DelegateContext) hasValues() bool {
	return c.Values != nil
}

func (c *DelegateContext) update(api *gio.ApiDefinition) error {
	if !c.hasManagement() {
		return nil
	}

	if _, ok := api.Status.Contexts[c.Location]; !ok {
		api.Status.Contexts[c.Location] = gio.StatusContext{}
	}

	spec := &api.Spec
	spec.ID = api.PickID(c.Location)
	spec.CrossID = api.PickCrossID(c.Location)

	spec.SetDefinitionContext()

	generateEmptyPlanCrossIds(spec)

	statusContext := api.Status.Contexts[c.Location]

	_, findErr := c.APIs.GetByCrossID(spec.CrossID)
	if errors.IgnoreNotFound(findErr) != nil {
		return newContextError(c.Location, findErr)
	}

	importMethod := http.MethodPost
	if findErr == nil {
		importMethod = http.MethodPut
	}

	mgmtApi, mgmtErr := c.APIs.Import(importMethod, &spec.Api)
	if mgmtErr != nil {
		return newContextError(c.Location, mgmtErr)
	}

	if mgmtApi.ShouldSetKubernetesContext() {
		if err := c.APIs.SetKubernetesContext(mgmtApi.ID); err != nil {
			return newContextError(c.Location, err)
		}
	}

	retrieveMgmtPlanIds(spec, mgmtApi)

	statusContext.ID = mgmtApi.ID
	statusContext.CrossID = spec.CrossID
	statusContext.State = spec.State
	statusContext.OrgID = c.OrgID()
	statusContext.EnvID = c.EnvID()
	statusContext.Status = gio.ProcessingStatusCompleted

	api.Status.Contexts[c.Location] = statusContext

	return nil
}

func (c *DelegateContext) delete(api *gio.ApiDefinition) error {
	if !c.hasManagement() {
		return nil
	}

	if _, ok := api.Status.Contexts[c.Location]; !ok {
		return nil
	}

	status := api.Status.Contexts[c.Location]

	return errors.IgnoreNotFound(c.APIs.Delete(status.ID))
}
