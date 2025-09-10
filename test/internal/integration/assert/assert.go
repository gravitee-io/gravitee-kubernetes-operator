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
	"slices"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/sort"
	v1 "k8s.io/api/core/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	. "github.com/onsi/gomega"
)

const reconcileStatus = "reconcile status"
const reconcileCondition = "reconcile condition"

func HasFinalizer(object client.Object, value string) error {
	if !controllerutil.ContainsFinalizer(object, value) {
		return fmt.Errorf(
			"expected %s %s to have finalizer %s",
			object.GetObjectKind(), object.GetName(), value,
		)
	}
	return nil
}

func StrEndingWithPath(str, path string) error {
	return Equals("path", path, str[strings.LastIndex(str, "/"):])
}

func ApiCompleted(apiDefinition *v1alpha1.ApiDefinition) error {
	return Equals(reconcileStatus, core.ProcessingStatusCompleted, apiDefinition.Status.ProcessingStatus)
}

<<<<<<< HEAD
=======
func IsAccepted(obj core.ConditionAwareObject) error {
	if !k8s.IsAccepted(obj) {
		return newAssertEqualError("Accepted condition", "True", "False")
	}
	return nil
}

func IsResolved(obj core.ConditionAwareObject) error {
	if !k8s.IsResolved(obj) {
		return newAssertEqualError("ResolvedRefs condition", "True", "False")
	}
	return nil
}

func IsUnresolved(obj core.ConditionAwareObject) error {
	if k8s.IsResolved(obj) {
		return newAssertEqualError("ResolvedRefs condition", "False", "True")
	}
	return nil
}

func ApiAccepted(apiDefinition *v1alpha1.ApiDefinition) error {
	return Equals(reconcileCondition, true,
		k8s.MapConditions(apiDefinition.Status.Conditions)[k8s.ConditionAccepted].Status == metav1.ConditionTrue)
}

func ApiRejected(apiDefinition *v1alpha1.ApiDefinition) error {
	return Equals(reconcileCondition, true,
		k8s.MapConditions(apiDefinition.Status.Conditions)[k8s.ConditionAccepted].Status == metav1.ConditionFalse)
}

>>>>>>> 28d59ae (refactor: do not mutate notificaction spec on updates)
func ApiV4Completed(apiDefinition *v1alpha1.ApiV4Definition) error {
	return Equals(reconcileStatus, core.ProcessingStatusCompleted, apiDefinition.Status.ProcessingStatus)
}

func ApplicationCompleted(app *v1alpha1.Application) error {
	return Equals(reconcileStatus, core.ProcessingStatusCompleted, app.Status.ProcessingStatus)
}

func ApplicationFailed(app *v1alpha1.Application) error {
	return Equals(reconcileStatus, core.ProcessingStatusFailed, app.Status.ProcessingStatus)
}

func SubscriptionCompleted(sub *v1alpha1.Subscription) error {
	return Equals(reconcileStatus, core.ProcessingStatusCompleted, sub.Status.ProcessingStatus)
}

func SubscriptionFailed(sub *v1alpha1.Subscription) error {
	return Equals(reconcileStatus, core.ProcessingStatusFailed, sub.Status.ProcessingStatus)
}

func SharedPolicyGroupCompleted(sub *v1alpha1.SharedPolicyGroup) error {
	return Equals(reconcileStatus, core.ProcessingStatusCompleted, sub.Status.ProcessingStatus)
}

func SharedPolicyGroupFailed(sub *v1alpha1.SharedPolicyGroup) error {
	return Equals(reconcileStatus, core.ProcessingStatusFailed, sub.Status.ProcessingStatus)
}

func GroupCompleted(group *v1alpha1.Group) error {
	return Equals(reconcileStatus, core.ProcessingStatusCompleted, group.Status.ProcessingStatus)
}

func GroupFailed(group *v1alpha1.Group) error {
	return Equals(reconcileStatus, core.ProcessingStatusFailed, group.Status.ProcessingStatus)
}

func NotificationCompleted(notification *v1alpha1.Notification) error {
	return Equals(reconcileCondition, false, notification.Status.IsFailed())
}

func NotificationFailed(notification *v1alpha1.Notification) error {
	return Equals(reconcileCondition, true, notification.Status.IsFailed())
}

func ApiFailed(apiDefinition *v1alpha1.ApiDefinition) error {
	return Equals(reconcileStatus, core.ProcessingStatusFailed, apiDefinition.Status.ProcessingStatus)
}

func ApiV4Failed(apiDefinition *v1alpha1.ApiV4Definition) error {
	return Equals(reconcileStatus, core.ProcessingStatusFailed, apiDefinition.Status.ProcessingStatus)
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

func StrStartingWith(str, prefix string) error {
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

func Deleted[T client.Object](ctx context.Context, kind string, obj T) error {
	err := manager.GetLatest(ctx, obj)
	if k8serr.IsNotFound(err) {
		return nil
	}
	return newAssertEqualError(kind, "NOT FOUND", err)
}

func SliceEqualsSorted[S ~[]E, E any](field string, expected S, given S, comp sort.Comparator[E]) error {
	ecp, gcp := make([]E, len(expected)), make([]E, len(given))
	copy(ecp, expected)
	slices.SortFunc(ecp, comp)
	copy(gcp, given)
	slices.SortFunc(gcp, comp)
	return Equals(field, callStringerIfExists(ecp), callStringerIfExists(gcp))
}

func callStringerIfExists[E any](stringersOrNot []E) []any {
	stringified := make([]any, len(stringersOrNot))
	for i := range stringersOrNot {
		stringerOrNot := stringersOrNot[i]
		if stringer, ok := any(stringerOrNot).(fmt.Stringer); ok {
			stringified[i] = stringer.String()
		} else {
			stringified[i] = stringerOrNot
		}
	}
	return stringified
}

func NotEmptySlice[T any](field string, value []T) error {
	if len(value) == 0 {
		return fmt.Errorf("expected %#v not to be empty", field)
	}
	return nil
}

func SliceOfSize[T any](field string, value []T, size int) error {
	if len(value) != size {
		return fmt.Errorf("expected %s to have len %d, got %d", field, size, len(value))
	}
	return nil
}

func NotEmptyString(field string, value string) error {
	if value == "" {
		return fmt.Errorf("expected %s not to be empty", field)
	}
	return nil
}

func Nil(field string, value any) error {
	if value != nil && !reflect.ValueOf(value).IsNil() {
		return fmt.Errorf("expected %s to be nil", field)
	}
	return nil
}

func NotNil(field string, value any) error {
	if value == nil || reflect.ValueOf(value).IsNil() {
		return fmt.Errorf("expected %s not to be nil", field)
	}
	return nil
}

func MapContaining[K comparable, V any](m map[K]V, key K, value V) error {
	val, ok := m[key]
	if !ok {
		return fmt.Errorf("expected map to contain key %v", key)
	}
	return Equals(fmt.Sprintf("map[%v]", key), value, val)
}

func MapNotContaining[K comparable, V any](m map[K]V, key K) error {
	if _, ok := m[key]; ok {
		return fmt.Errorf("expected map to not contain key %v", key)
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
	return fmt.Errorf("expected %s to be %#v, got %#v", field, expected, given)
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
