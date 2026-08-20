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

package model

type Metadata interface {
	GetName() string
	GetValue() string
	GetDefaultValue() string
}

type BaseMetadata struct {
	Name         string  `json:"name"`
	Value        *string `json:"value,omitempty"`
	DefaultValue *string `json:"defaultValue,omitempty"`
}

func (m BaseMetadata) GetName() string {
	return m.Name
}

func (m BaseMetadata) GetValue() string {
	return *m.Value
}

func (m BaseMetadata) GetDefaultValue() string {
	return *m.DefaultValue
}
