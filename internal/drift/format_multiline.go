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

	"github.com/pmezard/go-difflib/difflib"
)

// notEqOp is the aligned operator between CRD (left) and remote (right) lines.
// Subsequent differing rows keep the same column width with spaces so the
// right-hand side stays aligned under the first difference.
const notEqOpML = " " + notEqOp + " "

type lineRow struct {
	identical int
	left      string
	right     string
}

func formatMultilineStrings(builder *strings.Builder, indent int, property, left, right string) {
	if a, b := splitLines(left), splitLines(right); len(a) == 0 && len(b) == 0 {
		return
	} else if rows := diffLineRows(a, b); len(rows) > 0 {
		writeLineRows(builder, indent, property, rows)
	}
	return
}

func writeLineRows(b *strings.Builder, indent int, property string, rows []lineRow) {
	leftWidth := 0
	for i, row := range rows {
		if row.identical > 0 {
			continue
		}
		rows[i].left = displayLine(row.left)
		rows[i].right = displayLine(row.right)
		leftWidth = max(leftWidth, utf8.RuneCountInString(rows[i].left))
	}

	label := property + propSuffix + " "
	padding := strings.Repeat(" ", utf8.RuneCountInString(label))
	wroteNotEqOp := false

	for i, row := range rows {
		if i > 0 {
			b.WriteByte('\n')
		}
		addIndent(b, indent)
		if i == 0 {
			b.WriteString(label)
		} else {
			b.WriteString(padding)
		}
		writeLineRow(b, row, leftWidth, &wroteNotEqOp)
	}
}

func writeLineRow(b *strings.Builder, row lineRow, leftWidth int, wroteNotEqOp *bool) {
	if row.identical > 0 {
		b.WriteString(identicalSummary(row.identical))
		return
	}
	b.WriteString(row.left)
	b.WriteString(strings.Repeat(" ", leftWidth-utf8.RuneCountInString(row.left)))
	if !*wroteNotEqOp {
		b.WriteString(notEqOpML)
		*wroteNotEqOp = true
		b.WriteString(row.right)
		return
	}
	b.WriteString(strings.Repeat(" ", len(notEqOpML)))
	b.WriteString(row.right)
}

func identicalSummary(n int) string {
	if n == 1 {
		return "... 1 identical line"
	}
	return fmt.Sprintf("... %d identical lines", n)
}

func splitLines(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.Split(s, "\n")
}

// diffLineRows turns two line slices into display rows. Matching is done by
// go-difflib (equal / replace / delete / insert hunks); this function only maps
// those hunks onto collapsed identical runs and side-by-side changed lines.
func diffLineRows(a, b []string) []lineRow {
	var rows []lineRow
	for _, op := range difflib.NewMatcher(a, b).GetOpCodes() {
		// if equal
		if op.Tag == 'e' {
			if n := op.I2 - op.I1; n > 0 {
				rows = append(rows, lineRow{identical: n})
			}
			continue
		}
		left, right := a[op.I1:op.I2], b[op.J1:op.J2]
		for i := range max(len(left), len(right)) {
			row := lineRow{}
			if i < len(left) {
				row.left = left[i]
			}
			if i < len(right) {
				row.right = right[i]
			}
			rows = append(rows, row)
		}
	}
	return rows
}
