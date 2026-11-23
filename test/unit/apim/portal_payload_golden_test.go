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

package apim_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/yaml"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	apimportal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portal"
	apimportallink "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portallink"
	apimportaltheme "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portaltheme"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/service"
)

// The portal navigation payloads GKO sent before the APIM SDK, one golden file per example manifest.
// They were written once from the hand-written mappers (UPDATE_GOLDEN=1) and are never regenerated: the
// SDK mappers must send the same JSON, except for the differences each kind declares below.

const (
	portalExamplesDir   = "../../../examples/apim/portal"
	portalGoldenDir     = "testdata/portal_payloads"
	portalGoldenNs      = "default"
	emptyContentExample = "synthetic-documentation-empty-content"
)

type portalPayloadExample struct {
	name string
	kind string
	obj  any
}

func loadPortalPayloadExamples() []portalPayloadExample {
	files, err := filepath.Glob(filepath.Join(portalExamplesDir, "*.yml"))
	Expect(err).ToNot(HaveOccurred())
	Expect(files).ToNot(BeEmpty())

	examples := make([]portalPayloadExample, 0, len(files)+1)
	for _, file := range files {
		data, err := os.ReadFile(file)
		Expect(err).ToNot(HaveOccurred())
		var meta struct {
			Kind string `json:"kind"`
		}
		Expect(yaml.Unmarshal(data, &meta)).To(Succeed())
		obj := newPortalPayloadObject(meta.Kind)
		Expect(yaml.Unmarshal(data, obj)).To(Succeed(), file)
		setPortalGoldenNamespace(obj)
		name := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		examples = append(examples, portalPayloadExample{name: name, kind: meta.Kind, obj: obj})
	}

	// No example has an empty content, which the SDK sends where GKO used to omit it.
	empty := &v1alpha1.Documentation{}
	empty.Name = "empty-content"
	empty.Namespace = portalGoldenNs
	empty.Spec.Name = "Empty"
	empty.Spec.PageType = "GRAVITEE_MARKDOWN"
	empty.Spec.Portal = &refs.NamespacedName{Name: "default-portal"}
	examples = append(examples, portalPayloadExample{name: emptyContentExample, kind: "Documentation", obj: empty})
	return examples
}

func newPortalPayloadObject(kind string) any {
	switch kind {
	case "PortalTheme":
		return &v1alpha1.PortalTheme{}
	case "Portal":
		return &v1alpha1.Portal{}
	case "PortalLink":
		return &v1alpha1.PortalLink{}
	case "PortalListing":
		return &v1alpha1.PortalListing{}
	case "Documentation":
		return &v1alpha1.Documentation{}
	}
	Fail("unexpected kind in examples/apim/portal: " + kind)
	return nil
}

func setPortalGoldenNamespace(obj any) {
	type namespaced interface {
		GetNamespace() string
		SetNamespace(string)
	}
	o, ok := obj.(namespaced)
	Expect(ok).To(BeTrue())
	if o.GetNamespace() == "" {
		o.SetNamespace(portalGoldenNs)
	}
}

// legacyPortalPayload is the body the hand-written services PUT.
func legacyPortalPayload(obj any) any {
	switch o := obj.(type) {
	case *v1alpha1.PortalTheme:
		return model.ToPortalThemeDTO(o)
	case *v1alpha1.Portal:
		return model.ToPortalDTO(o.Spec.Type, refs.NewNamespacedNameFromObject(o).HRID(), o.ActiveThemeHRID())
	case *v1alpha1.PortalLink:
		return model.ToPortalLinkDTO(o.Spec.Type, refs.NewNamespacedNameFromObject(o).HRID())
	case *v1alpha1.PortalListing:
		return service.ToPortalListingDTO(o)
	case *v1alpha1.Documentation:
		return model.ToDocumentationDTO(o)
	}
	Fail("unexpected object")
	return nil
}

// sdkPortalPayload is the body the SDK packages PUT. Kinds not migrated yet still go through the
// hand-written mapper.
func sdkPortalPayload(obj any) any {
	switch o := obj.(type) {
	case *v1alpha1.PortalTheme:
		dto, err := apimportaltheme.ToDTO(o)
		Expect(err).ToNot(HaveOccurred())
		return dto
	case *v1alpha1.Portal:
		dto, err := apimportal.ToDTO(o)
		Expect(err).ToNot(HaveOccurred())
		return dto
	case *v1alpha1.PortalLink:
		dto, err := apimportallink.ToDTO(o)
		Expect(err).ToNot(HaveOccurred())
		return dto
	}
	return legacyPortalPayload(obj)
}

// withAllowedDifferences applies to a golden payload the only differences the SDK payload may have:
// a portal no longer sends the deprecated navigation[].order, which APIM ignores (Issue 07).
func withAllowedDifferences(kind string, golden map[string]any) map[string]any {
	if kind == "Portal" {
		if entries, ok := golden["navigation"].([]any); ok {
			for _, entry := range entries {
				if fields, ok := entry.(map[string]any); ok {
					delete(fields, "order")
				}
			}
		}
	}
	return golden
}

func marshalPortalPayload(payload any) map[string]any {
	data, err := json.Marshal(payload)
	Expect(err).ToNot(HaveOccurred())
	out := map[string]any{}
	Expect(json.Unmarshal(data, &out)).To(Succeed())
	return out
}

func readPortalGolden(name string) map[string]any {
	data, err := os.ReadFile(filepath.Join(portalGoldenDir, name+".json"))
	Expect(err).ToNot(HaveOccurred(), "missing golden file, run once with UPDATE_GOLDEN=1")
	out := map[string]any{}
	Expect(json.Unmarshal(data, &out)).To(Succeed())
	return out
}

func writePortalGolden(name string, payload any) {
	data, err := json.MarshalIndent(payload, "", "  ")
	Expect(err).ToNot(HaveOccurred())
	Expect(os.MkdirAll(portalGoldenDir, 0o755)).To(Succeed())
	Expect(os.WriteFile(filepath.Join(portalGoldenDir, name+".json"), append(data, '\n'), 0o600)).To(Succeed())
}

var _ = Describe("Portal navigation payloads", func() {
	It("should match the golden payloads GKO sent before the APIM SDK", func() {
		for _, example := range loadPortalPayloadExamples() {
			if os.Getenv("UPDATE_GOLDEN") == "1" {
				writePortalGolden(example.name, legacyPortalPayload(example.obj))
				continue
			}
			golden := withAllowedDifferences(example.kind, readPortalGolden(example.name))
			Expect(marshalPortalPayload(sdkPortalPayload(example.obj))).To(Equal(golden), example.name)
		}
	})
})
