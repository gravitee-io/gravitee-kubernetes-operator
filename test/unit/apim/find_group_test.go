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
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/service"
	xhttp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
)

const groupsPath = "/management/organizations/DEFAULT/environments/DEFAULT/configuration/groups/_paged"

// pagedGroups serves the given pages of groups the way the APIM groups _paged endpoint does,
// and records the requested page numbers.
func pagedGroups(pages [][]model.Group, requested *[]int) *httptest.Server {
	total := 0
	for _, p := range pages {
		total += len(p)
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer GinkgoRecover()
		Expect(r.URL.Path).To(Equal(groupsPath))
		// APIM defaults to the first page when none is requested
		page := 1
		if p := r.URL.Query().Get("page"); p != "" {
			var err error
			page, err = strconv.Atoi(p)
			Expect(err).ToNot(HaveOccurred())
		}
		*requested = append(*requested, page)

		data := []model.Group{}
		if page >= 1 && page <= len(pages) {
			data = pages[page-1]
		}
		w.Header().Set("Content-Type", "application/json")
		Expect(json.NewEncoder(w).Encode(map[string]any{
			"data": data,
			"page": map[string]any{
				"current":        page,
				"size":           len(data),
				"total_pages":    len(pages),
				"total_elements": total,
			},
		})).To(Succeed())
	}))
}

func newEnv(baseURL string) *service.Env {
	urls, err := client.NewURLs(baseURL, "/management", "/automation", "DEFAULT", "DEFAULT")
	Expect(err).ToNot(HaveOccurred())
	httpClient, err := xhttp.NewClient(context.Background(), nil)
	Expect(err).ToNot(HaveOccurred())
	return service.NewEnv(&client.Client{HTTP: httpClient, URLs: urls})
}

var _ = Describe("Env.FindGroup", func() {
	var requested []int

	BeforeEach(func() {
		requested = nil
	})

	It("returns the group named exactly as requested, not the first search result", func() {
		server := pagedGroups([][]model.Group{{
			{ID: "id-ops-admin", Name: "ops-admin"},
			{ID: "id-admin", Name: "admin"},
			{ID: "id-platform-admin", Name: "platform-admin"},
		}}, &requested)
		defer server.Close()

		group, err := newEnv(server.URL).FindGroup("admin")

		Expect(err).ToNot(HaveOccurred())
		Expect(group).ToNot(BeNil())
		Expect(group.ID).To(Equal("id-admin"))
		Expect(requested).To(Equal([]int{1}))
	})

	It("browses the next pages until the exact name is found", func() {
		server := pagedGroups([][]model.Group{
			{{ID: "id-admin-1", Name: "admin-1"}},
			{{ID: "id-admin", Name: "admin"}},
		}, &requested)
		defer server.Close()

		group, err := newEnv(server.URL).FindGroup("admin")

		Expect(err).ToNot(HaveOccurred())
		Expect(group).ToNot(BeNil())
		Expect(group.ID).To(Equal("id-admin"))
		Expect(requested).To(Equal([]int{1, 2}))
	})

	It("returns nil when only groups with a different name match the search", func() {
		server := pagedGroups([][]model.Group{
			{{ID: "id-platform-admin", Name: "platform-admin"}},
			{{ID: "id-Admin", Name: "Admin"}},
		}, &requested)
		defer server.Close()

		group, err := newEnv(server.URL).FindGroup("admin")

		Expect(err).ToNot(HaveOccurred())
		Expect(group).To(BeNil())
		Expect(requested).To(Equal([]int{1, 2}))
	})

	It("returns nil when the search has no result", func() {
		server := pagedGroups([][]model.Group{}, &requested)
		defer server.Close()

		group, err := newEnv(server.URL).FindGroup("admin")

		Expect(err).ToNot(HaveOccurred())
		Expect(group).To(BeNil())
		Expect(requested).To(Equal([]int{1}))
	})
})
