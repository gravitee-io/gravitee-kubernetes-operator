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

package model

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
)

type ApiEntity struct {
	ID                string             `json:"id"`
	CrossID           string             `json:"crossId"`
	Name              string             `json:"name"`
	State             string             `json:"state"`
	Visibility        string             `json:"visibility"`
	ApiLifecycleState string             `json:"lifecycle_state"`
	Plans             []*Plan            `json:"plans"`
	Resources         []*Resource        `json:"resources,omitempty"`
	DefinitionContext *DefinitionContext `json:"definition_context,omitempty"`
}

func (api *ApiEntity) ShouldSetKubernetesContext() bool {
	return api.DefinitionContext == nil || api.DefinitionContext.Origin == OriginManagement
}

type DefinitionContext struct {
	Origin string `json:"origin,omitempty"`
	Mode   string `json:"mode,omitempty"`
}

const (
	OriginManagement = "management"
	OriginKubernetes = "kubernetes"
	ModeFullyManaged = "fully_managed"
)

func NewKubernetesContext() *DefinitionContext {
	return &DefinitionContext{
		Origin: OriginKubernetes,
		Mode:   ModeFullyManaged,
	}
}

type ApiListItem struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	State             string `json:"state"`
	Visibility        string `json:"visibility"`
	ApiLifecycleState string `json:"lifecycle_state"`
}

type Action string

const (
	ActionStart Action = "START"
	ActionStop  Action = "STOP"
)

func ApiStateToAction(s base.ApiState) Action {
	switch s {
	case base.StateStarted:
		return ActionStart
	case base.StateStopped:
		return ActionStop
	default:
		return ActionStop
	}
}

type Plan struct {
	Id       string           `json:"id"`
	CrossId  string           `json:"crossId"`
	Name     string           `json:"name"`
	Security PlanSecurityType `json:"security"`
	Status   PlanStatus       `json:"status"`
	Api      string           `json:"api"`
}

type Resource struct {
	Enabled      bool   `json:"enabled"`
	Name         string `json:"name,omitempty"`
	ResourceType string `json:"type,omitempty"`
}

type PlanSecurityType string

type PlanStatus string

type ApiDeployment struct {
	DeploymentLabel string `json:"deploymentLabel"`
}

type ApiImport struct {
	*v2.Api `json:",inline"`
	Pages   []*PageImport `json:"pages,omitempty"`
}
