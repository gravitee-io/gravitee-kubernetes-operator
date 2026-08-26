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

package env

import "os"

// DriftPolicy controls how the operator reacts to drift detection outcomes.
type DriftPolicy string

const (
	// DriftPolicyDeny blocks the update with an error.
	DriftPolicyDeny DriftPolicy = "deny"
	// DriftPolicyWarn allows the update and surfaces a kubectl warning.
	DriftPolicyWarn DriftPolicy = "warn"
	// DriftPolicyAllow allows the update and only logs a warning in the controller logs.
	DriftPolicyAllow DriftPolicy = "allow"
)

func (p DriftPolicy) String() string {
	return string(p)
}

// DriftDetection holds global drift detection settings loaded from the environment.
type DriftDetection struct {
	Enabled         bool
	Policy          DriftPolicy
	OnRemoteMissing DriftPolicy
	OnFetchFailure  DriftPolicy
}

func parseDriftPolicy(key string, defaultValue DriftPolicy) DriftPolicy {
	switch p := DriftPolicy(os.Getenv(key)); p {
	case DriftPolicyDeny, DriftPolicyWarn, DriftPolicyAllow:
		return p
	default:
		return defaultValue
	}
}
