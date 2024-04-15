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

package base

// +kubebuilder:validation:Enum=STRING;NUMERIC;BOOLEAN;DATE;MAIL;URL;
type MetadataFormat string

type MetadataEntry struct {
	// MetaData Key
	Key string `json:"key"`
	// MetaData Name
	Name string `json:"name"`
	// MetaData Format
	Format MetadataFormat `json:"format"`
	// MetaData Value
	Value string `json:"value,omitempty"`
	// MetaData Default value
	DefaultValue string `json:"defaultValue,omitempty"`
}
