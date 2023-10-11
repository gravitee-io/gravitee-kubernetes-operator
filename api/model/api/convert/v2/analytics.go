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

package v2

import (
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
)

func toLogging(analytics *v4.Analytics) *v2.Logging {
	if analytics == nil || analytics.Logging == nil {
		return nil
	}

	return &v2.Logging{
		Condition: analytics.Logging.Condition,
		Content:   toLoggingContent(analytics.Logging.Content),
		Mode:      toV2LoggingMode(analytics.Logging.Mode),
		Scope:     toV2LoggingScope(analytics.Logging.Phase),
	}
}

func toLoggingContent(content *v4.LoggingContent) v2.LoggingContent {
	if content == nil {
		return v2.NoLoggingContent
	}

	switch {
	case content.Headers && !content.Payload:
		return v2.HeadersLoggingContent
	case !content.Headers && content.Payload:
		return v2.PayloadsLoggingContent
	case content.Headers && content.Payload:
		return v2.HeadersPayloadsLoggingContent
	default:
		return v2.NoLoggingContent
	}
}

func toV2LoggingMode(mode *v4.LoggingMode) v2.LoggingMode {
	if mode == nil {
		return v2.NoLoggingMode
	}

	switch {
	case mode.Entrypoint && !mode.Endpoint:
		return v2.ClientMode
	case !mode.Entrypoint && mode.Endpoint:
		return v2.ProxyMode
	case mode.Entrypoint && mode.Endpoint:
		return v2.ClientProxyMode
	default:
		return v2.NoLoggingMode
	}
}

func toV2LoggingScope(phase *v4.LoggingPhase) v2.LoggingScope {
	if phase == nil {
		return v2.NoLoggingScope
	}

	switch {
	case phase.Request && !phase.Response:
		return v2.RequestLoggingScope
	case !phase.Request && phase.Response:
		return v2.ResponseLoggingScope
	case phase.Request && phase.Response:
		return v2.RequestResponseLoggingScope
	default:
		return v2.NoLoggingScope
	}
}
