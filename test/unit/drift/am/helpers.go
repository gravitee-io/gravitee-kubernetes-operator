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

package am

import (
	"fmt"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func expectNoDrift(r drift.Result) {
	GinkgoHelper()
	Expect(r.String()).To(BeEmpty())
	Expect(r.DriftDetected()).To(BeFalse())
}

func expectDrift(r drift.Result, expected string) {
	GinkgoHelper()
	Expect(r.String()).To(Equal(expected))
	Expect(r.DriftDetected()).To(BeTrue())
}

// expectedEquivalentNotHavingAnyZeroValue is a helper to check that no property is not tested.
// It requires that no value is zero (nil, empty string, empty slice, false booleans etc.).
func expectedEquivalentNotHavingAnyZeroValue(crd, remote any) {
	GinkgoHelper()
	detect := drift.DetectWithNamespace(crd, remote, "")
	expectNoDrift(detect)
	doAssertNoResultHasZeroOrNilValue(detect, []string{})
}

func doAssertNoResultHasZeroOrNilValue(r drift.Result, ancestors []string) {
	GinkgoHelper()
	if r.Equivalent != drift.CannotCompare && r.Property != "" && (r.Index == nil || len(r.Children()) == 0) {
		Expect(r.CRDValue).NotTo(BeZero(),
			"%s.%s is not tested",
			strings.Join(ancestors, "."),
			r.Property)
	}
	if r.Children() != nil {
		if r.Property != "" {
			var index string
			if r.Index != nil {
				index = fmt.Sprintf("[%v]", *r.Index)
			}
			ancestors = append(ancestors, r.Property+index)
		}
		copyOfAncestors := make([]string, len(ancestors))
		copy(copyOfAncestors, ancestors)
		for _, child := range r.Children() {
			doAssertNoResultHasZeroOrNilValue(*child, copyOfAncestors)
		}
	}
}
