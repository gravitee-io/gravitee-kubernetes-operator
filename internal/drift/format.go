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
	"fmt"
	"strings"
)

const indentSpaces = 2

func format(this *Result, b *strings.Builder, indent int) {
	if len(this.children) > 0 {
		if this.DriftDetected() {
			indent += indentSpaces
		}
		formatChildren(this, b, indent)
	} else if len(this.children) == 0 && this.Equivalent == Inequivalent && this.Property != "" {
		formatValue(this, b, indent)
	}
}

func formatValue(this *Result, b *strings.Builder, indent int) {
	addIndent(b, indent)
	if this.Index != nil {
		b.WriteString(fmt.Sprintf("%s[%d]: %v != %v", this.Property, *this.Index, resolve(this.CRDValue), resolve(this.RemoteValue)))
	} else {
		b.WriteString(fmt.Sprintf("%s: %v != %v", this.Property, resolve(this.CRDValue), resolve(this.RemoteValue)))
	}
	switch r := this.Reason.(type) {
	case string:
		if r != "" {
			b.WriteString(fmt.Sprintf(" (%s)", this.Reason))
		}
	case error:
		b.WriteString(fmt.Sprintf(" (error: %s)", this.Reason))
	default:
	}
	b.WriteString("\n")
}

func formatChildren(this *Result, b *strings.Builder, indent int) {
	for _, child := range this.children {
		if child.DriftDetected() && len(child.children) > 0 {
			addIndent(b, indent)
			property := child.Property
			if child.Index != nil {
				property += fmt.Sprintf("[%d]", *child.Index)
			}
			b.WriteString(fmt.Sprintf("%s:\n", property))
		}
		format(child, b, indent)
	}
}

func resolve(v any) any {
	if s, ok := v.(string); ok {
		return fmt.Sprintf(`"%s"`, s)
	}
	return v
}

func addIndent(b *strings.Builder, amount int) {
	if amount > 0 {
		b.WriteString(strings.Repeat(" ", amount))
	}
}
