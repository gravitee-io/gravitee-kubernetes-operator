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
	"net/http"

	xhttp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
)

const notFoundStatusText = "Service not found."

type ResponseTemplate struct {
	Content     string
	ContentType string
}

type Opts struct {
	Templates map[int]ResponseTemplate
}

func NewOpts() Opts {
	return Opts{
		Templates: map[int]ResponseTemplate{},
	}
}

func mergeOpts(opts Opts) Opts {
	baseOpts := baseOpts()

	for status, template := range opts.Templates {
		baseOpts.Templates[status] = template
	}

	return baseOpts
}

func baseOpts() Opts {
	return Opts{
		Templates: map[int]ResponseTemplate{
			http.StatusNotFound: {
				Content:     notFoundStatusText,
				ContentType: xhttp.ContentTypeTextPlain,
			},
		},
	}
}
