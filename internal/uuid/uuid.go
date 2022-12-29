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

package uuid

import (
	"encoding/base64"

	uuid "github.com/satori/go.uuid" //nolint:gomodguard // to replace with google implementation
)

func FromString(seed string) string {
	encoded := base64.RawStdEncoding.EncodeToString([]byte(seed))
	return uuid.NewV3(uuid.NamespaceURL, encoded).String()
}

func NewV4String() string {
	return uuid.NewV4().String()
}
