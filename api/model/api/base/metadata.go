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
	// Metadata Key
	Key string `json:"key"`
	// Metadata Name
	Name string `json:"name"`
	// Metadata Format
	Format MetadataFormat `json:"format"`
	// Metadata Value
	Value string `json:"value,omitempty"`
	// Metadata Default value
	// +kubebuilder:validation:Optional
	DefaultValue *string `json:"defaultValue,omitempty"`
}
