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

package mapper

import (
	"encoding/json"
	"fmt"
)

// MapViaJSON converts src to T through its JSON form. An error means some field did not fit: the partial
// result must not be used.
func MapViaJSON[T any](src any) (T, error) {
	var dst T
	if src == nil {
		return dst, nil
	}

	data, err := json.Marshal(src)
	if err != nil {
		return dst, fmt.Errorf("mapping %T to %T: %w", src, dst, err)
	}

	if err := json.Unmarshal(data, &dst); err != nil {
		return dst, fmt.Errorf("mapping %T to %T: %w", src, dst, err)
	}

	return dst, nil
}
