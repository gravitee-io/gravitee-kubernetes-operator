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

package service

import (
	"strconv"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/dictionary"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

type Dictionaries struct {
	*client.Client
}

func NewDictionaries(client *client.Client) *Dictionaries {
	return &Dictionaries{Client: client}
}

func (svc *Dictionaries) CreateOrUpdate(dict *v1alpha1.Dictionary) (dictionary.Status, error) {
	return svc.createOrUpdate(dict, false)
}

func (svc *Dictionaries) DryRunCreateOrUpdate(dict *v1alpha1.Dictionary) (dictionary.Status, error) {
	return svc.createOrUpdate(dict, true)
}

func (svc *Dictionaries) createOrUpdate(
	dict *v1alpha1.Dictionary,
	dryRun bool,
) (dictionary.Status, error) {
	url := svc.AutomationTarget("dictionaries").
		WithQueryParam("dryRun", strconv.FormatBool(dryRun))

	dto := model.ToDictionaryDTO(dict.Spec.Type, refs.NewNamespacedNameFromObject(dict).HRID())
	importStatus := &dictionary.Status{}

	if err := svc.HTTP.Put(url.String(), dto, &importStatus); err != nil {
		return *importStatus, err
	}

	k8s.AddAutomationAPIManagedCondition(dict)

	return *importStatus, nil
}

func (svc *Dictionaries) Delete(dict *v1alpha1.Dictionary) error {
	hrid := refs.NewNamespacedNameFromObject(dict).HRID()
	url := svc.AutomationTarget("dictionaries").WithPath(hrid)
	return svc.HTTP.Delete(url.String(), nil)
}

// GetByHRID For test purposes only.
func (svc *Dictionaries) GetByHRID(hrid string) (*model.DictionaryState, error) {
	url := svc.AutomationTarget("dictionaries").WithPath(hrid)
	dict := new(model.DictionaryState)
	if err := svc.HTTP.Get(url.String(), dict); err != nil {
		return nil, err
	}
	return dict, nil
}
