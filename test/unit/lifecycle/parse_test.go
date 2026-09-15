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

package lifecycle_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle/ref"
)

var _ = Describe("ParseRefTag", func() {
	DescribeTable("parses valid tags",
		func(tag, field string, want ref.TagSpec) {
			got, err := ref.ParseRefTag(tag, field)
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(want))
		},
		Entry("kind only → fieldNameRef, no key",
			"secret", "Content",
			ref.TagSpec{Kind: "secret", RefField: "ContentRef", Key: ""},
		),
		Entry("kind + explicit ref field",
			"secret,SecretRef", "Content",
			ref.TagSpec{Kind: "secret", RefField: "SecretRef", Key: ""},
		),
		Entry("kind + ref field + key",
			"secret,SecretRef,tls.crt", "Content",
			ref.TagSpec{Kind: "secret", RefField: "SecretRef", Key: "tls.crt"},
		),
		Entry("trims spaces",
			" secret , SecretRef , tls.crt ", "Content",
			ref.TagSpec{Kind: "secret", RefField: "SecretRef", Key: "tls.crt"},
		),
		Entry("amcontext convention",
			"amcontext", "Context",
			ref.TagSpec{Kind: "amcontext", RefField: "ContextRef", Key: ""},
		),
		Entry("empty ref field defaults to fieldNameRef",
			"secret,,tls.crt", "Content",
			ref.TagSpec{Kind: "secret", RefField: "ContentRef", Key: "tls.crt"},
		),
		Entry("empty ref field, two parts, defaults to fieldNameRef",
			"secret,", "Content",
			ref.TagSpec{Kind: "secret", RefField: "ContentRef", Key: ""},
		),
	)

	DescribeTable("rejects invalid tags",
		func(tag string) {
			_, err := ref.ParseRefTag(tag, "Content")
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, ref.ErrRefTagInvalid)).To(BeTrue())
			Expect(errors.Is(err, &ref.Error{Op: "parse"})).To(BeTrue())
		},
		Entry("four parts", "secret,SecretRef,tls.crt,extra"),
		Entry("empty list after extra commas", "a,b,c,d"),
		Entry("empty tag", ""),
		Entry("whitespace only", "   "),
	)
})

var _ = Describe("Error", func() {
	It("Unwraps the sentinel and matches Is on Op", func() {
		err := &ref.Error{Op: "resolve", Kind: "secret", Err: ref.ErrUnknownKind}
		Expect(errors.Is(err, ref.ErrUnknownKind)).To(BeTrue())
		Expect(errors.Is(err, &ref.Error{Op: "resolve"})).To(BeTrue())
		Expect(errors.Is(err, &ref.Error{Op: "parse"})).To(BeFalse())
		Expect(errors.Is(err, &ref.Error{Kind: "secret"})).To(BeTrue())
	})
})
