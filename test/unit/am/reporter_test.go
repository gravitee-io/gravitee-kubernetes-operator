// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package am_test

import (
	"context"
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	amsdk "github.com/gravitee-io/gravitee-automation-sdk/am-sdk/v2/pkg/sdk"
	reportermodel "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/am/reporter"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/reporter"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// reporterWithReadyDomain registers a fake cluster holding a domain AM has created, and returns a file
// reporter under it.
func reporterWithReadyDomain() *v1alpha1.AMReporter {
	d := &v1alpha1.AMSecurityDomain{ObjectMeta: metav1.ObjectMeta{Name: "domain", Namespace: "ns"}}
	d.Status.Key = "ns-domain"
	scheme := runtime.NewScheme()
	Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
	k8s.RegisterClient(fake.NewClientBuilder().WithScheme(scheme).WithObjects(d).Build())

	r := &v1alpha1.AMReporter{ObjectMeta: metav1.ObjectMeta{Name: "audit", Namespace: "ns"}}
	r.Spec.DomainRef = refs.NamespacedName{Name: "domain"}
	r.Spec.Name = new("Audit to file")
	r.Spec.Type = new("reporter-am-file")
	r.Spec.Configuration = utils.NewGenericStringMap().
		Put("filename", "audit").
		Put("retainDays", 7)
	r.Spec.AttributeMappings = []reportermodel.ReporterAttributeMapping{
		{Expression: new("{#context.attributes['client']}"), ExportedName: new("client")},
	}
	r.Spec.AttributeMappingEventTypes = []string{"USER_LOGIN"}
	return r
}

// toRemote is what AM returns for a reporter it stored as sent.
func toRemote(obj *v1alpha1.AMReporter) reporter.Reporter {
	dto, err := reporter.ToReporterDTO(obj)
	Expect(err).ToNot(HaveOccurred())
	return dto
}

var _ = Describe("AMReporter DTO mapping", func() {
	It("keys the reporter and its domain by their HRIDs", func() {
		r := reporterWithReadyDomain()
		r.Spec.DomainRef.Namespace = "ignored"

		dto, err := reporter.ToReporterDTO(r)
		Expect(err).ToNot(HaveOccurred())

		Expect(dto.Key).To(Equal("ns-audit"))
		Expect(dto.DomainKey).To(Equal("ns-domain"))
	})

	It("sends configuration as a JSON string, with the attribute mappings", func() {
		dto, err := reporter.ToReporterDTO(reporterWithReadyDomain())
		Expect(err).ToNot(HaveOccurred())

		payload, err := json.Marshal(dto.Reporter)
		Expect(err).ToNot(HaveOccurred())
		var sent map[string]any
		Expect(json.Unmarshal(payload, &sent)).To(Succeed())
		sentConfiguration, ok := sent["configuration"].(string)
		Expect(ok).To(BeTrue(), "configuration sent as %T", sent["configuration"])
		var configuration map[string]any
		Expect(json.Unmarshal([]byte(sentConfiguration), &configuration)).To(Succeed())
		Expect(configuration).To(Equal(map[string]any{"filename": "audit", "retainDays": float64(7)}))
		Expect(sent["name"]).To(Equal("Audit to file"))
		Expect(sent["type"]).To(Equal("reporter-am-file"))
		Expect(sent["attributeMappings"]).To(Equal([]any{
			map[string]any{"expression": "{#context.attributes['client']}", "exportedName": "client"},
		}))
		Expect(sent["attributeMappingEventTypes"]).To(Equal([]any{"USER_LOGIN"}))
	})

	It("defaults enabled to true and system to false", func() {
		dto, err := reporter.ToReporterDTO(reporterWithReadyDomain())
		Expect(err).ToNot(HaveOccurred())

		Expect(dto.Enabled).To(Equal(new(true)))
		Expect(dto.System).To(Equal(new(false)))
	})
})

var _ = Describe("AMReporter drift mapping", func() {
	It("blanks the five fields AM ignores for the system reporter", func() {
		r := reporterWithReadyDomain()
		r.Spec.System = new(true)

		dto, err := reporter.ToReporterDTOForDrift(r)
		Expect(err).ToNot(HaveOccurred())

		Expect(dto.Name).To(BeNil())
		Expect(dto.Type).To(BeNil())
		Expect(dto.Configuration).To(BeNil())
		Expect(dto.AttributeMappings).To(BeNil())
		Expect(dto.AttributeMappingEventTypes).To(BeNil())
		Expect(dto.Key).To(Equal("ns-audit"))
	})

	It("keeps every field of a regular reporter", func() {
		r := reporterWithReadyDomain()

		dto, err := reporter.ToReporterDTOForDrift(r)
		Expect(err).ToNot(HaveOccurred())

		Expect(dto).To(Equal(toRemote(r)))
	})
})

var _ = Describe("AMReporter drift detection", func() {
	// detect compares an unchanged CR with AM, as admission does on an update (see drift.Merge): a field the
	// CR leaves unset is not compared.
	detect := func(local *v1alpha1.AMReporter, remote reporter.Reporter) drift.Result {
		dto, err := reporter.ToReporterDTOForDrift(local)
		Expect(err).ToNot(HaveOccurred())
		result := drift.DetectWithNamespace(dto, remote, local.Namespace)
		return drift.Merge(result, result)
	}

	It("detects nothing when AM holds the reporter as sent", func() {
		r := reporterWithReadyDomain()

		result := detect(r, toRemote(r))

		Expect(result.DriftDetected()).To(BeFalse(), result.String())
	})

	It("detects a change in enabled", func() {
		r := reporterWithReadyDomain()
		remote := toRemote(r)
		remote.Enabled = new(false)

		result := detect(r, remote)

		Expect(result.DriftDetected()).To(BeTrue())
	})

	It("detects a change in the attribute mappings", func() {
		r := reporterWithReadyDomain()
		remote := toRemote(r)
		remote.AttributeMappings = []amsdk.ReporterAttributeMapping{
			{Expression: new("{#context.attributes['user']}"), ExportedName: new("user")},
		}

		result := detect(r, remote)

		Expect(result.DriftDetected()).To(BeTrue())
	})

	It("detects a change in a configuration field", func() {
		r := reporterWithReadyDomain()
		changed := r.DeepCopy()
		changed.Spec.Configuration.Put("retainDays", 30)

		result := detect(r, toRemote(changed))

		Expect(result.DriftDetected()).To(BeTrue())
	})

	It("ignores the name, type, configuration and attribute mappings AM builds for the system reporter", func() {
		r := reporterWithReadyDomain()
		r.Spec.System = new(true)
		remote := toRemote(r)
		remote.Name = new("Default reporter")
		remote.Type = new("mongodb")
		remote.Configuration = toRemote(reporterWithReadyDomain()).Configuration
		remote.AttributeMappings = nil
		remote.AttributeMappingEventTypes = nil

		result := detect(r, remote)

		Expect(result.DriftDetected()).To(BeFalse(), result.String())
	})
})

var _ = Describe("AMReporter PreCheck", func() {
	ctx := context.Background()

	It("admits a valid reporter without warnings", func() {
		errs := reporter.PreCheck(ctx, reporterWithReadyDomain())

		Expect(errs.IsSevere()).To(BeFalse())
		Expect(warnings(errs)).To(BeEmpty())
	})

	It("rejects a key AM would refuse", func() {
		r := reporterWithReadyDomain()
		r.Name = "reporter.with.dots"

		Expect(reporter.PreCheck(ctx, r).IsSevere()).To(BeTrue())
	})

	It("rejects a domainRef in another namespace", func() {
		r := reporterWithReadyDomain()
		r.Spec.DomainRef.Namespace = "other"

		Expect(reporter.PreCheck(ctx, r).IsSevere()).To(BeTrue())
	})

	It("warns when the domain is missing", func() {
		r := reporterWithReadyDomain()
		r.Spec.DomainRef.Name = "missing"

		errs := reporter.PreCheck(ctx, r)

		Expect(errs.IsSevere()).To(BeFalse())
		Expect(warnings(errs)).To(ContainElement(ContainSubstring("not found or not yet created in AM")))
	})

	It("warns that name, type, configuration and the attribute mappings are ignored for the system reporter", func() {
		r := reporterWithReadyDomain()
		r.Spec.System = new(true)

		errs := reporter.PreCheck(ctx, r)

		Expect(errs.IsSevere()).To(BeFalse())
		Expect(warnings(errs)).To(ConsistOf(
			"'configuration', 'name', 'type', 'attributeMappings', 'attributeMappingEventTypes' " +
				"will be ignored when 'system' is 'true'."))
	})
})

var _ = Describe("AMReporter UpdateStatus", func() {
	It("copies the keys and the name and type AM stored", func() {
		r := &v1alpha1.AMReporter{}
		resp := reporter.Response{
			DomainSubResourceResponse: am.DomainSubResourceResponse{
				DomainKey: "ns-domain",
				BaseResponse: am.BaseResponse{
					Key:    "ns-audit",
					OrgEnv: am.OrgEnv{OrgID: "DEFAULT", EnvID: "DEFAULT"},
				},
			},
			Name: "Default",
			Type: "mongodb",
		}

		Expect(reporter.UpdateStatus(context.Background(), r, resp)).To(Succeed())

		Expect(r.Status.Key).To(Equal("ns-audit"))
		Expect(r.Status.DomainKey).To(Equal("ns-domain"))
		Expect(r.Status.OrgID).To(Equal("DEFAULT"))
		Expect(r.Status.EnvID).To(Equal("DEFAULT"))
		Expect(r.Status.Name).To(Equal("Default"))
		Expect(r.Status.Type).To(Equal("mongodb"))
	})
})
