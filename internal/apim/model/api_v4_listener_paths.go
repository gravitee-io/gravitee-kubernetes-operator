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

// WithoutOverrideAccess removes overrideAccess from every listener path, so that drift
// detection ignores it. A multi-tenant APIM forces it to true on the first path, and it
// only changes the entrypoint URL the console displays.
// Drift path only: the payload sent to APIM is unchanged.
func (dto APIV4DTO) WithoutOverrideAccess() APIV4DTO {
	for _, listener := range dto.Listeners {
		if listener == nil || listener.GenericStringMap == nil {
			continue
		}
		for _, item := range listener.GetSlice("paths") {
			if path, ok := item.(map[string]any); ok {
				delete(path, "overrideAccess")
			}
		}
	}
	return dto
}
