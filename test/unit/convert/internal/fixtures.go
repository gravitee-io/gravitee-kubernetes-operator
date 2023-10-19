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

package internal

import (
	"encoding/json"
	"os"

	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	. "github.com/onsi/gomega"
)

const dataBase = "internal/data"

func UnmarshalV4(path string) *v4.Api {
	api := &v4.Api{}
	data := ReadJSON(path + "/v4.json")
	err := json.Unmarshal(data, api)
	Expect(err).NotTo(HaveOccurred())
	return api
}

func UnmarshalV2(path string) *v2.Api {
	api := &v2.Api{}
	data := ReadJSON(path + "/v2.json")
	err := json.Unmarshal(data, api)
	Expect(err).NotTo(HaveOccurred())
	return api
}

func Marshal(api interface{}) []byte {
	data, err := json.Marshal(api)
	Expect(err).NotTo(HaveOccurred())
	return data
}

func ReadJSON(path string) []byte {
	data, err := os.ReadFile(dataBase + "/" + path)
	Expect(err).NotTo(HaveOccurred())
	return data
}
