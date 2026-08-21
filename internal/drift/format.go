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
	"unicode/utf8"
)

const (
	indentSpaces    = 2
	maxDisplayRunes = 60
	ellipsis        = "..."
	propSuffix      = ": "
	notEqOp         = " != "
)

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
	if left, right, ok := stringPair(this.CRDValue, this.RemoteValue); ok && isMultiline(left, right) {
		formatMultilineStrings(b, indent, propertyLabel(this), left, right)
		writeReason(this, b)
		b.WriteString("\n")
		return
	}
	addIndent(b, indent)
	b.WriteString(fmt.Sprintf("%s%s%v%s%v", propertyLabel(this), propSuffix, resolve(this.CRDValue), notEqOp, resolve(this.RemoteValue)))
	writeReason(this, b)
	b.WriteString("\n")
}

func propertyLabel(this *Result) string {
	if this.Index != nil {
		return fmt.Sprintf("%s[%d]", this.Property, *this.Index)
	}
	return this.Property
}

func stringPair(crd, remote any) (string, string, bool) {
	left, lok := crd.(string)
	right, rok := remote.(string)
	return left, right, lok && rok
}

func isMultiline(left, right string) bool {
	return strings.Contains(left, "\n") || strings.Contains(right, "\n")
}

func writeReason(this *Result, b *strings.Builder) {
	switch r := this.Reason.(type) {
	case string:
		if r != "" {
			b.WriteString(fmt.Sprintf(" (%s)", this.Reason))
		}
	case error:
		b.WriteString(fmt.Sprintf(" (error: %s)", this.Reason))
	default:
	}
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
		return displayQuoted(s)
	}
	return v
}

func displayQuoted(s string) string {
	visible, hidden := clipRunes(s)
	if hidden == 0 {
		return `"` + visible + `"`
	}
	return fmt.Sprintf(`"%s%s" (%s)`, visible, ellipsis, hiddenLabel(hidden))
}

func displayLine(s string) string {
	visible, hidden := clipRunes(s)
	if hidden == 0 {
		return visible
	}
	return fmt.Sprintf("%s%s (%s)", visible, ellipsis, hiddenLabel(hidden))
}

func clipRunes(s string) (string, int) {
	n := utf8.RuneCountInString(s)
	if n <= maxDisplayRunes {
		return s, 0
	}
	return takeRunes(s, maxDisplayRunes), n - maxDisplayRunes
}

func takeRunes(s string, n int) string {
	if n <= 0 {
		return ""
	}
	i := 0
	for idx := range s {
		if i == n {
			return s[:idx]
		}
		i++
	}
	return s
}

func hiddenLabel(n int) string {
	if n == 1 {
		return "1 char hidden"
	}
	return fmt.Sprintf("%d chars hidden", n)
}

func addIndent(b *strings.Builder, amount int) {
	if amount > 0 {
		b.WriteString(strings.Repeat(" ", amount))
	}
}
