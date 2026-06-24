package drift

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func expectDrift(r drift.Result, expected string) {
	GinkgoHelper()
	Expect(r.DriftDetected()).To(BeTrue())
	Expect(r.String()).To(Equal(expected))
}

func expectNoDrift(r drift.Result) {
	GinkgoHelper()
	Expect(r.String()).To(BeEmpty())
	Expect(r.DriftDetected()).To(BeFalse())
}


func ptr[T any](v T) *T {
	GinkgoHelper()
	return &v
}
