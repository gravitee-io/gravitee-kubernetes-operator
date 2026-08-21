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

package framework

import (
	"fmt"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
)

type withDescription struct {
	Description string `json:"description"`
}

var _ = Describe("multiline string formatting", func() {
	const (
		leftShort  = "same first line\nleft-two\nleft-three"
		rightShort = "same first line\nright-two\nright-three"

		leftFull  = leftShort + "\na\nb\nc\nd\ne\nleft-end"
		rightFull = rightShort + "\na\nb\nc\nd\ne\nright-end"
	)

	It("collapses identical lines and aligns differing lines", func() {
		expectDrift(drift.DetectWithNamespace(
			withDescription{Description: leftShort},
			withDescription{Description: rightShort},
			"",
		), `description:  ... 1 identical line
              left-two    !=  right-two
              left-three      right-three`)
	})

	It("collapses identical runs between later differences", func() {
		expectDrift(drift.DetectWithNamespace(
			withDescription{Description: leftFull},
			withDescription{Description: rightFull},
			"",
		), `description:  ... 1 identical line
              left-two    !=  right-two
              left-three      right-three
              ... 5 identical lines
              left-end        right-end`)
	})

	It("keeps tree indentation on continuation lines", func() {
		shared := Nested{Name: new("n"), Multiple: MultipleWithPtr{Name: new("m"), Order: new(1)}}
		expectDrift(drift.DetectWithNamespace(
			nestedDeep{Title: "t", First: shared, Second: &Nested{
				Name: new("n"), Multiple: MultipleWithPtr{Name: new("m"), Order: new(1)}, Description: new(leftShort),
			}},
			nestedDeep{Title: "t", First: shared, Second: &Nested{
				Name: new("n"), Multiple: MultipleWithPtr{Name: new("m"), Order: new(1)}, Description: new(rightShort),
			}},
			"",
		), `second:
  description:  ... 1 identical line
                left-two    !=  right-two
                left-three      right-three`)
	})

	It("keeps the property index in the label", func() {
		expectDrift(drift.DetectWithNamespace(
			withStringArray{Name: "test", Tags: []string{leftShort}},
			withStringArray{Name: "test", Tags: []string{rightShort}},
			"",
		), `tags[0]:  ... 1 identical line
          left-two    !=  right-two
          left-three      right-three`)
	})

	It("pairs an inserted line instead of shifting the rest", func() {
		expectDrift(drift.DetectWithNamespace(
			withDescription{Description: "hello\nworld"},
			withDescription{Description: "hello\ninserted\nworld"},
			"",
		), `description:  ... 1 identical line
                !=  inserted
              ... 1 identical line`)
	})
})

var _ = Describe("long string formatting", func() {
	It("truncates a long single-line string on both sides", func() {
		left := strings.Repeat("A", 180)
		right := strings.Repeat("B", 130)
		expectDrift(drift.DetectWithNamespace(
			withDescription{Description: left},
			withDescription{Description: right},
			"",
		), fmt.Sprintf(
			`description: "%s..." (120 chars hidden) != "%s..." (70 chars hidden)`,
			strings.Repeat("A", 60),
			strings.Repeat("B", 60),
		))
	})

	It("truncates only the long side of a single-line string", func() {
		expectDrift(drift.DetectWithNamespace(
			withDescription{Description: strings.Repeat("A", 180)},
			withDescription{Description: "short"},
			"",
		), fmt.Sprintf(
			`description: "%s..." (120 chars hidden) != "short"`,
			strings.Repeat("A", 60),
		))
	})

	It("does not truncate a string of exactly 60 characters", func() {
		expectDrift(drift.DetectWithNamespace(
			withDescription{Description: strings.Repeat("A", 60)},
			withDescription{Description: strings.Repeat("B", 60)},
			"",
		), fmt.Sprintf(
			`description: "%s" != "%s"`,
			strings.Repeat("A", 60),
			strings.Repeat("B", 60),
		))
	})

	It("reports a single hidden character", func() {
		expectDrift(drift.DetectWithNamespace(
			withDescription{Description: strings.Repeat("A", 61)},
			withDescription{Description: "short"},
			"",
		), fmt.Sprintf(
			`description: "%s..." (1 char hidden) != "short"`,
			strings.Repeat("A", 60),
		))
	})

	It("truncates long lines in a multiline string and keeps columns aligned", func() {
		leftLine := strings.Repeat("A", 180)
		rightLine := strings.Repeat("B", 130)
		expectDrift(drift.DetectWithNamespace(
			withDescription{Description: "same\n" + leftLine + "\nshort-left"},
			withDescription{Description: "same\n" + rightLine + "\nshort-right"},
			"",
		), fmt.Sprintf(
			`description:  ... 1 identical line
              %s... (120 chars hidden)  !=  %s... (70 chars hidden)
              short-left%s%s%s`,
			strings.Repeat("A", 60),
			strings.Repeat("B", 60),
			strings.Repeat(" ", 72),
			strings.Repeat(" ", 6),
			"short-right",
		))
	})
})
