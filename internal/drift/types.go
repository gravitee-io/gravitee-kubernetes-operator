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

package drift

import (
	"strings"
)

type EquivalentStatus byte

const (
	CannotCompare EquivalentStatus = iota
	Equivalent    EquivalentStatus = 1
	Inequivalent  EquivalentStatus = 2
)

type Equivalence struct {
	Equivalent EquivalentStatus
	FuncName   string
	Skip       bool
}

type EquivalenceFunc func(crd any, api any) Equivalence

type Result struct {
	Equivalence
	Property string
	Index    *int
	CRDValue any
	APIValue any
	Children []*Result
}

func (r *Result) String() string {
	var builder strings.Builder
	format(r, &builder, -2)
	return strings.TrimSpace(strings.TrimRight(builder.String(), "\n"))
}

func (r *Result) DriftDetected() bool {
	if r.Equivalent == Inequivalent {
		return true
	}
	if len(r.Children) > 0 {
		for _, child := range r.Children {
			if child.DriftDetected() {
				return true
			}
		}
	}
	return false
}
