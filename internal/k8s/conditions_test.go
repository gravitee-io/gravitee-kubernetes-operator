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

package k8s

import (
	"fmt"
	"testing"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestErrorToCondition_PreservesAutomationAPIManaged(t *testing.T) {
	api := &v1alpha1.ApiV4Definition{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "test-api",
			Namespace:  "default",
			Generation: 2,
		},
	}

	api.Status.Conditions = []metav1.Condition{
		{
			Type:               AutomationAPIManaged,
			Status:             metav1.ConditionTrue,
			Reason:             AutomationAPIManaged,
			ObservedGeneration: 1,
		},
	}

	err := gerrors.NewControlPlaneError(fmt.Errorf("some transient APIM error"))
	ErrorToCondition(api, err)

	conditions := api.GetConditions()

	if _, ok := conditions[AutomationAPIManaged]; !ok {
		t.Fatal("AutomationAPIManaged condition was wiped by ErrorToCondition")
	}
	if conditions[AutomationAPIManaged].Status != metav1.ConditionTrue {
		t.Fatalf("AutomationAPIManaged condition status changed to %s, want True",
			conditions[AutomationAPIManaged].Status)
	}

	if _, ok := conditions[ConditionAccepted]; !ok {
		t.Fatal("Accepted condition missing after ErrorToCondition")
	}
	if conditions[ConditionAccepted].Status != ConditionStatusFalse {
		t.Fatalf("Accepted condition status is %s, want False", conditions[ConditionAccepted].Status)
	}

	if _, ok := conditions[ConditionResolvedRefs]; !ok {
		t.Fatal("ResolvedRefs condition missing after ErrorToCondition")
	}
}

func TestErrorToCondition_WithoutAutomationAPIManaged(t *testing.T) {
	api := &v1alpha1.ApiV4Definition{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "test-api",
			Namespace:  "default",
			Generation: 1,
		},
	}
	api.Status.Conditions = []metav1.Condition{}

	err := gerrors.NewControlPlaneError(fmt.Errorf("some error"))
	ErrorToCondition(api, err)

	conditions := api.GetConditions()

	if _, ok := conditions[AutomationAPIManaged]; ok {
		t.Fatal("AutomationAPIManaged condition should not appear when it was not previously set")
	}

	if _, ok := conditions[ConditionAccepted]; !ok {
		t.Fatal("Accepted condition missing after ErrorToCondition")
	}

	if _, ok := conditions[ConditionResolvedRefs]; !ok {
		t.Fatal("ResolvedRefs condition missing after ErrorToCondition")
	}
}
