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

package v4

import (
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
)

func toAnalytics(proxy *v2.Proxy) *v4.Analytics {
	if proxy == nil || proxy.Logging == nil {
		return nil
	}
	analytics := v4.NewAnalytics()
	analytics.Logging = toLogging(proxy.Logging)
	return analytics
}

func toLogging(logging *v2.Logging) *v4.Logging {
	return &v4.Logging{
		Condition: logging.Condition,
		Content:   toLoggingContent(logging.Content),
		Mode:      toLoggingMode(logging.Mode),
		Phase:     toLoggingPhase(logging.Scope),
	}
}

func toLoggingContent(content v2.LoggingContent) *v4.LoggingContent {
	switch content {
	case v2.HeadersLoggingContent:
		return v4.NewLoggingContent(true, false, false, false, false)
	case v2.PayloadsLoggingContent:
		return v4.NewLoggingContent(false, false, true, false, false)
	case v2.HeadersPayloadsLoggingContent:
		return v4.NewLoggingContent(true, false, true, false, false)
	default:
		return nil
	}
}

func toLoggingMode(mode v2.LoggingMode) *v4.LoggingMode {
	switch mode {
	case v2.ClientMode:
		return v4.NewLoggingMode(true, false)
	case v2.ProxyMode:
		return v4.NewLoggingMode(false, true)
	case v2.ClientProxyMode:
		return v4.NewLoggingMode(true, true)
	default:
		return nil
	}
}

func toLoggingPhase(phase v2.LoggingScope) *v4.LoggingPhase {
	switch phase {
	case v2.RequestLoggingScope:
		return v4.NewLoggingPhase(true, false)
	case v2.ResponseLoggingScope:
		return v4.NewLoggingPhase(false, true)
	case v2.RequestResponseLoggingScope:
		return v4.NewLoggingPhase(true, true)
	default:
		return nil
	}
}
