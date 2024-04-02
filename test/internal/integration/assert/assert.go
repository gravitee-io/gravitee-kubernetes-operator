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

package assert

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	. "github.com/onsi/gomega"
)

func AssertFinalizer(object client.Object, value string) error {
	if !controllerutil.ContainsFinalizer(object, value) {
		return fmt.Errorf(
			"expected %s %s to have finalizer %s",
			object.GetObjectKind(), object.GetName(), value,
		)
	}
	return nil
}

func PathEquals(url, path string) error {
	return Equals("path", path, url[strings.LastIndex(url, "/"):])
}

func ApiCompleted(apiDefinition *v1alpha1.ApiDefinition) error {
	return Equals("reconcile status", v1alpha1.ProcessingStatusCompleted, apiDefinition.Status.Status)
}

func ApplicationCompleted(app *v1alpha1.Application) error {
	return Equals("reconcile status", v1alpha1.ProcessingStatusCompleted, app.Status.Status)
}

func ApiFailed(apiDefinition *v1alpha1.ApiDefinition) error {
	return Equals("reconcile status", v1alpha1.ProcessingStatusFailed, apiDefinition.Status.Status)
}

func NoErrorAndHTTPStatus(err error, res *http.Response, expectedStatus int) error {
	if err != nil {
		return err
	}
	if res.StatusCode != expectedStatus {
		return newAssertEqualError("status", expectedStatus, res.StatusCode)
	}
	return nil
}

func StrStartsWith(str, prefix string) error {
	if !strings.HasPrefix(str, prefix) {
		return fmt.Errorf("expected %s to start with %s", str, prefix)
	}
	return nil
}

func Equals(field string, expected, given any) error {
	if !reflect.DeepEqual(expected, given) {
		return newAssertEqualError(field, expected, given)
	}
	return nil
}

func NotEmptySlice[T any](field string, value []T) error {
	if len(value) == 0 {
		return fmt.Errorf("expected %s not to be empty", field)
	}
	return nil
}

func NotEmptyString(field string, value string) error {
	if value == "" {
		return fmt.Errorf("expected %s not to be empty", field)
	}
	return nil
}

func EventsEmitted(obj client.Object, reasons ...string) {
	Eventually(
		getObjectEvents(obj),
		constants.EventualTimeout, constants.Interval,
	).Should(
		ContainElements(reasons),
	)
}

func NotFoundError(err error) error {
	if !errors.IsNotFound(err) {
		return newAssertEqualError("error", errors.NewNotFoundError(), err)
	}
	return nil
}

func newAssertEqualError(field string, expected, given any) error {
	return fmt.Errorf("expected %s to be %v, got %v", field, expected, given)
}

func getObjectEvents(obj client.Object) func() []string {
	return func() []string {
		eventsReason := []string{}

		events := &v1.EventList{}
		k8sClient := manager.Client()
		ctx := context.Background()
		if err := k8sClient.List(
			ctx,
			events,
			&client.ListOptions{Namespace: obj.GetNamespace()},
			client.MatchingFields{"involvedObject.name": obj.GetName()},
		); err != nil {
			return nil
		}

		for _, event := range events.Items {
			eventsReason = append(eventsReason, event.Reason)
		}
		return eventsReason
	}
}