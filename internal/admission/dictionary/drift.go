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

package dictionary

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func validateUpdate(ctx context.Context, oldDict *v1alpha1.Dictionary, newDict *v1alpha1.Dictionary) *errors.AdmissionErrors {
	errs := validateCreate(ctx, newDict)
	if errs.IsSevere() {
		return errs
	}
	errs.MergeWith(drift.ValidateDrift(ctx, oldDict, newDict, resolveRefs, getRemoteDictionary,
		drift.MapDTO(toDictionaryDTO)))
	return errs
}

func resolveRefs(context.Context, *v1alpha1.Dictionary) error {
	return nil
}

func toDictionaryDTO(dict *v1alpha1.Dictionary) model.DictionaryDTO {
	return model.ToDictionaryDTO(
		dict.Spec.Type,
		refs.NewNamespacedNameFromObject(dict).HRID(),
	)
}

func getRemoteDictionary(apimClient *apim.APIM, dict *v1alpha1.Dictionary, errs *errors.AdmissionErrors) any {
	hrid := refs.NewNamespacedNameFromObject(dict).HRID()
	remote, err := apimClient.Dictionaries.GetByHRID(hrid)
	if err != nil {
		errs.AddSeveref("cannot fetch Dictionary during drift detection from HRID %s: %s", hrid, err.Error())
		return nil
	}
	return remote.DictionaryDTO
}
