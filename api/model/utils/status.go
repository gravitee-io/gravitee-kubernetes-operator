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

package utils

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func MapConditions(conditionsSlice []metav1.Condition) map[string]metav1.Condition {
	conditions := make(map[string]metav1.Condition)
	for _, condition := range conditionsSlice {
		conditions[condition.Type] = condition
	}
	return conditions
}

func ToConditions(conditionsMap map[string]metav1.Condition) []metav1.Condition {
	conditions := make([]metav1.Condition, 0)
	for _, condition := range conditionsMap {
		conditions = append(conditions, condition)
	}
	return conditions
}
